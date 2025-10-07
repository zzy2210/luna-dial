import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import useAuthStore from '../store/auth';
import TaskTree from '../components/TaskTree';
import taskService from '../services/task';
import { Task, TaskStatus, PeriodType } from '../types';
import '../styles/dashboard.css';

const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const { user, logout } = useAuthStore();
  const [currentPeriod, setCurrentPeriod] = useState<PeriodType>('day');
  const [tasks, setTasks] = useState<Task[]>([]);

  // 模拟的今日任务数据
  const [todayTasks] = useState([
    { id: '1', title: '完成项目文档', icon: '📝', status: TaskStatus.Completed, score: 8 },
    { id: '2', title: '代码审查', icon: '👨‍💻', status: TaskStatus.InProgress, score: 6 },
    { id: '3', title: '系统优化', icon: '🔧', status: TaskStatus.NotStarted, score: 0 },
  ]);

  // 模拟的统计数据
  const [stats] = useState({
    todayScore: 14,
    weekScore: 68,
    monthScore: 245,
    weekProgress: [60, 80, 45, 90, 70, 85, 50],
    taskStats: {
      notStarted: 3,
      inProgress: 5,
      completed: 12,
      cancelled: 1
    }
  });

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  const handleTaskStatusChange = (taskId: string, status: TaskStatus) => {
    console.log('Task status changed:', taskId, status);
    // TODO: 调用API更新任务状态
  };

  const handleTaskClick = (task: Task) => {
    console.log('Task clicked:', task);
    // TODO: 显示任务详情
  };

  useEffect(() => {
    // 加载初始数据
    const loadData = async () => {
      try {
        const response = await taskService.getTaskTree();
        setTasks(response.items);
      } catch (error) {
        console.error('Failed to load tasks:', error);
      }
    };

    loadData();
  }, []);

  const getCurrentDateString = () => {
    const date = new Date();
    const year = date.getFullYear();
    const month = date.getMonth() + 1;
    const day = date.getDate();
    const weekDay = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六'][date.getDay()];
    return `${year}年${month}月${day}日 ${weekDay}`;
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
            <button className="btn-create-task">+ 新建任务</button>
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
            <h3>今日任务</h3>
            <div className="current-date">{getCurrentDateString()}</div>

            <div className="today-tasks">
              {todayTasks.map(task => (
                <div key={task.id} className="daily-task">
                  <div className="task-info">
                    <span className="task-icon">{task.icon}</span>
                    <span className="task-text">{task.title}</span>
                  </div>
                  <div className="task-controls">
                    <select className="task-status-select" defaultValue={task.status}>
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
                        defaultValue={task.score}
                        disabled={task.status === TaskStatus.NotStarted}
                      />
                      <span className="score-display">
                        {task.status === TaskStatus.NotStarted ? '-' : task.score}
                      </span>
                    </div>
                  </div>
                </div>
              ))}
            </div>

            <button className="btn-add-task">+ 添加今日任务</button>
          </div>

          {/* 当前周期日志 */}
          <div className="journal-card">
            <div className="journal-header">
              <h3>今日日志</h3>
              <button className="btn-new-journal">+ 新建日志</button>
            </div>

            <div className="journal-list">
              <div className="journal-entry">
                <div className="journal-content">
                  <span className="journal-icon">🌅</span>
                  <div className="journal-info">
                    <h4>早晨计划</h4>
                    <span className="journal-time">09:00</span>
                  </div>
                </div>
                <div className="journal-actions">
                  <button>查看</button>
                  <button>编辑</button>
                  <button>删除</button>
                </div>
              </div>

              <div className="journal-entry">
                <div className="journal-content">
                  <span className="journal-icon">🌙</span>
                  <div className="journal-info">
                    <h4>晚间总结</h4>
                    <span className="journal-time">22:30</span>
                  </div>
                </div>
                <div className="journal-actions">
                  <button>查看</button>
                  <button>编辑</button>
                  <button>删除</button>
                </div>
              </div>
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
                <div key={day} className={`chart-bar ${index === 6 ? 'today' : ''}`}>
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
        <button className="action-btn primary">
          <span className="action-icon">➕</span>
          <span>创建任务</span>
        </button>
        <button className="action-btn">
          <span className="action-icon">📝</span>
          <span>写日志</span>
        </button>
        <button className="action-btn">
          <span className="action-icon">✅</span>
          <span>每日打卡</span>
        </button>
      </footer>
    </div>
  );
};

export default Dashboard;