#!/bin/bash

# Luna Dial Server Docker 启动脚本

set -e

echo "🚀 启动 Luna Dial Server..."

# 检查 Docker 和 Docker Compose
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    exit 1
fi

# 检查 Docker Compose（优先使用 V2 版本）
if docker compose version &> /dev/null; then
    COMPOSE_CMD="docker compose"
    echo "✅ 使用 Docker Compose V2"
elif command -v docker-compose &> /dev/null; then
    COMPOSE_CMD="docker-compose"
    echo "⚠️  使用旧版 Docker Compose V1，建议升级到 V2"
else
    echo "❌ Docker Compose 未安装，请先安装 Docker Compose"
    echo "💡 安装方法: https://docs.docker.com/compose/install/"
    exit 1
fi

# 创建数据目录
echo "📁 创建数据目录..."
sudo mkdir -p /opt/okr/postgres_data
sudo chown -R 999:999 /opt/okr/postgres_data  # PostgreSQL 容器内用户 ID

# 构建并启动服务
echo "📦 构建镜像..."
$COMPOSE_CMD build

echo "🔧 启动服务..."
$COMPOSE_CMD up -d

echo "⏳ 等待服务启动..."
sleep 15

# 检查服务状态
echo "🔍 检查服务状态..."
$COMPOSE_CMD ps

# 测试健康检查
echo "🏥 测试健康检查..."
if curl -f http://localhost:8081/health &> /dev/null; then
    echo "✅ 服务启动成功！"
    echo "🌐 API 地址: http://localhost:8081"
    echo "🔍 健康检查: http://localhost:8081/health"
    echo "📊 版本信息: http://localhost:8081/version"
    echo "🐘 PostgreSQL: localhost:15432 (用户: okr_user, 数据库: okr_db)"
    echo "💾 数据存储: /opt/okr/postgres_data"
else
    echo "❌ 服务启动失败，请检查日志:"
    $COMPOSE_CMD logs luna-dial-server
fi

echo ""
echo "📋 常用命令:"
echo "  查看日志: $COMPOSE_CMD logs -f"
echo "  查看特定服务日志: $COMPOSE_CMD logs -f luna-dial-server"
echo "  停止服务: $COMPOSE_CMD down"
echo "  重启服务: $COMPOSE_CMD restart"
echo "  连接数据库: psql -h localhost -p 15432 -U okr_user -d okr_db"
