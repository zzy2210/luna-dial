package service

import (
	"fmt"
	"luna_dial/internal/biz"
	"regexp"

	"github.com/labstack/echo/v4"
)

// 查看指定时间段内指定类型的任务
func (s *Service) handleListTasks(c echo.Context) error {
	var req ListTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, NewErrorResponse(400, "Invalid request data"))
	}

	// 获取当前用户ID
	userId, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, "User not found"))
	}

	period := biz.Period{
		Start: req.StartDate,
		End:   req.EndDate,
	}

	// 解析PeriodType
	periodType, err := PeriodTypeFromString(req.PeriodType)
	if err != nil {
		return c.JSON(400, NewErrorResponse(400, fmt.Sprintf("Invalid period type: %s", req.PeriodType)))
	}

	// 调用业务层获取任务列表
	tasks, err := s.taskUsecase.ListTaskByPeriod(c.Request().Context(), biz.ListTaskByPeriodParam{
		UserID:  userId,
		Period:  period,
		GroupBy: periodType,
	})
	if err != nil {
		return c.JSON(500, NewErrorResponse(500, "Failed to get tasks"))
	}

	// 直接返回任务列表
	response := NewSuccessResponse(tasks)
	return c.JSON(200, response)
}

// 创建任务
func (s *Service) handleCreateTask(c echo.Context) error {
	var req CreateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, NewErrorResponse(400, "Invalid request data"))
	}

	userID, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, "User not found"))
	}

	if req.Icon != "" && !IsIcon(req.Icon) {
		return c.JSON(400, NewErrorResponse(400, "Invalid icon format"))
	}

	pType, err := PeriodTypeFromString(req.PeriodType)
	if err != nil {
		return c.JSON(400, NewErrorResponse(400, fmt.Sprintf("Invalid period type: %s", req.PeriodType)))
	}

	priority, err := TaskPriorityFromString(req.Priority)
	if err != nil {
		return c.JSON(400, NewErrorResponse(400, fmt.Sprintf("Invalid priority: %s", req.Priority)))
	}

	task, err := s.taskUsecase.CreateTask(c.Request().Context(), biz.CreateTaskParam{
		Title:  req.Title,
		UserID: userID,
		Type:   pType,
		Period: biz.Period{
			Start: req.StartDate,
			End:   req.EndDate,
		},
		Icon:     req.Icon,
		Tags:     req.Tags,
		Priority: priority,
	})
	if err != nil {
		return c.JSON(500, NewErrorResponse(500, "Failed to create task"))
	}
	return c.JSON(200, NewSuccessResponseWithMessage("create task endpoint", task))
}

// 创建子任务
func (s *Service) handleCreateSubTask(c echo.Context) error {
	var req CreateSubTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, NewErrorResponse(400, "Invalid request data"))
	}

	userID, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, "User not found"))
	}

	if req.Icon != "" && !IsIcon(req.Icon) {
		return c.JSON(400, NewErrorResponse(400, "Invalid icon format"))
	}

	periodType, err := PeriodTypeFromString(req.PeriodType)
	if err != nil {
		return c.JSON(400, NewErrorResponse(400, fmt.Sprintf("Invalid period type: %s", req.PeriodType)))
	}

	priority, err := TaskPriorityFromString(req.Priority)
	if err != nil {
		return c.JSON(400, NewErrorResponse(400, fmt.Sprintf("Invalid priority: %s", req.Priority)))
	}

	subTask, err := s.taskUsecase.CreateSubTask(c.Request().Context(), biz.CreateSubTaskParam{
		ParentID: req.TaskID,
		UserID:   userID,
		Type:     periodType,
		Period: biz.Period{
			Start: req.StartDate,
			End:   req.EndDate,
		},
		Icon:     req.Icon,
		Priority: priority,
	})
	if err != nil {
		return c.JSON(500, NewErrorResponse(500, "Failed to create subtask"))
	}
	return c.JSON(200, NewSuccessResponseWithMessage("create subtask endpoint", subTask))
}

// 更新任务
func (s *Service) handleUpdateTask(c echo.Context) error {
	var req UpdateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, NewErrorResponse(400, "Invalid request data"))
	}

	userID, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, "User not found"))
	}

	if req.Icon != nil && *req.Icon != "" && !IsIcon(*req.Icon) {
		return c.JSON(400, NewErrorResponse(400, "Invalid icon format"))
	}

	// 构建更新参数
	updateParam := biz.UpdateTaskParam{
		TaskID: req.TaskID,
		UserID: userID,
	}

	// 只设置传递了的字段
	if req.Title != nil {
		updateParam.Title = req.Title
	}
	if req.StartDate != nil && req.EndDate != nil {
		updateParam.Period = &biz.Period{
			Start: *req.StartDate,
			End:   *req.EndDate,
		}
	}
	if req.Status != nil {
		status, err := TaskStatusFromString(*req.Status)
		if err != nil {
			return c.JSON(400, NewErrorResponse(400, fmt.Sprintf("Invalid status: %s", *req.Status)))
		}
		updateParam.Status = &status
	}
	if req.Priority != nil {
		priority, err := TaskPriorityFromString(*req.Priority)
		if err != nil {
			return c.JSON(400, NewErrorResponse(400, fmt.Sprintf("Invalid priority: %s", *req.Priority)))
		}
		updateParam.Priority = &priority
	}
	if req.Icon != nil {
		updateParam.Icon = req.Icon
	}
	if req.Tags != nil {
		updateParam.Tags = req.Tags
	}

	task, err := s.taskUsecase.UpdateTask(c.Request().Context(), updateParam)
	if err != nil {
		return c.JSON(500, NewErrorResponse(500, "Failed to update task"))
	}
	return c.JSON(200, NewSuccessResponseWithMessage("update task endpoint", task))
}

// 标记任务完成
func (s *Service) handleCompleteTask(c echo.Context) error {
	var req CompleteTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, NewErrorResponse(400, "Invalid request data"))
	}

	userID, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, "User not found"))
	}

	status := biz.TaskStatusCompleted
	s.taskUsecase.UpdateTask(c.Request().Context(), biz.UpdateTaskParam{
		TaskID: req.TaskID,
		UserID: userID,
		Status: &status,
	})
	return c.JSON(200, NewSuccessResponseWithMessage("complete task endpoint", nil))
}

// 更新任务分数
func (s *Service) handleUpdateTaskScore(c echo.Context) error {
	var req UpdateTaskScoreRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, NewErrorResponse(400, "Invalid request data"))
	}

	userID, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, "User not found"))
	}

	if req.Score < 0 {
		return c.JSON(400, NewErrorResponse(400, "Score must be non-negative"))
	}

	_, err = s.taskUsecase.SetTaskScore(c.Request().Context(), biz.SetTaskScoreParam{
		TaskID: req.TaskID,
		UserID: userID,
		Score:  req.Score,
	})
	if err != nil {
		return c.JSON(500, NewErrorResponse(500, "Failed to update task score"))
	}
	return c.JSON(200, NewSuccessResponseWithMessage("update task score endpoint", nil))
}

// 删除任务
func (s *Service) handleDeleteTask(c echo.Context) error {
	var req DeleteTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, NewErrorResponse(400, "Invalid request data"))
	}

	userID, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, "User not found"))
	}

	err = s.taskUsecase.DeleteTask(c.Request().Context(), biz.DeleteTaskParam{
		TaskID: req.TaskID,
		UserID: userID,
	})
	if err != nil {
		return c.JSON(500, NewErrorResponse(500, "Failed to delete task"))
	}
	return c.JSON(200, NewSuccessResponseWithMessage("delete task endpoint", nil))
}

// 检查 string 是否是 icon （emoji）
func IsIcon(s string) bool {
	if s == "" {
		return false
	}

	runes := []rune(s)
	if len(runes) == 0 || len(runes) > 4 {
		return false // emoji通常不会超过4个rune
	}

	// 使用正则表达式匹配emoji的Unicode范围
	// 这个正则表达式涵盖了大部分常见的emoji Unicode范围
	emojiRegex := regexp.MustCompile(`^[\x{1F600}-\x{1F64F}\x{1F300}-\x{1F5FF}\x{1F680}-\x{1F6FF}\x{1F1E0}-\x{1F1FF}\x{2600}-\x{26FF}\x{2700}-\x{27BF}\x{2B00}-\x{2BFF}\x{1F900}-\x{1F9FF}\x{1F018}-\x{1F270}\x{238C}-\x{2454}\x{20D0}-\x{20FF}\x{FE0F}\x{200D}]+$`)

	// 先用正则表达式快速检查
	if emojiRegex.MatchString(s) {
		return true
	}

	// 检查每个字符是否都是emoji相关的Unicode字符
	for _, r := range runes {
		// 检查是否是emoji相关的Unicode范围
		if !isEmojiRune(r) {
			return false
		}
	}

	return true
}

// 辅助函数：检查单个rune是否是emoji字符
func isEmojiRune(r rune) bool {
	return (r >= 0x1F600 && r <= 0x1F64F) || // 表情符号
		(r >= 0x1F300 && r <= 0x1F5FF) || // 符号和象形文字
		(r >= 0x1F680 && r <= 0x1F6FF) || // 交通和地图符号
		(r >= 0x1F1E0 && r <= 0x1F1FF) || // 区域指示符号
		(r >= 0x2600 && r <= 0x26FF) || // 杂项符号
		(r >= 0x2700 && r <= 0x27BF) || // 装饰符号
		(r >= 0x1F900 && r <= 0x1F9FF) || // 补充符号和象形文字
		(r >= 0x1F018 && r <= 0x1F270) || // 其他符号
		(r >= 0x238C && r <= 0x2454) || // 技术符号
		(r >= 0x20D0 && r <= 0x20FF) || // 组合变音符号
		(r >= 0x2B00 && r <= 0x2BFF) || // 杂项符号和箭头 (包含⭐)
		r == 0xFE0F || // 变异选择器-16 (emoji presentation)
		r == 0x200D // 零宽连接符 (用于组合emoji)
}

func BoolPtr(b bool) *bool {
	return &b
}
