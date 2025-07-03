package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// 测试默认配置
	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, "okr_db", cfg.Database.DBName)
	assert.Equal(t, "info", cfg.Log.Level)
}

func TestEnvironmentOverride(t *testing.T) {
	// 设置环境变量
	os.Setenv("SERVER_HOST", "127.0.0.1")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DATABASE_HOST", "db.example.com")
	os.Setenv("LOG_LEVEL", "debug")

	defer func() {
		// 清理环境变量
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DATABASE_HOST")
		os.Unsetenv("LOG_LEVEL")
	}()

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "127.0.0.1", cfg.Server.Host)
	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, "db.example.com", cfg.Database.Host)
	assert.Equal(t, "debug", cfg.Log.Level)
}

func TestValidation(t *testing.T) {
	cfg := &Config{}
	setDefaults(cfg)

	// 测试有效配置
	err := validate(cfg)
	assert.NoError(t, err)

	// 测试无效端口
	cfg.Server.Port = -1
	err = validate(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid server port")

	// 重置端口并测试空数据库名
	cfg.Server.Port = 8080
	cfg.Database.DBName = ""
	err = validate(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database name cannot be empty")
}

func TestDatabaseURL(t *testing.T) {
	cfg := &Config{}
	setDefaults(cfg)

	expected := "postgres://postgres:password@localhost:5432/okr_db?sslmode=disable"
	assert.Equal(t, expected, cfg.DatabaseURL())
}
