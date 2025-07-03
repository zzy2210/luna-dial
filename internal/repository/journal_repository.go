package repository

import (
	"context"
	"fmt"
	"time"

	"okr-web/ent"
	"okr-web/ent/journalentry"

	"github.com/google/uuid"
)

// journalRepository 日志Repository实现
type journalRepository struct {
	client *ent.Client
}

// NewJournalRepository 创建新的日志Repository
func NewJournalRepository(client *ent.Client) JournalRepository {
	return &journalRepository{client: client}
}

// Create 创建新日志条目
func (r *journalRepository) Create(ctx context.Context, builder func(*ent.JournalEntryCreate) *ent.JournalEntryCreate) (*ent.JournalEntry, error) {
	return builder(r.client.JournalEntry.Create()).Save(ctx)
}

// GetByID 根据ID获取日志条目
func (r *journalRepository) GetByID(ctx context.Context, id uuid.UUID) (*ent.JournalEntry, error) {
	j, err := r.client.JournalEntry.
		Query().
		Where(journalentry.ID(id)).
		WithUser().
		WithTasks().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("journal entry not found")
		}
		return nil, fmt.Errorf("failed to get journal entry: %w", err)
	}
	return j, nil
}

// GetByUserID 根据用户ID获取日志条目列表
func (r *journalRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*ent.JournalEntry, error) {
	query := r.client.JournalEntry.
		Query().
		Where(journalentry.UserID(userID)).
		WithUser().
		WithTasks().
		Order(ent.Desc(journalentry.FieldCreatedAt))

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	journals, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get journal entries by user ID: %w", err)
	}
	return journals, nil
}

// GetByTimeReference 根据时间引用获取日志条目
func (r *journalRepository) GetByTimeReference(ctx context.Context, userID uuid.UUID, timeRef string) ([]*ent.JournalEntry, error) {
	journals, err := r.client.JournalEntry.
		Query().
		Where(
			journalentry.UserID(userID),
			journalentry.TimeReference(timeRef),
		).
		WithUser().
		WithTasks().
		Order(ent.Desc(journalentry.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get journal entries by time reference: %w", err)
	}
	return journals, nil
}

// Update 更新日志条目
func (r *journalRepository) Update(ctx context.Context, id uuid.UUID, updater func(*ent.JournalEntryUpdateOne) *ent.JournalEntryUpdateOne) (*ent.JournalEntry, error) {
	updateOne := r.client.JournalEntry.UpdateOneID(id)
	updateOne = updater(updateOne)

	j, err := updateOne.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("journal entry not found")
		}
		return nil, fmt.Errorf("failed to update journal entry: %w", err)
	}
	return j, nil
}

// Delete 删除日志条目
func (r *journalRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.client.JournalEntry.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("journal entry not found")
		}
		return fmt.Errorf("failed to delete journal entry: %w", err)
	}
	return nil
}

// CountByUserID 统计用户的日志条目数量
func (r *journalRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	count, err := r.client.JournalEntry.
		Query().
		Where(journalentry.UserID(userID)).
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count journal entries: %w", err)
	}
	return count, nil
}

// GetByTimeRange 根据时间范围获取日志条目
func (r *journalRepository) GetByTimeRange(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*ent.JournalEntry, error) {
	journals, err := r.client.JournalEntry.
		Query().
		Where(
			journalentry.UserID(userID),
			journalentry.CreatedAtGTE(start),
			journalentry.CreatedAtLTE(end),
		).
		WithUser().
		WithTasks().
		Order(ent.Desc(journalentry.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get journal entries by time range: %w", err)
	}
	return journals, nil
}

// GetByTimeScaleAndReferences 根据时间尺度和 time_reference 列表获取日志条目
func (r *journalRepository) GetByTimeScaleAndReferences(ctx context.Context, userID uuid.UUID, timeScale journalentry.TimeScale, timeRefs []string) ([]*ent.JournalEntry, error) {
	journals, err := r.client.JournalEntry.
		Query().
		Where(
			journalentry.UserID(userID),
			journalentry.TimeScaleEQ(timeScale),
			journalentry.TimeReferenceIn(timeRefs...),
		).
		WithUser().
		WithTasks().
		Order(ent.Desc(journalentry.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get journal entries by time scale and references: %w", err)
	}
	return journals, nil
}
