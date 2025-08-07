-- 删除树结构优化相关的索引
DROP INDEX IF EXISTS idx_tasks_user_root;
DROP INDEX IF EXISTS idx_tasks_has_children;
DROP INDEX IF EXISTS idx_tasks_tree_depth;
DROP INDEX IF EXISTS idx_tasks_no_parent;
DROP INDEX IF EXISTS idx_tasks_root_task_id;

-- 删除树结构优化字段
ALTER TABLE tasks DROP COLUMN IF EXISTS tree_depth;
ALTER TABLE tasks DROP COLUMN IF EXISTS root_task_id;
ALTER TABLE tasks DROP COLUMN IF EXISTS children_count;
ALTER TABLE tasks DROP COLUMN IF EXISTS has_children;
