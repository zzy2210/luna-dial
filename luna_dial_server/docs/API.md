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
    "session_id": "string",
    "user_id": "string",
    "username": "string"
  }
}
```

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

**描述**: 获取当前登录用户的详细信息

**响应**:
```json
{
  "code": 200,
  "message": "Success",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "id": "string",
    "username": "string",
    "email": "string",
    "created_at": "2023-08-05T10:30:00Z",
    "updated_at": "2023-08-05T10:30:00Z"
  }
}
```

##### 2. 用户登出

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

#### 用户管理

##### 1. 获取当前用户信息

```http
GET /api/v1/users/me
```

**描述**: 获取当前登录用户的基本信息

**响应**:
```json
{
  "code": 200,
  "message": "Success",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "id": "string",
    "username": "string",
    "email": "string"
  }
}
```

#### 日志管理

##### 1. 获取日志列表

```http
GET /api/v1/journals?period=2023-08&page=1&limit=10
```

**描述**: 按时间周期获取日志列表

**查询参数**:
- `period` (string): 时间周期，格式 YYYY-MM
- `page` (int, 可选): 页码，默认 1
- `limit` (int, 可选): 每页数量，默认 10

**响应**:
```json
{
  "code": 200,
  "message": "Success",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "journals": [
      {
        "id": "string",
        "title": "string",
        "content": "string",
        "created_at": "2023-08-05T10:30:00Z",
        "updated_at": "2023-08-05T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 50,
      "total_pages": 5
    }
  }
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
  "title": "string",
  "content": "string",
  "tags": ["string"]
}
```

**响应**:
```json
{
  "code": 200,
  "message": "Journal created successfully",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "id": "string",
    "title": "string",
    "content": "string",
    "created_at": "2023-08-05T10:30:00Z"
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

**请求体**:
```json
{
  "title": "string",
  "content": "string",
  "tags": ["string"]
}
```

**响应**:
```json
{
  "code": 200,
  "message": "Journal updated successfully",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "id": "string",
    "title": "string",
    "content": "string",
    "updated_at": "2023-08-05T10:30:00Z"
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
```json
{
  "code": 200,
  "message": "Journal deleted successfully",
  "success": true,
  "timestamp": 1691234567
}
```

#### 任务管理

##### 1. 获取任务列表

```http
GET /api/v1/tasks?status=pending&page=1&limit=10
```

**描述**: 获取用户的任务列表

**查询参数**:
- `status` (string, 可选): 任务状态 (pending, completed, cancelled)
- `page` (int, 可选): 页码，默认 1
- `limit` (int, 可选): 每页数量，默认 10

**响应**:
```json
{
  "code": 200,
  "message": "Success",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "tasks": [
      {
        "id": "string",
        "title": "string",
        "description": "string",
        "status": "pending",
        "priority": "high",
        "due_date": "2023-08-10T18:00:00Z",
        "created_at": "2023-08-05T10:30:00Z",
        "updated_at": "2023-08-05T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 25,
      "total_pages": 3
    }
  }
}
```

##### 2. 创建任务

```http
POST /api/v1/tasks
```

**描述**: 创建新任务

**请求体**:
```json
{
  "title": "string",
  "description": "string",
  "priority": "high|medium|low",
  "due_date": "2023-08-10T18:00:00Z",
  "tags": ["string"]
}
```

**响应**:
```json
{
  "code": 200,
  "message": "Task created successfully",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "id": "string",
    "title": "string",
    "description": "string",
    "status": "pending",
    "priority": "high",
    "due_date": "2023-08-10T18:00:00Z",
    "created_at": "2023-08-05T10:30:00Z"
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

**请求体**:
```json
{
  "title": "string",
  "description": "string",
  "priority": "high|medium|low",
  "due_date": "2023-08-10T18:00:00Z",
  "status": "pending|completed|cancelled"
}
```

##### 4. 删除任务

```http
DELETE /api/v1/tasks/{task_id}
```

**描述**: 删除指定任务

**路径参数**:
- `task_id` (string): 任务 ID

##### 5. 完成任务

```http
POST /api/v1/tasks/{task_id}/complete
```

**描述**: 标记任务为已完成

**路径参数**:
- `task_id` (string): 任务 ID

**响应**:
```json
{
  "code": 200,
  "message": "Task completed successfully",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "id": "string",
    "status": "completed",
    "completed_at": "2023-08-05T15:30:00Z"
  }
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
  "title": "string",
  "description": "string",
  "priority": "high|medium|low",
  "due_date": "2023-08-10T18:00:00Z",
  "tags": ["string"]
}
```

**响应**:
```json
{
  "code": 200,
  "message": "Subtask created successfully",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "id": "string",
    "parent_task_id": "string",
    "title": "string",
    "description": "string",
    "status": "pending",
    "priority": "high",
    "due_date": "2023-08-10T18:00:00Z",
    "created_at": "2023-08-05T10:30:00Z"
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
  "score": 85,
  "comment": "string"
}
```

#### 计划管理

##### 1. 获取计划列表

```http
GET /api/v1/plans?type=daily&page=1&limit=10
```

**描述**: 获取用户的计划列表

**查询参数**:
- `type` (string, 可选): 计划类型 (daily, weekly, monthly)
- `page` (int, 可选): 页码，默认 1
- `limit` (int, 可选): 每页数量，默认 10

**响应**:
```json
{
  "code": 200,
  "message": "Success",
  "success": true,
  "timestamp": 1691234567,
  "data": {
    "plans": [
      {
        "id": "string",
        "title": "string",
        "description": "string",
        "type": "daily",
        "start_date": "2023-08-05T00:00:00Z",
        "end_date": "2023-08-05T23:59:59Z",
        "tasks": ["task_id_1", "task_id_2"],
        "created_at": "2023-08-05T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 15,
      "total_pages": 2
    }
  }
}
```

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

### 3. 创建子任务

```bash
curl -X POST http://localhost:8081/api/v1/tasks/parent_task_id/subtasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_session_id" \
  -d '{
    "title": "审查文档",
    "description": "审查 API 文档的准确性",
    "priority": "medium",
    "due_date": "2023-08-09T12:00:00Z"
  }'
```

---

## 部署信息

- **Docker 端口**: 8081
- **数据库**: PostgreSQL (端口 15432)
- **健康检查**: `/health`
- **配置文件**: `configs/config.ini`

更多详情请参考项目 README 和部署文档。
