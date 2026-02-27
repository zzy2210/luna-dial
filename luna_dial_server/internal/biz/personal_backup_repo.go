package biz

import "context"

// PersonalBackupRepo 个人数据备份的数据访问接口
type PersonalBackupRepo interface {
	// ListAllTasksForUser 按用户列出所有任务（用于导出；顺序由实现保证稳定）
	ListAllTasksForUser(ctx context.Context, userID string) ([]*Task, error)
	// ListAllJournalsForUser 按用户列出所有日志（用于导出；顺序由实现保证稳定）
	ListAllJournalsForUser(ctx context.Context, userID string) ([]*Journal, error)
	// ImportOverwrite 覆盖式导入：在同一事务内先清空再写入
	ImportOverwrite(ctx context.Context, userID string, tasks []*Task, journals []*Journal) error
}

