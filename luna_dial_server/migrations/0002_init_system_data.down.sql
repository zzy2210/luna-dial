-- 0002_init_system_data.down.sql
-- 回滚系统数据初始化

-- 删除系统配置数据
DELETE FROM system_configs WHERE config_key IN ('sm2_private_key', 'sm2_public_key', 'jwt_secret');

-- 删除管理员用户
DELETE FROM users WHERE user_name = 'admin';
