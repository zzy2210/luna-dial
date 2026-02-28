import React, { useState } from 'react';
import { Task, TaskStatus, TaskPriority, TaskType } from '../types';
import { parseTagsFromJsonString } from '../utils/tags';
import { TASK_STATUS_LABELS } from '../constants/task';
import '../styles/dialog.css';
import '../styles/task-view-dialog.css';

interface TaskViewDialogProps {
  task: Task;
  onClose: () => void;
  onEdit: (task: Task) => void;
  onDelete: (taskId: string) => void;
  onScoreUpdate?: (taskId: string, score: number) => void;
  onCreateSubtask?: (parentTaskId: string) => void;
}

const TaskViewDialog: React.FC<TaskViewDialogProps> = ({
  task,
  onClose,
  onEdit,
  onDelete,
  onScoreUpdate,
  onCreateSubtask
}) => {
  const [currentScore, setCurrentScore] = useState<number>(task.score);
  const [editingScore, setEditingScore] = useState<number>(task.score);
  const [isEditingScore, setIsEditingScore] = useState(false);

  const taskTypeLabels = {
    0: '日任务',
    1: '周任务',
    2: '月任务',
    3: '季度任务',
    4: '年度任务'
  };

  const priorityLabels = {
    [TaskPriority.Low]: '低',
    [TaskPriority.Medium]: '中',
    [TaskPriority.High]: '高',
    [TaskPriority.Urgent]: '紧急'
  };

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  const handleDelete = () => {
    if (window.confirm('确定要删除这个任务吗？')) {
      onDelete(task.id);
    }
  };

  const handleScoreChange = (newScore: number) => {
    setEditingScore(newScore);
  };

  const handleScoreSave = () => {
    if (onScoreUpdate && editingScore !== currentScore) {
      onScoreUpdate(task.id, editingScore);
      setCurrentScore(editingScore);
    }
    setIsEditingScore(false);
  };

  const handleScoreCancel = () => {
    setEditingScore(currentScore);
    setIsEditingScore(false);
  };

  const canEditScore = (task.status === TaskStatus.InProgress || task.status === TaskStatus.Completed) && task.task_type === TaskType.Day;

  const getStatusClassName = (status: TaskStatus): string => {
    switch (status) {
      case TaskStatus.NotStarted:
        return 'not-started';
      case TaskStatus.InProgress:
        return 'in-progress';
      case TaskStatus.Completed:
        return 'completed';
      case TaskStatus.Cancelled:
        return 'cancelled';
      default:
        return 'not-started';
    }
  };

  const getPriorityClassName = (priority: TaskPriority): string => {
    switch (priority) {
      case TaskPriority.Low:
        return 'low';
      case TaskPriority.Medium:
        return 'medium';
      case TaskPriority.High:
        return 'high';
      case TaskPriority.Urgent:
        return 'urgent';
      default:
        return 'medium';
    }
  };

  const renderScoreStars = (score: number) => {
    const stars = [];
    const fullStars = Math.floor(score / 2);
    for (let i = 0; i < 5; i++) {
      stars.push(
        <span key={i} className={`score-star ${i < fullStars ? '' : 'empty'}`}>
          ★
        </span>
      );
    }
    return stars;
  };

  const tags = parseTagsFromJsonString(task.tags);

  return (
    <div className="dialog-overlay" onClick={onClose}>
      <div className="dialog-container large" onClick={e => e.stopPropagation()}>
        <div className="dialog-header">
          <h2>任务详情</h2>
          <button className="dialog-close" onClick={onClose}>×</button>
        </div>

        <div className="task-view-content">
          {/* 任务头部卡片 */}
          <div className="task-header-card">
            <span className="task-icon-large">{task.icon || '📋'}</span>
            <div className="task-header-info">
              <h3 className="task-title">{task.title}</h3>
              <div className="task-meta">
                <span className="task-type-label">{taskTypeLabels[task.task_type]}</span>
                <span className="meta-divider">·</span>
                <span className={`priority-badge ${getPriorityClassName(task.priority)}`}>
                  {priorityLabels[task.priority]}
                </span>
                <span className="meta-divider">·</span>
                <span className={`status-badge ${getStatusClassName(task.status)}`}>
                  {TASK_STATUS_LABELS[task.status]}
                </span>
              </div>
            </div>
          </div>

          {/* 任务时间范围 */}
          <div className="info-section">
            <div className="info-section-title">
              <span className="info-section-icon">📅</span>
              时间范围
            </div>
            <div className="date-range">
              <span className="date-icon">📆</span>
              <span className="date-text">
                {formatDate(task.period.start)} → {formatDate(task.period.end)}
              </span>
            </div>
          </div>

          {/* 努力评分 */}
          <div className="info-section">
            <div className="info-section-title">
              <span className="info-section-icon">💪</span>
              努力评分
              {canEditScore && !isEditingScore && onScoreUpdate && (
                <button
                  type="button"
                  className="btn-edit-score"
                  onClick={() => setIsEditingScore(true)}
                  title="编辑评分"
                >
                  ✏️
                </button>
              )}
            </div>
            {task.status === TaskStatus.NotStarted ? (
              <div className="info-card">
                任务尚未开始,暂无评分
              </div>
            ) : task.task_type !== TaskType.Day ? (
              <div className="info-card">
                仅日任务支持评分功能
              </div>
            ) : isEditingScore ? (
              <div className="score-editor">
                <div className="score-input-container">
                  <input
                    type="range"
                    min="0"
                    max="10"
                    step="1"
                    value={editingScore}
                    onChange={(e) => handleScoreChange(parseInt(e.target.value))}
                    className="score-slider"
                  />
                  <input
                    type="number"
                    min="0"
                    max="10"
                    value={editingScore}
                    onChange={(e) => handleScoreChange(Math.min(10, Math.max(0, parseInt(e.target.value) || 0)))}
                    className="score-number-input"
                  />
                </div>
                <div className="score-preview">
                  <div className="score-display">
                    <span className="score-value">{editingScore}</span>
                    <span className="score-max">/ 10</span>
                  </div>
                  <div className="score-bar">
                    <div
                      className="score-bar-fill"
                      style={{ width: `${(editingScore / 10) * 100}%` }}
                    />
                  </div>
                  <div className="score-stars">
                    {renderScoreStars(editingScore)}
                  </div>
                </div>
                <div className="score-actions">
                  <button
                    type="button"
                    className="btn-cancel-score"
                    onClick={handleScoreCancel}
                  >
                    取消
                  </button>
                  <button
                    type="button"
                    className="btn-save-score"
                    onClick={handleScoreSave}
                  >
                    保存
                  </button>
                </div>
              </div>
            ) : (
              <div className="score-container">
                <div className="score-display">
                  <span className="score-value">{currentScore}</span>
                  <span className="score-max">/ 10</span>
                </div>
                <div className="score-bar">
                  <div
                    className="score-bar-fill"
                    style={{ width: `${(currentScore / 10) * 100}%` }}
                  />
                </div>
                <div className="score-stars">
                  {renderScoreStars(currentScore)}
                </div>
              </div>
            )}
          </div>

          {/* 标签 */}
          {tags.length > 0 && (
            <div className="info-section">
              <div className="info-section-title">
                <span className="info-section-icon">🏷️</span>
                标签
              </div>
              <div className="tag-list">
                {tags.map(tag => (
                  <span key={tag} className="task-tag">
                    {tag}
                  </span>
                ))}
              </div>
            </div>
          )}

          {/* 任务层级信息 */}
          {(task.parent_id || task.has_children) && (
            <div className="info-section">
              <div className="info-section-title">
                <span className="info-section-icon">🌲</span>
                任务层级
              </div>
              <div className="hierarchy-card">
                {task.parent_id && (
                  <div className="hierarchy-item">
                    <span className="hierarchy-icon">↑</span>
                    <span className="hierarchy-label">父任务ID:</span>
                    <span className="hierarchy-value">{task.parent_id}</span>
                  </div>
                )}
                {task.has_children && (
                  <div className="hierarchy-item">
                    <span className="hierarchy-icon">↓</span>
                    <span className="hierarchy-label">子任务数量:</span>
                    <span className="hierarchy-value">{task.children_count}</span>
                  </div>
                )}
                <div className="hierarchy-item">
                  <span className="hierarchy-icon">📊</span>
                  <span className="hierarchy-label">树深度:</span>
                  <span className="hierarchy-value">{task.tree_depth}</span>
                </div>
              </div>
            </div>
          )}

          {/* 创建时间 */}
          <div className="timestamp">
            <span className="timestamp-icon">🕐</span>
            创建于 {new Date(task.created_at).toLocaleString('zh-CN')}
            {task.updated_at !== task.created_at && (
              <> · 最后更新 {new Date(task.updated_at).toLocaleString('zh-CN')}</>
            )}
          </div>
        </div>

        {/* 操作按钮 */}
        <div className="task-actions">
          <button type="button" onClick={handleDelete} className="btn-delete">
            删除任务
          </button>
          <div className="action-buttons-right">
            <button type="button" onClick={onClose} className="btn-cancel">
              关闭
            </button>
            {onCreateSubtask && (
              <button
                type="button"
                onClick={() => onCreateSubtask(task.id)}
                className="btn-secondary"
              >
                创建子任务
              </button>
            )}
            <button
              type="button"
              onClick={() => onEdit(task)}
              className="btn-primary"
            >
              编辑
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default TaskViewDialog;
