import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import useAuthStore from '../store/auth';
import ParentTaskList from '../components/ParentTaskList';
import TimeNavigator from '../components/TimeNavigator';
import taskService from '../services/task';
import journalService from '../services/journal';
import planService from '../services/plan';
import { Task, TaskStatus, TaskPriority, PeriodType, Journal, PlanResponse } from '../types';
import TaskEditDialog from '../components/TaskEditDialog';
import TaskViewDialog from '../components/TaskViewDialog';
import JournalEditDialog from '../components/JournalEditDialog';
import JournalViewDialog from '../components/JournalViewDialog';
import ScoreEditDialog from '../components/ScoreEditDialog';
import { TASK_STATUS_TO_API } from '../constants/task';
import {
  formatLocalDate,
  getWeekNumber,
  getQuarterNumber,
  getPeriodDates,
  getWeekStartDate,
  formatPeriodRange,
  getCurrentDateString,
  getPeriodLabel
} from '../utils/dateUtils';
import '../styles/dashboard.css';

// 趋势数据项接口
interface TrendDataItem {
  label: string;
  score: number;
  percentage: number;
  isCurrent: boolean;
}

const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const { user, logout } = useAuthStore();
  const [currentPeriod, setCurrentPeriod] = useState<PeriodType>('day');
  const [currentDate, setCurrentDate] = useState<Date>(new Date());
  const [journals, setJournals] = useState<Journal[]>([]);
  const [planData, setPlanData] = useState<PlanResponse | null>(null);
  const [loading, setLoading] = useState(true);

  // 对话框状态
  const [showTaskDialog, setShowTaskDialog] = useState(false);
  const [showViewTaskDialog, setShowViewTaskDialog] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);
  const [viewingTask, setViewingTask] = useState<Task | null>(null);
  const [parentTaskIdForCreation, setParentTaskIdForCreation] = useState<string | undefined>(undefined);
  const [showJournalDialog, setShowJournalDialog] = useState(false);
  const [showViewJournalDialog, setShowViewJournalDialog] = useState(false);
  const [editingJournal, setEditingJournal] = useState<Journal | null>(null);
  const [viewingJournal, setViewingJournal] = useState<Journal | null>(null);
  const [showScoreDialog, setShowScoreDialog] = useState(false);
  const [editingScoreTask, setEditingScoreTask] = useState<Task | null>(null);

  // 统计数据
  const [stats, setStats] = useState({
    todayScore: 0,
    weekScore: 0,
    monthScore: 0,
    weekProgress: [0, 0, 0, 0, 0, 0, 0],
    trendData: [] as TrendDataItem[],
    taskStats: {
      notStarted: 0,
      inProgress: 0,
      completed: 0,
      cancelled: 0
    }
  });

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  const handleTaskStatusChange = async (taskId: string, status: TaskStatus) => {
    try {
      await taskService.updateTask(taskId, {
        status: TASK_STATUS_TO_API[status]
      });

      // 刷新数据
      loadPlanData();
    } catch (error) {
      console.error('Failed to update task status:', error);
    }
  };

  const handleTaskScoreChange = async (taskId: string, score: number) => {
    try {
      await taskService.updateScore(taskId, score);
      // 刷新数据
      loadPlanData();
    } catch (error) {
      console.error('Failed to update task score:', error);
    }
  };

  const handleTaskClick = (task: Task) => {
    setViewingTask(task);
    setShowViewTaskDialog(true);
  };

  const handleCreateTask = () => {
    setEditingTask(null);
    setParentTaskIdForCreation(undefined);
    setShowTaskDialog(true);
  };

  const handleCreateSubtask = (parentTaskId: string) => {
    setParentTaskIdForCreation(parentTaskId);
    setEditingTask(null);
    setShowViewTaskDialog(false);
    setShowTaskDialog(true);
  };

  const handleEditTask = (task: Task) => {
    setEditingTask(task);
    setShowViewTaskDialog(false);
    setShowTaskDialog(true);
  };

  const handleDeleteTask = async (taskId: string) => {
    try {
      await taskService.deleteTask(taskId);
      setShowViewTaskDialog(false);
      loadPlanData();
    } catch (error) {
      console.error('Failed to delete task:', error);
      alert('删除任务失败，请重试');
    }
  };

  const handleCreateJournal = () => {
    setEditingJournal(null);
    setShowJournalDialog(true);
  };

  const handleViewJournal = (journal: Journal) => {
    setViewingJournal(journal);
    setShowViewJournalDialog(true);
  };

  const handleEditJournal = (journal: Journal) => {
    setEditingJournal(journal);
    setShowViewJournalDialog(false);
    setShowJournalDialog(true);
  };

  const handleDeleteJournal = async (journalId: string) => {
    if (window.confirm('确定要删除这条日志吗？')) {
      try {
        await journalService.deleteJournal(journalId);
        loadPlanData();
      } catch (error) {
        console.error('Failed to delete journal:', error);
      }
    }
  };

  // 打开评分编辑对话框
  const handleOpenScoreDialog = (task: Task) => {
    setEditingScoreTask(task);
    setShowScoreDialog(true);
  };

  // 获取优先级标签文本
  const getPriorityLabel = (priority: TaskPriority): string => {
    const labels = {
      [TaskPriority.Low]: '低',
      [TaskPriority.Medium]: '中',
      [TaskPriority.High]: '高',
      [TaskPriority.Urgent]: '紧急'
    };
    return labels[priority] || '中';
  };

  // 获取优先级样式类名
  const getPriorityClass = (priority: TaskPriority): string => {
    const classes = {
      [TaskPriority.Low]: 'priority-low',
      [TaskPriority.Medium]: 'priority-medium',
      [TaskPriority.High]: 'priority-high',
      [TaskPriority.Urgent]: 'priority-urgent'
    };
    return classes[priority] || 'priority-medium';
  };

  // 处理日期变化
  const handleDateChange = (newDate: Date) => {
    setCurrentDate(newDate);
  };

  // 处理前进/后退导航
  const handleNavigate = (direction: 'prev' | 'next') => {
    const newDate = new Date(currentDate);
    const offset = direction === 'prev' ? -1 : 1;

    switch (currentPeriod) {
      case 'day':
        newDate.setDate(newDate.getDate() + offset);
        break;
      case 'week':
        newDate.setDate(newDate.getDate() + (offset * 7));
        break;
      case 'month':
        newDate.setMonth(newDate.getMonth() + offset);
        break;
      case 'quarter':
        newDate.setMonth(newDate.getMonth() + (offset * 3));
        break;
      case 'year':
        newDate.setFullYear(newDate.getFullYear() + offset);
        break;
    }

    setCurrentDate(newDate);
  };


  // 获取趋势图的分组粒度
  const getTrendGroupBy = (): PeriodType => {
    switch (currentPeriod) {
      case 'day':
      case 'week':
        return 'day';     // 日/周视图：按日分组
      case 'month':
        return 'week';    // 月视图：按周分组
      case 'quarter':
        return 'month';   // 季视图：按月分组
      case 'year':
        return 'quarter'; // 年视图：按季度分组
      default:
        return 'day';
    }
  };

  // 获取趋势图的时间范围
  const getTrendPeriod = () => {
    const trendPeriod = currentPeriod === 'day' ? 'week' : currentPeriod;
    return getPeriodDates(trendPeriod, currentDate);
  };

  // 解析本周趋势（7天）
  const parseWeekTrend = (groupStats: { group_key: string; score_total: number }[]): TrendDataItem[] => {
    const weekStart = getWeekStartDate(currentDate);

    const result: TrendDataItem[] = [];
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    // 创建7天的数据
    for (let i = 0; i < 7; i++) {
      const date = new Date(weekStart);
      date.setDate(weekStart.getDate() + i);
      const dateKey = formatLocalDate(date);

      const stat = groupStats?.find(s => s.group_key === dateKey);
      const score = stat?.score_total || 0;

      result.push({
        label: ['一', '二', '三', '四', '五', '六', '日'][i],
        score,
        percentage: 0, // 将在后面归一化
        isCurrent: date.getTime() === today.getTime()
      });
    }

    // 归一化百分比
    const maxScore = Math.max(...result.map(r => r.score), 1);
    result.forEach(item => {
      item.percentage = (item.score / maxScore) * 100;
    });

    return result;
  };

  // 解析本月趋势（各周）
  const parseMonthTrend = (groupStats: { group_key: string; score_total: number }[]): TrendDataItem[] => {
    const monthStart = new Date(currentDate);
    monthStart.setDate(1);
    monthStart.setHours(0, 0, 0, 0);

    const monthEnd = new Date(currentDate);
    monthEnd.setMonth(monthEnd.getMonth() + 1, 1);
    monthEnd.setHours(0, 0, 0, 0);

    const result: TrendDataItem[] = [];
    const currentWeekKey = (() => {
      const year = new Date().getFullYear();
      const isoWeek = getWeekNumber(new Date());
      return `${year}-W${String(isoWeek).padStart(2, '0')}`;
    })();

    // 收集本月内的周
    const weeks: { key: string; weekNum: number; monday: Date }[] = [];
    let currentMonday = new Date(monthStart);

    // 找到本月第一个周一
    const firstDay = currentMonday.getDay();
    if (firstDay !== 1) {
      const daysToMonday = firstDay === 0 ? 1 : (8 - firstDay);
      currentMonday.setDate(currentMonday.getDate() + daysToMonday);
    }

    // 如果第一个周一已经超出本月，说明本月1号之前就有周一，回退一周
    if (currentMonday >= monthEnd) {
      currentMonday.setDate(currentMonday.getDate() - 7);
    } else if (currentMonday > monthStart) {
      // 如果第一个周一不是1号，检查是否应该包含前一周
      const prevMonday = new Date(currentMonday);
      prevMonday.setDate(prevMonday.getDate() - 7);
      if (prevMonday >= monthStart || prevMonday.getMonth() === monthStart.getMonth() - 1) {
        currentMonday = prevMonday;
      }
    }

    // 收集所有周一在本月或之前，且不超过月末的周
    while (currentMonday < monthEnd) {
      const year = currentMonday.getFullYear();
      const weekNum = getWeekNumber(currentMonday);
      const weekKey = `${year}-W${String(weekNum).padStart(2, '0')}`;

      weeks.push({
        key: weekKey,
        weekNum,
        monday: new Date(currentMonday)
      });

      currentMonday.setDate(currentMonday.getDate() + 7);
    }

    // 生成趋势数据
    weeks.forEach((week, index) => {
      const stat = groupStats?.find(s => s.group_key === week.key);
      const score = stat?.score_total || 0;

      result.push({
        label: `第${index + 1}周`,
        score,
        percentage: 0,
        isCurrent: week.key === currentWeekKey
      });
    });

    // 归一化百分比
    const maxScore = Math.max(...result.map(r => r.score), 1);
    result.forEach(item => {
      item.percentage = (item.score / maxScore) * 100;
    });

    return result;
  };

  // 解析本季度趋势（各月）
  const parseQuarterTrend = (groupStats: { group_key: string; score_total: number }[]): TrendDataItem[] => {
    const quarter = getQuarterNumber(currentDate);
    const year = currentDate.getFullYear();
    const startMonth = (quarter - 1) * 3 + 1;

    const result: TrendDataItem[] = [];
    const currentMonth = new Date().getMonth() + 1;

    for (let i = 0; i < 3; i++) {
      const month = startMonth + i;
      const monthKey = `${year}-${String(month).padStart(2, '0')}`;

      const stat = groupStats?.find(s => s.group_key === monthKey);
      const score = stat?.score_total || 0;

      result.push({
        label: `${month}月`,
        score,
        percentage: 0,
        isCurrent: month === currentMonth && year === new Date().getFullYear()
      });
    }

    // 归一化百分比
    const maxScore = Math.max(...result.map(r => r.score), 1);
    result.forEach(item => {
      item.percentage = (item.score / maxScore) * 100;
    });

    return result;
  };

  // 解析本年趋势（各季度）
  const parseYearTrend = (groupStats: { group_key: string; score_total: number }[]): TrendDataItem[] => {
    const year = currentDate.getFullYear();
    const result: TrendDataItem[] = [];
    const currentQuarter = getQuarterNumber(new Date());

    for (let i = 1; i <= 4; i++) {
      const quarterKey = `${year}-Q${i}`;

      const stat = groupStats?.find(s => s.group_key === quarterKey);
      const score = stat?.score_total || 0;

      result.push({
        label: `Q${i}`,
        score,
        percentage: 0,
        isCurrent: i === currentQuarter && year === new Date().getFullYear()
      });
    }

    // 归一化百分比
    const maxScore = Math.max(...result.map(r => r.score), 1);
    result.forEach(item => {
      item.percentage = (item.score / maxScore) * 100;
    });

    return result;
  };

  const loadPlanData = async () => {
    setLoading(true);
    // 重置状态，确保从干净状态开始
    setJournals([]);

    try {
      const dates = getPeriodDates(currentPeriod, currentDate);

      // 获取计划数据
      try {
        const plan = await planService.getPlan({
          period_type: currentPeriod,
          ...dates
        });
        setPlanData(plan);

        // 更新任务列表
        if (plan.tasks) {
          // 统计任务状态
          const statusStats = {
            notStarted: 0,
            inProgress: 0,
            completed: 0,
            cancelled: 0
          };

          plan.tasks.forEach(task => {
            switch (task.status) {
              case TaskStatus.NotStarted:
                statusStats.notStarted++;
                break;
              case TaskStatus.InProgress:
                statusStats.inProgress++;
                break;
              case TaskStatus.Completed:
                statusStats.completed++;
                break;
              case TaskStatus.Cancelled:
                statusStats.cancelled++;
                break;
            }
          });

          setStats(prev => ({
            ...prev,
            taskStats: statusStats,
            todayScore: currentPeriod === 'day' ? plan.score_total : prev.todayScore,
            weekScore: currentPeriod === 'week' ? plan.score_total : prev.weekScore,
            monthScore: currentPeriod === 'month' ? plan.score_total : prev.monthScore
          }));
        }

        // 更新日志列表（总是设置，即使为空）
        setJournals(plan.journals || []);
      } catch (planError) {
        console.warn('Failed to load plan data, using empty data:', planError);
        // Plan API失败时，确保状态被清空
        setJournals([]);
      }

      // 获取趋势数据（独立的API调用）
      try {
        const trendGroupBy = getTrendGroupBy();
        const trendPeriod = getTrendPeriod();

        const groupStats = await planService.getPlanStats({
          group_by: trendGroupBy,
          ...trendPeriod
        });

        // 根据当前周期类型解析趋势数据
        let trendData: TrendDataItem[] = [];
        switch (currentPeriod) {
          case 'day':
          case 'week':
            trendData = parseWeekTrend(groupStats || []);
            break;
          case 'month':
            trendData = parseMonthTrend(groupStats || []);
            break;
          case 'quarter':
            trendData = parseQuarterTrend(groupStats || []);
            break;
          case 'year':
            trendData = parseYearTrend(groupStats || []);
            break;
        }

        setStats(prev => ({
          ...prev,
          trendData
        }));
      } catch (trendError) {
        console.warn('Failed to load trend data, using empty data:', trendError);
        setStats(prev => ({
          ...prev,
          trendData: []
        }));
      }

    } catch (error) {
      console.error('Failed to load data:', error);
      // 确保即使外层异常也设置空数据
      setJournals([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadPlanData();
  }, [currentPeriod, currentDate]);

  return (
    <div style={{ minHeight: '100vh', background: 'var(--bg-primary)' }}>
      {/* 顶部导航栏 */}
      <header className="navbar">
        <div className="navbar-brand">
          <span className="logo">🌙</span>
          <h1>Luna Dial</h1>
        </div>

        {/* 周期切换器 */}
        <div className="period-switcher">
          {(['day', 'week', 'month', 'quarter', 'year'] as PeriodType[]).map(period => (
            <button
              key={period}
              onClick={() => setCurrentPeriod(period)}
              className={`period-btn ${currentPeriod === period ? 'active' : ''}`}
            >
              {period === 'day' && '日'}
              {period === 'week' && '周'}
              {period === 'month' && '月'}
              {period === 'quarter' && '季'}
              {period === 'year' && '年'}
            </button>
          ))}
        </div>

        {/* 时间导航器 */}
        <TimeNavigator
          currentPeriod={currentPeriod}
          currentDate={currentDate}
          onDateChange={handleDateChange}
          onNavigate={handleNavigate}
        />

        <div className="user-info">
          <span className="user-name">{user?.name || user?.username}</span>
          <button className="btn-profile" onClick={() => navigate('/profile')}>
            设置
          </button>
          <button className="btn-logout" onClick={handleLogout}>
            登出
          </button>
        </div>
      </header>

      {/* 主内容区域 */}
      <main className={`dashboard-container ${currentPeriod === 'year' ? 'no-parent-panel' : 'has-parent-panel'}`}>
        {/* 左侧：上级任务列表（年任务时隐藏） */}
        {currentPeriod !== 'year' && (
          <section className="parent-task-panel">
            <div className="panel-header">
              <h2>上级任务</h2>
              <button className="btn-view-tree" onClick={() => navigate('/tasks/tree')}>
                任务树
              </button>
            </div>

            <ParentTaskList
              currentPeriod={currentPeriod}
              currentDate={currentDate}
              onTaskStatusChange={handleTaskStatusChange}
              onTaskClick={handleTaskClick}
              onTaskEdit={handleEditTask}
              onTaskDelete={handleDeleteTask}
            />
          </section>
        )}

        {/* 中间：当前周期概览 */}
        <section className="overview-panel">
          {/* 今日任务卡片 */}
          <div className="focus-card">
            <h3>{getPeriodLabel(currentPeriod, currentDate)}任务</h3>
            <div className="current-date">
              {getCurrentDateString(currentDate)}
              {formatPeriodRange(currentDate, currentPeriod) && (
                <span className="period-range"> ({formatPeriodRange(currentDate, currentPeriod)})</span>
              )}
            </div>

            {loading ? (
              <div className="loading">加载中...</div>
            ) : (
              <div className="today-tasks">
                {planData?.tasks?.filter(task =>
                  task.task_type === (currentPeriod === 'day' ? 0 :
                                      currentPeriod === 'week' ? 1 :
                                      currentPeriod === 'month' ? 2 :
                                      currentPeriod === 'quarter' ? 3 : 4)
                ).map(task => (
                  <div key={task.id} className="daily-task">
                    <div className="task-header">
                      <div className="task-info">
                        <span className="task-icon">{task.icon || '📝'}</span>
                        <span className="task-text">{task.title}</span>
                        <span className={`priority-badge ${getPriorityClass(task.priority)}`}>
                          {getPriorityLabel(task.priority)}
                        </span>
                      </div>
                    </div>
                    <div className="task-controls">
                      <div className="control-item">
                        <label className="control-label">状态</label>
                        <select
                          className="task-status-select"
                          value={task.status}
                          onChange={(e) => handleTaskStatusChange(task.id, Number(e.target.value) as TaskStatus)}
                        >
                          <option value={TaskStatus.NotStarted}>未开始</option>
                          <option value={TaskStatus.InProgress}>进行中</option>
                          <option value={TaskStatus.Completed}>已完成</option>
                          <option value={TaskStatus.Cancelled}>已取消</option>
                        </select>
                      </div>
                      <div className="control-item">
                        <label className="control-label">努力程度</label>
                        <div className="score-display-container">
                          <span className={`score-value ${task.status === TaskStatus.NotStarted || task.task_type !== 0 ? 'disabled' : ''}`}>
                            {task.status === TaskStatus.NotStarted || task.task_type !== 0 ? '-' : `${task.score}/10`}
                          </span>
                          <button
                            className="btn-edit-score"
                            onClick={() => handleOpenScoreDialog(task)}
                            disabled={task.status === TaskStatus.NotStarted || task.task_type !== 0}
                            title={task.task_type !== 0 ? "仅日任务支持评分" : "修改努力程度"}
                          >
                            ✏️
                          </button>
                        </div>
                      </div>
                      <div className="control-item task-actions">
                        <button
                          className="btn-delete-task"
                          onClick={() => {
                            if (window.confirm('确定要删除这个任务吗？')) {
                              handleDeleteTask(task.id);
                            }
                          }}
                          title="删除任务"
                        >
                          🗑️
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}

            <button className="btn-add-task" onClick={handleCreateTask}>
              + 添加{getPeriodLabel(currentPeriod, currentDate)}任务
            </button>
          </div>

          {/* 当前周期日志 */}
          <div className="journal-card">
            <div className="journal-header">
              <h3>{getPeriodLabel(currentPeriod, currentDate)}日志</h3>
              <button className="btn-new-journal" onClick={handleCreateJournal}>
                + 新建日志
              </button>
            </div>

            <div className="journal-list">
              {loading ? (
                <div className="loading">加载中...</div>
              ) : journals.length > 0 ? (
                journals.map(journal => (
                  <div key={journal.id} className="journal-entry">
                    <div className="journal-content">
                      <span className="journal-icon">{journal.icon || '📝'}</span>
                      <div className="journal-info">
                        <h4>{journal.title}</h4>
                        <span className="journal-time">
                          {new Date(journal.created_at).toLocaleTimeString('zh-CN', {
                            hour: '2-digit',
                            minute: '2-digit'
                          })}
                        </span>
                      </div>
                    </div>
                    <div className="journal-actions">
                      <button onClick={() => handleViewJournal(journal)}>查看</button>
                      <button onClick={() => handleEditJournal(journal)}>编辑</button>
                      <button onClick={() => handleDeleteJournal(journal.id)}>删除</button>
                    </div>
                  </div>
                ))
              ) : (
                <div className="empty-state">暂无日志</div>
              )}
            </div>
          </div>
        </section>

        {/* 右侧：成长数据 */}
        <section className="stats-panel">
          {/* 努力程度统计 */}
          <div className="score-card">
            <h3>努力程度统计</h3>
            <div className="score-display">
              <div className="score-item">
                <span className="score-label">今日努力总分</span>
                <span className="score-value">{stats.todayScore}</span>
              </div>
              <div className="score-item">
                <span className="score-label">本周累计</span>
                <span className="score-value">{stats.weekScore}</span>
              </div>
              <div className="score-item">
                <span className="score-label">本月累计</span>
                <span className="score-value">{stats.monthScore}</span>
              </div>
            </div>
          </div>

          {/* 努力趋势图表 */}
          <div className="progress-card">
            <h3>
              {currentPeriod === 'day' || currentPeriod === 'week' ? '本周' : ''}
              {currentPeriod === 'month' ? '本月' : ''}
              {currentPeriod === 'quarter' ? '本季' : ''}
              {currentPeriod === 'year' ? '本年' : ''}
              努力趋势
            </h3>
            <div className="trend-chart-container">
              {stats.trendData.length > 0 ? (() => {
                const width = 280;
                const height = 160;
                const padding = { top: 20, right: 10, bottom: 30, left: 35 };
                const chartWidth = width - padding.left - padding.right;
                const chartHeight = height - padding.top - padding.bottom;

                // 计算最大值和刻度
                const maxScore = Math.max(...stats.trendData.map(d => d.score), 1);
                const yMax = Math.ceil(maxScore / 5) * 5 || 10; // 向上取整到5的倍数
                const yTicks = 5;
                const yStep = yMax / yTicks;

                // 计算数据点坐标
                const points = stats.trendData.map((item, i) => {
                  const x = padding.left + (chartWidth / (stats.trendData.length - 1 || 1)) * i;
                  const y = padding.top + chartHeight - (item.score / yMax) * chartHeight;
                  return { x, y, ...item };
                });

                // 生成折线路径
                const linePath = points.map((p, i) =>
                  `${i === 0 ? 'M' : 'L'} ${p.x},${p.y}`
                ).join(' ');

                // 生成填充区域路径
                const areaPath = `${linePath} L ${points[points.length - 1].x},${height - padding.bottom} L ${padding.left},${height - padding.bottom} Z`;

                return (
                  <svg width={width} height={height} className="trend-chart">
                    {/* 背景网格线 */}
                    {Array.from({ length: yTicks + 1 }).map((_, i) => {
                      const y = padding.top + (chartHeight / yTicks) * i;
                      return (
                        <line
                          key={`grid-${i}`}
                          x1={padding.left}
                          y1={y}
                          x2={width - padding.right}
                          y2={y}
                          className="grid-line"
                        />
                      );
                    })}

                    {/* Y轴刻度标签 */}
                    {Array.from({ length: yTicks + 1 }).map((_, i) => {
                      const value = yMax - (yStep * i);
                      const y = padding.top + (chartHeight / yTicks) * i;
                      return (
                        <text
                          key={`y-label-${i}`}
                          x={padding.left - 8}
                          y={y + 4}
                          className="axis-label"
                          textAnchor="end"
                        >
                          {Math.round(value)}
                        </text>
                      );
                    })}

                    {/* 填充区域 */}
                    <path
                      d={areaPath}
                      className="area-fill"
                    />

                    {/* 折线 */}
                    <path
                      d={linePath}
                      className="trend-line"
                    />

                    {/* 数据点 */}
                    {points.map((point, i) => (
                      <g key={`point-${i}`}>
                        {/* 数据点圆圈 */}
                        <circle
                          cx={point.x}
                          cy={point.y}
                          r={point.isCurrent ? 5 : 4}
                          className={`data-point ${point.isCurrent ? 'current' : ''}`}
                        />
                        {/* 分数标签 */}
                        {point.score > 0 && (
                          <text
                            x={point.x}
                            y={point.y - 10}
                            className="score-label"
                            textAnchor="middle"
                          >
                            {point.score}
                          </text>
                        )}
                      </g>
                    ))}

                    {/* X轴标签 */}
                    {points.map((point, i) => (
                      <text
                        key={`x-label-${i}`}
                        x={point.x}
                        y={height - padding.bottom + 20}
                        className={`x-axis-label ${point.isCurrent ? 'current' : ''}`}
                        textAnchor="middle"
                      >
                        {point.label}
                      </text>
                    ))}
                  </svg>
                );
              })() : (
                <div className="empty-chart">暂无数据</div>
              )}
            </div>
          </div>

          {/* 任务状态分布 */}
          <div className="status-card">
            <h3>任务状态</h3>
            <div className="status-stats">
              <div className="status-item">
                <span className="status-dot status-notstarted"></span>
                <span className="status-label">未开始</span>
                <span className="status-count">{stats.taskStats.notStarted}</span>
              </div>
              <div className="status-item">
                <span className="status-dot status-inprogress"></span>
                <span className="status-label">进行中</span>
                <span className="status-count">{stats.taskStats.inProgress}</span>
              </div>
              <div className="status-item">
                <span className="status-dot status-completed"></span>
                <span className="status-label">已完成</span>
                <span className="status-count">{stats.taskStats.completed}</span>
              </div>
              <div className="status-item">
                <span className="status-dot status-cancelled"></span>
                <span className="status-label">已取消</span>
                <span className="status-count">{stats.taskStats.cancelled}</span>
              </div>
            </div>
          </div>
        </section>
      </main>

      {/* 快速操作栏 */}
      <footer className="quick-actions">
        <button className="action-btn primary" onClick={handleCreateTask}>
          <span className="action-icon">➕</span>
          <span>创建任务</span>
        </button>
        <button className="action-btn" onClick={handleCreateJournal}>
          <span className="action-icon">📝</span>
          <span>写日志</span>
        </button>
      </footer>

      {/* 对话框组件 */}
      {showTaskDialog && (
        <TaskEditDialog
          task={editingTask}
          onClose={() => {
            setShowTaskDialog(false);
            setParentTaskIdForCreation(undefined);
          }}
          onSuccess={() => {
            setShowTaskDialog(false);
            setParentTaskIdForCreation(undefined);
            loadPlanData();
          }}
          currentPeriod={currentPeriod}
          currentDate={currentDate}
          parentTaskId={parentTaskIdForCreation}
        />
      )}

      {showViewTaskDialog && viewingTask && (
        <TaskViewDialog
          task={viewingTask}
          onClose={() => setShowViewTaskDialog(false)}
          onEdit={(task) => handleEditTask(task)}
          onDelete={(taskId) => handleDeleteTask(taskId)}
          onScoreUpdate={handleTaskScoreChange}
          onCreateSubtask={handleCreateSubtask}
        />
      )}

      {showJournalDialog && (
        <JournalEditDialog
          journal={editingJournal}
          onClose={() => setShowJournalDialog(false)}
          onSuccess={() => {
            setShowJournalDialog(false);
            loadPlanData();
          }}
          currentPeriod={currentPeriod}
          currentDate={currentDate}
        />
      )}

      {showViewJournalDialog && viewingJournal && (
        <JournalViewDialog
          journal={viewingJournal}
          onClose={() => setShowViewJournalDialog(false)}
          onEdit={(journal) => handleEditJournal(journal)}
        />
      )}

      {showScoreDialog && editingScoreTask && (
        <ScoreEditDialog
          taskTitle={editingScoreTask.title}
          currentScore={editingScoreTask.score}
          onClose={() => {
            setShowScoreDialog(false);
            setEditingScoreTask(null);
          }}
          onSave={async (score) => {
            await handleTaskScoreChange(editingScoreTask.id, score);
            setShowScoreDialog(false);
            setEditingScoreTask(null);
          }}
        />
      )}
    </div>
  );
};

export default Dashboard;
