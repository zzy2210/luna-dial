"""客户端便捷方法测试"""

import unittest
from datetime import datetime
from unittest.mock import Mock, patch, MagicMock

from okr_client.client import OKRClient, TaskCreationError
from okr_client.models import Task, TaskType, TaskStatus, TimeScale


class TestOKRClientConvenienceMethods(unittest.TestCase):
    """OKR客户端便捷方法测试"""
    
    def setUp(self):
        """测试前设置"""
        self.client = OKRClient()
        # Mock掉网络请求
        self.client._request = Mock()
    
    @patch('okr_client.client.get_today_range')
    def test_create_today_task(self, mock_get_today_range):
        """测试创建今日任务"""
        # 模拟时间范围
        start = datetime(2025, 7, 11, 0, 0, 0)
        end = datetime(2025, 7, 11, 23, 59, 59)
        mock_get_today_range.return_value = (start, end)
        
        # 模拟API响应
        mock_task_data = {
            "id": "task123",
            "title": "测试任务",
            "description": "测试描述",
            "type": "day",
            "start_date": start,
            "end_date": end,
            "status": "pending",
            "score": 5,
            "parent_id": None,
            "user_id": "user123",
            "tags": None,
            "created_at": datetime.now(),
            "updated_at": datetime.now()
        }
        
        self.client._request.return_value = {"data": mock_task_data}
        
        # 调用方法
        task = self.client.create_today_task("测试任务", "测试描述", 5)
        
        # 验证结果
        self.assertIsInstance(task, Task)
        self.assertEqual(task.title, "测试任务")
        self.assertEqual(task.type, TaskType.DAY)
        
        # 验证API调用
        self.client._request.assert_called_once_with(
            "POST", "/tasks", json=unittest.mock.ANY
        )
    
    @patch('okr_client.client.get_this_week_range')
    def test_create_this_week_task(self, mock_get_this_week_range):
        """测试创建本周任务"""
        start = datetime(2025, 7, 7, 0, 0, 0)  # 周一
        end = datetime(2025, 7, 13, 23, 59, 59)  # 周日
        mock_get_this_week_range.return_value = (start, end)
        
        mock_task_data = {
            "id": "task123",
            "title": "周任务",
            "description": None,
            "type": "week",
            "start_date": start,
            "end_date": end,
            "status": "pending",
            "score": None,
            "parent_id": None,
            "user_id": "user123",
            "tags": None,
            "created_at": datetime.now(),
            "updated_at": datetime.now()
        }
        
        self.client._request.return_value = {"data": mock_task_data}
        
        task = self.client.create_this_week_task("周任务")
        
        self.assertEqual(task.title, "周任务")
        self.assertEqual(task.type, TaskType.WEEK)
        self.assertIsNone(task.score)
    
    @patch('okr_client.client.get_quarter_range')
    def test_create_quarter_task(self, mock_get_quarter_range):
        """测试创建指定季度任务"""
        start = datetime(2024, 10, 1, 0, 0, 0)
        end = datetime(2024, 12, 31, 23, 59, 59)
        mock_get_quarter_range.return_value = (start, end)
        
        mock_task_data = {
            "id": "task123",
            "title": "Q4目标",
            "description": "第四季度目标",
            "type": "quarter",
            "start_date": start,
            "end_date": end,
            "status": "pending",
            "score": 8,
            "parent_id": None,
            "user_id": "user123",
            "tags": None,
            "created_at": datetime.now(),
            "updated_at": datetime.now()
        }
        
        self.client._request.return_value = {"data": mock_task_data}
        
        task = self.client.create_quarter_task("Q4目标", 2024, 4, "第四季度目标", 8)
        
        self.assertEqual(task.title, "Q4目标")
        self.assertEqual(task.type, TaskType.QUARTER)
        self.assertEqual(task.score, 8)
        
        # 验证调用了正确的时间范围函数
        mock_get_quarter_range.assert_called_once_with(2024, 4)
    
    def test_create_task_with_invalid_params(self):
        """测试创建任务时参数错误处理"""
        # 模拟get_quarter_range抛出异常
        with patch('okr_client.client.get_quarter_range', side_effect=ValueError("季度必须在1-4之间")):
            with self.assertRaises(TaskCreationError):
                self.client.create_quarter_task("错误任务", 2024, 5)
    
    def test_get_plan_view(self):
        """测试获取计划视图"""
        mock_response_data = {
            "tasks": [],
            "journals": [],
            "time_range": {
                "start": datetime(2024, 10, 1),
                "end": datetime(2024, 12, 31)
            },
            "stats": {
                "total_tasks": 10,
                "completed_tasks": 5,
                "in_progress_tasks": 3,
                "pending_tasks": 2,
                "total_score": 50,
                "completed_score": 25
            }
        }
        
        self.client._request.return_value = {"data": mock_response_data}
        
        # 测试有效的时间参考
        with patch('okr_client.client.validate_time_reference', return_value=True):
            result = self.client.get_plan_view(TimeScale.QUARTER, "2024-Q4")
            
            self.client._request.assert_called_once_with(
                'GET', '/plan', params={'scale': 'quarter', 'time_ref': '2024-Q4'}
            )
        
        # 测试无效的时间参考
        with patch('okr_client.client.validate_time_reference', return_value=False):
            with self.assertRaises(Exception):  # 应该抛出PlanViewError
                self.client.get_plan_view(TimeScale.QUARTER, "invalid")
    
    def test_get_score_trend(self):
        """测试获取分数趋势"""
        mock_response_data = {
            "labels": ["2025-07-01", "2025-07-02", "2025-07-03"],
            "scores": [5, 8, 6],
            "counts": [2, 3, 2],
            "scale": "day",
            "time_ref": "2025-07",
            "time_range": {
                "start": datetime(2025, 7, 1),
                "end": datetime(2025, 7, 31)
            },
            "summary": {
                "total_score": 19,
                "total_tasks": 7,
                "average_score": 2.71,
                "average_task_count": 2.33,
                "max_score": 8,
                "max_tasks": 3,
                "min_score": 5,
                "min_tasks": 2
            }
        }
        
        self.client._request.return_value = {"data": mock_response_data}
        
        with patch('okr_client.client.validate_time_reference', return_value=True):
            result = self.client.get_score_trend(TimeScale.MONTH, "2025-07")
            
            self.client._request.assert_called_once_with(
                'GET', '/stats/score-trend-ref', params={'scale': 'month', 'time_ref': '2025-07'}
            )
    
    def test_convenience_plan_methods(self):
        """测试便捷计划视图方法"""
        self.client.get_plan_view = Mock()
        
        # 测试季度便捷方法
        self.client.get_plan_view_for_quarter(2024, 4)
        self.client.get_plan_view.assert_called_with(TimeScale.QUARTER, "2024-Q4")
        
        # 测试月份便捷方法
        self.client.get_plan_view_for_month(2025, 7)
        self.client.get_plan_view.assert_called_with(TimeScale.MONTH, "2025-07")
    
    def test_convenience_trend_methods(self):
        """测试便捷趋势方法"""
        self.client.get_score_trend = Mock()
        
        # 测试月度趋势便捷方法
        self.client.get_monthly_score_trend(2025, 7)
        self.client.get_score_trend.assert_called_with(TimeScale.MONTH, "2025-07")
        
        # 测试季度趋势便捷方法
        self.client.get_quarterly_score_trend(2024, 4)
        self.client.get_score_trend.assert_called_with(TimeScale.QUARTER, "2024-Q4")


if __name__ == '__main__':
    unittest.main()
