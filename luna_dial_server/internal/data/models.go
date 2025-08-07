package data

import "time"

// 采用gorm吧
// 用户数据模型
type User struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserName  string    `gorm:"uniqueIndex;type:varchar(50);not null" json:"user_name"`
	Name      string    `gorm:"type:varchar(100)" json:"name"`
	Email     string    `gorm:"uniqueIndex;type:varchar(100)" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// 任务数据模型
type Task struct {
	ID          string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID      string    `gorm:"type:varchar(36);index;not null" json:"user_id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	TaskType    int       `gorm:"type:int;not null" json:"task_type"`
	PeriodStart time.Time `gorm:"type:datetime" json:"period_start"`
	PeriodEnd   time.Time `gorm:"type:datetime" json:"period_end"`
	Tags        string    `gorm:"type:text" json:"tags"`
	Icon        string    `gorm:"type:varchar(10)" json:"icon"`
	Score       int       `gorm:"default:0" json:"score"`
	Status      int       `gorm:"default:0;not null" json:"status"`
	Priority    int       `gorm:"default:0;not null" json:"priority"`
	ParentID    string    `gorm:"type:varchar(36);index" json:"parent_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// 日志数据模型
type Journal struct {
	ID          string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID      string    `gorm:"type:varchar(36);index;not null" json:"user_id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	JournalType int       `gorm:"type:int;not null" json:"journal_type"`
	PeriodStart time.Time `gorm:"type:datetime" json:"period_start"`
	PeriodEnd   time.Time `gorm:"type:datetime" json:"period_end"`
	Icon        string    `gorm:"type:varchar(10)" json:"icon"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
