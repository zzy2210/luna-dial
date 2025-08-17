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
    #[serde(rename = "daily")]
    Daily,
    #[serde(rename = "weekly")]
    Weekly,
    #[serde(rename = "monthly")]
    Monthly,
    #[serde(rename = "quarterly")]
    Quarterly,
    #[serde(rename = "yearly")]
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
    #[serde(rename = "type")]
    pub task_type: PeriodType,
    #[serde(rename = "period")]
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

// API响应的通用结构
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ApiResponse<T> {
    pub code: u32,
    pub message: String,
    pub success: bool,
    pub timestamp: u64,
    pub data: T,
}

// 计划请求体
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PlanRequest {
    pub period_type: String, // "day", "week", "month", "quarter", "year"
    pub start_date: DateTime<Utc>,
    pub end_date: DateTime<Utc>,
}

// 计划响应数据
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PlanData {
    pub tasks: Vec<Task>,
    pub tasks_total: u32,
    pub journals: Vec<Journal>,
    pub journals_total: u32,
    pub plan_type: String,
    pub plan_period: Period,
    pub score_total: u32,
    pub group_stats: Vec<GroupStat>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GroupStat {
    pub group_key: String,
    pub task_count: u32,
    pub score_total: u32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Journal {
    pub id: String,
    pub title: String,
    pub content: String,
    pub journal_type: String,
    pub time_period: Period,
    pub icon: String,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    pub user_id: String,
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
