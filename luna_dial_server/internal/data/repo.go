package data

import (
	"context"
	"luna_dial/internal/biz"
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

func (r *taskRepo) ListTaskTree(ctx context.Context, taskID, userID string) ([]*biz.Task, error) {
	// 递归查询子任务树
	var dataTasks []*Task

	// 先查询当前任务
	var currentTask Task
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", taskID, userID).
		First(&currentTask).Error
	if err != nil {
		return nil, err
	}

	// 然后递归查询所有子任务
	err = r.db.WithContext(ctx).
		Where("user_id = ? AND (id = ? OR parent_id = ?)", userID, taskID, taskID).
		Find(&dataTasks).Error

	if err != nil {
		return nil, err
	}

	return r.converter.DataToBizList(dataTasks), nil
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
