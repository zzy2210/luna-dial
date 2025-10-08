import React from 'react';
import { Journal } from '../types';
import '../styles/dialog.css';

interface JournalViewDialogProps {
  journal: Journal;
  onClose: () => void;
  onEdit: (journal: Journal) => void;
}

const JournalViewDialog: React.FC<JournalViewDialogProps> = ({
  journal,
  onClose,
  onEdit
}) => {
  const journalTypeLabels = {
    0: '日志',
    1: '周志',
    2: '月志',
    3: '季志',
    4: '年志'
  };

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  return (
    <div className="dialog-overlay" onClick={onClose}>
      <div className="dialog-container large" onClick={e => e.stopPropagation()}>
        <div className="dialog-header">
          <h2>查看日志</h2>
          <button className="dialog-close" onClick={onClose}>×</button>
        </div>

        <div className="dialog-content">
          {/* 日志图标和标题 */}
          <div style={{
            display: 'flex',
            alignItems: 'center',
            gap: '12px',
            marginBottom: '24px',
            paddingBottom: '16px',
            borderBottom: '1px solid var(--border-color, #eee)'
          }}>
            <span style={{ fontSize: '32px' }}>{journal.icon || '📝'}</span>
            <div style={{ flex: 1 }}>
              <h3 style={{ margin: 0, fontSize: '20px', fontWeight: 600 }}>
                {journal.title}
              </h3>
              <div style={{
                display: 'flex',
                gap: '12px',
                marginTop: '8px',
                fontSize: '14px',
                color: 'var(--text-secondary, #666)'
              }}>
                <span>{journalTypeLabels[journal.journal_type]}</span>
                <span>•</span>
                <span>{formatDate(journal.time_period.start)} - {formatDate(journal.time_period.end)}</span>
              </div>
            </div>
          </div>

          {/* 日志内容 */}
          <div style={{ marginBottom: '24px' }}>
            <h4 style={{
              fontSize: '14px',
              fontWeight: 600,
              color: 'var(--text-secondary, #666)',
              marginBottom: '12px',
              textTransform: 'uppercase',
              letterSpacing: '0.5px'
            }}>
              内容
            </h4>
            <div style={{
              whiteSpace: 'pre-wrap',
              lineHeight: '1.6',
              padding: '16px',
              background: 'var(--bg-secondary, #f9f9f9)',
              borderRadius: '8px',
              fontSize: '15px'
            }}>
              {journal.content}
            </div>
          </div>

          {/* 创建时间 */}
          <div style={{
            fontSize: '13px',
            color: 'var(--text-tertiary, #999)',
            marginBottom: '24px'
          }}>
            创建于 {new Date(journal.created_at).toLocaleString('zh-CN')}
            {journal.updated_at !== journal.created_at && (
              <> · 最后更新 {new Date(journal.updated_at).toLocaleString('zh-CN')}</>
            )}
          </div>
        </div>

        {/* 操作按钮 */}
        <div className="dialog-actions">
          <button type="button" onClick={onClose} className="btn-cancel">
            关闭
          </button>
          <button
            type="button"
            onClick={() => onEdit(journal)}
            className="btn-primary"
          >
            编辑
          </button>
        </div>
      </div>
    </div>
  );
};

export default JournalViewDialog;
