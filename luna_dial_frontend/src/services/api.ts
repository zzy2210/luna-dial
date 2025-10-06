import axios, { AxiosInstance, AxiosError, InternalAxiosRequestConfig } from 'axios';
import { ApiResponse } from '../types';

// API基础URL
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081';

// 创建axios实例
const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 获取存储的session
export const getStoredSession = (): string | null => {
  return localStorage.getItem('session_id');
};

// 存储session
export const storeSession = (sessionId: string): void => {
  localStorage.setItem('session_id', sessionId);
};

// 清除session
export const clearSession = (): void => {
  localStorage.removeItem('session_id');
};

// 请求拦截器
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const sessionId = getStoredSession();
    if (sessionId) {
      config.headers.Authorization = `Bearer ${sessionId}`;
    }
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error: AxiosError<ApiResponse>) => {
    if (error.response?.status === 401) {
      // 未授权，清除session并跳转到登录页
      clearSession();
      window.location.href = '/login';
    }

    // 提取错误信息
    const message = error.response?.data?.message || error.message || '网络错误';
    console.error('API Error:', message);

    return Promise.reject(error);
  }
);

export default apiClient;