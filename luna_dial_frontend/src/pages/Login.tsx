import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import useAuthStore from '../store/auth';
import '../styles/auth.css';

const LoginPage: React.FC = () => {
  const navigate = useNavigate();
  const { login, isLoading, error, clearError } = useAuthStore();

  const [formData, setFormData] = useState({
    username: '',
    password: '',
    rememberMe: false
  });

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));

    // 清除错误信息
    if (error) {
      clearError();
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      await login(formData.username, formData.password);
      navigate('/dashboard');
    } catch (error) {
      // 错误已经在store中处理
      console.error('Login failed:', error);
    }
  };

  return (
    <div className="auth-container">
      {/* Logo和标题 */}
      <div className="auth-header">
        <div className="logo-container">
          <span className="logo-icon">🌙</span>
          <h1 className="logo-text">Luna Dial</h1>
        </div>
        <p className="tagline">积极向上，记录成长</p>
      </div>

      {/* 登录卡片 */}
      <div className="auth-card">
        <h2 className="auth-title">欢迎回来</h2>

        {/* 错误提示 */}
        {error && (
          <div className="error-message">
            <span className="error-icon">⚠️</span>
            <span className="error-text">{error}</span>
          </div>
        )}

        {/* 登录表单 */}
        <form onSubmit={handleSubmit} className="auth-form">
          <div className="form-group">
            <label htmlFor="username" className="form-label">
              用户名
            </label>
            <input
              type="text"
              id="username"
              name="username"
              value={formData.username}
              onChange={handleInputChange}
              className="form-input"
              placeholder="请输入用户名"
              required
              autoComplete="username"
              disabled={isLoading}
            />
          </div>

          <div className="form-group">
            <label htmlFor="password" className="form-label">
              密码
            </label>
            <input
              type="password"
              id="password"
              name="password"
              value={formData.password}
              onChange={handleInputChange}
              className="form-input"
              placeholder="请输入密码"
              required
              autoComplete="current-password"
              disabled={isLoading}
            />
          </div>

          <div className="form-options">
            <label className="checkbox-label">
              <input
                type="checkbox"
                id="rememberMe"
                name="rememberMe"
                checked={formData.rememberMe}
                onChange={handleInputChange}
                className="checkbox-input"
                disabled={isLoading}
              />
              <span className="checkbox-text">记住我</span>
            </label>
          </div>

          <button type="submit" className="btn-primary" disabled={isLoading}>
            {isLoading ? (
              <span className="btn-loading">
                <span className="spinner"></span>
                登录中...
              </span>
            ) : (
              <span className="btn-text">登录</span>
            )}
          </button>
        </form>

        {/* 其他选项 */}
        <div className="auth-footer">
          <p className="footer-text">
            首次使用？请联系管理员创建账户
          </p>
        </div>
      </div>

      {/* 底部信息 */}
      <div className="auth-bottom">
        <p className="copyright">© 2025 Luna Dial. All rights reserved.</p>
      </div>
    </div>
  );
};

export default LoginPage;