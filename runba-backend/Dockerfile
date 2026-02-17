# ========== 构建阶段 ==========
FROM golang:1.24-alpine AS builder

# 使用国内 Go 模块代理
ENV GOPROXY=https://goproxy.cn,direct

# 使用阿里云 Alpine 镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

WORKDIR /app

# 先复制依赖文件，利用 Docker 缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码并编译
COPY . .
RUN CGO_ENABLED=0 go build -o runba-backend main.go

# ========== 运行阶段 ==========
FROM alpine:3.20

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

COPY --from=builder /app/runba-backend .
COPY --from=builder /app/config.yaml .

EXPOSE 8080

CMD ["./runba-backend"]
