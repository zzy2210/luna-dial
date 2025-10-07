package service

import (
	"fmt"
	"luna_dial/internal/biz"
	"time"

	"github.com/labstack/echo/v4"
)

// 根据时间段与时间类型获取 无分页
func (s *Service) handleListJournalsByPeriod(c echo.Context) error {
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

	userID, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, "User not found"))
	}

	periodTypeEnum, err := PeriodTypeFromString(periodType)
	if err != nil {
		return c.JSON(400, NewErrorResponse(400, "Invalid period type"))
	}

	journalList, err := s.journalUsecase.ListJournalByPeriod(c.Request().Context(), biz.ListJournalByPeriodParam{
		UserID:  userID,
		Period:  biz.Period{Start: startDate, End: endDate},
		GroupBy: periodTypeEnum,
	})
	if err != nil {
		c.Logger().Error("Failed to get journals:", err)
		return c.JSON(500, NewErrorResponse(500, "Failed to get journals: " + err.Error()))
	}

	return c.JSON(200, NewSuccessResponse(journalList))
}

func (s *Service) handleCreateJournal(c echo.Context) error {
    var req CreateJournalRequest
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
	if req.Title == "" || req.Content == "" || req.JournalType == "" || req.StartDate.IsZero() || req.EndDate.IsZero() {
		return c.JSON(400, NewErrorResponse(400, "Title, content, journal type and time period are required"))
	}
	journalType, err := PeriodTypeFromString(req.JournalType)
	if err != nil {
		return c.JSON(400, NewErrorResponse(400, "Invalid journal type"))
	}

	journal, err := s.journalUsecase.CreateJournal(c.Request().Context(), biz.CreateJournalParam{
		UserID:      userID,
		Title:       req.Title,
		Content:     req.Content,
		JournalType: journalType,
		TimePeriod: biz.Period{
			Start: req.StartDate,
			End:   req.EndDate,
		},
		Icon: req.Icon,
	})
	if err != nil {
		return c.JSON(500, NewErrorResponse(500, "Failed to create journal"))
	}
	return c.JSON(201, NewSuccessResponse(journal))
}

// 更新
func (s *Service) handleUpdateJournal(c echo.Context) error {
    // 路径参数作为唯一 ID 来源
    journalID := c.Param("journal_id")
    if journalID == "" {
        return c.JSON(400, NewErrorResponse(400, "Journal ID is required"))
    }

    var req UpdateJournalRequest
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

    if req.Title == nil && req.Content == nil && req.JournalType == nil && req.Icon == nil {
        return c.JSON(400, NewErrorResponse(400, "At least one field must be provided for update"))
    }
    var journalType *biz.PeriodType
    if req.JournalType != nil && *req.JournalType != "" {
        if jt, err := PeriodTypeFromString(*req.JournalType); err != nil {
            return c.JSON(400, NewErrorResponse(400, "Invalid journal type"))
        } else {
            journalType = &jt
        }
    }

    journal, err := s.journalUsecase.UpdateJournal(c.Request().Context(), biz.UpdateJournalParam{
        JournalID:   journalID,
        UserID:      userID,
        Title:       req.Title,
        Content:     req.Content,
        JournalType: journalType,
        Icon:        req.Icon,
    })
    if err != nil {
        return c.JSON(500, NewErrorResponse(500, "Failed to update journal"))
    }
    return c.JSON(200, NewSuccessResponse(journal))
}

// 删除
func (s *Service) handleDeleteJournal(c echo.Context) error {
	journalID := c.Param("journal_id")
	if journalID == "" {
		return c.JSON(400, NewErrorResponse(400, "Journal ID is required"))
	}

	userID, _, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, "User not found"))
	}

	if err := s.journalUsecase.DeleteJournal(c.Request().Context(), biz.DeleteJournalParam{
		JournalID: journalID,
		UserID:    userID,
	}); err != nil {
		return c.JSON(500, NewErrorResponse(500, "Failed to delete journal"))
	}
	return c.NoContent(204)
}

// 分页查询日志列表（支持过滤）
func (s *Service) handleListJournalsWithPagination(c echo.Context) error {
    var req ListJournalsWithPaginationRequest
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

	// 转换日志类型过滤条件
	var journalType *int
	if req.JournalType != nil {
		pType, err := PeriodTypeFromString(*req.JournalType)
		if err != nil {
			return c.JSON(400, NewErrorResponse(400, fmt.Sprintf("Invalid journal type: %s", *req.JournalType)))
		}
		intType := int(pType)
		journalType = &intType
	}

	// 调用业务层
	journals, total, err := s.journalUsecase.ListJournalsWithPagination(c.Request().Context(), biz.ListJournalsWithPaginationParam{
		UserID:      userID,
		Page:        req.Page,
		PageSize:    req.PageSize,
		JournalType: journalType,
		PeriodStart: req.StartDate,
		PeriodEnd:   req.EndDate,
	})
	if err != nil {
		return c.JSON(500, NewErrorResponse(500, "Failed to get journals"))
	}

	// 返回分页响应
	return c.JSON(200, NewPaginatedResponse(journals, req.Page, req.PageSize, total))
}
