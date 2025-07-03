-- 删除任务表相关的触发器和索引
DROP TRIGGER IF EXISTS update_tasks_updated_at ON tasks;
DROP INDEX IF EXISTS idx_tasks_user_parent;
DROP INDEX IF EXISTS idx_tasks_user_status;
DROP INDEX IF EXISTS idx_tasks_user_type;
DROP INDEX IF EXISTS idx_tasks_created_at;
DROP INDEX IF EXISTS idx_tasks_end_date;
DROP INDEX IF EXISTS idx_tasks_start_date;
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_tasks_type;
DROP INDEX IF EXISTS idx_tasks_parent_id;
DROP INDEX IF EXISTS idx_tasks_user_id;
DROP TABLE IF EXISTS tasks;
