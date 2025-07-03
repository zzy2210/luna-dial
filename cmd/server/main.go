package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/labstack/echo/v4"

	"okr-web/ent"
	"okr-web/internal/config"
	"okr-web/internal/middleware"
	"okr-web/internal/repository"
	"okr-web/internal/service"
	"okr-web/internal/utils"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库连接
	db, err := config.NewDatabase(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Println("Continuing without database connection...")
	}
	defer func() {
		if db != nil {
			db.Close()
		}
	}()

	// 运行数据库迁移
	if db != nil {
		migrationManager := config.NewMigrationManager(db.DB, "migrations")
		if err := migrationManager.ApplyMigrations(); err != nil {
			log.Printf("Warning: Failed to apply migrations: %v", err)
		}

		// 创建默认用户
		if err := config.CreateDefaultUserIfNotExist(nil, cfg); err != nil {
			log.Printf("Warning: Failed to create default user: %v", err)
		}
	}

	// 创建健康检查器
	healthChecker := utils.NewHealthChecker(db)

	// 创建Echo实例
	e := echo.New()

	// 设置验证器
	e.Validator = middleware.NewValidator()

	// 设置自定义错误处理器
	e.HTTPErrorHandler = middleware.ErrorHandler()

	// 基础中间件
	e.Use(middleware.LoggerConfig())
	e.Use(middleware.RecoverConfig())
	e.Use(middleware.CORSConfig())

	// 健康检查端点
	e.GET("/health", healthChecker.Check)

	// 初始化仓库层和服务层
	var services *service.Services

	if db != nil {
		// 创建Ent客户端
		entClient := ent.NewClient(ent.Driver(sql.OpenDB(dialect.Postgres, db.DB)))

		// 创建repositories和services
		repos := repository.NewRepositories(entClient)
		services = service.NewServices(repos, cfg.JWT.Secret, cfg.JWT.ExpiryHour)

		log.Println("Ent client and services initialized successfully")
	}

	// 检查是否有 --debug 参数
	for _, arg := range os.Args[1:] {
		if arg == "--debug" || arg == "-debug" {
			log.Println("[DEBUG] Debug mode enabled: 请求与响应内容将被打印")
			middleware.DebugMode = true
			break
		}
	}

	// 设置路由
	if services != nil {
		if middleware.DebugMode {
			e.Use(middleware.DebugLoggerMiddleware)
		}
		// 不要全局 use JWT，中间件由 setupRoutes 内部分组控制
		setupRoutes(e, services, cfg.JWT.Secret)
	} else {
		// 如果服务未初始化，添加基础测试路由
		api := e.Group("/api")
		api.GET("/ping", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
		})
		api.GET("/error-test", func(c echo.Context) error {
			return echo.NewHTTPError(http.StatusBadRequest, "这是一个测试错误")
		})
	}

	// 启动服务器
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	// 优雅关闭
	go func() {
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on %s", serverAddr)

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
