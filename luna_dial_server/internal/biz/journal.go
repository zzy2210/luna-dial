package biz

import (
	"context"
	"errors"
	"time"

	"luna_dial/internal/model"
)

type Journal struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	JournalType PeriodType `json:"journal_type"`
	TimePeriod  Period     `json:"time_period"`
	Icon        string     `json:"icon"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UserID      string     `json:"user_id"`
}

// 创建日志参数
type CreateJournalParam struct {
	UserID      string
	Title       string
	Content     string
	JournalType PeriodType
	TimePeriod  Period
	Icon        string
}

// 编辑日志参数
type UpdateJournalParam struct {
	JournalID   string
	UserID      string
	Title       *string
	Content     *string
	JournalType *PeriodType
	TimePeriod  *Period
	Icon        *string
}

// 删除日志参数
type DeleteJournalParam struct {
	JournalID string
	UserID    string
}

// 获取日志详情参数
type GetJournalParam struct {
	JournalID string
	UserID    string
}

// 分页参数
type PaginationParam struct {
	PageNum  int
	PageSize int
}

// 获取全部日志列表参数
type ListAllJournalsParam struct {
	UserID     string
	Pagination PaginationParam
}

type JournalUsecase struct {
	repo JournalRepo
	// log  *log.Helper
}

// 获取指定时间的指定类型的日志列表参数
type ListJournalByPeriodParam struct {
	UserID  string
	Period  Period
	GroupBy PeriodType
}

func NewJournalUsecase(repo JournalRepo) *JournalUsecase {
	return &JournalUsecase{repo: repo}
}

// 创建日志
func (uc *JournalUsecase) CreateJournal(ctx context.Context, param CreateJournalParam) (*Journal, error) {
	if param.UserID == "" {
		return nil, ErrNoPermission
	}
	if param.Title == "" {
		return nil, ErrInvalidInput
	}
	if param.Content == "" {
		return nil, ErrJournalContentEmpty
	}

	if !param.TimePeriod.IsValid() {
		return nil, ErrJournalPeriodInvalid
	}

	if !param.TimePeriod.MatchesPeriodType(param.JournalType) {
		return nil, ErrJournalTypeInvalid
	}
	journal := &Journal{
		ID:          generateID(), // 生成ID逻辑待实现
		Title:       param.Title,
		Content:     param.Content,
		JournalType: param.JournalType,
		TimePeriod:  param.TimePeriod,
		Icon:        param.Icon,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UserID:      param.UserID,
	}
	if err := uc.repo.CreateJournal(ctx, journal); err != nil {
		return nil, err
	}

	return journal, nil
}

// 编辑日志
// 一样，暂时不支持更新时间类型
func (uc *JournalUsecase) UpdateJournal(ctx context.Context, param UpdateJournalParam) (*Journal, error) {
	if param.JournalID == "" || param.UserID == "" {
		return nil, ErrInvalidInput
	}
	if param.Title != nil && *param.Title == "" {
		return nil, ErrInvalidInput
	}
	if param.Content != nil && *param.Content == "" {
		return nil, ErrJournalContentEmpty
	}
	if param.TimePeriod != nil && !param.TimePeriod.IsValid() {
		return nil, ErrJournalPeriodInvalid
	}

	oldJournal, err := uc.repo.GetJournalWithAuth(ctx, param.JournalID, param.UserID)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			return nil, ErrJournalNotFound
		}
		return nil, err
	}
	if oldJournal == nil {
		return nil, ErrJournalNotFound
	}
	// 更新数据
	if param.Content != nil {
		oldJournal.Content = *param.Content
	}
	if param.Title != nil {
		oldJournal.Title = *param.Title
	}
	if param.TimePeriod != nil {
		oldJournal.TimePeriod = *param.TimePeriod
	}
	if param.Icon != nil {
		oldJournal.Icon = *param.Icon
	}

	if err := uc.repo.UpdateJournal(ctx, oldJournal); err != nil {
		return nil, err
	}
	return oldJournal, nil

}

// 删除日志
func (uc *JournalUsecase) DeleteJournal(ctx context.Context, param DeleteJournalParam) error {
	if param.JournalID == "" || param.UserID == "" {
		return ErrInvalidInput
	}

	err := uc.repo.DeleteJournalWithAuth(ctx, param.JournalID, param.UserID)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			return ErrJournalNotFound
		}
		return err
	}

	return nil
}

// 获取日志详情
func (uc *JournalUsecase) GetJournal(ctx context.Context, param GetJournalParam) (*Journal, error) {
	if param.JournalID == "" || param.UserID == "" {
		return nil, ErrInvalidInput
	}

	journal, err := uc.repo.GetJournalWithAuth(ctx, param.JournalID, param.UserID)
	if err != nil {
		// 将数据库层错误转换为业务层错误
		if errors.Is(err, model.ErrRecordNotFound) {
			return nil, ErrJournalNotFound
		}
		return nil, err
	}
	if journal == nil {
		return nil, ErrJournalNotFound
	}

	return journal, nil
}

// 获取指定时间的指定类型的日志列表
func (uc *JournalUsecase) ListJournalByPeriod(ctx context.Context, param ListJournalByPeriodParam) ([]*Journal, error) {
	if param.UserID == "" {
		return nil, ErrUserIDEmpty
	}
	if !param.Period.IsValid() {
		return nil, ErrJournalPeriodInvalid
	}
	// 如果时间范围与类型不匹配
	if !param.Period.MatchesPeriodType(param.GroupBy) {
		return nil, ErrJournalTypeInvalid
	}

	groupBy := int(param.GroupBy)
	journals, err := uc.repo.ListJournals(ctx, param.UserID, param.Period.Start, param.Period.End, groupBy)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			return nil, ErrJournalNotFound
		}
		return nil, err
	}

	return journals, nil
}

// 获取全部日志列表，带分页
func (uc *JournalUsecase) ListAllJournals(ctx context.Context, param ListAllJournalsParam) ([]*Journal, error) {
	if param.UserID == "" {
		return nil, ErrUserIDEmpty
	}
	if param.Pagination.PageNum <= 0 || param.Pagination.PageSize <= 0 {
		return nil, ErrInvalidInput
	}

	// 计算offset
	offset := (param.Pagination.PageNum - 1) * param.Pagination.PageSize
	limit := param.Pagination.PageSize

	journals, err := uc.repo.ListAllJournals(ctx, param.UserID, offset, limit)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			// 没有记录时返回空列表，而不是错误
			return []*Journal{}, nil
		}
		return nil, err
	}

	return journals, nil
}
