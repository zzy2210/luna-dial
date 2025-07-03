-- 创建日志条目表
CREATE TABLE IF NOT EXISTS journal_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content TEXT NOT NULL,
    time_reference VARCHAR(50) NOT NULL,
    time_scale VARCHAR(20) NOT NULL CHECK (time_scale IN ('day', 'week', 'month', 'quarter', 'year')),
    entry_type VARCHAR(20) NOT NULL CHECK (entry_type IN ('plan-start', 'reflection', 'summary')),
    user_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- 外键约束
    CONSTRAINT fk_journal_entries_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 创建日志条目与任务的关联表（多对多关系）
CREATE TABLE IF NOT EXISTS journal_entry_tasks (
    journal_entry_id UUID NOT NULL,
    task_id UUID NOT NULL,
    
    PRIMARY KEY (journal_entry_id, task_id),
    
    -- 外键约束
    CONSTRAINT fk_journal_entry_tasks_journal_id FOREIGN KEY (journal_entry_id) REFERENCES journal_entries(id) ON DELETE CASCADE,
    CONSTRAINT fk_journal_entry_tasks_task_id FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

-- 为日志条目表创建更新时间戳的触发器
CREATE TRIGGER update_journal_entries_updated_at
    BEFORE UPDATE ON journal_entries
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_journal_entries_user_id ON journal_entries(user_id);
CREATE INDEX IF NOT EXISTS idx_journal_entries_time_scale ON journal_entries(time_scale);
CREATE INDEX IF NOT EXISTS idx_journal_entries_entry_type ON journal_entries(entry_type);
CREATE INDEX IF NOT EXISTS idx_journal_entries_time_reference ON journal_entries(time_reference);
CREATE INDEX IF NOT EXISTS idx_journal_entries_created_at ON journal_entries(created_at);
-- 复合索引用于用户相关查询
CREATE INDEX IF NOT EXISTS idx_journal_entries_user_time_scale ON journal_entries(user_id, time_scale);
CREATE INDEX IF NOT EXISTS idx_journal_entries_user_type ON journal_entries(user_id, entry_type);

-- 关联表索引
CREATE INDEX IF NOT EXISTS idx_journal_entry_tasks_journal_id ON journal_entry_tasks(journal_entry_id);
CREATE INDEX IF NOT EXISTS idx_journal_entry_tasks_task_id ON journal_entry_tasks(task_id);
