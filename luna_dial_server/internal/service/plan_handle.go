package service

import (
	"time"

	"github.com/labstack/echo/v4"
)

func (s *Service) handleListPlans(c echo.Context) error {
	return c.JSON(200, &Response{
		Code:      200,
		Message:   "list plans endpoint",
		Success:   true,
		Timestamp: time.Now().Unix(),
	})
}

func (s *Service) handleCreatePlan(c echo.Context) error {
	return c.JSON(200, &Response{
		Code:      200,
		Message:   "create plan endpoint",
		Success:   true,
		Timestamp: time.Now().Unix(),
	})
}
