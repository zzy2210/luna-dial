// API 配置
const USE_MOCK = true; // 设置为true使用mock模式，false使用真实API
const API_BASE_URL = 'http://localhost:8081';
const API_ENDPOINTS = {
    login: '/api/v1/public/auth/login',
    logout: '/api/v1/auth/logout',
    logoutAll: '/api/v1/auth/logout-all',
    profile: '/api/v1/auth/profile',
    userDetail: '/api/v1/users/me'
};

// Mock数据
const MOCK_USERS = {
    'admin': {
        password: '123456',
        userInfo: {
            user_id: 'user_001',
            username: 'admin',
            name: '管理员',
            email: 'admin@lunadial.com',
            created_at: '2025-01-01T00:00:00Z',
            updated_at: '2025-01-06T10:00:00Z'
        }
    },
    'test': {
        password: 'test123',
        userInfo: {
            user_id: 'user_002',
            username: 'test',
            name: '测试用户',
            email: 'test@lunadial.com',
            created_at: '2025-01-02T00:00:00Z',
            updated_at: '2025-01-06T10:00:00Z'
        }
    },
    'demo': {
        password: 'demo',
        userInfo: {
            user_id: 'user_003',
            username: 'demo',
            name: '演示用户',
            email: 'demo@lunadial.com',
            created_at: '2025-01-03T00:00:00Z',
            updated_at: '2025-01-06T10:00:00Z'
        }
    }
};

// Session 管理
const SessionManager = {
    // 获取存储的session
    getSession() {
        const sessionData = localStorage.getItem('luna_session');
        if (!sessionData) return null;

        try {
            const session = JSON.parse(sessionData);
            // 检查是否过期
            if (session.expiresAt && new Date(session.expiresAt) < new Date()) {
                this.clearSession();
                return null;
            }
            return session;
        } catch (e) {
            this.clearSession();
            return null;
        }
    },

    // 保存session
    saveSession(sessionId, expiresIn) {
        const expiresAt = new Date(Date.now() + expiresIn * 1000);
        const sessionData = {
            sessionId,
            expiresAt: expiresAt.toISOString(),
            createdAt: new Date().toISOString()
        };
        localStorage.setItem('luna_session', JSON.stringify(sessionData));
    },

    // 清除session
    clearSession() {
        localStorage.removeItem('luna_session');
    },

    // 获取认证头
    getAuthHeader() {
        const session = this.getSession();
        if (!session) return null;
        return `Bearer ${session.sessionId}`;
    }
};

// API 请求封装
async function apiRequest(endpoint, options = {}) {
    // Mock模式处理
    if (USE_MOCK) {
        return mockApiRequest(endpoint, options);
    }

    const url = API_BASE_URL + endpoint;
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers
    };

    // 添加认证头（如果需要）
    if (!endpoint.includes('/public/')) {
        const authHeader = SessionManager.getAuthHeader();
        if (authHeader) {
            headers['Authorization'] = authHeader;
        } else {
            // 未登录，跳转到登录页
            window.location.href = 'login.html';
            return;
        }
    }

    try {
        const response = await fetch(url, {
            ...options,
            headers
        });

        const data = await response.json();

        if (!response.ok) {
            if (response.status === 401) {
                // 认证失败，清除session并跳转到登录页
                SessionManager.clearSession();
                window.location.href = 'login.html';
                return;
            }
            throw new Error(data.message || `HTTP error! status: ${response.status}`);
        }

        return data;
    } catch (error) {
        console.error('API request failed:', error);
        throw error;
    }
}

// Mock API请求
async function mockApiRequest(endpoint, options = {}) {
    // 模拟网络延迟
    await new Promise(resolve => setTimeout(resolve, 300));

    // 登录接口
    if (endpoint === API_ENDPOINTS.login && options.method === 'POST') {
        const { username, password } = JSON.parse(options.body);
        const user = MOCK_USERS[username];

        if (user && user.password === password) {
            return {
                success: true,
                code: 200,
                message: 'Login successful',
                data: {
                    session_id: 'mock_session_' + Date.now(),
                    expires_in: 86400
                }
            };
        } else {
            throw new Error('用户名或密码错误');
        }
    }

    // 获取用户资料
    if (endpoint === API_ENDPOINTS.profile) {
        const session = SessionManager.getSession();
        if (!session) {
            throw new Error('未登录');
        }

        // 从localStorage获取当前登录的用户名
        const currentUser = localStorage.getItem('mock_current_user');
        const user = MOCK_USERS[currentUser];

        if (user) {
            return {
                success: true,
                code: 200,
                message: 'success',
                data: user.userInfo
            };
        }
    }

    // 获取用户详细信息
    if (endpoint === API_ENDPOINTS.userDetail) {
        const session = SessionManager.getSession();
        if (!session) {
            throw new Error('未登录');
        }

        const currentUser = localStorage.getItem('mock_current_user');
        const user = MOCK_USERS[currentUser];

        if (user) {
            const sessionId = session.sessionId;
            const now = new Date();
            return {
                success: true,
                code: 200,
                message: 'success',
                data: {
                    ...user.userInfo,
                    session: {
                        session_id: sessionId,
                        last_access_at: now.toISOString(),
                        expires_at: session.expiresAt
                    }
                }
            };
        }
    }

    // 登出接口
    if (endpoint === API_ENDPOINTS.logout && options.method === 'POST') {
        return {
            success: true,
            code: 200,
            message: 'Logout successful'
        };
    }

    // 登出所有设备
    if (endpoint === API_ENDPOINTS.logoutAll && options.method === 'DELETE') {
        return {
            success: true,
            code: 200,
            message: 'Logged out from all devices'
        };
    }

    return {
        success: true,
        code: 200,
        message: 'Mock response'
    };
}

// ==================== 登录页面逻辑 ====================
if (document.getElementById('loginForm')) {
    // 检查是否已登录
    if (SessionManager.getSession()) {
        window.location.href = 'dashboard.html';
    }

    const loginForm = document.getElementById('loginForm');
    const errorMessage = document.getElementById('errorMessage');
    const errorText = document.getElementById('errorText');
    const loginButton = document.getElementById('loginButton');
    const btnText = loginButton.querySelector('.btn-text');
    const btnLoading = loginButton.querySelector('.btn-loading');

    // 处理登录表单提交
    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        // 清除之前的错误信息
        errorMessage.style.display = 'none';

        // 获取表单数据
        const username = document.getElementById('username').value.trim();
        const password = document.getElementById('password').value;
        const rememberMe = document.getElementById('rememberMe').checked;

        // 基本验证
        if (!username || !password) {
            showError('请输入用户名和密码');
            return;
        }

        // 显示加载状态
        btnText.style.display = 'none';
        btnLoading.style.display = 'flex';
        loginButton.disabled = true;

        try {
            // 发送登录请求
            const response = await apiRequest(API_ENDPOINTS.login, {
                method: 'POST',
                body: JSON.stringify({ username, password })
            });

            if (response.success) {
                // 保存session
                const { session_id, expires_in } = response.data;
                SessionManager.saveSession(session_id, expires_in);

                // Mock模式下保存用户名
                if (USE_MOCK) {
                    localStorage.setItem('mock_current_user', username);
                }

                // 如果勾选了记住我，延长session有效期（这里只是示例，实际需要后端支持）
                if (rememberMe) {
                    // TODO: 发送请求延长session有效期
                }

                // 跳转到主页
                window.location.href = 'dashboard.html';
            } else {
                showError(response.message || '登录失败');
            }
        } catch (error) {
            showError(error.message || '登录失败，请稍后重试');
        } finally {
            // 恢复按钮状态
            btnText.style.display = 'block';
            btnLoading.style.display = 'none';
            loginButton.disabled = false;
        }
    });

    function showError(message) {
        errorText.textContent = message;
        errorMessage.style.display = 'flex';
    }
}

// ==================== 个人设置页面逻辑 ====================
async function loadUserProfile() {
    // 检查是否已登录
    if (!SessionManager.getSession()) {
        window.location.href = 'login.html';
        return;
    }

    const loadingOverlay = document.getElementById('loadingOverlay');
    if (loadingOverlay) {
        loadingOverlay.style.display = 'flex';
    }

    try {
        // 获取用户详细信息
        const response = await apiRequest(API_ENDPOINTS.userDetail);

        if (response.success) {
            const userData = response.data;

            // 更新基本信息
            updateElementText('username', userData.username);
            updateElementText('name', userData.name);
            updateElementText('email', userData.email || '未设置');
            updateElementText('userId', userData.user_id);

            // 更新账户信息
            updateElementText('createdAt', formatDateTime(userData.created_at));
            updateElementText('updatedAt', formatDateTime(userData.updated_at));

            // 更新会话信息
            if (userData.session) {
                updateElementText('sessionId', userData.session.session_id);
                updateElementText('lastAccess', formatDateTime(userData.session.last_access_at));
                updateElementText('expiresAt', formatDateTime(userData.session.expires_at));

                // 计算剩余时间
                const expiresAt = new Date(userData.session.expires_at);
                const now = new Date();
                const remainingMs = expiresAt - now;
                if (remainingMs > 0) {
                    const hours = Math.floor(remainingMs / (1000 * 60 * 60));
                    const minutes = Math.floor((remainingMs % (1000 * 60 * 60)) / (1000 * 60));
                    updateElementText('remainingTime', `${hours}小时${minutes}分钟`);
                } else {
                    updateElementText('remainingTime', '已过期');
                }
            }

            // TODO: 获取统计信息（需要额外的API调用）
            updateElementText('totalTasks', '12');
            updateElementText('completedTasks', '8');
            updateElementText('totalJournals', '15');
            updateElementText('totalEffort', '156');
        }
    } catch (error) {
        console.error('Failed to load user profile:', error);
        alert('获取用户信息失败，请稍后重试');
    } finally {
        if (loadingOverlay) {
            loadingOverlay.style.display = 'none';
        }
    }
}

// 更新元素文本
function updateElementText(id, text) {
    const element = document.getElementById(id);
    if (element) {
        element.textContent = text || '-';
    }
}

// 格式化日期时间
function formatDateTime(dateString) {
    if (!dateString) return '-';

    const date = new Date(dateString);
    return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
}

// 登出功能
if (document.getElementById('btnLogout')) {
    document.getElementById('btnLogout').addEventListener('click', async () => {
        showConfirmDialog('确认退出', '确定要退出登录吗？', async () => {
            try {
                await apiRequest(API_ENDPOINTS.logout, { method: 'POST' });
            } catch (error) {
                console.error('Logout error:', error);
            } finally {
                SessionManager.clearSession();
                if (USE_MOCK) {
                    localStorage.removeItem('mock_current_user');
                }
                window.location.href = 'login.html';
            }
        });
    });
}

// 登出所有设备
if (document.getElementById('btnLogoutAll')) {
    document.getElementById('btnLogoutAll').addEventListener('click', async () => {
        showConfirmDialog('退出所有设备', '确定要退出所有设备的登录状态吗？此操作将清除所有会话。', async () => {
            try {
                const response = await apiRequest(API_ENDPOINTS.logoutAll, { method: 'DELETE' });
                if (response.success) {
                    alert('已成功退出所有设备');
                    SessionManager.clearSession();
                    if (USE_MOCK) {
                        localStorage.removeItem('mock_current_user');
                    }
                    window.location.href = 'login.html';
                }
            } catch (error) {
                console.error('Logout all error:', error);
                alert('操作失败，请稍后重试');
            }
        });
    });
}

// 确认对话框
function showConfirmDialog(title, message, onConfirm) {
    const modal = document.getElementById('confirmModal');
    const modalTitle = document.getElementById('modalTitle');
    const modalMessage = document.getElementById('modalMessage');
    const btnCancel = document.getElementById('btnCancel');
    const btnConfirm = document.getElementById('btnConfirm');

    if (!modal) return;

    modalTitle.textContent = title;
    modalMessage.textContent = message;
    modal.style.display = 'flex';

    // 取消按钮
    btnCancel.onclick = () => {
        modal.style.display = 'none';
    };

    // 确认按钮
    btnConfirm.onclick = async () => {
        modal.style.display = 'none';
        if (onConfirm) {
            await onConfirm();
        }
    };

    // 点击遮罩关闭
    modal.onclick = (e) => {
        if (e.target === modal) {
            modal.style.display = 'none';
        }
    };
}

// 检查认证状态（用于其他页面）
function checkAuth() {
    const session = SessionManager.getSession();
    if (!session) {
        // 未登录，跳转到登录页
        window.location.href = 'login.html';
        return false;
    }
    return true;
}

// 获取当前用户信息（用于dashboard等页面）
async function getCurrentUser() {
    try {
        const response = await apiRequest(API_ENDPOINTS.profile);
        if (response.success) {
            return response.data;
        }
    } catch (error) {
        console.error('Failed to get user info:', error);
    }
    return null;
}

// 导出供其他页面使用
window.LunaAuth = {
    SessionManager,
    apiRequest,
    checkAuth,
    getCurrentUser,
    API_BASE_URL,
    API_ENDPOINTS
};