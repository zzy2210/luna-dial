package service

import (
	"context"
	"strings"
	"time"

	"okr-web/internal/repository"
	"okr-web/internal/types"

	"github.com/google/uuid"
)

// StatsServiceImpl 统计服务实现
type StatsServiceImpl struct {
	taskRepo    repository.TaskRepository
	journalRepo repository.JournalRepository
	userRepo    repository.UserRepository
}

// NewStatsService 创建统计服务
func NewStatsService(taskRepo repository.TaskRepository, journalRepo repository.JournalRepository, userRepo repository.UserRepository) StatsService {
	return &StatsServiceImpl{
		taskRepo:    taskRepo,
		journalRepo: journalRepo,
		userRepo:    userRepo,
	}
}

// GetTaskCompletionStats 获取任务完成统计
func (s *StatsServiceImpl) GetTaskCompletionStats(ctx context.Context, userID uuid.UUID, timeRange TimeRangeRequest) (*TaskCompletionStats, error) {
	// 验证用户存在
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, &types.AppError{
			Code:    404,
			Message: "用户不存在",
			Type:    "USER_NOT_FOUND",
		}
	}

	// 验证时间范围
	if timeRange.StartDate.After(timeRange.EndDate) {
		return nil, &types.AppError{
			Code:    400,
			Message: "开始时间不能晚于结束时间",
			Type:    "INVALID_TIME_RANGE",
		}
	}

	// 获取用户所有任务（这里需要Repository支持时间范围过滤）
	tasks, err := s.taskRepo.GetByUserID(ctx, userID, 1000, 0)
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "获取任务失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// 统计数据
	stats := &TaskCompletionStats{
		TotalTasks:     0,
		CompletedTasks: 0,
		CompletionRate: 0.0,
		ByType:         make(map[types.TaskType]int64),
		ByStatus:       make(map[types.TaskStatus]int64),
	}

	// 在内存中过滤时间范围并统计
	for _, task := range tasks {
		// 检查任务是否在时间范围内
		if (task.CreatedAt.After(timeRange.StartDate) || task.CreatedAt.Equal(timeRange.StartDate)) &&
			(task.CreatedAt.Before(timeRange.EndDate) || task.CreatedAt.Equal(timeRange.EndDate)) {

			stats.TotalTasks++

			// 按类型统计
			taskType := types.TaskType(task.Type)
			stats.ByType[taskType]++

			// 按状态统计
			taskStatus := types.TaskStatus(task.Status)
			stats.ByStatus[taskStatus]++

			if taskStatus == types.TaskStatusCompleted {
				stats.CompletedTasks++
			}
		}
	}

	// 计算完成率
	if stats.TotalTasks > 0 {
		stats.CompletionRate = float64(stats.CompletedTasks) / float64(stats.TotalTasks) * 100
	}

	return stats, nil
}

// GetScoreTrend 获取评分趋势
func (s *StatsServiceImpl) GetScoreTrend(ctx context.Context, userID uuid.UUID, timeRange TimeRangeRequest) (*ScoreTrendStats, error) {
	// 验证用户存在
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, &types.AppError{
			Code:    404,
			Message: "用户不存在",
			Type:    "USER_NOT_FOUND",
		}
	}

	// 验证时间范围
	if timeRange.StartDate.After(timeRange.EndDate) {
		return nil, &types.AppError{
			Code:    400,
			Message: "开始时间不能晚于结束时间",
			Type:    "INVALID_TIME_RANGE",
		}
	}

	// 获取用户所有任务
	tasks, err := s.taskRepo.GetByUserID(ctx, userID, 1000, 0)
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "获取任务失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// 按日期分组统计评分
	dailyScores := make(map[string][]int)

	for _, task := range tasks {
		// 检查任务是否在时间范围内且已完成
		if (task.CreatedAt.After(timeRange.StartDate) || task.CreatedAt.Equal(timeRange.StartDate)) &&
			(task.CreatedAt.Before(timeRange.EndDate) || task.CreatedAt.Equal(timeRange.EndDate)) &&
			string(task.Status) == string(types.TaskStatusCompleted) &&
			task.Score > 0 {

			dateKey := task.CreatedAt.Format("2006-01-02")
			dailyScores[dateKey] = append(dailyScores[dateKey], task.Score)
		}
	}

	// 计算趋势数据
	var trend []ScoreTrendPoint
	var totalScore float64
	var totalCount int

	for dateStr, scores := range dailyScores {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		// 计算当日平均分
		var daySum int
		for _, score := range scores {
			daySum += score
			totalScore += float64(score)
			totalCount++
		}

		avgScore := float64(daySum) / float64(len(scores))
		trend = append(trend, ScoreTrendPoint{
			Date:  date,
			Score: avgScore,
		})
	}

	// 计算整体平均分
	var averageScore float64
	if totalCount > 0 {
		averageScore = totalScore / float64(totalCount)
	}

	return &ScoreTrendStats{
		AverageScore: averageScore,
		Trend:        trend,
	}, nil
}

// GetTimeDistribution 获取时间分布统计
func (s *StatsServiceImpl) GetTimeDistribution(ctx context.Context, userID uuid.UUID, timeRange TimeRangeRequest) (*TimeDistributionStats, error) {
	// 验证用户存在
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, &types.AppError{
			Code:    404,
			Message: "用户不存在",
			Type:    "USER_NOT_FOUND",
		}
	}

	// 验证时间范围
	if timeRange.StartDate.After(timeRange.EndDate) {
		return nil, &types.AppError{
			Code:    400,
			Message: "开始时间不能晚于结束时间",
			Type:    "INVALID_TIME_RANGE",
		}
	}

	// 获取用户所有任务
	tasks, err := s.taskRepo.GetByUserID(ctx, userID, 1000, 0)
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "获取任务失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// 获取用户所有日志
	journals, err := s.journalRepo.GetByUserID(ctx, userID, 1000, 0)
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "获取日志失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// 统计数据
	stats := &TimeDistributionStats{
		ByType:      make(map[types.TaskType]int64),
		ByTimeScale: make(map[types.TimeScale]int64),
	}

	// 统计任务类型分布
	for _, task := range tasks {
		if (task.CreatedAt.After(timeRange.StartDate) || task.CreatedAt.Equal(timeRange.StartDate)) &&
			(task.CreatedAt.Before(timeRange.EndDate) || task.CreatedAt.Equal(timeRange.EndDate)) {

			taskType := types.TaskType(task.Type)
			stats.ByType[taskType]++
		}
	}

	// 统计日志时间尺度分布
	for _, journal := range journals {
		if (journal.CreatedAt.After(timeRange.StartDate) || journal.CreatedAt.Equal(timeRange.StartDate)) &&
			(journal.CreatedAt.Before(timeRange.EndDate) || journal.CreatedAt.Equal(timeRange.EndDate)) {

			timeScale := types.TimeScale(journal.TimeScale)
			stats.ByTimeScale[timeScale]++
		}
	}

	return stats, nil
}

// GetUserOverview 获取用户概览统计
func (s *StatsServiceImpl) GetUserOverview(ctx context.Context, userID uuid.UUID) (*UserOverviewStats, error) {
	// 验证用户存在
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, &types.AppError{
			Code:    404,
			Message: "用户不存在",
			Type:    "USER_NOT_FOUND",
		}
	}

	// 获取任务总数
	taskCount, err := s.taskRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "获取任务统计失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// 获取日志总数
	journalCount, err := s.journalRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "获取日志统计失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	// 获取用户所有任务来计算详细统计
	tasks, err := s.taskRepo.GetByUserID(ctx, userID, 1000, 0)
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "获取任务失败",
			Type:    "INTERNAL_SERVER_ERROR",
		}
	}

	var completedTasks int64
	var totalScore float64
	var scoreCount int
	var activeTasks int64

	for _, task := range tasks {
		if string(task.Status) == string(types.TaskStatusCompleted) {
			completedTasks++
			if task.Score > 0 {
				totalScore += float64(task.Score)
				scoreCount++
			}
		} else if string(task.Status) == string(types.TaskStatusInProgress) {
			activeTasks++
		}
	}

	// 计算平均分
	var averageScore float64
	if scoreCount > 0 {
		averageScore = totalScore / float64(scoreCount)
	}

	return &UserOverviewStats{
		TotalTasks:       int64(taskCount),
		CompletedTasks:   completedTasks,
		AverageScore:     averageScore,
		TotalJournals:    int64(journalCount),
		ActiveTasksCount: activeTasks,
	}, nil
}

// GetScoreTrendByReference 获取基于时间参考的分数趋势统计
func (s *StatsServiceImpl) GetScoreTrendByReference(ctx context.Context, userID uuid.UUID, req ScoreTrendRequest) (*types.ScoreTrendResponse, error) {
	// 验证用户存在
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, &types.AppError{
			Code:    404,
			Message: "用户不存在",
			Type:    "USER_NOT_FOUND",
		}
	}

	// 解析时间参考
	timeRange, err := types.ParseTimeReference(req.TimeRef, req.Scale)
	if err != nil {
		return nil, &types.AppError{
			Code:    400,
			Message: "无效的时间参考: " + err.Error(),
			Type:    "INVALID_TIME_REFERENCE",
		}
	}

	// 获取分数统计数据（按天）
	scorePoints, err := s.taskRepo.GetScoreStatsInTimeRange(ctx, userID, timeRange.Start, timeRange.End)
	if err != nil {
		return nil, &types.AppError{
			Code:    500,
			Message: "获取分数统计失败: " + err.Error(),
			Type:    "DATABASE_ERROR",
		}
	}

	labels := types.GenerateTimeLabels(timeRange, req.Scale)

	// 先将所有分数点按天映射
	type dayStat struct{ Score, Count int }
	dayMap := make(map[string]dayStat)
	for _, point := range scorePoints {
		dateStr := point.Date.Format("2006-01-02")
		dayMap[dateStr] = dayStat{Score: point.Score, Count: point.Count}
	}

	scores := make([]int, len(labels))
	counts := make([]int, len(labels))

	totalScore, totalTasks := 0, 0
	maxScore, maxTasks := 0, 0
	minScore, minTasks := -1, -1

	for i, label := range labels {
		var sumScore, sumCount int
		//var daysInPeriod int // 已不再需要
		var periodDays []string
		// 解析 label 范围
		switch req.Scale {
		case types.TimeScaleDay:
			periodDays = []string{label}
		case types.TimeScaleWeek:
			// label: YYYY-MM-DD~YYYY-MM-DD
			parts := strings.Split(label, "~")
			if len(parts) == 2 {
				start, _ := time.Parse("2006-01-02", parts[0])
				end, _ := time.Parse("2006-01-02", parts[1])
				for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
					periodDays = append(periodDays, d.Format("2006-01-02"))
				}
			}
		case types.TimeScaleMonth:
			start, _ := time.Parse("2006-01", label)
			end := start.AddDate(0, 1, -1)
			for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
				periodDays = append(periodDays, d.Format("2006-01-02"))
			}
		case types.TimeScaleQuarter:
			tr, _ := types.ParseTimeReference(label, types.TimeScaleQuarter)
			for d := tr.Start; !d.After(tr.End); d = d.AddDate(0, 0, 1) {
				periodDays = append(periodDays, d.Format("2006-01-02"))
			}
		case types.TimeScaleYear:
			tr, _ := types.ParseTimeReference(label, types.TimeScaleYear)
			for d := tr.Start; !d.After(tr.End); d = d.AddDate(0, 0, 1) {
				periodDays = append(periodDays, d.Format("2006-01-02"))
			}
		}
		//daysInPeriod = len(periodDays) // 已不再需要
		for _, day := range periodDays {
			if stat, ok := dayMap[day]; ok {
				sumScore += stat.Score
				sumCount += stat.Count
			}
		}
		scores[i] = sumScore
		counts[i] = sumCount
		totalScore += sumScore
		totalTasks += sumCount
		if sumScore > maxScore {
			maxScore = sumScore
		}
		if sumCount > maxTasks {
			maxTasks = sumCount
		}
		if minScore == -1 || (sumScore > 0 && sumScore < minScore) {
			minScore = sumScore
		}
		if minTasks == -1 || (sumCount > 0 && sumCount < minTasks) {
			minTasks = sumCount
		}
	}
	if minScore == -1 {
		minScore = 0
	}
	if minTasks == -1 {
		minTasks = 0
	}
	avgScore := 0.0
	avgTasks := 0.0
	if len(labels) > 0 {
		avgScore = float64(totalScore) / float64(len(labels))
		avgTasks = float64(totalTasks) / float64(len(labels))
	}
	return &types.ScoreTrendResponse{
		Labels:    labels,
		Scores:    scores,
		Counts:    counts,
		Scale:     req.Scale,
		TimeRef:   req.TimeRef,
		TimeRange: timeRange,
		Summary: &types.TrendSummary{
			TotalScore:       totalScore,
			TotalTasks:       totalTasks,
			AverageScore:     avgScore,
			AverageTaskCount: avgTasks,
			MaxScore:         maxScore,
			MaxTasks:         maxTasks,
			MinScore:         minScore,
			MinTasks:         minTasks,
		},
	}, nil
}
