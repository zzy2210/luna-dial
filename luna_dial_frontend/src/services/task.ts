import apiClient from './api';
import {
  ApiResponse,
  Task,
  CreateTaskRequest,
  UpdateTaskRequest,
  PaginatedResponse,
  GetByPeriodRequest,
  TaskType,
  TaskStatus,
  TaskPriority
} from '../types';

// 任务服务类
class TaskService {
  // 获取任务列表（按时间周期）
  async getTasks(params: GetByPeriodRequest): Promise<Task[]> {
    const response = await apiClient.get<ApiResponse<Task[]>>(
      '/api/v1/tasks',
      { data: params }
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.message || '获取任务列表失败');
  }

  // 创建任务
  async createTask(task: CreateTaskRequest): Promise<Task> {
    const response = await apiClient.post<ApiResponse<Task>>(
      '/api/v1/tasks',
      task
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.message || '创建任务失败');
  }

  // 更新任务
  async updateTask(taskId: string, updates: UpdateTaskRequest): Promise<Task> {
    const response = await apiClient.put<ApiResponse<Task>>(
      `/api/v1/tasks/${taskId}`,
      updates
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.message || '更新任务失败');
  }

  // 删除任务
  async deleteTask(taskId: string): Promise<void> {
    const response = await apiClient.delete<ApiResponse>(
      `/api/v1/tasks/${taskId}`,
      { data: { task_id: taskId } }
    );

    if (!response.data.success) {
      throw new Error(response.data.message || '删除任务失败');
    }
  }

  // 完成任务
  async completeTask(taskId: string): Promise<void> {
    const response = await apiClient.post<ApiResponse>(
      `/api/v1/tasks/${taskId}/complete`,
      { task_id: taskId }
    );

    if (!response.data.success) {
      throw new Error(response.data.message || '完成任务失败');
    }
  }

  // 创建子任务
  async createSubtask(parentId: string, subtask: CreateTaskRequest): Promise<Task> {
    const response = await apiClient.post<ApiResponse<Task>>(
      `/api/v1/tasks/${parentId}/subtasks`,
      subtask
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.message || '创建子任务失败');
  }

  // 更新任务评分
  async updateScore(taskId: string, score: number): Promise<void> {
    const response = await apiClient.put<ApiResponse>(
      `/api/v1/tasks/${taskId}/score`,
      { score }
    );

    if (!response.data.success) {
      throw new Error(response.data.message || '更新任务评分失败');
    }
  }

  // 获取根任务分页
  async getRootTasks(
    page: number = 1,
    pageSize: number = 20,
    filters?: {
      status?: string[];
      priority?: string[];
      taskType?: string[];
    }
  ): Promise<PaginatedResponse<Task>> {
    const params = new URLSearchParams();
    params.append('page', page.toString());
    params.append('page_size', pageSize.toString());

    if (filters?.status) {
      filters.status.forEach(s => params.append('status', s));
    }
    if (filters?.priority) {
      filters.priority.forEach(p => params.append('priority', p));
    }
    if (filters?.taskType) {
      filters.taskType.forEach(t => params.append('task_type', t));
    }

    const response = await apiClient.get<ApiResponse<PaginatedResponse<Task>>>(
      `/api/v1/tasks/roots?${params.toString()}`
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.message || '获取根任务失败');
  }

  // 获取任务树
  async getTaskTree(
    page: number = 1,
    pageSize: number = 10,
    includeEmpty: boolean = true
  ): Promise<PaginatedResponse<Task>> {
    const params = new URLSearchParams();
    params.append('page', page.toString());
    params.append('page_size', pageSize.toString());
    params.append('include_empty', includeEmpty.toString());

    const response = await apiClient.get<ApiResponse<PaginatedResponse<Task>>>(
      `/api/v1/tasks/tree?${params.toString()}`
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.message || '获取任务树失败');
  }

  // 获取指定任务的完整任务树
  async getTaskSubtree(taskId: string): Promise<Task> {
    const response = await apiClient.get<ApiResponse<Task>>(
      `/api/v1/tasks/${taskId}/tree`
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.message || '获取任务子树失败');
  }

  // 获取任务的父任务链
  async getTaskParents(taskId: string): Promise<Task[]> {
    const response = await apiClient.get<ApiResponse<Task[]>>(
      `/api/v1/tasks/${taskId}/parents`
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.message || '获取父任务链失败');
  }

  // 移动任务
  async moveTask(taskId: string, newParentId?: string): Promise<void> {
    const response = await apiClient.put<ApiResponse>(
      `/api/v1/tasks/${taskId}/move`,
      {
        task_id: taskId,
        new_parent_id: newParentId
      }
    );

    if (!response.data.success) {
      throw new Error(response.data.message || '移动任务失败');
    }
  }

  // 工具方法：将任务类型字符串转换为枚举值
  taskTypeFromString(type: string): TaskType {
    const map: Record<string, TaskType> = {
      'day': TaskType.Day,
      'week': TaskType.Week,
      'month': TaskType.Month,
      'quarter': TaskType.Quarter,
      'year': TaskType.Year,
    };
    return map[type] ?? TaskType.Day;
  }

  // 工具方法：将任务状态字符串转换为枚举值
  taskStatusFromString(status: string): TaskStatus {
    const map: Record<string, TaskStatus> = {
      'not_started': TaskStatus.NotStarted,
      'in_progress': TaskStatus.InProgress,
      'completed': TaskStatus.Completed,
      'cancelled': TaskStatus.Cancelled,
    };
    return map[status] ?? TaskStatus.NotStarted;
  }

  // 工具方法：将任务优先级字符串转换为枚举值
  taskPriorityFromString(priority: string): TaskPriority {
    const map: Record<string, TaskPriority> = {
      'low': TaskPriority.Low,
      'medium': TaskPriority.Medium,
      'high': TaskPriority.High,
      'urgent': TaskPriority.Urgent,
    };
    return map[priority] ?? TaskPriority.Medium;
  }
}

export default new TaskService();