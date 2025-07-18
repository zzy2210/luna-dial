package model

import "errors"

// 数据库层错误定义
var (
	// 通用数据库错误
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateKey   = errors.New("duplicate key violation")
	ErrForeignKey     = errors.New("foreign key constraint violation")
	ErrConnection     = errors.New("database connection error")

	// 权限相关错误
	ErrUnauthorized = errors.New("unauthorized access")

	// 数据完整性错误
	ErrDataIntegrity = errors.New("data integrity violation")
)
