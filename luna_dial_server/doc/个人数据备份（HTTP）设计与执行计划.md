# 个人数据备份（HTTP）设计与执行计划

## 0. 背景与目标

### 目标

为 Luna Dial 后端新增“个人数据备份”能力，提供两项最小但可维护的能力：

- **导出备份**：当前登录用户可通过 HTTP (中文名：超文本传输协议接口) 导出自己的业务数据为一个可携带的备份包。
- **覆盖式导入**：当前登录用户可上传备份包并执行覆盖式恢复：先删除自己现有业务数据，再导入备份中的数据。

### 非目标（v1 明确不做）

- **不做合并导入**（merge，中文名：去重合并）：只做覆盖式导入。
- **不做全库备份**：只针对“当前会话用户”的数据。
- **不备份 `users` 表**：不导出用户基本信息、密码等。
- **不备份 `system_configs`**：不导出 JWT 密钥等系统级敏感配置。
- **不做前端 UI**：先落后端能力与文档，UI 以后再说。
- **不做加密备份**（encryption，中文名：加密）：v1 先保证正确性与兼容性，后续再评估口令加密。

## 1. 现状（代码与数据）

### 技术栈

- 后端：Go + Echo + GORM + PostgreSQL。
- 鉴权：已有 SessionMiddleware (中文名：会话鉴权中间件)，受保护路由挂载在 `/api/v1` 下。

### 业务数据范围（v1）

只处理这两类“个人业务数据”：

- `tasks`（任务）
- `journals`（日志）

注意：`tasks` 表里存在树优化冗余字段（如 `root_task_id`、`tree_depth`、`children_count`、`has_children`）。这些字段属于 derived fields（中文名：可重建的派生冗余数据），导入后应重建/校准，不应成为备份格式的硬依赖。

## 2. 用户故事与约束

### 用户故事

- 作为一个用户，我想把自己的任务/日志导出成文件，便于迁移或备份。
- 作为一个用户，我想把备份文件导入回来，并用它覆盖我当前的数据，恢复到备份时的状态。

### 约束（Fail Fast）

- 导入语义固定为覆盖式：仅提供“覆盖式导入”接口路径，避免伪选项。
- 导入必须是原子操作：要么全成功，要么全失败（事务）。
- 不允许跨用户写入：导入时忽略备份中携带的 `user_id`，统一覆盖为“当前会话用户 ID”。
- 默认防误导入：若备份主体用户名与当前会话用户名不一致，则拒绝导入；仅在显式 `force=true` 时允许继续。

## 3. API 设计（HTTP）

### 3.1 路由（建议）

使用受保护路由组 `/api/v1`（已挂载 SessionMiddleware）新增：

- `POST /api/v1/personal-backup/export`
- `POST /api/v1/personal-backup/import-overwrite[?force=true]`

命名说明：
- “personal-backup”对齐“个人数据备份”功能名称，避免未来出现“管理员全库备份”混用同一路径导致权限语义污染。

### 3.2 认证与授权

- 认证：必须通过 SessionMiddleware（复用现有登录态）。
- 授权：只允许操作“当前 session 对应 user_id”的数据；不提供 `user_id` 参数，避免引入管理员权限模型。

### 3.3 导出接口：`POST /api/v1/personal-backup/export`

#### 请求

- 方法：`POST`
- Body：无（可扩展为 JSON 参数，但 v1 不需要）

#### 响应

返回二进制流：

- `Content-Type: application/gzip`（或 `application/octet-stream`）
- `Content-Disposition: attachment; filename="luna-dial-personal-backup-YYYYMMDD-HHMMSS.tar.gz"`

失败时返回统一 JSON 错误结构（详见 3.5）。

### 3.4 导入接口：`POST /api/v1/personal-backup/import-overwrite[?force=true]`

#### 请求

- 方法：`POST`
- Query：
  - `force=true`（可选；默认 false。用于在“备份主体用户名 != 当前会话用户名”时强制导入）
- Body：`multipart/form-data`（中文名：表单文件上传）：
  - 字段 `file`：备份包（`.tar.gz`）

#### 响应

成功返回 JSON（建议）：

```json
{
  "success": true,
  "data": {
    "imported_tasks": 123,
    "imported_journals": 45
  }
}
```

### 3.5 错误返回（统一结构）

沿用当前服务端的 response 规范（如果已有统一响应体，则以现有为准），建议错误码至少包含：

- `400`：参数错误（缺少 overwrite、文件缺失、格式不对）
- `400`：参数错误（文件缺失、格式不对）
- `401`：未登录/会话失效
- `413`：文件过大（超出限制）
- `422`：备份包语义不合法（manifest 校验失败、hash 不匹配等）
- `500`：服务端错误

## 4. 备份包格式（可维护、可演进）

### 4.1 容器格式

使用 `tar.gz`（中文名：tar 打包 + gzip 压缩）作为备份容器。原因：

- 单文件易携带
- 可流式生成/解析（streaming，中文名：边读边处理）
- 工具链成熟

### 4.2 文件清单

备份包包含以下文件：

- `manifest.json`（manifest，中文名：清单/元数据）
- `tasks.jsonl`（JSONL，中文名：逐行 JSON）
- `journals.jsonl`

### 4.3 manifest.json（建议字段）

```json
{
  "schema_version": 1,
  "created_at": "2026-02-26T00:00:00Z",
  "subject": {
    "username": "admin"
  },
  "app": {
    "name": "luna-dial-server",
    "version": "1.0.0"
  },
  "content": {
    "tasks": { "filename": "tasks.jsonl", "sha256": "<hex>", "count": 0 },
    "journals": { "filename": "journals.jsonl", "sha256": "<hex>", "count": 0 }
  }
}
```

说明：
- `schema_version` 是备份格式版本，不是数据库 migration 版本。
- `sha256`（中文名：SHA-256 哈希校验）用于完整性校验，避免导入半截文件或被篡改文件。
- `subject.username`（中文名：备份主体用户名）用于防误导入：导入时若与当前会话用户名不一致，默认拒绝，除非显式 `force=true`。

### 4.4 tasks.jsonl（字段范围与类型）

每行一个 JSON 对象（JSONL，中文名：逐行 JSON），仅包含“业务必需字段”，不依赖可重建字段：

- `id`: string
- `title`: string
- `task_type`: string（枚举：`day|week|month|quarter|year`）
- `period_start`: string|null（RFC3339 时间；允许 null）
- `period_end`: string|null（RFC3339 时间；允许 null）
- `tags`: string[]（允许空数组）
- `icon`: string
- `score`: number
- `status`: string（枚举：`not_started|in_progress|completed|cancelled`）
- `priority`: string（枚举：`low|medium|high|urgent`）
- `parent_id`: string（根任务为空字符串）
- `created_at`: string（RFC3339 时间）
- `updated_at`: string（RFC3339 时间）

注意：
- 不写入 `user_id` 也可以（manifest 隐式指向“导出者”），但实现更简单的是允许存在该字段；无论是否存在，**导入时统一覆盖为当前 user_id**。
- `root_task_id`、`tree_depth`、`children_count`、`has_children` 不作为输入信任源；导入后重建。
 - 字段使用字符串枚举而非 int：避免未来枚举值调整导致旧备份不可导入。

示例（单行）：

```json
{"id":"<uuid>","title":"示例任务","task_type":"week","period_start":"2026-02-24T00:00:00Z","period_end":"2026-03-02T00:00:00Z","tags":["work","okr"],"icon":"","score":0,"status":"not_started","priority":"medium","parent_id":"","created_at":"2026-02-26T01:02:03Z","updated_at":"2026-02-26T01:02:03Z"}
```

### 4.5 journals.jsonl（字段范围与类型）

- `id`: string
- `title`: string
- `content`: string
- `journal_type`: string（枚举：`day|week|month|quarter|year`）
- `period_start`: string|null（RFC3339 时间；允许 null）
- `period_end`: string|null（RFC3339 时间；允许 null）
- `icon`: string
- `created_at`: string（RFC3339 时间）
- `updated_at`: string（RFC3339 时间）

同样：导入时统一覆盖 `user_id` 为当前用户。

## 5. 导入算法（覆盖式、事务、重建）

### 5.1 总体流程

1) 解析 tar.gz，读取 `manifest.json`
2) 校验 `schema_version` 是否支持
3) 计算 `tasks.jsonl` / `journals.jsonl` 的 sha256 并与 manifest 比对
4) 开启数据库事务（transaction，中文名：事务）
5) 删除当前用户数据：
   - `DELETE FROM tasks WHERE user_id = ?`
   - `DELETE FROM journals WHERE user_id = ?`
6) 插入备份数据（导入时覆盖 user_id 为当前用户）
7) 重建任务树冗余字段（root_task_id/tree_depth/children_count/has_children）
8) 提交事务；任何失败都回滚

### 5.2 重建任务树冗余字段（建议策略）

目标：保证导入后查询性能字段一致，不依赖备份内容是否携带这些字段。

建议两段式：

1) 全量重置：
   - `root_task_id = id`（先把每个任务视为根）
   - `tree_depth = 0`
   - `children_count = 0`
   - `has_children = false`

2) 计算真实值：
   - 对每个任务，沿 `parent_id` 向上追溯得到根（设置 `root_task_id`，并计算 `tree_depth`）
   - 对每个父任务，统计直接子任务数量更新 `children_count/has_children`

实现上应注意避免 O(n^2) 的天坑；可用一次性加载 tasks 到内存构建 map，再线性计算。

## 6. 安全与运维约束

### 6.1 文件大小与超时

- 建议限制上传大小，例如 50MB（之后可配置）。
- 导入/导出应有超时控制，避免长连接占满资源。

### 6.2 日志与隐私

- 日志中不得打印备份内容。
- 错误信息可以包含“第 N 行解析失败”，但不要回显整行内容（journals 可能含隐私）。

### 6.3 速率限制（可选）

如果部署在公网，建议对导入/导出接口加 rate limit（中文名：限流）。v1 可以先不做，但文档中明确风险。

## 7. 兼容性策略（必须写清楚的承诺）

- 新版本必须能导入旧版本导出的备份（至少同一 major 版本内）。
- 不承诺旧版本能导入新版本备份。
- `schema_version` 递增时必须提供迁移策略或拒绝导入并给出明确错误。

## 8. 测试计划（最小但有效）

### 单元测试

- manifest 生成与 sha256 校验一致性
- tasks/journals 的 JSONL 编解码（含特殊字符、换行、空字段）
- tar.gz 打包/解包完整性

### 集成测试（建议）

- 准备一组用户数据 -> 导出 -> 清空 -> 导入 -> 校验数量/关键字段一致
- 校验导入后树冗余字段正确（至少 root_task_id、tree_depth、children_count）

## 9. 里程碑（执行计划）

1) 文档定稿（本文件）
2) 后端新增路由与 handler（导出/导入）
3) 实现备份格式与校验
4) 增加测试（至少单测 + 一条集成测试）
5) 手工验证与回归（登录态、现有 API 不受影响）
