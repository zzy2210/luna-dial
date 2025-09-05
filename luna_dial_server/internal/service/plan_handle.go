package service

import (
	"luna_dial/internal/biz"

	"github.com/labstack/echo/v4"
)

func (s *Service) handleListPlans(c echo.Context) error {
    var req ListPlansRequest
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
    groupBy, err := PeriodTypeFromString(req.PeriodType)
    if err != nil {
        return c.JSON(400, NewErrorResponse(400, "Invalid period type"))
    }
    plan, err := s.planUsecase.GetPlanByPeriod(c.Request().Context(), biz.GetPlanByPeriodParam{
        UserID: userID,
        Period: biz.Period{
            Start: req.StartDate,
            End:   req.EndDate,
        },
        GroupBy: groupBy,
    })
    if err != nil {
        return c.JSON(500, NewErrorResponse(500, "Failed to get plan"))
    }
    return c.JSON(200, NewSuccessResponse(plan))
}
