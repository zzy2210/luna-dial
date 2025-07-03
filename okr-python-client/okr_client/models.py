"""数据模型定义"""

from datetime import datetime
from typing import Optional, List
from pydantic import BaseModel, Field
from enum import Enum


class TaskType(str, Enum):
    """任务类型枚举"""
    YEAR = "year"
    QUARTER = "quarter"
    MONTH = "month"
    WEEK = "week"
    DAY = "day"


class TaskStatus(str, Enum):
    """任务状态枚举"""
    PENDING = "pending"
    IN_PROGRESS = "in-progress"
    COMPLETED = "completed"


class TimeScale(str, Enum):
    """时间尺度枚举"""
    DAY = "day"
    WEEK = "week"
    MONTH = "month"
    QUARTER = "quarter"
    YEAR = "year"


class EntryType(str, Enum):
    """日志条目类型枚举"""
    PLAN_START = "plan-start"
    REFLECTION = "reflection"
    SUMMARY = "summary"


class User(BaseModel):
    """用户模型"""
    id: str
    username: str
    email: str
    created_at: datetime
    updated_at: datetime


class Task(BaseModel):
    """任务模型"""
    id: str
    title: str
    description: Optional[str] = None
    type: TaskType
    start_date: datetime
    end_date: datetime
    status: TaskStatus
    score: Optional[int] = Field(None, ge=0, le=10)
    parent_id: Optional[str] = None
    user_id: str
    tags: Optional[str] = None
    created_at: datetime
    updated_at: datetime


class JournalEntry(BaseModel):
    """日志条目模型"""
    id: str
    content: str
    time_reference: str
    time_scale: TimeScale
    entry_type: EntryType
    user_id: str
    created_at: datetime
    updated_at: datetime


# 新增扩展模型

class TimeRange(BaseModel):
    """时间范围模型"""
    start: datetime
    end: datetime


class TaskTree(Task):
    """任务树模型，兼容扁平结构"""
    children: List['TaskTree'] = []
    
    class Config:
        arbitrary_types_allowed = True


class PlanStats(BaseModel):
    """计划统计信息"""
    total_tasks: int
    completed_tasks: int
    in_progress_tasks: int
    pending_tasks: int
    total_score: int
    completed_score: int


class PlanRequest(BaseModel):
    """计划视图请求模型"""
    scale: TimeScale
    time_ref: str


class PlanResponse(BaseModel):
    """计划视图响应模型"""
    tasks: List[TaskTree] = []
    journals: List[JournalEntry]
    time_range: TimeRange
    stats: PlanStats


class ScoreTrendRequest(BaseModel):
    """分数趋势请求模型"""
    scale: TimeScale
    time_ref: str


class TrendSummary(BaseModel):
    """趋势摘要"""
    total_score: int
    total_tasks: int
    average_score: float
    average_task_count: float
    max_score: int
    max_tasks: int
    min_score: int
    min_tasks: int


class ScoreTrendResponse(BaseModel):
    """分数趋势响应模型"""
    labels: List[str]
    scores: List[int]
    counts: List[int]
    scale: str  # 兼容后端不规范响应
    time_ref: str
    time_range: TimeRange
    summary: Optional[TrendSummary]


# 现有模型

class TaskRequest(BaseModel):
    """创建/更新任务请求"""
    title: str
    description: Optional[str] = None
    type: TaskType
    start_date: datetime
    end_date: datetime
    status: TaskStatus = TaskStatus.PENDING
    score: Optional[int] = Field(None, ge=0, le=10)
    parent_id: Optional[str] = None
    tags: Optional[str] = None


class JournalRequest(BaseModel):
    """创建/更新日志请求"""
    content: str
    time_reference: Optional[str] = None
    time_scale: TimeScale = TimeScale.DAY
    entry_type: EntryType = EntryType.REFLECTION


class LoginRequest(BaseModel):
    """登录请求"""
    username: str
    password: str


class AuthResponse(BaseModel):
    """认证响应"""
    user: User
    token: str


class SuccessResponse(BaseModel):
    """成功响应格式"""
    success: bool = True
    data: Optional[dict] = None
    message: Optional[str] = None


class ErrorResponse(BaseModel):
    """错误响应格式"""
    success: bool = False
    error: str
    message: str


class PaginationResponse(BaseModel):
    """分页响应格式"""
    success: bool = True
    data: List[dict]
    total: int
    current_page: int
    page_size: int
    total_pages: int
