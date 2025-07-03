-- 删除日志条目表相关的触发器和索引
DROP TRIGGER IF EXISTS update_journal_entries_updated_at ON journal_entries;
DROP INDEX IF EXISTS idx_journal_entry_tasks_task_id;
DROP INDEX IF EXISTS idx_journal_entry_tasks_journal_id;
DROP INDEX IF EXISTS idx_journal_entries_user_type;
DROP INDEX IF EXISTS idx_journal_entries_user_time_scale;
DROP INDEX IF EXISTS idx_journal_entries_created_at;
DROP INDEX IF EXISTS idx_journal_entries_time_reference;
DROP INDEX IF EXISTS idx_journal_entries_entry_type;
DROP INDEX IF EXISTS idx_journal_entries_time_scale;
DROP INDEX IF EXISTS idx_journal_entries_user_id;
DROP TABLE IF EXISTS journal_entry_tasks;
DROP TABLE IF EXISTS journal_entries;
