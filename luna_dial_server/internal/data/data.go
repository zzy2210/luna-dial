package data

import (
	"time"

	"gorm.io/gorm"
)

// Data 数据层结构，管理基础设施和共享服务
type Data struct {
	DB             *gorm.DB
	SystemConfig   *SystemConfig   // 导出SystemConfig供service层使用
	SessionManager SessionManager  // 导出SessionManager供service层使用
}

// NewData 创建数据层实例
func NewData(db *gorm.DB) (*Data, func(), error) {
	// 创建Session管理器，90分钟超时
	sessionManager := NewMemorySessionManager(90 * time.Minute)
	
	d := &Data{
		DB:             db,
		SystemConfig:   NewSystemConfig(db),
		SessionManager: sessionManager,
	}

	cleanup := func() {
		// 关闭Session管理器
		if sessionManager != nil {
			sessionManager.Close()
		}
		
		// 关闭数据库连接
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}

	return d, cleanup, nil
}
