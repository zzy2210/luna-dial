# OKR管理系统 - AI设计规范文档

## 文档元信息
- **项目名称**: OKR年度计划管理系统
- **文档版本**: v1.0
- **创建日期**: 2025-07-10
- **技术栈**: Go + Echo + Ent ORM + PostgreSQL
- **目标**: 为AI实现提供结构化、可执行的设计规范

## 1. 项目结构定义

### 1.1 目录结构
```
okr-web/
├── server/                    # 后端Go项目根目录
│   ├── cmd/                   # 应用程序入口
│   │   └── main.go           # 主程序入口
│   ├── internal/             # 内部代码，不对外暴露
│   │   ├── config/           # 配置管理
│   │   ├── controller/       # 控制器层
│   │   ├── service/          # 服务层
│   │   ├── repository/       # 数据访问层
│   │   ├── middleware/       # 中间件
│   │   ├── types/            # 枚举类型定义
│   │   └── utils/            # 工具函数
│   ├── ent/                  # Ent ORM生成的代码
│   │   ├── schema/           # 数据模型定义
│   │   └── (generated files) # 生成的ORM代码
│   ├── migrations/           # 数据库迁移文件
│   ├── tests/               # 测试文件
│   ├── config.ini           # 配置文件
│   ├── go.mod              # Go模块定义
│   └── go.sum              # 依赖锁定
├── docs/                    # 项目文档
└── README.md               # 项目说明
```

### 1.2 依赖包清单
```go
// 必需的Go依赖
require (
    github.com/labstack/echo/v4           // Web框架
    entgo.io/ent                          // ORM框架
    github.com/google/uuid                // UUID生成
    github.com/lib/pq                     // PostgreSQL驱动
    gopkg.in/ini.v1                       // INI配置文件解析
    golang.org/x/crypto                   // 密码加密
    github.com/golang-jwt/jwt/v5          // JWT认证
    github.com/labstack/gommon/log        // 日志
    github.com/stretchr/testify           // 测试框架
)
```

## 2. 数据模型规范

### 2.1 枚举类型定义
```go
// 必须在 internal/types/ 目录下实现以下枚举类型

// TaskType - 任务类型枚举
type TaskType string
const (
    TaskTypeYear    TaskType = "year"     // 年度任务
    TaskTypeQuarter TaskType = "quarter"  // 季度任务
    TaskTypeMonth   TaskType = "month"    // 月度任务
    TaskTypeWeek    TaskType = "week"     // 周任务
    TaskTypeDay     TaskType = "day"      // 日任务
)

// TaskStatus - 任务状态枚举
type TaskStatus string
const (
    TaskStatusPending    TaskStatus = "pending"     // 待开始
    TaskStatusInProgress TaskStatus = "in-progress" // 进行中
    TaskStatusCompleted  TaskStatus = "completed"   // 已完成
)

// TimeScale - 时间尺度枚举
type TimeScale string
const (
    TimeScaleDay     TimeScale = "day"     // 日
    TimeScaleWeek    TimeScale = "week"    // 周
    TimeScaleMonth   TimeScale = "month"   // 月
    TimeScaleQuarter TimeScale = "quarter" // 季
    TimeScaleYear    TimeScale = "year"    // 年
)

// EntryType - 日志条目类型枚举
type EntryType string
const (
    EntryTypePlanStart     EntryType = "plan-start"     // 开始计划
    EntryTypeReflection    EntryType = "reflection"     // 阶段反思
    EntryTypeSummary       EntryType = "summary"        // 结束总结
)
```

### 2.2 Ent Schema定义要求

#### User Schema
```go
// 位置: server/ent/schema/user.go
// 字段: ID(UUID), Username(string), Password(string), Email(string), CreatedAt, UpdatedAt
// 约束: Username唯一, Email唯一
// 关系: 一对多Tasks, 一对多JournalEntries
```

#### Task Schema
```go
// 位置: server/ent/schema/task.go
// 字段: ID(UUID), Title(string), Description(string), Type(TaskType), 
//       StartDate(time.Time), EndDate(time.Time), Status(TaskStatus), 
//       Score(int), ParentID(*UUID), UserID(UUID), Tags(string), CreatedAt, UpdatedAt
// 约束: Title非空, Type使用枚举, Score范围1-10
// 关系: 多对一User, 自关联Parent/Children, 多对多JournalEntries
// 索引: UserID, ParentID, Type, StartDate, EndDate
```

#### JournalEntry Schema
```go
// 位置: server/ent/schema/journalentry.go
// 字段: ID(UUID), Content(string), TimeReference(string), TimeScale(TimeScale), 
//       EntryType(EntryType), UserID(UUID), CreatedAt, UpdatedAt
// 约束: Content非空, TimeScale使用枚举, EntryType使用枚举
// 关系: 多对一User, 多对多Tasks
// 索引: UserID, TimeScale
```

## 3. API接口规范

### 3.1 响应格式标准
```go
// 成功响应格式
type SuccessResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data"`
    Message string      `json:"message,omitempty"`
}

// 错误响应格式
type ErrorResponse struct {
    Success bool   `json:"success"`
    Error   string `json:"error"`
    Message string `json:"message"`
}

// 分页响应格式
type PaginationResponse struct {
    Success     bool        `json:"success"`
    Data        interface{} `json:"data"`
    Total       int64       `json:"total"`
    CurrentPage int         `json:"current_page"`
    PageSize    int         `json:"page_size"`
    TotalPages  int         `json:"total_pages"`
}
```

### 3.2 API路由分组
```go
// 路由分组定义
API_PREFIX = "/api"

// 认证相关 - /api/auth
POST   /auth/register     // 用户注册
POST   /auth/login        // 用户登录
GET    /auth/me          // 获取当前用户
POST   /auth/logout      // 用户登出

// 任务相关 - /api/tasks (需要JWT认证)
GET    /tasks                    // 获取任务列表(分页+过滤)
GET    /tasks/:id               // 获取单个任务
POST   /tasks                   // 创建任务
PUT    /tasks/:id               // 更新任务
DELETE /tasks/:id               // 删除任务
GET    /tasks/:id/children      // 获取子任务
GET    /tasks/:id/full-tree     // 获取完整任务树
POST   /tasks/:id/sub-task      // 创建子任务
PUT    /tasks/:id/score         // 更新任务分数
GET    /tasks/context-view      // 上下文视图
GET    /tasks/global-view       // 全局树视图

// 日志相关 - /api/journals (需要JWT认证)
GET    /journals              // 获取日志列表
GET    /journals/:id          // 获取日志详情
POST   /journals              // 创建日志
PUT    /journals/:id          // 更新日志
DELETE /journals/:id          // 删除日志
GET    /journals/by-time      // 按时间查询日志

// 统计相关 - /api/stats (需要JWT认证)
GET    /stats/task-completion    // 任务完成统计
GET    /stats/score-trend        // 评分趋势
GET    /stats/time-distribution  // 任务时间分布
```

## 4. 服务层架构规范

### 4.1 Service接口定义
```go
// 每个Service必须定义接口，实现依赖注入

// UserService接口
type UserService interface {
    Register(ctx context.Context, req RegisterRequest) (*User, error)
    Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
    GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
    // ... 其他方法
}

// TaskService接口
type TaskService interface {
    CreateTask(ctx context.Context, userID uuid.UUID, req TaskRequest) (*Task, error)
    GetTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*Task, error)
    UpdateTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, req TaskRequest) (*Task, error)
    DeleteTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) error
    GetTaskChildren(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) ([]*Task, error)
    GetTaskTree(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*TaskTree, error)
    GetGlobalView(ctx context.Context, userID uuid.UUID) ([]*TaskTree, error)
    UpdateTaskScore(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, score int) error
    // ... 其他方法
}

// JournalService接口
type JournalService interface {
    CreateJournal(ctx context.Context, userID uuid.UUID, req JournalRequest) (*JournalEntry, error)
    GetJournal(ctx context.Context, userID uuid.UUID, journalID uuid.UUID) (*JournalEntry, error)
    UpdateJournal(ctx context.Context, userID uuid.UUID, journalID uuid.UUID, req JournalRequest) (*JournalEntry, error)
    DeleteJournal(ctx context.Context, userID uuid.UUID, journalID uuid.UUID) error
    GetJournalsByTime(ctx context.Context, userID uuid.UUID, timeRange TimeRangeRequest) ([]*JournalEntry, error)
    // ... 其他方法
}
```

### 4.2 Repository层模式
```go
// Repository层封装Ent操作，提供统一接口

type UserRepository interface {
    Create(ctx context.Context, user *ent.UserCreate) (*ent.User, error)
    GetByID(ctx context.Context, id uuid.UUID) (*ent.User, error)
    GetByUsername(ctx context.Context, username string) (*ent.User, error)
    Update(ctx context.Context, id uuid.UUID, update *ent.UserUpdateOne) (*ent.User, error)
    Delete(ctx context.Context, id uuid.UUID) error
}

// 类似地定义TaskRepository, JournalRepository
```

## 5. 配置管理规范

### 5.1 配置结构定义
```go
// 位置: internal/config/config.go
type Config struct {
    Server   ServerConfig   `ini:"server"`
    Database DatabaseConfig `ini:"database"`
    Log      LogConfig      `ini:"log"`
    JWT      JWTConfig      `ini:"jwt"`
}

type ServerConfig struct {
    Host string `ini:"host"`
    Port int    `ini:"port"`
}

type DatabaseConfig struct {
    Host     string `ini:"host"`
    Port     int    `ini:"port"`
    User     string `ini:"user"`
    Password string `ini:"password"`
    DBName   string `ini:"dbname"`
    SSLMode  string `ini:"sslmode"`
}

type LogConfig struct {
    Level string `ini:"level"`
}

type JWTConfig struct {
    Secret     string `ini:"secret"`
    ExpiryHour int    `ini:"expiry_hour"`
}
```

### 5.2 环境变量覆盖规则
```
// 环境变量命名规则: SECTION_KEY
SERVER_HOST     -> config.Server.Host
SERVER_PORT     -> config.Server.Port
DATABASE_HOST   -> config.Database.Host
DATABASE_PORT   -> config.Database.Port
// ... 以此类推
```

## 6. 中间件规范

### 6.1 必需中间件
```go
// 中间件执行顺序 (从外到内)
1. Logger中间件      // 请求日志记录
2. CORS中间件        // 跨域处理
3. Recover中间件     // 恐慌恢复
4. JWT中间件         // JWT认证 (仅受保护路由)
5. RateLimit中间件   // 限流 (可选)
```

### 6.2 JWT中间件要求
```go
// JWT Claims结构
type JWTClaims struct {
    UserID   uuid.UUID `json:"user_id"`
    Username string    `json:"username"`
    jwt.RegisteredClaims
}

// 中间件应将用户信息存储到Context中
// 在Handler中通过 c.Get("user_id") 获取当前用户ID
```

## 7. 错误处理规范

### 7.1 错误类型定义
```go
// 位置: internal/types/errors.go
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Type    string `json:"type"`
}

// 预定义错误类型
var (
    ErrUserNotFound     = &AppError{Code: 404, Message: "用户不存在", Type: "USER_NOT_FOUND"}
    ErrTaskNotFound     = &AppError{Code: 404, Message: "任务不存在", Type: "TASK_NOT_FOUND"}
    ErrUnauthorized     = &AppError{Code: 401, Message: "未授权访问", Type: "UNAUTHORIZED"}
    ErrInvalidRequest   = &AppError{Code: 400, Message: "请求参数无效", Type: "INVALID_REQUEST"}
    ErrInternalServer   = &AppError{Code: 500, Message: "服务器内部错误", Type: "INTERNAL_SERVER_ERROR"}
)
```

### 7.2 统一错误处理中间件
```go
// 全局错误处理中间件，统一返回错误格式
func ErrorHandlerMiddleware() echo.MiddlewareFunc {
    return middleware.Recover()
    // 自定义错误处理逻辑
}
```

## 8. 测试规范

### 8.1 测试文件结构
```
tests/
├── unit/              # 单元测试
│   ├── service/       # Service层测试
│   ├── repository/    # Repository层测试
│   └── utils/         # 工具函数测试
├── integration/       # 集成测试
│   ├── api/          # API集成测试
│   └── database/     # 数据库集成测试
└── testdata/         # 测试数据
    ├── fixtures/     # 测试夹具
    └── mocks/        # Mock数据
```

### 8.2 测试覆盖要求
```go
// 每个Service方法必须有对应的单元测试
// 每个API端点必须有集成测试
// 测试覆盖率目标: 80%以上
// 使用testify框架进行断言
// 使用database/sql事务进行数据库测试隔离
```

## 9. 数据库迁移规范

### 9.1 迁移文件命名
```
migrations/
├── 001_create_users_table.up.sql
├── 001_create_users_table.down.sql
├── 002_create_tasks_table.up.sql
├── 002_create_tasks_table.down.sql
└── ...
```

### 9.2 迁移执行策略
```go
// 应用启动时自动执行未执行的迁移
// 提供CLI命令手动执行迁移
// 支持回滚操作
// 迁移操作必须是幂等的
```

## 10. 实现检查清单

### 10.1 Phase 1: 基础设施
- [ ] 项目结构创建
- [ ] Go模块初始化
- [ ] 依赖包安装
- [ ] 配置系统实现
- [ ] 基础中间件实现

### 10.2 Phase 2: 数据层
- [ ] Ent Schema定义
- [ ] 枚举类型实现
- [ ] Repository层实现
- [ ] 数据库连接配置
- [ ] 迁移脚本编写

### 10.3 Phase 3: 业务层
- [ ] Service接口定义
- [ ] Service实现
- [ ] 业务逻辑单元测试
- [ ] 错误处理实现

### 10.4 Phase 4: 接口层
- [ ] Controller实现
- [ ] 路由配置
- [ ] JWT认证实现
- [ ] API集成测试

### 10.5 Phase 5: 完善
- [ ] 日志配置
- [ ] 性能优化
- [ ] 文档完善
- [ ] 部署配置

---

**注意**: 此文档为AI实现指南，包含了所有必要的技术细节和约束条件。实现时必须严格遵循此规范，确保代码质量和一致性。
