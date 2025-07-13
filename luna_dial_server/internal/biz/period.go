package biz

import "time"

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

// 获取某个时间点在指定周期类型下的区间
func (pt PeriodType) RangeOf(t time.Time) Period {
	// TODO: 实现
	return Period{}
}

// 判断当前 Period 是否完全在另一个周期内
func (p Period) In(ref Period) bool {
	// TODO: 实现
	return false
}

// 判断某个时间点是否在当前 Period 内
func (p Period) Contains(t time.Time) bool {
	// TODO: 实现
	return false
}

// 获取当前 Period 的周期类型（如自动判断是月、周等）
func (p Period) Type() PeriodType {
	// TODO: 实现
	return PeriodDay
}
