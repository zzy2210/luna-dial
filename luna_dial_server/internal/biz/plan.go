package biz

import "context"

type Plan struct {
	Tasks         []*Task     `json:"tasks"`
	TasksTotal    int         `json:"tasks_total"`
	Journals      []*Journal  `json:"journals"`
	JournalsTotal int         `json:"journals_total"`
	PlanType      PeriodType  `json:"plan_type"`
	PlanPeriod    Period      `json:"plan_period"`
	ScoreTotal    int         `json:"score_total"`
	GroupStats    []GroupStat `json:"group_stats"`
}

type GroupStat struct {
	GroupKey   string `json:"group_key"`   // 分组键：日(2025-01-15)、周(2025-W03)、月(2025-01)、季度(2025-Q1)、年(2025)
	TaskCount  int    `json:"task_count"`  // 该分组内的任务总数
	ScoreTotal int    `json:"score_total"` // 该分组内的分数总和
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
	if param.UserID == "" {
		return nil, ErrNoPermission
	}

	tasks, err := uc.taskUsecase.ListTaskByPeriod(ctx, ListTaskByPeriodParam(param))
	if err != nil {
		return nil, err
	}

	journals, err := uc.journalUsecase.ListJournalByPeriod(ctx, ListJournalByPeriodParam(param))
	if err != nil {
		return nil, err
	}

	// 将[]Task转换为[]*Task
	taskPointers := make([]*Task, len(tasks))
	for i := range tasks {
		taskPointers[i] = &tasks[i]
	}

	// 获取统计信息
	groupStats, err := uc.GetPlanStats(ctx, GetPlanStatsParam(param))
	if err != nil {
		return nil, err
	}

	// 计算总分数
	var scoreTotal int
	for _, stat := range groupStats {
		scoreTotal += stat.ScoreTotal
	}

	plan := &Plan{
		Tasks:         taskPointers,
		TasksTotal:    len(tasks),
		Journals:      journals,
		JournalsTotal: len(journals),
		PlanType:      param.GroupBy,
		PlanPeriod:    param.Period,
		ScoreTotal:    scoreTotal,
		GroupStats:    groupStats,
	}

	return plan, nil
}

// 获取指定时间的统计
func (uc *PlanUsecase) GetPlanStats(ctx context.Context, param GetPlanStatsParam) ([]GroupStat, error) {
	if param.UserID == "" {
		return nil, ErrNoPermission
	}
	if !param.Period.IsValid() {
		return nil, ErrPlanPeriodInvalid
	}

	// 直接调用TaskUsecase的GetTaskStats方法
	return uc.taskUsecase.GetTaskStats(ctx, GetTaskStatsParam(param))
}
