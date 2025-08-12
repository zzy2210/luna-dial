use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum TaskStatus {
    NotStarted,
    InProgress,
    Completed,
    Cancelled,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum TaskPriority {
    Low,
    Medium,
    High,
    Urgent,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum PeriodType {
    Daily,
    Weekly,
    Monthly,
    Quarterly,
    Yearly,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Period {
    pub start: DateTime<Utc>,
    pub end: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Task {
    pub id: String,
    pub title: String,
    pub task_type: PeriodType,
    pub time_period: Period,
    pub status: TaskStatus,
    pub tags: Vec<String>,
    pub icon: String,
    pub score: u32,
    pub priority: TaskPriority,
    pub parent_id: Option<String>,
    pub user_id: String,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,

    // 后端的树结构字段
    pub has_children: bool,
    pub children_count: u32,
    pub tree_depth: u32,
    pub root_task_id: Option<String>,
    pub children: Vec<Task>,
}

impl TaskStatus {
    pub fn icon(&self) -> &'static str {
        match self {
            TaskStatus::NotStarted => "⭕",
            TaskStatus::InProgress => "⏳",
            TaskStatus::Completed => "✓",
            TaskStatus::Cancelled => "⛔",
        }
    }

    pub fn next(&self) -> TaskStatus {
        match self {
            TaskStatus::NotStarted => TaskStatus::InProgress,
            TaskStatus::InProgress => TaskStatus::Completed,
            TaskStatus::Completed => TaskStatus::Cancelled,
            TaskStatus::Cancelled => TaskStatus::NotStarted,
        }
    }
}

impl TaskPriority {
    pub fn display(&self) -> &'static str {
        match self {
            TaskPriority::Low => "[低]",
            TaskPriority::Medium => "[中]",
            TaskPriority::High => "[高]",
            TaskPriority::Urgent => "[急]",
        }
    }
}
