package service

import (
	"fmt"
	"luna_dial/internal/biz"
	"time"

	"github.com/labstack/echo/v4"
)

func (s *Service) handleListPlans(c echo.Context) error {
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

    groupBy, err := PeriodTypeFromString(periodType)
    if err != nil {
        return c.JSON(400, NewErrorResponse(400, "Invalid period type"))
    }

    plan, err := s.planUsecase.GetPlanByPeriod(c.Request().Context(), biz.GetPlanByPeriodParam{
        UserID: userID,
        Period: biz.Period{
            Start: startDate,
            End:   endDate,
        },
        GroupBy: groupBy,
    })
    if err != nil {
        // 打印详细错误信息到日志
        c.Logger().Error("Failed to get plan:", err)

        // 根据不同错误类型返回更明确的错误信息
        switch err {
        case biz.ErrInvalidInput:
            // 检查是否是时间区间问题
            if !startDate.Before(endDate) {
                return c.JSON(400, NewErrorResponse(400,
                    fmt.Sprintf("Invalid time period: end_date must be after start_date. Got start=%s, end=%s",
                        startDateStr, endDateStr)))
            }
            return c.JSON(400, NewErrorResponse(400, "Invalid input parameters"))
        case biz.ErrPlanPeriodInvalid:
            return c.JSON(400, NewErrorResponse(400,
                fmt.Sprintf("Invalid period: start_date must be before end_date. Got start=%s, end=%s",
                    startDateStr, endDateStr)))
        default:
            return c.JSON(500, NewErrorResponse(500, "Failed to get plan: " + err.Error()))
        }
    }
    return c.JSON(200, NewSuccessResponse(plan))
}

func (s *Service) handleGetPlanStats(c echo.Context) error {
    // 从查询参数获取值
    groupBy := c.QueryParam("group_by")
    startDateStr := c.QueryParam("start_date")
    endDateStr := c.QueryParam("end_date")

    // 手动验证必填字段
    if groupBy == "" {
        return c.JSON(400, NewErrorResponse(400, "field group_by is required"))
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

    groupByType, err := PeriodTypeFromString(groupBy)
    if err != nil {
        return c.JSON(400, NewErrorResponse(400, "Invalid group_by type"))
    }

    stats, err := s.planUsecase.GetPlanStats(c.Request().Context(), biz.GetPlanStatsParam{
        UserID: userID,
        Period: biz.Period{
            Start: startDate,
            End:   endDate,
        },
        GroupBy: groupByType,
    })
    if err != nil {
        c.Logger().Error("Failed to get plan stats:", err)

        switch err {
        case biz.ErrInvalidInput:
            if !startDate.Before(endDate) {
                return c.JSON(400, NewErrorResponse(400,
                    fmt.Sprintf("Invalid time period: end_date must be after start_date. Got start=%s, end=%s",
                        startDateStr, endDateStr)))
            }
            return c.JSON(400, NewErrorResponse(400, "Invalid input parameters"))
        case biz.ErrPlanPeriodInvalid:
            return c.JSON(400, NewErrorResponse(400,
                fmt.Sprintf("Invalid period: start_date must be before end_date. Got start=%s, end=%s",
                    startDateStr, endDateStr)))
        default:
            return c.JSON(500, NewErrorResponse(500, "Failed to get plan stats: "+err.Error()))
        }
    }

    return c.JSON(200, NewSuccessResponse(stats))
}
