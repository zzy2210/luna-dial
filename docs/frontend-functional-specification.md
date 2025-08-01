# Luna-Dial OKR管理系统 - 前端功能设计文档

## 文档概述

**文档版本**: v1.0  
**创建日期**: 2025年7月28日  
**目标受众**: 前端开发团队  
**系统名称**: Luna-Dial OKR年度计划管理系统  
**后端技术栈**: Go + Echo + Ent ORM + PostgreSQL  

## 系统概述

Luna-Dial是一个基于OKR（目标和关键成果）方法论的个人/团队目标管理系统，支持多层级任务管理、日志记录、数据统计分析等功能。系统支持从年度到日任务的五级时间层次管理，提供完整的任务生命周期追踪和分析。

## 核心功能模块

### 1. 用户认证模块

#### 1.1 用户注册功能
- **接口**: `POST /api/auth/register`
- **功能描述**: 用户账号注册
- **请求参数**:
  ```json
  {
    "username": "string (3-50字符)",
    "email": "string (邮箱格式)",
    "password": "string (最少6字符)"
  }
  ```
- **响应格式**:
  ```json
  {
    "success": true,
    "data": {
      "id": "uuid",
      "username": "string",
      "email": "string",
      "created_at": "timestamp"
    }
  }
  ```
- **前端实现要点**:
  - 表单验证（用户名长度、邮箱格式、密码强度）
  - 重复密码确认
  - 注册成功后自动跳转登录

#### 1.2 用户登录功能
- **接口**: `POST /api/auth/login`
- **功能描述**: 用户登录认证
- **请求参数**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **响应格式**:
  ```json
  {
    "success": true,
    "data": {
      "token": "jwt_token",
      "user": {
        "id": "uuid",
        "username": "string",
        "email": "string"
      }
    }
  }
  ```
- **前端实现要点**:
  - JWT Token本地存储
  - 自动登录状态保持
  - 登录失败错误提示

#### 1.3 用户信息管理
- **获取当前用户信息**: `GET /api/auth/me`
- **更新用户信息**: `PUT /api/auth/me`
- **用户登出**: `POST /api/auth/logout`

### 2. 任务管理模块

#### 2.1 任务类型说明
系统支持5种时间层级的任务类型：
- **年度任务** (year): 年度级别的大目标
- **季度任务** (quarter): 季度级别的关键成果
- **月度任务** (month): 月度级别的具体目标
- **周任务** (week): 周级别的执行计划
- **日任务** (day): 日级别的具体行动

#### 2.2 任务状态管理
- **待开始** (pending): 任务创建但未开始执行
- **进行中** (in-progress): 任务正在执行中
- **已完成** (completed): 任务已完成

#### 2.3 任务CRUD操作

##### 2.3.1 创建任务
- **接口**: `POST /api/tasks`
- **功能描述**: 创建新任务
- **请求参数**:
  ```json
  {
    "title": "string (1-200字符)",
    "description": "string (可选)",
    "type": "year|quarter|month|week|day",
    "start_date": "timestamp (可选)",
    "end_date": "timestamp (可选)",
    "status": "pending|in-progress|completed",
    "score": "number (1-10分，可选)",
    "tags": "string (可选)"
  }
  ```

##### 2.3.2 获取任务列表
- **接口**: `GET /api/tasks`
- **功能描述**: 获取用户任务列表，支持筛选和分页
- **查询参数**:
  - `type`: 任务类型筛选
  - `status`: 任务状态筛选
  - `parent_id`: 父任务ID筛选
  - `page`: 页码 (默认1)
  - `page_size`: 每页条数 (默认20，最大100)

##### 2.3.3 获取任务详情
- **接口**: `GET /api/tasks/{id}`
- **功能描述**: 获取指定任务的详细信息

##### 2.3.4 更新任务
- **接口**: `PUT /api/tasks/{id}`
- **功能描述**: 更新任务信息

##### 2.3.5 删除任务
- **接口**: `DELETE /api/tasks/{id}`
- **功能描述**: 删除指定任务

#### 2.4 任务层级管理

##### 2.4.1 创建子任务
- **接口**: `POST /api/tasks/{id}/subtasks`
- **功能描述**: 为指定任务创建子任务
- **前端实现要点**:
  - 支持任务的层级展示
  - 拖拽式任务层级调整
  - 子任务完成度对父任务的影响展示

##### 2.4.2 获取子任务列表
- **接口**: `GET /api/tasks/{id}/children`
- **功能描述**: 获取指定任务的直接子任务

##### 2.4.3 获取完整任务树
- **接口**: `GET /api/tasks/{id}/tree`
- **功能描述**: 获取以指定任务为根的完整任务树

##### 2.4.4 全局任务视图
- **接口**: `GET /api/tasks/global`
- **功能描述**: 获取用户的全局任务树视图
- **前端实现要点**:
  - 树形结构展示
  - 任务状态可视化
  - 支持折叠/展开

#### 2.5 任务评分功能
- **接口**: `PUT /api/tasks/{id}/score`
- **功能描述**: 更新任务评分 (1-10分)
- **前端实现要点**:
  - 星级评分组件
  - 评分历史记录
  - 评分统计图表

#### 2.6 任务上下文视图
- **接口**: `GET /api/tasks/{id}/context`
- **功能描述**: 获取任务的上下文信息，包括父任务链和兄弟任务
- **前端实现要点**:
  - 面包屑导航显示任务层级
  - 相关任务推荐
  - 任务关联关系图

### 3. 日志管理模块

#### 3.1 日志类型说明
- **开始计划** (plan-start): 计划开始时的记录
- **阶段反思** (reflection): 执行过程中的反思记录
- **结束总结** (summary): 计划结束时的总结记录

#### 3.2 日志CRUD操作

##### 3.2.1 创建日志
- **接口**: `POST /api/journals`
- **功能描述**: 创建新的日志条目
- **请求参数**:
  ```json
  {
    "content": "string (日志内容)",
    "time_reference": "string (时间参考)",
    "time_scale": "day|week|month|quarter|year",
    "entry_type": "plan-start|reflection|summary",
    "task_ids": ["uuid"] (关联任务ID列表，可选)
  }
  ```

##### 3.2.2 获取日志列表
- **接口**: `GET /api/journals`
- **功能描述**: 获取用户日志列表
- **查询参数**:
  - `time_scale`: 时间尺度筛选
  - `entry_type`: 日志类型筛选
  - `start_date`: 开始日期
  - `end_date`: 结束日期
  - `page`: 页码
  - `page_size`: 每页条数

##### 3.2.3 按时间查询日志
- **接口**: `GET /api/journals/time-range`
- **功能描述**: 按时间范围查询日志

#### 3.3 日志与任务关联
- **接口**: `POST /api/journals/{id}/tasks`
- **功能描述**: 将日志关联到多个任务
- **前端实现要点**:
  - 富文本编辑器
  - 任务标签选择器
  - 关联任务快速跳转

### 4. 计划视图模块

#### 4.1 计划视图功能
- **接口**: `GET /api/plan`
- **功能描述**: 获取指定时间范围的计划视图
- **查询参数**:
  - `scale`: 时间尺度 (day|week|month|quarter|year)
  - `time_ref`: 时间参考 (如 "2024-Q4", "2025-07", "2025-07-15")

- **响应格式**:
  ```json
  {
    "success": true,
    "data": {
      "tasks": [{
        "id": "uuid",
        "title": "string",
        "description": "string",
        "type": "string",
        "status": "string",
        "score": "number",
        "ancestors": [], // 完整父级链
        "children": [],  // 直接子任务
        "depth": "number" // 任务层级深度
      }],
      "journals": [], // 相关日志
      "time_range": {
        "start": "timestamp",
        "end": "timestamp"
      },
      "stats": {
        "total_tasks": "number",
        "completed_tasks": "number",
        "in_progress_tasks": "number",
        "pending_tasks": "number",
        "total_score": "number",
        "completed_score": "number"
      }
    }
  }
  ```

#### 4.2 前端实现要点
- **时间导航器**: 支持日/周/月/季/年切换
- **甘特图视图**: 显示任务时间线
- **看板视图**: 按状态分组显示任务
- **日历视图**: 在日历上展示任务和日志
- **统计面板**: 显示完成度、得分等统计信息

### 5. 统计分析模块

#### 5.1 用户概览统计
- **接口**: `GET /api/stats/overview`
- **功能描述**: 获取用户的整体统计概览
- **前端展示**: 仪表盘卡片式布局

#### 5.2 任务完成度统计
- **接口**: `GET /api/stats/completion`
- **功能描述**: 获取任务完成度统计
- **查询参数**: `start`, `end` (时间范围)
- **前端展示**: 环形图、柱状图

#### 5.3 评分趋势统计
- **接口**: `GET /api/stats/score-trend`
- **功能描述**: 获取评分趋势统计
- **前端展示**: 折线图展示分数变化趋势

#### 5.4 基于时间参考的评分趋势
- **接口**: `GET /api/stats/score-trend-ref`
- **功能描述**: 基于时间参考获取评分趋势
- **查询参数**:
  - `scale`: 统计尺度
  - `time_ref`: 时间参考
- **响应格式**:
  ```json
  {
    "success": true,
    "data": {
      "labels": ["2025-07-01", "2025-07-02"],
      "scores": [85, 92],
      "counts": [3, 4],
      "scale": "day",
      "time_ref": "2025-07",
      "summary": {
        "total_score": 177,
        "total_tasks": 7,
        "average_score": 88.5,
        "max_score": 92,
        "min_score": 85
      }
    }
  }
  ```

#### 5.5 时间分布统计
- **接口**: `GET /api/stats/time-distribution`
- **功能描述**: 获取任务时间分布统计
- **前端展示**: 热力图、时间轴图

### 6. 前端架构建议

#### 6.1 技术栈推荐
- **框架**: React/Vue.js/Angular
- **状态管理**: Redux/Vuex/NgRx
- **UI组件库**: Ant Design/Element Plus/Material UI
- **图表库**: ECharts/Chart.js/D3.js
- **日期处理**: Day.js/Moment.js
- **HTTP客户端**: Axios

#### 6.2 页面结构设计

##### 主布局
```
┌─────────────────────────────────────────┐
│ Header (导航栏 + 用户信息)               │
├─────────────────────────────────────────┤
│ Sidebar │ Main Content Area             │
│ (菜单)   │                              │
│         │                               │
│         │                               │
│         │                               │
├─────────┴───────────────────────────────┤
│ Footer (状态栏)                          │
└─────────────────────────────────────────┘
```

##### 核心页面
1. **仪表盘页面** (`/dashboard`)
   - 用户概览统计
   - 近期任务列表
   - 评分趋势图表
   - 快速操作入口

2. **任务管理页面** (`/tasks`)
   - 任务列表视图
   - 任务创建/编辑表单
   - 任务搜索筛选
   - 批量操作功能

3. **任务详情页面** (`/tasks/:id`)
   - 任务详细信息
   - 子任务管理
   - 相关日志显示
   - 任务操作历史

4. **计划视图页面** (`/plan`)
   - 时间维度切换器
   - 多种视图模式切换
   - 任务拖拽编辑
   - 进度统计面板

5. **日志管理页面** (`/journals`)
   - 日志列表展示
   - 富文本编辑器
   - 任务关联选择
   - 日志搜索功能

6. **统计分析页面** (`/analytics`)
   - 各类统计图表
   - 时间范围选择器
   - 数据导出功能
   - 自定义报表

7. **用户设置页面** (`/settings`)
   - 个人信息管理
   - 系统偏好设置
   - 账号安全设置

#### 6.3 状态管理设计

##### 全局状态
```javascript
{
  user: {
    info: {},
    isAuthenticated: boolean,
    token: string
  },
  tasks: {
    list: [],
    currentTask: {},
    filters: {},
    pagination: {}
  },
  journals: {
    list: [],
    currentJournal: {},
    filters: {}
  },
  stats: {
    overview: {},
    trends: {},
    completion: {}
  },
  ui: {
    loading: boolean,
    error: string,
    currentView: string
  }
}
```

#### 6.4 API调用封装

##### HTTP拦截器配置
```javascript
// 请求拦截器 - 自动添加认证头
request.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// 响应拦截器 - 统一错误处理
request.interceptors.response.use(
  response => response.data,
  error => {
    if (error.response?.status === 401) {
      // 清除认证信息，跳转登录
      localStorage.removeItem('token');
      router.push('/login');
    }
    return Promise.reject(error);
  }
);
```

##### API服务封装
```javascript
// userAPI.js
export const userAPI = {
  login: (data) => request.post('/api/auth/login', data),
  register: (data) => request.post('/api/auth/register', data),
  getCurrentUser: () => request.get('/api/auth/me'),
  updateUser: (data) => request.put('/api/auth/me', data),
  logout: () => request.post('/api/auth/logout')
};

// taskAPI.js  
export const taskAPI = {
  getTasks: (params) => request.get('/api/tasks', { params }),
  getTask: (id) => request.get(`/api/tasks/${id}`),
  createTask: (data) => request.post('/api/tasks', data),
  updateTask: (id, data) => request.put(`/api/tasks/${id}`, data),
  deleteTask: (id) => request.delete(`/api/tasks/${id}`),
  getGlobalView: () => request.get('/api/tasks/global'),
  updateTaskScore: (id, score) => request.put(`/api/tasks/${id}/score`, { score })
};
```

### 7. 错误处理规范

#### 7.1 标准错误响应格式
```json
{
  "success": false,
  "error": "ERROR_CODE",
  "message": "错误描述信息",
  "code": 400
}
```

#### 7.2 常见错误码
- `USER_NOT_FOUND`: 用户不存在
- `TASK_NOT_FOUND`: 任务不存在
- `JOURNAL_NOT_FOUND`: 日志不存在
- `UNAUTHORIZED`: 未授权访问
- `INVALID_REQUEST`: 请求参数无效
- `DUPLICATE_USER`: 用户名或邮箱已存在
- `INVALID_PASSWORD`: 密码错误

#### 7.3 前端错误处理
- 网络错误统一提示
- 表单验证错误高亮显示
- 权限错误自动跳转登录
- 操作失败友好提示

### 8. 数据格式规范

#### 8.1 日期时间格式
- **前端显示**: 根据用户偏好本地化
- **API传输**: ISO 8601格式 (YYYY-MM-DDTHH:mm:ssZ)
- **时间参考格式**:
  - 年: "2025"
  - 季度: "2024-Q4"
  - 月: "2025-07"
  - 周: "2025-W30"
  - 日: "2025-07-15"

#### 8.2 分页数据格式
```json
{
  "success": true,
  "data": [],
  "total": 100,
  "current_page": 1,
  "page_size": 20,
  "total_pages": 5
}
```

### 9. 性能优化建议

#### 9.1 数据加载优化
- 列表页面使用虚拟滚动
- 图表数据按需加载
- 实现缓存机制避免重复请求

#### 9.2 用户体验优化
- 骨架屏加载状态
- 操作反馈及时响应
- 离线数据缓存
- 页面切换动画

### 10. 开发规范

#### 10.1 代码规范
- 统一的代码格式化工具
- 组件命名规范
- 注释规范要求
- TypeScript类型定义

#### 10.2 测试要求
- 单元测试覆盖核心业务逻辑
- 集成测试覆盖关键用户流程
- E2E测试覆盖主要功能场景

#### 10.3 部署要求
- 支持多环境配置
- 自动化构建部署
- 版本管理规范

---

**文档维护**: 本文档将随着后端API的更新而持续维护，请及时同步最新版本。  
**联系方式**: 如有疑问，请联系后端开发团队获取技术支持。
