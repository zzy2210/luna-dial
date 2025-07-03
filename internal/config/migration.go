package config

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/argon2"
)

// Migration 表示一个数据库迁移
type Migration struct {
	Version int
	Name    string
	UpSQL   string
	DownSQL string
}

// MigrationManager 迁移管理器
type MigrationManager struct {
	db            *sql.DB
	migrationsDir string
}

// NewMigrationManager 创建新的迁移管理器
func NewMigrationManager(db *sql.DB, migrationsDir string) *MigrationManager {
	return &MigrationManager{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// createMigrationsTable 创建迁移记录表
func (m *MigrationManager) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := m.db.Exec(query)
	return err
}

// getAppliedMigrations 获取已应用的迁移
func (m *MigrationManager) getAppliedMigrations() (map[int]bool, error) {
	applied := make(map[int]bool)

	rows, err := m.db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return applied, err
	}
	defer rows.Close()

	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return applied, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// loadMigrations 加载所有迁移文件
func (m *MigrationManager) loadMigrations() ([]*Migration, error) {
	files, err := ioutil.ReadDir(m.migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	migrations := make(map[int]*Migration)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		if !strings.HasSuffix(filename, ".sql") {
			continue
		}

		// 解析文件名: 001_create_users_table.up.sql
		parts := strings.Split(filename, "_")
		if len(parts) < 2 {
			continue
		}

		versionStr := parts[0]
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			continue
		}

		// 确定是up还是down迁移
		isUp := strings.Contains(filename, ".up.sql")
		isDown := strings.Contains(filename, ".down.sql")

		if !isUp && !isDown {
			continue
		}

		// 读取SQL内容
		filePath := filepath.Join(m.migrationsDir, filename)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		// 获取或创建migration
		migration := migrations[version]
		if migration == nil {
			// 从文件名提取名称
			nameParts := parts[1:]
			name := strings.Join(nameParts, "_")
			name = strings.TrimSuffix(name, ".up.sql")
			name = strings.TrimSuffix(name, ".down.sql")

			migration = &Migration{
				Version: version,
				Name:    name,
			}
			migrations[version] = migration
		}

		// 设置SQL内容
		if isUp {
			migration.UpSQL = string(content)
		} else {
			migration.DownSQL = string(content)
		}
	}

	// 转换为切片并排序
	result := make([]*Migration, 0, len(migrations))
	for _, migration := range migrations {
		if migration.UpSQL != "" { // 只包含有up脚本的迁移
			result = append(result, migration)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Version < result[j].Version
	})

	return result, nil
}

// ApplyMigrations 应用所有未执行的迁移
func (m *MigrationManager) ApplyMigrations() error {
	// 创建迁移记录表
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// 获取已应用的迁移
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// 加载所有迁移文件
	migrations, err := m.loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// 应用未执行的迁移
	for _, migration := range migrations {
		if applied[migration.Version] {
			fmt.Printf("Migration %d (%s) already applied, skipping\n", migration.Version, migration.Name)
			continue
		}

		fmt.Printf("Applying migration %d (%s)...\n", migration.Version, migration.Name)

		// 开始事务
		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %d: %w", migration.Version, err)
		}

		// 执行迁移SQL
		if _, err := tx.Exec(migration.UpSQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %d: %w", migration.Version, err)
		}

		// 记录迁移已应用
		if _, err := tx.Exec("INSERT INTO schema_migrations (version, name) VALUES ($1, $2)", migration.Version, migration.Name); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}

		// 提交事务
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
		}

		fmt.Printf("Migration %d (%s) applied successfully\n", migration.Version, migration.Name)
	}

	fmt.Println("All migrations applied successfully")
	return nil
}

// RollbackMigration 回滚指定版本的迁移
func (m *MigrationManager) RollbackMigration(version int) error {
	// 检查迁移是否已应用
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if !applied[version] {
		return fmt.Errorf("migration %d is not applied", version)
	}

	// 加载迁移文件
	migrations, err := m.loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// 找到要回滚的迁移
	var migration *Migration
	for _, m := range migrations {
		if m.Version == version {
			migration = m
			break
		}
	}

	if migration == nil {
		return fmt.Errorf("migration %d not found", version)
	}

	if migration.DownSQL == "" {
		return fmt.Errorf("migration %d has no down script", version)
	}

	fmt.Printf("Rolling back migration %d (%s)...\n", migration.Version, migration.Name)

	// 开始事务
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction for rollback %d: %w", migration.Version, err)
	}

	// 执行回滚SQL
	if _, err := tx.Exec(migration.DownSQL); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute rollback %d: %w", migration.Version, err)
	}

	// 删除迁移记录
	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", migration.Version); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record %d: %w", migration.Version, err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit rollback %d: %w", migration.Version, err)
	}

	fmt.Printf("Migration %d (%s) rolled back successfully\n", migration.Version, migration.Name)
	return nil
}

// CreateDefaultUserIfNotExist 创建默认管理员用户（如果不存在）
func CreateDefaultUserIfNotExist(client interface{}, cfg *Config) error {
	// 这里需要使用实际的 Ent 客户端
	// 由于当前代码中 Ent 客户端可能还未完全配置，我们先使用 SQL 方式实现

	// 注意：这个函数需要在有 Ent 客户端的情况下重新实现
	// 当前使用 SQL 的临时实现
	return createDefaultUserWithSQL(cfg)
}

// createDefaultUserWithSQL 使用原生 SQL 创建默认用户的临时实现
func createDefaultUserWithSQL(cfg *Config) error {
	// 连接数据库
	db, err := NewDatabase(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// 检查用户是否已存在
	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", "admin").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check if admin user exists: %w", err)
	}

	if count > 0 {
		fmt.Println("Admin user already exists, skipping creation")
		return nil
	}

	// 使用数据库密码作为 admin 用户密码
	password := cfg.Database.Password

	// 使用与 user_service.go 相同的密码加密逻辑
	hashedPassword, err := hashPasswordForDefault(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 生成 UUID
	userID := generateUUID()

	// 创建用户
	_, err = db.DB.Exec(`
		INSERT INTO users (id, username, email, password, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`, userID, "admin", "admin@okr.local", hashedPassword)

	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	fmt.Println("✅ Admin user created successfully with username: admin")
	return nil
}

// hashPasswordForDefault 为默认用户创建密码哈希（复制自 user_service.go 的逻辑）
func hashPasswordForDefault(password string) (string, error) {
	// 密码配置（与 user_service.go 保持一致）
	time := uint32(1)
	memory := uint32(64 * 1024)
	threads := uint8(4)
	keyLen := uint32(32)

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	return fmt.Sprintf(format, argon2.Version, memory, time, threads, b64Salt, b64Hash), nil
}

// generateUUID 生成 UUID 字符串
func generateUUID() string {
	return uuid.New().String()
}
