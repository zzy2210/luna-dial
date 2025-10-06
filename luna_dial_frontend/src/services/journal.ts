import apiClient from './api';
import {
  ApiResponse,
  Journal,
  CreateJournalRequest,
  UpdateJournalRequest,
  GetByPeriodRequest,
  PaginatedResponse,
  JournalType
} from '../types';

// 日志服务类
class JournalService {
  // 获取日志列表（按时间周期）
  async getJournals(params: GetByPeriodRequest): Promise<Journal[]> {
    const response = await apiClient.get<ApiResponse<Journal[]>>(
      '/api/v1/journals',
      { data: params }
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.message || '获取日志列表失败');
  }

  // 创建日志
  async createJournal(journal: CreateJournalRequest): Promise<Journal> {
    const response = await apiClient.post<ApiResponse<Journal>>(
      '/api/v1/journals',
      journal
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.message || '创建日志失败');
  }

  // 更新日志
  async updateJournal(journalId: string, updates: UpdateJournalRequest): Promise<Journal> {
    const response = await apiClient.put<ApiResponse<Journal>>(
      `/api/v1/journals/${journalId}`,
      {
        ...updates,
        journal_id: journalId
      }
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.message || '更新日志失败');
  }

  // 删除日志
  async deleteJournal(journalId: string): Promise<void> {
    const response = await apiClient.delete<ApiResponse>(
      `/api/v1/journals/${journalId}`
    );

    if (!response.data.success) {
      throw new Error(response.data.message || '删除日志失败');
    }
  }

  // 分页查询日志
  async getJournalsPaginated(
    page: number = 1,
    pageSize: number = 20,
    filters?: {
      journalType?: string;
      startDate?: string;
      endDate?: string;
    }
  ): Promise<PaginatedResponse<Journal>> {
    const params = new URLSearchParams();
    params.append('page', page.toString());
    params.append('page_size', pageSize.toString());

    if (filters?.journalType) {
      params.append('journal_type', filters.journalType);
    }
    if (filters?.startDate) {
      params.append('start_date', filters.startDate);
    }
    if (filters?.endDate) {
      params.append('end_date', filters.endDate);
    }

    const response = await apiClient.get<ApiResponse<PaginatedResponse<Journal>>>(
      `/api/v1/journals/paginated?${params.toString()}`
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.message || '获取日志分页失败');
  }

  // 工具方法：将日志类型字符串转换为枚举值
  journalTypeFromString(type: string): JournalType {
    const map: Record<string, JournalType> = {
      'day': JournalType.Day,
      'week': JournalType.Week,
      'month': JournalType.Month,
      'quarter': JournalType.Quarter,
      'year': JournalType.Year,
    };
    return map[type] ?? JournalType.Day;
  }

  // 工具方法：将日志类型枚举转换为字符串
  journalTypeToString(type: JournalType): string {
    const map: Record<JournalType, string> = {
      [JournalType.Day]: 'day',
      [JournalType.Week]: 'week',
      [JournalType.Month]: 'month',
      [JournalType.Quarter]: 'quarter',
      [JournalType.Year]: 'year',
    };
    return map[type];
  }

  // 工具方法：获取日志类型的中文名称
  getJournalTypeLabel(type: JournalType): string {
    const labels: Record<JournalType, string> = {
      [JournalType.Day]: '日志',
      [JournalType.Week]: '周志',
      [JournalType.Month]: '月志',
      [JournalType.Quarter]: '季志',
      [JournalType.Year]: '年志',
    };
    return labels[type];
  }
}

export default new JournalService();