package service

import (
	"okr-web/internal/repository"
)

// Services 服务集合
type Services struct {
	User    UserService
	Task    TaskService
	Journal JournalService
	Stats   StatsService
}

// NewServices 创建服务集合
func NewServices(repos *repository.Repositories, jwtSecret string, jwtExpiryHours int) *Services {
	return &Services{
		User:    NewUserService(repos.User, jwtSecret, jwtExpiryHours),
		Task:    NewTaskService(repos.Task, repos.User, repos.Journal),
		Journal: NewJournalService(repos.Journal, repos.Task, repos.User),
		Stats:   NewStatsService(repos.Task, repos.Journal, repos.User),
	}
}
