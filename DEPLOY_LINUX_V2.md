# RunBa 部署方案 V2（Mac 开发机 -> Linux 生产机）

适用项目：
- `runba-backend`
- `runba-frontend`

目标：你在 MacBook 开发，远程 Ubuntu/CentOS Linux 服务器运行。

部署原则（关键）：
1. 只把源码同步到 Linux，不要把 Mac 构建产物直接上传运行。
2. 后端二进制在 Linux 上编译。
3. 前端 `node_modules` 和 `.next` 在 Linux 上重新安装/构建。

## 1. Linux 服务器准备（Ubuntu 示例）

```bash
sudo apt update
sudo apt install -y nginx mysql-server redis-server curl git
```

如果你的服务器是 CentOS（或 Rocky/Alma），使用：

```bash
sudo dnf -y install epel-release
sudo dnf -y install nginx mariadb-server redis curl git
sudo systemctl enable --now nginx mariadb redis
```

如果你是 CentOS 7（没有 `dnf`），使用：

```bash
sudo yum -y install epel-release
sudo yum -y install nginx mariadb-server redis curl git
sudo systemctl enable --now nginx mariadb redis
```

安装 Node.js 20：

```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs
node -v
npm -v
```

CentOS（或 Rocky/Alma）安装 Node.js 20：

```bash
curl -fsSL https://rpm.nodesource.com/setup_20.x | sudo bash -
sudo dnf -y install nodejs
node -v
npm -v
```

CentOS 7（`yum`）安装 Node.js 20：

```bash
curl -fsSL https://rpm.nodesource.com/setup_20.x | sudo bash -
sudo yum -y install nodejs
node -v
npm -v
```

安装 Go（建议官方包，保证 >= 1.24，按服务器架构选择）：

```bash
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then GO_PKG="go1.24.5.linux-amd64.tar.gz"; fi
if [ "$ARCH" = "aarch64" ]; then GO_PKG="go1.24.5.linux-arm64.tar.gz"; fi
if [ -z "$GO_PKG" ]; then echo "Unsupported arch: $ARCH" && exit 1; fi

cd /tmp
curl -LO "https://go.dev/dl/${GO_PKG}"
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf "${GO_PKG}"
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee /etc/profile.d/go.sh
source /etc/profile.d/go.sh
go version
```

## 2. 目录与账号

```bash
sudo useradd -r -s /usr/sbin/nologin runba || true
sudo mkdir -p /opt/runba
sudo chown -R $USER:$USER /opt/runba
```

建议目录：
- `/opt/runba/runba-backend`
- `/opt/runba/runba-frontend`

## 3. 从 Mac 同步到 Linux（推荐两种）

### 方式 A：Git（推荐）
在 Mac 提交并推送，然后在 Linux 拉取：

```bash
cd /opt/runba
git clone <your-repo-url> .
# 或已有仓库时
# git pull
```

### 方式 B：rsync 直传源码
在 Mac 执行：

```bash
rsync -avz --delete \
  --exclude '.git' \
  --exclude 'runba-frontend/node_modules' \
  --exclude 'runba-frontend/.next' \
  --exclude 'runba-backend/runba-backend' \
  /Users/jio/codespace/go/src/RunBa/ \
  <user>@<server_ip>:/opt/runba/
```

## 4. 后端配置与数据库

编辑 `/opt/runba/runba-backend/config.yaml`：
- `server.mode: release`
- `database.*` 改为生产库配置
- `jwt.secret` 改为强随机字符串（建议 32 字符以上）
- 不要使用开发环境密码

创建数据库：

```sql
CREATE DATABASE rb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

导入表：

```bash
mysql -u root -p rb < /opt/runba/runba-backend/forges_swordsmiths.sql
mysql -u root -p rb < /opt/runba/runba-backend/certificates.sql
```

补充 `users` 表（仓库当前无自动迁移）：

```sql
CREATE TABLE users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(50) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  nickname VARCHAR(50) DEFAULT '',
  status TINYINT NOT NULL DEFAULT 1,
  last_login DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  KEY idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

初始化管理员（上线后立即改密码）：

```sql
INSERT INTO users(username,password,nickname,status,created_at,updated_at)
VALUES('admin','$2a$10$Q7zEbOkd728rwWywoJzxuOTaPoWrhlWdnPP8gemPxNFzNLYIsZ.1m','admin',1,NOW(),NOW());
```

## 5. 后端（Linux 上编译 + systemd）

```bash
cd /opt/runba/runba-backend
/usr/local/go/bin/go mod download
/usr/local/go/bin/go build -o /opt/runba/runba-backend/runba-backend main.go
```

创建 `/etc/systemd/system/runba-backend.service`：

```ini
[Unit]
Description=RunBa Backend
After=network.target mysql.service redis-server.service

[Service]
Type=simple
WorkingDirectory=/opt/runba/runba-backend
ExecStart=/opt/runba/runba-backend/runba-backend
Restart=always
RestartSec=3
User=runba
Group=runba

[Install]
WantedBy=multi-user.target
```

授权并启动：

```bash
sudo chown -R runba:runba /opt/runba/runba-backend
sudo systemctl daemon-reload
sudo systemctl enable --now runba-backend
```

## 6. 前端（Linux 上安装依赖 + 构建 + systemd）

```bash
cd /opt/runba/runba-frontend
npm ci
npm run build
```

创建 `/etc/systemd/system/runba-frontend.service`：

```ini
[Unit]
Description=RunBa Frontend (Next.js)
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/runba/runba-frontend
ExecStart=/usr/bin/npm run start -- -p 3000 -H 127.0.0.1
Restart=always
RestartSec=3
User=runba
Group=runba
Environment=NODE_ENV=production

[Install]
WantedBy=multi-user.target
```

授权并启动：

```bash
sudo chown -R runba:runba /opt/runba/runba-frontend
sudo systemctl daemon-reload
sudo systemctl enable --now runba-frontend
```

## 7. Nginx 反向代理

创建 `/etc/nginx/sites-available/runba.conf`：

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

启用并重载：

```bash
sudo ln -sf /etc/nginx/sites-available/runba.conf /etc/nginx/sites-enabled/runba.conf
sudo nginx -t
sudo systemctl reload nginx
```

## 8. HTTPS（Let's Encrypt）

```bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
sudo systemctl status certbot.timer
```

## 9. 验证与排障

```bash
curl http://127.0.0.1:8080/health
sudo systemctl status runba-backend
sudo systemctl status runba-frontend
sudo journalctl -u runba-backend -n 100 --no-pager
sudo journalctl -u runba-frontend -n 100 --no-pager
```

## 10. 日常发布流程（Mac -> Linux）

每次发布都按这个顺序：

```bash
# Linux 服务器
cd /opt/runba
# git pull 或 rsync 同步后执行

cd /opt/runba/runba-backend
/usr/local/go/bin/go build -o /opt/runba/runba-backend/runba-backend main.go
sudo systemctl restart runba-backend

cd /opt/runba/runba-frontend
npm ci
npm run build
sudo systemctl restart runba-frontend
```

## 11. Mac 与 Linux 跨平台注意事项

1. 不要把 Mac 编译的 `runba-backend` 直接传到 Linux 运行。
2. 不要把 Mac 的 `node_modules`、`.next` 上传到 Linux 直接跑。
3. 如果你必须在 Mac 交叉编译后端，至少要显式指定：

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o runba-backend-linux main.go
```

4. 但前端仍建议在 Linux 上执行 `npm ci && npm run build`。
