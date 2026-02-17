# RunBa Docker Compose 部署文档

适用于 Ubuntu 24 服务器（中国大陆网络环境）。

---

## 1. 安装 Docker

```bash
sudo apt update
sudo apt install -y ca-certificates curl gnupg

# 使用阿里云 Docker 镜像源安装
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://mirrors.aliyun.com/docker-ce/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://mirrors.aliyun.com/docker-ce/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# 将当前用户加入 docker 组（免 sudo）
sudo usermod -aG docker $USER
newgrp docker

# 验证
docker --version
docker compose version
```

## 2. 配置 Docker 镜像加速

```bash
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<'EOF'
{
  "registry-mirrors": [
    "https://docker.1ms.run",
    "https://docker.xuanyuan.me"
  ]
}
EOF

sudo systemctl daemon-reload
sudo systemctl restart docker
```

> 镜像加速地址可能会失效，如果拉取镜像缓慢，请搜索最新可用的国内 Docker 镜像源替换。

## 3. 同步代码到服务器

### 方式 A：Git

```bash
cd /opt
sudo mkdir -p runba && sudo chown $USER:$USER runba
cd runba
git clone <your-repo-url> .
```

### 方式 B：rsync（从开发机）

```bash
rsync -avz --delete \
  --exclude '.git' \
  --exclude 'runba-frontend/node_modules' \
  --exclude 'runba-frontend/.next' \
  --exclude 'runba-backend/runba-backend' \
  /Users/jio/codespace/go/src/RunBa/ \
  <user>@<server_ip>:/opt/runba/
```

## 4. 修改配置

编辑 `docker/config.yaml`，修改数据库密码和 JWT Secret：

```bash
cd /opt/runba
nano docker/config.yaml
```

同时确保 `docker-compose.yml` 中 MySQL 的密码与 `docker/config.yaml` 一致：

```yaml
# docker-compose.yml 中的 MySQL 环境变量
MYSQL_PASSWORD: runba123    # ← 与 docker/config.yaml 中 database.password 一致
```

## 5. 启动

```bash
cd /opt/runba
docker compose up -d --build
```

首次启动会自动：
- 拉取 MySQL 8.0、Nginx 镜像
- 编译 Go 后端
- 构建 Next.js 前端
- 创建数据库表并插入默认管理员

查看启动状态：

```bash
docker compose ps
docker compose logs -f
```

等所有容器状态为 `running` 后，访问 `http://<服务器IP>` 即可。

默认管理员：`admin` / `ChangeMe123!`

## 6. 日常操作

### 更新部署

```bash
cd /opt/runba
git pull origin main
docker compose up -d --build
```

### 查看日志

```bash
# 所有服务
docker compose logs -f

# 单个服务
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f mysql
docker compose logs -f nginx
```

### 重启单个服务

```bash
docker compose restart backend
```

### 停止所有服务

```bash
docker compose down
```

### 备份数据库

```bash
docker compose exec mysql mysqldump -u runba -prunba123 rb > backup_$(date +%Y%m%d).sql
```

### 恢复数据库

```bash
docker compose exec -T mysql mysql -u runba -prunba123 rb < backup_20260217.sql
```

### 完全重置（清除数据库数据）

```bash
docker compose down -v   # -v 会删除数据卷
docker compose up -d --build
```

## 7. 防火墙

```bash
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

## 8. 文件结构说明

```
RunBa/
├── docker-compose.yml              # 编排文件
├── docker/
│   ├── config.yaml                 # 后端生产配置（挂载进容器）
│   ├── mysql/
│   │   └── init.sql                # 数据库初始化（首次启动自动执行）
│   └── nginx/
│       └── default.conf            # Nginx 反向代理配置
├── runba-backend/
│   ├── Dockerfile                  # 后端构建（多阶段：编译 + alpine 运行）
│   └── ...
└── runba-frontend/
    ├── Dockerfile                  # 前端构建（多阶段：依赖 + 构建 + standalone 运行）
    └── ...
```

## 9. 常见问题

### Docker 镜像拉取慢或失败

更换 `/etc/docker/daemon.json` 中的镜像源，然后：

```bash
sudo systemctl restart docker
```

### 后端连不上数据库

MySQL 容器首次启动需要初始化，后端可能会先启动失败然后自动重启。等 MySQL 健康检查通过后会自动恢复。查看状态：

```bash
docker compose ps
docker compose logs mysql
```

### 数据库 init.sql 没有执行

MySQL 的 `docker-entrypoint-initdb.d` 只在**首次创建数据卷时**执行。如果需要重新初始化：

```bash
docker compose down -v   # 删除数据卷
docker compose up -d --build
```

### 前端构建时 npm install 慢

Dockerfile 中已配置 `registry.npmmirror.com` 国内镜像。如果仍然慢，检查服务器 DNS 是否正常：

```bash
nslookup registry.npmmirror.com
```

### 修改了 Nginx 配置后如何生效

```bash
docker compose restart nginx
```
