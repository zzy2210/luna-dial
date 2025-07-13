# 开头

md 废物AI 浪费了我一堆资源,最还是不太行.有坑

彳亍 我还是古法编程吧.

## 定义

不过大致的定义做的差不多了,所以我现在梳理一下:

### DDD核心概念说明

- **实体（Entity）**：有唯一标识（如id），生命周期明确。
- **值对象（Value Object）**：只关心值本身，无唯一标识，通常不可变。
- **聚合根（Aggregate Root）**：聚合的入口，外部只能通过聚合根访问聚合内对象。

### 分层

这里参考 kratos 的分层思维:

- server (路由注册)
- service （handle）
- biz （实际业务逻辑）
- data （repo）

### 领域行为与业务规则设计（目录）

#### 领域行为

- User
  - 注册用户
  - 用户登录
  - 编辑用户信息
  - 删除用户
  - 查询用户信息

- Task
  - 创建任务（需指定 userID）
  - 编辑任务（需校验归属 userID）
  - 删除任务（需校验归属 userID）
  - 完成任务
  - 设置任务分数（只给日类型任务）
  - 新建子任务（继承tag，需指定 userID）
  - 添加/移除tag
  - 查看指定时间/时间段 指定类型的任务与父任务聚合（需指定 userID）
  - 查看指定任务的整体任务树（需指定 userID）
  - 添加 emoji icon （需指定 userID）

- Journal
  - 新建日志（需指定 userID）
  - 删除日志（需校验归属 userID）
  - 编辑日志（需校验归属 userID）
  - 查看指定时间/时间段 指定类型的日志聚合（需指定 userID）
  - 添加 emoji icon （需指定 userID）

- Plan
  - 查看指定时间/时间段的任务与日志聚合（需指定 userID）
- 查看指定时间/时间段下，按日/周/月/季度分组的日任务分数总和与总数（需指定 userID，分组粒度可选）

#### 业务规则

- Task
  - 只能被 所属userID的用户查看，修改，删除
  - 当任务完成后，将只可以修改它的分数
  - 子任务userID 与父任务必须一致
  - 任务的 开始时间 早于结束时间 （左闭右开）
  - 标题不能为空，可以重复
- Journal
  - 日志只能被所属 userID 的用户查看、编辑、删除。
  - 日志内容不能为空。
  - 日志类型必须是枚举值之一。
  - 日志的 periodStart 必须早于 periodEnd。
- Plan
  - 只能查看自己（userID）的任务与日志聚合。
  - 统计分数时，只统计“日”类型的任务。
  - 分组统计时，分组粒度（日/周/月/季度）可选，统计结果按分组输出。
  - 聚合的时间区间必须合法（periodStart <= periodEnd）。
- User
  - 用户名必须唯一且不能为空。
  - 密码必须加密存储，且长度不少于8位。
  - 邮箱格式必须合法，且唯一。
  - 用户删除前，必须先删除其所有任务和日志。


#### 聚合与边界

- **Task 聚合**
  - 聚合根：Task
  - 包含：子任务（Task）、标签（Tag，值对象）、任务周期（Period，值对象）、icon（值对象或简单字段）
  - 规则：所有对子任务、标签等的操作，必须通过Task聚合根进行；子任务的userID必须与父任务一致；外部不能直接操作子任务，必须通过父Task。
  - 边界：Task聚合只负责与自己直接相关的子任务、标签等，不直接操作Journal、User等其他聚合。

- **Journal 聚合**
  - 聚合根：Journal
  - 包含：日志周期（Period，值对象）、icon（值对象或简单字段）
  - 规则：日志的所有编辑、删除等操作，必须通过Journal聚合根进行；日志归属userID不可变。
  - 边界：Journal聚合只负责自身内容，不直接操作Task、User等其他聚合。

- **Plan 视图**
  - 不是聚合根，是聚合后的只读视图（DTO/值对象），不落库
  - 来源于Task和Journal聚合的数据聚合
  - 规则：只做数据统计和聚合，不负责业务写操作
  - 支持按日/周/月/季度分组统计日任务分数与数量，分组方式由参数指定。

- **User 聚合**
  - 聚合根：User
  - 包含：用户基本信息
  - 规则：用户的所有信息修改、删除等操作，必须通过User聚合根进行；删除用户前，需先删除其所有Task和Journal
  - 边界：User聚合不直接操作Task、Journal，但可以通过应用服务协调

---

#### task（任务实体，聚合根）

任务模型（实体，聚合根）

- id  uuid（唯一标识，实体ID）
- userID uuid（所属用户ID）
- title string
- type enum 枚举 日/周/月/季/年
- parentID 父任务ID 下级索引上级
- score 分数 （我尽力完成程度）
- isCompleted 是否完成
- periodStart 起始时间（可与 periodEnd 组成值对象Period）
- periodEnd 结束时间  
- tag  标签（可建为值对象）
- icon string emoji 表情
- creatAt
- updateAt

#### journal（日志实体，聚合根）

日志模型（实体，聚合根）

- id  uuid（唯一标识，实体ID）
- userID uuid（所属用户ID）
- title string
- content string 日志内容
- type enum 枚举 日/周/月/季/年
- periodStart 起始时间（可与 periodEnd 组成值对象Period）
- periodEnd 结束时间  
- icon string emoji 表情
- creatAt
- updateAt

#### plan（计划视图，值对象/DTO）

计划视图，不落库，是聚合后的展示模型（值对象/DTO）

- tasks
- taskTotal
- journals
- journalTotal
- type 类型enum 枚举 日/周/月/季/年
- periodStart 起始时间（可与 periodEnd 组成值对象Period）
- periodEnd 结束时间  
- groupStats 分组统计（数组，每个元素包含分组key、任务数、分数总和等）
  - groupKey string（如2025-01、2025-W01等，分组粒度可为日/周/月/季度）
  - taskCount int
  - scoreTotal int


#### user（用户实体，聚合根）

- id uuid（唯一标识，实体ID）
- username string （登录用）
- name string （自定义）
- password string（建议加密存储）
- email string
- createdAt
- updatedAt