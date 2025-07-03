package main

import (
	"net/http"

	"okr-web/internal/controller"
	"okr-web/internal/middleware"
	"okr-web/internal/service"

	"github.com/labstack/echo/v4"
)

// setupRoutes 设置所有API路由
func setupRoutes(e *echo.Echo, services *service.Services, jwtSecret string) {
	// 创建控制器实例
	userController := controller.NewUserController(services.User)
	taskController := controller.NewTaskController(services.Task)
	journalController := controller.NewJournalController(services.Journal)
	statsController := controller.NewStatsController(services.Stats)
	planController := controller.NewPlanController(services.Task)

	// JWT中间件
	jwtMiddleware := middleware.JWT(jwtSecret, services.User)

	// API路由组
	api := e.Group("/api")

	// 公开路由（无需认证）
	setupPublicRoutes(api, userController)

	// 受保护路由（需要认证）
	setupProtectedRoutes(api, jwtMiddleware, userController, taskController, journalController, statsController, planController)
}

// setupPublicRoutes 设置公开路由
func setupPublicRoutes(api *echo.Group, userController *controller.UserController) {
	// 测试路由
	api.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
	})

	// 用户认证相关路由
	auth := api.Group("/auth")
	auth.POST("/register", userController.Register)
	auth.POST("/login", userController.Login)
}

// setupProtectedRoutes 设置受保护路由
func setupProtectedRoutes(
	api *echo.Group,
	jwtMiddleware echo.MiddlewareFunc,
	userController *controller.UserController,
	taskController *controller.TaskController,
	journalController *controller.JournalController,
	statsController *controller.StatsController,
	planController *controller.PlanController,
) {
	// 应用JWT和Auth中间件，顺序保证先JWT后Auth
	protected := api.Group("", jwtMiddleware, middleware.AuthMiddleware())

	// 用户相关路由
	setupUserRoutes(protected, userController)

	// 任务相关路由
	setupTaskRoutes(protected, taskController)

	// 日志相关路由
	setupJournalRoutes(protected, journalController)

	// 统计相关路由
	setupStatsRoutes(protected, statsController)

	// 计划视图相关路由
	setupPlanRoutes(protected, planController)
}

// setupUserRoutes 设置用户相关路由
func setupUserRoutes(protected *echo.Group, userController *controller.UserController) {
	users := protected.Group("/users")
	users.GET("/me", userController.GetUserInfo)
	users.PUT("/me", userController.UpdateUser)
	users.POST("/logout", userController.Logout)
}

// setupTaskRoutes 设置任务相关路由
func setupTaskRoutes(protected *echo.Group, taskController *controller.TaskController) {
	tasks := protected.Group("/tasks")

	// 基本CRUD操作
	tasks.POST("", taskController.CreateTask)
	tasks.GET("/:id", taskController.GetTask)
	tasks.PUT("/:id", taskController.UpdateTask)
	tasks.DELETE("/:id", taskController.DeleteTask)

	// 任务列表和查询
	tasks.GET("", taskController.GetTasks)
	tasks.GET("/global", taskController.GetGlobalView)
	tasks.GET("/:id/context", taskController.GetContextView)
	tasks.GET("/:id/children", taskController.GetTaskChildren)
	tasks.GET("/:id/tree", taskController.GetTaskTree)

	// 子任务操作
	tasks.POST("/:id/subtasks", taskController.CreateSubTask)

	// 任务评分
	tasks.PUT("/:id/score", taskController.UpdateTaskScore)
}

// setupJournalRoutes 设置日志相关路由
func setupJournalRoutes(protected *echo.Group, journalController *controller.JournalController) {
	journals := protected.Group("/journals")

	// 基本CRUD操作
	journals.POST("", journalController.CreateJournal)
	journals.GET("/:id", journalController.GetJournal)
	journals.PUT("/:id", journalController.UpdateJournal)
	journals.DELETE("/:id", journalController.DeleteJournal)

	// 日志查询
	journals.GET("", journalController.GetJournalsByUser)
	journals.GET("/time-range", journalController.GetJournalsByTime)

	// 日志与任务关联
	journals.POST("/:id/tasks", journalController.LinkJournalToTasks)
}

// setupStatsRoutes 设置统计相关路由
func setupStatsRoutes(protected *echo.Group, statsController *controller.StatsController) {
	stats := protected.Group("/stats")

	// 各种统计查询
	stats.GET("/overview", statsController.GetUserOverview)
	stats.GET("/completion", statsController.GetTaskCompletionStats)
	stats.GET("/score-trend", statsController.GetScoreTrend)
	stats.GET("/score-trend-ref", statsController.GetScoreTrendByReference)
	stats.GET("/time-distribution", statsController.GetTimeDistribution)
}

// setupPlanRoutes 设置计划视图相关路由
func setupPlanRoutes(protected *echo.Group, planController *controller.PlanController) {
	// 计划视图路由
	protected.GET("/plan", planController.GetPlanView)
}
