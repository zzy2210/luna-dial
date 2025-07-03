package controller

import (
	"net/http"
	"time"

	"okr-web/internal/middleware"
	"okr-web/internal/service"
	"okr-web/internal/types"

	"github.com/labstack/echo/v4"
)

// StatsController 统计控制器
type StatsController struct {
	statsService service.StatsService
}

// NewStatsController 创建统计控制器
func NewStatsController(statsService service.StatsService) *StatsController {
	return &StatsController{
		statsService: statsService,
	}
}

// GetTaskCompletionStats 获取任务完成度统计
func (c *StatsController) GetTaskCompletionStats(ctx echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return middleware.HandleUnauthorized(ctx, err)
	}

	// 解析时间范围参数
	timeRange, err := c.parseTimeRange(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
			Message: "时间范围参数错误",
			Code:    http.StatusBadRequest,
		})
	}

	stats, err := c.statsService.GetTaskCompletionStats(ctx.Request().Context(), userID, timeRange)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Message: "获取任务完成度统计失败",
			Code:    http.StatusInternalServerError,
		})
	}

	return ctx.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    stats,
	})
}

// GetScoreTrend 获取评分趋势统计
func (c *StatsController) GetScoreTrend(ctx echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return middleware.HandleUnauthorized(ctx, err)
	}

	// 解析时间范围参数
	timeRange, err := c.parseTimeRange(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
			Message: "时间范围参数错误",
			Code:    http.StatusBadRequest,
		})
	}

	stats, err := c.statsService.GetScoreTrend(ctx.Request().Context(), userID, timeRange)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Message: "获取评分趋势统计失败",
			Code:    http.StatusInternalServerError,
		})
	}

	return ctx.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    stats,
	})
}

// GetTimeDistribution 获取时间分布统计
func (c *StatsController) GetTimeDistribution(ctx echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return middleware.HandleUnauthorized(ctx, err)
	}

	// 解析时间范围参数
	timeRange, err := c.parseTimeRange(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
			Message: "时间范围参数错误",
			Code:    http.StatusBadRequest,
		})
	}

	stats, err := c.statsService.GetTimeDistribution(ctx.Request().Context(), userID, timeRange)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Message: "获取时间分布统计失败",
			Code:    http.StatusInternalServerError,
		})
	}

	return ctx.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    stats,
	})
}

// GetUserOverview 获取用户概览统计
func (c *StatsController) GetUserOverview(ctx echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return middleware.HandleUnauthorized(ctx, err)
	}

	overview, err := c.statsService.GetUserOverview(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Message: "获取用户概览统计失败",
			Code:    http.StatusInternalServerError,
		})
	}

	return ctx.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    overview,
	})
}

// parseTimeRange 解析时间范围参数
func (c *StatsController) parseTimeRange(ctx echo.Context) (service.TimeRangeRequest, error) {
	var timeRange service.TimeRangeRequest

	startStr := ctx.QueryParam("start")
	endStr := ctx.QueryParam("end")

	if startStr != "" {
		start, err := time.Parse("2006-01-02", startStr)
		if err != nil {
			return timeRange, err
		}
		timeRange.StartDate = start
	}

	if endStr != "" {
		end, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			return timeRange, err
		}
		timeRange.EndDate = end
	}

	// 如果没有指定时间范围，默认为过去30天
	if startStr == "" && endStr == "" {
		now := time.Now()
		start := now.AddDate(0, 0, -30)
		timeRange.StartDate = start
		timeRange.EndDate = now
	}

	return timeRange, nil
}

// GetScoreTrendByReference 获取分数趋势统计（基于时间参考）
// @Summary 获取分数趋势统计
// @Description 根据时间尺度和时间参考获取分数趋势数据
// @Tags stats
// @Accept json
// @Produce json
// @Param scale query string true "统计尺度" Enums(day,week,month,quarter,year)
// @Param time_ref query string true "时间参考" example("2024-Q4", "2025-07", "2025-07-15")
// @Success 200 {object} types.SuccessResponse{data=service.ScoreTrendResponse}
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security BearerAuth
// @Router /api/stats/score-trend [get]
func (c *StatsController) GetScoreTrendByReference(ctx echo.Context) error {
	// 从上下文获取用户ID
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return middleware.HandleUnauthorized(ctx, err)
	}

	// 解析查询参数
	var req service.ScoreTrendRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_REQUEST",
			Message: "请求参数格式错误: " + err.Error(),
		})
	}

	// 验证请求参数
	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "VALIDATION_FAILED",
			Message: "请求参数验证失败: " + err.Error(),
		})
	}

	// 调用服务获取分数趋势
	response, err := c.statsService.GetScoreTrendByReference(ctx.Request().Context(), userID, req)
	if err != nil {
		if appErr, ok := err.(*types.AppError); ok {
			return ctx.JSON(appErr.Code, types.ErrorResponse{
				Success: false,
				Error:   appErr.Type,
				Message: appErr.Message,
			})
		}
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Success: false,
			Error:   "INTERNAL_SERVER_ERROR",
			Message: "获取分数趋势失败: " + err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    response,
		Message: "获取分数趋势成功",
	})
}
