#!/bin/bash

# Luna Dial Docker 启动脚本
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "======================================"
echo "Luna Dial Docker 环境启动"
echo "======================================"

# 检查 .env 文件
if [ ! -f .env ]; then
    echo "错误：.env 文件不存在！"
    echo "请先配置 .env 文件"
    exit 1
fi

# 创建必要的目录
echo "创建必要的目录..."
mkdir -p data/postgres
mkdir -p config

# 构建镜像
echo "构建 Docker 镜像..."
docker compose build

# 启动服务
echo "启动服务..."
docker compose up -d

# 等待服务启动
echo "等待服务启动..."
sleep 5

# 检查服务状态
echo ""
echo "服务状态："
docker compose ps

echo ""
echo "======================================"
echo "Luna Dial 已启动！"
echo "访问地址：http://localhost:10755"
echo "API 地址：http://localhost:10755/api/v1"
echo ""
echo "默认管理员账号："
echo "  用户名：admin"
echo "  密码：admin@123"
echo ""
echo "常用命令："
echo "  查看日志：docker compose logs -f"
echo "  停止服务：docker compose down"
echo "  重启服务：docker compose restart"
echo "======================================"