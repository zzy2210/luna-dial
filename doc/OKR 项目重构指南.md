# OKR 项目重构指南

## 项目概述

本项目包含：
- Go 语言服务端（RESTful API）
- Python 客户端（从单次CLI重构为胖应用）
- 未来支持 Rust 客户端

## 架构设计

### Go 服务端架构

```
okr-server/
├── cmd/server/main.go              # 入口点
├── internal/
│   ├── config/                     # 配置管理
│   ├── domain/                     # 业务实体(task.go, user.go)
│   ├── handler/                    # HTTP处理器
│   ├── repository/                 # 数据访问层
│   │   ├── sqlite/                 # SQLite实现
│   │   └── postgres/               # PostgreSQL实现
│   ├── service/                    # 业务逻辑层
│   └── middleware/                 # 中间件
├── pkg/                           # 公共包
└── api/openapi.yaml               # API文档
```

### Python 客户端架构

```
okr-client/
├── okr_client/
│   ├── main.py                    # 主入口
│   ├── config/                    # 配置管理
│   ├── core/                      # 核心功能(API客户端、本地存储、同步)
│   ├── models/                    # 数据模型
│   ├── ui/
│   │   ├── cli/                   # 命令行接口
│   │   └── tui/                   # 终端UI(可选)
│   └── services/                  # 业务服务
├── data/                          # 用户数据目录
│   ├── app.db                     # 本地SQLite
│   └── config.toml                # 用户配置
└── requirements.txt
```

## API 设计

### RESTful 接口

```
GET    /api/v1/tasks              # 获取任务列表
POST   /api/v1/tasks              # 创建任务
GET    /api/v1/tasks/{id}         # 获取单个任务
PUT    /api/v1/tasks/{id}         # 更新任务
DELETE /api/v1/tasks/{id}         # 删除任务

GET    /health                    # 健康检查
GET    /ready                     # 就绪检查
```

### 数据模型

**Task 实体**：
- id, title, description
- status (pending/in_progress/completed/cancelled)
- priority (low/medium/high)
- due_date, created_at, updated_at
- user_id

## 数据库设计

### 服务端表结构

```sql
CREATE TABLE tasks (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    priority VARCHAR(10) DEFAULT 'medium',
    due_date DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    user_id VARCHAR(36) NOT NULL
);
```

### 客户端本地表结构

```sql
CREATE TABLE local_tasks (
    -- 基础字段同服务端
    id, title, description, status, priority, 
    due_date, created_at, updated_at, user_id,
    
    -- 同步状态字段
    is_synced BOOLEAN DEFAULT FALSE,
    last_sync_at DATETIME,
    local_changes BOOLEAN DEFAULT FALSE
);

CREATE TABLE sync_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    last_sync_time DATETIME,
    sync_type VARCHAR(20),  -- 'full', 'incremental'
    status VARCHAR(20),     -- 'success', 'failed'
    error_message TEXT
);
```

## Python 客户端 UI 设计

### CLI 界面结构

```python
# 主命令组
okr                           # 进入TUI模式(默认)
okr --mode=cli               # 传统CLI模式

# CLI子命令
okr list                     # 列出任务
okr add "任务标题"            # 添加任务
okr done <task_id>           # 完成任务
okr sync                     # 手动同步
okr status                   # 查看同步状态
```

### TUI 界面设计 (使用 Textual)

### TUI 视图结构总览

```python
# 登录与用户信息
┌─ 登录 ───────────────────────────────┐
│ 请输入用户名和密码                   │
└─────────────────────────────────────┘

┌─ 用户信息 ───────────────────────────┐
│ 用户名: xxx   邮箱: xxx@xxx.com      │
└─────────────────────────────────────┘

# 计划视图（支持今日/本周/本月/本季度/本年/自定义）
┌─ 2025-07-07 ~ 2025-07-13 计划视图 ──────────────┐
│ 📊 统计概览:                                   │
│ • 总任务数: 3   • 总分: 18   • 完成分数: 12    │
│                                               │
│ 📋 任务树:                                     │
│ 年任务A                                       │
│ └─ 季度任务B                                  │
│    └─ 月任务C                                 │
│       └─ 周任务D                              │
│          ├─ [✓] 日任务E1   分数: 5            │
│          └─ [ ] 日任务E2   分数: 3            │
│                                               │
│ 📒 日志:                                       │
│  - 2025-07-08: 完成了文档初稿                 │
│                                               │
│ [N]新建任务 [J]新建日志 [T]任务树 [Q]退出      │
└───────────────────────────────────────────────┘

# 任务树视图
┌─ 任务树 ──────────────────────────────────────┐
│ 年任务A                                       │
│ └─ 季度任务B                                  │
│    └─ 月任务C                                 │
│       └─ 周任务D                              │
│          ├─ 日任务E1                          │
│          └─ 日任务E2                          │
│                                               │
│ [↑/↓]移动 [Enter]详情 [N]新建子任务 [Q]返回    │
└───────────────────────────────────────────────┘

# 任务详情视图
┌─ 任务详情 ─────────────────────────────────────┐
│ 标题: [重构服务端代码                        ] │
│ 描述: [重新设计架构，提升代码质量...          ] │
│ 状态: [○ 待办 ● 进行中 ○ 已完成 ○ 已取消]     │
│ 优先级: [○ 低 ● 中 ○ 高]                     │
│ 截止日期: [2025-07-20                       ] │
│ 父级任务: [月任务C]                           │
│                                              │
│ [保存] [完成] [打分] [删除] [返回]            │
└───────────────────────────────────────────────┘

# 分数趋势视图
┌─ 2025-07-01 ~ 2025-07-31 分数趋势 ─────────────┐
│ 📊 趋势摘要:                                   │
│ • 总分: 48   • 总任务: 12                      │
│ • 平均分: 12.00   • 平均任务数: 3.00           │
│ • 最高分: 20   • 最低分: 5                     │
│                                               │
│ 📈 趋势图:（本月每周）                         │
│  2025-W27 ▓▓▓▓▓▓▓▓▓▓ 20分 (5任务)              │
│  2025-W28 ▓▓▓▓▓▓▓     12分 (3任务)             │
│  2025-W29 ▓▓▓         8分 (2任务)              │
│  2025-W30 ▓▓▓▓        8分 (2任务)              │
│                                               │
│ [←/→]切换到每日/季度 [Q]退出                   │
└───────────────────────────────────────────────┘

# 日志管理视图
┌─ 日志列表 ────────────────────────────────────┐
│ 2025-07-08: 完成了文档初稿   [E]编辑 [D]删除   │
│ 2025-07-10: 服务端重构遇到问题 [E]编辑 [D]删除 │
│                                               │
│ [N]新建日志 [Q]返回                            │
└───────────────────────────────────────────────┘

# 设置与帮助
┌─ 设置 ────────────────────────────────────────┐
│ API地址: http://localhost:8081/api            │
│ 用户配置路径: ~/.okr/config                   │
│                                               │
│ [E]编辑配置 [Q]返回                           │
└───────────────────────────────────────────────┘

┌─ 帮助 ────────────────────────────────────────┐
│ F1: 帮助  F2: 新建  F3: 编辑  F4: 删除 ...    │
└───────────────────────────────────────────────┘
```

### 关键组件实现


## 同步机制设计

### 同步策略
1. **启动时同步**: 应用启动时检查并同步
2. **定时同步**: 每5分钟自动同步一次
3. **操作后同步**: 用户操作后延迟同步
4. **手动同步**: 用户主动触发同步

### 冲突解决
- **时间戳比较**: 以最新修改时间为准
- **用户选择**: 冲突时提示用户选择
- **备份策略**: 冲突数据保存到备份表

## 重构执行计划

### Phase 1: 基础架构 (1-2周)
- Go服务端重构：清理架构，实现核心API
- Python客户端：建立新架构，实现本地存储

### Phase 2: 核心功能 (3-4周)  
- 完善CRUD操作
- 实现同步机制
- 添加TUI界面

### Phase 3: 优化提升 (5-6周)
- 性能优化
- 用户体验提升
- 错误处理完善

### Phase 4: 测试部署 (7-8周)
- 集成测试
- 文档编写
- 部署准备

## 关键注意事项

1. **数据兼容性**: 确保现有数据可以平滑迁移
2. **向后兼容**: 保持API接口稳定
3. **错误处理**: 统一错误码和用户友好的提示
4. **性能考虑**: 数据库查询优化，网络请求优化
5. **离线支持**: 确保客户端在离线时仍可基本使用

## 技术栈选择

### Go 服务端
- Echo/Gin (HTTP框架)
- GORM/原生SQL (数据库)
- Viper (配置管理)
- Logrus/Zap (日志)

### Python 客户端  
- Click (CLI框架)
- Textual (TUI框架)
- SQLAlchemy (本地数据库)
- aiohttp (HTTP客户端)
- Pydantic (数据验证)

这个重构将显著提升项目的可维护性和用户体验，为未来的多语言客户端支持奠定基础。