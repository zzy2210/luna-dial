package biz

import (
	"context"
	"time"
)

type JournalRepo interface {
	CreateJournal(ctx context.Context, journal *Journal) error
	UpdateJournal(ctx context.Context, journal *Journal) error
	DeleteJournal(ctx context.Context, journalID, userID string) error
	GetJournal(ctx context.Context, journalID, userID string) (*Journal, error)
	ListJournals(ctx context.Context, userID string, periodStart, periodEnd time.Time, journalType string) ([]*Journal, error)
}
