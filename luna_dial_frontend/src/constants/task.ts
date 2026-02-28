import { TaskPriority, TaskStatus } from '../types';
import type { TaskPriority as TaskPriorityValue, TaskStatus as TaskStatusValue, UpdateTaskRequest } from '../types';

// TaskStatus 对应的 API status 字符串类型
type TaskStatusApiString = NonNullable<UpdateTaskRequest['status']>;

// 任务状态到 API status 字符串映射
export const TASK_STATUS_TO_API: Record<TaskStatusValue, TaskStatusApiString> = {
  [TaskStatus.NotStarted]: 'not_started',
  [TaskStatus.InProgress]: 'in_progress',
  [TaskStatus.Completed]: 'completed',
  [TaskStatus.Cancelled]: 'cancelled',
} as const;

// 任务状态到中文文案映射
export const TASK_STATUS_LABELS: Record<TaskStatusValue, string> = {
  [TaskStatus.NotStarted]: '未开始',
  [TaskStatus.InProgress]: '进行中',
  [TaskStatus.Completed]: '已完成',
  [TaskStatus.Cancelled]: '已取消',
} as const;

// 任务优先级到中文文案映射
export const TASK_PRIORITY_LABELS: Record<TaskPriorityValue, string> = {
  [TaskPriority.Low]: '低',
  [TaskPriority.Medium]: '中',
  [TaskPriority.High]: '高',
  [TaskPriority.Urgent]: '紧急',
} as const;

