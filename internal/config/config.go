package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/ini.v1"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `ini:"server"`
	Database DatabaseConfig `ini:"database"`
	Log      LogConfig      `ini:"log"`
	JWT      JWTConfig      `ini:"jwt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `ini:"host"`
	Port int    `ini:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	User     string `ini:"user"`
	Password string `ini:"password"`
	DBName   string `ini:"dbname"`
	SSLMode  string `ini:"sslmode"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `ini:"level"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string `ini:"secret"`
	ExpiryHour int    `ini:"expiry_hour"`
}

// Load 加载配置文件
func Load() (*Config, error) {
	cfg := &Config{}

	// 设置默认值
	setDefaults(cfg)

	// 尝试加载配置文件
	if _, err := os.Stat("config.ini"); err == nil {
		iniFile, err := ini.Load("config.ini")
		if err != nil {
			return nil, fmt.Errorf("failed to load config.ini: %v", err)
		}

		if err := iniFile.MapTo(cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config.ini: %v", err)
		}
	}

	// 环境变量覆盖
	overrideWithEnv(cfg)

	// 验证配置
	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %v", err)
	}

	return cfg, nil
}

// setDefaults 设置默认配置值
func setDefaults(cfg *Config) {
	cfg.Server.Host = "0.0.0.0"
	cfg.Server.Port = 8080

	cfg.Database.Host = "localhost"
	cfg.Database.Port = 5432
	cfg.Database.User = "postgres"
	cfg.Database.Password = "password"
	cfg.Database.DBName = "okr_db"
	cfg.Database.SSLMode = "disable"

	cfg.Log.Level = "info"

	cfg.JWT.Secret = "your-secret-key-change-in-production"
	cfg.JWT.ExpiryHour = 24
}

// overrideWithEnv 使用环境变量覆盖配置
func overrideWithEnv(cfg *Config) {
	// Server配置
	if host := os.Getenv("SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Server.Port = p
		}
	}

	// Database配置
	if host := os.Getenv("DATABASE_HOST"); host != "" {
		cfg.Database.Host = host
	}
	if port := os.Getenv("DATABASE_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Database.Port = p
		}
	}
	if user := os.Getenv("DATABASE_USER"); user != "" {
		cfg.Database.User = user
	}
	if password := os.Getenv("DATABASE_PASSWORD"); password != "" {
		cfg.Database.Password = password
	}
	if dbname := os.Getenv("DATABASE_DBNAME"); dbname != "" {
		cfg.Database.DBName = dbname
	}
	if sslmode := os.Getenv("DATABASE_SSLMODE"); sslmode != "" {
		cfg.Database.SSLMode = sslmode
	}

	// Log配置
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cfg.Log.Level = level
	}

	// JWT配置
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		cfg.JWT.Secret = secret
	}
	if expiry := os.Getenv("JWT_EXPIRY_HOUR"); expiry != "" {
		if h, err := strconv.Atoi(expiry); err == nil {
			cfg.JWT.ExpiryHour = h
		}
	}
}

// validate 验证配置
func validate(cfg *Config) error {
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}

	if cfg.Database.Port < 1 || cfg.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", cfg.Database.Port)
	}

	if cfg.Database.DBName == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[cfg.Log.Level] {
		return fmt.Errorf("invalid log level: %s", cfg.Log.Level)
	}

	if len(cfg.JWT.Secret) < 10 {
		return fmt.Errorf("JWT secret must be at least 10 characters long")
	}

	if cfg.JWT.ExpiryHour < 1 {
		return fmt.Errorf("JWT expiry hour must be positive")
	}

	return nil
}

// DatabaseURL 返回数据库连接URL
func (cfg *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)
}
