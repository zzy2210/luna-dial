package biz

import (
	"testing"
	"time"
)

func TestPeriod_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		period Period
		want   bool
	}{
		{
			name: "正常情况 - Start < End",
			period: Period{
				Start: time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 7, 17, 0, 0, 0, 0, time.UTC),
			},
			want: true,
		},
		{
			name: "Start == End",
			period: Period{
				Start: time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC),
			},
			want: false,
		},
		{
			name: "Start > End",
			period: Period{
				Start: time.Date(2025, 7, 17, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC),
			},
			want: false,
		},
		{
			name: "Start 为零值",
			period: Period{
				Start: time.Time{},
				End:   time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC),
			},
			want: false,
		},
		{
			name: "End 为零值",
			period: Period{
				Start: time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC),
				End:   time.Time{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.period.IsValid(); got != tt.want {
				t.Errorf("Period.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriod_MatchesPeriodType(t *testing.T) {
	// 测试用的时区
	utc := time.UTC

	tests := []struct {
		name       string
		period     Period
		periodType PeriodType
		want       bool
	}{
		// PeriodDay 测试
		{
			name: "PeriodDay - 完全匹配的一天",
			period: Period{
				Start: time.Date(2025, 7, 16, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 17, 0, 0, 0, 0, utc),
			},
			periodType: PeriodDay,
			want:       true,
		},
		{
			name: "PeriodDay - 不是从00:00:00开始",
			period: Period{
				Start: time.Date(2025, 7, 16, 1, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 17, 0, 0, 0, 0, utc),
			},
			periodType: PeriodDay,
			want:       false,
		},
		{
			name: "PeriodDay - 跨越多天",
			period: Period{
				Start: time.Date(2025, 7, 16, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 18, 0, 0, 0, 0, utc),
			},
			periodType: PeriodDay,
			want:       false,
		},
		// PeriodWeek 测试
		{
			name: "PeriodWeek - 完全匹配的ISO周",
			period: Period{
				Start: time.Date(2025, 7, 14, 0, 0, 0, 0, utc), // 2025年7月14日是周一
				End:   time.Date(2025, 7, 21, 0, 0, 0, 0, utc), // 下周一
			},
			periodType: PeriodWeek,
			want:       true,
		},
		{
			name: "PeriodWeek - 不是从周一开始",
			period: Period{
				Start: time.Date(2025, 7, 15, 0, 0, 0, 0, utc), // 周二开始
				End:   time.Date(2025, 7, 22, 0, 0, 0, 0, utc),
			},
			periodType: PeriodWeek,
			want:       false,
		},
		// PeriodMonth 测试
		{
			name: "PeriodMonth - 完全匹配的自然月",
			period: Period{
				Start: time.Date(2025, 7, 1, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 8, 1, 0, 0, 0, 0, utc),
			},
			periodType: PeriodMonth,
			want:       true,
		},
		{
			name: "PeriodMonth - 不是从1号开始",
			period: Period{
				Start: time.Date(2025, 7, 2, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 8, 1, 0, 0, 0, 0, utc),
			},
			periodType: PeriodMonth,
			want:       false,
		},
		{
			name: "PeriodMonth - 2月份闰年",
			period: Period{
				Start: time.Date(2024, 2, 1, 0, 0, 0, 0, utc), // 2024是闰年
				End:   time.Date(2024, 3, 1, 0, 0, 0, 0, utc),
			},
			periodType: PeriodMonth,
			want:       true,
		},
		// PeriodQuarter 测试
		{
			name: "PeriodQuarter - Q1",
			period: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 4, 1, 0, 0, 0, 0, utc),
			},
			periodType: PeriodQuarter,
			want:       true,
		},
		{
			name: "PeriodQuarter - Q2",
			period: Period{
				Start: time.Date(2025, 4, 1, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 1, 0, 0, 0, 0, utc),
			},
			periodType: PeriodQuarter,
			want:       true,
		},
		{
			name: "PeriodQuarter - Q3",
			period: Period{
				Start: time.Date(2025, 7, 1, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 10, 1, 0, 0, 0, 0, utc),
			},
			periodType: PeriodQuarter,
			want:       true,
		},
		{
			name: "PeriodQuarter - Q4",
			period: Period{
				Start: time.Date(2025, 10, 1, 0, 0, 0, 0, utc),
				End:   time.Date(2026, 1, 1, 0, 0, 0, 0, utc),
			},
			periodType: PeriodQuarter,
			want:       true,
		},
		{
			name: "PeriodQuarter - 不完整的季度",
			period: Period{
				Start: time.Date(2025, 1, 15, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 4, 1, 0, 0, 0, 0, utc),
			},
			periodType: PeriodQuarter,
			want:       false,
		},
		// PeriodYear 测试
		{
			name: "PeriodYear - 完全匹配的自然年",
			period: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, utc),
				End:   time.Date(2026, 1, 1, 0, 0, 0, 0, utc),
			},
			periodType: PeriodYear,
			want:       true,
		},
		{
			name: "PeriodYear - 闰年",
			period: Period{
				Start: time.Date(2024, 1, 1, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 1, 1, 0, 0, 0, 0, utc),
			},
			periodType: PeriodYear,
			want:       true,
		},
		{
			name: "PeriodYear - 不完整的年份",
			period: Period{
				Start: time.Date(2025, 2, 1, 0, 0, 0, 0, utc),
				End:   time.Date(2026, 1, 1, 0, 0, 0, 0, utc),
			},
			periodType: PeriodYear,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.period.MatchesPeriodType(tt.periodType); got != tt.want {
				t.Errorf("Period.MatchesPeriodType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPeriod(t *testing.T) {
	utc := time.UTC

	tests := []struct {
		name       string
		start      time.Time
		end        time.Time
		wantPeriod Period
		wantError  bool
	}{
		{
			name:  "正常创建 - start < end",
			start: time.Date(2025, 7, 16, 10, 30, 45, 0, utc),
			end:   time.Date(2025, 7, 17, 15, 20, 30, 0, utc),
			wantPeriod: Period{
				Start: time.Date(2025, 7, 16, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 17, 0, 0, 0, 0, utc),
			},
			wantError: false,
		},
		{
			name:       "错误情况 - start == end",
			start:      time.Date(2025, 7, 16, 10, 30, 45, 0, utc),
			end:        time.Date(2025, 7, 16, 15, 20, 30, 0, utc),
			wantPeriod: Period{},
			wantError:  true,
		},
		{
			name:       "错误情况 - start > end",
			start:      time.Date(2025, 7, 17, 10, 30, 45, 0, utc),
			end:        time.Date(2025, 7, 16, 15, 20, 30, 0, utc),
			wantPeriod: Period{},
			wantError:  true,
		},
		{
			name:  "时间规范化测试",
			start: time.Date(2025, 7, 16, 23, 59, 59, 999999999, utc),
			end:   time.Date(2025, 7, 18, 1, 2, 3, 123456789, utc),
			wantPeriod: Period{
				Start: time.Date(2025, 7, 16, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 18, 0, 0, 0, 0, utc),
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPeriod(tt.start, tt.end)
			if (err != nil) != tt.wantError {
				t.Errorf("NewPeriod() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				if !got.Start.Equal(tt.wantPeriod.Start) || !got.End.Equal(tt.wantPeriod.End) {
					t.Errorf("NewPeriod() = %v, want %v", got, tt.wantPeriod)
				}
			}
		})
	}
}

func TestNewPeriodFromPeriodType(t *testing.T) {
	utc := time.UTC

	tests := []struct {
		name          string
		periodType    PeriodType
		referenceTime time.Time
		wantStart     time.Time
		wantEnd       time.Time
	}{
		// PeriodDay 测试
		{
			name:          "PeriodDay - 任意时间",
			periodType:    PeriodDay,
			referenceTime: time.Date(2025, 7, 16, 14, 30, 45, 0, utc),
			wantStart:     time.Date(2025, 7, 16, 0, 0, 0, 0, utc),
			wantEnd:       time.Date(2025, 7, 17, 0, 0, 0, 0, utc),
		},
		{
			name:          "PeriodDay - 月末",
			periodType:    PeriodDay,
			referenceTime: time.Date(2025, 7, 31, 23, 59, 59, 0, utc),
			wantStart:     time.Date(2025, 7, 31, 0, 0, 0, 0, utc),
			wantEnd:       time.Date(2025, 8, 1, 0, 0, 0, 0, utc),
		},
		{
			name:          "PeriodDay - 年末",
			periodType:    PeriodDay,
			referenceTime: time.Date(2025, 12, 31, 12, 0, 0, 0, utc),
			wantStart:     time.Date(2025, 12, 31, 0, 0, 0, 0, utc),
			wantEnd:       time.Date(2026, 1, 1, 0, 0, 0, 0, utc),
		},
		// PeriodWeek 测试
		{
			name:          "PeriodWeek - 周一",
			periodType:    PeriodWeek,
			referenceTime: time.Date(2025, 7, 14, 10, 0, 0, 0, utc), // 2025年7月14日是周一
			wantStart:     time.Date(2025, 7, 14, 0, 0, 0, 0, utc),
			wantEnd:       time.Date(2025, 7, 21, 0, 0, 0, 0, utc),
		},
		{
			name:          "PeriodWeek - 周日",
			periodType:    PeriodWeek,
			referenceTime: time.Date(2025, 7, 20, 15, 0, 0, 0, utc), // 2025年7月20日是周日
			wantStart:     time.Date(2025, 7, 14, 0, 0, 0, 0, utc),  // 本周的周一
			wantEnd:       time.Date(2025, 7, 21, 0, 0, 0, 0, utc),
		},
		{
			name:          "PeriodWeek - 周三",
			periodType:    PeriodWeek,
			referenceTime: time.Date(2025, 7, 16, 12, 0, 0, 0, utc), // 2025年7月16日是周三
			wantStart:     time.Date(2025, 7, 14, 0, 0, 0, 0, utc),  // 本周的周一
			wantEnd:       time.Date(2025, 7, 21, 0, 0, 0, 0, utc),
		},
		{
			name:          "PeriodWeek - 跨月的周",
			periodType:    PeriodWeek,
			referenceTime: time.Date(2025, 7, 31, 10, 0, 0, 0, utc), // 2025年7月31日是周四
			wantStart:     time.Date(2025, 7, 28, 0, 0, 0, 0, utc),  // 本周的周一
			wantEnd:       time.Date(2025, 8, 4, 0, 0, 0, 0, utc),   // 下周一
		},
		// PeriodMonth 测试
		{
			name:          "PeriodMonth - 月中任意日期",
			periodType:    PeriodMonth,
			referenceTime: time.Date(2025, 7, 16, 14, 30, 0, 0, utc),
			wantStart:     time.Date(2025, 7, 1, 0, 0, 0, 0, utc),
			wantEnd:       time.Date(2025, 8, 1, 0, 0, 0, 0, utc),
		},
		{
			name:          "PeriodMonth - 2月份闰年",
			periodType:    PeriodMonth,
			referenceTime: time.Date(2024, 2, 15, 10, 0, 0, 0, utc), // 2024是闰年
			wantStart:     time.Date(2024, 2, 1, 0, 0, 0, 0, utc),
			wantEnd:       time.Date(2024, 3, 1, 0, 0, 0, 0, utc),
		},
		{
			name:          "PeriodMonth - 12月跨年",
			periodType:    PeriodMonth,
			referenceTime: time.Date(2025, 12, 25, 18, 0, 0, 0, utc),
			wantStart:     time.Date(2025, 12, 1, 0, 0, 0, 0, utc),
			wantEnd:       time.Date(2026, 1, 1, 0, 0, 0, 0, utc),
		},
		// PeriodQuarter 测试
		{
			name:          "PeriodQuarter - Q1中的日期",
			periodType:    PeriodQuarter,
			referenceTime: time.Date(2025, 2, 15, 10, 0, 0, 0, utc),
			wantStart:     time.Date(2025, 1, 1, 0, 0, 0, 0, utc),
			wantEnd:       time.Date(2025, 4, 1, 0, 0, 0, 0, utc),
		},
		{
			name:          "PeriodQuarter - Q2中的日期",
			periodType:    PeriodQuarter,
			referenceTime: time.Date(2025, 5, 10, 14, 0, 0, 0, utc),
			wantStart:     time.Date(2025, 4, 1, 0, 0, 0, 0, utc),
			wantEnd:       time.Date(2025, 7, 1, 0, 0, 0, 0, utc),
		},
		{
			name:          "PeriodQuarter - Q3中的日期",
			periodType:    PeriodQuarter,
			referenceTime: time.Date(2025, 8, 20, 16, 0, 0, 0, utc),
			wantStart:     time.Date(2025, 7, 1, 0, 0, 0, 0, utc),
			wantEnd:       time.Date(2025, 10, 1, 0, 0, 0, 0, utc),
		},
		{
			name:          "PeriodQuarter - Q4中的日期",
			periodType:    PeriodQuarter,
			referenceTime: time.Date(2025, 11, 5, 9, 0, 0, 0, utc),
			wantStart:     time.Date(2025, 10, 1, 0, 0, 0, 0, utc),
			wantEnd:       time.Date(2026, 1, 1, 0, 0, 0, 0, utc),
		},
		// PeriodYear 测试
		{
			name:          "PeriodYear - 年中任意日期",
			periodType:    PeriodYear,
			referenceTime: time.Date(2025, 7, 16, 14, 30, 0, 0, utc),
			wantStart:     time.Date(2025, 1, 1, 0, 0, 0, 0, utc),
			wantEnd:       time.Date(2026, 1, 1, 0, 0, 0, 0, utc),
		},
		{
			name:          "PeriodYear - 闰年",
			periodType:    PeriodYear,
			referenceTime: time.Date(2024, 6, 15, 12, 0, 0, 0, utc), // 2024是闰年
			wantStart:     time.Date(2024, 1, 1, 0, 0, 0, 0, utc),
			wantEnd:       time.Date(2025, 1, 1, 0, 0, 0, 0, utc),
		},
		// Default 测试
		{
			name:          "Default - 无效的PeriodType",
			periodType:    PeriodType(999), // 无效的类型
			referenceTime: time.Date(2025, 7, 16, 14, 30, 0, 0, utc),
			wantStart:     time.Date(2025, 7, 16, 0, 0, 0, 0, utc),
			wantEnd:       time.Date(2025, 7, 17, 0, 0, 0, 0, utc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPeriodFromPeriodType(tt.periodType, tt.referenceTime)
			if !got.Start.Equal(tt.wantStart) {
				t.Errorf("NewPeriodFromPeriodType() Start = %v, want %v", got.Start, tt.wantStart)
			}
			if !got.End.Equal(tt.wantEnd) {
				t.Errorf("NewPeriodFromPeriodType() End = %v, want %v", got.End, tt.wantEnd)
			}
		})
	}
}

// 集成测试：NewPeriodFromPeriodType 创建的 Period 应该能通过对应的 MatchesPeriodType 验证
func TestIntegration_NewPeriodFromPeriodType_MatchesPeriodType(t *testing.T) {
	utc := time.UTC
	referenceTime := time.Date(2025, 7, 16, 14, 30, 45, 0, utc)

	periodTypes := []PeriodType{
		PeriodDay,
		PeriodWeek,
		PeriodMonth,
		PeriodQuarter,
		PeriodYear,
	}

	for _, pt := range periodTypes {
		t.Run(pt.String(), func(t *testing.T) {
			period := NewPeriodFromPeriodType(pt, referenceTime)
			if !period.MatchesPeriodType(pt) {
				t.Errorf("Period created by NewPeriodFromPeriodType(%v) should match PeriodType %v", pt, pt)
			}
		})
	}
}

// 边界情况测试
func TestBoundaryConditions(t *testing.T) {
	utc := time.UTC

	t.Run("闰年2月29日", func(t *testing.T) {
		leapDay := time.Date(2024, 2, 29, 12, 0, 0, 0, utc) // 2024年2月29日
		period := NewPeriodFromPeriodType(PeriodMonth, leapDay)
		expected := Period{
			Start: time.Date(2024, 2, 1, 0, 0, 0, 0, utc),
			End:   time.Date(2024, 3, 1, 0, 0, 0, 0, utc),
		}
		if !period.Start.Equal(expected.Start) || !period.End.Equal(expected.End) {
			t.Errorf("NewPeriodFromPeriodType(PeriodMonth, leapDay) = %v, want %v", period, expected)
		}
		if !period.MatchesPeriodType(PeriodMonth) {
			t.Error("闰年2月份的Period应该匹配PeriodMonth")
		}
	})

	t.Run("年末12月31日", func(t *testing.T) {
		yearEnd := time.Date(2025, 12, 31, 23, 59, 59, 0, utc)
		period := NewPeriodFromPeriodType(PeriodYear, yearEnd)
		expected := Period{
			Start: time.Date(2025, 1, 1, 0, 0, 0, 0, utc),
			End:   time.Date(2026, 1, 1, 0, 0, 0, 0, utc),
		}
		if !period.Start.Equal(expected.Start) || !period.End.Equal(expected.End) {
			t.Errorf("NewPeriodFromPeriodType(PeriodYear, yearEnd) = %v, want %v", period, expected)
		}
		if !period.MatchesPeriodType(PeriodYear) {
			t.Error("年末的Period应该匹配PeriodYear")
		}
	})

	t.Run("年初1月1日", func(t *testing.T) {
		yearStart := time.Date(2025, 1, 1, 0, 0, 0, 0, utc)
		period := NewPeriodFromPeriodType(PeriodYear, yearStart)
		expected := Period{
			Start: time.Date(2025, 1, 1, 0, 0, 0, 0, utc),
			End:   time.Date(2026, 1, 1, 0, 0, 0, 0, utc),
		}
		if !period.Start.Equal(expected.Start) || !period.End.Equal(expected.End) {
			t.Errorf("NewPeriodFromPeriodType(PeriodYear, yearStart) = %v, want %v", period, expected)
		}
		if !period.MatchesPeriodType(PeriodYear) {
			t.Error("年初的Period应该匹配PeriodYear")
		}
	})
}

// PeriodType.String() 方法用于测试输出
func (pt PeriodType) String() string {
	switch pt {
	case PeriodDay:
		return "PeriodDay"
	case PeriodWeek:
		return "PeriodWeek"
	case PeriodMonth:
		return "PeriodMonth"
	case PeriodQuarter:
		return "PeriodQuarter"
	case PeriodYear:
		return "PeriodYear"
	default:
		return "Unknown"
	}
}

// 测试 Period.IsWithin 方法
func TestPeriod_IsWithin(t *testing.T) {
	utc := time.UTC

	tests := []struct {
		name      string
		period    Period
		refPeriod Period
		want      bool
	}{
		{
			name: "完全包含",
			period: Period{
				Start: time.Date(2025, 7, 16, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 17, 0, 0, 0, 0, utc),
			},
			refPeriod: Period{
				Start: time.Date(2025, 7, 15, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 18, 0, 0, 0, 0, utc),
			},
			want: true,
		},
		{
			name: "边界相等 - 完全相同",
			period: Period{
				Start: time.Date(2025, 7, 16, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 17, 0, 0, 0, 0, utc),
			},
			refPeriod: Period{
				Start: time.Date(2025, 7, 16, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 17, 0, 0, 0, 0, utc),
			},
			want: true,
		},
		{
			name: "左边界相等",
			period: Period{
				Start: time.Date(2025, 7, 16, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 17, 0, 0, 0, 0, utc),
			},
			refPeriod: Period{
				Start: time.Date(2025, 7, 16, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 18, 0, 0, 0, 0, utc),
			},
			want: true,
		},
		{
			name: "右边界相等",
			period: Period{
				Start: time.Date(2025, 7, 16, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 17, 0, 0, 0, 0, utc),
			},
			refPeriod: Period{
				Start: time.Date(2025, 7, 15, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 17, 0, 0, 0, 0, utc),
			},
			want: true,
		},
		{
			name: "部分重叠 - 左边超出",
			period: Period{
				Start: time.Date(2025, 7, 14, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 17, 0, 0, 0, 0, utc),
			},
			refPeriod: Period{
				Start: time.Date(2025, 7, 15, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 18, 0, 0, 0, 0, utc),
			},
			want: false,
		},
		{
			name: "部分重叠 - 右边超出",
			period: Period{
				Start: time.Date(2025, 7, 16, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 19, 0, 0, 0, 0, utc),
			},
			refPeriod: Period{
				Start: time.Date(2025, 7, 15, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 18, 0, 0, 0, 0, utc),
			},
			want: false,
		},
		{
			name: "完全不重叠",
			period: Period{
				Start: time.Date(2025, 7, 20, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 21, 0, 0, 0, 0, utc),
			},
			refPeriod: Period{
				Start: time.Date(2025, 7, 15, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 18, 0, 0, 0, 0, utc),
			},
			want: false,
		},
		{
			name: "当前Period无效",
			period: Period{
				Start: time.Date(2025, 7, 17, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 16, 0, 0, 0, 0, utc), // Start > End
			},
			refPeriod: Period{
				Start: time.Date(2025, 7, 15, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 18, 0, 0, 0, 0, utc),
			},
			want: false,
		},
		{
			name: "参考Period无效",
			period: Period{
				Start: time.Date(2025, 7, 16, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 17, 0, 0, 0, 0, utc),
			},
			refPeriod: Period{
				Start: time.Date(2025, 7, 18, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 15, 0, 0, 0, 0, utc), // Start > End
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.period.IsWithin(tt.refPeriod); got != tt.want {
				t.Errorf("Period.IsWithin() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试 Period.ContainsTime 方法
func TestPeriod_ContainsTime(t *testing.T) {
	utc := time.UTC

	period := Period{
		Start: time.Date(2025, 7, 16, 0, 0, 0, 0, utc),
		End:   time.Date(2025, 7, 17, 0, 0, 0, 0, utc),
	}

	tests := []struct {
		name   string
		period Period
		time   time.Time
		want   bool
	}{
		{
			name:   "时间点在区间内",
			period: period,
			time:   time.Date(2025, 7, 16, 12, 0, 0, 0, utc),
			want:   true,
		},
		{
			name:   "时间点等于Start（左闭）",
			period: period,
			time:   time.Date(2025, 7, 16, 0, 0, 0, 0, utc),
			want:   true,
		},
		{
			name:   "时间点等于End（右开）",
			period: period,
			time:   time.Date(2025, 7, 17, 0, 0, 0, 0, utc),
			want:   false,
		},
		{
			name:   "时间点在Start之前",
			period: period,
			time:   time.Date(2025, 7, 15, 23, 59, 59, 0, utc),
			want:   false,
		},
		{
			name:   "时间点在End之后",
			period: period,
			time:   time.Date(2025, 7, 17, 0, 0, 1, 0, utc),
			want:   false,
		},
		{
			name: "Period无效",
			period: Period{
				Start: time.Date(2025, 7, 17, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 16, 0, 0, 0, 0, utc), // Start > End
			},
			time: time.Date(2025, 7, 16, 12, 0, 0, 0, utc),
			want: false,
		},
		{
			name:   "边界测试 - Start前一纳秒",
			period: period,
			time:   time.Date(2025, 7, 16, 0, 0, 0, 0, utc).Add(-time.Nanosecond),
			want:   false,
		},
		{
			name:   "边界测试 - End前一纳秒",
			period: period,
			time:   time.Date(2025, 7, 17, 0, 0, 0, 0, utc).Add(-time.Nanosecond),
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.period.ContainsTime(tt.time); got != tt.want {
				t.Errorf("Period.ContainsTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试 Period.DetectType 方法
func TestPeriod_DetectType(t *testing.T) {
	utc := time.UTC

	tests := []struct {
		name   string
		period Period
		want   PeriodType
	}{
		// 标准类型测试
		{
			name: "标准日周期",
			period: Period{
				Start: time.Date(2025, 7, 16, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 17, 0, 0, 0, 0, utc),
			},
			want: PeriodDay,
		},
		{
			name: "标准周周期",
			period: Period{
				Start: time.Date(2025, 7, 14, 0, 0, 0, 0, utc), // 周一
				End:   time.Date(2025, 7, 21, 0, 0, 0, 0, utc), // 下周一
			},
			want: PeriodWeek,
		},
		{
			name: "标准月周期",
			period: Period{
				Start: time.Date(2025, 7, 1, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 8, 1, 0, 0, 0, 0, utc),
			},
			want: PeriodMonth,
		},
		{
			name: "标准季度周期",
			period: Period{
				Start: time.Date(2025, 7, 1, 0, 0, 0, 0, utc),  // Q3开始
				End:   time.Date(2025, 10, 1, 0, 0, 0, 0, utc), // Q4开始
			},
			want: PeriodQuarter,
		},
		{
			name: "标准年周期",
			period: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, utc),
				End:   time.Date(2026, 1, 1, 0, 0, 0, 0, utc),
			},
			want: PeriodYear,
		},
		// 非标准但根据时长判断的测试
		{
			name: "非标准1天期间",
			period: Period{
				Start: time.Date(2025, 7, 16, 8, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 16, 18, 0, 0, 0, utc), // 10小时
			},
			want: PeriodDay,
		},
		{
			name: "非标准3天期间",
			period: Period{
				Start: time.Date(2025, 7, 16, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 19, 0, 0, 0, 0, utc), // 3天
			},
			want: PeriodWeek,
		},
		{
			name: "非标准15天期间",
			period: Period{
				Start: time.Date(2025, 7, 1, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 16, 0, 0, 0, 0, utc), // 15天
			},
			want: PeriodMonth,
		},
		{
			name: "非标准60天期间",
			period: Period{
				Start: time.Date(2025, 6, 1, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 8, 1, 0, 0, 0, 0, utc), // 约60天
			},
			want: PeriodQuarter,
		},
		{
			name: "非标准200天期间",
			period: Period{
				Start: time.Date(2025, 1, 1, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 8, 1, 0, 0, 0, 0, utc), // 约200天
			},
			want: PeriodYear,
		},
		// 边界情况
		{
			name: "无效Period",
			period: Period{
				Start: time.Date(2025, 7, 17, 0, 0, 0, 0, utc),
				End:   time.Date(2025, 7, 16, 0, 0, 0, 0, utc), // Start > End
			},
			want: PeriodDay, // 默认值
		},
		{
			name: "零值Period",
			period: Period{
				Start: time.Time{},
				End:   time.Time{},
			},
			want: PeriodDay, // 默认值
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.period.DetectType(); got != tt.want {
				t.Errorf("Period.DetectType() = %v, want %v", got, tt.want)
			}
		})
	}
}
