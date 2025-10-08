import React, { useState, useEffect } from 'react';
import { PeriodType } from '../../types';
import './styles.css';

interface TimeNavigatorProps {
  currentPeriod: PeriodType;
  currentDate: Date;
  onDateChange: (date: Date) => void;
  onNavigate: (direction: 'prev' | 'next') => void;
}

const TimeNavigator: React.FC<TimeNavigatorProps> = ({
  currentPeriod,
  currentDate,
  onDateChange,
  onNavigate
}) => {
  const [showPicker, setShowPicker] = useState(false);
  const [weekOptions, setWeekOptions] = useState<Array<{ value: string; label: string }>>([]);

  // 生成周选项
  const generateWeekOptions = (year: number) => {
    const options: Array<{ value: string; label: string }> = [];
    const yearStart = new Date(year, 0, 1);

    // 找到第一个周一（ISO周的开始）
    let firstMonday = new Date(yearStart);
    const startDay = yearStart.getDay();
    if (startDay === 0) { // 周日
      firstMonday.setDate(yearStart.getDate() + 1);
    } else if (startDay > 1) { // 周二到周六
      firstMonday.setDate(yearStart.getDate() + (8 - startDay));
    }

    const yearEnd = new Date(year, 11, 31);
    let currentMonday = new Date(firstMonday);
    let weekNum = 1;

    while (currentMonday.getFullYear() <= year && weekNum <= 53) {
      if (currentMonday > yearEnd) break;

      const weekEnd = new Date(currentMonday);
      weekEnd.setDate(currentMonday.getDate() + 6);

      const startMonth = currentMonday.getMonth() + 1;
      const startDay = currentMonday.getDate();
      const endMonth = weekEnd.getMonth() + 1;
      const endDay = weekEnd.getDate();

      let rangeText;
      if (startMonth === endMonth) {
        rangeText = `${startMonth}月${startDay}日-${endDay}日`;
      } else {
        rangeText = `${startMonth}月${startDay}日-${endMonth}月${endDay}日`;
      }

      options.push({
        value: `${year}-W${weekNum}`,
        label: `第${weekNum}周 (${rangeText})`
      });

      currentMonday.setDate(currentMonday.getDate() + 7);
      weekNum++;
    }

    return options;
  };

  // 根据年份和周数获取日期
  const getDateFromWeek = (year: number, week: number): Date => {
    const yearStart = new Date(year, 0, 1);

    let firstMonday = new Date(yearStart);
    const startDay = yearStart.getDay();
    if (startDay === 0) {
      firstMonday.setDate(yearStart.getDate() + 1);
    } else if (startDay > 1) {
      firstMonday.setDate(yearStart.getDate() + (8 - startDay));
    }

    const targetMonday = new Date(firstMonday);
    targetMonday.setDate(firstMonday.getDate() + (week - 1) * 7);

    return targetMonday;
  };

  // 更新周选项
  useEffect(() => {
    if (currentPeriod === 'week') {
      const year = currentDate.getFullYear();
      setWeekOptions(generateWeekOptions(year));
    }
  }, [currentPeriod, currentDate]);

  // 处理日期选择器变化
  const handleDateChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newDate = new Date(e.target.value);
    onDateChange(newDate);
    setShowPicker(false);
  };

  // 处理周选择器变化
  const handleWeekChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    if (e.target.value) {
      const [year, weekStr] = e.target.value.split('-W');
      const week = parseInt(weekStr);
      const newDate = getDateFromWeek(parseInt(year), week);
      onDateChange(newDate);
      setShowPicker(false);
    }
  };

  // 处理月份选择器变化
  const handleMonthChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const month = parseInt(e.target.value);
    if (month) {
      const year = currentDate.getFullYear();
      const newDate = new Date(year, month - 1, 1);
      onDateChange(newDate);
      setShowPicker(false);
    }
  };

  // 处理季度选择器变化
  const handleQuarterChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const quarter = parseInt(e.target.value);
    if (quarter) {
      const year = currentDate.getFullYear();
      const startMonth = (quarter - 1) * 3;
      const newDate = new Date(year, startMonth, 1);
      onDateChange(newDate);
      setShowPicker(false);
    }
  };

  // 处理年份选择器变化
  const handleYearChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const year = parseInt(e.target.value);
    if (year) {
      const newDate = new Date(year, 0, 1);
      onDateChange(newDate);
      // 如果是周视图，更新周选项
      if (currentPeriod === 'week') {
        setWeekOptions(generateWeekOptions(year));
      }
    }
  };

  // 格式化日期为 YYYY-MM-DD
  const formatDate = (date: Date): string => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  };

  // 获取当前周数（ISO周）
  const getCurrentWeek = (date: Date): number => {
    const year = date.getFullYear();
    const yearStart = new Date(year, 0, 1);

    let firstMonday = new Date(yearStart);
    const startDay = yearStart.getDay();
    if (startDay === 0) {
      firstMonday.setDate(yearStart.getDate() + 1);
    } else if (startDay > 1) {
      firstMonday.setDate(yearStart.getDate() + (8 - startDay));
    }

    const diff = date.getTime() - firstMonday.getTime();
    const weekNum = Math.floor(diff / (7 * 24 * 60 * 60 * 1000)) + 1;

    return Math.max(1, Math.min(weekNum, 53));
  };

  // 生成年份选项
  const generateYearOptions = () => {
    const currentYear = new Date().getFullYear();
    const years = [];
    for (let year = currentYear - 5; year <= currentYear + 5; year++) {
      years.push(year);
    }
    return years;
  };

  // 获取当前季度
  const getCurrentQuarter = (date: Date): number => {
    return Math.floor(date.getMonth() / 3) + 1;
  };

  const renderPicker = () => {
    switch (currentPeriod) {
      case 'day':
        return (
          <input
            type="date"
            className="date-picker"
            value={formatDate(currentDate)}
            onChange={handleDateChange}
          />
        );

      case 'week':
        return (
          <>
            <select
              className="year-picker"
              value={currentDate.getFullYear()}
              onChange={handleYearChange}
            >
              {generateYearOptions().map(year => (
                <option key={year} value={year}>{year}年</option>
              ))}
            </select>
            <select
              className="week-picker"
              value={`${currentDate.getFullYear()}-W${getCurrentWeek(currentDate)}`}
              onChange={handleWeekChange}
            >
              <option value="">选择周</option>
              {weekOptions.map(option => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </>
        );

      case 'month':
        return (
          <>
            <select
              className="year-picker"
              value={currentDate.getFullYear()}
              onChange={handleYearChange}
            >
              {generateYearOptions().map(year => (
                <option key={year} value={year}>{year}年</option>
              ))}
            </select>
            <select
              className="month-picker"
              value={currentDate.getMonth() + 1}
              onChange={handleMonthChange}
            >
              {['一月', '二月', '三月', '四月', '五月', '六月',
                '七月', '八月', '九月', '十月', '十一月', '十二月'].map((month, index) => (
                <option key={index + 1} value={index + 1}>
                  {month} ({index + 1}月)
                </option>
              ))}
            </select>
          </>
        );

      case 'quarter':
        return (
          <>
            <select
              className="year-picker"
              value={currentDate.getFullYear()}
              onChange={handleYearChange}
            >
              {generateYearOptions().map(year => (
                <option key={year} value={year}>{year}年</option>
              ))}
            </select>
            <select
              className="quarter-picker"
              value={getCurrentQuarter(currentDate)}
              onChange={handleQuarterChange}
            >
              <option value="1">Q1 第一季度 (1-3月)</option>
              <option value="2">Q2 第二季度 (4-6月)</option>
              <option value="3">Q3 第三季度 (7-9月)</option>
              <option value="4">Q4 第四季度 (10-12月)</option>
            </select>
          </>
        );

      case 'year':
        return (
          <select
            className="year-picker"
            value={currentDate.getFullYear()}
            onChange={handleYearChange}
          >
            {generateYearOptions().map(year => (
              <option key={year} value={year}>{year}年</option>
            ))}
          </select>
        );

      default:
        return null;
    }
  };

  return (
    <div className="time-navigator">
      {/* 前进/后退按钮 */}
      <button
        className="nav-btn nav-prev"
        onClick={() => onNavigate('prev')}
        title="上一个周期"
      >
        ◀
      </button>

      <button
        className="nav-btn nav-next"
        onClick={() => onNavigate('next')}
        title="下一个周期"
      >
        ▶
      </button>

      {/* 选择器容器 */}
      {showPicker && (
        <div className="picker-container">
          {renderPicker()}
        </div>
      )}

      {/* 跳转按钮 */}
      <button
        className="btn-time-jump"
        onClick={() => setShowPicker(!showPicker)}
      >
        跳转到...
      </button>
    </div>
  );
};

export default TimeNavigator;
