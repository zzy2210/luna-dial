import React, { useEffect } from 'react';
import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { Loader2 } from 'lucide-react';
import useAuthStore from '../store/auth';
import { getStoredSession } from '../services/api';

interface ProtectedRouteProps {
  children?: React.ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
  const location = useLocation();
  const { isAuthenticated, isLoading, checkAuth } = useAuthStore();
  const session = getStoredSession();

  useEffect(() => {
    // 如果有session但还没有验证，尝试验证
    if (session && !isAuthenticated && !isLoading) {
      checkAuth();
    }
  }, [session, isAuthenticated, isLoading, checkAuth]);

  // 加载中状态
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <Loader2 className="w-8 h-8 text-indigo-600 animate-spin mx-auto mb-4" />
          <p className="text-gray-600">加载中...</p>
        </div>
      </div>
    );
  }

  // 未认证，重定向到登录页
  if (!session || !isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  // 已认证，渲染子组件
  return children ? <>{children}</> : <Outlet />;
};

export default ProtectedRoute;