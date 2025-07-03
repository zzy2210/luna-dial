# OKR管理系统 - AI设计规范文档 (补充)

本文档是对主设计文档 `ai-design-specification.md` 的补充，旨在定义新增的“计划视图”和“分数统计”功能。

## 1. 新增功能：计划视图 (Plan View)

### 1.1 功能描述
“计划视图”功能旨在为用户提供一个特定时间周期（如年、季、月）的综合视图。该视图将展示该周期内的所有任务、这些任务的完整父级链（递归至顶层），以及与该周期相关的日志条目。

### 1.2 API 路由补充

在 `API路由分组` 中增加一个新的分组：

```go
// ... (原有路由)

// 计划视图相关 - /api/plan (需要JWT认证)
GET    /plan                     // 获取计划视图 (查询参数: scale, time_ref)

// 统计相关 - /api/stats (需要JWT认证)
// ... (原有统计路由)
```

- **GET /api/plan**: 通过查询参数 `scale` (e.g., `quarter`, `month`) 和 `time_ref` (e.g., `2024-Q4`, `2025-07`) 来获取指定周期的计划数据。

### 1.3 服务层架构补充

#### 1.3.1 新增数据传输对象 (DTOs)

```go
// PlanRequest 定义了获取计划视图的请求参数
type PlanRequest struct {
    Scale   TimeScale `query:"scale"`    // 时间尺度 (e.g., "quarter", "month")
    TimeRef string    `query:"time_ref"` // 时间参考 (e.g., "2024-Q4", "2025-07")
}

// PlanResponse 定义了计划视图的响应数据
type PlanResponse struct {
    Tasks    []*TaskTree     `json:"tasks"`    // 包含完整父级链的任务树列表
    Journals []*JournalEntry `json:"journals"` // 相关时间周期的日志列表
}
```

#### 1.3.2 更新 Service 接口

在 `TaskService` 接口中增加新方法：

```go
// TaskService接口
type TaskService interface {
    // ... (原有方法)
    GetPlanView(ctx context.Context, userID uuid.UUID, req PlanRequest) (*PlanResponse, error)
}
```

## 2. 新增功能：分数趋势统计 (Score Trend Statistics)

### 2.1 功能描述
“分数趋势统计”功能用于按指定的时间粒度（如日、月、季）聚合和展示任务的分数总和与任务数量。例如，查询2025年7月的月度趋势，将返回该月每一天的任务总分和任务总数。

### 2.2 API 路由补充

明确 `/api/stats/score-trend` 路由的功能和参数：

```go
// ... (原有路由)

// 统计相关 - /api/stats (需要JWT认证)
GET    /stats/task-completion    // 任务完成统计
GET    /stats/score-trend        // 评分趋势统计 (查询参数: scale, time_ref)
GET    /stats/time-distribution  // 任务时间分布
```

- **GET /api/stats/score-trend**: 通过查询参数 `scale` (e.g., `month`, `quarter`) 和 `time_ref` (e.g., `2025-07`, `2024-Q4`) 来获取分数趋势数据。

### 2.3 服务层架构补充

#### 2.3.1 新增数据传输对象 (DTOs)

```go
// ScoreTrendRequest 定义了获取分数趋势的请求参数
type ScoreTrendRequest struct {
    Scale   TimeScale `query:"scale"`    // 统计尺度 (e.g., "month", "quarter")
    TimeRef string    `query:"time_ref"` // 时间参考 (e.g., "2025-07", "2024-Q4")
}

// ScoreTrendResponse 定义了分数趋势的响应数据
type ScoreTrendResponse struct {
    Labels []string `json:"labels"` // 时间标签 (e.g., ["2025-07-01", "2025-07-02"])
    Scores []int    `json:"scores"` // 对应标签的分数总和
    Counts []int    `json:"counts"` // 对应标签的任务数量
}
```

#### 2.3.2 新增 Service 接口

建议创建一个新的 `StatsService` 来处理所有统计相关的业务逻辑。

```go
// StatsService接口
type StatsService interface {
    GetScoreTrend(ctx context.Context, userID uuid.UUID, req ScoreTrendRequest) (*ScoreTrendResponse, error)
    // ... 其他统计相关方法
}
```
---
