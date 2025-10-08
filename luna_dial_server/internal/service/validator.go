package service

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator 是 Echo Validator 接口的适配器
// 使用 go-playground/validator/v10 进行数据验证
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator 创建一个新的 CustomValidator 实例
func NewValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

// Validate 实现 Echo 的 Validator 接口
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}
