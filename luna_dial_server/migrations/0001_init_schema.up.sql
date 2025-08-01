-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    user_name VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_users_user_name (user_name),
    INDEX idx_users_email (email)
);

-- 创建任务表
CREATE TABLE IF NOT EXISTS tasks (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    task_type INT NOT NULL,
    period_start DATETIME,
    period_end DATETIME,
    tags TEXT,
    icon VARCHAR(10),
    score INT DEFAULT 0,
    is_completed BOOLEAN DEFAULT FALSE,
    parent_id VARCHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_tasks_user_id (user_id),
    INDEX idx_tasks_parent_id (parent_id),
    INDEX idx_tasks_period (period_start, period_end),
    INDEX idx_tasks_type (task_type),
    INDEX idx_tasks_user_period (user_id, period_start, period_end)
);

-- 创建日志表
CREATE TABLE IF NOT EXISTS journals (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    journal_type INT NOT NULL,
    period_start DATETIME,
    period_end DATETIME,
    icon VARCHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_journals_user_id (user_id),
    INDEX idx_journals_period (period_start, period_end),
    INDEX idx_journals_type (journal_type),
    INDEX idx_journals_user_period (user_id, period_start, period_end)
);


-- 创建系统配置表
CREATE TABLE IF NOT EXISTS system_configs (
    id VARCHAR(36) PRIMARY KEY,
    config_key VARCHAR(100) UNIQUE NOT NULL,
    config_value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_configs_key (config_key)
);



-- 创建迁移记录表（如果不存在）
CREATE TABLE IF NOT EXISTS schema_migrations (
    version BIGINT PRIMARY KEY,
    dirty BOOLEAN NOT NULL
);

