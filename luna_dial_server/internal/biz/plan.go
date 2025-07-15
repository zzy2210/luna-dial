package biz

import "context"

type Plan struct {
	Tasks         []Task      `json:"tasks"`
	TasksTotal    int         `json:"tasks_total"`
	Journals      []Journal   `json:"journals"`
	JournalsTotal int         `json:"journals_total"`
	PlanType      PeriodType  `json:"plan_type"`
	PlanPeriod    Period      `json:"plan_period"`
	ScoreTotal    int         `json:"score_total"`
	GroupStats    []GroupStat `json:"group_stats"`
}

type GroupStat struct {
	GroupKey   string `json:"group_key"`   // 分组键，如2025-01、2025-W01等
	TaskCount  int    `json:"task_count"`  // 日任务总和
	ScoreTotal int    `json:"score_total"` // 分数总和
}

// 获取指定时间的计划参数
type GetPlanByPeriodParam struct {
	UserID  string
	Period  Period
	GroupBy PeriodType
}

type GetPlanStatsParam struct {
	UserID  string
	Period  Period
	GroupBy PeriodType
}

type PlanUsecase struct {
	taskUsecase    *TaskUsecase
	journalUsecase *JournalUsecase
}

func NewPlanUsecase(taskUsecase *TaskUsecase, journalUsecase *JournalUsecase) *PlanUsecase {
	return &PlanUsecase{
		taskUsecase:    taskUsecase,
		journalUsecase: journalUsecase,
	}
}

// 获取指定时间的计划
func (uc *PlanUsecase) GetPlanByPeriod(ctx context.Context, param GetPlanByPeriodParam) (*Plan, error) {
	return nil, ErrNoPermission // TODO: 实现
}

// 获取指定时间的统计
func (uc *PlanUsecase) GetPlanStats(ctx context.Context, param GetPlanStatsParam) ([]GroupStat, error) {
	return nil, ErrNoPermission // TODO: 实现
}
