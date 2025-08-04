# VS Code REST Client 使用指南

## 🚀 快速开始

1. **安装插件**：在 VS Code 中搜索并安装 "REST Client" 插件
2. **打开文件**：打开 `httpclient.http` 或 `advanced.http` 文件
3. **执行请求**：点击请求上方的 "Send Request" 按钮

## 📁 文件说明

### `httpclient.http`
- 基础版本，包含所有 API 接口
- 手动设置 Session ID
- 适合快速测试单个接口

### `advanced.http`
- 高级版本，支持变量提取和链式调用
- 自动从登录响应中提取 Session ID
- 包含完整的测试流程和错误测试

### `http-client.env.json`
- 环境配置文件
- 支持 dev/staging/production 多环境切换

## 🔧 使用技巧

### 1. 变量定义
```http
@baseUrl = http://localhost:8081
@sessionId = your_session_id_here

GET {{baseUrl}}/health
```

### 2. 响应变量提取
```http
# @name login
POST {{baseUrl}}/api/v1/public/auth/login
Content-Type: application/json

{
    "username": "testuser",
    "password": "testpassword"
}

### 使用登录响应中的数据
GET {{baseUrl}}/api/v1/tasks
Cookie: session_id={{login.response.body.data.session_id}}
```

### 3. 环境切换
在 `http-client.env.json` 文件中配置多个环境，然后在请求中使用：
```http
@baseUrl = {{$dotenv %baseUrl}}
@username = {{$dotenv %username}}
```

## 📋 常用操作

### 完整测试流程（推荐使用 advanced.http）
1. 健康检查 → 版本信息 → 登录
2. 创建任务 → 创建子任务 → 更新任务
3. 创建日志 → 更新日志
4. 获取列表数据
5. 登出

### 快速测试单个接口
1. 先执行登录获取 Session ID
2. 复制 Session ID 到变量 `@sessionId`
3. 执行需要测试的接口

## 🎯 快捷键

- `Ctrl+Alt+R` (Windows/Linux) 或 `Cmd+Alt+R` (Mac)：发送请求
- `Ctrl+Alt+L` (Windows/Linux) 或 `Cmd+Alt+L` (Mac)：发送请求并切换到响应视图

## 💡 提示

1. **请求分隔**：使用 `###` 分隔不同的请求
2. **注释**：使用 `#` 添加注释
3. **变量作用域**：变量在整个文件中有效
4. **响应查看**：请求执行后会在右侧显示响应
5. **历史记录**：可以查看请求历史和响应

## 🐛 调试技巧

1. **检查响应**：注意查看响应状态码和错误信息
2. **Session 过期**：如果收到 401 错误，重新执行登录
3. **网络问题**：确保服务器正在运行 (`./start.sh`)
4. **数据格式**：确保 JSON 格式正确

## 📊 对比 Postman

| 特性 | REST Client | Postman |
|------|-------------|---------|
| 启动速度 | ⚡ 快 | 🐌 慢 |
| 内存占用 | 💚 低 | 🔴 高 |
| 文本化配置 | ✅ 支持 | ❌ 图形界面 |
| 版本控制 | ✅ Git 友好 | ❌ 需要导出 |
| 协作 | ✅ 直接分享文件 | ❌ 需要导入导出 |
| 高级功能 | ❌ 基础 | ✅ 丰富 |

REST Client 更适合开发阶段的快速测试和团队协作！
