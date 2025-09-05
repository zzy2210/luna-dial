# Luna Dial Server API 文档

## 概述

Luna Dial Server 是一个任务和日志管理系统的后端服务，提供用户认证、任务管理、日志记录和计划管理功能。

**服务地址**: `http://localhost:8081`  
**API 版本**: v1  
**认证方式**: Session-based Authentication  

---

## 认证说明

本 API 使用基于 Session 的认证机制：

1. **登录**: 通过 `/api/v1/public/auth/login` 获取 Session ID
2. **受保护的接口**: 需要在请求头中包含 `Authorization: Bearer <session_id>`
3. **登出**: 通过 `/api/v1/auth/logout` 或 `/api/v1/auth/logout-all` 终止 Session

**认证格式**:
```
Authorization: Bearer <session_id>
```

**示例**:
```
Authorization: Bearer 9e936d7b20c034cad9ca192c108a7ae45a0bc40df9256d87a6bed145f47e5f62
```

---

## API 端点

### 🔓 公开接口

#### 1. 健康检查

```http
GET /health
```

**描述**: 检查服务运行状态

**响应**:
```
200 OK
Content-Type: text/plain

Service is running
```

#### 2. 版本信息

```http
GET /version
```

**描述**: 获取服务版本信息

**响应**:
```
200 OK
Content-Type: text/plain

Version 1.0.0
```

#### 3. 用户登录

```http
POST /api/v1/public/auth/login
```

**描述**: 用户登录，获取 Session

**请求体**:
```json
{
  "username": "string",
  "password": "string"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "Login successful",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "session_id": "9e936d7b20c034cad9ca192c108a7ae45a0bc40df9256d87a6bed145f47e5f62",
    "expires_in": 86400
  }
}
```

**字段说明**:
- `session_id`: 会话ID，用于后续认证
- `expires_in`: 会话过期时间（秒）

**错误响应**:
```json
{
  "code": 401,
  "message": "Invalid username or password",
  "success": false,
  "timestamp": 1691234567
}
```

---

### 🔒 受保护接口

> **注意**: 以下接口需要在请求头中包含有效的 Session 信息

#### 认证管理

##### 1. 获取用户资料

```http
GET /api/v1/auth/profile
```

**描述**: 获取当前登录用户的基本信息

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "user_id": "user_456",
    "username": "john_doe",
    "name": "John Doe",
    "email": "john.doe@example.com"
  }
}
```

##### 2. 获取当前用户详细信息

```http
GET /api/v1/users/me
```

**描述**: 获取当前登录用户的详细信息，包含会话信息

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "user_id": "user_456",
    "username": "john_doe",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "created_at": "2023-08-01T10:30:00Z",
    "updated_at": "2023-08-05T15:45:00Z",
    "session": {
      "session_id": "9e936d7b20c034cad9ca192c108a7ae45a0bc40df9256d87a6bed145f47e5f62",
      "last_access_at": "2023-08-05T16:20:00Z",
      "expires_at": "2023-08-06T10:30:00Z"
    }
  }
}
```

##### 3. 用户登出

```http
POST /api/v1/auth/logout
```

**描述**: 登出当前 Session

**响应**:
```json
{
  "code": 200,
  "message": "Logout successful",
  "success": true,
  "timestamp": 1691234567
}
```

##### 3. 登出所有设备

```http
DELETE /api/v1/auth/logout-all
```

**描述**: 登出该用户的所有 Session

**响应**:
```json
{
  "code": 200,
  "message": "All sessions logged out",
  "success": true,
  "timestamp": 1691234567
}
```

#### 日志管理

##### 1. 获取日志列表（按时间周期）

```http
GET /api/v1/journals
```

**描述**: 根据指定的时间周期类型和时间范围获取日志列表

**请求体**:
```json
{
  "period_type": "day|week|month|quarter|year",
  "start_date": "2023-08-05T00:00:00Z",
  "end_date": "2023-08-06T00:00:00Z"
}
```

**查询参数说明**:
- `period_type` (string, 必填): 时间周期类型
  - `day`: 日志，时间范围必须是完整的一天（00:00:00 到次日 00:00:00）
  - `week`: 周志，时间范围必须是完整的一周（周一 00:00:00 到下周一 00:00:00）
  - `month`: 月志，时间范围必须是完整的一个月
  - `quarter`: 季志，时间范围必须是完整的一个季度
  - `year`: 年志，时间范围必须是完整的一年
- `start_date` (string, 必填): 开始时间，ISO 8601 格式
- `end_date` (string, 必填): 结束时间，ISO 8601 格式，必须大于开始时间

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "success": true,
  "timestamp": 1691234567,
  "data": [
    {
      "id": "journal_123",
      "title": "每日总结",
      "content": "今天完成了项目需求分析...",
      "journal_type": "day",
      "time_period": {
        "start": "2023-08-05T00:00:00Z",
        "end": "2023-08-06T00:00:00Z"
      },
      "icon": "📝",
      "created_at": "2023-08-05T10:30:00Z",
      "updated_at": "2023-08-05T15:45:00Z",
      "user_id": "user_456"
    }
  ]
}
```

##### 2. 创建日志

```http
POST /api/v1/journals
```

**描述**: 创建新的日志条目

**请求体**:
```json
{
  "title": "每日工作总结",
  "content": "今天主要完成了以下工作：\n1. 完成了API文档编写\n2. 修复了数据库连接问题",
  "journal_type": "day",
  "start_date": "2023-08-05T00:00:00Z",
  "end_date": "2023-08-06T00:00:00Z",
  "icon": "📝"
}
```

**字段说明**:
- `title` (string, 必填): 日志标题
- `content` (string, 必填): 日志内容
- `journal_type` (string, 必填): 日志类型 (`day`|`week`|`month`|`quarter`|`year`)
- `start_date` (string, 必填): 日志时间段开始时间
- `end_date` (string, 必填): 日志时间段结束时间
- `icon` (string, 可选): 日志图标

**响应**:
```json
{
  "code": 201,
  "message": "success",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "id": "journal_123",
    "title": "每日工作总结",
    "content": "今天主要完成了以下工作：\n1. 完成了API文档编写\n2. 修复了数据库连接问题",
    "journal_type": "day",
    "time_period": {
      "start": "2023-08-05T00:00:00Z",
      "end": "2023-08-06T00:00:00Z"
    },
    "icon": "📝",
    "created_at": "2023-08-05T10:30:00Z",
    "updated_at": "2023-08-05T10:30:00Z",
    "user_id": "user_456"
  }
}
```

##### 3. 更新日志

```http
PUT /api/v1/journals/{journal_id}
```

**描述**: 更新指定的日志条目

**路径参数**:
- `journal_id` (string): 日志 ID

**请求体** (所有字段均为可选):
```json
{
  "journal_id": "journal_123",
  "title": "更新后的标题",
  "content": "更新后的内容",
  "journal_type": "day",
  "icon": "📖"
}
```

**字段说明**:
- `journal_id` (string, 必填): 日志 ID（在请求体中）
- `title` (string, 可选): 新的日志标题
- `content` (string, 可选): 新的日志内容  
- `journal_type` (string, 可选): 新的日志类型
- `icon` (string, 可选): 新的日志图标

**注意**: 至少需要提供一个要更新的字段

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "id": "journal_123",
    "title": "更新后的标题",
    "content": "更新后的内容",
    "journal_type": "day",
    "time_period": {
      "start": "2023-08-05T00:00:00Z",
      "end": "2023-08-06T00:00:00Z"
    },
    "icon": "📖",
    "created_at": "2023-08-05T10:30:00Z",
    "updated_at": "2023-08-05T16:20:00Z",
    "user_id": "user_456"
  }
}
```

##### 4. 删除日志

```http
DELETE /api/v1/journals/{journal_id}
```

**描述**: 删除指定的日志条目

**路径参数**:
- `journal_id` (string): 日志 ID

**响应**:
```
HTTP/1.1 204 No Content
```

#### 任务管理

##### 1. 获取任务列表（按时间周期）

```http
GET /api/v1/tasks
```

**描述**: 根据指定的时间周期类型和时间范围获取任务列表

**请求体**:
```json
{
  "period_type": "day|week|month|quarter|year",
  "start_date": "2023-08-05T00:00:00Z",
  "end_date": "2023-08-12T00:00:00Z"
}
```

**查询参数说明**:
- `period_type` (string, 必填): 时间周期类型，支持的值：
  - `day`: 获取指定时间范围内的日任务
  - `week`: 获取指定时间范围内的周任务
  - `month`: 获取指定时间范围内的月任务
  - `quarter`: 获取指定时间范围内的季度任务
  - `year`: 获取指定时间范围内的年任务
- `start_date` (string, 必填): 开始时间，ISO 8601 格式
- `end_date` (string, 必填): 结束时间，ISO 8601 格式

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "success": true,
  "timestamp": 1691234567,
  "data": [
    {
      "id": "task_123",
      "title": "完成API文档编写",
      "type": "day",
      "period": {
        "start": "2023-08-05T09:00:00Z",
        "end": "2023-08-05T18:00:00Z"
      },
      "tags": ["开发", "文档"],
      "icon": "📝",
      "score": 85,
      "status": 2,
      "priority": 2,
      "parent_id": "",
      "user_id": "user_456",
      "created_at": "2023-08-05T08:00:00Z",
      "updated_at": "2023-08-05T17:30:00Z"
    }
  ]
}
```

**状态码说明**:
- `status`: 0=未开始, 1=进行中, 2=已完成, 3=已取消
- `priority`: 0=低, 1=中, 2=高, 3=紧急

##### 2. 创建任务

```http
POST /api/v1/tasks
```

**描述**: 创建新任务

**请求体**:
```json
{
  "title": "完成项目文档",
  "description": "编写 API 文档和用户手册",
  "start_date": "2023-08-05T09:00:00Z",
  "end_date": "2023-08-10T18:00:00Z",
  "priority": "high",
  "icon": "📝",
  "tags": ["开发", "文档"]
}
```

**字段说明**:
- `title` (string, 必填): 任务标题
- `description` (string, 可选): 任务描述
- `start_date` (string, 必填): 任务开始时间
- `end_date` (string, 必填): 任务结束时间
- `priority` (string, 必填): 优先级 (`low`|`medium`|`high`|`urgent`)
- `icon` (string, 可选): 任务图标（emoji）
- `tags` (array, 可选): 任务标签数组

**响应**:
```json
{
  "code": 200,
  "message": "create task endpoint",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "id": "task_123",
    "title": "完成项目文档",
    "type": "day",
    "period": {
      "start": "2023-08-05T09:00:00Z",
      "end": "2023-08-10T18:00:00Z"
    },
    "tags": ["开发", "文档"],
    "icon": "📝",
    "score": 0,
    "status": 0,
    "priority": 2,
    "parent_id": "",
    "user_id": "user_456",
    "created_at": "2023-08-05T10:30:00Z",
    "updated_at": "2023-08-05T10:30:00Z"
  }
}
```

##### 3. 更新任务

```http
PUT /api/v1/tasks/{task_id}
```

**描述**: 更新指定任务

**路径参数**:
- `task_id` (string): 任务 ID

**请求体** (所有字段均为可选):
```json
{
  "title": "更新后的任务标题",
  "description": "更新后的任务描述",
  "start_date": "2023-08-05T09:00:00Z",
  "end_date": "2023-08-10T18:00:00Z",
  "priority": "urgent",
  "status": "in_progress",
  "icon": "🚀",
  "tags": ["开发", "文档", "紧急"]
}
```

**字段说明**:
- 任务 ID 从路径参数 `task_id` 提供
- `title` (string, 可选): 新的任务标题
- `description` (string, 可选): 新的任务描述
- `start_date` (string, 可选): 新的开始时间
- `end_date` (string, 可选): 新的结束时间（与start_date必须同时提供）
- `priority` (string, 可选): 新的优先级
- `status` (string, 可选): 新的状态 (`not_started`|`in_progress`|`completed`|`cancelled`)
- `icon` (string, 可选): 新的图标
- `tags` (array, 可选): 新的标签数组

**响应**:
```json
{
  "code": 200,
  "message": "update task endpoint",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "id": "task_123",
    "title": "更新后的任务标题",
    "type": "day",
    "period": {
      "start": "2023-08-05T09:00:00Z",
      "end": "2023-08-10T18:00:00Z"
    },
    "tags": ["开发", "文档", "紧急"],
    "icon": "🚀",
    "score": 0,
    "status": 1,
    "priority": 3,
    "parent_id": "",
    "user_id": "user_456",
    "created_at": "2023-08-05T10:30:00Z",
    "updated_at": "2023-08-05T16:20:00Z"
  }
}
```

##### 4. 删除任务

```http
DELETE /api/v1/tasks/{task_id}
```

**描述**: 删除指定任务

**路径参数**:
- `task_id` (string): 任务 ID

**请求体**:
```json
{
  "task_id": "task_123"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "delete task endpoint",
  "success": true,
  "timestamp": 1691234567
}
```

##### 5. 完成任务

```http
POST /api/v1/tasks/{task_id}/complete
```

**描述**: 标记任务为已完成

**路径参数**:
- `task_id` (string): 任务 ID

**请求体**:
```json
{
  "task_id": "task_123"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "complete task endpoint",
  "success": true,
  "timestamp": 1691234567
}
```

##### 6. 创建子任务

```http
POST /api/v1/tasks/{task_id}/subtasks
```

**描述**: 为指定任务创建子任务

**路径参数**:
- `task_id` (string): 父任务 ID

**请求体**:
```json
{
  "title": "子任务标题",
  "description": "子任务描述",
  "start_date": "2023-08-06T09:00:00Z",
  "end_date": "2023-08-06T18:00:00Z",
  "priority": "medium",
  "icon": "📋",
  "tags": ["子任务"]
}
```

**字段说明**:
- 父任务 ID 从路径参数 `task_id` 提供
- 其他字段与创建任务相同

**响应**:
```json
{
  "code": 200,
  "message": "create subtask endpoint",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "id": "subtask_456",
    "title": "子任务标题",
    "type": "day",
    "period": {
      "start": "2023-08-06T09:00:00Z",
      "end": "2023-08-06T18:00:00Z"
    },
    "tags": ["子任务"],
    "icon": "📋",
    "score": 0,
    "status": 0,
    "priority": 1,
    "parent_id": "parent_task_123",
    "user_id": "user_456",
    "created_at": "2023-08-05T11:00:00Z",
    "updated_at": "2023-08-05T11:00:00Z"
  }
}
```

##### 7. 更新任务评分

```http
PUT /api/v1/tasks/{task_id}/score
```

**描述**: 更新任务的完成评分

**路径参数**:
- `task_id` (string): 任务 ID

**请求体**:
```json
{
  "score": 85
}
```

**字段说明**:
- `score` (int, 必填): 评分（非负整数）

**响应**:
```json
{
  "code": 200,
  "message": "update task score endpoint",
  "success": true,
  "timestamp": 1691234567
}
```

#### 🆕 任务树优化API（阶段五新增）

##### 9. 分页查询根任务

```http
GET /api/v1/tasks/roots
```

**描述**: 分页获取用户的根任务列表（无父任务的任务）

**查询参数**:
- `page` (int, 可选): 页码，默认为1
- `page_size` (int, 可选): 每页大小，默认为20，最大100
- `status` (string[], 可选): 状态过滤，可选值：`not_started`, `in_progress`, `completed`, `cancelled`
- `priority` (string[], 可选): 优先级过滤，可选值：`low`, `medium`, `high`, `urgent`
- `task_type` (string[], 可选): 任务类型过滤，可选值：`day`, `week`, `month`, `quarter`, `year`

**请求示例**:
```http
GET /api/v1/tasks/roots?page=1&page_size=10&status=not_started&status=in_progress
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "items": [
      {
        "id": "task_123",
        "title": "2024年度目标",
        "type": "year",
        "has_children": true,
        "children_count": 4,
        "tree_depth": 0,
        "root_task_id": "task_123",
        "status": 1,
        "priority": 2,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 10,
      "total": 25,
      "total_pages": 3,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

##### 10. 全局任务树视图（分页）

```http
GET /api/v1/tasks/tree
```

**描述**: 分页获取用户的完整任务树，返回嵌套的树形结构

**查询参数**:
- `page` (int, 可选): 页码，默认为1
- `page_size` (int, 可选): 每页大小，默认为10，最大50
- `status` (string[], 可选): 状态过滤
- `include_empty` (bool, 可选): 是否包含无子任务的根任务，默认true

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "items": [
      {
        "id": "task_123",
        "title": "2024年度目标",
        "type": "year",
        "has_children": true,
        "children_count": 2,
        "tree_depth": 0,
        "children": [
          {
            "id": "task_124",
            "title": "Q1目标",
            "type": "quarter",
            "parent_id": "task_123",
            "tree_depth": 1,
            "has_children": true,
            "children_count": 1,
            "children": [
              {
                "id": "task_125",
                "title": "1月任务",
                "type": "month",
                "parent_id": "task_124",
                "tree_depth": 2,
                "has_children": false,
                "children_count": 0,
                "children": []
              }
            ]
          }
        ]
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 10,
      "total": 5,
      "total_pages": 1,
      "has_next": false,
      "has_prev": false
    }
  }
}
```

##### 11. 获取指定任务的完整任务树

```http
GET /api/v1/tasks/{task_id}/tree
```

**描述**: 获取指定任务及其所有子任务的完整树形结构

**路径参数**:
- `task_id` (string): 任务ID

**查询参数**:
- `status` (string, 可选): 状态过滤

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "id": "task_123",
    "title": "2024年度目标",
    "type": "year",
    "has_children": true,
    "children_count": 2,
    "children": [
      {
        "id": "task_124",
        "title": "Q1目标",
        "parent_id": "task_123",
        "children": [...]
      }
    ]
  }
}
```

##### 12. 获取任务的父任务链

```http
GET /api/v1/tasks/{task_id}/parents
```

**描述**: 获取指定任务的所有父任务，从根任务到直接父任务的链路

**路径参数**:
- `task_id` (string): 任务ID

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "success": true,
  "timestamp": 1691234567,
  "data": [
    {
      "id": "task_123",
      "title": "2024年度目标",
      "type": "year",
      "tree_depth": 0
    },
    {
      "id": "task_124", 
      "title": "Q1目标",
      "type": "quarter",
      "tree_depth": 1,
      "parent_id": "task_123"
    }
  ]
}
```

##### 13. 移动任务

```http
PUT /api/v1/tasks/{task_id}/move
```

**描述**: 移动任务到新的父任务下或移动到根级别

**路径参数**:
- `task_id` (string): 要移动的任务ID

**请求体**:
```json
{
  "task_id": "task_125",
  "new_parent_id": "task_456"
}
```

**字段说明**:
- `task_id` (string, 必填): 要移动的任务ID
- `new_parent_id` (string, 可选): 新的父任务ID，为空或null表示移动到根级别

**响应**:
```json
{
  "code": 501,
  "message": "Task move functionality is not yet implemented",
  "success": false,
  "timestamp": 1691234567
}
```

**注意**: 任务移动功能的业务逻辑尚未实现，将在后续版本中提供。

##### 14. 优化的任务创建

```http
POST /api/v1/tasks/optimized
```

**描述**: 使用树结构优化的任务创建方法，自动维护树形结构的冗余字段

**请求体**: 与创建任务接口相同

**响应**: 
```json
{
  "code": 200,
  "message": "Task created with tree optimization",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "id": "task_789",
    "title": "新任务",
    "has_children": false,
    "children_count": 0,
    "tree_depth": 0,
    "root_task_id": "task_789"
  }
}
```

##### 15. 分页查询日志（支持过滤）

```http
GET /api/v1/journals/paginated
```

**描述**: 分页查询用户日志，支持按类型和时间范围过滤

**查询参数**:
- `page` (int, 可选): 页码，默认为1
- `page_size` (int, 可选): 每页大小，默认为20，最大100
- `journal_type` (string, 可选): 日志类型过滤，可选值：`day`, `week`, `month`, `quarter`, `year`
- `start_date` (string, 可选): 开始时间过滤，ISO 8601格式
- `end_date` (string, 可选): 结束时间过滤，ISO 8601格式

**请求示例**:
```http
GET /api/v1/journals/paginated?page=1&page_size=20&journal_type=week&start_date=2024-01-01T00:00:00Z
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "items": [
      {
        "id": "journal_123",
        "title": "第1周工作总结",
        "content": "本周完成了...",
        "journal_type": "week",
        "time_period": {
          "start": "2024-01-01T00:00:00Z",
          "end": "2024-01-07T23:59:59Z"
        },
        "icon": "📝",
        "created_at": "2024-01-07T20:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 15,
      "total_pages": 1,
      "has_next": false,
      "has_prev": false
    }
  }
}
```

#### 计划管理

##### 1. 获取计划列表（按时间周期）

```http
GET /api/v1/plans
```

**描述**: 根据指定的时间周期类型和时间范围获取计划信息，包含该时间段内的任务、日志和统计信息

**请求体**:
```json
{
  "period_type": "day|week|month|quarter|year",
  "start_date": "2023-08-05T00:00:00Z",
  "end_date": "2023-08-12T00:00:00Z"
}
```

**查询参数说明**:
- `period_type` (string, 必填): 时间周期类型
  - `day`: 日计划
  - `week`: 周计划
  - `month`: 月计划
  - `quarter`: 季度计划
  - `year`: 年度计划
- `start_date` (string, 必填): 开始时间，ISO 8601 格式
- `end_date` (string, 必填): 结束时间，ISO 8601 格式

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "tasks": [
      {
        "id": "task_123",
        "title": "完成API文档编写",
        "type": "day",
        "period": {
          "start": "2023-08-05T09:00:00Z",
          "end": "2023-08-05T18:00:00Z"
        },
        "tags": ["开发", "文档"],
        "icon": "📝",
        "score": 85,
        "status": 2,
        "priority": 2,
        "parent_id": "",
        "user_id": "user_456",
        "created_at": "2023-08-05T08:00:00Z",
        "updated_at": "2023-08-05T17:30:00Z"
      }
    ],
    "tasks_total": 5,
    "journals": [
      {
        "id": "journal_456",
        "title": "工作日志",
        "content": "今天完成了API文档编写...",
        "journal_type": "day",
        "time_period": {
          "start": "2023-08-05T00:00:00Z",
          "end": "2023-08-06T00:00:00Z"
        },
        "icon": "📖",
        "created_at": "2023-08-05T18:00:00Z",
        "updated_at": "2023-08-05T18:30:00Z",
        "user_id": "user_456"
      }
    ],
    "journals_total": 3,
    "plan_type": "week",
    "plan_period": {
      "start": "2023-08-05T00:00:00Z",
      "end": "2023-08-12T00:00:00Z"
    },
    "score_total": 425,
    "group_stats": [
      {
        "group_key": "2023-08-05",
        "task_count": 2,
        "score_total": 85
      },
      {
        "group_key": "2023-08-06", 
        "task_count": 1,
        "score_total": 92
      }
    ]
  }
}
```

**响应字段说明**:
- `tasks`: 该时间段内的任务列表
- `tasks_total`: 任务总数
- `journals`: 该时间段内的日志列表  
- `journals_total`: 日志总数
- `plan_type`: 计划类型（与请求的period_type相同）
- `plan_period`: 计划时间段
- `score_total`: 总分数（所有任务分数之和）
- `group_stats`: 分组统计信息
  - `group_key`: 分组键（根据plan_type不同格式不同）
    - day: "2023-08-05" (日期)
    - week: "2023-W32" (ISO周)  
    - month: "2023-08" (年月)
    - quarter: "2023-Q3" (季度)
    - year: "2023" (年份)
  - `task_count`: 该分组内的任务数量
  - `score_total`: 该分组内的分数总和

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权（未登录或 Session 无效） |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 通用响应格式

所有 API 响应都遵循以下格式：

```json
{
  "code": 200,
  "message": "Success",
  "success": true,
  "timestamp": 1691234567,
  "data": {}
}
```

**字段说明**:
- `code`: HTTP 状态码
- `message`: 响应消息
- `success`: 操作是否成功
- `timestamp`: 响应时间戳
- `data`: 响应数据（可选）

---

## 使用示例

### 1. 登录并获取任务列表

```bash
# 1. 登录
curl -X POST http://localhost:8081/api/v1/public/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your_username",
    "password": "your_password"
  }'

# 2. 使用返回的 Session 获取任务列表
curl -X GET http://localhost:8081/api/v1/tasks \
  -H "Authorization: Bearer your_session_id"
```

### 2. 创建任务

```bash
curl -X POST http://localhost:8081/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_session_id" \
  -d '{
    "title": "完成项目文档",
    "description": "编写 API 文档和用户手册",
    "priority": "high",
    "due_date": "2023-08-10T18:00:00Z"
  }'
```

### 3. 创建日志

```bash
# 创建一个日志条目
curl -X POST http://localhost:8081/api/v1/journals \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_session_id" \
  -d '{
    "title": "每日工作总结",
    "content": "今天完成了API文档的编写工作，修复了3个bug，完成度85%",
    "journal_type": "day",
    "start_date": "2023-08-05T00:00:00Z",
    "end_date": "2023-08-06T00:00:00Z",
    "icon": "📝"
  }'
```

### 4. 获取日志列表

```bash
# 获取某一天的日志
curl -X GET http://localhost:8081/api/v1/journals \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_session_id" \
  -d '{
    "period_type": "day",
    "start_date": "2023-08-05T00:00:00Z",
    "end_date": "2023-08-06T00:00:00Z"
  }'

# 获取某一周的日志
curl -X GET http://localhost:8081/api/v1/journals \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_session_id" \
  -d '{
    "period_type": "week", 
    "start_date": "2023-07-31T00:00:00Z",
    "end_date": "2023-08-07T00:00:00Z"
  }'
```

### 5. 更新日志

```bash
# 更新日志内容
curl -X PUT http://localhost:8081/api/v1/journals/journal_123 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_session_id" \
  -d '{
    "journal_id": "journal_123",
    "title": "每日工作总结（已更新）",
    "content": "今天完成了API文档的编写工作，修复了5个bug，完成度90%"
  }'
```

### 6. 创建任务

```bash
curl -X POST http://localhost:8081/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_session_id" \
  -d '{
    "title": "完成项目文档",
    "description": "编写 API 文档和用户手册",
    "priority": "high",
    "start_date": "2023-08-05T09:00:00Z",
    "end_date": "2023-08-10T18:00:00Z"
  }'
```

### 6. 创建任务

```bash
curl -X POST http://localhost:8081/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_session_id" \
  -d '{
    "title": "完成项目文档",
    "description": "编写 API 文档和用户手册",
    "priority": "high",
    "start_date": "2023-08-05T09:00:00Z",
    "end_date": "2023-08-10T18:00:00Z",
    "icon": "📝",
    "tags": ["开发", "文档"]
  }'
```

### 7. 获取任务列表

```bash
# 获取某一周的任务
curl -X GET http://localhost:8081/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_session_id" \
  -d '{
    "period_type": "week",
    "start_date": "2023-07-31T00:00:00Z",
    "end_date": "2023-08-07T00:00:00Z"
  }'
```

### 8. 获取计划信息

```bash
# 获取某一周的计划（包含任务和日志）
curl -X GET http://localhost:8081/api/v1/plans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_session_id" \
  -d '{
    "period_type": "week",
    "start_date": "2023-07-31T00:00:00Z", 
    "end_date": "2023-08-07T00:00:00Z"
  }'
```

### 9. 创建子任务

```bash
curl -X POST http://localhost:8081/api/v1/tasks/parent_task_123/subtasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_session_id" \
  -d '{
    "title": "审查API文档",
    "description": "审查 API 文档的准确性和完整性",
    "priority": "medium",
    "start_date": "2023-08-08T09:00:00Z",
    "end_date": "2023-08-08T17:00:00Z",
    "task_id": "parent_task_123",
    "icon": "🔍",
    "tags": ["审查", "文档"]
  }'
```

---

## 部署信息

- **Docker 端口**: 8081
- **数据库**: PostgreSQL (端口 15432)
- **健康检查**: `/health`
- **配置文件**: `configs/config.ini`

更多详情请参考项目 README 和部署文档。
