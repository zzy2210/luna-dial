import React, { useState, useEffect } from 'react';
import { Task, TaskStatus, TaskPriority, PeriodType } from '../../types';
import planService from '../../services/plan';
import { getPeriodDates, getParentPeriodConfig, getPeriodLabel } from '../../utils/dateUtils';
import '../../styles/parent-task-list.css';

interface ParentTaskListProps {
  currentPeriod: PeriodType;
  currentDate: Date;
  onTaskStatusChange?: (taskId: string, status: TaskStatus) => void;
  onTaskClick?: (task: Task) => void;
  onTaskEdit?: (task: Task) => void;
  onTaskDelete?: (taskId: string) => void;
}

const ParentTaskList: React.FC<ParentTaskListProps> = ({
  currentPeriod,
  currentDate,
  onTaskStatusChange,
  onTaskClick,
  onTaskEdit,
  onTaskDelete
}) => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [statusFilter, setStatusFilter] = useState<TaskStatus[]>([]);

  const parentConfig = getParentPeriodConfig(currentPeriod);

  useEffect(() => {
    if (!parentConfig) {
      setTasks([]);
      setLoading(false);
      return;
    }

    loadParentTasks();
  }, [currentPeriod, currentDate]);

  const loadParentTasks = async () => {
    if (!parentConfig) return;

    setLoading(true);
    try {
      // 计算上级周期的日期范围
      const dates = getPeriodDates(parentConfig.parent, currentDate);

      // 获取上级周期的计划数据
      const plan = await planService.getPlan({
        period_type: parentConfig.parent,
        ...dates
      });

      // 筛选出对应类型的任务
      const filteredTasks = (plan.tasks || []).filter(
        task => task.task_type === parentConfig.taskType
      );

      setTasks(filteredTasks);
    } catch (error) {
      console.error('Failed to load parent tasks:', error);
      setTasks([]);
    } finally {
      setLoading(false);
    }
  };

  const getStatusText = (status: TaskStatus) => {
    const statusMap = {
      [TaskStatus.NotStarted]: '未开始',
      [TaskStatus.InProgress]: '进行中',
      [TaskStatus.Completed]: '已完成',
      [TaskStatus.Cancelled]: '已取消',
    };
    return statusMap[status];
  };

  const getPriorityLabel = (priority: TaskPriority): string => {
    const labels = {
      [TaskPriority.Low]: '低',
      [TaskPriority.Medium]: '中',
      [TaskPriority.High]: '高',
      [TaskPriority.Urgent]: '紧急'
    };
    return labels[priority] || '中';
  };

  const getPriorityClass = (priority: TaskPriority): string => {
    const classes = {
      [TaskPriority.Low]: 'priority-low',
      [TaskPriority.Medium]: 'priority-medium',
      [TaskPriority.High]: 'priority-high',
      [TaskPriority.Urgent]: 'priority-urgent'
    };
    return classes[priority] || 'priority-medium';
  };

  const handleStatusChange = async (taskId: string, status: TaskStatus) => {
    if (onTaskStatusChange) {
      await onTaskStatusChange(taskId, status);
      // 刷新父任务列表
      loadParentTasks();
    }
  };

  const handleDelete = async (taskId: string, e?: React.MouseEvent) => {
    if (e) {
      e.stopPropagation();
    }
    if (window.confirm('确定要删除这个任务吗？')) {
      if (onTaskDelete) {
        await onTaskDelete(taskId);
        // 刷新父任务列表
        loadParentTasks();
      }
    }
  };

  const handleCardClick = (task: Task) => {
    if (onTaskClick) {
      onTaskClick(task);
    }
  };

  const handleEditClick = (task: Task, e: React.MouseEvent) => {
    e.stopPropagation();
    if (onTaskEdit) {
      onTaskEdit(task);
    }
  };

  const toggleStatusFilter = (status: TaskStatus) => {
    setStatusFilter(prev => {
      if (prev.includes(status)) {
        return prev.filter(s => s !== status);
      } else {
        return [...prev, status];
      }
    });
  };

  const filteredTasks = statusFilter.length > 0
    ? tasks.filter(task => statusFilter.includes(task.status))
    : tasks;

  // 按状态排序：未开始 → 进行中 → 已完成 → 已取消
  const sortedTasks = React.useMemo(() => {
    return [...filteredTasks].sort((a, b) => {
      if (a.status !== b.status) {
        return a.status - b.status;
      }
      return 0;
    });
  }, [filteredTasks]);

  if (!parentConfig) {
    return null;
  }

  return (
    <div className="parent-task-list">
      <div className="parent-task-header">
        <h3>{getPeriodLabel(parentConfig.parent, currentDate)}任务</h3>
        <div className="task-count">{tasks.length} 个任务</div>
      </div>

      {/* 状态筛选器 */}
      <div className="status-filter">
        <button
          className={`filter-btn ${statusFilter.includes(TaskStatus.NotStarted) ? 'active' : ''}`}
          onClick={() => toggleStatusFilter(TaskStatus.NotStarted)}
        >
          未开始
        </button>
        <button
          className={`filter-btn ${statusFilter.includes(TaskStatus.InProgress) ? 'active' : ''}`}
          onClick={() => toggleStatusFilter(TaskStatus.InProgress)}
        >
          进行中
        </button>
        <button
          className={`filter-btn ${statusFilter.includes(TaskStatus.Completed) ? 'active' : ''}`}
          onClick={() => toggleStatusFilter(TaskStatus.Completed)}
        >
          已完成
        </button>
        <button
          className={`filter-btn ${statusFilter.includes(TaskStatus.Cancelled) ? 'active' : ''}`}
          onClick={() => toggleStatusFilter(TaskStatus.Cancelled)}
        >
          已取消
        </button>
      </div>

      {/* 任务列表 */}
      <div className="parent-tasks-container">
        {loading ? (
          <div className="loading-state">加载中...</div>
        ) : sortedTasks.length > 0 ? (
          sortedTasks.map(task => (
            <div
              key={task.id}
              className="parent-task-item"
              onClick={() => handleCardClick(task)}
            >
              <div className="task-main">
                <span className="task-icon">{task.icon || '📋'}</span>
                <div className="task-info">
                  <div className="task-title">
                    {task.title}
                  </div>
                  <span className={`priority-badge ${getPriorityClass(task.priority)}`}>
                    {getPriorityLabel(task.priority)}
                  </span>
                </div>
              </div>

              <div className="task-controls">
                <select
                  className="task-status-select"
                  value={task.status}
                  onChange={(e) => handleStatusChange(task.id, Number(e.target.value) as TaskStatus)}
                  onClick={(e) => e.stopPropagation()}
                >
                  <option value={TaskStatus.NotStarted}>{getStatusText(TaskStatus.NotStarted)}</option>
                  <option value={TaskStatus.InProgress}>{getStatusText(TaskStatus.InProgress)}</option>
                  <option value={TaskStatus.Completed}>{getStatusText(TaskStatus.Completed)}</option>
                  <option value={TaskStatus.Cancelled}>{getStatusText(TaskStatus.Cancelled)}</option>
                </select>

                <div className="task-actions">
                  <button
                    className="btn-action btn-edit"
                    onClick={(e) => handleEditClick(task, e)}
                    title="编辑"
                  >
                    ✏️
                  </button>
                  <button
                    className="btn-action btn-delete"
                    onClick={(e) => handleDelete(task.id, e)}
                    title="删除"
                  >
                    🗑️
                  </button>
                </div>
              </div>
            </div>
          ))
        ) : (
          <div className="empty-state">
            {statusFilter.length > 0 ? '没有符合筛选条件的任务' : '暂无任务'}
          </div>
        )}
      </div>
    </div>
  );
};

export default ParentTaskList;
