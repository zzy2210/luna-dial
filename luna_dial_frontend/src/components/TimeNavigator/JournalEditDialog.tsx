import React, { useState, useEffect } from 'react';
import journalService from '../services/journal';
import { Journal, CreateJournalRequest, UpdateJournalRequest, PeriodType } from '../types';
import '../styles/dialog.css';

interface JournalEditDialogProps {
  journal?: Journal | null;
  onClose: () => void;
  onSuccess: () => void;
  currentPeriod?: PeriodType;
  currentDate?: Date;
}

const JournalEditDialog: React.FC<JournalEditDialogProps> = ({
  journal,
  onClose,
  onSuccess,
  currentPeriod = 'day',
  currentDate
}) => {
  const [loading, setLoading] = useState(false);
  const isEdit = !!journal;

  // 本地时间格式化函数，避免 toISOString() 的时区问题
  const formatLocalDate = (date: Date): string => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  };

  const [formData, setFormData] = useState<CreateJournalRequest>({
    title: '',
    content: '',
    journal_type: currentPeriod,
    start_date: formatLocalDate(new Date()),
    end_date: formatLocalDate(new Date()),
    icon: '📝'
  });

  // 根据当前周期设置默认日期
  const getDefaultDates = () => {
    const today = currentDate || new Date();
    const startDate = new Date(today);
    const endDate = new Date(today);

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
      start_date: formatLocalDate(startDate),
      end_date: formatLocalDate(endDate)
    };
  };

  useEffect(() => {
    if (journal) {
      // 编辑模式，加载现有日志数据
      const journalTypeMap = {
        0: 'day' as const,
        1: 'week' as const,
        2: 'month' as const,
        3: 'quarter' as const,
        4: 'year' as const
      };

      setFormData({
        title: journal.title,
        content: journal.content,
        journal_type: journalTypeMap[journal.journal_type] || 'day',
        start_date: journal.time_period?.start || formatLocalDate(new Date()),
        end_date: journal.time_period?.end || formatLocalDate(new Date()),
        icon: journal.icon || '📝'
      });
    } else {
      // 新建模式，设置默认值
      const dates = getDefaultDates();
      setFormData(prev => ({
        ...prev,
        journal_type: currentPeriod,
        ...dates
      }));
    }
  }, [journal, currentPeriod]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.title.trim()) {
      alert('请输入日志标题');
      return;
    }

    if (!formData.content.trim()) {
      alert('请输入日志内容');
      return;
    }

    setLoading(true);
    try {
      if (isEdit) {
        const updateData: UpdateJournalRequest = {
          journal_id: journal.id,
          title: formData.title,
          content: formData.content,
          journal_type: formData.journal_type,
          icon: formData.icon
        };
        await journalService.updateJournal(journal.id, updateData);
      } else {
        await journalService.createJournal(formData);
      }
      onSuccess();
    } catch (error) {
      console.error('Failed to save journal:', error);
      alert(isEdit ? '更新日志失败，请重试' : '创建日志失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  const handleIconSelect = (icon: string) => {
    setFormData(prev => ({ ...prev, icon }));
  };

  const icons = ['📝', '🌅', '🌙', '💭', '🎯', '📚', '💡', '🌟', '📊', '✨'];

  const journalTypeLabels = {
    day: '日志',
    week: '周志',
    month: '月志',
    quarter: '季志',
    year: '年志'
  };

  return (
    <div className="dialog-overlay" onClick={onClose}>
      <div className="dialog-container large" onClick={e => e.stopPropagation()}>
        <div className="dialog-header">
          <h2>{isEdit ? '编辑日志' : '新建日志'}</h2>
          <button className="dialog-close" onClick={onClose}>×</button>
        </div>

        <form onSubmit={handleSubmit} className="dialog-form">
          <div className="form-group">
            <label>日志标题 *</label>
            <input
              type="text"
              value={formData.title}
              onChange={e => setFormData(prev => ({ ...prev, title: e.target.value }))}
              placeholder="输入日志标题"
              maxLength={100}
              required
            />
          </div>

          <div className="form-group">
            <label>日志内容 *</label>
            <textarea
              value={formData.content}
              onChange={e => setFormData(prev => ({ ...prev, content: e.target.value }))}
              placeholder="记录你的想法、计划或总结..."
              rows={10}
              required
            />
            <div className="char-count">
              {formData.content.length} 字
            </div>
          </div>

          <div className="form-row">
            <div className="form-group">
              <label>日志类型</label>
              <select
                value={formData.journal_type}
                onChange={e => setFormData(prev => ({ ...prev, journal_type: e.target.value as PeriodType }))}
                disabled={isEdit}
              >
                <option value="day">日志</option>
                <option value="week">周志</option>
                <option value="month">月志</option>
                <option value="quarter">季志</option>
                <option value="year">年志</option>
              </select>
            </div>

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
            <label>日志图标</label>
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

          <div className="form-tips">
            <p>💡 提示：{journalTypeLabels[formData.journal_type]}用于记录{
              formData.journal_type === 'day' ? '每日' :
              formData.journal_type === 'week' ? '每周' :
              formData.journal_type === 'month' ? '每月' :
              formData.journal_type === 'quarter' ? '每季度' : '每年'
            }的计划、思考和总结。</p>
          </div>

          <div className="dialog-actions">
            <button type="button" onClick={onClose} className="btn-cancel">
              取消
            </button>
            <button type="submit" disabled={loading} className="btn-primary">
              {loading ? (isEdit ? '更新中...' : '创建中...') : (isEdit ? '更新日志' : '创建日志')}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default JournalEditDialog;