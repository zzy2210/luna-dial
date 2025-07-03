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

// JournalController 日志控制器
type JournalController struct {
	journalService service.JournalService
}

// NewJournalController 创建日志控制器
func NewJournalController(journalService service.JournalService) *JournalController {
	return &JournalController{
		journalService: journalService,
	}
}

// CreateJournal 创建日志
func (c *JournalController) CreateJournal(ctx echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return middleware.HandleUnauthorized(ctx, err)
	}

	var req service.JournalRequest
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

	journal, err := c.journalService.CreateJournal(ctx.Request().Context(), userID, req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Message: "创建日志失败",
			Code:    http.StatusInternalServerError,
		})
	}

	return ctx.JSON(http.StatusCreated, types.SuccessResponse{
		Success: true,
		Data:    journal,
	})
}

// GetJournal 获取日志详情
func (c *JournalController) GetJournal(ctx echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return middleware.HandleUnauthorized(ctx, err)
	}

	journalIDStr := ctx.Param("id")
	journalID, err := uuid.Parse(journalIDStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
			Message: "无效的日志ID",
			Code:    http.StatusBadRequest,
		})
	}

	journal, err := c.journalService.GetJournal(ctx.Request().Context(), userID, journalID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, types.ErrorResponse{
			Message: "日志不存在",
			Code:    http.StatusNotFound,
		})
	}

	return ctx.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    journal,
	})
}

// UpdateJournal 更新日志
func (c *JournalController) UpdateJournal(ctx echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return middleware.HandleUnauthorized(ctx, err)
	}

	journalIDStr := ctx.Param("id")
	journalID, err := uuid.Parse(journalIDStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
			Message: "无效的日志ID",
			Code:    http.StatusBadRequest,
		})
	}

	var req service.JournalRequest
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

	journal, err := c.journalService.UpdateJournal(ctx.Request().Context(), userID, journalID, req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Message: "更新日志失败",
			Code:    http.StatusInternalServerError,
		})
	}

	return ctx.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    journal,
	})
}

// DeleteJournal 删除日志
func (c *JournalController) DeleteJournal(ctx echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return middleware.HandleUnauthorized(ctx, err)
	}

	journalIDStr := ctx.Param("id")
	journalID, err := uuid.Parse(journalIDStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
			Message: "无效的日志ID",
			Code:    http.StatusBadRequest,
		})
	}

	err = c.journalService.DeleteJournal(ctx.Request().Context(), userID, journalID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Message: "删除日志失败",
			Code:    http.StatusInternalServerError,
		})
	}

	return ctx.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "日志删除成功",
	})
}

// GetJournalsByTime 按时间范围获取日志
func (c *JournalController) GetJournalsByTime(ctx echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return middleware.HandleUnauthorized(ctx, err)
	}

	// 解析时间范围参数
	startStr := ctx.QueryParam("start")
	endStr := ctx.QueryParam("end")

	var timeRange service.TimeRangeRequest
	if startStr != "" {
		start, err := time.Parse("2006-01-02", startStr)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
				Message: "开始时间格式错误",
				Code:    http.StatusBadRequest,
			})
		}
		timeRange.StartDate = start
	}

	if endStr != "" {
		end, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
				Message: "结束时间格式错误",
				Code:    http.StatusBadRequest,
			})
		}
		timeRange.EndDate = end
	}

	journals, err := c.journalService.GetJournalsByTime(ctx.Request().Context(), userID, timeRange)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Message: "获取日志失败",
			Code:    http.StatusInternalServerError,
		})
	}

	return ctx.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    journals,
	})
}

// GetJournalsByUser 获取用户的日志列表（分页）
func (c *JournalController) GetJournalsByUser(ctx echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return middleware.HandleUnauthorized(ctx, err)
	}

	// 解析分页参数
	page := 1
	limit := 20

	if pageStr := ctx.QueryParam("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := ctx.QueryParam("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// 构建过滤条件
	filters := service.JournalFilters{
		Page:     page,
		PageSize: limit,
	}

	// 解析其他过滤条件
	if entryType := ctx.QueryParam("entry_type"); entryType != "" {
		if et, err := types.ParseEntryType(entryType); err == nil {
			filters.EntryType = &et
		}
	}

	if startStr := ctx.QueryParam("start"); startStr != "" {
		if start, err := time.Parse("2006-01-02", startStr); err == nil {
			filters.StartDate = &start
		}
	}

	if endStr := ctx.QueryParam("end"); endStr != "" {
		if end, err := time.Parse("2006-01-02", endStr); err == nil {
			filters.EndDate = &end
		}
	}

	result, err := c.journalService.GetJournalsByUser(ctx.Request().Context(), userID, filters)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Message: "获取日志列表失败",
			Code:    http.StatusInternalServerError,
		})
	}

	return ctx.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    result,
	})
}

// LinkJournalToTasks 关联日志到任务
func (c *JournalController) LinkJournalToTasks(ctx echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return middleware.HandleUnauthorized(ctx, err)
	}

	journalIDStr := ctx.Param("id")
	journalID, err := uuid.Parse(journalIDStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
			Message: "无效的日志ID",
			Code:    http.StatusBadRequest,
		})
	}

	var req struct {
		TaskIDs []uuid.UUID `json:"task_ids" validate:"required"`
	}

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

	err = c.journalService.LinkJournalToTasks(ctx.Request().Context(), userID, journalID, req.TaskIDs)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Message: "关联任务失败",
			Code:    http.StatusInternalServerError,
		})
	}

	return ctx.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "任务关联成功",
	})
}
