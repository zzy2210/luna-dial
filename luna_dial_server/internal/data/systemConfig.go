package data

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/tjfoc/gmsm/sm3"
	"gorm.io/gorm"
)

// SystemConfig 系统配置结构
type SystemConfig struct {
	db *gorm.DB
}

// CryptoKeys 加密密钥配置
type CryptoKeys struct {
	// SM2密钥对，用于可能的加密需求
	SM2PrivateKey string `json:"sm2_private_key"`
	SM2PublicKey  string `json:"sm2_public_key"`
	// JWT密钥
	JWTSecret string `json:"jwt_secret"`
}

// SystemConfigRecord 系统配置数据模型
type SystemConfigRecord struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Key       string    `gorm:"uniqueIndex;type:varchar(100);not null;column:config_key" json:"key"`
	Value     string    `gorm:"type:text;not null;column:config_value" json:"value"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (SystemConfigRecord) TableName() string {
	return "system_configs"
}

func NewSystemConfig(db *gorm.DB) *SystemConfig {
	return &SystemConfig{db: db}
}

// InitializeSystem 初始化系统配置
func (sc *SystemConfig) InitializeSystem(ctx context.Context) error {
	log.Println("开始初始化系统...")

	// 1. 总是运行数据库迁移（让migrate自己判断是否需要执行）
	if err := sc.runMigrations(ctx); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	// 2. 检查并初始化基础数据（密钥、管理员用户）
	if !sc.IsBasicDataInitialized(ctx) {
		if err := sc.initializeBasicData(ctx); err != nil {
			return fmt.Errorf("初始化基础数据失败: %w", err)
		}
	} else {
		log.Println("基础数据已存在，跳过初始化")
	}

	log.Println("系统初始化完成")
	return nil
}

// IsBasicDataInitialized 检查基础数据是否已初始化
func (sc *SystemConfig) IsBasicDataInitialized(ctx context.Context) bool {
	// 检查JWT密钥
	if _, err := sc.GetJWTSecret(ctx); err != nil {
		return false
	}

	// 检查管理员用户
	var adminUser User
	err := sc.db.WithContext(ctx).Where("user_name = ?", "admin").First(&adminUser).Error
	if err != nil {
		return false
	}

	log.Println("基础数据已存在")
	return true
}

// initializeBasicData 初始化基础数据
func (sc *SystemConfig) initializeBasicData(ctx context.Context) error {
	log.Println("初始化基础数据...")

	// 生成密钥
	if _, err := sc.generateCryptoKeys(ctx); err != nil {
		return fmt.Errorf("生成密钥失败: %w", err)
	}

	// 创建管理员用户
	if err := sc.createAdminUser(ctx); err != nil {
		return fmt.Errorf("创建管理员用户失败: %w", err)
	}

	log.Println("基础数据初始化完成")
	return nil
}

// runMigrations 运行数据库迁移
func (sc *SystemConfig) runMigrations(ctx context.Context) error {
	log.Println("检查并运行数据库迁移...")

	// 获取数据库底层连接
	sqlDB, err := sc.db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 创建postgres驱动实例
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("创建数据库驱动失败: %w", err)
	}

	// 获取迁移文件路径
	migrationsPath := "file://migrations"

	// 创建迁移实例
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("创建迁移实例失败: %w", err)
	}
	defer m.Close()

	// 运行迁移 - migrate会自动判断哪些需要执行
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("执行迁移失败: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("数据库已是最新版本，无需迁移")
	} else {
		log.Println("数据库迁移完成")
	}

	return nil
}

// generateCryptoKeys 生成加密密钥
func (sc *SystemConfig) generateCryptoKeys(ctx context.Context) (*CryptoKeys, error) {
	log.Println("生成加密密钥...")

	// 检查是否已经存在密钥
	if keys, err := sc.getCryptoKeys(ctx); err == nil && keys != nil {
		log.Println("密钥已存在，跳过生成")
		return keys, nil
	}

	// 生成JWT密钥
	jwtSecret := make([]byte, 32)
	if _, err := rand.Read(jwtSecret); err != nil {
		return nil, fmt.Errorf("生成JWT密钥失败: %w", err)
	}

	keys := &CryptoKeys{
		JWTSecret: hex.EncodeToString(jwtSecret),
	}

	// 保存密钥到数据库
	if err := sc.saveCryptoKeys(ctx, keys); err != nil {
		return nil, fmt.Errorf("保存密钥失败: %w", err)
	}

	log.Println("密钥生成并保存完成")
	return keys, nil
}

// getCryptoKeys 获取加密密钥
func (sc *SystemConfig) getCryptoKeys(ctx context.Context) (*CryptoKeys, error) {
	var configs []SystemConfigRecord
	err := sc.db.WithContext(ctx).
		Where("config_key IN (?)", []string{"jwt_secret"}).
		Find(&configs).Error

	if err != nil {
		return nil, err
	}

	if len(configs) != 1 {
		return nil, fmt.Errorf("密钥配置不完整")
	}

	keys := &CryptoKeys{}
	for _, config := range configs {
		switch config.Key {
		case "jwt_secret":
			keys.JWTSecret = config.Value
		}
	}

	return keys, nil
}

// saveCryptoKeys 保存加密密钥
func (sc *SystemConfig) saveCryptoKeys(ctx context.Context, keys *CryptoKeys) error {
	config := SystemConfigRecord{
		ID:    uuid.New().String(),
		Key:   "jwt_secret",
		Value: keys.JWTSecret,
	}

	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()
	if err := sc.db.WithContext(ctx).Create(&config).Error; err != nil {
		return err
	}

	return nil
}

// createAdminUser 创建管理员用户
func (sc *SystemConfig) createAdminUser(ctx context.Context) error {
	log.Println("创建管理员用户...")

	// 检查是否已存在管理员用户
	var existingUser User
	err := sc.db.WithContext(ctx).Where("user_name = ?", "admin").First(&existingUser).Error
	if err == nil {
		log.Println("管理员用户已存在，跳过创建")
		return nil
	}

	// 使用与biz层相同的密码哈希方式
	hasher := sm3.New()
	hasher.Write([]byte("admin@123"))
	hashedPassword := hasher.Sum(nil)

	// 创建管理员用户
	admin := User{
		ID:        uuid.New().String(),
		UserName:  "admin",
		Name:      "系统管理员",
		Email:     "admin@luna-dial.com",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := sc.db.WithContext(ctx).Create(&admin).Error; err != nil {
		return err
	}

	log.Println("管理员用户创建完成")
	return nil
}

// GetJWTSecret 获取JWT密钥
func (sc *SystemConfig) GetJWTSecret(ctx context.Context) (string, error) {
	keys, err := sc.getCryptoKeys(ctx)
	if err != nil {
		return "", err
	}
	return keys.JWTSecret, nil
}

// IsSystemInitialized 检查系统是否已初始化
func (sc *SystemConfig) IsSystemInitialized(ctx context.Context) bool {
	return sc.IsBasicDataInitialized(ctx)
}
