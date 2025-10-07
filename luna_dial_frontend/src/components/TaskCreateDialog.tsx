import React, { useState } from 'react';
import taskService from '../services/task';
import { CreateTaskRequest, PeriodType } from '../types';
import '../styles/dialog.css';

interface TaskCreateDialogProps {
  onClose: () => void;
  onSuccess: () => void;
  currentPeriod?: PeriodType;
  parentTaskId?: string;
}

const TaskCreateDialog: React.FC<TaskCreateDialogProps> = ({
  onClose,
  onSuccess,
  currentPeriod = 'day',
  parentTaskId
}) => {
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState<CreateTaskRequest>({
    title: '',
    description: '',
    start_date: new Date().toISOString().split('T')[0],
    end_date: new Date().toISOString().split('T')[0],
    priority: 'medium',
    icon: '📝',
    tags: [],
    parent_id: parentTaskId
  });

  const [tagInput, setTagInput] = useState('');

  // 根据当前周期设置默认日期
  const getDefaultDates = () => {
    const today = new Date();
    const startDate = new Date();
    const endDate = new Date();

    switch (currentPeriod) {
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

  React.useEffect(() => {
    const dates = getDefaultDates();
    setFormData(prev => ({
      ...prev,
      ...dates
    }));
  }, [currentPeriod]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.title.trim()) {
      alert('请输入任务标题');
      return;
    }

    setLoading(true);
    try {
      if (parentTaskId) {
        await taskService.createSubtask(parentTaskId, formData);
      } else {
        await taskService.createTask(formData);
      }
      onSuccess();
    } catch (error) {
      console.error('Failed to create task:', error);
      alert('创建任务失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  const handleAddTag = () => {
    if (tagInput.trim() && !formData.tags?.includes(tagInput.trim())) {
      setFormData(prev => ({
        ...prev,
        tags: [...(prev.tags || []), tagInput.trim()]
      }));
      setTagInput('');
    }
  };

  const handleRemoveTag = (tag: string) => {
    setFormData(prev => ({
      ...prev,
      tags: prev.tags?.filter(t => t !== tag) || []
    }));
  };

  const handleIconSelect = (icon: string) => {
    setFormData(prev => ({ ...prev, icon }));
  };

  const icons = ['📝', '💡', '🎯', '📚', '💻', '🏃', '🎨', '🌟', '🔧', '📊'];

  return (
    <div className="dialog-overlay" onClick={onClose}>
      <div className="dialog-container" onClick={e => e.stopPropagation()}>
        <div className="dialog-header">
          <h2>{parentTaskId ? '创建子任务' : '创建新任务'}</h2>
          <button className="dialog-close" onClick={onClose}>×</button>
        </div>

        <form onSubmit={handleSubmit} className="dialog-form">
          <div className="form-group">
            <label>任务标题 *</label>
            <input
              type="text"
              value={formData.title}
              onChange={e => setFormData(prev => ({ ...prev, title: e.target.value }))}
              placeholder="输入任务标题"
              maxLength={100}
              required
            />
          </div>

          <div className="form-group">
            <label>任务描述</label>
            <textarea
              value={formData.description}
              onChange={e => setFormData(prev => ({ ...prev, description: e.target.value }))}
              placeholder="输入任务描述（可选）"
              rows={3}
            />
          </div>

          <div className="form-row">
            <div className="form-group">
              <label>开始日期</label>
              <input
                type="date"
                value={formData.start_date}
                onChange={e => setFormData(prev => ({ ...prev, start_date: e.target.value }))}
                required
              />
            </div>

            <div className="form-group">
              <label>结束日期</label>
              <input
                type="date"
                value={formData.end_date}
                onChange={e => setFormData(prev => ({ ...prev, end_date: e.target.value }))}
                min={formData.start_date}
                required
              />
            </div>
          </div>

          <div className="form-group">
            <label>优先级</label>
            <select
              value={formData.priority}
              onChange={e => setFormData(prev => ({ ...prev, priority: e.target.value as any }))}
            >
              <option value="low">低</option>
              <option value="medium">中</option>
              <option value="high">高</option>
              <option value="urgent">紧急</option>
            </select>
          </div>

          <div className="form-group">
            <label>任务图标</label>
            <div className="icon-selector">
              {icons.map(icon => (
                <button
                  key={icon}
                  type="button"
                  className={`icon-btn ${formData.icon === icon ? 'selected' : ''}`}
                  onClick={() => handleIconSelect(icon)}
                >
                  {icon}
                </button>
              ))}
            </div>
          </div>

          <div className="form-group">
            <label>标签</label>
            <div className="tag-input-container">
              <input
                type="text"
                value={tagInput}
                onChange={e => setTagInput(e.target.value)}
                placeholder="输入标签后按回车或点击添加"
                onKeyPress={e => {
                  if (e.key === 'Enter') {
                    e.preventDefault();
                    handleAddTag();
                  }
                }}
              />
              <button type="button" onClick={handleAddTag} className="btn-add-tag">
                添加
              </button>
            </div>
            {formData.tags && formData.tags.length > 0 && (
              <div className="tag-list">
                {formData.tags.map(tag => (
                  <span key={tag} className="tag">
                    {tag}
                    <button
                      type="button"
                      onClick={() => handleRemoveTag(tag)}
                      className="tag-remove"
                    >
                      ×
                    </button>
                  </span>
                ))}
              </div>
            )}
          </div>

          <div className="dialog-actions">
            <button type="button" onClick={onClose} className="btn-cancel">
              取消
            </button>
            <button type="submit" disabled={loading} className="btn-primary">
              {loading ? '创建中...' : '创建任务'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default TaskCreateDialog;