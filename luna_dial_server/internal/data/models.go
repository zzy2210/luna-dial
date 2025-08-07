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
	
	// 新增：树结构优化字段
	// 设计思路：通过冗余字段减少递归查询，提升性能
	// 查询策略：利用 root_task_id 批量查询整个树，然后在内存中构建父子关系
	
	HasChildren   bool   `gorm:"default:false" json:"has_children"`     // 是否有子任务：快速判断节点类型，避免额外查询
	ChildrenCount int    `gorm:"default:0" json:"children_count"`       // 直接子任务数量：用于统计和分页计算
	RootTaskID    string `gorm:"type:varchar(36);index" json:"root_task_id"` // 根任务ID：批量查询整个树的关键字段
	TreeDepth     int    `gorm:"default:0" json:"tree_depth"`           // 树深度：排序和层级控制，根任务depth=0
	
	// 查询示例：
	// 1. 获取指定任务的完整任务树（任意层级的taskID）：
	//    步骤1：SELECT root_task_id FROM tasks WHERE id = ? AND user_id = ?
	//    步骤2：SELECT * FROM tasks WHERE user_id = ? AND (id = root_task_id OR root_task_id = root_task_id) ORDER BY tree_depth
	// 2. 获取以指定任务为根的子树：SELECT * FROM tasks WHERE user_id = ? AND root_task_id = ? ORDER BY tree_depth
	// 3. 获取根任务分页：SELECT * FROM tasks WHERE user_id = ? AND (parent_id IS NULL OR parent_id = '') LIMIT ? OFFSET ?
	// 4. 批量获取子任务：SELECT * FROM tasks WHERE user_id = ? AND root_task_id IN (?, ?, ?) ORDER BY root_task_id, tree_depth
	
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
