#!/bin/bash

# Docker Compose V2 安装/升级脚本

set -e

echo "🔧 Docker Compose V2 安装/升级脚本"
echo "================================="

# 检查是否已安装 Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    echo "💡 安装方法: https://docs.docker.com/engine/install/"
    exit 1
fi

# 检查当前 Docker Compose 状态
echo "🔍 检查当前 Docker Compose 状态..."

if docker compose version &> /dev/null; then
    echo "✅ Docker Compose V2 已安装"
    docker compose version
    echo ""
    echo "如果您想使用最新版本，可以继续执行升级。"
    read -p "是否继续升级？(y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "取消升级"
        exit 0
    fi
elif command -v docker-compose &> /dev/null; then
    echo "⚠️  检测到旧版 Docker Compose V1"
    docker-compose version
    echo ""
    echo "建议升级到 Docker Compose V2 (Go 版本)"
    read -p "是否继续安装 V2 版本？(Y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Nn]$ ]]; then
        echo "取消安装"
        exit 0
    fi
else
    echo "📦 未检测到 Docker Compose，将安装 V2 版本"
fi

# 检测系统架构
ARCH=$(uname -m)
case $ARCH in
    x86_64)
        ARCH="x86_64"
        ;;
    aarch64|arm64)
        ARCH="aarch64"
        ;;
    *)
        echo "❌ 不支持的架构: $ARCH"
        exit 1
        ;;
esac

echo "🖥️  检测到系统架构: $ARCH"

# 获取最新版本号
echo "🔍 获取最新版本信息..."
LATEST_VERSION=$(curl -s https://api.github.com/repos/docker/compose/releases/latest | grep -Po '"tag_name": "\K[^"]*')

if [ -z "$LATEST_VERSION" ]; then
    echo "❌ 无法获取最新版本信息，使用默认版本 v2.21.0"
    LATEST_VERSION="v2.21.0"
fi

echo "📦 最新版本: $LATEST_VERSION"

# 创建插件目录
echo "📁 创建 Docker CLI 插件目录..."
mkdir -p ~/.docker/cli-plugins/

# 下载 Docker Compose V2
echo "⬇️  下载 Docker Compose V2..."
DOWNLOAD_URL="https://github.com/docker/compose/releases/download/${LATEST_VERSION}/docker-compose-linux-${ARCH}"

if curl -L "$DOWNLOAD_URL" -o ~/.docker/cli-plugins/docker-compose; then
    echo "✅ 下载完成"
else
    echo "❌ 下载失败，请检查网络连接"
    exit 1
fi

# 设置可执行权限
echo "🔐 设置可执行权限..."
chmod +x ~/.docker/cli-plugins/docker-compose

# 验证安装
echo "🧪 验证安装..."
if docker compose version &> /dev/null; then
    echo "✅ Docker Compose V2 安装成功！"
    echo ""
    docker compose version
    echo ""
    echo "🎉 现在可以使用 'docker compose' 命令了！"
    echo ""
    echo "📋 常用命令对比:"
    echo "  旧版: docker-compose up -d"
    echo "  新版: docker compose up -d"
    echo ""
    echo "💡 提示: 新版本速度更快，功能更强大！"
    
    # 检查是否有旧版本
    if command -v docker-compose &> /dev/null; then
        echo ""
        echo "⚠️  检测到旧版 docker-compose 仍然存在"
        echo "建议卸载旧版本以避免混淆:"
        echo "  sudo apt remove docker-compose  # Ubuntu/Debian"
        echo "  sudo yum remove docker-compose  # CentOS/RHEL"
        echo "  pip uninstall docker-compose    # 如果是通过 pip 安装的"
    fi
else
    echo "❌ 安装验证失败"
    exit 1
fi

echo ""
echo "🚀 安装完成！现在可以运行 Luna Dial Server 了："
echo "  ./start.sh"
