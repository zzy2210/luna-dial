package controller

import (
	"net/http"

	"okr-web/internal/middleware"
	"okr-web/internal/service"
	"okr-web/internal/types"

	"github.com/labstack/echo/v4"
)

// UserController 用户控制器
type UserController struct {
	userService service.UserService
}

// NewUserController 创建用户控制器
func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UpdateUserRequest 更新用户请求结构
type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6"`
}

// Register 用户注册
func (ctrl *UserController) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_REQUEST_BODY",
			Message: "请求体格式无效",
		})
	}

	// 参数验证
	if req.Username == "" {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "MISSING_USERNAME",
			Message: "用户名不能为空",
		})
	}
	if req.Email == "" {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "MISSING_EMAIL",
			Message: "邮箱不能为空",
		})
	}
	if req.Password == "" {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "MISSING_PASSWORD",
			Message: "密码不能为空",
		})
	}
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_USERNAME_LENGTH",
			Message: "用户名长度必须在3-50字符之间",
		})
	}
	if len(req.Password) < 6 {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_PASSWORD_LENGTH",
			Message: "密码长度至少6个字符",
		})
	}

	// 调用服务层
	serviceReq := service.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := ctrl.userService.Register(c.Request().Context(), serviceReq)
	if err != nil {
		if appErr, ok := err.(*types.AppError); ok {
			return c.JSON(appErr.Code, types.ErrorResponse{
				Success: false,
				Error:   appErr.Type,
				Message: appErr.Message,
			})
		}
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Success: false,
			Error:   "INTERNAL_SERVER_ERROR",
			Message: "服务器内部错误",
		})
	}

	// 隐藏密码字段
	user.Password = ""

	return c.JSON(http.StatusCreated, types.SuccessResponse{
		Success: true,
		Data:    user,
		Message: "用户注册成功",
	})
}

// Login 用户登录
func (ctrl *UserController) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_REQUEST_BODY",
			Message: "请求体格式无效",
		})
	}

	// 参数验证
	if req.Username == "" {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "MISSING_USERNAME",
			Message: "用户名不能为空",
		})
	}
	if req.Password == "" {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "MISSING_PASSWORD",
			Message: "密码不能为空",
		})
	}

	// 调用服务层
	serviceReq := service.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	authResponse, err := ctrl.userService.Login(c.Request().Context(), serviceReq)
	if err != nil {
		if appErr, ok := err.(*types.AppError); ok {
			return c.JSON(appErr.Code, types.ErrorResponse{
				Success: false,
				Error:   appErr.Type,
				Message: appErr.Message,
			})
		}
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Success: false,
			Error:   "INTERNAL_SERVER_ERROR",
			Message: "服务器内部错误",
		})
	}

	// 隐藏密码字段
	authResponse.User.Password = ""

	return c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    authResponse,
		Message: "登录成功",
	})
}

// GetMe 获取当前用户信息
func (ctrl *UserController) GetMe(c echo.Context) error {
	// 从AuthMiddleware获取已解析的用户ID
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return middleware.HandleUnauthorized(c, err)
	}

	// 调用服务层
	user, err := ctrl.userService.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		if appErr, ok := err.(*types.AppError); ok {
			return c.JSON(appErr.Code, types.ErrorResponse{
				Success: false,
				Error:   appErr.Type,
				Message: appErr.Message,
			})
		}
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Success: false,
			Error:   "INTERNAL_SERVER_ERROR",
			Message: "服务器内部错误",
		})
	}

	// 隐藏密码字段
	user.Password = ""

	return c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    user,
		Message: "获取用户信息成功",
	})
}

// UpdateMe 更新当前用户信息
func (ctrl *UserController) UpdateMe(c echo.Context) error {
	// 从AuthMiddleware获取已解析的用户ID
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return middleware.HandleUnauthorized(c, err)
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_REQUEST_BODY",
			Message: "请求体格式无效",
		})
	}

	// 参数验证
	if req.Password != nil && len(*req.Password) < 6 {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_PASSWORD_LENGTH",
			Message: "密码长度至少6个字符",
		})
	}

	// 调用服务层
	serviceReq := service.UpdateUserRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := ctrl.userService.UpdateUser(c.Request().Context(), userID, serviceReq)
	if err != nil {
		if appErr, ok := err.(*types.AppError); ok {
			return c.JSON(appErr.Code, types.ErrorResponse{
				Success: false,
				Error:   appErr.Type,
				Message: appErr.Message,
			})
		}
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Success: false,
			Error:   "INTERNAL_SERVER_ERROR",
			Message: "服务器内部错误",
		})
	}

	// 隐藏密码字段
	user.Password = ""

	return c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    user,
		Message: "用户信息更新成功",
	})
}

// Logout 用户登出（客户端处理，服务端只返回成功）
func (ctrl *UserController) Logout(c echo.Context) error {
	return c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    nil,
		Message: "登出成功",
	})
}

// GetUserInfo 获取用户信息
func (c *UserController) GetUserInfo(ctx echo.Context) error {
	// 从AuthMiddleware获取已解析的用户ID
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return middleware.HandleUnauthorized(ctx, err)
	}

	user, err := c.userService.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, types.ErrorResponse{
			Message: "用户不存在",
			Code:    http.StatusNotFound,
		})
	}

	return ctx.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    user,
	})
}

// UpdateUser 更新用户信息
func (c *UserController) UpdateUser(ctx echo.Context) error {
	// 从AuthMiddleware获取已解析的用户ID
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return middleware.HandleUnauthorized(ctx, err)
	}

	var req service.UpdateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
			Message: "请求参数格式错误",
			Code:    http.StatusBadRequest,
		})
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
			Message: "请求参数验证失败",
			Code:    http.StatusBadRequest,
		})
	}

	user, err := c.userService.UpdateUser(ctx.Request().Context(), userID, req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Message: "更新用户信息失败",
			Code:    http.StatusInternalServerError,
		})
	}

	return ctx.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    user,
	})
}

// RegisterRoutes 注册用户相关路由
func (ctrl *UserController) RegisterRoutes(e *echo.Echo, jwtMiddleware echo.MiddlewareFunc) {
	// 公开路由（不需要认证）
	auth := e.Group("/api/auth")
	auth.POST("/register", ctrl.Register)
	auth.POST("/login", ctrl.Login)

	// 需要认证的路由
	authProtected := e.Group("/api/auth")
	authProtected.Use(jwtMiddleware)
	authProtected.GET("/me", ctrl.GetMe)
	authProtected.PUT("/me", ctrl.UpdateMe)
	authProtected.POST("/logout", ctrl.Logout)
}
