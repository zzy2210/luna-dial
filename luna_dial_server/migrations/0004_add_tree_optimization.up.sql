-- 添加树结构优化字段
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS has_children BOOLEAN DEFAULT FALSE;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS children_count INT DEFAULT 0;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS root_task_id VARCHAR(36);
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS tree_depth INT DEFAULT 0;

-- 创建优化索引
CREATE INDEX IF NOT EXISTS idx_tasks_root_task_id ON tasks(root_task_id);
CREATE INDEX IF NOT EXISTS idx_tasks_tree_depth ON tasks(tree_depth);
CREATE INDEX IF NOT EXISTS idx_tasks_no_parent ON tasks(user_id) WHERE parent_id IS NULL OR parent_id = '';
CREATE INDEX IF NOT EXISTS idx_tasks_has_children ON tasks(has_children);

-- 为根任务索引添加复合索引（用户ID + 根任务条件）
CREATE INDEX IF NOT EXISTS idx_tasks_user_root ON tasks(user_id, root_task_id);

-- 初始化现有数据的树结构字段
-- 注意：由于当前没有老数据，这些更新语句主要用于文档说明和未来数据迁移参考

-- 为根任务设置初始值
UPDATE tasks SET 
    root_task_id = id,
    tree_depth = 0,
    has_children = (SELECT COUNT(*) > 0 FROM tasks t2 WHERE t2.parent_id = tasks.id),
    children_count = (SELECT COUNT(*) FROM tasks t2 WHERE t2.parent_id = tasks.id)
WHERE parent_id IS NULL OR parent_id = '';