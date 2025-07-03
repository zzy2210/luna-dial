# OKR & 日志管理系统 API

基于Go和Echo框架实现的RESTful API系统，支持OKR（目标和关键结果）管理以及日志记录功能。

## 功能特性

- **计划管理**: 支持创建、查询、更新、删除计划，包括时间容器（年/季/月/周/日）和原子任务
- **周期计划**: 管理周期性任务模板，支持每日、每周、每月频率
- **日志记录**: 支持各种类型的日志条目（日记、周总结、月总结等）
- **聚合视图**: 提供指定日期的完整视图，聚合所有相关计划和日志

## 项目结构

```
.
├── cmd/server/main.go          # 应用入口点
├── internal/
│   ├── domain/                 # 领域模型层
│   │   ├── journal_entry.go    # 日志条目实体
│   │   ├── plan.go            # 计划实体
│   │   ├── recurring_plan.go   # 周期计划实体
│   │   └── repository.go       # 仓储接口
│   ├── dto/                    # 数据传输对象
│   │   └── dto.go
│   ├── handler/                # HTTP处理器层
│   │   ├── daily_view_handler.go
│   │   ├── journal_handler.go
│   │   ├── plan_handler.go
│   │   └── recurring_plan_handler.go
│   ├── mock/                   # 模拟数据层
│   │   ├── journal_repository.go
│   │   ├── plan_repository.go
│   │   └── recurring_plan_repository.go
│   └── service/                # 业务服务层
│       └── daily_view_service.go
├── docs/                       # 文档
│   ├── plan.md
│   └── 实际设计文档.md
├── go.mod
└── go.sum
```

## 快速开始

### 运行服务器

```bash
go run cmd/server/main.go
```

服务器将在 `http://localhost:8080` 启动。

### 健康检查

```bash
curl http://localhost:8080/health
```

## API 端点

### 计划 API

- `POST /api/v1/plans` - 创建计划
- `GET /api/v1/plans` - 获取计划列表（支持parentID筛选和分页）
- `GET /api/v1/plans/{id}` - 获取计划详情
- `PUT /api/v1/plans/{id}` - 更新计划
- `DELETE /api/v1/plans/{id}` - 删除计划

### 周期计划 API

- `POST /api/v1/recurring-plans` - 创建周期计划
- `GET /api/v1/recurring-plans` - 获取周期计划列表
- `PUT /api/v1/recurring-plans/{id}` - 更新周期计划
- `DELETE /api/v1/recurring-plans/{id}` - 删除周期计划

### 日志 API

- `POST /api/v1/journal-entries` - 创建日志条目
- `GET /api/v1/journal-entries` - 获取日志列表
- `PUT /api/v1/journal-entries/{id}` - 更新日志
- `DELETE /api/v1/journal-entries/{id}` - 删除日志

### 聚合视图 API

- `GET /api/v1/daily-view?date=2025-07-03` - 获取指定日期的完整视图

## API 使用示例

### 创建计划

```bash
curl -X POST http://localhost:8080/api/v1/plans \
  -H "Content-Type: application/json" \
  -d '{
    "objective": "学习 Go 语言",
    "planType": "Daily",
    "parentID": null
  }'
```

### 创建周期计划

```bash
curl -X POST http://localhost:8080/api/v1/recurring-plans \
  -H "Content-Type: application/json" \
  -d '{
    "objective": "每日冥想10分钟",
    "frequency": "Daily",
    "activePeriod": {
      "startDate": "2025-01-01T00:00:00Z",
      "endDate": "2025-12-31T23:59:59Z"
    }
  }'
```

### 创建日志条目

```bash
curl -X POST http://localhost:8080/api/v1/journal-entries \
  -H "Content-Type: application/json" \
  -d '{
    "content": "今天完成了 API 设计，状态很好。",
    "entryType": "Daily",
    "period": {
      "planType": "Daily",
      "date": "2025-07-03"
    }
  }'
```

### 获取日常视图

```bash
curl "http://localhost:8080/api/v1/daily-view?date=2025-07-03"
```

## 开发说明

项目使用内存中的模拟仓储实现，已预填充一些测试数据：

- 时间容器：2025年 -> 2025年7月 -> 2025-07-03
- 原子任务：完成API设计文档、三分化训练
- 周期计划：每日冥想、每周总结、每日站会总结
- 日志条目：日记、周总结、月总结

## 技术栈

- **Go 1.23+**: 编程语言
- **Echo v4**: Web框架
- **UUID**: 唯一标识符生成
- **领域驱动设计**: 架构模式

## 下一步开发

- 集成真实数据库（PostgreSQL/MySQL）
- 添加用户认证和授权
- 实现数据验证
- 添加单元测试和集成测试
- API文档生成（Swagger）
