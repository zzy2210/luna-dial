package main

import (
	"context"
	"fmt"
	"log"
	"luna_dial/internal/config"
	"luna_dial/internal/data"
	"luna_dial/internal/server"

	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	configFile string
)

var rootCmd = &cobra.Command{
	Use:   "luna-dial-server",
	Short: "Luna Dial Server - A task and journal management system",
	Long: `Luna Dial Server is a backend service for task and journal management.
It provides APIs for user authentication, task management, journaling, and planning.`,
	Run: runServer,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (default is configs/config.ini)")
}

func runServer(cmd *cobra.Command, args []string) {
	// Initialize config
	config.InitConfig(configFile)

	// Initialize database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.Cfg.Database.Host,
		config.Cfg.Database.User,
		config.Cfg.Database.Password,
		config.Cfg.Database.DBName,
		config.Cfg.Database.Port,
		config.Cfg.Database.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Initialize data
	dataInstance, sessionCleanup, err := data.NewData(db)
	if err != nil {
		log.Fatalf("failed to initialize data: %v", err)
	}

	// Create context
	ctx := context.Background()

	// Initialize system (migrations, admin user, etc.)
	log.Println("Checking system initialization...")
	if err := dataInstance.SystemConfig.InitializeSystem(ctx); err != nil {
		log.Fatalf("failed to initialize system: %v", err)
	}

	// Create server
	e := server.NewServer(ctx, dataInstance)

	log.Printf("Starting server on %s:%d", config.Cfg.Server.Host, config.Cfg.Server.Port)

	// Create comprehensive cleanup function
	cleanup := func() {
		// 首先关闭会话管理器等资源
		sessionCleanup()

		// 然后关闭数据库连接
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}

	// Start server with graceful shutdown
	server.Start(e, cleanup)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
