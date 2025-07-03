"""测试 OKR 客户端"""

import pytest
from unittest.mock import Mock, patch

from okr_client.client import OKRClient, OKRClientError
from okr_client.models import User, Task, TaskType, TaskStatus


class TestOKRClient:
    """OKR 客户端测试"""
    
    def test_init(self):
        """测试客户端初始化"""
        client = OKRClient("http://test.com/api")
        assert client.base_url == "http://test.com/api"
    
    @patch('okr_client.client.requests.Session.request')
    def test_login_success(self, mock_request):
        """测试登录成功"""
        # 模拟成功响应
        mock_response = Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            "success": True,
            "data": {
                "user": {
                    "id": "123",
                    "username": "admin",
                    "email": "admin@test.com",
                    "created_at": "2025-07-11T10:00:00Z",
                    "updated_at": "2025-07-11T10:00:00Z"
                },
                "token": "test_token"
            }
        }
        mock_request.return_value = mock_response
        
        client = OKRClient()
        auth_response = client.login("admin", "password")
        
        assert auth_response.user.username == "admin"
        assert auth_response.token == "test_token"
    
    @patch('okr_client.client.requests.Session.request')
    def test_get_tasks(self, mock_request):
        """测试获取任务列表"""
        # 模拟成功响应
        mock_response = Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            "success": True,
            "data": [
                {
                    "id": "task1",
                    "title": "Test Task",
                    "description": "Test Description",
                    "type": "day",
                    "start_date": "2025-07-11T09:00:00Z",
                    "end_date": "2025-07-11T18:00:00Z",
                    "status": "pending",
                    "score": 5,
                    "parent_id": None,
                    "user_id": "user1",
                    "tags": "test",
                    "created_at": "2025-07-11T10:00:00Z",
                    "updated_at": "2025-07-11T10:00:00Z"
                }
            ]
        }
        mock_request.return_value = mock_response
        
        client = OKRClient()
        tasks = client.get_tasks()
        
        assert len(tasks) == 1
        assert tasks[0].title == "Test Task"
        assert tasks[0].type == TaskType.DAY
