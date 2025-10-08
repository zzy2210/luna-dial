import React, { useState, useEffect } from 'react';
import '../styles/score-edit-dialog.css';

interface ScoreEditDialogProps {
  taskTitle: string;
  currentScore: number;
  onClose: () => void;
  onSave: (score: number) => void;
}

const ScoreEditDialog: React.FC<ScoreEditDialogProps> = ({
  taskTitle,
  currentScore,
  onClose,
  onSave,
}) => {
  const [score, setScore] = useState(currentScore);

  useEffect(() => {
    setScore(currentScore);
  }, [currentScore]);

  const handleSave = () => {
    onSave(score);
    onClose();
  };

  const handleOverlayClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      onClose();
    }
  };

  return (
    <div className="score-dialog-overlay" onClick={handleOverlayClick}>
      <div className="score-dialog">
        <div className="score-dialog-header">
          <h3>修改努力程度</h3>
          <button className="close-btn" onClick={onClose}>×</button>
        </div>

        <div className="score-dialog-body">
          <div className="task-title-display">
            <span className="task-label">任务：</span>
            <span className="task-name">{taskTitle}</span>
          </div>

          <div className="score-input-section">
            <label className="score-label">努力程度评分 (0-10)</label>

            {/* 滑块控件 */}
            <div className="slider-container">
              <input
                type="range"
                min="0"
                max="10"
                step="1"
                value={score}
                onChange={(e) => setScore(Number(e.target.value))}
                className="score-slider"
              />
              <div className="score-value-display">{score}</div>
            </div>

            {/* 快速选择按钮 */}
            <div className="quick-score-buttons">
              {[0, 2, 4, 6, 8, 10].map((value) => (
                <button
                  key={value}
                  className={`quick-score-btn ${score === value ? 'active' : ''}`}
                  onClick={() => setScore(value)}
                >
                  {value}
                </button>
              ))}
            </div>

            {/* 数字输入框（备选） */}
            <div className="number-input-container">
              <label>或直接输入：</label>
              <input
                type="number"
                min="0"
                max="10"
                value={score}
                onChange={(e) => {
                  const val = Number(e.target.value);
                  if (val >= 0 && val <= 10) {
                    setScore(val);
                  }
                }}
                className="score-number-input"
              />
            </div>
          </div>
        </div>

        <div className="score-dialog-footer">
          <button className="btn-cancel" onClick={onClose}>
            取消
          </button>
          <button className="btn-save" onClick={handleSave}>
            保存
          </button>
        </div>
      </div>
    </div>
  );
};

export default ScoreEditDialog;
