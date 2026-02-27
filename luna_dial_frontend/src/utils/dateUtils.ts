import { PeriodType } from '../types';

/**
 * 本地时间格式化函数，避免 toISOString() 的时区问题
 */
export const formatLocalDate = (date: Date): string => {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
};

/**
 * 获取ISO周数
 */
export const getWeekNumber = (date: Date): number => {
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

/**
 * 获取季度数
 */
export const getQuarterNumber = (date: Date): number => {
  return Math.floor(date.getMonth() / 3) + 1;
};

/**
 * 判断两个日期是否在同一周期
 */
export const isSamePeriod = (date1: Date, date2: Date, period: PeriodType): boolean => {
  const y1 = date1.getFullYear();
  const y2 = date2.getFullYear();
  const m1 = date1.getMonth();
  const m2 = date2.getMonth();
  const d1 = date1.getDate();
  const d2 = date2.getDate();

  switch (period) {
    case 'day':
      return y1 === y2 && m1 === m2 && d1 === d2;
    case 'week':
      return y1 === y2 && getWeekNumber(date1) === getWeekNumber(date2);
    case 'month':
      return y1 === y2 && m1 === m2;
    case 'quarter':
      return y1 === y2 && getQuarterNumber(date1) === getQuarterNumber(date2);
    case 'year':
      return y1 === y2;
    default:
      return false;
  }
};

/**
 * 获取周期的起止日期
 */
export const getPeriodDates = (period: PeriodType, baseDate: Date = new Date()) => {
  const startDate = new Date(baseDate);
  const endDate = new Date(baseDate);

  switch (period) {
    case 'day':
      // 指定日期 [date 00:00, date+1 00:00)
      startDate.setHours(0, 0, 0, 0);
      endDate.setDate(endDate.getDate() + 1);
      endDate.setHours(0, 0, 0, 0);
      break;
    case 'week':
      // 指定周 ISO Week [Monday 00:00, Next Monday 00:00)
      const day = startDate.getDay();
      const diff = startDate.getDate() - day + (day === 0 ? -6 : 1);
      startDate.setDate(diff);
      startDate.setHours(0, 0, 0, 0);
      endDate.setDate(startDate.getDate() + 7);
      endDate.setHours(0, 0, 0, 0);
      break;
    case 'month':
      // 指定月 [1st 00:00, Next Month 1st 00:00)
      startDate.setDate(1);
      startDate.setHours(0, 0, 0, 0);
      endDate.setMonth(endDate.getMonth() + 1, 1);
      endDate.setHours(0, 0, 0, 0);
      break;
    case 'quarter':
      // 指定季度 [Quarter Start 00:00, Next Quarter Start 00:00)
      const quarter = Math.floor(startDate.getMonth() / 3);
      startDate.setMonth(quarter * 3, 1);
      startDate.setHours(0, 0, 0, 0);
      endDate.setMonth((quarter + 1) * 3, 1);
      endDate.setHours(0, 0, 0, 0);
      break;
    case 'year':
      // 指定年 [Jan 1 00:00, Next Year Jan 1 00:00)
      startDate.setMonth(0, 1);
      startDate.setHours(0, 0, 0, 0);
      endDate.setFullYear(endDate.getFullYear() + 1, 0, 1);
      endDate.setHours(0, 0, 0, 0);
      break;
  }

  return {
    start_date: formatLocalDate(startDate),
    end_date: formatLocalDate(endDate)
  };
};

/**
 * 格式化周期范围显示
 */
export const formatPeriodRange = (date: Date, period: PeriodType): string => {
  const year = date.getFullYear();
  const month = date.getMonth() + 1;

  switch (period) {
    case 'day':
      return '';  // 日期在外面单独显示
    case 'week': {
      const weekNum = getWeekNumber(date);
      // 计算周的起止日期
      const yearStart = new Date(year, 0, 1);
      let firstMonday = new Date(yearStart);
      const firstDayOfYear = yearStart.getDay();
      if (firstDayOfYear === 0) {
        firstMonday.setDate(yearStart.getDate() + 1);
      } else if (firstDayOfYear > 1) {
        firstMonday.setDate(yearStart.getDate() + (8 - firstDayOfYear));
      }
      const weekStart = new Date(firstMonday);
      weekStart.setDate(firstMonday.getDate() + (weekNum - 1) * 7);
      const weekEnd = new Date(weekStart);
      weekEnd.setDate(weekStart.getDate() + 6);

      const startMonth = weekStart.getMonth() + 1;
      const weekStartDay = weekStart.getDate();
      const endMonth = weekEnd.getMonth() + 1;
      const weekEndDay = weekEnd.getDate();

      if (startMonth === endMonth) {
        return `第${weekNum}周: ${startMonth}月${weekStartDay}日-${weekEndDay}日`;
      } else {
        return `第${weekNum}周: ${startMonth}月${weekStartDay}日-${endMonth}月${weekEndDay}日`;
      }
    }
    case 'month':
      return `${month}月`;
    case 'quarter': {
      const quarter = getQuarterNumber(date);
      const startMonth = (quarter - 1) * 3 + 1;
      const endMonth = quarter * 3;
      return `第${quarter}季度 (${startMonth}-${endMonth}月)`;
    }
    case 'year':
      return '';  // 年份在外面单独显示
    default:
      return '';
  }
};

/**
 * 获取当前日期字符串
 */
export const getCurrentDateString = (date: Date = new Date()) => {
  const year = date.getFullYear();
  const month = date.getMonth() + 1;
  const day = date.getDate();
  const weekDay = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六'][date.getDay()];
  return `${year}年${month}月${day}日 ${weekDay}`;
};

/**
 * 获取周期标签
 */
export const getPeriodLabel = (period: PeriodType, date: Date = new Date()) => {
  const now = new Date();
  const isCurrentPeriod = isSamePeriod(date, now, period);
  const year = date.getFullYear();

  if (isCurrentPeriod) {
    const labels = {
      day: '今日',
      week: '本周',
      month: '本月',
      quarter: '本季度',
      year: '本年'
    };
    return labels[period];
  } else {
    // 非当前周期，显示具体时间
    switch (period) {
      case 'day':
        return getCurrentDateString(date).replace(/\s星期./, '');  // 移除星期部分
      case 'week':
        return `${year}年${formatPeriodRange(date, period)}`;
      case 'month':
        return `${year}年${formatPeriodRange(date, period)}`;
      case 'quarter':
        return `${year}年${formatPeriodRange(date, period)}`;
      case 'year':
        return `${year}年`;
      default:
        return '';
    }
  }
};

/**
 * 获取上级周期配置
 */
export interface ParentPeriodConfig {
  parent: PeriodType;
  taskType: number;
}

export const getParentPeriodConfig = (period: PeriodType): ParentPeriodConfig | null => {
  const periodMap: Record<PeriodType, ParentPeriodConfig | null> = {
    'day': { parent: 'week', taskType: 1 },
    'week': { parent: 'month', taskType: 2 },
    'month': { parent: 'quarter', taskType: 3 },
    'quarter': { parent: 'year', taskType: 4 },
    'year': null  // 不显示
  };

  return periodMap[period];
};
