import React, { useMemo, useRef, useState, useEffect } from 'react';
import journalService from '../services/journal';
import { Journal, CreateJournalRequest, UpdateJournalRequest, PeriodType } from '../types';
import { formatLocalDate, getPeriodDates } from '../utils/dateUtils';
import { JOURNAL_ICONS } from '../constants/icons';
import '../styles/dialog.css';

// 日志草稿缓存版本号
const JOURNAL_DRAFT_VERSION = 1 as const;

// 日志草稿缓存结构（v1）
interface JournalDraftV1 {
  v: typeof JOURNAL_DRAFT_VERSION;
  formData: CreateJournalRequest;
  updatedAt: number;
}

// 生成日志草稿缓存键（避免刷新/跳转导致内容丢失）
const buildJournalDraftKey = (params: {
  isEdit: boolean;
  journalId: string | null;
  formData: Pick<CreateJournalRequest, 'journal_type' | 'start_date' | 'end_date'>;
}): string => {
  if (params.isEdit && params.journalId) {
    return `journal_draft:v${JOURNAL_DRAFT_VERSION}:edit:${params.journalId}`;
  }
  return `journal_draft:v${JOURNAL_DRAFT_VERSION}:new:${params.formData.journal_type}:${params.formData.start_date}:${params.formData.end_date}`;
};

// 解析日志草稿缓存内容（失败则返回 null）
const parseJournalDraft = (raw: string | null): JournalDraftV1 | null => {
  if (!raw) return null;
  try {
    const obj = JSON.parse(raw) as Partial<JournalDraftV1>;
    if (!obj || obj.v !== JOURNAL_DRAFT_VERSION) return null;
    if (!obj.formData) return null;
    if (typeof obj.updatedAt !== 'number') return null;
    return obj as JournalDraftV1;
  } catch {
    return null;
  }
};

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

  const hasPromptedDraftRestoreRef = useRef(false);
  const previousDraftKeyRef = useRef<string | null>(null);
  const [hasLoadedInitialFormData, setHasLoadedInitialFormData] = useState(false);

  const [formData, setFormData] = useState<CreateJournalRequest>(() => {
    const dates = getPeriodDates(currentPeriod, currentDate || new Date());
    return {
      title: '',
      content: '',
      journal_type: currentPeriod,
      start_date: dates.start_date,
      end_date: dates.end_date,
      icon: '📝'
    };
  });

  useEffect(() => {
    setHasLoadedInitialFormData(false);
    hasPromptedDraftRestoreRef.current = false;

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
      const dates = getPeriodDates(currentPeriod, currentDate || new Date());
      setFormData(prev => ({
        ...prev,
        journal_type: currentPeriod,
        ...dates
      }));
    }
    setHasLoadedInitialFormData(true);
  }, [journal, currentPeriod, currentDate]);

  const draftKey = useMemo(() => {
    return buildJournalDraftKey({
      isEdit,
      journalId: journal?.id ?? null,
      formData: {
        journal_type: formData.journal_type,
        start_date: formData.start_date,
        end_date: formData.end_date
      }
    });
  }, [isEdit, journal?.id, formData.journal_type, formData.start_date, formData.end_date]);

  useEffect(() => {
    // 初始化时记录一次 key，确保后续能正确迁移
    if (!previousDraftKeyRef.current) {
      previousDraftKeyRef.current = draftKey;
    }
  }, [draftKey]);

  useEffect(() => {
    if (isEdit) return;
    const previousKey = previousDraftKeyRef.current;
    if (!previousKey) {
      previousDraftKeyRef.current = draftKey;
      return;
    }
    if (previousKey === draftKey) return;

    const raw = localStorage.getItem(previousKey);
    if (raw && !localStorage.getItem(draftKey)) {
      localStorage.setItem(draftKey, raw);
    }
    localStorage.removeItem(previousKey);
    previousDraftKeyRef.current = draftKey;
  }, [isEdit, draftKey]);

  useEffect(() => {
    if (!hasLoadedInitialFormData) return;
    if (hasPromptedDraftRestoreRef.current) return;

    const draft = parseJournalDraft(localStorage.getItem(draftKey));
    if (!draft) {
      hasPromptedDraftRestoreRef.current = true;
      return;
    }

    const draftSnapshot = JSON.stringify(draft.formData);
    const currentSnapshot = JSON.stringify(formData);
    if (draftSnapshot === currentSnapshot) {
      hasPromptedDraftRestoreRef.current = true;
      return;
    }

    const ok = window.confirm('检测到未保存的日志草稿，是否恢复？');
    hasPromptedDraftRestoreRef.current = true;
    if (!ok) return;

    if (isEdit) {
      setFormData(prev => ({
        ...prev,
        title: draft.formData.title,
        content: draft.formData.content,
        icon: draft.formData.icon ?? prev.icon
      }));
      return;
    }

    setFormData(draft.formData);
  }, [draftKey, formData, hasLoadedInitialFormData, isEdit]);

  useEffect(() => {
    const handle = window.setTimeout(() => {
      const draft: JournalDraftV1 = {
        v: JOURNAL_DRAFT_VERSION,
        formData,
        updatedAt: Date.now()
      };
      localStorage.setItem(draftKey, JSON.stringify(draft));
    }, 350);

    return () => window.clearTimeout(handle);
  }, [draftKey, formData]);

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
      // 提交前先写入草稿，避免提交过程中被 401 跳转导致完全丢失
      const draft: JournalDraftV1 = {
        v: JOURNAL_DRAFT_VERSION,
        formData,
        updatedAt: Date.now()
      };
      localStorage.setItem(draftKey, JSON.stringify(draft));

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
      localStorage.removeItem(draftKey);
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
              {JOURNAL_ICONS.map(icon => (
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
