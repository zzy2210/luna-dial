package biz

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TaskStatus 任务状态枚举
type TaskStatus int

const (
	TaskStatusNotStarted TaskStatus = iota // 未开始
	TaskStatusInProgress                   // 进行中
	TaskStatusCompleted                    // 已完成
	TaskStatusCancelled                    // 已取消
)

// TaskPriority 任务优先级枚举
type TaskPriority int

const (
	TaskPriorityLow    TaskPriority = iota // 低
	TaskPriorityMedium                     // 中
	TaskPriorityHigh                       // 高
	TaskPriorityUrgent                     // 紧急
)

type Task struct {
	ID         string       `json:"id"`
	Title      string       `json:"title"`
	TaskType   PeriodType   `json:"type"`
	TimePeriod Period       `json:"period"`
	Tags       []string     `json:"tags"`
	Icon       string       `json:"icon"`
	Score      int          `json:"score"`
	Status     TaskStatus   `json:"status"`
	Priority   TaskPriority `json:"priority"`
	ParentID   string       `json:"parent_id"`
	UserID     string       `json:"user_id"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

// 创建任务参数
type CreateTaskParam struct {
	UserID   string
	Title    string
	Type     PeriodType
	Period   Period
	Tags     []string
	Icon     string
	Score    int
	Priority TaskPriority
	ParentID string
}

// 编辑任务参数
type UpdateTaskParam struct {
	TaskID string
	UserID string
	Title  *string
	// Type        *PeriodType 暂时不运行修改任务类型吧
	Period   *Period
	Tags     *[]string
	Icon     *string
	Score    *int
	Status   *TaskStatus
	Priority *TaskPriority
}

// 删除任务参数
type DeleteTaskParam struct {
	TaskID string
	UserID string
}

// 设置任务分数参数
type SetTaskScoreParam struct {
	TaskID string
	UserID string
	Score  int
}

// 创建子任务参数
type CreateSubTaskParam struct {
	ParentID string
	UserID   string
	Title    string
	Type     PeriodType
	Period   Period
	Tags     []string
	Icon     string
	Score    int
}

// 修改标签参数
type EditTagParam struct {
	TaskID string
	UserID string
	Tags   []string
}

// 设置任务icon参数
type SetTaskIconParam struct {
	TaskID string
	UserID string
	Icon   string
}

// 获取指定时间的指定类型的任务列表参数
type ListTaskByPeriodParam struct {
	UserID  string
	Period  Period
	GroupBy PeriodType
}

// 获取某个任务的父任务树列表参数
type ListTaskParentTreeParam struct {
	UserID string
	TaskID string
}

// 获取某个任务的整个任务树参数
type ListTaskTreeParam struct {
	UserID string
	TaskID string
}

type GetTaskStatsParam struct {
	UserID  string
	Period  Period
	GroupBy PeriodType
}

type TaskUsecase struct {
	repo TaskRepo
	// log *log.Helper
}

func NewTaskUsecase(repo TaskRepo) *TaskUsecase {
	return &TaskUsecase{repo: repo}
}

// 创建任务
// 必填 类型，时间，名称
func (uc *TaskUsecase) CreateTask(ctx context.Context, param CreateTaskParam) (*Task, error) {
	if param.UserID == "" || param.Title == "" {
		return nil, ErrInvalidInput // 参数不合法
	}

	tags := []string{}
	if param.Tags != nil {
		tags = param.Tags
	}
	if param.ParentID != "" {
		// 检查父任务是否存在
		parentTask, err := uc.repo.GetTask(ctx, param.ParentID, param.UserID)
		if err != nil {
			return nil, err // 返回仓库层的错误
		}
		if parentTask == nil {
			return nil, ErrInvalidInput // 父任务不存在
		}
		if len(parentTask.Tags) > 0 {
			tags = append(tags, parentTask.Tags...) // 继承父任务的标签
		}
	}
	if !param.Period.IsValid() {
		return nil, ErrInvalidInput // 时间段不合法
	}
	// 确保周期类型匹配
	if !param.Period.MatchesPeriodType(param.Type) {
		return nil, ErrInvalidInput // 周期类型不匹配
	}

	task := &Task{
		ID:         generateID(), // 假设有一个生成ID的函数
		Title:      param.Title,
		TaskType:   param.Type,
		TimePeriod: param.Period,
		Tags:       tags,
		Icon:       param.Icon,
		Score:      param.Score,
		Status:     TaskStatusNotStarted, // 默认状态为未开始
		Priority:   param.Priority,       // 使用传入的优先级，如果为0则默认为低优先级
		UserID:     param.UserID,
		ParentID:   param.ParentID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := uc.repo.CreateTask(ctx, task)
	if err != nil {
		return nil, err // 返回仓库层的错误
	}
	return task, nil
}

// 更新任务
// 禁止修改 period type，其实不难做，但是先不放开
// 如果传递了标题，必须不为空
// 如果传递了时间段，必须不为空且合法
func (uc *TaskUsecase) UpdateTask(ctx context.Context, param UpdateTaskParam) (*Task, error) {
	if param.TaskID == "" || param.UserID == "" {
		return nil, ErrInvalidInput // 参数不合法
	}
	// title 如果传递了，必须不为空
	if param.Title != nil && *param.Title == "" {
		return nil, ErrInvalidInput // 标题不能为空
	}

	if param.Period != nil && (param.Period.Start.IsZero() || param.Period.End.IsZero()) {
		return nil, ErrInvalidInput // 时间段不能为空
	}
	// 校验period 是否合法
	if param.Period != nil && !param.Period.IsValid() {
		return nil, ErrInvalidInput // 时间段不合法
	}

	task, err := uc.repo.GetTask(ctx, param.TaskID, param.UserID)
	if err != nil {
		return nil, err // 返回仓库层的错误
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	task.UpdatedAt = time.Now()
	if param.Title != nil {
		task.Title = *param.Title
	}
	if param.Period != nil {
		task.TimePeriod = *param.Period
	}
	if param.Tags != nil {
		task.Tags = *param.Tags
	}
	if param.Icon != nil {
		task.Icon = *param.Icon
	}
	if param.Score != nil {
		task.Score = *param.Score
	}
	if param.Status != nil {
		task.Status = *param.Status
	}
	if param.Priority != nil {
		task.Priority = *param.Priority
	}
	err = uc.repo.UpdateTask(ctx, task)
	if err != nil {
		return nil, err // 返回仓库层的错误
	}
	return task, nil
}

// 根据ID删除任务
// 检查USERID
func (uc *TaskUsecase) DeleteTask(ctx context.Context, param DeleteTaskParam) error {
	if param.TaskID == "" || param.UserID == "" {
		return ErrInvalidInput // 参数不合法
	}

	// 检查任务是否存在并且属于该用户
	task, err := uc.repo.GetTask(ctx, param.TaskID, param.UserID)
	if err != nil {
		return err // 返回仓库层的错误
	}
	if task == nil {
		return ErrTaskNotFound // 任务不存在
	}

	// 删除任务
	err = uc.repo.DeleteTask(ctx, param.TaskID, param.UserID)
	if err != nil {
		return err // 返回仓库层的错误
	}

	return nil
}

// 根据ID更新任务分数
// 检查USERID
func (uc *TaskUsecase) SetTaskScore(ctx context.Context, param SetTaskScoreParam) (*Task, error) {
	if param.TaskID == "" || param.UserID == "" {
		return nil, ErrInvalidInput // 参数不合法
	}

	if param.Score < 0 {
		return nil, ErrInvalidInput // 分数不能为负数
	}

	// 检查任务是否存在并且属于该用户
	task, err := uc.repo.GetTask(ctx, param.TaskID, param.UserID)
	if err != nil {
		return nil, err // 返回仓库层的错误
	}
	if task == nil {
		return nil, ErrTaskNotFound // 任务不存在
	}

	// 更新分数
	task.Score = param.Score
	task.UpdatedAt = time.Now()

	err = uc.repo.UpdateTask(ctx, task)
	if err != nil {
		return nil, err // 返回仓库层的错误
	}

	return task, nil
}

// 创建子任务
// 检查USERID和ParentID
// 子任务的时间段和类型必须与父任务一致
func (uc *TaskUsecase) CreateSubTask(ctx context.Context, param CreateSubTaskParam) (*Task, error) {
	if param.ParentID == "" || param.UserID == "" || param.Title == "" {
		return nil, ErrInvalidInput // 参数不合法
	}

	// 检查父任务是否存在并且属于该用户
	parentTask, err := uc.repo.GetTask(ctx, param.ParentID, param.UserID)
	if err != nil {
		return nil, err // 返回仓库层的错误
	}
	if parentTask == nil {
		return nil, ErrTaskNotFound // 父任务不存在
	}

	// 验证子任务的时间段和类型必须与父任务一致
	if param.Type != parentTask.TaskType {
		return nil, ErrInvalidInput // 子任务类型必须与父任务一致
	}

	if !param.Period.IsValid() {
		return nil, ErrInvalidInput // 时间段不合法
	}

	// 子任务的时间段必须在父任务时间段内
	if param.Period.Start.Before(parentTask.TimePeriod.Start) || param.Period.End.After(parentTask.TimePeriod.End) {
		return nil, ErrInvalidInput // 子任务时间段必须在父任务时间段内
	}

	// 继承父任务的标签
	tags := param.Tags
	if len(parentTask.Tags) > 0 {
		tags = append(tags, parentTask.Tags...)
	}

	// 创建子任务
	task := &Task{
		ID:         generateID(),
		Title:      param.Title,
		TaskType:   param.Type,
		TimePeriod: param.Period,
		Tags:       tags,
		Icon:       param.Icon,
		Score:      param.Score,
		Status:     TaskStatusNotStarted, // 子任务默认状态为未开始
		Priority:   parentTask.Priority,  // 等级与父任务一致
		UserID:     param.UserID,
		ParentID:   param.ParentID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = uc.repo.CreateTask(ctx, task)
	if err != nil {
		return nil, err // 返回仓库层的错误
	}

	return task, nil
}

// 修改标签 - 直接覆盖替换任务的所有标签
func (uc *TaskUsecase) EditTag(ctx context.Context, param EditTagParam) (*Task, error) {
	if param.TaskID == "" || param.UserID == "" {
		return nil, ErrInvalidInput // 参数不合法
	}

	// 检查任务是否存在并且属于该用户
	task, err := uc.repo.GetTask(ctx, param.TaskID, param.UserID)
	if err != nil {
		return nil, err // 返回仓库层的错误
	}
	if task == nil {
		return nil, ErrTaskNotFound // 任务不存在
	}

	// 直接覆盖替换所有标签
	task.Tags = param.Tags
	task.UpdatedAt = time.Now()

	err = uc.repo.UpdateTask(ctx, task)
	if err != nil {
		return nil, err // 返回仓库层的错误
	}

	return task, nil
}

// 设置icon（覆盖写）
func (uc *TaskUsecase) SetTaskIcon(ctx context.Context, param SetTaskIconParam) (*Task, error) {
	if param.TaskID == "" || param.UserID == "" {
		return nil, ErrInvalidInput // 参数不合法
	}

	// 检查任务是否存在并且属于该用户
	task, err := uc.repo.GetTask(ctx, param.TaskID, param.UserID)
	if err != nil {
		return nil, err // 返回仓库层的错误
	}
	if task == nil {
		return nil, ErrTaskNotFound // 任务不存在
	}

	// 设置图标
	task.Icon = param.Icon
	task.UpdatedAt = time.Now()

	err = uc.repo.UpdateTask(ctx, task)
	if err != nil {
		return nil, err // 返回仓库层的错误
	}

	return task, nil
}

// 按照时间周期获取任务
// xx时间段内，xx类型的任务
func (uc *TaskUsecase) ListTaskByPeriod(ctx context.Context, param ListTaskByPeriodParam) ([]Task, error) {
	if param.UserID == "" {
		return nil, ErrInvalidInput // 参数不合法
	}

	if !param.Period.IsValid() {
		return nil, ErrInvalidInput // 时间段不合法
	}

	// 调用仓库层获取任务列表
	// 注意：这里假设 GroupBy 参数用于过滤任务类型
	taskType := int(param.GroupBy)
	tasks, err := uc.repo.ListTasks(ctx, param.UserID, param.Period.Start, param.Period.End, taskType)
	if err != nil {
		return nil, err // 返回仓库层的错误
	}

	// 将 []*Task 转换为 []Task
	result := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		if task != nil {
			result = append(result, *task)
		}
	}

	return result, nil
}

// 获取某个任务的任务树 (整个树)
func (uc *TaskUsecase) ListTaskTree(ctx context.Context, param ListTaskTreeParam) ([]Task, error) {
	if param.TaskID == "" || param.UserID == "" {
		return nil, ErrInvalidInput // 参数不合法
	}

	// 首先检查根任务是否存在并且属于该用户
	rootTask, err := uc.repo.GetTask(ctx, param.TaskID, param.UserID)
	if err != nil {
		return nil, err // 返回仓库层的错误
	}
	if rootTask == nil {
		return nil, ErrTaskNotFound // 任务不存在
	}

	// 调用仓库层获取任务树
	tasks, err := uc.repo.ListTaskTree(ctx, param.TaskID, param.UserID)
	if err != nil {
		return nil, err // 返回仓库层的错误
	}

	// 将 []*Task 转换为 []Task
	result := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		if task != nil {
			result = append(result, *task)
		}
	}

	return result, nil
}

// 获取某个任务的父任务树 (从根节点到该任务)
func (uc *TaskUsecase) ListTaskParentTree(ctx context.Context, param ListTaskParentTreeParam) ([]Task, error) {
	if param.TaskID == "" || param.UserID == "" {
		return nil, ErrInvalidInput // 参数不合法
	}

	// 首先检查任务是否存在并且属于该用户
	currentTask, err := uc.repo.GetTask(ctx, param.TaskID, param.UserID)
	if err != nil {
		return nil, err // 返回仓库层的错误
	}
	if currentTask == nil {
		return nil, ErrTaskNotFound // 任务不存在
	}

	// 调用仓库层获取父任务树
	tasks, err := uc.repo.ListTaskParentTree(ctx, param.TaskID, param.UserID)
	if err != nil {
		return nil, err // 返回仓库层的错误
	}

	// 将 []*Task 转换为 []Task
	result := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		if task != nil {
			result = append(result, *task)
		}
	}

	return result, nil
}

// 获取某个时间范围内，所有period 为day 的任务集合
// group by 代表分组方式 日/周/月/季度/年
func (uc *TaskUsecase) GetTaskStats(ctx context.Context, param GetTaskStatsParam) ([]GroupStat, error) {
	if param.UserID == "" {
		return nil, ErrInvalidInput // 参数不合法
	}

	if !param.Period.IsValid() {
		return nil, ErrInvalidInput // 时间段不合法
	}

	// 获取指定时间范围内的所有日任务
	tasks, err := uc.repo.ListTasks(ctx, param.UserID, param.Period.Start, param.Period.End, int(PeriodDay))
	if err != nil {
		return nil, err // 返回仓库层的错误
	}

	// 按照 GroupBy 参数进行分组统计
	statsMap := make(map[string]*GroupStat)

	for _, task := range tasks {
		if task == nil {
			continue
		}

		// 根据分组方式生成分组键
		groupKey := uc.generateGroupKey(task.TimePeriod.Start, param.GroupBy)

		// 如果该分组不存在，创建新的统计对象
		if _, exists := statsMap[groupKey]; !exists {
			statsMap[groupKey] = &GroupStat{
				GroupKey:   groupKey,
				TaskCount:  0,
				ScoreTotal: 0,
			}
		}

		// 累加统计数据
		statsMap[groupKey].TaskCount++
		statsMap[groupKey].ScoreTotal += task.Score
	}

	// 将 map 转换为切片
	var result []GroupStat
	for _, stat := range statsMap {
		result = append(result, *stat)
	}

	return result, nil
}

// generateGroupKey 根据时间和分组类型生成分组键
func (uc *TaskUsecase) generateGroupKey(t time.Time, groupBy PeriodType) string {
	switch groupBy {
	case PeriodDay:
		return t.Format("2006-01-02") // 2025-01-15
	case PeriodWeek:
		year, week := t.ISOWeek()
		return fmt.Sprintf("%d-W%02d", year, week) // 2025-W03
	case PeriodMonth:
		return t.Format("2006-01") // 2025-01
	case PeriodQuarter:
		quarter := (int(t.Month())-1)/3 + 1
		return fmt.Sprintf("%d-Q%d", t.Year(), quarter) // 2025-Q1
	case PeriodYear:
		return t.Format("2006") // 2025
	default:
		return t.Format("2006-01-02") // 默认按日分组
	}
}

func generateID() string {
	// 生成UUID并去除连字符，符合项目规范
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}
