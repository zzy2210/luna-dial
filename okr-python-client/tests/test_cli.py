"""测试命令行接口"""

import pytest
from click.testing import CliRunner
from unittest.mock import Mock, patch

from okr_client.cli import cli
from okr_client.models import User


class TestCLI:
    """命令行接口测试"""
    
    def setup_method(self):
        """测试设置"""
        self.runner = CliRunner()
    
    @patch('okr_client.cli.get_client')
    def test_me_command(self, mock_get_client):
        """测试 me 命令"""
        # 模拟客户端和用户数据
        mock_client = Mock()
        mock_user = User(
            id="123",
            username="admin",
            email="admin@test.com",
            created_at="2025-07-11T10:00:00Z",
            updated_at="2025-07-11T10:00:00Z"
        )
        mock_client.get_current_user.return_value = mock_user
        mock_get_client.return_value = mock_client
        
        result = self.runner.invoke(cli, ['me'])
        
        assert result.exit_code == 0
        assert "admin" in result.output
    
    @patch('okr_client.cli.get_client')
    def test_task_list_command(self, mock_get_client):
        """测试任务列表命令"""
        # 模拟客户端和任务数据
        mock_client = Mock()
        mock_client.get_tasks.return_value = []
        mock_get_client.return_value = mock_client
        
        result = self.runner.invoke(cli, ['task', 'list'])
        
        assert result.exit_code == 0
        assert "没有找到任务" in result.output
