# RunBa Ubuntu 24 部署文档

本文档提供 RunBa 项目在 Ubuntu 24 服务器上的完整部署方案（非 Docker 方式）。

## 目录

- [1. 服务器环境准备](#1-服务器环境准备)
- [2. 安装系统依赖](#2-安装系统依赖)
- [3. 创建部署目录和用户](#3-创建部署目录和用户)
- [4. 同步项目代码](#4-同步项目代码)
- [5. 数据库配置](#5-数据库配置)
- [6. 后端部署](#6-后端部署)
- [7. 前端部署](#7-前端部署)
- [8. Nginx 反向代理配置](#8-nginx-反向代理配置)
- [9. HTTPS 配置（可选）](#9-https-配置可选)
- [10. 服务验证](#10-服务验证)
- [11. 日常发布流程](#11-日常发布流程)
- [12. 常见问题排查](#12-常见问题排查)

---

## 1. 服务器环境准备

### 1.1 更新系统

```bash
sudo apt update
sudo apt upgrade -y
```

### 1.2 开放防火墙端口

如果启用了 UFW 防火墙：

```bash
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 22/tcp
sudo ufw enable
sudo ufw status
```

---

## 2. 安装系统依赖

### 2.1 安装基础工具

```bash
sudo apt install -y nginx mysql-server redis-server curl git build-essential
```

### 2.2 安装 Node.js 20

```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs
node -v
npm -v
```

### 2.3 安装 Go 1.24+

```bash
# 下载并安装 Go
cd /tmp
curl -LO https://go.dev/dl/go1.24.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24.5.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee /etc/profile.d/go.sh
source /etc/profile.d/go.sh

# 验证安装
go version
```

### 2.4 启动并配置系统服务

```bash
# 启动 MySQL
sudo systemctl enable --now mysql

# 启动 Redis
sudo systemctl enable --now redis-server

# 启动 Nginx
sudo systemctl enable --now nginx

# 验证服务状态
sudo systemctl status mysql
sudo systemctl status redis-server
sudo systemctl status nginx
```

---

## 3. 创建部署目录和用户

### 3.1 创建专用运行用户

```bash
sudo useradd -r -s /usr/sbin/nologin -m runba
```

### 3.2 创建部署目录

```bash
sudo mkdir -p /opt/runba
sudo chown -R $USER:$USER /opt/runba
```

部署目录结构：
```
/opt/runba/
├── runba-backend/    # 后端项目
└── runba-frontend/   # 前端项目
```

---

## 4. 同步项目代码

### 方式 A：使用 Git（推荐）

```bash
cd /opt/runba
git clone <your-repo-url> .

# 或者如果已有仓库
git pull origin main
```

### 方式 B：使用 rsync 从开发机同步

在开发机（Mac）上执行：

```bash
rsync -avz --delete \
  --exclude '.git' \
  --exclude 'runba-frontend/node_modules' \
  --exclude 'runba-frontend/.next' \
  --exclude 'runba-backend/runba-backend' \
  /Users/jio/codespace/go/src/RunBa/ \
  <user>@<server_ip>:/opt/runba/
```

---

## 5. 数据库配置

### 5.1 配置 MySQL

```bash
# 首次安装需要安全配置
sudo mysql_secure_installation
```

按提示设置：
- 设置 root 密码
- 移除匿名用户
- 禁止 root 远程登录
- 删除测试数据库

### 5.2 创建数据库和用户

登录 MySQL：

```bash
sudo mysql -u root -p
```

执行以下 SQL：

```sql
-- 创建数据库
CREATE DATABASE rb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建专用用户（请修改密码）
CREATE USER 'runba'@'localhost' IDENTIFIED BY 'YourStrongPassword123!';
GRANT ALL PRIVILEGES ON rb.* TO 'runba'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

### 5.3 导入数据表结构

```bash
# 导入锻刀所和刀匠数据
mysql -u runba -p rb < /opt/runba/runba-backend/forges_swordsmiths.sql

# 导入证书数据
mysql -u runba -p rb < /opt/runba/runba-backend/certificates.sql
```

### 5.4 创建 users 表

登录 MySQL：

```bash
mysql -u runba -p rb
```

执行以下 SQL：

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

-- 插入默认管理员（账号：admin，密码：ChangeMe123!）
-- ⚠️ 上线后请立即修改密码
INSERT INTO users(username, password, nickname, status, created_at, updated_at)
VALUES('admin', '$2a$10$Q7zEbOkd728rwWywoJzxuOTaPoWrhlWdnPP8gemPxNFzNLYIsZ.1m', 'admin', 1, NOW(), NOW());

EXIT;
```

---

## 6. 后端部署

### 6.1 配置后端

编辑配置文件：

```bash
cd /opt/runba/runba-backend
cp config.yaml config.yaml.bak
nano config.yaml
```

关键配置项（请根据实际情况修改）：

```yaml
app:
  name: "RunBa Backend"
  version: "1.0.0"

server:
  port: 8080
  mode: release  # 生产环境必须设置为 release

database:
  host: localhost
  port: 3306
  username: runba
  password: YourStrongPassword123!  # 与上面创建的数据库密码一致
  dbname: rb
  max_idle_conns: 10
  max_open_conns: 100

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "Your-Very-Strong-JWT-Secret-At-Least-32-Characters-Long!"  # ⚠️ 必须修改为强随机字符串
  expiration: 168  # Token 过期时间（小时）

log:
  level: info  # debug, info, warn, error
```

### 6.2 编译后端

```bash
cd /opt/runba/runba-backend

# 下载依赖
go mod download

# 编译二进制文件
go build -o runba-backend main.go

# 验证编译结果
./runba-backend --help
```

### 6.3 创建 systemd 服务

创建服务文件：

```bash
sudo nano /etc/systemd/system/runba-backend.service
```

写入以下内容：

```ini
[Unit]
Description=RunBa Backend API Service
Documentation=https://github.com/your-org/runba
After=network.target mysql.service redis-server.service
Wants=mysql.service redis-server.service

[Service]
Type=simple
User=runba
Group=runba
WorkingDirectory=/opt/runba/runba-backend
ExecStart=/opt/runba/runba-backend/runba-backend
Restart=always
RestartSec=3
StandardOutput=journal
StandardError=journal

# 安全加固
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

### 6.4 授权并启动后端服务

```bash
# 修改所有权
sudo chown -R runba:runba /opt/runba/runba-backend

# 重载 systemd 配置
sudo systemctl daemon-reload

# 启动并设置开机自启
sudo systemctl enable --now runba-backend

# 查看服务状态
sudo systemctl status runba-backend

# 查看实时日志
sudo journalctl -u runba-backend -f
```

---

## 7. 前端部署

### 7.1 安装依赖

```bash
cd /opt/runba/runba-frontend
npm ci
```

### 7.2 配置前端 API 地址（可选）

如果需要修改 API 配置，编辑：

```bash
nano src/lib/api.ts
```

确认 API 配置正确（默认使用本地代理）：

```typescript
export const API_CONFIG = {
  LOCAL: '/api/v1',  // 通过 Nginx 反向代理
  REMOTE: 'http://your-server-ip:8080/api/v1'
}
```

### 7.3 构建前端

```bash
npm run build
```

构建成功后会生成 `.next` 目录。

### 7.4 创建 systemd 服务

创建服务文件：

```bash
sudo nano /etc/systemd/system/runba-frontend.service
```

写入以下内容：

```ini
[Unit]
Description=RunBa Frontend (Next.js)
Documentation=https://github.com/your-org/runba
After=network.target

[Service]
Type=simple
User=runba
Group=runba
WorkingDirectory=/opt/runba/runba-frontend
ExecStart=/usr/bin/npm run start -- -p 3000 -H 127.0.0.1
Restart=always
RestartSec=3
Environment=NODE_ENV=production
StandardOutput=journal
StandardError=journal

# 安全加固
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

### 7.5 授权并启动前端服务

```bash
# 修改所有权
sudo chown -R runba:runba /opt/runba/runba-frontend

# 重载 systemd 配置
sudo systemctl daemon-reload

# 启动并设置开机自启
sudo systemctl enable --now runba-frontend

# 查看服务状态
sudo systemctl status runba-frontend

# 查看实时日志
sudo journalctl -u runba-frontend -f
```

---

## 8. Nginx 反向代理配置

### 8.1 创建站点配置文件

```bash
sudo nano /etc/nginx/sites-available/runba.conf
```

写入以下内容：

```nginx
server {
    listen 80;
    server_name your-domain.com;  # 修改为你的域名或服务器 IP

    client_max_body_size 20m;

    # 后端 API 代理
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # 前端应用代理
    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;

        # WebSocket 支持
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 8.2 启用站点配置

```bash
# 创建软链接
sudo ln -sf /etc/nginx/sites-available/runba.conf /etc/nginx/sites-enabled/runba.conf

# 移除默认站点配置（可选）
sudo rm -f /etc/nginx/sites-enabled/default

# 测试配置文件语法
sudo nginx -t

# 重载 Nginx 配置
sudo systemctl reload nginx
```

---

## 9. HTTPS 配置（可选）

### 9.1 安装 Certbot

```bash
sudo apt install -y certbot python3-certbot-nginx
```

### 9.2 获取 SSL 证书

```bash
# 确保域名已解析到服务器 IP
sudo certbot --nginx -d your-domain.com
```

按提示操作：
1. 输入邮箱地址
2. 同意服务条款
3. 选择是否重定向 HTTP 到 HTTPS（推荐选择重定向）

### 9.3 验证自动续期

```bash
# 测试续期命令
sudo certbot renew --dry-run

# 查看续期定时任务
sudo systemctl status certbot.timer
```

---

## 10. 服务验证

### 10.1 检查后端健康状态

```bash
curl http://127.0.0.1:8080/health
```

预期输出：
```json
{"status":"ok"}
```

### 10.2 检查服务状态

```bash
sudo systemctl status runba-backend
sudo systemctl status runba-frontend
sudo systemctl status nginx
sudo systemctl status mysql
sudo systemctl status redis-server
```

### 10.3 查看服务日志

```bash
# 后端日志
sudo journalctl -u runba-backend -n 50 --no-pager

# 前端日志
sudo journalctl -u runba-frontend -n 50 --no-pager

# Nginx 访问日志
sudo tail -f /var/log/nginx/access.log

# Nginx 错误日志
sudo tail -f /var/log/nginx/error.log
```

### 10.4 浏览器访问测试

```
http://your-domain.com        # 访问前端
http://your-domain.com/api/v1 # 访问后端 API
```

### 10.5 登录测试

使用默认管理员账号登录：
- 用户名：`admin`
- 密码：`ChangeMe123!`

**⚠️ 重要：首次登录后请立即修改密码！**

---

## 11. 日常发布流程

### 11.1 从开发机推送代码

```bash
# 在开发机上
cd /Users/jio/codespace/go/src/RunBa
git add .
git commit -m "更新描述"
git push origin main
```

### 11.2 在服务器上更新部署

```bash
# SSH 登录服务器
ssh user@your-server-ip

# 拉取最新代码
cd /opt/runba
git pull origin main

# 更新后端
cd /opt/runba/runba-backend
go build -o runba-backend main.go
sudo systemctl restart runba-backend
sudo systemctl status runba-backend

# 更新前端
cd /opt/runba/runba-frontend
npm ci
npm run build
sudo systemctl restart runba-frontend
sudo systemctl status runba-frontend

# 查看服务日志确认无误
sudo journalctl -u runba-backend -f
sudo journalctl -u runba-frontend -f
```

### 11.3 回滚操作（如遇问题）

```bash
# 回滚代码
cd /opt/runba
git log --oneline -n 5  # 查看最近提交
git reset --hard <commit-hash>

# 重新编译和重启
# （按上述更新步骤操作）
```

---

## 12. 常见问题排查

### 12.1 后端无法启动

```bash
# 查看详细错误日志
sudo journalctl -u runba-backend -n 100 --no-pager

# 检查配置文件
nano /opt/runba/runba-backend/config.yaml

# 检查数据库连接
mysql -u runba -p -h localhost rb

# 检查端口占用
sudo lsof -i :8080
```

### 12.2 前端构建失败

```bash
# 清理缓存重新构建
cd /opt/runba/runba-frontend
rm -rf node_modules .next
npm ci
npm run build

# 检查 Node.js 版本
node -v  # 应该是 v20.x
```

### 12.3 数据库连接失败

```bash
# 检查 MySQL 服务状态
sudo systemctl status mysql

# 检查数据库用户权限
sudo mysql -u root -p
# 执行：SELECT User, Host FROM mysql.user WHERE User='runba';

# 测试连接
mysql -u runba -p -h localhost rb
```

### 12.4 Nginx 502 错误

```bash
# 检查后端和前端服务是否运行
sudo systemctl status runba-backend
sudo systemctl status runba-frontend

# 检查端口监听
sudo netstat -tlnp | grep -E '8080|3000'

# 查看 Nginx 错误日志
sudo tail -f /var/log/nginx/error.log
```

### 12.5 权限问题

```bash
# 确保文件所有权正确
sudo chown -R runba:runba /opt/runba/runba-backend
sudo chown -R runba:runba /opt/runba/runba-frontend

# 检查可执行权限
chmod +x /opt/runba/runba-backend/runba-backend
```

---

## 附录：安全建议

1. **修改默认密码**
   - 管理员账号密码
   - 数据库密码
   - JWT Secret

2. **定期备份数据库**
   ```bash
   # 创建备份脚本
   mysqldump -u runba -p rb > /backup/rb_$(date +%Y%m%d_%H%M%S).sql
   ```

3. **配置防火墙规则**
   - 只开放必要端口（80, 443, 22）
   - 限制 SSH 访问

4. **启用日志审计**
   ```bash
   # 设置日志轮转
   sudo nano /etc/logrotate.d/runba
   ```

5. **监控服务状态**
   - 配置监控告警（如 Prometheus + Grafana）
   - 定期检查系统资源使用情况

---

## 技术支持

如遇问题请查看：
- 后端日志：`sudo journalctl -u runba-backend -f`
- 前端日志：`sudo journalctl -u runba-frontend -f`
- Nginx 日志：`/var/log/nginx/error.log`

项目仓库：https://github.com/your-org/runba
