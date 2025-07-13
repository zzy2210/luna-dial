package biz

import "time"

type JournalRepo interface {
	CreateJournal(journal *Journal) error
	UpdateJournal(journal *Journal) error
	DeleteJournal(journalID, userID string) error
	GetJournal(journalID, userID string) (*Journal, error)
	ListJournals(userID string, periodStart, periodEnd time.Time, journalType string) ([]*Journal, error)
}
