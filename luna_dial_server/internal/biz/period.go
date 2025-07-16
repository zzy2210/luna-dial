package biz

import (
	"fmt"
	"time"
)

type PeriodType int

const (
	PeriodDay PeriodType = iota
	PeriodWeek
	PeriodMonth
	PeriodQuarter
	PeriodYear
)

type Period struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func (p Period) IsValid() bool {
	return !p.Start.IsZero() && !p.End.IsZero() && p.Start.Before(p.End)
}

// 检查当前 Period 是否符合指定的周期类型要求（左闭右开）
func (p Period) MatchesPeriodType(pt PeriodType) bool {
	if !p.IsValid() {
		return false
	}

	switch pt {
	case PeriodDay:
		// 必须是同一天的开始到次日开始（左闭右开）
		dayStart := time.Date(p.Start.Year(), p.Start.Month(), p.Start.Day(), 0, 0, 0, 0, p.Start.Location())
		dayEnd := dayStart.AddDate(0, 0, 1)
		return p.Start.Equal(dayStart) && p.End.Equal(dayEnd)

	case PeriodWeek:
		// 必须是ISO周的开始（周一）到下周一（左闭右开）
		weekStart := p.Start
		// 获取ISO周的周一
		weekday := int(weekStart.Weekday())
		if weekday == 0 { // 周日
			weekday = 7
		}
		weekStart = weekStart.AddDate(0, 0, 1-weekday)
		weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())
		weekEnd := weekStart.AddDate(0, 0, 7)
		return p.Start.Equal(weekStart) && p.End.Equal(weekEnd)

	case PeriodMonth:
		// 必须是自然月的开始到下月开始（左闭右开）
		monthStart := time.Date(p.Start.Year(), p.Start.Month(), 1, 0, 0, 0, 0, p.Start.Location())
		monthEnd := monthStart.AddDate(0, 1, 0)
		return p.Start.Equal(monthStart) && p.End.Equal(monthEnd)

	case PeriodQuarter:
		// 必须是季度的开始到下季度开始（左闭右开）
		month := p.Start.Month()
		var quarterStartMonth time.Month
		switch {
		case month >= 1 && month <= 3:
			quarterStartMonth = 1
		case month >= 4 && month <= 6:
			quarterStartMonth = 4
		case month >= 7 && month <= 9:
			quarterStartMonth = 7
		case month >= 10 && month <= 12:
			quarterStartMonth = 10
		}
		quarterStart := time.Date(p.Start.Year(), quarterStartMonth, 1, 0, 0, 0, 0, p.Start.Location())
		quarterEnd := quarterStart.AddDate(0, 3, 0)
		return p.Start.Equal(quarterStart) && p.End.Equal(quarterEnd)

	case PeriodYear:
		// 必须是自然年的开始到下年开始（左闭右开）
		yearStart := time.Date(p.Start.Year(), 1, 1, 0, 0, 0, 0, p.Start.Location())
		yearEnd := yearStart.AddDate(1, 0, 0)
		return p.Start.Equal(yearStart) && p.End.Equal(yearEnd)

	default:
		return false
	}
}

// 判断当前Period是否被另一个Period包含
func (p Period) IsWithin(ref Period) bool {
	if !p.IsValid() || !ref.IsValid() {
		return false
	}
	// 当前Period完全在ref内：p.Start >= ref.Start && p.End <= ref.End
	return (p.Start.Equal(ref.Start) || p.Start.After(ref.Start)) &&
		(p.End.Equal(ref.End) || p.End.Before(ref.End))
}

// 判断某个时间点是否在当前Period内（左闭右开）
func (p Period) ContainsTime(t time.Time) bool {
	if !p.IsValid() {
		return false
	}
	// 左闭右开：t >= Start && t < End
	return (t.Equal(p.Start) || t.After(p.Start)) && t.Before(p.End)
}

// 自动检测当前Period的周期类型
func (p Period) DetectType() PeriodType {
	if !p.IsValid() {
		return PeriodDay // 默认值
	}

	// 先检查是否符合各种标准类型
	if p.MatchesPeriodType(PeriodDay) {
		return PeriodDay
	}
	if p.MatchesPeriodType(PeriodWeek) {
		return PeriodWeek
	}
	if p.MatchesPeriodType(PeriodMonth) {
		return PeriodMonth
	}
	if p.MatchesPeriodType(PeriodQuarter) {
		return PeriodQuarter
	}
	if p.MatchesPeriodType(PeriodYear) {
		return PeriodYear
	}

	// 如果都不匹配，根据时间跨度粗略判断
	duration := p.End.Sub(p.Start)
	days := int(duration.Hours() / 24)

	if days <= 1 {
		return PeriodDay
	} else if days <= 7 {
		return PeriodWeek
	} else if days <= 31 {
		return PeriodMonth
	} else if days <= 93 { // 约3个月
		return PeriodQuarter
	} else {
		return PeriodYear
	}
}

// NewPeriod 创建新的时间周期，确保时间格式规范化（左闭右开）
func NewPeriod(start, end time.Time) (Period, error) {
	// 规范化时间：去除时分秒，只保留日期部分
	normalizedStart := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	normalizedEnd := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

	period := Period{
		Start: normalizedStart,
		End:   normalizedEnd,
	}

	if !period.IsValid() {
		return Period{}, fmt.Errorf("invalid period range: start must be before end")
	}

	return period, nil
}

// NewPeriodFromPeriodType 根据周期类型和参考时间创建标准的时间周期（左闭右开）
// 例如：
// - PeriodDay: 输入某天，返回该天 00:00:00 到次日 00:00:00
// - PeriodWeek: 输入某天，返回包含该天的 ISO 周（周一到下周一）
// - PeriodMonth: 输入某天，返回包含该天的自然月（月初到下月初）
// - PeriodQuarter: 输入某天，返回包含该天的季度（季度初到下季度初）
// - PeriodYear: 输入某天，返回包含该天的自然年（年初到下年初）
func NewPeriodFromPeriodType(pt PeriodType, referenceTime time.Time) Period {
	switch pt {
	case PeriodDay:
		start := time.Date(referenceTime.Year(), referenceTime.Month(), referenceTime.Day(), 0, 0, 0, 0, referenceTime.Location())
		end := start.AddDate(0, 0, 1)
		return Period{Start: start, End: end}

	case PeriodWeek:
		// ISO周：从周一开始
		weekday := int(referenceTime.Weekday())
		if weekday == 0 { // 周日
			weekday = 7
		}
		start := referenceTime.AddDate(0, 0, 1-weekday)
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		end := start.AddDate(0, 0, 7)
		return Period{Start: start, End: end}

	case PeriodMonth:
		start := time.Date(referenceTime.Year(), referenceTime.Month(), 1, 0, 0, 0, 0, referenceTime.Location())
		end := start.AddDate(0, 1, 0)
		return Period{Start: start, End: end}

	case PeriodQuarter:
		month := referenceTime.Month()
		var quarterStartMonth time.Month
		switch {
		case month >= 1 && month <= 3:
			quarterStartMonth = 1
		case month >= 4 && month <= 6:
			quarterStartMonth = 4
		case month >= 7 && month <= 9:
			quarterStartMonth = 7
		case month >= 10 && month <= 12:
			quarterStartMonth = 10
		}
		start := time.Date(referenceTime.Year(), quarterStartMonth, 1, 0, 0, 0, 0, referenceTime.Location())
		end := start.AddDate(0, 3, 0)
		return Period{Start: start, End: end}

	case PeriodYear:
		start := time.Date(referenceTime.Year(), 1, 1, 0, 0, 0, 0, referenceTime.Location())
		end := start.AddDate(1, 0, 0)
		return Period{Start: start, End: end}

	default:
		// 默认返回一天
		start := time.Date(referenceTime.Year(), referenceTime.Month(), referenceTime.Day(), 0, 0, 0, 0, referenceTime.Location())
		end := start.AddDate(0, 0, 1)
		return Period{Start: start, End: end}
	}
}
