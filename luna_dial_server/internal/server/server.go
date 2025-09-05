package server

import (
	"context"
	"fmt"
	"luna_dial/internal/config"
	"luna_dial/internal/data"
	"luna_dial/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewServer(ctx context.Context, dataInstance *data.Data) *echo.Echo {
    e := echo.New()

    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Validator
    e.Validator = service.NewSimpleValidator()

	// Initialize service
	svc := service.NewService(ctx, e, dataInstance)

	// Setup routes
	svc.SetupRouter()

	return e
}

func Start(e *echo.Echo, cleanup func()) {
	serverAddr := fmt.Sprintf("%s:%d", config.Cfg.Server.Host, config.Cfg.Server.Port)

	// Start server in a goroutine
	go func() {
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	e.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	// Call cleanup function
	cleanup()

	e.Logger.Info("Server gracefully stopped")
}
