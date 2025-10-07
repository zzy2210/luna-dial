# ==========================================
# 阶段 1: 构建前端
# ==========================================
FROM node:20-alpine AS frontend-builder

WORKDIR /app

# 复制前端项目文件
COPY luna_dial_frontend/package*.json ./
RUN npm ci --only=production

COPY luna_dial_frontend/ ./

# 设置 API 地址为相对路径（通过 Caddy 代理）
ENV VITE_API_BASE_URL=

# 构建前端
RUN npm run build

# ==========================================
# 阶段 2: Caddy 运行时
# ==========================================
FROM caddy:2-alpine

# 复制前端构建产物
COPY --from=frontend-builder /app/dist /srv

# 复制 Caddyfile（将通过 volume 挂载，这里作为备份）
COPY docker/Caddyfile /etc/caddy/Caddyfile

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 -O - http://localhost:80/health || exit 1

EXPOSE 80