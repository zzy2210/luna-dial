-- 添加新字段：状态和优先级
-- 状态：0=未开始, 1=进行中, 2=已完成, 3=已取消
-- 优先级：0=低, 1=中, 2=高, 3=紧急
ALTER TABLE tasks ADD COLUMN status INT DEFAULT 0 NOT NULL;
ALTER TABLE tasks ADD COLUMN priority INT DEFAULT 0 NOT NULL;

-- 数据迁移：将 is_completed 转换为 status
-- 所有未完成的任务设为"未开始"状态 (0)
UPDATE tasks SET status = 0 WHERE is_completed = false;
-- 所有已完成的任务设为"已完成"状态 (2)  
UPDATE tasks SET status = 2 WHERE is_completed = true;

-- 删除旧的 is_completed 字段
ALTER TABLE tasks DROP COLUMN is_completed;
