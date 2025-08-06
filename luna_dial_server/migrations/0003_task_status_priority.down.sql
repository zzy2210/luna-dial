-- 回滚迁移：恢复 is_completed 字段，移除 status 和 priority 字段

-- 重新添加 is_completed 字段
ALTER TABLE tasks ADD COLUMN is_completed BOOLEAN DEFAULT false NOT NULL;

-- 数据恢复：将 status 转换回 is_completed  
-- 已完成状态 (status = 2) 转换为 true
UPDATE tasks SET is_completed = true WHERE status = 2;
-- 其他状态 (status != 2) 转换为 false
UPDATE tasks SET is_completed = false WHERE status != 2;

-- 删除新增的字段
ALTER TABLE tasks DROP COLUMN priority;
ALTER TABLE tasks DROP COLUMN status;
