package biz

import (
	"context"
	"time"
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
	return nil, ErrNoPermission // TODO: 实现
}

// 编辑日志
func (uc *JournalUsecase) UpdateJournal(ctx context.Context, param UpdateJournalParam) (*Journal, error) {
	return nil, ErrNoPermission // TODO: 实现
}

// 删除日志
func (uc *JournalUsecase) DeleteJournal(ctx context.Context, param DeleteJournalParam) error {
	return ErrNoPermission // TODO: 实现
}

// 获取日志详情
func (uc *JournalUsecase) GetJournal(ctx context.Context, param GetJournalParam) (*Journal, error) {
	return nil, ErrNoPermission // TODO: 实现
}

// 获取指定时间的指定类型的日志列表
func (uc *JournalUsecase) ListJournalByPeriod(ctx context.Context, param ListJournalByPeriodParam) ([]Journal, error) {
	return nil, ErrNoPermission // TODO: 实现
}

// 获取全部日志列表，带分页
func (uc *JournalUsecase) ListAllJournals(ctx context.Context, param ListAllJournalsParam) ([]Journal, error) {
	return nil, ErrNoPermission // TODO: 实现
}
