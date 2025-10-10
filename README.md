# Luna Dial 🌙

基于领域驱动设计（DDD）的个人任务与日志管理系统，支持 OKR 层级任务体系和积极成长导向的自我评估机制。

## 技术栈

- **后端**: Go + GORM + PostgreSQL + DDD 分层架构
- **前端**: React 19 + TypeScript + Vite + TailwindCSS + Zustand + React Query
- **部署**: Docker Compose + Caddy (反向代理)

## 核心特性

### 1. OKR 层级任务体系
任务支持树形结构，从上向下逐层拆解：年度目标 → 季度目标 → 月度目标 → 周目标 → 日任务。每个层级可独立管理，也可关联追溯到上级目标。

### 2. 积极状态设计
摒弃传统的"失败"或"未完成"状态，任务状态只有：
- **未开始**: 计划中的任务
- **进行中**: 正在执行的任务
- **已完成**: 达成目标的任务
- **已取消**: 因情况变化而不再需要的任务

这种设计鼓励用户关注过程而非结果，避免因任务积压产生焦虑。

### 3. 努力分数评估系统
每日任务可设置"努力分数"（0-10分），用于记录当天在该任务上的投入程度，而非完成度。即使任务未完成，高努力分数也能反映真实付出，系统会统计并展示长期的努力趋势。

### 4. 多维度日志系统
支持为不同时间周期（日/周/月/季/年）记录多条日志：
- **展望日志**: 期初规划与目标设定
- **过程日志**: 中期回顾与调整
- **总结日志**: 期末复盘与经验沉淀

### 5. DDD 分层架构
后端采用严格的 DDD 分层设计：
- **领域层**（Business Layer）: 核心业务逻辑与领域模型
- **服务层**（Service Layer）: HTTP 请求处理与数据校验
- **数据层**（Data Layer）: 仓储实现与数据持久化
- **路由层**（Server Layer）: API 路由与请求分发

## 部署方法

### Docker 部署（推荐）

适用于生产环境或快速体验，一键启动完整服务：

```bash
# 在项目根目录执行
cd docker
docker compose up -d

# 查看服务状态
docker compose ps

# 查看日志
docker compose logs -f
```

服务启动后访问 `http://localhost:10755`

**服务说明**:
- PostgreSQL 数据库会自动初始化并应用 migrations
- 后端 API 服务端口: 8081 (内部)
- Caddy 反向代理端口: 10755 (对外)
- 数据持久化在 `docker/data/postgres/`

### 本地开发部署

适用于开发调试，前后端分离启动：

#### 1. 启动后端服务

```bash
cd luna_dial_server

# 方式A: 使用 Docker 运行数据库和后端
./start.sh

# 方式B: 本地开发（需先启动 PostgreSQL）
go mod download
go run cmd/main.go

# 运行测试
go test ./...
```

后端服务运行在 `http://localhost:8081`

#### 2. 启动前端服务

```bash
cd luna_dial_frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

前端开发服务器运行在 `http://localhost:5173`

## 项目结构

```
luna-dial/
├── docker/                    # Docker 部署配置
│   ├── docker-compose.yml    # 服务编排配置
│   ├── config/               # 应用配置文件
│   └── data/                 # 数据持久化目录
├── luna_dial_server/         # 后端服务 (Go)
│   ├── cmd/                  # 程序入口
│   ├── internal/
│   │   ├── biz/             # 业务逻辑层
│   │   ├── service/         # 服务层 (HTTP handlers)
│   │   ├── data/            # 数据层 (Repository)
│   │   └── server/          # 路由层
│   ├── migrations/          # 数据库迁移文件
│   └── docs/                # API 文档
└── luna_dial_frontend/       # 前端应用 (React)
    ├── src/
    │   ├── components/      # UI 组件
    │   ├── pages/           # 页面组件
    │   ├── store/           # Zustand 状态管理
    │   └── api/             # API 请求封装
    └── dist/                # 构建输出目录
```

## 许可证

本项目采用 **CC BY-NC 4.0** 许可证 - 详见 [LICENSE](LICENSE) 文件

- 允许个人使用、学习、研究、教育等非商业用途
- 允许自由分享、修改和基于本项目创作
- 使用时需注明原作者和来源
- 禁止用于商业目的（如需商业授权请联系作者）

## 联系方式

- GitHub Issues: [提交问题](https://github.com/zzy2210/luna-dial/issues)
- 项目维护者: [@zzy2210](https://github.com/zzy2210)
