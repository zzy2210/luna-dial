"""OKR API 客户端"""

import json
import os
from datetime import datetime
from pathlib import Path
from typing import Optional, List, Dict, Any
from urllib.parse import urljoin

import requests
from requests import Response

from .models import (
    User, Task, JournalEntry, TaskRequest, JournalRequest,
    LoginRequest, AuthResponse, SuccessResponse, ErrorResponse,
    TaskType, TaskStatus, TimeScale, EntryType,
    # 新增模型
    PlanRequest, PlanResponse, ScoreTrendRequest, ScoreTrendResponse,
    TimeRange, TaskTree, PlanStats, TrendSummary
)
from .utils import (
    validate_time_reference, parse_time_reference_to_range,
    get_today_range, get_this_week_range, get_this_month_range,
    get_this_quarter_range, get_this_year_range,
    get_quarter_range, get_month_range, get_week_range
)


class OKRClientError(Exception):
    """OKR 客户端异常"""
    pass


class PlanViewError(OKRClientError):
    """计划视图相关错误"""
    pass


class ScoreTrendError(OKRClientError):
    """分数趋势相关错误"""
    pass


class TaskCreationError(OKRClientError):
    """任务创建相关错误"""
    pass


class OKRClient:
    """OKR API 客户端"""
    
    def __init__(self, base_url: str = None, config_path: str = None):
        self.base_url = base_url or os.getenv("OKR_API_BASE_URL", "http://localhost:8081/api")
        self.config_path = config_path or os.path.expanduser("~/.okr/config")
        self.session = requests.Session()
        
        # 加载已保存的认证信息
        self._load_auth()
    
    def _load_auth(self):
        """加载保存的认证信息"""
        try:
            if os.path.exists(self.config_path):
                with open(self.config_path, 'r') as f:
                    config = json.load(f)
                    token = config.get('token')
                    if token:
                        self.session.headers.update({
                            'Authorization': f'Bearer {token}'
                        })
        except Exception:
            # 忽略配置加载错误
            pass
    
    def _save_auth(self, token: str):
        """保存认证信息"""
        os.makedirs(os.path.dirname(self.config_path), exist_ok=True)
        config = {'token': token}
        with open(self.config_path, 'w') as f:
            json.dump(config, f)
        
        # 更新会话头
        self.session.headers.update({
            'Authorization': f'Bearer {token}'
        })
    
    def _clear_auth(self):
        """清除认证信息"""
        if os.path.exists(self.config_path):
            os.remove(self.config_path)
        
        # 清除会话头
        self.session.headers.pop('Authorization', None)
    
    def _request(self, method: str, endpoint: str, **kwargs) -> Dict[Any, Any]:
        """发送 HTTP 请求"""
        url = urljoin(self.base_url + "/", endpoint.lstrip("/"))
        
        try:
            response: Response = self.session.request(method, url, **kwargs)
            response.raise_for_status()
            
            return response.json()
        except requests.exceptions.HTTPError as e:
            try:
                error_data = response.json()
                raise OKRClientError(f"API Error: {error_data.get('message', str(e))}")
            except json.JSONDecodeError:
                raise OKRClientError(f"HTTP Error: {e}")
        except requests.exceptions.RequestException as e:
            raise OKRClientError(f"Request Error: {e}")
    
    # 认证相关方法
    def login(self, username: str, password: str) -> AuthResponse:
        """用户登录"""
        request_data = LoginRequest(username=username, password=password)
        response = self._request("POST", "/auth/login", json=request_data.model_dump())
        
        auth_response = AuthResponse(**response["data"])
        self._save_auth(auth_response.token)
        
        return auth_response
    
    def logout(self):
        """用户登出"""
        try:
            self._request("POST", "/users/logout")
        finally:
            self._clear_auth()
    
    def get_current_user(self) -> User:
        """获取当前用户信息"""
        response = self._request("GET", "/users/me")
        return User(**response["data"])
    
    # 任务相关方法
    def get_tasks(self, 
                  task_type: Optional[TaskType] = None,
                  start_date: Optional[datetime] = None,
                  end_date: Optional[datetime] = None,
                  status: Optional[TaskStatus] = None,
                  page: int = 1,
                  page_size: int = 20) -> List[Task]:
        """获取任务列表"""
        params = {
            "page": page,
            "page_size": page_size
        }
        
        if task_type:
            params["type"] = task_type.value
        if start_date:
            params["start_date"] = start_date.isoformat()
        if end_date:
            params["end_date"] = end_date.isoformat()
        if status:
            params["status"] = status.value

        response = self._request("GET", "/tasks", params=params)
        # 修正：只遍历 data["tasks"]
        return [Task(**task) for task in response["data"]["tasks"]]
    
    def get_task(self, task_id: str) -> Task:
        """获取单个任务"""
        response = self._request("GET", f"/tasks/{task_id}")
        return Task(**response["data"])
    
    def create_task(self, task_request: TaskRequest) -> Task:
        """创建任务"""
        # 只在明确需要时才传 parent_id
        data = task_request.model_dump(mode="json", exclude_none=True)
        if "parent_id" in data and data["parent_id"] is None:
            del data["parent_id"]
        response = self._request("POST", "/tasks", json=data)
        return Task(**response["data"])
    
    def update_task(self, task_id: str, task_request: TaskRequest) -> Task:
        """更新任务"""
        response = self._request("PUT", f"/tasks/{task_id}", json=task_request.model_dump(mode="json"))
        return Task(**response["data"])
    
    def delete_task(self, task_id: str):
        """删除任务"""
        self._request("DELETE", f"/tasks/{task_id}")
    
    def complete_task(self, task_id: str) -> Task:
        """完成任务"""
        # 先获取任务信息
        task = self.get_task(task_id)
        
        # 更新状态为已完成
        task_request = TaskRequest(
            title=task.title,
            description=task.description,
            type=task.type,
            start_date=task.start_date,
            end_date=task.end_date,
            status=TaskStatus.COMPLETED,
            score=task.score,
            parent_id=task.parent_id,
            tags=task.tags
        )
        
        return self.update_task(task_id, task_request)
    
    def update_task_score(self, task_id: str, score: int) -> Task:
        """更新任务分数"""
        response = self._request("PUT", f"/tasks/{task_id}/score", json={"score": score})
        return Task(**response["data"])
    
    def get_task_children(self, task_id: str) -> List[Task]:
        """获取子任务"""
        response = self._request("GET", f"/tasks/{task_id}/children")
        return [Task(**task) for task in response["data"]]
    
    # 新增：计划视图方法
    def get_plan_view(self, scale: TimeScale, time_ref: str) -> PlanResponse:
        """获取计划视图
        
        Args:
            scale: 时间尺度
            time_ref: 时间参考字符串
            
        Returns:
            PlanResponse: 计划视图响应
            
        Raises:
            PlanViewError: 计划视图相关错误
        """
        try:
            # 验证时间参考格式
            if not validate_time_reference(time_ref, scale):
                raise PlanViewError(f"时间参考格式错误: {time_ref}，期望格式为 {scale.value}")
            params = {
                'scale': scale.value,
                'time_ref': time_ref
            }
            response = self._request('GET', '/plan', params=params)
            data = response['data']
            if data.get('tasks') is None:
                data['tasks'] = []
            return PlanResponse(**data)
        except OKRClientError as e:
            raise PlanViewError(f"获取计划视图失败: {e}")
    
    def get_score_trend(self, scale: TimeScale, time_ref: str) -> ScoreTrendResponse:
        """获取分数趋势统计
        
        Args:
            scale: 统计尺度
            time_ref: 时间参考字符串
            
        Returns:
            ScoreTrendResponse: 分数趋势响应
            
        Raises:
            ScoreTrendError: 分数趋势相关错误
        """
        try:
            # 验证时间参考格式
            if not validate_time_reference(time_ref, scale):
                raise ScoreTrendError(f"时间参考格式错误: {time_ref}，期望格式为 {scale.value}")
            
            params = {
                'scale': scale.value,
                'time_ref': time_ref
            }
            
            response = self._request('GET', '/stats/score-trend-ref', params=params)
            return ScoreTrendResponse(**response['data'])
            
        except OKRClientError as e:
            raise ScoreTrendError(f"获取分数趋势失败: {e}")
    
    # 新增：便捷计划视图方法
    def get_plan_view_for_quarter(self, year: int, quarter: int) -> PlanResponse:
        """获取指定季度的计划视图（便捷方法）
        
        Args:
            year: 年份
            quarter: 季度 (1-4)
            
        Returns:
            PlanResponse: 计划视图响应
        """
        time_ref = f"{year}-Q{quarter}"
        return self.get_plan_view(TimeScale.QUARTER, time_ref)
    
    def get_plan_view_for_month(self, year: int, month: int) -> PlanResponse:
        """获取指定月份的计划视图（便捷方法）
        
        Args:
            year: 年份
            month: 月份 (1-12)
            
        Returns:
            PlanResponse: 计划视图响应
        """
        time_ref = f"{year}-{month:02d}"
        return self.get_plan_view(TimeScale.MONTH, time_ref)
    
    def get_monthly_score_trend(self, year: int, month: int) -> ScoreTrendResponse:
        """获取月度分数趋势（便捷方法）
        
        Args:
            year: 年份
            month: 月份 (1-12)
            
        Returns:
            ScoreTrendResponse: 分数趋势响应
        """
        time_ref = f"{year}-{month:02d}"
        return self.get_score_trend(TimeScale.MONTH, time_ref)
    
    def get_quarterly_score_trend(self, year: int, quarter: int) -> ScoreTrendResponse:
        """获取季度分数趋势（便捷方法）
        
        Args:
            year: 年份
            quarter: 季度 (1-4)
            
        Returns:
            ScoreTrendResponse: 分数趋势响应
        """
        time_ref = f"{year}-Q{quarter}"
        return self.get_score_trend(TimeScale.QUARTER, time_ref)
    
    # 新增：便捷任务创建方法
    def _create_task_with_time_range(self, title: str, start: datetime, end: datetime, 
                                   task_type: TaskType, description: str = None, 
                                   score: int = None) -> Task:
        """创建指定时间范围的任务
        
        Args:
            title: 任务标题
            start: 开始时间
            end: 结束时间
            task_type: 任务类型
            description: 任务描述
            score: 任务分数
            
        Returns:
            Task: 创建的任务
        """
        task_request = TaskRequest(
            title=title,
            description=description,
            type=task_type,
            start_date=start.isoformat() + 'Z' if isinstance(start, datetime) and start.tzinfo is None else start,
            end_date=end.isoformat() + 'Z' if isinstance(end, datetime) and end.tzinfo is None else end,
            status=TaskStatus.PENDING,
            score=None,  # 始终为 None
        )
        
        return self.create_task(task_request)
    
    def create_today_task(self, title: str, description: str = None, score: int = None) -> Task:
        """创建今日任务
        
        Args:
            title: 任务标题
            description: 任务描述
            score: 任务分数
            
        Returns:
            Task: 创建的任务
        """
        try:
            start, end = get_today_range()
            return self._create_task_with_time_range(
                title, start, end, TaskType.DAY, description, score
            )
        except Exception as e:
            raise TaskCreationError(f"创建今日任务失败: {e}")
    
    def create_this_week_task(self, title: str, description: str = None, score: int = None) -> Task:
        """创建本周任务
        
        Args:
            title: 任务标题
            description: 任务描述
            score: 任务分数
            
        Returns:
            Task: 创建的任务
        """
        try:
            start, end = get_this_week_range()
            return self._create_task_with_time_range(
                title, start, end, TaskType.WEEK, description, score
            )
        except Exception as e:
            raise TaskCreationError(f"创建本周任务失败: {e}")
    
    def create_this_month_task(self, title: str, description: str = None, score: int = None) -> Task:
        """创建本月任务
        
        Args:
            title: 任务标题
            description: 任务描述
            score: 任务分数
            
        Returns:
            Task: 创建的任务
        """
        try:
            start, end = get_this_month_range()
            return self._create_task_with_time_range(
                title, start, end, TaskType.MONTH, description, score
            )
        except Exception as e:
            raise TaskCreationError(f"创建本月任务失败: {e}")
    
    def create_this_quarter_task(self, title: str, description: str = None, score: int = None) -> Task:
        """创建本季度任务
        
        Args:
            title: 任务标题
            description: 任务描述
            score: 任务分数
            
        Returns:
            Task: 创建的任务
        """
        try:
            start, end = get_this_quarter_range()
            return self._create_task_with_time_range(
                title, start, end, TaskType.QUARTER, description, score
            )
        except Exception as e:
            raise TaskCreationError(f"创建本季度任务失败: {e}")
    
    def create_this_year_task(self, title: str, description: str = None, score: int = None) -> Task:
        """创建本年任务
        
        Args:
            title: 任务标题
            description: 任务描述
            score: 任务分数
            
        Returns:
            Task: 创建的任务
        """
        try:
            start, end = get_this_year_range()
            return self._create_task_with_time_range(
                title, start, end, TaskType.YEAR, description, score
            )
        except Exception as e:
            raise TaskCreationError(f"创建本年任务失败: {e}")
    
    def create_quarter_task(self, title: str, year: int, quarter: int, 
                          description: str = None, score: int = None) -> Task:
        """创建指定季度任务
        
        Args:
            title: 任务标题
            year: 年份
            quarter: 季度 (1-4)
            description: 任务描述
            score: 任务分数
            
        Returns:
            Task: 创建的任务
        """
        try:
            start, end = get_quarter_range(year, quarter)
            return self._create_task_with_time_range(
                title, start, end, TaskType.QUARTER, description, score
            )
        except Exception as e:
            raise TaskCreationError(f"创建指定季度任务失败: {e}")
    
    def create_month_task(self, title: str, year: int, month: int, 
                        description: str = None, score: int = None) -> Task:
        """创建指定月份任务
        
        Args:
            title: 任务标题
            year: 年份
            month: 月份 (1-12)
            description: 任务描述
            score: 任务分数
            
        Returns:
            Task: 创建的任务
        """
        try:
            start, end = get_month_range(year, month)
            return self._create_task_with_time_range(
                title, start, end, TaskType.MONTH, description, score
            )
        except Exception as e:
            raise TaskCreationError(f"创建指定月份任务失败: {e}")
    
    def create_week_task(self, title: str, year: int, week: int, 
                       description: str = None, score: int = None) -> Task:
        """创建指定周任务
        
        Args:
            title: 任务标题
            year: 年份
            week: 周数 (1-53)
            description: 任务描述
            score: 任务分数
            
        Returns:
            Task: 创建的任务
        """
        try:
            start, end = get_week_range(year, week)
            return self._create_task_with_time_range(
                title, start, end, TaskType.WEEK, description, score
            )
        except Exception as e:
            raise TaskCreationError(f"创建指定周任务失败: {e}")
    
    # 日志相关方法
    def get_journals(self,
                     time_scale: Optional[TimeScale] = None,
                     start_time: Optional[datetime] = None,
                     end_time: Optional[datetime] = None,
                     page: int = 1,
                     page_size: int = 20) -> List[JournalEntry]:
        """获取日志列表"""
        params = {
            "page": page,
            "page_size": page_size
        }

        if time_scale:
            params["time_scale"] = time_scale.value
        if start_time:
            params["start_time"] = start_time.isoformat()
        if end_time:
            params["end_time"] = end_time.isoformat()

        response = self._request("GET", "/journals", params=params)

        # 正确处理 journals 列表
        return [JournalEntry(**journal) for journal in response["data"]["journals"]]
    
    def get_journal(self, journal_id: str) -> JournalEntry:
        """获取单个日志"""
        response = self._request("GET", f"/journals/{journal_id}")
        return JournalEntry(**response["data"])
    
    def create_journal(self, journal_request: JournalRequest) -> JournalEntry:
        """创建日志"""
        response = self._request("POST", "/journals", json=journal_request.model_dump(mode="json"))
        return JournalEntry(**response["data"])
    
    def update_journal(self, journal_id: str, journal_request: JournalRequest) -> JournalEntry:
        """更新日志"""
        response = self._request("PUT", f"/journals/{journal_id}", json=journal_request.model_dump(mode="json"))
        return JournalEntry(**response["data"])
    
    def delete_journal(self, journal_id: str):
        """删除日志"""
        self._request("DELETE", f"/journals/{journal_id}")
    
    def get_journals_by_time(self, time_reference: str, time_scale: TimeScale) -> List[JournalEntry]:
        """按时间查询日志"""
        params = {
            "time_reference": time_reference,
            "time_scale": time_scale.value
        }
        response = self._request("GET", "/journals/by-time", params=params)
        return [JournalEntry(**journal) for journal in response["data"]]
    
    # 基于当前时间的计划视图便捷方法
    def get_plan_view_for_today(self) -> PlanResponse:
        """获取今日计划视图
        
        Returns:
            PlanResponse: 计划视图响应
        """
        from datetime import datetime
        today = datetime.now().strftime('%Y-%m-%d')
        return self.get_plan_view(TimeScale.DAY, today)
    
    def get_plan_view_for_this_week(self) -> PlanResponse:
        """获取本周计划视图（ISO周编号）
        Returns:
            PlanResponse: 计划视图响应
        """
        from datetime import datetime
        now = datetime.now()
        year, week, _ = now.isocalendar()
        time_ref = f"{year}-W{week:02d}"
        return self.get_plan_view(TimeScale.WEEK, time_ref)
    
    def get_plan_view_for_this_month(self) -> PlanResponse:
        """获取本月计划视图
        
        Returns:
            PlanResponse: 计划视图响应
        """
        from datetime import datetime
        this_month = datetime.now().strftime('%Y-%m')
        return self.get_plan_view(TimeScale.MONTH, this_month)
    
    def get_plan_view_for_this_quarter(self) -> PlanResponse:
        """获取本季度计划视图
        
        Returns:
            PlanResponse: 计划视图响应
        """
        from datetime import datetime
        now = datetime.now()
        quarter = (now.month - 1) // 3 + 1
        time_ref = f"{now.year}-Q{quarter}"
        return self.get_plan_view(TimeScale.QUARTER, time_ref)
    
    def get_plan_view_for_this_year(self) -> PlanResponse:
        """获取本年计划视图
        
        Returns:
            PlanResponse: 计划视图响应
        """
        from datetime import datetime
        this_year = str(datetime.now().year)
        return self.get_plan_view(TimeScale.YEAR, this_year)
    
    # 基于当前时间的分数趋势便捷方法
    def get_score_trend_for_today(self) -> ScoreTrendResponse:
        """获取今日分数趋势
        
        Returns:
            ScoreTrendResponse: 分数趋势响应
        """
        from datetime import datetime
        today = datetime.now().strftime('%Y-%m-%d')
        return self.get_score_trend(TimeScale.DAY, today)
    
    def get_score_trend_for_this_week(self) -> ScoreTrendResponse:
        """获取本周分数趋势（ISO周编号）
        Returns:
            ScoreTrendResponse: 分数趋势响应
        """
        from datetime import datetime
        now = datetime.now()
        year, week, _ = now.isocalendar()
        time_ref = f"{year}-W{week:02d}"
        return self.get_score_trend(TimeScale.WEEK, time_ref)
    
    def get_score_trend_for_this_month(self) -> ScoreTrendResponse:
        """获取本月分数趋势
        
        Returns:
            ScoreTrendResponse: 分数趋势响应
        """
        from datetime import datetime
        this_month = datetime.now().strftime('%Y-%m')
        return self.get_score_trend(TimeScale.MONTH, this_month)
    
    def get_score_trend_for_this_quarter(self) -> ScoreTrendResponse:
        """获取本季度分数趋势
        
        Returns:
            ScoreTrendResponse: 分数趋势响应
        """
        from datetime import datetime
        now = datetime.now()
        quarter = (now.month - 1) // 3 + 1
        time_ref = f"{now.year}-Q{quarter}"
        return self.get_score_trend(TimeScale.QUARTER, time_ref)
    
    def get_score_trend_for_this_year(self) -> ScoreTrendResponse:
        """获取本年分数趋势
        
        Returns:
            ScoreTrendResponse: 分数趋势响应
        """
        from datetime import datetime
        this_year = str(datetime.now().year)
        return self.get_score_trend(TimeScale.YEAR, this_year)
