package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"okr-web/internal/service"
	"okr-web/internal/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// PlanController 计划视图控制器
type PlanController struct {
	taskService service.TaskService
}

// NewPlanController 创建计划视图控制器
func NewPlanController(taskService service.TaskService) *PlanController {
	return &PlanController{
		taskService: taskService,
	}
}

// GetPlanView 获取计划视图
// @Summary 获取计划视图
// @Description 根据时间尺度和时间参考获取指定周期的计划数据
// @Tags plan
// @Accept json
// @Produce json
// @Param scale query string true "时间尺度" Enums(day,week,month,quarter,year)
// @Param time_ref query string true "时间参考" example("2024-Q4", "2025-07", "2025-07-15")
// @Success 200 {object} types.ApiResponse{data=service.PlanResponse}
// @Failure 400 {object} types.ApiResponse
// @Failure 401 {object} types.ApiResponse
// @Failure 500 {object} types.ApiResponse
// @Security BearerAuth
// @Router /api/plan [get]
func (pc *PlanController) GetPlanView(c echo.Context) error {
	// 从上下文获取已验证的用户ID（由全局AuthMiddleware设置）
	userID := c.Get("parsed_user_id").(uuid.UUID)

	// 解析查询参数
	var req service.PlanRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "INVALID_REQUEST",
			Message: "请求参数格式错误: " + err.Error(),
		})
	}

	// 验证请求参数
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Error:   "VALIDATION_FAILED",
			Message: "请求参数验证失败: " + err.Error(),
		})
	}

	// 调用服务获取计划视图
	response, err := pc.taskService.GetPlanView(c.Request().Context(), userID, req)
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
			Message: "获取计划视图失败: " + err.Error(),
		})
	}

	// 调试输出响应体（JSON 格式）
	if respJson, err := json.MarshalIndent(response, "", "  "); err == nil {
		log.Printf("[PlanController] PlanView response JSON: %s", respJson)
	} else {
		log.Printf("[PlanController] PlanView response Marshal error: %v", err)
	}

	return c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Data:    response,
		Message: "获取计划视图成功",
	})
}
