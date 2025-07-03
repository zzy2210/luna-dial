package repository

import (
	"okr-web/ent"
)

// Repositories Repository管理器，包含所有Repository
type Repositories struct {
	User    UserRepository
	Task    TaskRepository
	Journal JournalRepository
}

// NewRepositories 创建新的Repository管理器
func NewRepositories(client *ent.Client) *Repositories {
	return &Repositories{
		User:    NewUserRepository(client),
		Task:    NewTaskRepository(client),
		Journal: NewJournalRepository(client),
	}
}
