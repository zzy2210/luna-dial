# Luna Dial Server API 文档

本目录包含 Luna Dial Server 的完整 API 文档和测试工具。

## 📋 文档列表

### 1. API 文档
- **[API.md](./API.md)** - 完整的 API 接口文档
  - 包含所有端点的详细说明
  - 请求/响应示例
  - 错误码说明
  - 使用示例

### 2. Postman 集合
- **[Luna_Dial_API.postman_collection.json](./Luna_Dial_API.postman_collection.json)** - Postman 测试集合
  - 预配置的所有 API 请求
  - 自动提取和使用 Session ID
  - 包含示例请求体

## 🚀 快速开始

### 使用 Postman 测试

1. **导入集合**
   ```bash
   # 打开 Postman，点击 Import 按钮
   # 选择文件: Luna_Dial_API.postman_collection.json
   ```

2. **设置环境变量**
   - `baseUrl`: `http://localhost:8081` (默认)
   - `sessionId`: 自动从登录响应中提取

3. **测试流程**
   ```
   1. 执行 Health Check → 确认服务运行
   2. 执行 Login → 获取 Session ID (自动保存)
   3. 执行其他受保护的接口
   ```

### 使用 curl 测试

```bash
# 1. 健康检查
curl http://localhost:8081/health

# 2. 登录
curl -X POST http://localhost:8081/api/v1/public/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "testpassword"}'

# 3. 使用 Session 访问受保护接口
curl -X GET http://localhost:8081/api/v1/tasks \
  -H "Cookie: session_id=YOUR_SESSION_ID"
```

## 🔑 认证说明

Luna Dial Server 使用 **Session-based Authentication**：

1. **登录**: `POST /api/v1/public/auth/login`
2. **携带 Session**: 在后续请求中通过 `Cookie: session_id=xxx` 携带
3. **登出**: `POST /api/v1/auth/logout`

## 📊 API 结构

```
Luna Dial Server API
├── 🔓 公开接口
│   ├── /health (健康检查)
│   ├── /version (版本信息)
│   └── /api/v1/public/auth/login (用户登录)
│
└── 🔒 受保护接口 (需要 Session)
    ├── /api/v1/auth/* (认证管理)
    ├── /api/v1/users/* (用户管理)
    ├── /api/v1/journals/* (日志管理)
    ├── /api/v1/tasks/* (任务管理)
    └── /api/v1/plans/* (计划管理)
```

## 🐛 错误处理

所有 API 响应统一格式：

```json
{
  "code": 200,
  "message": "Success",
  "success": true,
  "timestamp": 1691234567,
  "data": {}
}
```

常见错误码：
- `200` - 成功
- `400` - 请求参数错误
- `401` - 未授权（未登录或 Session 无效）
- `404` - 资源不存在
- `500` - 服务器错误

## 🔧 开发测试

### 启动服务

```bash
# 使用 Docker Compose
./start.sh

# 或者直接运行
go run cmd/main.go

# 或者使用 Docker
docker-compose up -d
```

### 验证服务

```bash
# 健康检查
curl http://localhost:8081/health
# 预期响应: "Service is running"

# 版本信息
curl http://localhost:8081/version
# 预期响应: "Version 1.0.0"
```

## 📝 注意事项

1. **Session 过期**: Session 默认 90 分钟过期
2. **数据格式**: 所有请求/响应使用 JSON 格式
3. **时区**: 服务器使用 Asia/Shanghai 时区
4. **分页**: 列表接口支持 `page` 和 `limit` 参数

## 迭代计划
- 任务进度
  - 实现方式
    - 统计下一级的所有子任务（比如年任务统计子任务中的季度任务）
    - 计算他们的加权总和百分比

## 🤝 贡献

如发现 API 文档有误或需要补充，请提交 Issue 或 Pull Request。

---

**项目地址**: [Luna Dial Server](https://github.com/zzy2210/luna-dial)  
**文档更新**: 2025年8月4日
