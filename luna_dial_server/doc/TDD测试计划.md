# TDD测试计划

## 1. 领域对象测试点



### Plan

#### 1.1 计划聚合视图
**函数定义**
- 方法：func GetPlan(param GetPlanParam) (*Plan, error)
- 描述：获取指定用户、指定时间区间的任务与日志聚合视图。
- 参数：param - 包含 userID、periodStart、periodEnd、分组粒度（可选，日/周/月/季度）。
- 返回：Plan 视图对象，或 error。
- 规则：只能查看自己（userID）的任务与日志聚合，聚合的时间区间必须合法（periodStart <= periodEnd）。

**单元测试用例**
- 用例1：成功获取
  - Given：userID 合法，period 合法。
  - When：调用 GetPlan。
  - Then：返回 Plan 视图，包含 tasks、journals、groupStats 等，error 为 nil。
- 用例2：非法 userID
  - Given：userID 非法。
  - When：调用 GetPlan。
  - Then：返回 nil, error，且 error.Error() == "userID is required"。
- 用例3：时间区间非法
  - Given：periodStart > periodEnd。
  - When：调用 GetPlan。
  - Then：返回 nil, error，且 error.Error() == "invalid period range"。

#### 1.2 计划分组统计
**函数定义**
- 方法：func GetPlanGroupStats(param GetPlanGroupStatsParam) ([]GroupStat, error)
- 描述：按指定分组粒度（日/周/月/季度）统计指定时间区间内“日”类型任务的数量和分数总和。
- 参数：param - 包含 userID、periodStart、periodEnd、groupBy（分组粒度）。
- 返回：分组统计数组，或 error。
- 规则：只统计“日”类型任务，分组方式由 groupBy 参数指定。

**单元测试用例**
- 用例1：成功分组统计（月分组）
  - Given：userID 合法，period 合法，groupBy=month。
  - When：调用 GetPlanGroupStats。
  - Then：返回每月的 groupStats，taskCount、scoreTotal 正确，error 为 nil。
- 用例2：成功分组统计（周分组）
  - Given：userID 合法，period 合法，groupBy=week。
  - When：调用 GetPlanGroupStats。
  - Then：返回每周的 groupStats，taskCount、scoreTotal 正确，error 为 nil。
- 用例3：非法分组参数
  - Given：groupBy 非法。
  - When：调用 GetPlanGroupStats。
  - Then：返回 nil, error，且 error.Error() == "invalid group by"。
- 用例4：只统计日类型任务
  - Given：存在多种类型任务。
  - When：调用 GetPlanGroupStats。
  - Then：只统计 type=day 的任务。


### Journal

#### 1.1 新建日志
**函数定义**
- 方法：func CreateJournal(req CreateJournalRequest) (*Journal, error)
- 描述：为指定用户新建日志。
- 参数：param - 包含 userID、title、content、type、periodStart、periodEnd、icon。
- 返回：新建的 Journal 实体对象，或 error。
- 规则：userID、content 不能为空，type 必须为枚举值，periodStart < periodEnd。

**单元测试用例**
- 用例1：成功新建
  - Given：各字段合法。
  - When：调用 CreateJournal。
  - Then：返回新建 Journal，error 为 nil。
- 用例2：userID 为空
  - Given：userID 为空。
  - When：调用 CreateJournal。
  - Then：返回 nil, error，且 error.Error() == "userID is required"。
- 用例3：内容为空
  - Given：content 为空。
  - When：调用 CreateJournal。
  - Then：返回 nil, error，且 error.Error() == "content is required"。
- 用例4：类型非法
  - Given：type 非法。
  - When：调用 CreateJournal。
  - Then：返回 nil, error，且 error.Error() == "invalid journal type"。
- 用例5：period 非法
  - Given：periodStart >= periodEnd。
  - When：调用 CreateJournal。
  - Then：返回 nil, error，且 error.Error() == "invalid period range"。

#### 1.2 编辑日志
**函数定义**
- 方法：func UpdateJournal(req UpdateJournalRequest) (*Journal, error)
- 描述：编辑指定日志，只有归属用户可操作。
- 参数：param - 包含 journalID、userID、可更新字段。
- 返回：更新后的 Journal，或 error。
- 规则：只能归属用户操作，内容不能为空，类型必须为枚举值，periodStart < periodEnd。

**单元测试用例**
- 用例1：成功编辑
  - Given：归属用户，字段合法。
  - When：调用 UpdateJournal。
  - Then：返回更新后的 Journal，error 为 nil。
- 用例2：非归属用户
  - Given：userID 与日志归属不符。
  - When：调用 UpdateJournal。
  - Then：返回 nil, error，且 error.Error() == "no permission"。
- 用例3：内容为空
  - Given：content 为空。
  - When：调用 UpdateJournal。
  - Then：返回 nil, error，且 error.Error() == "content is required"。
- 用例4：类型非法
  - Given：type 非法。
  - When：调用 UpdateJournal。
  - Then：返回 nil, error，且 error.Error() == "invalid journal type"。
- 用例5：period 非法
  - Given：periodStart >= periodEnd。
  - When：调用 UpdateJournal。
  - Then：返回 nil, error，且 error.Error() == "invalid period range"。

#### 1.3 删除日志
**函数定义**
- 方法：func DeleteJournal(req DeleteJournalRequest) error
- 描述：删除指定日志，只有归属用户可操作。
- 参数：param - 包含 journalID、userID。
- 返回：error。非归属用户操作返回 error.New("无权操作")。

**单元测试用例**
- 用例1：成功删除
  - Given：归属用户。
  - When：调用 DeleteJournal。
  - Then：日志被删除，error 为 nil。
- 用例2：非归属用户
  - Given：userID 与日志归属不符。
  - When：调用 DeleteJournal。
  - Then：返回 error，且 error.Error() == "no permission"。

#### 1.4 添加emoji icon
**函数定义**
- 方法：func SetJournalIcon(req SetJournalIconRequest) (*Journal, error)
- 描述：为日志设置emoji icon。
- 参数：param - 包含 journalID、userID、icon。
- 返回：更新后的 Journal，或 error。非归属用户操作返回 error.New("无权操作")。

**单元测试用例**
- 用例1：成功设置
  - Given：归属用户。
  - When：调用 SetJournalIcon。
  - Then：icon 被正确设置，error 为 nil。
- 用例2：非归属用户
  - Given：userID 与日志归属不符。
  - When：调用 SetJournalIcon。
  - Then：返回 nil, error，且 error.Error() == "no permission"。

### Period（值对象）

#### 1.1 时间区间合法性
**函数定义**
- 方法：func NewPeriod(start, end time.Time) (Period, error)
- 描述：创建时间区间值对象，校验合法性。
- 参数：start、end。
- 返回：Period 值对象，或 error。
- 规则：start 必须早于 end。

**单元测试用例**
- 用例1：合法区间
  - Given：start < end。
  - When：调用 NewPeriod。
  - Then：返回 Period，error 为 nil。
- 用例2：非法区间
  - Given：start >= end。
  - When：调用 NewPeriod。
  - Then：返回 error，且 error.Error() == "invalid period range"。

### Task

#### 1.1 创建任务
**函数定义**
- 方法：func CreateTask(req CreateTaskRequest) (*Task, error)
- 描述：为指定用户创建一个新任务。
- 参数：param - 包含 userID、标题、周期、类型等信息的请求对象。
- 返回：新建的 Task 实体对象，或 error。
- 错误：当 userID 为空、period 非法、标题为空等，返回 error.New("userID不能为空")、error.New("时间区间非法")、error.New("标题不能为空") 等。

**单元测试用例**
- 用例1：成功创建
  - Given：CreateTaskRequest 各字段合法。
  - When：调用 CreateTask。
  - Then：返回新建的 Task，userID、period、标题等与请求一致，error 为 nil。
- 用例2：userID为空
  - Given：CreateTaskRequest.UserID 为空。
  - When：调用 CreateTask。
  - Then：返回 nil, error，且 error.Error() == "userID is required"。
- 用例3：period非法
  - Given：periodStart >= periodEnd。
  - When：调用 CreateTask。
  - Then：返回 nil, error，且 error.Error() == "invalid period range"。
- 用例4：标题为空
  - Given：标题为空字符串。
  - When：调用 CreateTask。
  - Then：返回 nil, error，且 error.Error() == "title is required"。

#### 1.2 编辑任务
**函数定义**
- 方法：func UpdateTask(req UpdateTaskRequest) (*Task, error)
- 描述：编辑指定任务，只有归属用户可操作。
- 参数：param - 包含 taskID、userID、可更新字段的请求对象。
- 返回：更新后的 Task 实体对象，或 error。
- 错误：非归属用户操作返回 error.New("无权操作")，period 非法返回 error.New("时间区间非法")。

**单元测试用例**
- 用例1：成功编辑
  - Given：归属用户，period合法。
  - When：调用 UpdateTask。
  - Then：返回更新后的 Task，字段变更生效，error 为 nil。
- 用例2：非归属用户
  - Given：userID与任务归属不符。
  - When：调用 UpdateTask。
  - Then：返回 nil, error，且 error.Error() == "no permission"。
- 用例3：已完成任务编辑
  - Given：任务已完成。
  - When：调用 UpdateTask。
  - Then：只能修改分数，其他字段变更无效或返回 error。

#### 1.3 删除任务
**函数定义**
- 方法：func DeleteTask(req DeleteTaskRequest) error
- 描述：删除指定任务，只有归属用户可操作。
- 参数：param - 包含 taskID、userID。
- 返回：error。非归属用户操作返回 error.New("无权操作")。

**单元测试用例**
- 用例1：成功删除
  - Given：归属用户。
  - When：调用 DeleteTask。
  - Then：任务被删除，error 为 nil。
- 用例2：非归属用户
  - Given：userID与任务归属不符。
  - When：调用 DeleteTask。
  - Then：返回 error，且 error.Error() == "no permission"。

#### 1.4 任务分数设置
**函数定义**
- 方法：func SetTaskScore(req SetTaskScoreRequest) (*Task, error)
- 描述：仅“日”类型任务可设置分数。
- 参数：param - 包含 taskID、userID、score。
- 返回：更新后的 Task，或 error。非“日”类型任务返回 error.New("仅日类型任务可设置分数")。

**单元测试用例**
- 用例1：成功设置
  - Given：日类型任务。
  - When：调用 SetTaskScore。
  - Then：分数被正确设置，error 为 nil。
- 用例2：非日类型任务
  - Given：周/月等类型任务。
  - When：调用 SetTaskScore。
  - Then：返回 nil, error，且 error.Error() == "only day type task can set score"。

#### 1.5 新建子任务
**函数定义**
- 方法：func CreateSubTask(req CreateSubTaskRequest) (*Task, error)
- 描述：为指定父任务创建子任务，userID必须与父任务一致。
- 参数：param - 包含 parentTaskID、userID、标题等。
- 返回：新建的子任务 Task，或 error。userID与父任务不一致返回 error.New("userID与父任务不一致")。

**单元测试用例**
- 用例1：成功创建
  - Given：userID与父任务一致。
  - When：调用 CreateSubTask。
  - Then：返回新建的子任务，error 为 nil。
- 用例2：userID与父任务不一致
  - Given：userID与父任务不一致。
  - When：调用 CreateSubTask。
  - Then：返回 nil, error，且 error.Error() == "userID does not match parent task"。

#### 1.6 标签操作
**函数定义**
- 方法：func AddTag(req AddTagRequest) (*Task, error)
-        func RemoveTag(req RemoveTagRequest) (*Task, error)
- 描述：为任务添加/移除标签。
- 参数：param - 包含 taskID、userID、tag。
- 返回：更新后的 Task，或 error。非归属用户操作返回 error.New("无权操作")。

**单元测试用例**
- 用例1：成功添加/移除
  - Given：归属用户。
  - When：调用 AddTag/RemoveTag。
  - Then：标签被正确添加/移除，error 为 nil。
- 用例2：非归属用户
  - Given：userID与任务归属不符。
  - When：调用 AddTag/RemoveTag。
  - Then：返回 nil, error，且 error.Error() == "no permission"。

#### 1.7 任务树/聚合查询
**函数定义**
- 方法：func GetTaskTree(req GetTaskTreeRequest) (*TaskTree, error)
- 描述：获取指定任务的整体任务树，校验权限。
- 参数：param - 包含 taskID、userID。
- 返回：任务树结构，或 error。非归属用户操作返回 error.New("无权操作")。

**单元测试用例**
- 用例1：成功获取
  - Given：归属用户。
  - When：调用 GetTaskTree。
  - Then：返回正确的任务树，error 为 nil。
- 用例2：非归属用户
  - Given：userID与任务归属不符。
  - When：调用 GetTaskTree。
  - Then：返回 nil, error，且 error.Error() == "no permission"。

#### 1.8 添加emoji icon
**函数定义**
- 方法：func SetTaskIcon(req SetTaskIconRequest) (*Task, error)
- 描述：为任务设置emoji icon。
- 参数：param - 包含 taskID、userID、icon。
- 返回：更新后的 Task，或 error。非归属用户操作返回 error.New("无权操作")。

**单元测试用例**
- 用例1：成功设置
  - Given：归属用户。
  - When：调用 SetTaskIcon。
  - Then：icon被正确设置，error 为 nil。
- 用例2：非归属用户
  - Given：userID与任务归属不符。
  - When：调用 SetTaskIcon。
  - Then：返回 nil, error，且 error.Error() == "no permission"。

