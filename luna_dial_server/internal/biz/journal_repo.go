package biz

import (
	"context"
	"time"
)

type JournalRepo interface {
	CreateJournal(ctx context.Context, journal *Journal) error
	UpdateJournal(ctx context.Context, journal *Journal) error
	DeleteJournalWithAuth(ctx context.Context, journalID, userID string) error
	GetJournalWithAuth(ctx context.Context, journalID, userID string) (*Journal, error)
	ListJournals(ctx context.Context, userID string, periodStart, periodEnd time.Time, journalType int) ([]*Journal, error)
	ListAllJournals(ctx context.Context, userID string, offset, limit int) ([]*Journal, error)
	ListJournalsWithPagination(ctx context.Context, userID string, page, pageSize int, journalType *int, periodStart, periodEnd *time.Time) ([]*Journal, int64, error)
}
