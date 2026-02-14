# RunBa Linux 部署方案

适用项目：
- `runba-backend`
- `runba-frontend`

部署方式：同机部署 + `systemd` + `nginx` 反向代理。

## 1. 服务器准备（Ubuntu 示例）

```bash
sudo apt update
sudo apt install -y nginx mysql-server curl git
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs
```

## 2. 部署目录

```bash
sudo mkdir -p /opt/runba
sudo chown -R $USER:$USER /opt/runba
cd /opt/runba
```

把两个项目放到：
- `/opt/runba/runba-backend`
- `/opt/runba/runba-frontend`

## 3. 后端配置与数据库

先修改 `runba-backend/config.yaml`：
- `server.mode: release`
- `database.*` 改为服务器真实配置
- `jwt.secret` 改为强随机字符串

创建数据库：

```sql
CREATE DATABASE rb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

导入已有表结构：

```bash
mysql -u root -p rb < /opt/runba/runba-backend/forges_swordsmiths.sql
mysql -u root -p rb < /opt/runba/runba-backend/certificates.sql
```

补充 `users` 表（当前仓库无自动迁移）：

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

插入管理员（账号 `admin`，密码 `ChangeMe123!`，上线后请立刻修改密码）：

```sql
INSERT INTO users(username,password,nickname,status,created_at,updated_at)
VALUES('admin','$2a$10$Q7zEbOkd728rwWywoJzxuOTaPoWrhlWdnPP8gemPxNFzNLYIsZ.1m','admin',1,NOW(),NOW());
```

## 4. 后端打包并做服务

```bash
cd /opt/runba/runba-backend
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o runba-backend main.go
```

创建 `systemd` 服务文件：`/etc/systemd/system/runba-backend.service`

```ini
[Unit]
Description=RunBa Backend
After=network.target

[Service]
WorkingDirectory=/opt/runba/runba-backend
ExecStart=/opt/runba/runba-backend/runba-backend
Restart=always
RestartSec=3
User=www-data
Group=www-data

[Install]
WantedBy=multi-user.target
```

启动后端：

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now runba-backend
```

## 5. 前端构建并做服务

```bash
cd /opt/runba/runba-frontend
npm ci
npm run build
```

创建 `systemd` 服务文件：`/etc/systemd/system/runba-frontend.service`

```ini
[Unit]
Description=RunBa Frontend
After=network.target

[Service]
WorkingDirectory=/opt/runba/runba-frontend
ExecStart=/usr/bin/npm run start -- -p 3000 -H 127.0.0.1
Restart=always
RestartSec=3
User=www-data
Group=www-data
Environment=NODE_ENV=production

[Install]
WantedBy=multi-user.target
```

启动前端：

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now runba-frontend
```

## 6. Nginx 反向代理（同域名）

创建配置：`/etc/nginx/sites-available/runba.conf`

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
    }

    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

启用并加载：

```bash
sudo ln -s /etc/nginx/sites-available/runba.conf /etc/nginx/sites-enabled/runba.conf
sudo nginx -t
sudo systemctl reload nginx
```

## 7. 最后检查

```bash
curl http://127.0.0.1:8080/health
sudo systemctl status runba-backend
sudo systemctl status runba-frontend
```

## 8. 上线后建议

1. 使用 HTTPS（建议接入 Let's Encrypt）。
2. 及时修改默认管理员密码。
3. 把 `config.yaml` 中敏感信息（数据库密码、JWT 密钥）迁移为环境变量。
4. 给 MySQL 做定时备份。
