package data

import (
	"context"

	"gorm.io/gorm"

	"luna_dial/internal/biz"
)

// personalBackupRepo 个人数据备份的 GORM 仓库实现
type personalBackupRepo struct {
	// db 数据库连接
	db *gorm.DB
	// taskConverter 任务模型转换器
	taskConverter *TaskConverter
	// journalConverter 日志模型转换器
	journalConverter *JournalConverter
}

// NewPersonalBackupRepo 创建个人数据备份仓库
func NewPersonalBackupRepo(db *gorm.DB) biz.PersonalBackupRepo {
	return &personalBackupRepo{
		db:              db,
		taskConverter:   NewTaskConverter(),
		journalConverter: NewJournalConverter(),
	}
}

func (r *personalBackupRepo) ListAllTasksForUser(ctx context.Context, userID string) ([]*biz.Task, error) {
	var rows []Task
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]*biz.Task, 0, len(rows))
	for i := range rows {
		out = append(out, r.taskConverter.DataToBiz(&rows[i]))
	}
	return out, nil
}

func (r *personalBackupRepo) ListAllJournalsForUser(ctx context.Context, userID string) ([]*biz.Journal, error) {
	var rows []Journal
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]*biz.Journal, 0, len(rows))
	for i := range rows {
		out = append(out, r.journalConverter.DataToBiz(&rows[i]))
	}
	return out, nil
}

func (r *personalBackupRepo) ImportOverwrite(ctx context.Context, userID string, tasks []*biz.Task, journals []*biz.Journal) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&Task{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", userID).Delete(&Journal{}).Error; err != nil {
			return err
		}

		if len(tasks) > 0 {
			dataTasks := make([]Task, 0, len(tasks))
			for i := range tasks {
				if tasks[i] == nil {
					continue
				}
				tt := *tasks[i]
				tt.UserID = userID
				dt := r.taskConverter.BizToData(&tt)
				if dt == nil {
					continue
				}
				dataTasks = append(dataTasks, *dt)
			}

			if len(dataTasks) > 0 {
				if err := tx.CreateInBatches(dataTasks, 200).Error; err != nil {
					return err
				}
			}
		}

		if len(journals) > 0 {
			dataJournals := make([]Journal, 0, len(journals))
			for i := range journals {
				if journals[i] == nil {
					continue
				}
				jj := *journals[i]
				jj.UserID = userID
				dj := r.journalConverter.BizToData(&jj)
				if dj == nil {
					continue
				}
				dataJournals = append(dataJournals, *dj)
			}

			if len(dataJournals) > 0 {
				if err := tx.CreateInBatches(dataJournals, 200).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

