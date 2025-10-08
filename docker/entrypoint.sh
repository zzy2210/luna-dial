#!/bin/sh
# Docker entrypoint script for luna-dial backend

# 使用环境变量或默认配置文件路径
CONFIG_FILE=${CONFIG_FILE:-/app/config/backend.ini}

# 启动服务器，传入配置文件参数
exec ./luna-dial-server --config "$CONFIG_FILE"