package main

import (
	"context"
	"log"

	"luna_dial/internal/data"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// 数据库连接配置
	dsn := "host=localhost user=okr_user password=your-password-word dbname=okr_db port=15432 sslmode=disable TimeZone=Asia/Shanghai"

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 创建系统配置实例
	systemConfig := data.NewSystemConfig(db)
	ctx := context.Background()

	// 检查并初始化系统
	if !systemConfig.IsSystemInitialized(ctx) {
		log.Println("系统未初始化，开始初始化...")
		if err := systemConfig.InitializeSystem(ctx); err != nil {
			log.Fatal("系统初始化失败:", err)
		}
		log.Println("系统初始化完成")
	} else {
		log.Println("系统已初始化")
	}

	// 现在可以创建你的Repository实例并启动应用
	_ = data.NewTaskRepo(db)
	_ = data.NewJournalRepo(db)
	_ = data.NewUserRepo(db)

	log.Println("Repository创建完成，应用可以启动...")

	// 这里可以继续你的应用逻辑
	// 例如：启动HTTP服务器等
}
