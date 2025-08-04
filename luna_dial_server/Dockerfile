# ==========================================
# 阶段 1: 构建阶段 (Build Stage)
# ==========================================
FROM golang:1.23-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum（利用 Docker 层缓存）
COPY go.mod go.sum ./

# 下载依赖（这一层会被缓存，除非 go.mod/go.sum 改变）
RUN go mod download && go mod verify && go mod tidy

# 复制源代码
COPY . .

# 构建应用（优化编译参数）
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o luna-dial-server \
    ./cmd/main.go

# ==========================================
# 阶段 2: 运行阶段 (Runtime Stage)
# ==========================================
FROM alpine:3.18

# 安装运行时依赖（最小化）
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    wget \
    && adduser -D -s /bin/sh -u 1001 appuser

# 设置时区
ENV TZ=Asia/Shanghai

# 切换到非 root 用户
USER appuser
WORKDIR /home/appuser

# 从构建阶段复制应用二进制文件
COPY --from=builder --chown=appuser:appuser /app/luna-dial-server ./luna-dial-server

# 复制配置文件和迁移文件
COPY --from=builder --chown=appuser:appuser /app/configs ./configs
COPY --from=builder --chown=appuser:appuser /app/migrations ./migrations

# 暴露端口
EXPOSE 8081

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# 设置入口点
ENTRYPOINT ["./luna-dial-server"]
