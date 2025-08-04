package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type ServerConfig struct {
	Host string `ini:"host"`
	Port int    `ini:"port"`
}

type DatabaseConfig struct {
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	User     string `ini:"user"`
	Password string `ini:"password"`
	DBName   string `ini:"dbname"`
	SSLMode  string `ini:"sslmode"`
}

type LogConfig struct {
	Level string `ini:"level"`
}

type Config struct {
	Server   ServerConfig   `ini:"server"`
	Database DatabaseConfig `ini:"database"`
	Log      LogConfig      `ini:"log"`
}

var Cfg *Config

func InitConfig(configPath string) {
	if configPath == "" {
		configPath = "configs/config.ini"
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			// for test
			configPath = "../configs/config.ini"
		}
	}

	cfg, err := ini.Load(configPath)
	if err != nil {
		log.Fatalf("Fail to read config file %s: %v", configPath, err)
	}

	Cfg = &Config{}
	err = cfg.MapTo(Cfg)
	if err != nil {
		log.Fatalf("Fail to map config: %v", err)
	}

	log.Printf("Config loaded from: %s", configPath)
}
