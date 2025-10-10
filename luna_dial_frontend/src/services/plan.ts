import apiClient from './api';
import {
  ApiResponse,
  PlanResponse,
  GetByPeriodRequest
} from '../types';

// 分组统计数据接口
export interface GroupStat {
  group_key: string;
  task_count: number;
  score_total: number;
}

// 计划服务类
class PlanService {
  // 获取计划信息（包含任务、日志和统计）
  async getPlan(params: GetByPeriodRequest): Promise<PlanResponse> {
    const response = await apiClient.get<ApiResponse<PlanResponse>>(
      '/api/v1/plans',
      { params }
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.message || '获取计划信息失败');
  }

  // 获取计划统计数据（按时间范围分组）
  async getPlanStats(params: {
    group_by: string;
    start_date: string;
    end_date: string;
  }): Promise<GroupStat[]> {
    const response = await apiClient.get<ApiResponse<GroupStat[]>>(
      '/api/v1/plans/stats',
      { params }
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.message || '获取统计数据失败');
  }
}

export default new PlanService();