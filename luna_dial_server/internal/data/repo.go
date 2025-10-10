package data

import (
    "context"
    "errors"
    "fmt"
    "luna_dial/internal/biz"
    "luna_dial/internal/model"
    "time"

    "gorm.io/gorm"
)

// TaskRepo 任务仓库实现
type taskRepo struct {
	db        *gorm.DB
	converter *TaskConverter
}

func NewTaskRepo(db *gorm.DB) biz.TaskRepo {
	return &taskRepo{
		db:        db,
		converter: NewTaskConverter(),
	}
}

func (r *taskRepo) CreateTask(ctx context.Context, bizTask *biz.Task) error {
	dataTask := r.converter.BizToData(bizTask)
	return r.db.WithContext(ctx).Create(dataTask).Error
}

func (r *taskRepo) GetTask(ctx context.Context, taskID, userID string) (*biz.Task, error) {
    var dataTask Task
    err := r.db.WithContext(ctx).
        Where("id = ? AND user_id = ?", taskID, userID).
        First(&dataTask).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // 对于任务：返回 nil, nil，让上层根据 nil 判断为未找到
            return nil, nil
        }
        return nil, err
    }

    return r.converter.DataToBiz(&dataTask), nil
}

func (r *taskRepo) UpdateTask(ctx context.Context, bizTask *biz.Task) error {
	dataTask := r.converter.BizToData(bizTask)
	return r.db.WithContext(ctx).Save(dataTask).Error
}

func (r *taskRepo) DeleteTask(ctx context.Context, taskID, userID string) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", taskID, userID).
		Delete(&Task{}).Error
}

func (r *taskRepo) ListTasks(ctx context.Context, userID string, periodStart, periodEnd time.Time, taskType int) ([]*biz.Task, error) {
	var dataTasks []*Task
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND task_type = ? AND period_start >= ? AND period_end <= ?",
			userID, taskType, periodStart, periodEnd).
		Find(&dataTasks).Error

	if err != nil {
		return nil, err
	}

	return r.converter.DataToBizList(dataTasks), nil
}

// buildTreeStructure 在内存中构建树形结构
// 输入：已按 tree_depth 排序的任务列表
// 输出：构建好父子关系的任务树
func (r *taskRepo) buildTreeStructure(tasks []*biz.Task) []*biz.Task {
	if len(tasks) == 0 {
		return tasks
	}

	// 创建任务映射表，便于快速查找
	taskMap := make(map[string]*biz.Task)
	for _, task := range tasks {
		// 初始化 Children 切片
		task.Children = make([]*biz.Task, 0)
		taskMap[task.ID] = task
	}

	// 构建父子关系
	var rootTasks []*biz.Task
	for _, task := range tasks {
		if task.ParentID == "" {
			// 根任务
			rootTasks = append(rootTasks, task)
		} else {
			// 子任务：添加到父任务的 Children 中
			if parent, exists := taskMap[task.ParentID]; exists {
				parent.Children = append(parent.Children, task)
			}
		}
	}

	return rootTasks
}

func (r *taskRepo) ListTaskParentTree(ctx context.Context, taskID, userID string) ([]*biz.Task, error) {
	// 递归查询父任务树
	var dataTasks []*Task

	// 从当前任务开始，向上查询父任务
	currentTaskID := taskID
	for currentTaskID != "" {
		var task Task
		err := r.db.WithContext(ctx).
			Where("id = ? AND user_id = ?", currentTaskID, userID).
			First(&task).Error
		if err != nil {
			break
		}

		dataTasks = append(dataTasks, &task)
		currentTaskID = task.ParentID
	}

	return r.converter.DataToBizList(dataTasks), nil
}

// ListRootTasksWithPagination 分页查询根任务
// 用于全局任务树视图的第一步：获取根任务列表
func (r *taskRepo) ListRootTasksWithPagination(ctx context.Context, userID string, page, pageSize int, includeStatus []biz.TaskStatus) ([]*biz.Task, int64, error) {
	// 构建查询条件
	query := r.db.WithContext(ctx).Model(&Task{}).
		Where("user_id = ? AND (parent_id IS NULL OR parent_id = '')", userID)

	// 状态过滤：默认排除已取消的任务
	if len(includeStatus) > 0 {
		statusInts := make([]int, len(includeStatus))
		for i, status := range includeStatus {
			statusInts[i] = int(status)
		}
		query = query.Where("status IN ?", statusInts)
	} else {
		// 默认排除已取消状态(3)
		query = query.Where("status != ?", int(biz.TaskStatusCancelled))
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var dataTasks []*Task
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&dataTasks).Error

	if err != nil {
		return nil, 0, err
	}

	return r.converter.DataToBizList(dataTasks), total, nil
}

// ListTasksByRootIDs 根据根任务ID列表批量查询子任务
// 用于全局任务树视图的第二步：批量获取每个根任务的完整子树
func (r *taskRepo) ListTasksByRootIDs(ctx context.Context, userID string, rootTaskIDs []string, includeStatus []biz.TaskStatus) ([]*biz.Task, error) {
	if len(rootTaskIDs) == 0 {
		return []*biz.Task{}, nil
	}

	// 构建查询条件
	query := r.db.WithContext(ctx).
		Where("user_id = ? AND root_task_id IN ?", userID, rootTaskIDs)

	// 状态过滤
	if len(includeStatus) > 0 {
		statusInts := make([]int, len(includeStatus))
		for i, status := range includeStatus {
			statusInts[i] = int(status)
		}
		query = query.Where("status IN ?", statusInts)
	} else {
		// 默认排除已取消状态(3)
		query = query.Where("status != ?", int(biz.TaskStatusCancelled))
	}

	var dataTasks []*Task
	err := query.Order("root_task_id, tree_depth, created_at").
		Find(&dataTasks).Error

	if err != nil {
		return nil, err
	}

	// 转换为业务模型并构建树结构
	bizTasks := r.converter.DataToBizList(dataTasks)
	return r.buildTreeStructure(bizTasks), nil
}

// GetCompleteTaskTree 获取包含指定任务的完整任务树
// 支持传入任意层级的任务ID，先找到根任务，然后返回完整树
func (r *taskRepo) GetCompleteTaskTree(ctx context.Context, taskID, userID string, includeStatus []biz.TaskStatus) ([]*biz.Task, error) {
	// 步骤1：获取指定任务的根任务ID
	var rootTaskID string
	err := r.db.WithContext(ctx).Model(&Task{}).
		Select("root_task_id").
		Where("id = ? AND user_id = ?", taskID, userID).
		Scan(&rootTaskID).Error
	if err != nil {
		return nil, err
	}

	// 步骤2：获取完整的任务树
	query := r.db.WithContext(ctx).
		Where("user_id = ? AND (id = ? OR root_task_id = ?)", userID, rootTaskID, rootTaskID)

	// 状态过滤（父任务链查询时包含所有状态，便于理解完整层级关系）
	if len(includeStatus) > 0 {
		statusInts := make([]int, len(includeStatus))
		for i, status := range includeStatus {
			statusInts[i] = int(status)
		}
		query = query.Where("status IN ?", statusInts)
	}

	var dataTasks []*Task
	err = query.Order("tree_depth, created_at").Find(&dataTasks).Error
	if err != nil {
		return nil, err
	}

	// 转换为业务模型并构建树结构
	bizTasks := r.converter.DataToBizList(dataTasks)
	return r.buildTreeStructure(bizTasks), nil
}

// GetTaskParentChain 获取任务的父级链路
// 从指定任务开始向上遍历，返回完整的父级链路（包含自身）
func (r *taskRepo) GetTaskParentChain(ctx context.Context, taskID, userID string) ([]*biz.Task, error) {
	var parentChain []*biz.Task
	currentTaskID := taskID

	// 循环向上查找父级任务，最多遍历5层（防止死循环）
	for i := 0; i < 5 && currentTaskID != ""; i++ {
		var dataTask Task
		err := r.db.WithContext(ctx).
			Where("id = ? AND user_id = ?", currentTaskID, userID).
			First(&dataTask).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				break // 任务不存在，结束查找
			}
			return nil, err
		}

		// 转换为业务模型并添加到链路前端
		bizTask := r.converter.DataToBiz(&dataTask)
		parentChain = append([]*biz.Task{bizTask}, parentChain...)

		// 设置下一个要查找的父级任务ID
		currentTaskID = dataTask.ParentID
	}

	return parentChain, nil
}

// UpdateTreeOptimizationFields 更新任务的树优化字段
// 用于维护树结构的冗余字段，确保查询性能
func (r *taskRepo) UpdateTreeOptimizationFields(ctx context.Context, taskID, userID string) error {
	// 获取任务详情
	var task Task
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", taskID, userID).
		First(&task).Error
	if err != nil {
		return err
	}

	// 计算树深度和根任务ID
	treeDepth := 0
	rootTaskID := taskID
	currentParentID := task.ParentID

	// 向上遍历找到根任务并计算深度
	for currentParentID != "" && treeDepth < 5 {
		var parentTask Task
		err := r.db.WithContext(ctx).
			Select("id, parent_id").
			Where("id = ? AND user_id = ?", currentParentID, userID).
			First(&parentTask).Error
		if err != nil {
			// 父任务查询失败，返回错误避免写入脏数据
			return fmt.Errorf("failed to query parent task %s: %w", currentParentID, err)
		}

		treeDepth++
		if parentTask.ParentID == "" {
			rootTaskID = parentTask.ID
			break
		}
		currentParentID = parentTask.ParentID
	}

	// 计算子任务数量
	var childrenCount int64
	err = r.db.WithContext(ctx).Model(&Task{}).
		Where("parent_id = ? AND user_id = ?", taskID, userID).
		Count(&childrenCount).Error
	if err != nil {
		return err
	}

	// 更新优化字段
	updates := map[string]interface{}{
		"tree_depth":     treeDepth,
		"root_task_id":   rootTaskID,
		"children_count": childrenCount,
		"has_children":   childrenCount > 0,
	}

	return r.db.WithContext(ctx).Model(&Task{}).
		Where("id = ? AND user_id = ?", taskID, userID).
		Updates(updates).Error
}

// JournalRepo 日志仓库实现
type journalRepo struct {
	db        *gorm.DB
	converter *JournalConverter
}

func NewJournalRepo(db *gorm.DB) biz.JournalRepo {
	return &journalRepo{
		db:        db,
		converter: NewJournalConverter(),
	}
}

func (r *journalRepo) CreateJournal(ctx context.Context, bizJournal *biz.Journal) error {
	dataJournal := r.converter.BizToData(bizJournal)
	return r.db.WithContext(ctx).Create(dataJournal).Error
}

func (r *journalRepo) GetJournalWithAuth(ctx context.Context, journalID, userID string) (*biz.Journal, error) {
    var dataJournal Journal
    err := r.db.WithContext(ctx).
        Where("id = ? AND user_id = ?", journalID, userID).
        First(&dataJournal).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, model.ErrRecordNotFound
        }
        return nil, err
    }

    return r.converter.DataToBiz(&dataJournal), nil
}

func (r *journalRepo) UpdateJournal(ctx context.Context, bizJournal *biz.Journal) error {
	dataJournal := r.converter.BizToData(bizJournal)
	return r.db.WithContext(ctx).Save(dataJournal).Error
}

func (r *journalRepo) DeleteJournalWithAuth(ctx context.Context, journalID, userID string) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", journalID, userID).
		Delete(&Journal{}).Error
}

func (r *journalRepo) ListJournals(ctx context.Context, userID string, periodStart, periodEnd time.Time, journalType int) ([]*biz.Journal, error) {
	var dataJournals []*Journal
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND journal_type = ? AND period_start >= ? AND period_end <= ?",
			userID, journalType, periodStart, periodEnd).
		Find(&dataJournals).Error

	if err != nil {
		return nil, err
	}

	return r.converter.DataToBizList(dataJournals), nil
}

func (r *journalRepo) ListAllJournals(ctx context.Context, userID string, offset, limit int) ([]*biz.Journal, error) {
	var dataJournals []*Journal
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&dataJournals).Error

	if err != nil {
		return nil, err
	}

	return r.converter.DataToBizList(dataJournals), nil
}

// ListJournalsWithPagination 分页查询日志并返回总数
// 支持按日志类型过滤和时间范围过滤
func (r *journalRepo) ListJournalsWithPagination(ctx context.Context, userID string, page, pageSize int, journalType *int, periodStart, periodEnd *time.Time) ([]*biz.Journal, int64, error) {
	// 构建基础查询
	query := r.db.WithContext(ctx).Model(&Journal{}).Where("user_id = ?", userID)

	// 日志类型过滤
	if journalType != nil {
		query = query.Where("journal_type = ?", *journalType)
	}

	// 时间范围过滤
	if periodStart != nil {
		query = query.Where("period_start >= ?", *periodStart)
	}
	if periodEnd != nil {
		query = query.Where("period_end <= ?", *periodEnd)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var dataJournals []*Journal
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&dataJournals).Error

	if err != nil {
		return nil, 0, err
	}

	return r.converter.DataToBizList(dataJournals), total, nil
}

// UserRepo 用户仓库实现
type userRepo struct {
	db        *gorm.DB
	converter *UserConverter
}

func NewUserRepo(db *gorm.DB) biz.UserRepo {
	return &userRepo{
		db:        db,
		converter: NewUserConverter(),
	}
}

func (r *userRepo) CreateUser(ctx context.Context, bizUser *biz.User) error {
	dataUser := r.converter.BizToData(bizUser)
	return r.db.WithContext(ctx).Create(dataUser).Error
}

func (r *userRepo) GetUserByID(ctx context.Context, userID string) (*biz.User, error) {
    var dataUser User
    err := r.db.WithContext(ctx).
        Where("id = ?", userID).
        First(&dataUser).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, biz.ErrUserNotFound
        }
        return nil, err
    }

    return r.converter.DataToBiz(&dataUser), nil
}

func (r *userRepo) GetUserByUserName(ctx context.Context, username string) (*biz.User, error) {
    var dataUser User
    err := r.db.WithContext(ctx).
        Where("user_name = ?", username).
        First(&dataUser).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, biz.ErrUserNotFound
        }
        return nil, err
    }

    return r.converter.DataToBiz(&dataUser), nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*biz.User, error) {
    var dataUser User
    err := r.db.WithContext(ctx).
        Where("email = ?", email).
        First(&dataUser).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, biz.ErrUserNotFound
        }
        return nil, err
    }

    return r.converter.DataToBiz(&dataUser), nil
}

func (r *userRepo) UpdateUser(ctx context.Context, bizUser *biz.User) error {
	dataUser := r.converter.BizToData(bizUser)
	return r.db.WithContext(ctx).Save(dataUser).Error
}

func (r *userRepo) DeleteUser(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Where("id = ?", userID).
		Delete(&User{}).Error
}
