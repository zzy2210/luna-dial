package main

import (
	"context"
	"fmt"
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

	// 创建数据层实例
	dataInstance, cleanup, err := data.NewData(db)
	if err != nil {
		log.Fatal("创建数据层实例失败:", err)
	}
	defer cleanup() // 确保资源清理

	ctx := context.Background()

	// 总是运行初始化（包括迁移检查和基础数据初始化）
	fmt.Println("开始系统初始化检查...")
	if err := dataInstance.SystemConfig.InitializeSystem(ctx); err != nil {
		log.Fatal("系统初始化失败:", err)
	}

	// 显示系统信息
	jwtSecret, err := dataInstance.SystemConfig.GetJWTSecret(ctx)
	if err != nil {
		log.Printf("获取JWT密钥失败: %v", err)
	} else {
		fmt.Printf("JWT密钥长度: %d\n", len(jwtSecret))
	}

	fmt.Println("管理员账号: admin")
	fmt.Println("管理员密码: admin@123")
	fmt.Println("系统就绪!")
}
