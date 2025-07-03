package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"okr-web/internal/types"
)

// ErrorHandler 统一错误处理中间件
func ErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var (
			code = http.StatusInternalServerError
			msg  interface{}
		)

		// 检查是否是自定义应用错误
		if appErr, ok := err.(*types.AppError); ok {
			code = appErr.Code
			response := &types.ErrorResponse{
				Success: false,
				Error:   appErr.Type,
				Message: appErr.Message,
				Code:    appErr.Code,
			}
			msg = response
		} else if he, ok := err.(*echo.HTTPError); ok {
			// Echo HTTP错误
			code = he.Code
			response := &types.ErrorResponse{
				Success: false,
				Error:   "HTTP_ERROR",
				Message: he.Message.(string),
				Code:    he.Code,
			}
			msg = response
		} else {
			// 其他未知错误
			response := &types.ErrorResponse{
				Success: false,
				Error:   "INTERNAL_SERVER_ERROR",
				Message: "服务器内部错误",
				Code:    http.StatusInternalServerError,
			}
			msg = response
		}

		// 如果响应已经发送，跳过
		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				err = c.NoContent(code)
			} else {
				err = c.JSON(code, msg)
			}
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}
}
