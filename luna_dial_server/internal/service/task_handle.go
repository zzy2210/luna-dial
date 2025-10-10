package service

import (
	"fmt"
	"luna_dial/internal/biz"
	"regexp"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	// 用于分割逗号分隔的状态参数的正则表达式
	statusSeparatorRegex = regexp.MustCompile(`\s*,\s*`)
)

// 查看指定时间段内指定类型的任务
func (s *Service) handleListTasks(c echo.Context) error {
    // 手动从查询参数获取值
    periodType := c.QueryParam("period_type")
    startDateStr := c.QueryParam("start_date")
    endDateStr := c.QueryParam("end_date")

    // 手动验证必填字段
    if periodType == "" {
        return c.JSON(400, NewErrorResponse(400, "field period_type is required"))
    }
    if startDateStr == "" {
        return c.JSON(400, NewErrorResponse(400, "field start_date is required"))
    }
    if endDateStr == "" {
        return c.JSON(400, NewErrorResponse(400, "field end_date is required"))
    }

    // 解析时间
    startDate, err := time.Parse("2006-01-02", startDateStr)
    if err != nil {
        return c.JSON(400, NewErrorResponse(400, "Invalid start_date format, expected YYYY-MM-DD"))
    }
    endDate, err := time.Parse("2006-01-02", endDateStr)
    if err != nil {
        return c.JSON(400, NewErrorResponse(400, "Invalid end_date format, expected YYYY-MM-DD"))
    }

	// 获取当前用户ID
	userId, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, "User not found"))
	}

	period := biz.Period{
		Start: startDate,
		End:   endDate,
	}

	// 解析PeriodType
	periodTypeEnum, err := PeriodTypeFromString(periodType)
	if err != nil {
		return c.JSON(400, NewErrorResponse(400, fmt.Sprintf("Invalid period type: %s", periodType)))
	}

	// 调用业务层获取任务列表
	tasks, err := s.taskUsecase.ListTaskByPeriod(c.Request().Context(), biz.ListTaskByPeriodParam{
		UserID:  userId,
		Period:  period,
		GroupBy: periodTypeEnum,
	})
	if err != nil {
		c.Logger().Error("Failed to get tasks:", err)
		return c.JSON(500, NewErrorResponse(500, "Failed to get tasks: " + err.Error()))
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
    if err := c.Validate(&req); err != nil {
        return c.JSON(400, NewErrorResponse(400, err.Error()))
    }

	userID, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, "User not found"))
	}

	if req.Icon != "" && !IsIcon(req.Icon) {
		return c.JSON(400, NewErrorResponse(400, "Invalid icon format"))
	}

	// 解析日期字符串
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return c.JSON(400, NewErrorResponse(400, "Invalid start_date format, expected YYYY-MM-DD"))
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return c.JSON(400, NewErrorResponse(400, "Invalid end_date format, expected YYYY-MM-DD"))
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
			Start: startDate,
			End:   endDate,
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
    // 路径参数作为父任务ID的单一可信来源
    parentID := c.Param("task_id")
    if parentID == "" {
        return c.JSON(400, NewErrorResponse(400, "Task ID is required"))
    }

    var req CreateSubTaskRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(400, NewErrorResponse(400, "Invalid request data"))
    }
    if err := c.Validate(&req); err != nil {
        return c.JSON(400, NewErrorResponse(400, err.Error()))
    }

    userID, _, err := GetUserFromContext(c)
    if err != nil {
        return c.JSON(401, NewErrorResponse(401, "User not found"))
    }

    if req.Icon != "" && !IsIcon(req.Icon) {
        return c.JSON(400, NewErrorResponse(400, "Invalid icon format"))
    }

    // 解析日期字符串
    startDate, err := time.Parse("2006-01-02", req.StartDate)
    if err != nil {
        return c.JSON(400, NewErrorResponse(400, "Invalid start_date format, expected YYYY-MM-DD"))
    }
    endDate, err := time.Parse("2006-01-02", req.EndDate)
    if err != nil {
        return c.JSON(400, NewErrorResponse(400, "Invalid end_date format, expected YYYY-MM-DD"))
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
        ParentID: parentID,
        UserID:   userID,
        Title:    req.Title,
        Type:     periodType,
        Period: biz.Period{
            Start: startDate,
            End:   endDate,
        },
        Icon:     req.Icon,
        Priority: priority,
        Tags:     req.Tags,
    })
    if err != nil {
        return c.JSON(500, NewErrorResponse(500, fmt.Sprintf("Failed to create subtask: %v", err)))
    }
    return c.JSON(200, NewSuccessResponseWithMessage("create subtask endpoint", subTask))
}

// 更新任务
func (s *Service) handleUpdateTask(c echo.Context) error {
    // 路径参数为唯一任务ID来源
    taskID := c.Param("task_id")
    if taskID == "" {
        return c.JSON(400, NewErrorResponse(400, "Task ID is required"))
    }

    var req UpdateTaskRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(400, NewErrorResponse(400, "Invalid request data"))
    }
    if err := c.Validate(&req); err != nil {
        return c.JSON(400, NewErrorResponse(400, err.Error()))
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
        TaskID: taskID,
        UserID: userID,
    }

	// 只设置传递了的字段
	if req.Title != nil {
		updateParam.Title = req.Title
	}
	if req.StartDate != nil && req.EndDate != nil {
		// 解析日期字符串
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return c.JSON(400, NewErrorResponse(400, "Invalid start_date format, expected YYYY-MM-DD"))
		}
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return c.JSON(400, NewErrorResponse(400, "Invalid end_date format, expected YYYY-MM-DD"))
		}
		updateParam.Period = &biz.Period{
			Start: startDate,
			End:   endDate,
		}
	}
	if req.Status != "" {
		status, err := TaskStatusFromString(req.Status)
		if err != nil {
			return c.JSON(400, NewErrorResponse(400, fmt.Sprintf("Invalid status: %s", req.Status)))
		}
		updateParam.Status = &status
	}
	if req.Priority != "" {
		priority, err := TaskPriorityFromString(req.Priority)
		if err != nil {
			return c.JSON(400, NewErrorResponse(400, fmt.Sprintf("Invalid priority: %s", req.Priority)))
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
    // 路径参数为唯一任务ID来源
    taskID := c.Param("task_id")
    if taskID == "" {
        return c.JSON(400, NewErrorResponse(400, "Task ID is required"))
    }

    // body 可以为空
    _ = c.Request()

    userID, _, err := GetUserFromContext(c)
    if err != nil {
        return c.JSON(401, NewErrorResponse(401, "User not found"))
    }

    status := biz.TaskStatusCompleted
    _, err = s.taskUsecase.UpdateTask(c.Request().Context(), biz.UpdateTaskParam{
        TaskID: taskID,
        UserID: userID,
        Status: &status,
    })
    if err != nil {
        return c.JSON(500, NewErrorResponse(500, "Failed to complete task"))
    }
    return c.JSON(200, NewSuccessResponseWithMessage("complete task endpoint", nil))
}

// 更新任务分数
func (s *Service) handleUpdateTaskScore(c echo.Context) error {
    // 路径参数为唯一任务ID来源
    taskID := c.Param("task_id")
    if taskID == "" {
        return c.JSON(400, NewErrorResponse(400, "Task ID is required"))
    }

    var req UpdateTaskScoreRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(400, NewErrorResponse(400, "Invalid request data"))
    }
    if err := c.Validate(&req); err != nil {
        return c.JSON(400, NewErrorResponse(400, err.Error()))
    }

    userID, _, err := GetUserFromContext(c)
    if err != nil {
        return c.JSON(401, NewErrorResponse(401, "User not found"))
    }

    if req.Score < 0 {
        return c.JSON(400, NewErrorResponse(400, "Score must be non-negative"))
    }

    _, err = s.taskUsecase.SetTaskScore(c.Request().Context(), biz.SetTaskScoreParam{
        TaskID: taskID,
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
    // 路径参数为唯一任务ID来源
    taskID := c.Param("task_id")
    if taskID == "" {
        return c.JSON(400, NewErrorResponse(400, "Task ID is required"))
    }

    userID, _, err := GetUserFromContext(c)
    if err != nil {
        return c.JSON(401, NewErrorResponse(401, "User not found"))
    }

    err = s.taskUsecase.DeleteTask(c.Request().Context(), biz.DeleteTaskParam{
        TaskID: taskID,
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

// 分页获取根任务列表
func (s *Service) handleListRootTasks(c echo.Context) error {
    var req ListRootTasksRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(400, NewErrorResponse(400, "Invalid request data"))
    }
    if err := c.Validate(&req); err != nil {
        return c.JSON(400, NewErrorResponse(400, err.Error()))
    }

	// 获取当前用户ID
	userID, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, "User not found"))
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 转换状态过滤条件
	var statusFilters []biz.TaskStatus
	for _, statusStr := range req.Status {
		status, err := TaskStatusFromString(statusStr)
		if err != nil {
			return c.JSON(400, NewErrorResponse(400, fmt.Sprintf("Invalid status: %s", statusStr)))
		}
		statusFilters = append(statusFilters, status)
	}

	// 注意：当前业务层不支持优先级和任务类型过滤，这些过滤条件将被忽略
	// TODO: 未来版本可以扩展业务层支持更多过滤条件

	// 调用业务层
	tasks, total, err := s.taskUsecase.ListRootTasks(c.Request().Context(), biz.ListRootTasksParam{
		UserID:        userID,
		Page:          req.Page,
		PageSize:      req.PageSize,
		IncludeStatus: statusFilters,
	})
	if err != nil {
		return c.JSON(500, NewErrorResponse(500, "Failed to get root tasks"))
	}

	// 返回分页响应
	return c.JSON(200, NewPaginatedResponse(tasks, req.Page, req.PageSize, total))
}

// 获取全局任务树（分页）
func (s *Service) handleListGlobalTaskTree(c echo.Context) error {
	var req ListGlobalTaskTreeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, NewErrorResponse(400, "Invalid request data"))
	}

	// 获取当前用户ID
	userID, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, "User not found"))
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 转换状态过滤条件
	var statusFilters []biz.TaskStatus
	for _, statusStr := range req.Status {
		status, err := TaskStatusFromString(statusStr)
		if err != nil {
			return c.JSON(400, NewErrorResponse(400, fmt.Sprintf("Invalid status: %s", statusStr)))
		}
		statusFilters = append(statusFilters, status)
	}

    // 调用业务层
    taskTrees, total, err := s.taskUsecase.ListGlobalTaskTree(c.Request().Context(), biz.ListGlobalTaskTreeParam{
        UserID:        userID,
        Page:          req.Page,
        PageSize:      req.PageSize,
        IncludeStatus: statusFilters,
    })
    if err != nil {
        return c.JSON(500, NewErrorResponse(500, "Failed to get global task tree"))
    }

	// 返回分页响应
	return c.JSON(200, NewPaginatedResponse(taskTrees, req.Page, req.PageSize, total))
}

// 获取指定任务的完整任务树
func (s *Service) handleGetTaskTree(c echo.Context) error {
	taskID := c.Param("task_id")
	if taskID == "" {
		return c.JSON(400, NewErrorResponse(400, "Task ID is required"))
	}

	// 获取当前用户ID
	userID, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, "User not found"))
	}

	// 可选的状态过滤
	statusFilters := []biz.TaskStatus{}
	if statusParam := c.QueryParam("status"); statusParam != "" {
		// 支持多个状态过滤，用逗号分隔
		for _, statusStr := range statusSeparatorRegex.Split(statusParam, -1) {
			status, err := TaskStatusFromString(statusStr)
			if err != nil {
				return c.JSON(400, NewErrorResponse(400, fmt.Sprintf("Invalid status: %s", statusStr)))
			}
			statusFilters = append(statusFilters, status)
		}
	}

	// 调用业务层
	taskTree, err := s.taskUsecase.GetCompleteTaskTree(c.Request().Context(), biz.GetCompleteTaskTreeParam{
		UserID:        userID,
		TaskID:        taskID,
		IncludeStatus: statusFilters,
	})
	if err != nil {
		return c.JSON(500, NewErrorResponse(500, "Failed to get task tree"))
	}

	// 返回树形结构响应
	return c.JSON(200, NewSuccessResponse(taskTree))
}

// 获取任务的父任务链
func (s *Service) handleGetTaskParents(c echo.Context) error {
	taskID := c.Param("task_id")
	if taskID == "" {
		return c.JSON(400, NewErrorResponse(400, "Task ID is required"))
	}

	// 获取当前用户ID
	userID, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, "User not found"))
	}

	// 调用业务层
	parentChain, err := s.taskUsecase.GetTaskParentChain(c.Request().Context(), biz.GetTaskParentChainParam{
		UserID: userID,
		TaskID: taskID,
	})
	if err != nil {
		return c.JSON(500, NewErrorResponse(500, "Failed to get task parent chain"))
	}

	// 返回父任务链
	return c.JSON(200, NewSuccessResponse(parentChain))
}

// 移动任务
func (s *Service) handleMoveTask(c echo.Context) error {
	var req MoveTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, NewErrorResponse(400, "Invalid request data"))
	}

	// 获取当前用户ID
	_, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, "User not found"))
	}

	// TODO: 实现任务移动逻辑（当前业务层尚未实现MoveTask方法）
	// 这里先返回一个提示信息
	return c.JSON(501, NewErrorResponse(501, "Task move functionality is not yet implemented"))
}

// 使用优化的任务创建方法
func (s *Service) handleCreateTaskWithOptimization(c echo.Context) error {
    var req CreateTaskRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(400, NewErrorResponse(400, "Invalid request data"))
    }
    if err := c.Validate(&req); err != nil {
        return c.JSON(400, NewErrorResponse(400, err.Error()))
    }

	userID, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, "User not found"))
	}

	if req.Icon != "" && !IsIcon(req.Icon) {
		return c.JSON(400, NewErrorResponse(400, "Invalid icon format"))
	}

	// 解析日期字符串
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return c.JSON(400, NewErrorResponse(400, "Invalid start_date format, expected YYYY-MM-DD"))
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return c.JSON(400, NewErrorResponse(400, "Invalid end_date format, expected YYYY-MM-DD"))
	}

	pType, err := PeriodTypeFromString(req.PeriodType)
	if err != nil {
		return c.JSON(400, NewErrorResponse(400, fmt.Sprintf("Invalid period type: %s", req.PeriodType)))
	}

	priority, err := TaskPriorityFromString(req.Priority)
	if err != nil {
		return c.JSON(400, NewErrorResponse(400, fmt.Sprintf("Invalid priority: %s", req.Priority)))
	}

	// 使用优化的创建方法
	task, err := s.taskUsecase.CreateTaskWithTreeOptimization(c.Request().Context(), biz.CreateTaskParam{
		Title:  req.Title,
		UserID: userID,
		Type:   pType,
		Period: biz.Period{
			Start: startDate,
			End:   endDate,
		},
		Icon:     req.Icon,
		Tags:     req.Tags,
		Priority: priority,
	})
	if err != nil {
		return c.JSON(500, NewErrorResponse(500, "Failed to create task"))
	}
	return c.JSON(200, NewSuccessResponseWithMessage("Task created with tree optimization", task))
}
