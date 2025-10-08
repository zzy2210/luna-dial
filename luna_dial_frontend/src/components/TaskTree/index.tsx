import React, { useState } from 'react';
import { Task, TaskStatus } from '../../types';

interface TaskTreeProps {
  tasks: Task[];
  onTaskStatusChange?: (taskId: string, status: TaskStatus) => void;
  onTaskClick?: (task: Task) => void;
}

interface TaskNodeProps {
  task: Task;
  level: number;
  onStatusChange?: (taskId: string, status: TaskStatus) => void;
  onClick?: (task: Task) => void;
}

const TaskNode: React.FC<TaskNodeProps> = ({ task, level, onStatusChange, onClick }) => {
  const [isExpanded, setIsExpanded] = useState(false);

  const getTaskTypeClass = (level: number) => {
    const classes = ['task-year', 'task-quarter', 'task-month', 'task-week', 'task-day'];
    return classes[Math.min(level, classes.length - 1)];
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

  const handleToggle = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (task.has_children) {
      setIsExpanded(!isExpanded);
    }
  };

  const handleStatusChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    e.stopPropagation();
    if (onStatusChange) {
      onStatusChange(task.id, Number(e.target.value) as TaskStatus);
    }
  };

  return (
    <>
      <div
        className={`task-node ${getTaskTypeClass(level)}`}
        onClick={() => onClick && onClick(task)}
      >
        {task.has_children ? (
          <span
            className={`task-toggle ${isExpanded ? 'expanded' : ''}`}
            onClick={handleToggle}
          >
            ▶
          </span>
        ) : (
          <span style={{ width: '20px', display: 'inline-block' }}></span>
        )}

        <span className="task-icon">{task.icon || '📋'}</span>
        <span className="task-title">{task.title}</span>

        <select
          className="task-status-mini"
          value={task.status}
          onChange={handleStatusChange}
          onClick={(e) => e.stopPropagation()}
        >
          <option value={TaskStatus.NotStarted}>{getStatusText(TaskStatus.NotStarted)}</option>
          <option value={TaskStatus.InProgress}>{getStatusText(TaskStatus.InProgress)}</option>
          <option value={TaskStatus.Completed}>{getStatusText(TaskStatus.Completed)}</option>
          <option value={TaskStatus.Cancelled}>{getStatusText(TaskStatus.Cancelled)}</option>
        </select>
      </div>

      {isExpanded && task.children && (
        <div className="task-children">
          {task.children.map(child => (
            <TaskNode
              key={child.id}
              task={child}
              level={level + 1}
              onStatusChange={onStatusChange}
              onClick={onClick}
            />
          ))}
        </div>
      )}
    </>
  );
};

const TaskTree: React.FC<TaskTreeProps> = ({ tasks, onTaskStatusChange, onTaskClick }) => {
  return (
    <div className="task-tree">
      {tasks.length > 0 ? (
        tasks.map(task => (
          <TaskNode
            key={task.id}
            task={task}
            level={0}
            onStatusChange={onTaskStatusChange}
            onClick={onTaskClick}
          />
        ))
      ) : (
        <div style={{
          textAlign: 'center',
          padding: '2rem',
          color: 'var(--text-tertiary)'
        }}>
          暂无任务
        </div>
      )}
    </div>
  );
};

export default TaskTree;