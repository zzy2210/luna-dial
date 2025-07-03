package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator 自定义验证器
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator 创建新的验证器实例
func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

// Validate 验证结构体
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// 可以在这里自定义错误信息格式
		return echo.NewHTTPError(400, err.Error())
	}
	return nil
}
