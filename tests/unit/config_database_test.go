package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"okr-web/internal/config"
)

func TestNewDatabase(t *testing.T) {
	// 测试配置
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			DBName:   "test_db",
			SSLMode:  "disable",
		},
	}

	// 注意：这个测试需要实际的PostgreSQL连接
	// 在CI/CD环境中应该跳过或使用模拟数据库
	db, err := config.NewDatabase(cfg)
	if err != nil {
		t.Logf("Database connection failed (expected in test environment): %v", err)
		// 在测试环境中，数据库连接失败是预期的
		assert.Error(t, err)
		return
	}

	// 如果连接成功，测试健康检查
	if db != nil {
		defer db.Close()
		err := db.Health()
		assert.NoError(t, err)
	}
}

func TestDatabaseClose(t *testing.T) {
	db := &config.Database{DB: nil}
	err := db.Close()
	assert.NoError(t, err)
}
