"""OKR Python Client - 命令行工具和API客户端"""

__version__ = "1.1.0"
__author__ = "OKR Team"
__description__ = "Python client for OKR management system"

# 导入核心组件
from .client import OKRClient
from .models import (
    # 枚举类型
    TaskType, TaskStatus, TimeScale, EntryType,
    
    # 基础模型
    User, Task, JournalEntry,
    
    # 请求模型
    TaskRequest, JournalRequest, LoginRequest,
    PlanRequest, ScoreTrendRequest,
    
    # 响应模型
    AuthResponse, SuccessResponse, ErrorResponse, PaginationResponse,
    PlanResponse, ScoreTrendResponse,
    
    # 新增模型
    TimeRange, TaskTree, PlanStats, TrendSummary
)

__all__ = [
    # 客户端
    'OKRClient',
    
    # 枚举类型
    'TaskType', 'TaskStatus', 'TimeScale', 'EntryType',
    
    # 基础模型
    'User', 'Task', 'JournalEntry',
    
    # 请求模型
    'TaskRequest', 'JournalRequest', 'LoginRequest',
    'PlanRequest', 'ScoreTrendRequest',
    
    # 响应模型
    'AuthResponse', 'SuccessResponse', 'ErrorResponse', 'PaginationResponse',
    'PlanResponse', 'ScoreTrendResponse',
    
    # 新增模型
    'TimeRange', 'TaskTree', 'PlanStats', 'TrendSummary'
]
