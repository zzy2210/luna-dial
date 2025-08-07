package service

import (
	"context"
	"luna_dial/internal/biz"
	"luna_dial/internal/data"

	"github.com/labstack/echo/v4"
)

type Service struct {
	e *echo.Echo

	systemConfig   *data.SystemConfig
	sessionManager data.SessionManager
	journalUsecase *biz.JournalUsecase
	userUsecase    *biz.UserUsecase
	taskUsecase    *biz.TaskUsecase
	planUsecase    *biz.PlanUsecase
}

func NewService(ctx context.Context, e *echo.Echo, dataInstance *data.Data) *Service {
	// 创建各个 repo
	taskRepo := data.NewTaskRepo(dataInstance.DB)
	journalRepo := data.NewJournalRepo(dataInstance.DB)
	userRepo := data.NewUserRepo(dataInstance.DB)

	s := &Service{
		e:              e,
		systemConfig:   dataInstance.SystemConfig,
		sessionManager: dataInstance.SessionManager,
		journalUsecase: biz.NewJournalUsecase(journalRepo),
		userUsecase:    biz.NewUserUsecase(userRepo),
		taskUsecase:    biz.NewTaskUsecase(taskRepo),
	}
	s.planUsecase = biz.NewPlanUsecase(s.taskUsecase, s.journalUsecase)
	return s
}

func (s *Service) SetupRouter() {
	// 设置Session相关路由
	s.setupSessionRoutes()

	// 保留原有的公开路由
	s.setupPublicRoutes()
}

func (s *Service) setupPublicRoutes() {
	s.e.GET("/health", func(c echo.Context) error {
		return c.String(200, "Service is running")
	})
	s.e.GET("/version", func(c echo.Context) error {
		return c.String(200, "Version 1.0.0")
	})

	public := s.e.Group("/api/v1/public")
	public.POST("/auth/login", s.handleSessionLogin)

}

func (s *Service) setupSessionRoutes() {

	// 受保护的路由 - 需要Session认证
	protected := s.e.Group("/api/v1")
	protected.Use(s.SessionMiddleware())

	// 用户相关接口
	protected.GET("/auth/profile", s.handleGetProfile)
	protected.POST("/auth/logout", s.handleSessionLogout)
	protected.DELETE("/auth/logout-all", s.handleLogoutAllSessions)

	// 其他业务接口...
	userGroup := protected.Group("/users")
	userGroup.GET("/me", s.handleGetCurrentUser)

	journalGroup := protected.Group("/journals")
	journalGroup.GET("", s.handleListJournalsByPeriod)
	journalGroup.POST("", s.handleCreateJournal)
	journalGroup.PUT("/:journal_id", s.handleUpdateJournal)
	journalGroup.DELETE("/:journal_id", s.handleDeleteJournal)
	// 阶段五新增：分页查询日志
	journalGroup.GET("/paginated", s.handleListJournalsWithPagination)

	taskGroup := protected.Group("/tasks")
	taskGroup.GET("", s.handleListTasks)
	taskGroup.POST("", s.handleCreateTask)
	taskGroup.POST("/:task_id/subtasks", s.handleCreateSubTask)
	taskGroup.PUT("/:task_id", s.handleUpdateTask)
	taskGroup.DELETE("/:task_id", s.handleDeleteTask)
	taskGroup.POST("/:task_id/complete", s.handleCompleteTask)
	taskGroup.PUT("/:task_id/score", s.handleUpdateTaskScore)
	// 阶段五新增：任务树优化相关API
	taskGroup.GET("/roots", s.handleListRootTasks)                   // 分页查询根任务
	taskGroup.GET("/tree", s.handleListGlobalTaskTree)               // 全局任务树视图（分页）
	taskGroup.GET("/:task_id/tree", s.handleGetTaskTree)             // 获取指定任务的完整任务树
	taskGroup.GET("/:task_id/parents", s.handleGetTaskParents)       // 获取任务的父任务链
	taskGroup.PUT("/:task_id/move", s.handleMoveTask)                // 移动任务
	taskGroup.POST("/optimized", s.handleCreateTaskWithOptimization) // 使用优化的任务创建

	planGroup := protected.Group("/plans")
	planGroup.GET("", s.handleListPlans)
}
