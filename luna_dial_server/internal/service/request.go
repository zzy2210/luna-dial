package service

import (
	"fmt"
	"luna_dial/internal/biz"
	"time"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ListTaskRequest struct {
	PeriodType string    `json:"period_type" validate:"required,oneof=daily weekly monthly yearly"`
	StartDate  time.Time `json:"start_date" validate:"required"`
	EndDate    time.Time `json:"end_date" validate:"required"`
}

type CreateTaskRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	EndDate     time.Time `json:"end_date" validate:"required"`
	Priority    string    `json:"priority" validate:"required,oneof=low medium high"`
	Icon        string    `json:"icon"`
	Tags        []string  `json:"tags"`
	Status      string    `json:"status" validate:"required,oneof=pending in-progress completed"`
}

func PeriodTypeFromString(s string) (biz.PeriodType, error) {
	switch s {
	case "day":
		return biz.PeriodDay, nil
	case "week":
		return biz.PeriodWeek, nil
	case "month":
		return biz.PeriodMonth, nil
	case "quarter":
		return biz.PeriodQuarter, nil
	case "year":
		return biz.PeriodYear, nil
	default:
		return 0, fmt.Errorf("unknown period type: %s", s)
	}
}
