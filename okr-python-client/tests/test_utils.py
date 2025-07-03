"""便捷任务创建功能测试"""

import unittest
from datetime import datetime, timedelta
from unittest.mock import Mock, patch

from okr_client.utils import (
    get_today_range, get_this_week_range, get_this_month_range,
    get_this_quarter_range, get_this_year_range,
    get_quarter_range, get_month_range, get_week_range,
    validate_time_reference, parse_time_reference_to_range
)
from okr_client.models import TimeScale


class TestTimeUtilFunctions(unittest.TestCase):
    """时间计算函数测试"""
    
    def test_get_today_range(self):
        """测试获取今日时间范围"""
        start, end = get_today_range()
        
        # 检查是否是同一天
        self.assertEqual(start.date(), end.date())
        
        # 检查开始时间是否是当天0点
        self.assertEqual(start.hour, 0)
        self.assertEqual(start.minute, 0)
        self.assertEqual(start.second, 0)
        self.assertEqual(start.microsecond, 0)
        
        # 检查结束时间是否是当天23:59:59
        self.assertEqual(end.hour, 23)
        self.assertEqual(end.minute, 59)
        self.assertEqual(end.second, 59)
    
    def test_get_this_week_range(self):
        """测试获取本周时间范围"""
        start, end = get_this_week_range()
        
        # 检查开始时间是否是周一
        self.assertEqual(start.weekday(), 0)  # 0 = 周一
        
        # 检查结束时间是否是周日
        self.assertEqual(end.weekday(), 6)  # 6 = 周日
        
        # 检查时间跨度是否是7天
        self.assertEqual((end.date() - start.date()).days, 6)
    
    def test_get_this_month_range(self):
        """测试获取本月时间范围"""
        start, end = get_this_month_range()
        
        # 检查开始时间是否是月初
        self.assertEqual(start.day, 1)
        
        # 检查是否是同一个月
        self.assertEqual(start.month, end.month)
        self.assertEqual(start.year, end.year)
    
    def test_get_this_quarter_range(self):
        """测试获取本季度时间范围"""
        start, end = get_this_quarter_range()
        
        # 检查时间跨度大致是3个月
        months_diff = (end.year - start.year) * 12 + (end.month - start.month) + 1
        self.assertGreaterEqual(months_diff, 3)
        self.assertLessEqual(months_diff, 3)
    
    def test_get_this_year_range(self):
        """测试获取本年时间范围"""
        start, end = get_this_year_range()
        
        # 检查是否是同一年
        self.assertEqual(start.year, end.year)
        
        # 检查开始时间是否是1月1日
        self.assertEqual(start.month, 1)
        self.assertEqual(start.day, 1)
        
        # 检查结束时间是否是12月31日
        self.assertEqual(end.month, 12)
        self.assertEqual(end.day, 31)
    
    def test_get_quarter_range(self):
        """测试获取指定季度时间范围"""
        # 测试2024年第4季度
        start, end = get_quarter_range(2024, 4)
        
        self.assertEqual(start.year, 2024)
        self.assertEqual(start.month, 10)  # Q4从10月开始
        self.assertEqual(start.day, 1)
        
        self.assertEqual(end.year, 2024)
        self.assertEqual(end.month, 12)  # Q4在12月结束
        self.assertEqual(end.day, 31)
        
        # 测试无效季度
        with self.assertRaises(ValueError):
            get_quarter_range(2024, 5)
    
    def test_get_month_range(self):
        """测试获取指定月份时间范围"""
        # 测试2025年7月
        start, end = get_month_range(2025, 7)
        
        self.assertEqual(start.year, 2025)
        self.assertEqual(start.month, 7)
        self.assertEqual(start.day, 1)
        
        self.assertEqual(end.year, 2025)
        self.assertEqual(end.month, 7)
        self.assertEqual(end.day, 31)
        
        # 测试无效月份
        with self.assertRaises(ValueError):
            get_month_range(2025, 13)
    
    def test_get_week_range(self):
        """测试获取指定周时间范围"""
        # 测试基本功能
        start, end = get_week_range(2025, 15)
        
        # 检查时间跨度是7天
        self.assertEqual((end.date() - start.date()).days, 6)
        
        # 检查开始时间是周一
        self.assertEqual(start.weekday(), 0)
        
        # 检查结束时间是周日
        self.assertEqual(end.weekday(), 6)
        
        # 测试无效周数
        with self.assertRaises(ValueError):
            get_week_range(2025, 54)
    
    def test_validate_time_reference(self):
        """测试时间参考格式验证"""
        # 测试年份格式
        self.assertTrue(validate_time_reference("2024", TimeScale.YEAR))
        self.assertFalse(validate_time_reference("24", TimeScale.YEAR))
        
        # 测试季度格式
        self.assertTrue(validate_time_reference("2024-Q4", TimeScale.QUARTER))
        self.assertFalse(validate_time_reference("2024-Q5", TimeScale.QUARTER))
        
        # 测试月份格式
        self.assertTrue(validate_time_reference("2025-07", TimeScale.MONTH))
        self.assertFalse(validate_time_reference("2025-13", TimeScale.MONTH))
        
        # 测试周格式
        self.assertTrue(validate_time_reference("2025-W15", TimeScale.WEEK))
        self.assertFalse(validate_time_reference("2025-W54", TimeScale.WEEK))
        
        # 测试日期格式
        self.assertTrue(validate_time_reference("2025-07-15", TimeScale.DAY))
        self.assertFalse(validate_time_reference("2025-13-15", TimeScale.DAY))
    
    def test_parse_time_reference_to_range(self):
        """测试时间参考解析为时间范围"""
        # 测试年份解析
        start, end = parse_time_reference_to_range("2024", TimeScale.YEAR)
        self.assertEqual(start.year, 2024)
        self.assertEqual(start.month, 1)
        self.assertEqual(start.day, 1)
        self.assertEqual(end.year, 2024)
        self.assertEqual(end.month, 12)
        self.assertEqual(end.day, 31)
        
        # 测试季度解析
        start, end = parse_time_reference_to_range("2024-Q4", TimeScale.QUARTER)
        self.assertEqual(start.month, 10)
        self.assertEqual(end.month, 12)
        
        # 测试月份解析
        start, end = parse_time_reference_to_range("2025-07", TimeScale.MONTH)
        self.assertEqual(start.month, 7)
        self.assertEqual(end.month, 7)
        
        # 测试日期解析
        start, end = parse_time_reference_to_range("2025-07-15", TimeScale.DAY)
        self.assertEqual(start.month, 7)
        self.assertEqual(start.day, 15)
        self.assertEqual(end.month, 7)
        self.assertEqual(end.day, 15)
        
        # 测试无效格式
        with self.assertRaises(ValueError):
            parse_time_reference_to_range("invalid", TimeScale.YEAR)


if __name__ == '__main__':
    unittest.main()
