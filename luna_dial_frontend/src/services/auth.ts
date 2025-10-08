import apiClient, { storeSession, clearSession } from './api';
import {
  ApiResponse,
  LoginRequest,
  LoginResponse,
  User
} from '../types';

// 认证服务类
class AuthService {
  // 登录
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await apiClient.post<ApiResponse<LoginResponse>>(
      '/api/v1/public/auth/login',
      credentials
    );

    if (response.data.success && response.data.data) {
      // 存储session
      storeSession(response.data.data.session_id);
      return response.data.data;
    }

    throw new Error(response.data.message || '登录失败');
  }

  // 登出
  async logout(): Promise<void> {
    try {
      await apiClient.post<ApiResponse>('/api/v1/auth/logout');
    } finally {
      // 无论请求是否成功，都清除本地session
      clearSession();
    }
  }

  // 登出所有设备
  async logoutAll(): Promise<void> {
    try {
      await apiClient.delete<ApiResponse>('/api/v1/auth/logout-all');
    } finally {
      clearSession();
    }
  }

  // 获取用户资料
  async getProfile(): Promise<User> {
    const response = await apiClient.get<ApiResponse<User>>(
      '/api/v1/auth/profile'
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.message || '获取用户资料失败');
  }

  // 获取当前用户详细信息
  async getCurrentUser(): Promise<User & { session: any }> {
    const response = await apiClient.get<ApiResponse<User & { session: any }>>(
      '/api/v1/users/me'
    );

    if (response.data.success && response.data.data) {
      return response.data.data;
    }

    throw new Error(response.data.message || '获取用户信息失败');
  }
}

export default new AuthService();