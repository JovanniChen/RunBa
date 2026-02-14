# RunBa Docker Compose 部署方案（Mac 开发 -> Linux 服务器）

适用目录：
- `runba-backend`
- `runba-frontend`
- `docker-compose.yml`

目标：
- 你在 MacBook 开发；
- Linux 服务器上使用 Docker Compose 一键部署；
- 避免 Mac 与 Linux 构建产物不兼容问题（全部在 Linux 构建镜像）。

## 1. 服务器安装 Docker / Compose（CentOS 示例）

```bash
sudo dnf -y install dnf-plugins-core
sudo dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo dnf -y install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
sudo systemctl enable --now docker
docker --version
docker compose version
```

如果你是 CentOS 7（没有 `dnf`），使用：

```bash
sudo yum -y install yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum -y install docker-ce docker-ce-cli containerd.io
sudo systemctl enable --now docker
docker --version
```

CentOS 7 常见情况是没有 `docker compose` 插件，可安装独立版 Compose：

```bash
sudo curl -L "https://github.com/docker/compose/releases/download/v2.27.0/docker-compose-linux-x86_64" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
docker-compose version
```

说明：如果你安装的是独立版，请把文档里的 `docker compose ...` 改成 `docker-compose ...`。

可选（免 sudo 执行 docker）：

```bash
sudo usermod -aG docker $USER
newgrp docker
```

防火墙开放 80/443（如果启用了 firewalld）：

```bash
sudo systemctl enable --now firewalld
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

如果你启用了 SELinux 且遇到容器挂载权限错误，可执行：

```bash
sudo chcon -Rt svirt_sandbox_file_t /opt/runba
```

## 2. 从 Mac 同步源码到 Linux

推荐 Git：

```bash
mkdir -p /opt/runba && cd /opt/runba
git clone <your-repo-url> .
```

或从 Mac 用 `rsync`：

```bash
rsync -avz --delete \
  --exclude '.git' \
  --exclude 'runba-frontend/node_modules' \
  --exclude 'runba-frontend/.next' \
  /Users/jio/codespace/go/src/RunBa/ \
  <user>@<server_ip>:/opt/runba/
```

## 3. 上线前必须改的配置

编辑 `runba-backend/config.docker.yaml`：
1. `database.password` 必须和 `docker-compose.yml` 中 `MYSQL_PASSWORD` 一致。
2. `jwt.secret` 改为你自己的强随机字符串（>= 32 字符）。
3. 根据需要调整日志级别。

编辑 `docker-compose.yml`：
1. `MYSQL_ROOT_PASSWORD` 改成强密码。
2. `MYSQL_PASSWORD` 改成强密码（并同步到 `config.docker.yaml`）。

## 4. 启动

在 Linux 服务器执行：

```bash
cd /opt/runba
docker compose build
docker compose up -d
```

查看状态：

```bash
docker compose ps
docker compose logs -f --tail=200
```

如果你要重置数据库并重新执行初始化 SQL（会清空数据）：

```bash
docker compose down -v
docker compose up -d --build
```

## 5. 访问与端口

- 对外入口：`http://<server_ip_or_domain>/`
- Nginx 容器暴露：`80`
- 内部服务：
  - frontend: `3000`（仅容器网络）
  - backend: `8080`（仅容器网络）
  - mysql: `3306`（仅容器网络）
  - redis: `6379`（仅容器网络）

## 6. 健康检查与验证

```bash
curl -I http://127.0.0.1/
docker compose exec backend sh -lc 'wget -qO- http://127.0.0.1:8080/health'
docker compose logs backend --tail=100
docker compose logs frontend --tail=100
docker compose logs nginx --tail=100
```

## 7. 日常发布流程（Mac -> Linux）

```bash
# 1) Mac 提交并推送代码

# 2) Linux 拉取并重建
cd /opt/runba
git pull
docker compose build
docker compose up -d
```

仅重启某个服务：

```bash
docker compose up -d --build backend
docker compose up -d --build frontend
```

## 8. 数据持久化与备份

MySQL 和 Redis 使用命名卷：
- `runba_mysql_data`
- `runba_redis_data`

导出 MySQL 备份示例：

```bash
docker compose exec -T mysql sh -lc 'mysqldump -uroot -p"$MYSQL_ROOT_PASSWORD" rb' > rb_backup.sql
```

恢复示例：

```bash
cat rb_backup.sql | docker compose exec -T mysql sh -lc 'mysql -uroot -p"$MYSQL_ROOT_PASSWORD" rb'
```

## 9. HTTPS（生产建议）

当前 Compose 方案默认 HTTP。生产建议：
1. 在宿主机再加一层 Nginx/Caddy 做 HTTPS 终止并反代到 `127.0.0.1:80`。
2. 或接入云厂商负载均衡证书。

## 10. 首次上线后立即执行

1. 登录后台改默认管理员密码（默认 admin / ChangeMe123!）。
2. 更换 `jwt.secret`。
3. 收敛服务器防火墙端口（只开放 80/443 和 SSH）。
