package service

import (
	"context"
	"fmt"
	"math"

	"okr-web/ent"
	"okr-web/ent/journalentry"
	taskent "okr-web/ent/task"
	"okr-web/internal/repository"
	"okr-web/internal/types"

	"github.com/google/uuid"
)

// TaskServiceImpl 任务服务实现
type TaskServiceImpl struct {
	taskRepo    repository.TaskRepository
	userRepo    repository.UserRepository
	journalRepo repository.JournalRepository
}

// NewTaskService 创建任务服务
func NewTaskService(taskRepo repository.TaskRepository, userRepo repository.UserRepository, journalRepo repository.JournalRepository) TaskService {
	return &TaskServiceImpl{
		taskRepo:    taskRepo,
		userRepo:    userRepo,
		journalRepo: journalRepo,
	}
}

// CreateTask 创建任务
func (s *TaskServiceImpl) CreateTask(ctx context.Context, userID uuid.UUID, req TaskRequest) (*ent.Task, error) {
	// 调试：打印收到的请求体
	fmt.Printf("[DEBUG] CreateTask received req: %+v\n", req)
	// 验证用户是否存在
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, &types.AppError{
			Code:    404,
			Message: "用户不存在",
			Type:    "USER_NOT_FOUND",
		}
	}

	// 验证任务类型
	if !req.Type.IsValid() {
		return nil, &types.AppError{
			Code:    400,
			Message: "无效的任务类型",
			Type:    "INVALID_TASK_TYPE",
		}
	}

	// 验证任务状态
	if !req.Status.IsValid() {
		return nil, &types.AppError{
			Code:    400,
			Message: "无效的任务状态",
			Type:    "INVALID_TASK_STATUS",
		}
	}

	// 验证评分范围
	if req.Score != nil && (*req.Score < 0 || *req.Score > 10) {
		return nil, &types.AppError{
			Code:    400,
			Message: "评分必须在0-10之间",
			Type:    "INVALID_SCORE",
		}
	}

	// 验证时间范围
	if req.StartDate != nil && req.EndDate != nil && req.StartDate.After(*req.EndDate) {
		return nil, &types.AppError{
			Code:    400,
			Message: "开始时间不能晚于结束时间",
			Type:    "INVALID_TIME_RANGE",
		}
	}

	// 创建任务
	task, err := s.taskRepo.Create(ctx, func(create *ent.TaskCreate) *ent.TaskCreate {
		create = create.
			SetTitle(req.Title).
			SetType(taskent.Type(req.Type)).
			SetStatus(taskent.Status(req.Status)).
			SetUserID(userID)

		if req.Description != nil {
			create = create.SetDescription(*req.Description)
		}
		if req.StartDate != nil {
			create = create.SetStartDate(*req.StartDate)
		}
		if req.EndDate != nil {
			create = create.SetEndDate(*req.EndDate)
		}
		if req.Score != nil {
			create = create.SetScore(*req.Score)
		}
		if req.Tags != nil {
			create = create.SetTags(*req.Tags)
		}

		return create
	})

	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "任务创建失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return task, nil
}

// GetTask 获取任务
func (s *TaskServiceImpl) GetTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*ent.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, &types.AppError{
			Code:    404,
			Message: "任务不存在",
			Type:    "TASK_NOT_FOUND",
		}
	}

	// 检查权限
	if task.UserID != userID {
		return nil, &types.AppError{
			Code:    403,
			Message: "无权限访问此任务",
			Type:    "FORBIDDEN",
		}
	}

	return task, nil
}

// UpdateTask 更新任务
func (s *TaskServiceImpl) UpdateTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, req TaskRequest) (*ent.Task, error) {
	// 先获取任务并验证权限
	_, err := s.GetTask(ctx, userID, taskID)
	if err != nil {
		return nil, err
	}

	// 验证任务类型
	if !req.Type.IsValid() {
		return nil, &types.AppError{
			Code:    400,
			Message: "无效的任务类型",
			Type:    "INVALID_TASK_TYPE",
		}
	}

	// 验证任务状态
	if !req.Status.IsValid() {
		return nil, &types.AppError{
			Code:    400,
			Message: "无效的任务状态",
			Type:    "INVALID_TASK_STATUS",
		}
	}

	// 验证评分范围
	if req.Score != nil && (*req.Score < 0 || *req.Score > 10) {
		return nil, &types.AppError{
			Code:    400,
			Message: "评分必须在0-10之间",
			Type:    "INVALID_SCORE",
		}
	}

	// 验证时间范围
	if req.StartDate != nil && req.EndDate != nil && req.StartDate.After(*req.EndDate) {
		return nil, &types.AppError{
			Code:    400,
			Message: "开始时间不能晚于结束时间",
			Type:    "INVALID_TIME_RANGE",
		}
	}

	// 更新任务
	updatedTask, err := s.taskRepo.Update(ctx, taskID, func(update *ent.TaskUpdateOne) *ent.TaskUpdateOne {
		update = update.
			SetTitle(req.Title).
			SetType(taskent.Type(req.Type)).
			SetStatus(taskent.Status(req.Status))

		if req.Description != nil {
			update = update.SetDescription(*req.Description)
		}
		if req.StartDate != nil {
			update = update.SetStartDate(*req.StartDate)
		}
		if req.EndDate != nil {
			update = update.SetEndDate(*req.EndDate)
		}
		if req.Score != nil {
			update = update.SetScore(*req.Score)
		}
		if req.Tags != nil {
			update = update.SetTags(*req.Tags)
		}

		return update
	})

	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "任务更新失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return updatedTask, nil
}

// DeleteTask 删除任务
func (s *TaskServiceImpl) DeleteTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) error {
	// 先获取任务并验证权限
	_, err := s.GetTask(ctx, userID, taskID)
	if err != nil {
		return err
	}

	// 检查是否有子任务
	children, err := s.taskRepo.GetChildren(ctx, taskID)
	if err != nil {
		return &types.AppError{
			Code:    500,
			Message: "检查子任务失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	if len(children) > 0 {
		return &types.AppError{
			Code:    400,
			Message: "无法删除有子任务的任务",
			Type:    "TASK_HAS_CHILDREN",
		}
	}

	// 删除任务
	err = s.taskRepo.Delete(ctx, taskID)
	if err != nil {
		return &types.AppError{
			Code:    500,
			Message: "任务删除失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return nil
}

// GetTaskChildren 获取子任务
func (s *TaskServiceImpl) GetTaskChildren(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) ([]*ent.Task, error) {
	// 先验证父任务权限
	_, err := s.GetTask(ctx, userID, taskID)
	if err != nil {
		return nil, err
	}

	// 获取子任务
	children, err := s.taskRepo.GetChildren(ctx, taskID)
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "获取子任务失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return children, nil
}

// GetTaskTree 获取任务树
func (s *TaskServiceImpl) GetTaskTree(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*types.TaskTree, error) {
	// 先验证任务权限
	rootTask, err := s.GetTask(ctx, userID, taskID)
	if err != nil {
		return nil, err
	}

	// 构建任务树
	tree, err := s.buildTaskTree(ctx, rootTask)
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "构建任务树失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return tree, nil
}

// buildTaskTree 递归构建任务树
func (s *TaskServiceImpl) buildTaskTree(ctx context.Context, task *ent.Task) (*types.TaskTree, error) {
	tree := &types.TaskTree{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Type:        task.Type.String(),
		StartDate:   task.StartDate.Format("2006-01-02"),
		EndDate:     task.EndDate.Format("2006-01-02"),
		Status:      task.Status.String(),
		Score:       task.Score,
		ParentID:    task.ParentID,
		UserID:      task.UserID,
		Tags:        task.Tags,
		CreatedAt:   task.CreatedAt.Format("2006-01-02"),
		UpdatedAt:   task.UpdatedAt.Format("2006-01-02"),
		Ancestors:   nil, // 如有需要可递归补充
		Children:    make([]*types.TaskTree, 0),
		Depth:       0, // 如有需要可补充
	}

	// 获取子任务
	children, err := s.taskRepo.GetChildren(ctx, task.ID)
	if err != nil {
		return nil, err
	}

	// 递归构建子树
	for _, child := range children {
		childTree, err := s.buildTaskTree(ctx, child)
		if err != nil {
			return nil, err
		}
		tree.Children = append(tree.Children, childTree)
	}

	return tree, nil
}

// GetGlobalView 获取全局视图
func (s *TaskServiceImpl) GetGlobalView(ctx context.Context, userID uuid.UUID) ([]*types.TaskTree, error) {
	// 获取所有根任务
	rootTasks, err := s.taskRepo.GetRootTasks(ctx, userID)
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "获取根任务失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// 构建所有任务树
	trees := make([]*types.TaskTree, 0, len(rootTasks))
	for _, rootTask := range rootTasks {
		tree, err := s.buildTaskTree(ctx, rootTask)
		if err != nil {
			return nil, err
		}
		trees = append(trees, tree)
	}

	return trees, nil
}

// UpdateTaskScore 更新任务分数
func (s *TaskServiceImpl) UpdateTaskScore(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, score int) error {
	// 验证分数范围
	if score < 0 || score > 10 {
		return &types.AppError{
			Code:    400,
			Message: "评分必须在0-10之间",
			Type:    "INVALID_SCORE",
		}
	}

	// 先验证任务权限
	_, err := s.GetTask(ctx, userID, taskID)
	if err != nil {
		return err
	}

	// 更新分数
	_, err = s.taskRepo.Update(ctx, taskID, func(update *ent.TaskUpdateOne) *ent.TaskUpdateOne {
		return update.SetScore(score)
	})

	if err != nil {
		return &types.AppError{
			Code:    500,
			Message: "更新任务分数失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return nil
}

// GetTasksByUser 获取用户的任务列表（分页）
func (s *TaskServiceImpl) GetTasksByUser(ctx context.Context, userID uuid.UUID, filters TaskFilters) (*TaskListResponse, error) {
	// 验证分页参数
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 || filters.PageSize > 100 {
		filters.PageSize = 20
	}

	// 计算偏移量
	offset := (filters.Page - 1) * filters.PageSize

	// 获取任务列表（这里需要repository支持复杂查询）
	tasks, err := s.taskRepo.GetByUserID(ctx, userID, filters.PageSize, offset)
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "获取任务列表失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// 获取总数
	total, err := s.taskRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "获取任务总数失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(filters.PageSize)))

	return &TaskListResponse{
		Tasks:       tasks,
		Total:       int64(total),
		CurrentPage: filters.Page,
		PageSize:    filters.PageSize,
		TotalPages:  totalPages,
	}, nil
}

// CreateSubTask 创建子任务
func (s *TaskServiceImpl) CreateSubTask(ctx context.Context, userID uuid.UUID, parentID uuid.UUID, req TaskRequest) (*ent.Task, error) {
	// 验证父任务权限
	_, err := s.GetTask(ctx, userID, parentID)
	if err != nil {
		return nil, err
	}

	// 创建子任务
	task, err := s.taskRepo.Create(ctx, func(create *ent.TaskCreate) *ent.TaskCreate {
		create = create.
			SetTitle(req.Title).
			SetType(taskent.Type(req.Type)).
			SetStatus(taskent.Status(req.Status)).
			SetUserID(userID).
			SetParentID(parentID)

		if req.Description != nil {
			create = create.SetDescription(*req.Description)
		}
		if req.StartDate != nil {
			create = create.SetStartDate(*req.StartDate)
		}
		if req.EndDate != nil {
			create = create.SetEndDate(*req.EndDate)
		}
		if req.Score != nil {
			create = create.SetScore(*req.Score)
		}
		if req.Tags != nil {
			create = create.SetTags(*req.Tags)
		}

		return create
	})

	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "子任务创建失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return task, nil
}

// GetContextView 获取上下文视图
func (s *TaskServiceImpl) GetContextView(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*ContextView, error) {
	// 先验证任务权限
	currentTask, err := s.GetTask(ctx, userID, taskID)
	if err != nil {
		return nil, err
	}

	contextView := &ContextView{
		Current: currentTask,
	}

	// 获取父任务
	if currentTask.ParentID != nil {
		parentTask, err := s.taskRepo.GetByID(ctx, *currentTask.ParentID)
		if err == nil {
			contextView.Parent = parentTask
		}
	}

	// 获取子任务
	children, err := s.taskRepo.GetChildren(ctx, taskID)
	if err == nil {
		contextView.Children = children
	}

	// 获取兄弟任务
	if currentTask.ParentID != nil {
		siblings, err := s.taskRepo.GetChildren(ctx, *currentTask.ParentID)
		if err == nil {
			// 过滤掉自己
			filteredSiblings := make([]*ent.Task, 0, len(siblings))
			for _, sibling := range siblings {
				if sibling.ID != taskID {
					filteredSiblings = append(filteredSiblings, sibling)
				}
			}
			contextView.Siblings = filteredSiblings
		}
	}

	return contextView, nil
}

// GetPlanView 获取计划视图
func (s *TaskServiceImpl) GetPlanView(ctx context.Context, userID uuid.UUID, req PlanRequest) (*types.PlanResponse, error) {
	// 验证用户是否存在
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, &types.AppError{
			Code:    404,
			Message: "用户不存在",
			Type:    "USER_NOT_FOUND",
		}
	}

	// 解析时间参考
	timeRange, err := types.ParseTimeReference(req.TimeRef, req.Scale)
	if err != nil {
		return nil, &types.AppError{
			Code:    400,
			Message: "无效的时间参考: " + err.Error(),
			Type:    "INVALID_TIME_REFERENCE",
		}
	}

	   // 获取时间范围内的任务及其完整父级链（按类型过滤，type=scale）
	   tasks, err := s.taskRepo.GetTasksWithAncestors(ctx, userID, timeRange.Start, timeRange.End, string(req.Scale))
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "获取任务失败: " + err.Error(),
			Type:    "DATABASE_ERROR",
		}
	}

	// 构建任务树结构
	taskTrees := s.buildTaskTrees(tasks)

	// 生成本周期所有 time_reference 列表
	timeRefs := types.GenerateTimeLabels(timeRange, req.Scale)

	// 将 types.TimeScale 转为 ent/journalentry.TimeScale
	var journalTimeScale journalentry.TimeScale
	switch req.Scale {
	case types.TimeScaleDay:
		journalTimeScale = journalentry.TimeScaleDay
	case types.TimeScaleWeek:
		journalTimeScale = journalentry.TimeScaleWeek
	case types.TimeScaleMonth:
		journalTimeScale = journalentry.TimeScaleMonth
	case types.TimeScaleQuarter:
		journalTimeScale = journalentry.TimeScaleQuarter
	case types.TimeScaleYear:
		journalTimeScale = journalentry.TimeScaleYear
	default:
		return nil, &types.AppError{
			Code:    400,
			Message: "无效的时间尺度",
			Type:    "INVALID_TIME_SCALE",
		}
	}

	// 获取本周期所有相关日志（按 time_scale + time_reference 批量）
	journals, err := s.journalRepo.GetByTimeScaleAndReferences(ctx, userID, journalTimeScale, timeRefs)
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "获取日志失败: " + err.Error(),
			Type:    "DATABASE_ERROR",
		}
	}

	response := &types.PlanResponse{
		Tasks:     taskTrees, // taskTrees 应为 []*types.TaskTree
		Journals:  journals,
		TimeRange: timeRange,
		Stats:     &types.PlanStats{},
	}

	// 填充统计数据
	stats := &types.PlanStats{}
	for _, t := range tasks {
		stats.TotalTasks++
		stats.TotalScore += t.Score
		switch t.Status {
		case "completed":
			stats.CompletedTasks++
			stats.CompletedScore += t.Score
		case "in-progress":
			stats.InProgressTasks++
		case "pending":
			stats.PendingTasks++
		}
	}
	response.Stats = stats

	return response, nil
}

// buildTaskTrees 构建任务树结构
func (s *TaskServiceImpl) buildTaskTrees(tasks []*ent.Task) []*types.TaskTree {
	// 构建任务映射
	taskMap := make(map[uuid.UUID]*ent.Task)
	for _, task := range tasks {
		taskMap[task.ID] = task
	}

	// 构建任务树节点映射
	treeMap := make(map[uuid.UUID]*types.TaskTree)
	for _, task := range tasks {
		treeMap[task.ID] = &types.TaskTree{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Type:        string(task.Type),
			StartDate:   task.StartDate.Format("2006-01-02T15:04:05Z"),
			EndDate:     task.EndDate.Format("2006-01-02T15:04:05Z"),
			Status:      string(task.Status),
			Score:       task.Score,
			ParentID:    task.ParentID,
			UserID:      task.UserID,
			Tags:        task.Tags,
			CreatedAt:   task.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   task.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			Children:    make([]*types.TaskTree, 0),
		}
	}

	// 构建父子关系
	var roots []*types.TaskTree
	for _, task := range tasks {
		treeNode := treeMap[task.ID]

		if task.ParentID != nil {
			if parentTree, exists := treeMap[*task.ParentID]; exists {
				parentTree.Children = append(parentTree.Children, treeNode)
			} else {
				roots = append(roots, treeNode)
			}
		} else {
			roots = append(roots, treeNode)
		}
	}

	return roots
}
