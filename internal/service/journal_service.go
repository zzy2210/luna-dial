package service

import (
	"context"
	"math"

	"okr-web/ent"
	journalent "okr-web/ent/journalentry"
	"okr-web/internal/repository"
	"okr-web/internal/types"

	"github.com/google/uuid"
)

// JournalServiceImpl 日志服务实现
type JournalServiceImpl struct {
	journalRepo repository.JournalRepository
	taskRepo    repository.TaskRepository
	userRepo    repository.UserRepository
}

// NewJournalService 创建日志服务
func NewJournalService(journalRepo repository.JournalRepository, taskRepo repository.TaskRepository, userRepo repository.UserRepository) JournalService {
	return &JournalServiceImpl{
		journalRepo: journalRepo,
		taskRepo:    taskRepo,
		userRepo:    userRepo,
	}
}

// CreateJournal 创建日志
func (s *JournalServiceImpl) CreateJournal(ctx context.Context, userID uuid.UUID, req JournalRequest) (*ent.JournalEntry, error) {
	// 验证用户是否存在
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, &types.AppError{
			Code:    404,
			Message: "用户不存在",
			Type:    "USER_NOT_FOUND",
		}
	}

	// 验证时间尺度
	if !req.TimeScale.IsValid() {
		return nil, &types.AppError{
			Code:    400,
			Message: "无效的时间尺度",
			Type:    "INVALID_TIME_SCALE",
		}
	}

	// 验证条目类型
	if !req.EntryType.IsValid() {
		return nil, &types.AppError{
			Code:    400,
			Message: "无效的条目类型",
			Type:    "INVALID_ENTRY_TYPE",
		}
	}

	// 如果有关联任务，验证任务是否存在且属于当前用户
	if len(req.TaskIDs) > 0 {
		for _, taskID := range req.TaskIDs {
			task, err := s.taskRepo.GetByID(ctx, taskID)
			if err != nil {
				return nil, &types.AppError{
					Code:    404,
					Message: "关联任务不存在",
					Type:    "TASK_NOT_FOUND",
				}
			}
			if task.UserID != userID {
				return nil, &types.AppError{
					Code:    403,
					Message: "无权限关联此任务",
					Type:    "FORBIDDEN",
				}
			}
		}
	}

	// 创建日志
	journal, err := s.journalRepo.Create(ctx, func(create *ent.JournalEntryCreate) *ent.JournalEntryCreate {
		create = create.
			SetContent(req.Content).
			SetTimeReference(req.TimeReference).
			SetTimeScale(journalent.TimeScale(req.TimeScale)).
			SetEntryType(journalent.EntryType(req.EntryType)).
			SetUserID(userID)

		return create
	})

	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "日志创建失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// 如果有关联任务，建立关联关系
	if len(req.TaskIDs) > 0 {
		err = s.LinkJournalToTasks(ctx, userID, journal.ID, req.TaskIDs)
		if err != nil {
			return nil, err
		}
	}

	return journal, nil
}

// GetJournal 获取日志
func (s *JournalServiceImpl) GetJournal(ctx context.Context, userID uuid.UUID, journalID uuid.UUID) (*ent.JournalEntry, error) {
	journal, err := s.journalRepo.GetByID(ctx, journalID)
	if err != nil {
		return nil, &types.AppError{
			Code:    404,
			Message: "日志不存在",
			Type:    "JOURNAL_NOT_FOUND",
		}
	}

	// 检查权限
	if journal.UserID != userID {
		return nil, &types.AppError{
			Code:    403,
			Message: "无权限访问此日志",
			Type:    "FORBIDDEN",
		}
	}

	return journal, nil
}

// UpdateJournal 更新日志
func (s *JournalServiceImpl) UpdateJournal(ctx context.Context, userID uuid.UUID, journalID uuid.UUID, req JournalRequest) (*ent.JournalEntry, error) {
	// 先获取日志并验证权限
	_, err := s.GetJournal(ctx, userID, journalID)
	if err != nil {
		return nil, err
	}

	// 验证时间尺度
	if !req.TimeScale.IsValid() {
		return nil, &types.AppError{
			Code:    400,
			Message: "无效的时间尺度",
			Type:    "INVALID_TIME_SCALE",
		}
	}

	// 验证条目类型
	if !req.EntryType.IsValid() {
		return nil, &types.AppError{
			Code:    400,
			Message: "无效的条目类型",
			Type:    "INVALID_ENTRY_TYPE",
		}
	}

	// 如果有关联任务，验证任务是否存在且属于当前用户
	if len(req.TaskIDs) > 0 {
		for _, taskID := range req.TaskIDs {
			task, err := s.taskRepo.GetByID(ctx, taskID)
			if err != nil {
				return nil, &types.AppError{
					Code:    404,
					Message: "关联任务不存在",
					Type:    "TASK_NOT_FOUND",
				}
			}
			if task.UserID != userID {
				return nil, &types.AppError{
					Code:    403,
					Message: "无权限关联此任务",
					Type:    "FORBIDDEN",
				}
			}
		}
	}

	// 更新日志
	updatedJournal, err := s.journalRepo.Update(ctx, journalID, func(update *ent.JournalEntryUpdateOne) *ent.JournalEntryUpdateOne {
		update = update.
			SetContent(req.Content).
			SetTimeReference(req.TimeReference).
			SetTimeScale(journalent.TimeScale(req.TimeScale)).
			SetEntryType(journalent.EntryType(req.EntryType))

		return update
	})

	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "日志更新失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// 如果有关联任务，重新建立关联关系
	if len(req.TaskIDs) > 0 {
		err = s.LinkJournalToTasks(ctx, userID, journalID, req.TaskIDs)
		if err != nil {
			return nil, err
		}
	}

	return updatedJournal, nil
}

// DeleteJournal 删除日志
func (s *JournalServiceImpl) DeleteJournal(ctx context.Context, userID uuid.UUID, journalID uuid.UUID) error {
	// 先获取日志并验证权限
	_, err := s.GetJournal(ctx, userID, journalID)
	if err != nil {
		return err
	}

	// 删除日志
	err = s.journalRepo.Delete(ctx, journalID)
	if err != nil {
		return &types.AppError{
			Code:    500,
			Message: "日志删除失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	return nil
}

// GetJournalsByTime 根据时间范围获取日志
func (s *JournalServiceImpl) GetJournalsByTime(ctx context.Context, userID uuid.UUID, timeRange TimeRangeRequest) ([]*ent.JournalEntry, error) {
	// 验证时间范围
	if timeRange.StartDate.After(timeRange.EndDate) {
		return nil, &types.AppError{
			Code:    400,
			Message: "开始时间不能晚于结束时间",
			Type:    "INVALID_TIME_RANGE",
		}
	}

	// 这里需要Repository层支持时间范围查询
	// 暂时使用简单的用户ID查询
	journals, err := s.journalRepo.GetByUserID(ctx, userID, 100, 0)
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "获取日志失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// 在内存中过滤时间范围（这里应该在数据库层面优化）
	var filteredJournals []*ent.JournalEntry
	for _, journal := range journals {
		if (journal.CreatedAt.After(timeRange.StartDate) || journal.CreatedAt.Equal(timeRange.StartDate)) &&
			(journal.CreatedAt.Before(timeRange.EndDate) || journal.CreatedAt.Equal(timeRange.EndDate)) {
			filteredJournals = append(filteredJournals, journal)
		}
	}

	return filteredJournals, nil
}

// GetJournalsByUser 获取用户的日志列表（分页）
func (s *JournalServiceImpl) GetJournalsByUser(ctx context.Context, userID uuid.UUID, filters JournalFilters) (*JournalListResponse, error) {
	// 验证分页参数
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 || filters.PageSize > 100 {
		filters.PageSize = 20
	}

	// 计算偏移量
	offset := (filters.Page - 1) * filters.PageSize

	// 获取日志列表
	journals, err := s.journalRepo.GetByUserID(ctx, userID, filters.PageSize, offset)
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "获取日志列表失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// 获取总数
	total, err := s.journalRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "获取日志总数失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(filters.PageSize)))

	return &JournalListResponse{
		Journals:    journals,
		Total:       int64(total),
		CurrentPage: filters.Page,
		PageSize:    filters.PageSize,
		TotalPages:  totalPages,
	}, nil
}

// LinkJournalToTasks 关联日志到任务
func (s *JournalServiceImpl) LinkJournalToTasks(ctx context.Context, userID uuid.UUID, journalID uuid.UUID, taskIDs []uuid.UUID) error {
	// 这里需要通过Repository层实现多对多关联
	// 由于Ent的复杂性，这里先简化处理
	// 实际实现时需要使用Ent的边(Edge)功能

	// 验证日志权限
	_, err := s.GetJournal(ctx, userID, journalID)
	if err != nil {
		return err
	}

	// 验证所有任务权限
	for _, taskID := range taskIDs {
		task, err := s.taskRepo.GetByID(ctx, taskID)
		if err != nil {
			return &types.AppError{
				Code:    404,
				Message: "任务不存在",
				Type:    "TASK_NOT_FOUND",
			}
		}
		if task.UserID != userID {
			return &types.AppError{
				Code:    403,
				Message: "无权限关联此任务",
				Type:    "FORBIDDEN",
			}
		}
	}

	// TODO: 实现实际的多对多关联逻辑
	// 这里需要使用Ent的Update().AddTasks()方法
	// 但需要先获取Journal的完整实体

	return nil
}
