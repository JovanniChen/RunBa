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
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
sudo ufw status
```

---

## 2. 安装系统依赖

### 2.1 安装基础工具

```bash
sudo apt install -y nginx mysql-server curl git build-essential
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

# 国内服务器配置 Go 模块代理（海外服务器可跳过）
go env -w GOPROXY=https://goproxy.cn,direct

# 验证安装
go version
```

### 2.4 启动并配置系统服务

```bash
# 启动 MySQL
sudo systemctl enable --now mysql

# 启动 Nginx
sudo systemctl enable --now nginx

# 验证服务状态
sudo systemctl status mysql
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
# 部署目录（当前用户拥有，用于拉取代码和构建）
sudo mkdir -p /opt/runba
sudo chown -R $USER:$USER /opt/runba

# 备份目录
sudo mkdir -p /opt/runba/backups
```

部署目录结构：
```
/opt/runba/
├── runba-backend/    # 后端项目
├── runba-frontend/   # 前端项目
└── backups/          # 数据库备份
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
-- 上线后请立即修改密码
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
nano config.yaml
```

关键配置项（请根据实际情况修改）：

```yaml
# 应用基本配置
app:
  name: 'RunBa Backend API'
  version: '1.0.0'
  environment: 'production'

# 服务器配置
server:
  host: '0.0.0.0'
  port: '8080'
  mode: 'release'  # 生产环境必须设置为 release
  read_timeout: 30
  write_timeout: 30
  idle_timeout: 60

# 数据库配置
database:
  host: '127.0.0.1'
  port: '3306'
  user: 'runba'
  password: 'YourStrongPassword123!'  # 与上面创建的数据库密码一致
  name: 'rb'
  charset: 'utf8mb4'
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600
  log_level: 'warn'

# JWT 配置
jwt:
  secret: 'Your-JWT-Secret-Key'
  expiration: 24
  refresh_expiration: 168

# 日志配置
log:
  level: 'info'
  format: 'json'
  output: 'stdout'
```

### 6.2 编译后端

```bash
cd /opt/runba/runba-backend

# 下载依赖
go mod download

# 编译二进制文件（保留上一个版本用于快速回滚）
[ -f runba-backend ] && cp runba-backend runba-backend.prev
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
After=network.target mysql.service
Wants=mysql.service

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

# 资源限制
LimitNOFILE=65535

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

### 7.1 配置 Next.js standalone 输出

编辑 `next.config.ts`，添加 `output: 'standalone'`：

```typescript
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: 'standalone',
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/api/:path*',
      },
    ]
  },
};

export default nextConfig;
```

> standalone 模式会在构建时生成独立运行目录，不依赖完整的 node_modules，部署体积更小。

### 7.2 安装依赖并构建

```bash
cd /opt/runba/runba-frontend
npm ci
npm run build
```

构建成功后会生成 `.next/standalone` 目录。

### 7.3 创建 systemd 服务

创建服务文件：

```bash
sudo nano /etc/systemd/system/runba-frontend.service
```

写入以下内容：

```ini
[Unit]
Description=RunBa Frontend (Next.js)
After=network.target

[Service]
Type=simple
User=runba
Group=runba
WorkingDirectory=/opt/runba/runba-frontend
ExecStart=/usr/bin/node .next/standalone/server.js
Restart=always
RestartSec=3
Environment=NODE_ENV=production
Environment=PORT=3000
Environment=HOSTNAME=127.0.0.1
StandardOutput=journal
StandardError=journal

# 资源限制
LimitNOFILE=65535

# 安全加固
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

### 7.4 准备 standalone 静态资源

standalone 模式不会自动复制 `public` 和 `.next/static` 目录，需要手动链接：

```bash
cd /opt/runba/runba-frontend

# 链接静态资源（构建后执行）
cp -r public .next/standalone/public
cp -r .next/static .next/standalone/.next/static
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
# API 限流：每个 IP 每秒 20 个请求，突发 40 个
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=20r/s;

server {
    listen 80;
    server_name your-domain.com;  # 修改为你的域名或服务器 IP

    client_max_body_size 20m;

    # --- 安全响应头 ---
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # --- Gzip 压缩 ---
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript image/svg+xml;

    # --- Next.js 静态资源（带 hash，长缓存）---
    location /_next/static/ {
        proxy_pass http://127.0.0.1:3000;
        expires 365d;
        add_header Cache-Control "public, immutable";
    }

    # --- 后端 API 代理 ---
    location /api/ {
        limit_req zone=api_limit burst=40 nodelay;

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

    # --- 前端应用代理 ---
    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;

        # WebSocket 支持（Next.js HMR 等）
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

### 10.2 检查所有服务状态

```bash
sudo systemctl status runba-backend
sudo systemctl status runba-frontend
sudo systemctl status nginx
sudo systemctl status mysql
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

**首次登录后请立即修改密码。**

---

## 11. 日常发布流程

### 11.1 从开发机推送代码

```bash
cd /Users/jio/codespace/go/src/RunBa
git add .
git commit -m "更新描述"
git push origin main
```

### 11.2 在服务器上更新部署

可以手动执行以下步骤，也可以使用 [11.3 节的部署脚本](#113-一键部署脚本)。

```bash
# SSH 登录服务器
ssh user@your-server-ip

# 备份数据库
mysqldump -u runba -p rb > /opt/runba/backups/rb_$(date +%Y%m%d_%H%M%S).sql

# 拉取最新代码
cd /opt/runba
git pull origin main

# 更新后端（保留上一个二进制用于回滚）
cd /opt/runba/runba-backend
cp runba-backend runba-backend.prev
go build -o runba-backend main.go
sudo systemctl restart runba-backend
sudo systemctl status runba-backend

# 更新前端
cd /opt/runba/runba-frontend
npm ci
npm run build
cp -r public .next/standalone/public
cp -r .next/static .next/standalone/.next/static
sudo chown -R runba:runba /opt/runba/runba-frontend
sudo systemctl restart runba-frontend
sudo systemctl status runba-frontend

# 查看服务日志确认无误
sudo journalctl -u runba-backend -n 20 --no-pager
sudo journalctl -u runba-frontend -n 20 --no-pager
```

### 11.3 一键部署脚本

创建部署脚本：

```bash
nano /opt/runba/deploy.sh
chmod +x /opt/runba/deploy.sh
```

写入以下内容：

```bash
#!/bin/bash
set -e

DEPLOY_DIR="/opt/runba"
BACKUP_DIR="$DEPLOY_DIR/backups"

echo "=== RunBa 部署开始 ==="

# 1. 备份数据库
echo "[1/6] 备份数据库..."
mysqldump -u runba -p"$1" rb > "$BACKUP_DIR/rb_$(date +%Y%m%d_%H%M%S).sql"
# 保留最近 10 个备份
ls -t "$BACKUP_DIR"/rb_*.sql | tail -n +11 | xargs -r rm

# 2. 拉取代码
echo "[2/6] 拉取最新代码..."
cd "$DEPLOY_DIR"
git pull origin main

# 3. 编译后端
echo "[3/6] 编译后端..."
cd "$DEPLOY_DIR/runba-backend"
[ -f runba-backend ] && cp runba-backend runba-backend.prev
go build -o runba-backend main.go

# 4. 构建前端
echo "[4/6] 构建前端..."
cd "$DEPLOY_DIR/runba-frontend"
npm ci
npm run build
cp -r public .next/standalone/public
cp -r .next/static .next/standalone/.next/static

# 5. 修复权限并重启服务
echo "[5/6] 重启服务..."
sudo chown -R runba:runba "$DEPLOY_DIR/runba-backend"
sudo chown -R runba:runba "$DEPLOY_DIR/runba-frontend"
sudo systemctl restart runba-backend
sudo systemctl restart runba-frontend

# 6. 验证
echo "[6/6] 验证服务状态..."
sleep 3
sudo systemctl is-active runba-backend
sudo systemctl is-active runba-frontend

echo "=== 部署完成 ==="
```

使用方式：

```bash
# 参数为数据库密码
bash /opt/runba/deploy.sh 'YourDBPassword'
```

### 11.4 回滚操作（如遇问题）

**快速回滚**（使用上一个二进制，无需重新编译）：

```bash
# 回滚后端
cd /opt/runba/runba-backend
cp runba-backend.prev runba-backend
sudo systemctl restart runba-backend

# 回滚数据库（如有需要）
mysql -u runba -p rb < /opt/runba/backups/rb_<timestamp>.sql
```

**代码回滚**：

```bash
cd /opt/runba
git log --oneline -n 5  # 查看最近提交
git revert <commit-hash> # 生成一个撤销提交，不会丢失历史

# 重新编译和重启（按上述更新步骤操作）
```

---

## 12. 常见问题排查

### 12.1 后端无法启动

```bash
# 查看详细错误日志
sudo journalctl -u runba-backend -n 100 --no-pager

# 检查配置文件
cat /opt/runba/runba-backend/config.yaml

# 检查数据库连接
mysql -u runba -p -h localhost rb

# 检查端口占用
sudo ss -tlnp | grep 8080
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
sudo ss -tlnp | grep -E '8080|3000'

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

## 附录 A：日志轮转配置

创建 logrotate 配置，防止 journal 日志过大：

```bash
sudo nano /etc/logrotate.d/runba-nginx
```

写入：

```
/var/log/nginx/access.log
/var/log/nginx/error.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    sharedscripts
    postrotate
        [ -f /var/run/nginx.pid ] && kill -USR1 $(cat /var/run/nginx.pid)
    endscript
}
```

配置 systemd journal 日志大小限制：

```bash
sudo nano /etc/systemd/journald.conf
```

设置：

```ini
[Journal]
SystemMaxUse=500M
```

重启 journald：

```bash
sudo systemctl restart systemd-journald
```

## 附录 B：数据库自动备份

创建 crontab 定时备份任务：

```bash
sudo crontab -e
```

添加每天凌晨 3 点备份：

```
0 3 * * * mysqldump -u runba -p'YourDBPassword' rb | gzip > /opt/runba/backups/rb_$(date +\%Y\%m\%d).sql.gz && find /opt/runba/backups -name "rb_*.sql.gz" -mtime +30 -delete
```

## 附录 C：安全检查清单

部署完成后逐项确认：

- [ ] 管理员默认密码已修改
- [ ] config.yaml 中的数据库密码已修改
- [ ] config.yaml 中的 JWT Secret 已修改
- [ ] config.yaml 中 server.mode 为 `release`
- [ ] UFW 防火墙已启用，仅开放 22/80/443
- [ ] MySQL 已禁止 root 远程登录
- [ ] HTTPS 已配置（如有域名）
