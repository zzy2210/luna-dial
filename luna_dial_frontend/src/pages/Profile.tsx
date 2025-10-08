import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import useAuthStore from '../store/auth';
import '../styles/profile.css';

const Profile: React.FC = () => {
  const navigate = useNavigate();
  const { user, logout } = useAuthStore();
  const [activeTab, setActiveTab] = useState('profile');

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  const handleGoBack = () => {
    navigate('/dashboard');
  };

  return (
    <div className="profile-page">
      {/* 顶部导航 */}
      <header className="profile-header">
        <button className="btn-back" onClick={handleGoBack}>
          ← 返回
        </button>
        <h1>个人设置</h1>
        <div className="header-spacer"></div>
      </header>

      <div className="profile-container">
        {/* 侧边栏 */}
        <aside className="profile-sidebar">
          <div className="user-info-card">
            <div className="user-avatar">
              {user?.name?.charAt(0) || user?.username?.charAt(0) || 'U'}
            </div>
            <h2>{user?.name || user?.username}</h2>
            <p>{user?.email || `@${user?.username}`}</p>
          </div>

          <nav className="settings-nav">
            <button
              className={`nav-item ${activeTab === 'profile' ? 'active' : ''}`}
              onClick={() => setActiveTab('profile')}
            >
              <span className="nav-icon">👤</span>
              个人信息
            </button>
            <button
              className={`nav-item ${activeTab === 'preferences' ? 'active' : ''}`}
              onClick={() => setActiveTab('preferences')}
            >
              <span className="nav-icon">⚙️</span>
              偏好设置
            </button>
            <button
              className={`nav-item ${activeTab === 'security' ? 'active' : ''}`}
              onClick={() => setActiveTab('security')}
            >
              <span className="nav-icon">🔒</span>
              账号安全
            </button>
            <button
              className={`nav-item ${activeTab === 'about' ? 'active' : ''}`}
              onClick={() => setActiveTab('about')}
            >
              <span className="nav-icon">ℹ️</span>
              关于
            </button>
          </nav>

          <button className="btn-logout" onClick={handleLogout}>
            登出账号
          </button>
        </aside>

        {/* 主内容区 */}
        <main className="profile-content">
          {activeTab === 'profile' && (
            <section className="settings-section">
              <h2>个人信息</h2>
              <div className="settings-card">
                <div className="form-group">
                  <label>用户名</label>
                  <input type="text" value={user?.username || ''} disabled />
                  <span className="form-hint">用户名创建后不可更改</span>
                </div>

                <div className="form-group">
                  <label>显示名称</label>
                  <input
                    type="text"
                    placeholder="输入您的显示名称"
                    defaultValue={user?.name || ''}
                  />
                </div>

                <div className="form-group">
                  <label>邮箱地址</label>
                  <input
                    type="email"
                    placeholder="your@email.com"
                    defaultValue={user?.email || ''}
                  />
                </div>

                <div className="form-group">
                  <label>个人简介</label>
                  <textarea
                    placeholder="介绍一下自己..."
                    rows={3}
                  />
                </div>

                <div className="form-actions">
                  <button className="btn-primary">保存更改</button>
                </div>
              </div>
            </section>
          )}

          {activeTab === 'preferences' && (
            <section className="settings-section">
              <h2>偏好设置</h2>
              <div className="settings-card">
                <div className="setting-item">
                  <div className="setting-info">
                    <h3>主题模式</h3>
                    <p>选择您喜欢的界面主题</p>
                  </div>
                  <select className="setting-control">
                    <option value="auto">跟随系统</option>
                    <option value="light">浅色模式</option>
                    <option value="dark">深色模式</option>
                  </select>
                </div>

                <div className="setting-item">
                  <div className="setting-info">
                    <h3>语言</h3>
                    <p>选择界面显示语言</p>
                  </div>
                  <select className="setting-control">
                    <option value="zh-CN">简体中文</option>
                    <option value="zh-TW">繁體中文</option>
                    <option value="en">English</option>
                  </select>
                </div>

                <div className="setting-item">
                  <div className="setting-info">
                    <h3>默认任务视图</h3>
                    <p>登录后默认显示的任务周期</p>
                  </div>
                  <select className="setting-control">
                    <option value="day">日视图</option>
                    <option value="week">周视图</option>
                    <option value="month">月视图</option>
                  </select>
                </div>

                <div className="setting-item">
                  <div className="setting-info">
                    <h3>提醒通知</h3>
                    <p>接收任务提醒和系统通知</p>
                  </div>
                  <label className="switch">
                    <input type="checkbox" defaultChecked />
                    <span className="slider"></span>
                  </label>
                </div>

                <div className="form-actions">
                  <button className="btn-primary">保存设置</button>
                </div>
              </div>
            </section>
          )}

          {activeTab === 'security' && (
            <section className="settings-section">
              <h2>账号安全</h2>
              <div className="settings-card">
                <div className="security-item">
                  <h3>修改密码</h3>
                  <p>定期更新密码有助于保护您的账号安全</p>
                  <button className="btn-secondary">更改密码</button>
                </div>

                <div className="security-item">
                  <h3>登录历史</h3>
                  <p>查看您的账号最近登录记录</p>
                  <button className="btn-secondary">查看历史</button>
                </div>

                <div className="security-item danger">
                  <h3>删除账号</h3>
                  <p>永久删除您的账号和所有相关数据</p>
                  <button className="btn-danger">删除账号</button>
                </div>
              </div>
            </section>
          )}

          {activeTab === 'about' && (
            <section className="settings-section">
              <h2>关于</h2>
              <div className="settings-card">
                <div className="about-content">
                  <div className="app-info">
                    <h3>🌙 Luna Dial</h3>
                    <p>个人任务与日志管理系统</p>
                    <p className="version">版本 1.0.0</p>
                  </div>

                  <div className="about-section">
                    <h4>系统特色</h4>
                    <ul>
                      <li>📅 分层任务管理（年→季→月→周→日）</li>
                      <li>📝 多维度日志记录</li>
                      <li>📊 成长数据可视化</li>
                      <li>🎯 积极导向的设计理念</li>
                    </ul>
                  </div>

                  <div className="about-section">
                    <h4>技术栈</h4>
                    <ul>
                      <li>前端：React + TypeScript + Vite</li>
                      <li>后端：Go + Gin + GORM</li>
                      <li>数据库：PostgreSQL</li>
                      <li>架构：DDD领域驱动设计</li>
                    </ul>
                  </div>

                  <div className="about-section">
                    <h4>开发团队</h4>
                    <p>Luna Dial 团队</p>
                    <p className="copyright">© 2024 Luna Dial. All rights reserved.</p>
                  </div>
                </div>
              </div>
            </section>
          )}
        </main>
      </div>
    </div>
  );
};

export default Profile;