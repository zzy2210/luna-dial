import React from 'react';
import { Task, TaskStatus, TaskPriority } from '../types';
import '../styles/dialog.css';

interface TaskViewDialogProps {
  task: Task;
  onClose: () => void;
  onEdit: (task: Task) => void;
  onDelete: (taskId: string) => void;
}

const TaskViewDialog: React.FC<TaskViewDialogProps> = ({
  task,
  onClose,
  onEdit,
  onDelete
}) => {
  const taskTypeLabels = {
    0: '日任务',
    1: '周任务',
    2: '月任务',
    3: '季度任务',
    4: '年度任务'
  };

  const statusLabels = {
    [TaskStatus.NotStarted]: '未开始',
    [TaskStatus.InProgress]: '进行中',
    [TaskStatus.Completed]: '已完成',
    [TaskStatus.Cancelled]: '已取消'
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

  const parseTags = (tags?: string): string[] => {
    if (!tags) return [];
    try {
      return JSON.parse(tags);
    } catch {
      return [];
    }
  };

  const handleDelete = () => {
    if (window.confirm('确定要删除这个任务吗？')) {
      onDelete(task.id);
    }
  };

  const tags = parseTags(task.tags);

  return (
    <div className="dialog-overlay" onClick={onClose}>
      <div className="dialog-container large" onClick={e => e.stopPropagation()}>
        <div className="dialog-header">
          <h2>任务详情</h2>
          <button className="dialog-close" onClick={onClose}>×</button>
        </div>

        <div className="dialog-content">
          {/* 任务图标和标题 */}
          <div style={{
            display: 'flex',
            alignItems: 'center',
            gap: '12px',
            marginBottom: '24px',
            paddingBottom: '16px',
            borderBottom: '1px solid var(--border-color, #eee)'
          }}>
            <span style={{ fontSize: '32px' }}>{task.icon || '📋'}</span>
            <div style={{ flex: 1 }}>
              <h3 style={{ margin: 0, fontSize: '20px', fontWeight: 600 }}>
                {task.title}
              </h3>
              <div style={{
                display: 'flex',
                gap: '12px',
                marginTop: '8px',
                fontSize: '14px',
                color: 'var(--text-secondary, #666)'
              }}>
                <span>{taskTypeLabels[task.task_type]}</span>
                <span>•</span>
                <span>优先级: {priorityLabels[task.priority]}</span>
                <span>•</span>
                <span>状态: {statusLabels[task.status]}</span>
              </div>
            </div>
          </div>

          {/* 任务描述 */}
          {task.description && (
            <div style={{ marginBottom: '24px' }}>
              <h4 style={{
                fontSize: '14px',
                fontWeight: 600,
                color: 'var(--text-secondary, #666)',
                marginBottom: '12px',
                textTransform: 'uppercase',
                letterSpacing: '0.5px'
              }}>
                描述
              </h4>
              <div style={{
                whiteSpace: 'pre-wrap',
                lineHeight: '1.6',
                padding: '16px',
                background: 'var(--bg-secondary, #f9f9f9)',
                borderRadius: '8px',
                fontSize: '15px'
              }}>
                {task.description}
              </div>
            </div>
          )}

          {/* 任务时间范围 */}
          <div style={{ marginBottom: '24px' }}>
            <h4 style={{
              fontSize: '14px',
              fontWeight: 600,
              color: 'var(--text-secondary, #666)',
              marginBottom: '12px',
              textTransform: 'uppercase',
              letterSpacing: '0.5px'
            }}>
              时间范围
            </h4>
            <div style={{
              padding: '16px',
              background: 'var(--bg-secondary, #f9f9f9)',
              borderRadius: '8px',
              fontSize: '15px'
            }}>
              {formatDate(task.period_start)} - {formatDate(task.period_end)}
            </div>
          </div>

          {/* 努力评分 */}
          <div style={{ marginBottom: '24px' }}>
            <h4 style={{
              fontSize: '14px',
              fontWeight: 600,
              color: 'var(--text-secondary, #666)',
              marginBottom: '12px',
              textTransform: 'uppercase',
              letterSpacing: '0.5px'
            }}>
              努力评分
            </h4>
            <div style={{
              padding: '16px',
              background: 'var(--bg-secondary, #f9f9f9)',
              borderRadius: '8px',
              fontSize: '15px'
            }}>
              {task.status === TaskStatus.NotStarted ? '未开始' : `${task.score} / 10`}
            </div>
          </div>

          {/* 标签 */}
          {tags.length > 0 && (
            <div style={{ marginBottom: '24px' }}>
              <h4 style={{
                fontSize: '14px',
                fontWeight: 600,
                color: 'var(--text-secondary, #666)',
                marginBottom: '12px',
                textTransform: 'uppercase',
                letterSpacing: '0.5px'
              }}>
                标签
              </h4>
              <div style={{
                display: 'flex',
                gap: '8px',
                flexWrap: 'wrap'
              }}>
                {tags.map(tag => (
                  <span key={tag} style={{
                    padding: '4px 12px',
                    background: 'var(--bg-secondary, #f0f0f0)',
                    borderRadius: '12px',
                    fontSize: '13px',
                    color: 'var(--text-primary, #333)'
                  }}>
                    {tag}
                  </span>
                ))}
              </div>
            </div>
          )}

          {/* 任务层级信息 */}
          {(task.parent_id || task.has_children) && (
            <div style={{ marginBottom: '24px' }}>
              <h4 style={{
                fontSize: '14px',
                fontWeight: 600,
                color: 'var(--text-secondary, #666)',
                marginBottom: '12px',
                textTransform: 'uppercase',
                letterSpacing: '0.5px'
              }}>
                任务层级
              </h4>
              <div style={{
                padding: '16px',
                background: 'var(--bg-secondary, #f9f9f9)',
                borderRadius: '8px',
                fontSize: '15px'
              }}>
                {task.parent_id && <div>父任务 ID: {task.parent_id}</div>}
                {task.has_children && (
                  <div>子任务数量: {task.children_count}</div>
                )}
                <div>树深度: {task.tree_depth}</div>
              </div>
            </div>
          )}

          {/* 创建时间 */}
          <div style={{
            fontSize: '13px',
            color: 'var(--text-tertiary, #999)',
            marginBottom: '24px'
          }}>
            创建于 {new Date(task.created_at).toLocaleString('zh-CN')}
            {task.updated_at !== task.created_at && (
              <> · 最后更新 {new Date(task.updated_at).toLocaleString('zh-CN')}</>
            )}
          </div>
        </div>

        {/* 操作按钮 */}
        <div className="dialog-actions">
          <button type="button" onClick={handleDelete} className="btn-cancel" style={{
            marginRight: 'auto',
            background: '#dc3545',
            color: 'white'
          }}>
            删除
          </button>
          <button type="button" onClick={onClose} className="btn-cancel">
            关闭
          </button>
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
  );
};

export default TaskViewDialog;
