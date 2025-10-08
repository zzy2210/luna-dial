import apiClient from './api';
import {
  ApiResponse,
  PlanResponse,
  GetByPeriodRequest
} from '../types';

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
}

export default new PlanService();