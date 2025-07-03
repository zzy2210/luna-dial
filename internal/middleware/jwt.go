package middleware

import (
	"log"
	"net/http"
	"strings"

	"okr-web/internal/service"
	"okr-web/internal/types"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTConfig JWT中间件配置
type JWTConfig struct {
	Skipper      func(c echo.Context) bool
	SigningKey   string
	UserService  service.UserService
	TokenLookup  string
	ContextKey   string
	ErrorHandler func(c echo.Context, err error) error
}

// JWTWithConfig 创建JWT中间件
func JWTWithConfig(config JWTConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = func(c echo.Context) bool {
			return false
		}
	}
	if config.TokenLookup == "" {
		config.TokenLookup = "header:Authorization"
	}
	if config.ContextKey == "" {
		config.ContextKey = "user"
	}
	if config.ErrorHandler == nil {
		config.ErrorHandler = defaultJWTErrorHandler
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			// 提取token
			token, err := extractToken(c, config.TokenLookup)
			if err != nil {
				return config.ErrorHandler(c, err)
			}

			// 验证token
			claims, err := validateToken(token, config.SigningKey)
			if err != nil {
				return config.ErrorHandler(c, err)
			}

			// 日志输出 claims.UserID 及类型
			log.Printf("[JWT] claims.UserID: %v, type: %T", claims.UserID, claims.UserID)

			// 统一用 string 注入 context
			c.Set("user_id", claims.UserID.String())
			c.Set("username", claims.Username)

			// 日志输出 context 取值类型
			log.Printf("[JWT] context user_id type: %T", c.Get("user_id"))

			return next(c)
		}
	}
}

// JWT 使用默认配置的JWT中间件
func JWT(signingKey string, userService service.UserService) echo.MiddlewareFunc {
	return JWTWithConfig(JWTConfig{
		SigningKey:  signingKey,
		UserService: userService,
	})
}

// extractToken 从请求中提取token
func extractToken(c echo.Context, tokenLookup string) (string, error) {
	parts := strings.Split(tokenLookup, ":")
	if len(parts) != 2 {
		return "", &types.AppError{
			Code:    500,
			Message: "无效的token查找配置",
			Type:    "INVALID_TOKEN_LOOKUP",
		}
	}

	switch parts[0] {
	case "header":
		auth := c.Request().Header.Get(parts[1])
		if auth == "" {
			return "", &types.AppError{
				Code:    401,
				Message: "缺少Authorization头",
				Type:    "MISSING_AUTH_HEADER",
			}
		}

		// 检查Bearer前缀
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(auth, bearerPrefix) {
			return "", &types.AppError{
				Code:    401,
				Message: "无效的Authorization格式",
				Type:    "INVALID_AUTH_FORMAT",
			}
		}

		return auth[len(bearerPrefix):], nil

	case "query":
		token := c.QueryParam(parts[1])
		if token == "" {
			return "", &types.AppError{
				Code:    401,
				Message: "缺少token参数",
				Type:    "MISSING_TOKEN_PARAM",
			}
		}
		return token, nil

	case "cookie":
		cookie, err := c.Cookie(parts[1])
		if err != nil {
			return "", &types.AppError{
				Code:    401,
				Message: "缺少token cookie",
				Type:    "MISSING_TOKEN_COOKIE",
			}
		}
		return cookie.Value, nil

	default:
		return "", &types.AppError{
			Code:    500,
			Message: "不支持的token提取方式",
			Type:    "UNSUPPORTED_TOKEN_EXTRACTION",
		}
	}
}

// validateToken 验证JWT token
func validateToken(tokenString, signingKey string) (*service.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &service.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &types.AppError{
				Code:    401,
				Message: "无效的签名方法",
				Type:    "INVALID_SIGNING_METHOD",
			}
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return nil, &types.AppError{
			Code:    401,
			Message: "无效的token",
			Type:    "INVALID_TOKEN",
		}
	}

	if claims, ok := token.Claims.(*service.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, &types.AppError{
		Code:    401,
		Message: "token验证失败",
		Type:    "TOKEN_VALIDATION_FAILED",
	}
}

// defaultJWTErrorHandler 默认JWT错误处理器
func defaultJWTErrorHandler(c echo.Context, err error) error {
	if appErr, ok := err.(*types.AppError); ok {
		return c.JSON(appErr.Code, types.ErrorResponse{
			Success: false,
			Error:   appErr.Type,
			Message: appErr.Message,
		})
	}

	return c.JSON(http.StatusUnauthorized, types.ErrorResponse{
		Success: false,
		Error:   "UNAUTHORIZED",
		Message: "未授权访问2",
	})
}

// SkipAuth 跳过认证的路径
func SkipAuth(path ...string) func(c echo.Context) bool {
	skipPaths := make(map[string]bool)
	for _, p := range path {
		skipPaths[p] = true
	}

	return func(c echo.Context) bool {
		return skipPaths[c.Path()]
	}
}

// RequireAuth 要求认证的中间件，用于特定路由
func RequireAuth(signingKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 提取token
			token, err := extractToken(c, "header:Authorization")
			if err != nil {
				return defaultJWTErrorHandler(c, err)
			}

			// 验证token
			claims, err := validateToken(token, signingKey)
			if err != nil {
				return defaultJWTErrorHandler(c, err)
			}

			// 将用户信息存储到context
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)

			return next(c)
		}
	}
}
