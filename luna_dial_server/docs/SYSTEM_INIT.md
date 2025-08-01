# 系统初始化使用说明

## 功能特性

1. **数据库迁移管理**: 使用 `golang-migrate` 进行数据库结构管理，支持版本控制
2. **密钥生成**: 自动生成 JWT 密钥用于认证
3. **管理员用户**: 自动创建默认管理员账户，使用与biz层相同的SM3哈希
4. **升级友好**: 支持后续添加新的迁移文件，总是检查最新版本

## 核心设计原则

### 分离关注点
- **数据库迁移**: 总是检查并运行最新的迁移文件
- **基础数据**: 仅在首次运行时创建（JWT密钥、管理员用户）

### 幂等性保证
- 可以反复运行初始化脚本
- 迁移由 `golang-migrate` 自动管理版本
- 基础数据检查存在性，避免重复创建

## 使用方法

### 1. 运行初始化

```bash
# 进入项目目录
cd /home/y1nhui/work/github_own/luna-dial/luna_dial_server

# 运行初始化脚本
go run cmd/init/main.go
```

### 2. 配置数据库连接

修改 `cmd/init/main.go` 中的数据库连接字符串：

```go
dsn := "host=localhost user=okr_user password=your-password-word dbname=okr_db port=15432 sslmode=disable TimeZone=Asia/Shanghai"
```

### 3. 默认管理员账户

- 用户名: `admin`
- 密码: `admin@123` (使用SM3哈希存储，与biz层一致)
- 邮箱: `admin@luna-dial.com`

## 密钥管理

系统会自动生成：

1. **JWT密钥**: 用于生成和验证JWT令牌

密钥存储在 `system_configs` 表中，key为 `jwt_secret`。

## 数据库迁移

### 迁移文件管理

迁移文件位于 `migrations/` 目录：

- `0001_init_schema.up.sql`: 创建基础表结构
- `0002_init_system_data.up.sql`: 系统数据初始化标记
- 未来可添加: `0003_xxx.up.sql`, `0004_xxx.up.sql` 等

### 升级流程

1. 添加新的迁移文件到 `migrations/` 目录
2. 运行 `go run cmd/init/main.go`
3. 系统会自动检测并执行新的迁移

```bash
# 示例：添加新功能后
echo "ALTER TABLE users ADD COLUMN avatar VARCHAR(255);" > migrations/0003_add_user_avatar.up.sql
echo "ALTER TABLE users DROP COLUMN avatar;" > migrations/0003_add_user_avatar.down.sql

# 运行初始化，会自动执行新迁移
go run cmd/init/main.go
```

## API 使用示例

### 在业务代码中使用

```go
package main

import (
    "context"
    "luna_dial/internal/data"
    "gorm.io/gorm"
)

func main() {
    // 假设你已经有了数据库连接
    var db *gorm.DB
    
    // 创建系统配置实例
    systemConfig := data.NewSystemConfig(db)
    ctx := context.Background()
    
    // 总是运行初始化（迁移检查 + 基础数据检查）
    if err := systemConfig.InitializeSystem(ctx); err != nil {
        log.Fatal("系统初始化失败:", err)
    }
    
    // 获取JWT密钥
    jwtSecret, err := systemConfig.GetJWTSecret(ctx)
    if err != nil {
        log.Printf("获取JWT密钥失败: %v", err)
    }
}
```

## 设计优势

### 1. 迁移总是检查
每次运行都会检查是否有新的迁移文件需要执行，确保数据库结构始终是最新的。

### 2. 基础数据幂等
JWT密钥和管理员用户只在首次运行时创建，后续运行会跳过。

### 3. 升级友好
添加新的迁移文件后，只需运行初始化脚本即可自动升级数据库。

### 4. 一致的密码处理
管理员用户使用与biz层相同的SM3哈希算法，确保登录验证一致。

## 注意事项

1. **迁移文件命名**: 遵循 `NNNN_description.up.sql` 格式
2. **密钥安全**: JWT密钥存储在数据库中，请确保数据库访问的安全性
3. **生产环境**: 建议在生产环境中修改默认管理员密码
4. **备份**: 在生产环境执行迁移前，请先备份数据库
