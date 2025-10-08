import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import useAuthStore from '../store/auth';
import TaskTree from '../components/TaskTree';
import taskService from '../services/task';
import journalService from '../services/journal';
import planService from '../services/plan';
import { Task, TaskStatus, PeriodType, Journal, PlanResponse } from '../types';
import TaskCreateDialog from '../components/TaskCreateDialog';
import JournalEditDialog from '../components/JournalEditDialog';
import JournalViewDialog from '../components/JournalViewDialog';
import '../styles/dashboard.css';

const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const { user, logout } = useAuthStore();
  const [currentPeriod, setCurrentPeriod] = useState<PeriodType>('day');
  const [tasks, setTasks] = useState<Task[]>([]);
  const [journals, setJournals] = useState<Journal[]>([]);
  const [planData, setPlanData] = useState<PlanResponse | null>(null);
  const [loading, setLoading] = useState(true);

  // 对话框状态
  const [showTaskDialog, setShowTaskDialog] = useState(false);
  const [showJournalDialog, setShowJournalDialog] = useState(false);
  const [showViewJournalDialog, setShowViewJournalDialog] = useState(false);
  const [editingJournal, setEditingJournal] = useState<Journal | null>(null);
  const [viewingJournal, setViewingJournal] = useState<Journal | null>(null);

  // 统计数据
  const [stats, setStats] = useState({
    todayScore: 0,
    weekScore: 0,
    monthScore: 0,
    weekProgress: [0, 0, 0, 0, 0, 0, 0],
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
      const statusMap = {
        [TaskStatus.NotStarted]: 'not_started',
        [TaskStatus.InProgress]: 'in_progress',
        [TaskStatus.Completed]: 'completed',
        [TaskStatus.Cancelled]: 'cancelled'
      };

      await taskService.updateTask(taskId, {
        status: statusMap[status] as any
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
    console.log('Task clicked:', task);
    // TODO: 显示任务详情
  };

  const handleCreateTask = () => {
    setShowTaskDialog(true);
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

  const getPeriodDates = (period: PeriodType) => {
    const today = new Date();
    const startDate = new Date();
    const endDate = new Date();

    switch (period) {
      case 'day':
        // 今天 [today 00:00, tomorrow 00:00)
        startDate.setHours(0, 0, 0, 0);
        endDate.setDate(endDate.getDate() + 1);
        endDate.setHours(0, 0, 0, 0);
        break;
      case 'week':
        // 本周 ISO Week [Monday 00:00, Next Monday 00:00)
        const day = today.getDay();
        const diff = today.getDate() - day + (day === 0 ? -6 : 1);
        startDate.setDate(diff);
        startDate.setHours(0, 0, 0, 0);
        endDate.setDate(startDate.getDate() + 7);
        endDate.setHours(0, 0, 0, 0);
        break;
      case 'month':
        // 本月 [1st 00:00, Next Month 1st 00:00)
        startDate.setDate(1);
        startDate.setHours(0, 0, 0, 0);
        endDate.setMonth(endDate.getMonth() + 1, 1);
        endDate.setHours(0, 0, 0, 0);
        break;
      case 'quarter':
        // 本季度 [Quarter Start 00:00, Next Quarter Start 00:00)
        const quarter = Math.floor(today.getMonth() / 3);
        startDate.setMonth(quarter * 3, 1);
        startDate.setHours(0, 0, 0, 0);
        endDate.setMonth((quarter + 1) * 3, 1);
        endDate.setHours(0, 0, 0, 0);
        break;
      case 'year':
        // 本年 [Jan 1 00:00, Next Year Jan 1 00:00)
        startDate.setMonth(0, 1);
        startDate.setHours(0, 0, 0, 0);
        endDate.setFullYear(endDate.getFullYear() + 1, 0, 1);
        endDate.setHours(0, 0, 0, 0);
        break;
    }

    return {
      start_date: startDate.toISOString().split('T')[0],
      end_date: endDate.toISOString().split('T')[0]
    };
  };

  const loadPlanData = async () => {
    setLoading(true);
    try {
      const dates = getPeriodDates(currentPeriod);

      // 获取计划数据
      try {
        const plan = await planService.getPlan({
          period_type: currentPeriod,
          ...dates
        });
        setPlanData(plan);

        // 更新任务列表
        if (plan.tasks) {
          setTasks(plan.tasks);

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

        // 更新日志列表
        if (plan.journals) {
          setJournals(plan.journals);
        }
      } catch (planError) {
        console.warn('Failed to load plan data, using empty data:', planError);
        // 使用空数据，防止页面崩溃
      }

      // 获取任务树数据（用于左侧面板）
      try {
        const treeResponse = await taskService.getTaskTree();
        setTasks(treeResponse.items);
      } catch (treeError) {
        console.warn('Failed to load task tree, using empty data:', treeError);
        setTasks([]);
      }

    } catch (error) {
      console.error('Failed to load data:', error);
      // 确保即使API失败也设置空数据
      setTasks([]);
      setJournals([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadPlanData();
  }, [currentPeriod]);

  const getCurrentDateString = () => {
    const date = new Date();
    const year = date.getFullYear();
    const month = date.getMonth() + 1;
    const day = date.getDate();
    const weekDay = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六'][date.getDay()];
    return `${year}年${month}月${day}日 ${weekDay}`;
  };

  const getPeriodLabel = (period: PeriodType) => {
    const labels = {
      day: '今日',
      week: '本周',
      month: '本月',
      quarter: '本季度',
      year: '本年'
    };
    return labels[period];
  };

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
      <main className="dashboard-container">
        {/* 左侧：任务管理区 */}
        <section className="task-panel">
          <div className="panel-header">
            <h2>任务树</h2>
            <button className="btn-create-task" onClick={handleCreateTask}>
              + 新建任务
            </button>
          </div>

          <TaskTree
            tasks={tasks}
            onTaskStatusChange={handleTaskStatusChange}
            onTaskClick={handleTaskClick}
          />
        </section>

        {/* 中间：当前周期概览 */}
        <section className="overview-panel">
          {/* 今日任务卡片 */}
          <div className="focus-card">
            <h3>{getPeriodLabel(currentPeriod)}任务</h3>
            <div className="current-date">{getCurrentDateString()}</div>

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
                    <div className="task-info">
                      <span className="task-icon">{task.icon || '📝'}</span>
                      <span className="task-text">{task.title}</span>
                    </div>
                    <div className="task-controls">
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
                      <div className={`score-control ${task.status === TaskStatus.NotStarted ? 'disabled' : ''}`}>
                        <label>努力程度:</label>
                        <input
                          type="number"
                          className="score-input"
                          min="0"
                          max="10"
                          value={task.score}
                          disabled={task.status === TaskStatus.NotStarted}
                          onChange={(e) => handleTaskScoreChange(task.id, Number(e.target.value))}
                        />
                        <span className="score-display">
                          {task.status === TaskStatus.NotStarted ? '-' : task.score}
                        </span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}

            <button className="btn-add-task" onClick={handleCreateTask}>
              + 添加{getPeriodLabel(currentPeriod)}任务
            </button>
          </div>

          {/* 当前周期日志 */}
          <div className="journal-card">
            <div className="journal-header">
              <h3>{getPeriodLabel(currentPeriod)}日志</h3>
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
            <h3>本周努力趋势</h3>
            <div className="progress-chart">
              {['一', '二', '三', '四', '五', '六', '日'].map((day, index) => (
                <div key={day} className={`chart-bar ${index === new Date().getDay() - 1 ? 'today' : ''}`}>
                  <div className="bar-fill" style={{ height: `${stats.weekProgress[index]}%` }}></div>
                  <span className="bar-label">{day}</span>
                </div>
              ))}
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
        <button className="action-btn" onClick={() => console.log('Daily check-in')}>
          <span className="action-icon">✅</span>
          <span>每日打卡</span>
        </button>
      </footer>

      {/* 对话框组件 */}
      {showTaskDialog && (
        <TaskCreateDialog
          onClose={() => setShowTaskDialog(false)}
          onSuccess={() => {
            setShowTaskDialog(false);
            loadPlanData();
          }}
          currentPeriod={currentPeriod}
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
        />
      )}

      {showViewJournalDialog && viewingJournal && (
        <JournalViewDialog
          journal={viewingJournal}
          onClose={() => setShowViewJournalDialog(false)}
          onEdit={(journal) => handleEditJournal(journal)}
        />
      )}
    </div>
  );
};

export default Dashboard;