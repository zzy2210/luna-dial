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
}

func NewService(ctx context.Context, e *echo.Echo, dataInstance *data.Data) *Service {
	// 创建各个 repo
	taskRepo := data.NewTaskRepo(dataInstance.DB)
	journalRepo := data.NewJournalRepo(dataInstance.DB)
	userRepo := data.NewUserRepo(dataInstance.DB)

	return &Service{
		e:              e,
		systemConfig:   dataInstance.SystemConfig,
		sessionManager: dataInstance.SessionManager,
		journalUsecase: biz.NewJournalUsecase(journalRepo),
		userUsecase:    biz.NewUserUsecase(userRepo),
		taskUsecase:    biz.NewTaskUsecase(taskRepo),
	}
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
	public.GET("/system-info", func(c echo.Context) error {
		jwtSecret, err := s.systemConfig.GetJWTSecret(c.Request().Context())
		if err != nil {
			return c.JSON(500, map[string]string{"error": "获取JWT密钥失败"})
		}
		return c.JSON(200, map[string]interface{}{
			"jwt_secret_length": len(jwtSecret),
			"admin_account":     "admin",
			"admin_password":    "admin@123",
		})
	})
	public.POST("/auth/login", s.handleSessionLogin)

}
