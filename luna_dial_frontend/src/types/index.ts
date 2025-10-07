// 用户相关类型
export interface User {
  id: string;
  username: string;
  name: string;
  email: string;
  created_at: string;
  updated_at: string;
}

// Session相关类型
export interface Session {
  session_id: string;
  expires_in: number;
  last_access_at?: string;
  expires_at?: string;
}

// 登录请求类型
export interface LoginRequest {
  username: string;
  password: string;
}

// 登录响应类型
export interface LoginResponse {
  session_id: string;
  expires_in: number;
}

// API响应通用类型
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  success: boolean;
  timestamp: number;
  data?: T;
}

// 任务状态枚举
export const TaskStatus = {
  NotStarted: 0,
  InProgress: 1,
  Completed: 2,
  Cancelled: 3,
} as const;
export type TaskStatus = typeof TaskStatus[keyof typeof TaskStatus];

// 任务优先级枚举
export const TaskPriority = {
  Low: 0,
  Medium: 1,
  High: 2,
  Urgent: 3,
} as const;
export type TaskPriority = typeof TaskPriority[keyof typeof TaskPriority];

// 任务类型枚举
export const TaskType = {
  Day: 0,
  Week: 1,
  Month: 2,
  Quarter: 3,
  Year: 4,
} as const;
export type TaskType = typeof TaskType[keyof typeof TaskType];

// 任务类型
export interface Task {
  id: string;
  user_id: string;
  title: string;
  task_type: TaskType;
  period_start: string;
  period_end: string;
  tags?: string;
  icon?: string;
  score: number;
  status: TaskStatus;
  priority: TaskPriority;
  parent_id?: string;
  has_children: boolean;
  children_count: number;
  root_task_id: string;
  tree_depth: number;
  created_at: string;
  updated_at: string;
  children?: Task[];
}

// 任务创建请求
export interface CreateTaskRequest {
  title: string;
  description?: string;
  start_date: string;
  end_date: string;
  priority: 'low' | 'medium' | 'high' | 'urgent';
  icon?: string;
  tags?: string[];
  parent_id?: string;
}

// 任务更新请求
export interface UpdateTaskRequest {
  title?: string;
  description?: string;
  start_date?: string;
  end_date?: string;
  priority?: 'low' | 'medium' | 'high' | 'urgent';
  status?: 'not_started' | 'in_progress' | 'completed' | 'cancelled';
  icon?: string;
  tags?: string[];
}

// 日志类型枚举
export const JournalType = {
  Day: 0,
  Week: 1,
  Month: 2,
  Quarter: 3,
  Year: 4,
} as const;
export type JournalType = typeof JournalType[keyof typeof JournalType];

// 日志类型
export interface Journal {
  id: string;
  user_id: string;
  title: string;
  content: string;
  journal_type: JournalType;
  time_period: {
    start: string;
    end: string;
  };
  icon?: string;
  created_at: string;
  updated_at: string;
}

// 日志创建请求
export interface CreateJournalRequest {
  title: string;
  content: string;
  journal_type: 'day' | 'week' | 'month' | 'quarter' | 'year';
  start_date: string;
  end_date: string;
  icon?: string;
}

// 日志更新请求
export interface UpdateJournalRequest {
  journal_id: string;
  title?: string;
  content?: string;
  journal_type?: 'day' | 'week' | 'month' | 'quarter' | 'year';
  icon?: string;
}

// 周期类型
export type PeriodType = 'day' | 'week' | 'month' | 'quarter' | 'year';

// 获取任务/日志/计划请求
export interface GetByPeriodRequest {
  period_type: PeriodType;
  start_date: string;
  end_date: string;
}

// 分页信息
export interface Pagination {
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
  has_next: boolean;
  has_prev: boolean;
}

// 分页响应
export interface PaginatedResponse<T> {
  items: T[];
  pagination: Pagination;
}

// 计划响应
export interface PlanResponse {
  tasks: Task[];
  tasks_total: number;
  journals: Journal[];
  journals_total: number;
  plan_type: PeriodType;
  plan_period: {
    start: string;
    end: string;
  };
  score_total: number;
  group_stats: GroupStats[];
}

// 分组统计
export interface GroupStats {
  group_key: string;
  task_count: number;
  score_total: number;
}