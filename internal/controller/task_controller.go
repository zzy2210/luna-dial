package controller

import (
	"net/http"
	"strconv"
	"time"

	"okr-web/internal/middleware"
	"okr-web/internal/service"
	"okr-web/internal/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// TaskController 任务控制器
type TaskController struct {
	taskService service.TaskService
}

// NewTaskController 创建任务控制器
func NewTaskController(taskService service.TaskService) *TaskController {
	return &TaskController{
		taskService: taskService,
	}
}

// TaskRequest 任务请求结构
type TaskRequest struct {
	Title       string           `json:"title" validate:"required,min=1,max=200"`
	Description *string          `json:"description,omitempty"`
	Type        types.TaskType   `json:"type" validate:"required"`
	StartDate   *time.Time       `json:"start_date,omitempty"`
	EndDate     *time.Time       `json:"end_date,omitempty"`
	Status      types.TaskStatus `json:"status" validate:"required"`
	Score       *int             `json:"score,omitempty" validate:"omitempty,min=1,max=10"`
	Tags        *string          `json:"tags,omitempty"`
}

// TaskFiltersRequest 任务过滤请求结构
type TaskFiltersRequest struct {
	Type      *types.TaskType   `query:"type"`
	Status    *types.TaskStatus `query:"status"`
	ParentID  *uuid.UUID        `query:"parent_id"`
	StartDate *time.Time        `query:"start_date"`
	EndDate   *time.Time        `query:"end_date"`
	Page      int               `query:"page"`
	PageSize  int               `query:"page_size"`
}

// UpdateScoreRequest 更新分数请求结构
type UpdateScoreRequest struct {
	Score int `json:"score" validate:"required,min=1,max=10"`
}

// CreateTask 创建任务
func (ctrl *TaskController) CreateTask(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return middleware.HandleUnauthorized(c, err)
	}

	var req TaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_REQUEST_BODY",
			Message: "请求体格式无效",
		})
	}

	// 参数验证
	if err := validateTaskRequest(req); err != nil {
		return err
	}

	// 调用服务层
	serviceReq := service.TaskRequest{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      req.Status,
		Score:       req.Score,
		Tags:        req.Tags,
	}

	task, err := ctrl.taskService.CreateTask(c.Request().Context(), userID, serviceReq)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusCreated, types.SuccessResponse{
		Success: true,
		Data:    task,
		Message: "任务创建成功",
	})
}

// GetTask 获取单个任务
func (ctrl *TaskController) GetTask(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return middleware.HandleUnauthorized(c, err)
	}

	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_TASK_ID",
			Message: "无效的任务ID",
		})
	}

	task, err := ctrl.taskService.GetTask(c.Request().Context(), userID, taskID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    task,
		Message: "获取任务成功",
	})
}

// UpdateTask 更新任务
func (ctrl *TaskController) UpdateTask(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_TASK_ID",
			Message: "无效的任务ID",
		})
	}

	var req TaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_REQUEST_BODY",
			Message: "请求体格式无效",
		})
	}

	// 参数验证
	if err := validateTaskRequest(req); err != nil {
		return err
	}

	// 调用服务层
	serviceReq := service.TaskRequest{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      req.Status,
		Score:       req.Score,
		Tags:        req.Tags,
	}

	task, err := ctrl.taskService.UpdateTask(c.Request().Context(), userID, taskID, serviceReq)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    task,
		Message: "任务更新成功",
	})
}

// DeleteTask 删除任务
func (ctrl *TaskController) DeleteTask(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_TASK_ID",
			Message: "无效的任务ID",
		})
	}

	err = ctrl.taskService.DeleteTask(c.Request().Context(), userID, taskID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    nil,
		Message: "任务删除成功",
	})
}

// GetTasksByUser 获取用户任务列表
func (ctrl *TaskController) GetTasksByUser(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	// 解析查询参数
	filters := TaskFiltersRequest{
		Page:     1,
		PageSize: 20,
	}

	if err := c.Bind(&filters); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_QUERY_PARAMS",
			Message: "查询参数格式无效",
		})
	}

	// 参数验证
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 || filters.PageSize > 100 {
		filters.PageSize = 20
	}

	// 调用服务层
	serviceFilters := service.TaskFilters{
		Type:      filters.Type,
		Status:    filters.Status,
		ParentID:  filters.ParentID,
		StartDate: filters.StartDate,
		EndDate:   filters.EndDate,
		Page:      filters.Page,
		PageSize:  filters.PageSize,
	}

	result, err := ctrl.taskService.GetTasksByUser(c.Request().Context(), userID, serviceFilters)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, types.PaginationResponse{
		Success:     true,
		Data:        result.Tasks,
		Total:       result.Total,
		CurrentPage: result.CurrentPage,
		PageSize:    result.PageSize,
		TotalPages:  result.TotalPages,
	})
}

// GetTaskChildren 获取子任务
func (ctrl *TaskController) GetTaskChildren(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_TASK_ID",
			Message: "无效的任务ID",
		})
	}

	children, err := ctrl.taskService.GetTaskChildren(c.Request().Context(), userID, taskID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    children,
		Message: "获取子任务成功",
	})
}

// GetTaskTree 获取任务树
func (ctrl *TaskController) GetTaskTree(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_TASK_ID",
			Message: "无效的任务ID",
		})
	}

	tree, err := ctrl.taskService.GetTaskTree(c.Request().Context(), userID, taskID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    tree,
		Message: "获取任务树成功",
	})
}

// GetGlobalView 获取全局视图
func (ctrl *TaskController) GetGlobalView(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	trees, err := ctrl.taskService.GetGlobalView(c.Request().Context(), userID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    trees,
		Message: "获取全局视图成功",
	})
}

// CreateSubTask 创建子任务
func (ctrl *TaskController) CreateSubTask(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	parentIDStr := c.Param("id")
	parentID, err := uuid.Parse(parentIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_PARENT_ID",
			Message: "无效的父任务ID",
		})
	}

	var req TaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_REQUEST_BODY",
			Message: "请求体格式无效",
		})
	}

	// 参数验证
	if err := validateTaskRequest(req); err != nil {
		return err
	}

	// 调用服务层
	serviceReq := service.TaskRequest{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      req.Status,
		Score:       req.Score,
		Tags:        req.Tags,
	}

	task, err := ctrl.taskService.CreateSubTask(c.Request().Context(), userID, parentID, serviceReq)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusCreated, types.SuccessResponse{
		Success: true,
		Data:    task,
		Message: "子任务创建成功",
	})
}

// UpdateTaskScore 更新任务分数
func (ctrl *TaskController) UpdateTaskScore(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_TASK_ID",
			Message: "无效的任务ID",
		})
	}

	var req UpdateScoreRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_REQUEST_BODY",
			Message: "请求体格式无效",
		})
	}

	// 参数验证
	if req.Score < 1 || req.Score > 10 {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_SCORE",
			Message: "评分必须在1-10之间",
		})
	}

	err = ctrl.taskService.UpdateTaskScore(c.Request().Context(), userID, taskID, req.Score)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    nil,
		Message: "任务分数更新成功",
	})
}

// GetContextView 获取上下文视图
func (ctrl *TaskController) GetContextView(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_TASK_ID",
			Message: "无效的任务ID",
		})
	}

	contextView, err := ctrl.taskService.GetContextView(c.Request().Context(), userID, taskID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    contextView,
		Message: "获取上下文视图成功",
	})
}

// GetTasks 获取任务列表
func (c *TaskController) GetTasks(ctx echo.Context) error {
	// 从AuthMiddleware获取已解析的用户ID
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return middleware.HandleUnauthorized(ctx, err)
	}

	// 解析分页参数
	page := 1
	pageSize := 20

	if pageStr := ctx.QueryParam("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if sizeStr := ctx.QueryParam("page_size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 && s <= 100 {
			pageSize = s
		}
	}

	// 构建过滤条件
	filters := service.TaskFilters{
		Page:     page,
		PageSize: pageSize,
	}

	// 解析其他过滤条件
	if taskType := ctx.QueryParam("type"); taskType != "" {
		if tt, err := types.ParseTaskType(taskType); err == nil {
			filters.Type = &tt
		}
	}

	if taskStatus := ctx.QueryParam("status"); taskStatus != "" {
		if ts, err := types.ParseTaskStatus(taskStatus); err == nil {
			filters.Status = &ts
		}
	}

	if parentIDStr := ctx.QueryParam("parent_id"); parentIDStr != "" {
		if parentID, err := uuid.Parse(parentIDStr); err == nil {
			filters.ParentID = &parentID
		}
	}

	if startStr := ctx.QueryParam("start_date"); startStr != "" {
		if start, err := time.Parse("2006-01-02", startStr); err == nil {
			filters.StartDate = &start
		}
	}

	if endStr := ctx.QueryParam("end_date"); endStr != "" {
		if end, err := time.Parse("2006-01-02", endStr); err == nil {
			filters.EndDate = &end
		}
	}

	result, err := c.taskService.GetTasksByUser(ctx.Request().Context(), userID, filters)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Message: "获取任务列表失败",
			Code:    http.StatusInternalServerError,
		})
	}

	return ctx.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    result,
	})
}

// RegisterRoutes 注册任务相关路由
func (ctrl *TaskController) RegisterRoutes(e *echo.Echo, jwtMiddleware echo.MiddlewareFunc) {
	// 所有任务相关路由都需要认证
	tasks := e.Group("/api/tasks")
	tasks.Use(jwtMiddleware)

	// 基本CRUD
	tasks.GET("", ctrl.GetTasksByUser)    // 获取任务列表
	tasks.POST("", ctrl.CreateTask)       // 创建任务
	tasks.GET("/:id", ctrl.GetTask)       // 获取单个任务
	tasks.PUT("/:id", ctrl.UpdateTask)    // 更新任务
	tasks.DELETE("/:id", ctrl.DeleteTask) // 删除任务

	// 任务关系和视图
	tasks.GET("/:id/children", ctrl.GetTaskChildren)    // 获取子任务
	tasks.POST("/:id/sub-task", ctrl.CreateSubTask)     // 创建子任务
	tasks.GET("/:id/full-tree", ctrl.GetTaskTree)       // 获取完整任务树
	tasks.PUT("/:id/score", ctrl.UpdateTaskScore)       // 更新任务分数
	tasks.GET("/:id/context-view", ctrl.GetContextView) // 获取上下文视图
	tasks.GET("/global-view", ctrl.GetGlobalView)       // 获取全局树视图
}

// 辅助函数

// getUserIDFromContext 从context中获取用户ID
// validateTaskRequest 验证任务请求
func validateTaskRequest(req TaskRequest) error {
	if req.Title == "" {
		return echo.NewHTTPError(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "MISSING_TITLE",
			Message: "任务标题不能为空",
		})
	}

	if len(req.Title) > 200 {
		return echo.NewHTTPError(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "TITLE_TOO_LONG",
			Message: "任务标题不能超过200字符",
		})
	}

	if !req.Type.IsValid() {
		return echo.NewHTTPError(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_TASK_TYPE",
			Message: "无效的任务类型",
		})
	}

	if !req.Status.IsValid() {
		return echo.NewHTTPError(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_TASK_STATUS",
			Message: "无效的任务状态",
		})
	}

	if req.Score != nil && (*req.Score < 1 || *req.Score > 10) {
		return echo.NewHTTPError(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_SCORE",
			Message: "评分必须在1-10之间",
		})
	}

	if req.StartDate != nil && req.EndDate != nil && req.StartDate.After(*req.EndDate) {
		return echo.NewHTTPError(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_TIME_RANGE",
			Message: "开始时间不能晚于结束时间",
		})
	}

	return nil
}

// handleServiceError 处理服务层错误
func handleServiceError(c echo.Context, err error) error {
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
