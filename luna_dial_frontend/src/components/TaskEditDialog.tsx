import React, { useState } from 'react';
import taskService from '../services/task';
import { CreateTaskRequest, UpdateTaskRequest, PeriodType, Task } from '../types';
import { getPeriodDates } from '../utils/dateUtils';
import '../styles/dialog.css';

interface TaskEditDialogProps {
  task?: Task | null;
  onClose: () => void;
  onSuccess: () => void;
  currentPeriod?: PeriodType;
   currentDate?: Date;
  parentTaskId?: string;
}

const TaskEditDialog: React.FC<TaskEditDialogProps> = ({
  task,
  onClose,
  onSuccess,
  currentPeriod = 'day',
   currentDate,
  parentTaskId
}) => {
  const [loading, setLoading] = useState(false);
  const [tagInput, setTagInput] = useState('');
  const isEdit = !!task;

  // 将ISO时间字符串转换为YYYY-MM-DD格式
  const isoToDateInput = (isoString: string): string => {
    return isoString.split('T')[0];
  };

  // 根据周期类型计算默认日期（左闭右开）
  const getDefaultDates = (periodType: PeriodType) => {
    // 优先使用外部传入的当前界面日期，否则回退到今天
    const baseDate = currentDate || new Date();
    const dates = getPeriodDates(periodType, baseDate);
    // Task 创建接口期望的是左闭右开的日期区间的日期字符串，这里与 Dashboard 保持一致
    return {
      start_date: dates.start_date,
      end_date: dates.end_date
    };
  };

  // 解析标签
  const parseTags = (tags?: string): string[] => {
    if (!tags) return [];
    try {
      return JSON.parse(tags);
    } catch {
      return [];
    }
  };

  // 将任务类型数字转换为字符串
  const taskTypeToPeriodType = (taskType: number): PeriodType => {
    const map: Record<number, PeriodType> = {
      0: 'day',
      1: 'week',
      2: 'month',
      3: 'quarter',
      4: 'year'
    };
    return map[taskType] || 'day';
  };

  // 将优先级数字转换为字符串
  const priorityToString = (priority: number): 'low' | 'medium' | 'high' | 'urgent' => {
    const map: Record<number, 'low' | 'medium' | 'high' | 'urgent'> = {
      0: 'low',
      1: 'medium',
      2: 'high',
      3: 'urgent'
    };
    return map[priority] || 'medium';
  };

  // 初始化表单数据
  const [formData, setFormData] = useState<CreateTaskRequest>(() => {
    if (task) {
      // 编辑模式：从现有任务加载数据
      return {
        title: task.title,
        start_date: isoToDateInput(task.period.start),
        end_date: isoToDateInput(task.period.end),
        period_type: taskTypeToPeriodType(task.task_type),
        priority: priorityToString(task.priority),
        icon: task.icon || '📝',
        tags: parseTags(task.tags),
        parent_id: task.parent_id
      };
    } else {
      // 新建模式：使用默认值
      const dates = getDefaultDates(currentPeriod);
      return {
        title: '',
        start_date: dates.start_date,
        end_date: dates.end_date,
        period_type: currentPeriod,
        priority: 'medium',
        icon: '📝',
        tags: [],
        parent_id: parentTaskId
      };
    }
  });

  const [score, setScore] = useState<number>(task?.score || 0);

  // 处理周期类型改变：自动更新日期为符合规范的格式
  const handlePeriodTypeChange = (newPeriodType: PeriodType) => {
    const dates = getDefaultDates(newPeriodType);
    setFormData(prev => ({
      ...prev,
      period_type: newPeriodType,
      start_date: dates.start_date,
      end_date: dates.end_date
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.title.trim()) {
      alert('请输入任务标题');
      return;
    }

    setLoading(true);
    try {
      if (isEdit && task) {
        // 编辑模式
        const updateData: UpdateTaskRequest = {
          title: formData.title,
          priority: formData.priority,
          icon: formData.icon,
          tags: formData.tags
        };
        await taskService.updateTask(task.id, updateData);

        // 如果评分有变化，单独更新评分
        if (score !== task.score) {
          await taskService.updateScore(task.id, score);
        }
      } else {
        // 新建模式
        if (parentTaskId) {
          await taskService.createSubtask(parentTaskId, formData);
        } else {
          await taskService.createTask(formData);
        }
      }
      onSuccess();
    } catch (error) {
      console.error('Failed to save task:', error);
      alert(isEdit ? '更新任务失败，请重试' : '创建任务失败，请重试');
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
          <h2>{isEdit ? '编辑任务' : (parentTaskId ? '创建子任务' : '创建新任务')}</h2>
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

          <div className="form-row">
            <div className="form-group">
              <label>开始日期</label>
              <input
                type="date"
                value={formData.start_date}
                onChange={e => setFormData(prev => ({ ...prev, start_date: e.target.value }))}
                disabled={isEdit}
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
                disabled={isEdit}
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

          {isEdit && (
            <div className="form-group">
              <label>努力评分 (0-10)</label>
              <div className="score-input-group">
                <input
                  type="range"
                  min="0"
                  max="10"
                  step="1"
                  value={score}
                  onChange={e => setScore(parseInt(e.target.value))}
                  className="score-slider-edit"
                />
                <input
                  type="number"
                  min="0"
                  max="10"
                  value={score}
                  onChange={e => setScore(Math.min(10, Math.max(0, parseInt(e.target.value) || 0)))}
                  className="score-number-edit"
                />
                <span className="score-display-text">{score} / 10</span>
              </div>
            </div>
          )}

          <div className="form-group">
            <label>任务周期</label>
            <select
              value={formData.period_type}
              onChange={e => handlePeriodTypeChange(e.target.value as PeriodType)}
              disabled={isEdit}
            >
              <option value="day">日任务</option>
              <option value="week">周任务</option>
              <option value="month">月任务</option>
              <option value="quarter">季度任务</option>
              <option value="year">年度任务</option>
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
                onKeyDown={e => {
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
              {loading ? (isEdit ? '更新中...' : '创建中...') : (isEdit ? '更新任务' : '创建任务')}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default TaskEditDialog;
