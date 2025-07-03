# OKR Web 项目 Makefile

# 默认目标
.DEFAULT_GOAL := help

# 变量定义
GO_CMD = go
SERVER_PATH = cmd/server
BINARY_NAME = okr-server
PORT = 8081
LOG_DIR = log
PID_FILE = $(LOG_DIR)/server.pid
LOG_FILE = $(LOG_DIR)/server.log

# 帮助信息
.PHONY: help
help: ## 显示帮助信息
	@echo "OKR Web 项目可用命令:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# 运行开发服务器
.PHONY: run
run: ## 后台运行开发服务器，日志输出到log目录
	@echo "🚀 启动 OKR 服务器（后台运行）..."
	@mkdir -p $(LOG_DIR)
	@# 检查是否已有服务器在运行
	@if [ -f $(PID_FILE) ]; then \
		if ps -p `cat $(PID_FILE)` > /dev/null 2>&1; then \
			echo "❌ 服务器已在运行 (PID: `cat $(PID_FILE)`)"; \
			echo "   使用 'make stop' 停止服务器"; \
			exit 1; \
		else \
			rm -f $(PID_FILE); \
		fi \
	fi
	@# 检查端口是否被占用
	@if lsof -ti :$(PORT) > /dev/null 2>&1; then \
		echo "❌ 端口$(PORT)已被占用，请先停止相关进程"; \
		echo "   占用进程: $$(lsof -ti :$(PORT))"; \
		echo "   使用 'make stop' 尝试停止"; \
		exit 1; \
	fi
	@# 启动服务器
	@nohup $(GO_CMD) run $(SERVER_PATH)/*.go > $(LOG_FILE) 2>&1 & echo $$! > $(PID_FILE)
	@sleep 3
	@# 验证启动是否成功
	@if ps -p `cat $(PID_FILE)` > /dev/null 2>&1; then \
		echo "✅ 服务器启动成功 (PID: `cat $(PID_FILE)`)"; \
		echo "📋 日志文件: $(LOG_FILE)"; \
		echo "🔍 查看日志: make logs"; \
		echo "🛑 停止服务: make stop"; \
		echo "🏥 健康检查: make health"; \
	else \
		echo "❌ 服务器启动失败，请查看日志: $(LOG_FILE)"; \
		rm -f $(PID_FILE); \
		echo "📋 最近日志:"; \
		tail -20 $(LOG_FILE) 2>/dev/null || echo "无法读取日志文件"; \
		exit 1; \
	fi

# 前台运行开发服务器
.PHONY: run-fg
run-fg: ## 前台运行开发服务器
	@echo "🚀 启动 OKR 服务器（前台运行）..."
	$(GO_CMD) run $(SERVER_PATH)/*.go

# 停止服务器
.PHONY: stop
stop: ## 停止后台运行的服务器
	@echo "🛑 停止 OKR 服务器..."
	@# 首先尝试通过PID文件停止
	@if [ -f $(PID_FILE) ]; then \
		PID=`cat $(PID_FILE)`; \
		if ps -p $$PID > /dev/null 2>&1; then \
			echo "   停止主进程 (PID: $$PID)..."; \
			kill -TERM $$PID 2>/dev/null || true; \
			sleep 3; \
			if ps -p $$PID > /dev/null 2>&1; then \
				echo "   强制停止主进程..."; \
				kill -9 $$PID 2>/dev/null || true; \
			fi; \
		fi; \
		rm -f $(PID_FILE); \
	fi
	@# 停止所有相关的go run进程
	@PIDS=$$(pgrep -f "go run.*cmd/server" 2>/dev/null || true); \
	if [ -n "$$PIDS" ]; then \
		echo "   停止 go run 进程: $$PIDS"; \
		echo $$PIDS | xargs kill -TERM 2>/dev/null || true; \
		sleep 2; \
		PIDS=$$(pgrep -f "go run.*cmd/server" 2>/dev/null || true); \
		if [ -n "$$PIDS" ]; then \
			echo "   强制停止 go run 进程: $$PIDS"; \
			echo $$PIDS | xargs kill -9 2>/dev/null || true; \
		fi; \
	fi
	@# 停止所有Go临时可执行文件进程（通过命令行特征识别）
	@PIDS=$$(pgrep -f "/tmp/go-build.*exe/main" 2>/dev/null || true); \
	if [ -n "$$PIDS" ]; then \
		echo "   停止 Go 临时进程: $$PIDS"; \
		echo $$PIDS | xargs kill -TERM 2>/dev/null || true; \
		sleep 2; \
		PIDS=$$(pgrep -f "/tmp/go-build.*exe/main" 2>/dev/null || true); \
		if [ -n "$$PIDS" ]; then \
			echo "   强制停止 Go 临时进程: $$PIDS"; \
			echo $$PIDS | xargs kill -9 2>/dev/null || true; \
		fi; \
	fi
	@# 最后停止所有监听8081端口的进程
	@PIDS=$$(lsof -ti :$(PORT) 2>/dev/null || true); \
	if [ -n "$$PIDS" ]; then \
		echo "   停止监听端口$(PORT)的进程: $$PIDS"; \
		echo $$PIDS | xargs kill -TERM 2>/dev/null || true; \
		sleep 2; \
		PIDS=$$(lsof -ti :$(PORT) 2>/dev/null || true); \
		if [ -n "$$PIDS" ]; then \
			echo "   强制停止监听端口$(PORT)的进程: $$PIDS"; \
			echo $$PIDS | xargs kill -9 2>/dev/null || true; \
		fi; \
	fi
	@echo "✅ 服务器停止完成"

# 强制停止所有相关进程  
.PHONY: force-stop
force-stop: ## 强制停止所有相关的服务器进程
	@echo "💀 强制停止所有 OKR 相关进程..."
	@# 停止所有go run相关进程
	-@pkill -9 -f "go run.*cmd/server" 2>/dev/null
	@# 停止所有Go临时可执行文件进程
	-@pkill -9 -f "/tmp/go-build.*exe/main" 2>/dev/null
	@# 停止所有监听8081端口的进程
	-@lsof -ti :$(PORT) 2>/dev/null | xargs kill -9 2>/dev/null
	@# 清理PID文件
	@rm -f $(PID_FILE)
	@echo "✅ 强制停止完成"

# 重启服务器
.PHONY: restart
restart: stop run ## 重启服务器

# 查看服务器状态
.PHONY: status
status: ## 查看服务器运行状态
	@if [ -f $(PID_FILE) ]; then \
		PID=`cat $(PID_FILE)`; \
		if ps -p $$PID > /dev/null 2>&1; then \
			echo "✅ 服务器正在运行 (PID: $$PID)"; \
			echo "📋 日志文件: $(LOG_FILE)"; \
			echo "🌐 服务地址: http://localhost:$(PORT)"; \
		else \
			echo "❌ 服务器未运行（PID文件存在但进程不存在）"; \
			rm -f $(PID_FILE); \
		fi \
	else \
		echo "❌ 服务器未运行"; \
	fi

# 构建项目
.PHONY: build
build: ## 构建项目
	@echo "🔨 构建项目..."
	$(GO_CMD) build -o bin/$(BINARY_NAME) $(SERVER_PATH)/*.go
	@echo "✅ 构建完成: bin/$(BINARY_NAME)"

# 运行构建后的二进制文件
.PHONY: start
start: build ## 构建并运行服务器
	@echo "🚀 启动服务器..."
	./bin/$(BINARY_NAME)

# 清理构建文件
.PHONY: clean
clean: ## 清理构建文件和日志
	@echo "🧹 清理构建文件和日志..."
	rm -rf bin/
	rm -rf $(LOG_DIR)/
	$(GO_CMD) clean
	@echo "✅ 清理完成"

# 清理日志
.PHONY: clean-logs
clean-logs: ## 清理日志文件
	@echo "🧹 清理日志文件..."
	rm -rf $(LOG_DIR)/*.log
	@echo "✅ 日志清理完成"

# 格式化代码
.PHONY: fmt
fmt: ## 格式化Go代码
	@echo "🎨 格式化代码..."
	$(GO_CMD) fmt ./...
	@echo "✅ 代码格式化完成"

# 代码检查
.PHONY: vet
vet: ## 运行go vet检查
	@echo "🔍 运行代码检查..."
	$(GO_CMD) vet ./...
	@echo "✅ 代码检查完成"

# 运行测试
.PHONY: test
test: ## 运行测试
	@echo "🧪 运行测试..."
	$(GO_CMD) test ./... -v
	@echo "✅ 测试完成"

# 运行测试覆盖率
.PHONY: test-cover
test-cover: ## 运行测试并生成覆盖率报告
	@echo "🧪 运行测试覆盖率..."
	$(GO_CMD) test ./... -coverprofile=coverage.out
	$(GO_CMD) tool cover -html=coverage.out -o coverage.html
	@echo "✅ 测试覆盖率报告生成: coverage.html"

# 下载依赖
.PHONY: deps
deps: ## 下载项目依赖
	@echo "📦 下载依赖..."
	$(GO_CMD) mod download
	$(GO_CMD) mod tidy
	@echo "✅ 依赖下载完成"

# 数据库迁移
.PHONY: migrate
migrate: ## 运行数据库迁移
	@echo "🗃️ 运行数据库迁移..."
	$(GO_CMD) run $(SERVER_PATH)/*.go migrate
	@echo "✅ 数据库迁移完成"

# 生成实体代码
.PHONY: generate
generate: ## 生成ent实体代码
	@echo "⚙️ 生成实体代码..."
	$(GO_CMD) generate ./ent
	@echo "✅ 实体代码生成完成"

# 开发环境设置
.PHONY: dev-setup
dev-setup: deps generate ## 设置开发环境
	@echo "🛠️ 设置开发环境..."
	@echo "✅ 开发环境设置完成"

# 检查所有（格式化、检查、测试）
.PHONY: check
check: fmt vet test ## 运行所有检查（格式化、代码检查、测试）
	@echo "✅ 所有检查完成"

# 创建发布版本
.PHONY: release
release: check build ## 创建发布版本
	@echo "🎉 发布版本创建完成"

# 热重载开发（需要安装air：go install github.com/cosmtrek/air@latest）
.PHONY: dev
dev: ## 使用热重载运行开发服务器（需要air）
	@echo "🔥 启动热重载开发服务器..."
	@which air > /dev/null || (echo "❌ 请先安装air: go install github.com/cosmtrek/air@latest" && exit 1)
	air

# Docker相关命令
.PHONY: docker-build
docker-build: ## 构建Docker镜像
	@echo "🐳 构建Docker镜像..."
	docker build -t okr-web .
	@echo "✅ Docker镜像构建完成"

.PHONY: docker-run
docker-run: ## 运行Docker容器
	@echo "🐳 运行Docker容器..."
	docker run -p $(PORT):$(PORT) okr-web

# Python客户端相关
.PHONY: client-test
client-test: ## 运行Python客户端测试
	@echo "🐍 运行Python客户端测试..."
	cd okr-python-client && python -m pytest tests/ -v
	@echo "✅ Python客户端测试完成"

.PHONY: client-install
client-install: ## 安装Python客户端依赖
	@echo "🐍 安装Python客户端依赖..."
	cd okr-python-client && pip install -r requirements.txt
	@echo "✅ Python客户端依赖安装完成"

# 日志查看
.PHONY: logs
logs: ## 查看服务器日志
	@if [ -f $(LOG_FILE) ]; then \
		echo "📋 查看服务器日志 (按Ctrl+C退出)..."; \
		tail -f $(LOG_FILE); \
	else \
		echo "❌ 日志文件不存在: $(LOG_FILE)"; \
	fi

# 查看最近日志
.PHONY: logs-tail
logs-tail: ## 查看最近的服务器日志
	@if [ -f $(LOG_FILE) ]; then \
		echo "📋 最近的服务器日志:"; \
		tail -n 50 $(LOG_FILE); \
	else \
		echo "❌ 日志文件不存在: $(LOG_FILE)"; \
	fi

# 健康检查
.PHONY: health
health: ## 检查服务器健康状态
	@echo "🏥 检查服务器健康状态..."
	@curl -s http://localhost:$(PORT)/health > /dev/null && echo "✅ 服务器运行正常" || echo "❌ 服务器未响应"

# 完整的开发流程
.PHONY: full-dev
full-dev: clean dev-setup run status ## 完整的开发流程（清理、设置、运行、查看状态）

# 生产环境部署
.PHONY: deploy
deploy: release ## 部署到生产环境
	@echo "🚀 部署到生产环境..."
	# 这里可以添加具体的部署命令
	@echo "✅ 部署完成"
