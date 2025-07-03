"""时间计算工具函数"""

import re
from datetime import datetime, timedelta
from typing import Tuple, Dict
from okr_client.models import TimeScale


# 时间参考格式验证模式
TIME_PATTERNS = {
    'year': r'^\d{4}$',
    'quarter': r'^\d{4}-Q[1-4]$',  
    'month': r'^\d{4}-(0[1-9]|1[0-2])$',
    'week': r'^\d{4}-W([1-4]\d|5[0-3]|\d)$',  # ISO周编号
    'day': r'^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$'
}


def validate_time_reference(time_ref: str, scale: TimeScale) -> bool:
    """验证时间参考格式是否正确
    
    Args:
        time_ref: 时间参考字符串
        scale: 时间尺度
        
    Returns:
        bool: 格式是否正确
    """
    pattern = TIME_PATTERNS.get(scale.value)
    if not pattern:
        return False
    return bool(re.match(pattern, time_ref))


def get_today_range() -> Tuple[datetime, datetime]:
    """获取今日的开始和结束时间
    
    Returns:
        Tuple[datetime, datetime]: (开始时间, 结束时间)
    """
    now = datetime.now()
    start = now.replace(hour=0, minute=0, second=0, microsecond=0)
    end = now.replace(hour=23, minute=59, second=59, microsecond=999999)
    return start, end


def get_this_week_range() -> Tuple[datetime, datetime]:
    """获取本周的开始和结束时间（ISO周：周一到周日）
    Returns:
        Tuple[datetime, datetime]: (开始时间, 结束时间)
    """
    now = datetime.now()
    # ISO周一为一周第一天
    days_since_monday = now.weekday()  # 0=周一
    start_of_week = now - timedelta(days=days_since_monday)
    start = start_of_week.replace(hour=0, minute=0, second=0, microsecond=0)
    end_of_week = start_of_week + timedelta(days=6)
    end = end_of_week.replace(hour=23, minute=59, second=59, microsecond=999999)
    return start, end



def get_this_month_range() -> Tuple[datetime, datetime]:
    """获取本月的开始和结束时间
    
    Returns:
        Tuple[datetime, datetime]: (开始时间, 结束时间)
    """
    now = datetime.now()
    start = now.replace(day=1, hour=0, minute=0, second=0, microsecond=0)
    
    # 计算下月第一天，然后减去一天得到本月最后一天
    if now.month == 12:
        next_month = now.replace(year=now.year + 1, month=1, day=1)
    else:
        next_month = now.replace(month=now.month + 1, day=1)
    
    last_day = next_month - timedelta(days=1)
    end = last_day.replace(hour=23, minute=59, second=59, microsecond=999999)
    
    return start, end


def get_this_quarter_range() -> Tuple[datetime, datetime]:
    """获取本季度的开始和结束时间
    
    Returns:
        Tuple[datetime, datetime]: (开始时间, 结束时间)
    """
    now = datetime.now()
    quarter = (now.month - 1) // 3 + 1
    return get_quarter_range(now.year, quarter)


def get_this_year_range() -> Tuple[datetime, datetime]:
    """获取本年的开始和结束时间
    
    Returns:
        Tuple[datetime, datetime]: (开始时间, 结束时间)
    """
    now = datetime.now()
    start = now.replace(month=1, day=1, hour=0, minute=0, second=0, microsecond=0)
    end = now.replace(month=12, day=31, hour=23, minute=59, second=59, microsecond=999999)
    return start, end


def get_quarter_range(year: int, quarter: int) -> Tuple[datetime, datetime]:
    """获取指定季度的时间范围
    
    Args:
        year: 年份
        quarter: 季度 (1-4)
        
    Returns:
        Tuple[datetime, datetime]: (开始时间, 结束时间)
        
    Raises:
        ValueError: 季度参数不在1-4范围内
    """
    if quarter not in [1, 2, 3, 4]:
        raise ValueError(f"季度必须在1-4之间，得到: {quarter}")
    
    # 季度起始月份
    quarter_start_months = {1: 1, 2: 4, 3: 7, 4: 10}
    start_month = quarter_start_months[quarter]
    
    start = datetime(year, start_month, 1, 0, 0, 0, 0)
    
    # 计算季度结束月份
    end_month = start_month + 2
    
    # 计算季度结束日期
    if end_month == 12:
        next_quarter_start = datetime(year + 1, 1, 1, 0, 0, 0, 0)
    else:
        next_quarter_start = datetime(year, end_month + 1, 1, 0, 0, 0, 0)
    
    end = next_quarter_start - timedelta(microseconds=1)
    
    return start, end


def get_month_range(year: int, month: int) -> Tuple[datetime, datetime]:
    """获取指定月份的时间范围
    
    Args:
        year: 年份
        month: 月份 (1-12)
        
    Returns:
        Tuple[datetime, datetime]: (开始时间, 结束时间)
        
    Raises:
        ValueError: 月份参数不在1-12范围内
    """
    if month not in range(1, 13):
        raise ValueError(f"月份必须在1-12之间，得到: {month}")
    
    start = datetime(year, month, 1, 0, 0, 0, 0)
    
    # 计算下月第一天，然后减去一微秒得到本月最后时刻
    if month == 12:
        next_month_start = datetime(year + 1, 1, 1, 0, 0, 0, 0)
    else:
        next_month_start = datetime(year, month + 1, 1, 0, 0, 0, 0)
    
    end = next_month_start - timedelta(microseconds=1)
    
    return start, end


def get_week_range(year: int, week: int) -> Tuple[datetime, datetime]:
    """获取指定周的时间范围（基于ISO周）
    
    Args:
        year: 年份
        week: 周数 (1-53)
        
    Returns:
        Tuple[datetime, datetime]: (开始时间, 结束时间)
        
    Raises:
        ValueError: 周数参数不在有效范围内
    """
    if week not in range(1, 54):
        raise ValueError(f"周数必须在1-53之间，得到: {week}")
    
    # 使用ISO周日期
    # 获取该年第一个周一
    jan1 = datetime(year, 1, 1)
    days_to_monday = (7 - jan1.weekday()) % 7
    if days_to_monday == 0 and jan1.weekday() != 0:
        days_to_monday = 7
    
    first_monday = jan1 + timedelta(days=days_to_monday)
    
    # 如果1月1日是周二到周日，第一周从去年开始
    if jan1.weekday() > 0:
        week_start = first_monday + timedelta(weeks=week-2)
    else:
        week_start = first_monday + timedelta(weeks=week-1)
    
    start = week_start.replace(hour=0, minute=0, second=0, microsecond=0)
    end = (week_start + timedelta(days=6)).replace(hour=23, minute=59, second=59, microsecond=999999)
    
    return start, end


def parse_time_reference_to_range(time_ref: str, scale: TimeScale) -> Tuple[datetime, datetime]:
    """解析时间参考字符串为具体的时间范围
    
    Args:
        time_ref: 时间参考字符串
        scale: 时间尺度
        
    Returns:
        Tuple[datetime, datetime]: (开始时间, 结束时间)
        
    Raises:
        ValueError: 时间参考格式错误
    """
    if not validate_time_reference(time_ref, scale):
        raise ValueError(f"时间参考格式错误: {time_ref}，期望格式为 {scale.value}")
    
    if scale == TimeScale.YEAR:
        year = int(time_ref)
        start = datetime(year, 1, 1, 0, 0, 0, 0)
        end = datetime(year, 12, 31, 23, 59, 59, 999999)
        return start, end
    
    elif scale == TimeScale.QUARTER:
        # 格式: "2024-Q4"
        year, quarter_str = time_ref.split('-Q')
        year = int(year)
        quarter = int(quarter_str)
        return get_quarter_range(year, quarter)
    
    elif scale == TimeScale.MONTH:
        # 格式: "2025-07"
        year, month = map(int, time_ref.split('-'))
        return get_month_range(year, month)
    
    elif scale == TimeScale.WEEK:
        # 格式: "2025-W15"
        year, week_str = time_ref.split('-W')
        year = int(year)
        week = int(week_str)
        return get_week_range(year, week)
    elif scale == TimeScale.DAY:
        # 格式: "2025-07-15"
        date = datetime.strptime(time_ref, '%Y-%m-%d')
        start = date.replace(hour=0, minute=0, second=0, microsecond=0)
        end = date.replace(hour=23, minute=59, second=59, microsecond=999999)
        return start, end
    else:
        raise ValueError(f"不支持的时间尺度: {scale}")

