// 模拟数据
const mockData = {
    currentUser: '用户名',
    currentPeriod: 'day',
    currentDate: new Date(),

    // 任务数据 - 更新score含义为努力程度(0-10)
    tasks: {
        // 年任务
        'task-year-1': {
            id: 'task-year-1',
            title: '2025年度目标',
            type: 'year',
            status: 'inprogress',
            icon: '🎯',
            score: null,
            children: ['task-quarter-1', 'task-quarter-2']
        },
        'task-year-2': {
            id: 'task-year-2',
            title: '技术成长规划',
            type: 'year',
            status: 'notstarted',
            icon: '📚',
            score: null
        },

        // 季度任务
        'task-quarter-1': {
            id: 'task-quarter-1',
            title: 'Q1 第一季度计划',
            type: 'quarter',
            status: 'inprogress',
            icon: '📅',
            score: null,
            parentId: 'task-year-1',
            children: ['task-month-1', 'task-month-2', 'task-month-3']
        },
        'task-quarter-2': {
            id: 'task-quarter-2',
            title: 'Q2 产品发布',
            type: 'quarter',
            status: 'notstarted',
            icon: '🚀',
            score: null,
            parentId: 'task-year-1'
        },

        // 月任务
        'task-month-1': {
            id: 'task-month-1',
            title: '1月份任务',
            type: 'month',
            status: 'inprogress',
            icon: '📆',
            score: null,
            parentId: 'task-quarter-1',
            children: ['task-week-1', 'task-week-2', 'task-week-3', 'task-week-4']
        },
        'task-month-2': {
            id: 'task-month-2',
            title: '2月份计划',
            type: 'month',
            status: 'notstarted',
            icon: '📋',
            score: null,
            parentId: 'task-quarter-1'
        },

        // 周任务
        'task-week-1': {
            id: 'task-week-1',
            title: '第1周任务',
            type: 'week',
            status: 'inprogress',
            icon: '📋',
            score: null,
            parentId: 'task-month-1',
            children: ['task-day-1', 'task-day-2', 'task-day-3']
        },
        'task-week-2': {
            id: 'task-week-2',
            title: '第2周规划',
            type: 'week',
            status: 'notstarted',
            icon: '📝',
            score: null,
            parentId: 'task-month-1'
        },

        // 日任务
        'task-day-1': {
            id: 'task-day-1',
            title: '完成项目文档',
            type: 'day',
            status: 'completed',
            icon: '📝',
            score: 8,
            parentId: 'task-week-1'
        },
        'task-day-2': {
            id: 'task-day-2',
            title: '代码审查',
            type: 'day',
            status: 'inprogress',
            icon: '👨‍💻',
            score: 6,
            parentId: 'task-week-1'
        },
        'task-day-3': {
            id: 'task-day-3',
            title: '系统优化',
            type: 'day',
            status: 'notstarted',
            icon: '🔧',
            score: 0,
            parentId: 'task-week-1'
        },
        'task-day-4': {
            id: 'task-day-4',
            title: '会议准备',
            type: 'day',
            status: 'completed',
            icon: '📊',
            score: 7,
            parentId: 'task-week-1'
        }
    },

    // 日志数据 - 用户自定义的内容
    journals: {
        day: [
            {
                id: 'journal-day-1',
                title: '早晨计划',
                content: '今天要完成三个主要任务...',
                icon: '🌅',
                type: 'day',
                createdAt: new Date('2025-01-06 09:00'),
                updatedAt: new Date('2025-01-06 09:00')
            },
            {
                id: 'journal-day-2',
                title: '晚间总结',
                content: '今天完成了文档编写，代码审查进行了一半...',
                icon: '🌙',
                type: 'day',
                createdAt: new Date('2025-01-06 22:30'),
                updatedAt: new Date('2025-01-06 22:30')
            }
        ],
        week: [
            {
                id: 'journal-week-1',
                title: '第一周回顾',
                content: '本周完成了项目的主要架构设计...',
                icon: '📊',
                type: 'week',
                createdAt: new Date('2025-01-05 20:00'),
                updatedAt: new Date('2025-01-05 20:00')
            }
        ],
        month: [],
        quarter: [],
        year: []
    },

    // 努力程度累计统计
    effortScores: {
        today: 14, // 今日努力总分 (8+6)
        week: 68,  // 本周累计
        month: 245, // 本月累计
        quarter: 780, // 本季累计
        year: 2890 // 本年累计
    },

    // 每日努力程度历史（用于图表）
    dailyEffortHistory: {
        '2025-01-01': 22, // 当天所有任务努力程度总和
        '2025-01-02': 18,
        '2025-01-03': 15,
        '2025-01-04': 25,
        '2025-01-05': 20,
        '2025-01-06': 14
    }
};

// DOM 元素
const periodButtons = document.querySelectorAll('.period-btn');
const taskNodes = document.querySelectorAll('.task-node');

// 初始化
document.addEventListener('DOMContentLoaded', function() {
    // 检查登录状态
    if (window.LunaAuth && !window.LunaAuth.checkAuth()) {
        return; // checkAuth会自动跳转到登录页
    }

    // 加载用户信息
    loadUserInfo();

    // 初始化各个模块
    initPeriodSwitcher();
    initTaskTree();
    initCurrentPeriodTasks();
    initJournals();
    initTimeNavigator();
    updateProgressChart();
    updateDateTime();
    updateScoreDisplay();

    // 初始化登出按钮
    initLogoutButton();
});

// 加载用户信息
async function loadUserInfo() {
    if (window.LunaAuth) {
        try {
            const user = await window.LunaAuth.getCurrentUser();
            if (user) {
                document.getElementById('currentUsername').textContent = user.name || user.username;
            }
        } catch (error) {
            console.error('Failed to load user info:', error);
        }
    }
}

// 初始化登出按钮
function initLogoutButton() {
    const btnLogout = document.getElementById('btnLogout');
    if (btnLogout && window.LunaAuth) {
        btnLogout.addEventListener('click', async () => {
            if (confirm('确定要退出登录吗？')) {
                try {
                    await window.LunaAuth.apiRequest(window.LunaAuth.API_ENDPOINTS.logout, {
                        method: 'POST'
                    });
                } catch (error) {
                    console.error('Logout error:', error);
                } finally {
                    window.LunaAuth.SessionManager.clearSession();
                    window.location.href = 'login.html';
                }
            }
        });
    }
}

// 初始化时间导航器
function initTimeNavigator() {
    const currentYear = new Date().getFullYear();

    // 初始化年份选择器
    const yearPicker = document.getElementById('yearPicker');
    for (let year = currentYear - 5; year <= currentYear + 5; year++) {
        const option = document.createElement('option');
        option.value = year;
        option.textContent = `${year}年`;
        if (year === currentYear) option.selected = true;
        yearPicker.appendChild(option);
    }

    // 初始化月份选择器
    const monthPicker = document.getElementById('monthPicker');
    const months = ['一月', '二月', '三月', '四月', '五月', '六月',
                   '七月', '八月', '九月', '十月', '十一月', '十二月'];
    for (let month = 1; month <= 12; month++) {
        const option = document.createElement('option');
        option.value = month;
        option.textContent = `${months[month-1]} (${month}月)`;
        if (month === new Date().getMonth() + 1) option.selected = true;
        monthPicker.appendChild(option);
    }

    // 初始化季度选择器
    const quarterPicker = document.getElementById('quarterPicker');
    const quarters = [
        { name: 'Q1 第一季度', months: '1-3月' },
        { name: 'Q2 第二季度', months: '4-6月' },
        { name: 'Q3 第三季度', months: '7-9月' },
        { name: 'Q4 第四季度', months: '10-12月' }
    ];
    quarters.forEach((quarter, index) => {
        const option = document.createElement('option');
        option.value = index + 1;
        option.textContent = `${quarter.name} (${quarter.months})`;
        quarterPicker.appendChild(option);
    });

    // 初始化周选择器
    updateWeekPicker();

    // 跳转按钮事件
    document.querySelector('.btn-time-jump').addEventListener('click', function() {
        const period = mockData.currentPeriod;
        showTimeNavigator(period);
    });

    // 日期选择器变化事件
    document.getElementById('datePicker').addEventListener('change', function() {
        const selectedDate = new Date(this.value);
        mockData.currentDate = selectedDate;
        updateDateTime();
        updateDashboard('day');
    });

    // 周选择器变化事件
    document.getElementById('weekPicker').addEventListener('change', function() {
        if (this.value) {
            const [year, week] = this.value.split('-W');
            // 跳转到指定周
            const date = getDateFromWeek(parseInt(year), parseInt(week));
            mockData.currentDate = date;
            updateDateTime();
            updateDashboard('week');
        }
    });

    // 月份选择器变化事件
    document.getElementById('monthPicker').addEventListener('change', function() {
        const year = document.getElementById('yearPicker').value || currentYear;
        const month = this.value;
        if (month) {
            const date = new Date(year, month - 1, 1);
            mockData.currentDate = date;
            updateDateTime();
            updateDashboard('month');
        }
    });

    // 季度选择器变化事件
    document.getElementById('quarterPicker').addEventListener('change', function() {
        const year = document.getElementById('yearPicker').value || currentYear;
        const quarter = this.value;
        if (quarter) {
            const startMonth = (quarter - 1) * 3;
            const date = new Date(year, startMonth, 1);
            mockData.currentDate = date;
            updateDateTime();
            updateDashboard('quarter');
        }
    });

    // 年份选择器变化事件
    document.getElementById('yearPicker').addEventListener('change', function() {
        const year = this.value;
        if (year) {
            const date = new Date(year, 0, 1);
            mockData.currentDate = date;
            updateDateTime();
            updateWeekPicker(); // 更新周选择器
            updateDashboard('year');
        }
    });
}

// 显示对应的时间导航器
function showTimeNavigator(period) {
    // 隐藏所有选择器
    document.querySelectorAll('.time-navigator input, .time-navigator select').forEach(el => {
        el.style.display = 'none';
    });

    // 根据当前周期显示对应的选择器
    switch(period) {
        case 'day':
            document.getElementById('datePicker').style.display = 'block';
            document.getElementById('datePicker').value = formatDate(mockData.currentDate);
            break;
        case 'week':
            document.getElementById('weekPicker').style.display = 'block';
            document.getElementById('yearPicker').style.display = 'block';
            break;
        case 'month':
            document.getElementById('monthPicker').style.display = 'block';
            document.getElementById('yearPicker').style.display = 'block';
            break;
        case 'quarter':
            document.getElementById('quarterPicker').style.display = 'block';
            document.getElementById('yearPicker').style.display = 'block';
            break;
        case 'year':
            document.getElementById('yearPicker').style.display = 'block';
            break;
    }
}

// 更新周选择器
function updateWeekPicker() {
    const weekPicker = document.getElementById('weekPicker');
    const year = parseInt(document.getElementById('yearPicker').value || new Date().getFullYear());

    // 清空现有选项
    weekPicker.innerHTML = '<option value="">选择周</option>';

    // 获取该年第一天
    const yearStart = new Date(year, 0, 1);

    // 找到第一个周一（ISO周的开始）
    let firstMonday = new Date(yearStart);
    const startDay = yearStart.getDay();
    if (startDay === 0) { // 周日
        firstMonday.setDate(yearStart.getDate() + 1);
    } else if (startDay > 1) { // 周二到周六
        firstMonday.setDate(yearStart.getDate() + (8 - startDay));
    }
    // 如果1月1日是周一，那就是第一周的开始

    // 计算该年有多少周
    const yearEnd = new Date(year, 11, 31);
    let weekCount = 0;
    let currentMonday = new Date(firstMonday);

    // 添加每周的选项
    let weekNum = 1;
    while (currentMonday.getFullYear() <= year) {
        if (currentMonday > yearEnd) break;

        const weekEnd = new Date(currentMonday);
        weekEnd.setDate(currentMonday.getDate() + 6); // 周日

        // 格式化日期范围
        const startMonth = currentMonday.getMonth() + 1;
        const startDay = currentMonday.getDate();
        const endMonth = weekEnd.getMonth() + 1;
        const endDay = weekEnd.getDate();

        let rangeText;
        if (startMonth === endMonth) {
            rangeText = `${startMonth}月${startDay}日-${endDay}日`;
        } else {
            rangeText = `${startMonth}月${startDay}日-${endMonth}月${endDay}日`;
        }

        const option = document.createElement('option');
        option.value = `${year}-W${weekNum}`;
        option.textContent = `第${weekNum}周 (${rangeText})`;
        weekPicker.appendChild(option);

        // 移动到下一周
        currentMonday.setDate(currentMonday.getDate() + 7);
        weekNum++;

        // 最多52或53周
        if (weekNum > 53) break;
    }
}

// 根据年份和周数获取日期（ISO周标准）
function getDateFromWeek(year, week) {
    // 获取该年第一天
    const yearStart = new Date(year, 0, 1);

    // 找到第一个周一
    let firstMonday = new Date(yearStart);
    const startDay = yearStart.getDay();
    if (startDay === 0) { // 周日
        firstMonday.setDate(yearStart.getDate() + 1);
    } else if (startDay > 1) { // 周二到周六
        firstMonday.setDate(yearStart.getDate() + (8 - startDay));
    }

    // 计算目标周的周一
    const targetMonday = new Date(firstMonday);
    targetMonday.setDate(firstMonday.getDate() + (week - 1) * 7);

    return targetMonday;
}

// 周期切换功能
function initPeriodSwitcher() {
    periodButtons.forEach(btn => {
        btn.addEventListener('click', function() {
            periodButtons.forEach(b => b.classList.remove('active'));
            this.classList.add('active');

            const period = this.dataset.period;
            mockData.currentPeriod = period;
            updateDashboard(period);
        });
    });
}

// 初始化任务树
function initTaskTree() {
    const taskToggles = document.querySelectorAll('.task-toggle');

    taskToggles.forEach(toggle => {
        toggle.addEventListener('click', function(e) {
            e.stopPropagation();

            const taskNode = this.parentElement;
            const nextElement = taskNode.nextElementSibling;

            if (nextElement && nextElement.classList.contains('task-children')) {
                nextElement.classList.toggle('hidden');
                this.textContent = nextElement.classList.contains('hidden') ? '▶' : '▼';
                this.classList.toggle('expanded');
            }
        });
    });

    // 为任务树中的状态选择器添加事件
    const treeStatusSelects = document.querySelectorAll('.task-status-mini');
    treeStatusSelects.forEach(select => {
        select.addEventListener('click', function(e) {
            e.stopPropagation(); // 防止触发任务节点的点击事件
        });

        select.addEventListener('change', function() {
            const taskId = this.dataset.taskId;
            const newStatus = this.value;

            // 更新数据
            if (mockData.tasks[taskId]) {
                mockData.tasks[taskId].status = newStatus;
            }

            // 如果切换的任务显示在当前周期任务列表中，更新对应的显示
            updateCurrentPeriodTasks();
        });
    });
}

// 初始化当前周期任务（原本的日任务）
function initCurrentPeriodTasks() {
    // 状态选择器
    const statusSelects = document.querySelectorAll('.daily-task .task-status-select');
    statusSelects.forEach(select => {
        select.addEventListener('change', function() {
            const taskId = this.dataset.taskId;
            const newStatus = this.value;
            const scoreControl = this.parentElement.querySelector('.score-control');

            // 更新任务状态
            if (mockData.tasks[taskId]) {
                mockData.tasks[taskId].status = newStatus;

                // 对于日任务，根据状态启用/禁用评分
                if (mockData.tasks[taskId].type === 'day') {
                    const scoreInput = scoreControl.querySelector('.score-input');
                    if (newStatus === 'notstarted' || newStatus === 'cancelled') {
                        scoreControl.classList.add('disabled');
                        scoreInput.disabled = true;
                        scoreInput.value = 0;
                        scoreControl.querySelector('.score-display').textContent = '-';
                    } else {
                        scoreControl.classList.remove('disabled');
                        scoreInput.disabled = false;
                    }
                }
            }

            updateTotalScore();
        });
    });

    // 努力程度评分输入（只对日任务）
    const scoreInputs = document.querySelectorAll('.score-input');
    scoreInputs.forEach(input => {
        input.addEventListener('input', function() {
            const value = Math.min(10, Math.max(0, parseInt(this.value) || 0));
            this.value = value;

            const taskId = this.dataset.taskId;
            const displayElement = this.parentElement.querySelector('.score-display');
            displayElement.textContent = value;

            // 更新数据
            if (mockData.tasks[taskId]) {
                mockData.tasks[taskId].score = value;
            }

            updateTotalScore();
        });
    });
}

// 初始化日志
function initJournals() {
    // 新建日志按钮
    document.querySelectorAll('.btn-new-journal').forEach(btn => {
        btn.addEventListener('click', function() {
            alert('新建日志功能 - 待实现\n用户可以输入：标题、内容、选择图标');
        });
    });

    // 查看日志按钮
    document.querySelectorAll('.btn-view-journal').forEach(btn => {
        btn.addEventListener('click', function() {
            const journalEntry = this.closest('.journal-entry');
            const title = journalEntry.querySelector('.journal-title').textContent;
            const time = journalEntry.querySelector('.journal-time').textContent;
            const icon = journalEntry.querySelector('.journal-icon').textContent;

            // 根据当前周期查找对应的日志数据
            const period = mockData.currentPeriod;
            const journals = mockData.journals[period] || [];
            const journal = journals.find(j => j.title === title);

            if (journal) {
                openJournalModal(journal);
            }
        });
    });

    // 编辑日志按钮
    document.querySelectorAll('.btn-edit-journal').forEach(btn => {
        btn.addEventListener('click', function() {
            const journalEntry = this.closest('.journal-entry');
            const title = journalEntry.querySelector('.journal-title').textContent;
            alert(`编辑日志: ${title}\n功能待实现`);
        });
    });

    // 删除日志按钮
    document.querySelectorAll('.btn-delete-journal').forEach(btn => {
        btn.addEventListener('click', function() {
            if (confirm('确定删除这篇日志吗？')) {
                const journalEntry = this.closest('.journal-entry');
                journalEntry.remove();

                // 检查是否还有日志
                const journalList = document.querySelector('.journal-list');
                if (journalList.children.length === 0) {
                    journalList.innerHTML = `
                        <div class="journal-empty">
                            <p>暂无日志记录</p>
                            <button class="btn-new-journal">创建第一篇日志</button>
                        </div>
                    `;
                }
            }
        });
    });
}

// 打开日志查看模态框
function openJournalModal(journal) {
    const modal = document.getElementById('journalModal');
    const modalIcon = document.getElementById('modalIcon');
    const modalTitle = document.getElementById('modalTitle');
    const modalTime = document.getElementById('modalTime');
    const modalContent = document.getElementById('modalContent');

    modalIcon.textContent = journal.icon;
    modalTitle.textContent = journal.title;
    modalTime.textContent = `创建于 ${formatDateTime(journal.createdAt)}`;
    modalContent.textContent = journal.content;

    modal.style.display = 'flex';

    // 点击遮罩关闭
    modal.onclick = function(e) {
        if (e.target === modal) {
            closeJournalModal();
        }
    };
}

// 关闭日志查看模态框
function closeJournalModal() {
    const modal = document.getElementById('journalModal');
    modal.style.display = 'none';
}

// 格式化日期时间
function formatDateTime(date) {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${year}-${month}-${day} ${hours}:${minutes}`;
}

// 更新总分数
function updateTotalScore() {
    let todayTotal = 0;

    // 计算所有日任务的努力程度总和
    document.querySelectorAll('.score-input:not(:disabled)').forEach(input => {
        todayTotal += parseInt(input.value) || 0;
    });

    mockData.effortScores.today = todayTotal;
    updateScoreDisplay();
}

// 更新分数显示
function updateScoreDisplay() {
    const scoreValues = document.querySelectorAll('.score-value');
    const period = mockData.currentPeriod;

    const periodScores = {
        day: [mockData.effortScores.today, mockData.effortScores.week, mockData.effortScores.month],
        week: [mockData.effortScores.week, mockData.effortScores.month, mockData.effortScores.quarter],
        month: [mockData.effortScores.month, mockData.effortScores.quarter, mockData.effortScores.year],
        quarter: [mockData.effortScores.quarter, mockData.effortScores.year, mockData.effortScores.year * 2],
        year: [mockData.effortScores.year, mockData.effortScores.year * 2, mockData.effortScores.year * 3]
    };

    scoreValues.forEach((value, index) => {
        if (periodScores[period] && periodScores[period][index]) {
            value.textContent = periodScores[period][index];
        }
    });
}

// 更新进度图表 - 显示努力程度趋势
function updateProgressChart() {
    const bars = document.querySelectorAll('.bar-fill');
    const weekDays = getWeekDays();

    bars.forEach((bar, index) => {
        const date = weekDays[index];
        const effort = mockData.dailyEffortHistory[date] || 0;
        // 假设每日最大努力值为30分（3个任务各10分）
        const height = Math.min(100, (effort / 30) * 100);
        bar.style.height = height + '%';
    });
}

// 获取本周日期
function getWeekDays() {
    const today = new Date();
    const currentDay = today.getDay();
    const weekStart = new Date(today);
    weekStart.setDate(today.getDate() - (currentDay === 0 ? 6 : currentDay - 1));

    const weekDays = [];
    for (let i = 0; i < 7; i++) {
        const date = new Date(weekStart);
        date.setDate(weekStart.getDate() + i);
        weekDays.push(formatDate(date));
    }

    return weekDays;
}

// 格式化日期
function formatDate(date) {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
}

// 更新日期时间显示
function updateDateTime() {
    const dateElement = document.querySelector('.current-date');
    const now = new Date();
    const weekDays = ['日', '一', '二', '三', '四', '五', '六'];

    const year = now.getFullYear();
    const month = now.getMonth() + 1;
    const date = now.getDate();
    const weekDay = weekDays[now.getDay()];

    dateElement.textContent = `${year}年${month}月${date}日 星期${weekDay}`;
}

// 更新 Dashboard 内容（根据周期）
function updateDashboard(period) {
    // 更新标题
    const focusTitle = document.querySelector('.focus-card h3');
    const periodTitles = {
        day: '今日任务',
        week: '本周任务',
        month: '本月任务',
        quarter: '本季任务',
        year: '本年任务'
    };
    focusTitle.textContent = periodTitles[period];

    // 更新任务列表
    updateCurrentPeriodTasks();

    // 更新分数标签
    const scoreLabels = document.querySelectorAll('.score-label');
    const periodNames = {
        day: ['今日努力总分', '本周累计', '本月累计'],
        week: ['本周努力总分', '本月累计', '本季累计'],
        month: ['本月努力总分', '本季累计', '本年累计'],
        quarter: ['本季努力总分', '本年累计', '总累计'],
        year: ['本年努力总分', '总累计', '历史最高']
    };

    scoreLabels.forEach((label, index) => {
        if (periodNames[period] && periodNames[period][index]) {
            label.textContent = periodNames[period][index];
        }
    });

    // 更新日志标题和内容
    const journalTitle = document.querySelector('.journal-card h3');
    const journalPeriodNames = {
        day: '今日日志',
        week: '本周日志',
        month: '本月日志',
        quarter: '本季日志',
        year: '本年日志'
    };
    journalTitle.textContent = journalPeriodNames[period];

    // 更新日志列表
    updateJournalList(period);

    // 更新图表标题
    const chartTitle = document.querySelector('.progress-card h3');
    const chartTitles = {
        day: '本周努力趋势',
        week: '本月努力趋势',
        month: '本季努力趋势',
        quarter: '本年努力趋势',
        year: '历年努力对比'
    };
    chartTitle.textContent = chartTitles[period];

    updateScoreDisplay();
}

// 更新当前周期的任务列表
function updateCurrentPeriodTasks() {
    const period = mockData.currentPeriod;
    const tasksContainer = document.querySelector('.today-tasks');

    // 根据周期类型筛选对应的任务
    const periodTaskTypes = {
        day: 'day',
        week: 'week',
        month: 'month',
        quarter: 'quarter',
        year: 'year'
    };

    const targetType = periodTaskTypes[period];
    const tasks = Object.values(mockData.tasks).filter(task => task.type === targetType);

    // 构建HTML
    let html = '';
    tasks.forEach(task => {
        if (task.type === 'day') {
            // 日任务有评分控件
            html += `
                <div class="daily-task">
                    <div class="task-info">
                        <span class="task-icon">${task.icon}</span>
                        <span class="task-text">${task.title}</span>
                    </div>
                    <div class="task-controls">
                        <select class="task-status-select" data-task-id="${task.id}">
                            <option value="notstarted" ${task.status === 'notstarted' ? 'selected' : ''}>未开始</option>
                            <option value="inprogress" ${task.status === 'inprogress' ? 'selected' : ''}>进行中</option>
                            <option value="completed" ${task.status === 'completed' ? 'selected' : ''}>已完成</option>
                            <option value="cancelled" ${task.status === 'cancelled' ? 'selected' : ''}>已取消</option>
                        </select>
                        <div class="score-control ${(task.status === 'notstarted' || task.status === 'cancelled') ? 'disabled' : ''}">
                            <label>努力程度:</label>
                            <input type="number" class="score-input" min="0" max="10" value="${task.score || 0}"
                                ${(task.status === 'notstarted' || task.status === 'cancelled') ? 'disabled' : ''}
                                data-task-id="${task.id}">
                            <span class="score-display">${task.score || '-'}</span>
                        </div>
                    </div>
                </div>
            `;
        } else {
            // 其他层级任务没有评分控件
            html += `
                <div class="daily-task">
                    <div class="task-info">
                        <span class="task-icon">${task.icon}</span>
                        <span class="task-text">${task.title}</span>
                    </div>
                    <div class="task-controls">
                        <select class="task-status-select" data-task-id="${task.id}">
                            <option value="notstarted" ${task.status === 'notstarted' ? 'selected' : ''}>未开始</option>
                            <option value="inprogress" ${task.status === 'inprogress' ? 'selected' : ''}>进行中</option>
                            <option value="completed" ${task.status === 'completed' ? 'selected' : ''}>已完成</option>
                            <option value="cancelled" ${task.status === 'cancelled' ? 'selected' : ''}>已取消</option>
                        </select>
                    </div>
                </div>
            `;
        }
    });

    // 如果没有任务，显示空状态
    if (tasks.length === 0) {
        html = '<div style="text-align: center; color: #6e6e73; padding: 2rem;">暂无任务</div>';
    }

    tasksContainer.innerHTML = html;

    // 重新绑定事件
    initCurrentPeriodTasks();
}

// 更新日志列表
function updateJournalList(period) {
    const journalList = document.querySelector('.journal-list');
    const journals = mockData.journals[period] || [];

    if (journals.length === 0) {
        journalList.innerHTML = `
            <div class="journal-empty">
                <p>暂无日志记录</p>
                <button class="btn-new-journal">创建第一篇日志</button>
            </div>
        `;
    } else {
        journalList.innerHTML = journals.map(journal => `
            <div class="journal-entry">
                <div class="journal-content">
                    <span class="journal-icon">${journal.icon}</span>
                    <div class="journal-info">
                        <h4 class="journal-title">${journal.title}</h4>
                        <span class="journal-time">${formatTime(journal.createdAt)}</span>
                    </div>
                </div>
                <div class="journal-actions">
                    <button class="btn-view-journal">查看</button>
                    <button class="btn-edit-journal">编辑</button>
                    <button class="btn-delete-journal">删除</button>
                </div>
            </div>
        `).join('');
    }

    // 重新绑定日志事件
    initJournals();
}

// 格式化时间
function formatTime(date) {
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${hours}:${minutes}`;
}

// 其他按钮事件
document.querySelector('.btn-create-task').addEventListener('click', function() {
    alert('创建任务功能 - 待实现');
});

document.querySelector('.btn-add-task').addEventListener('click', function() {
    alert('添加今日任务功能 - 待实现');
});

document.querySelectorAll('.action-btn').forEach(btn => {
    btn.addEventListener('click', function() {
        const action = this.querySelector('span:last-child').textContent;
        alert(`${action}功能 - 待实现`);
    });
});

document.querySelector('.btn-logout').addEventListener('click', function() {
    if (confirm('确定要退出登录吗？')) {
        alert('已退出登录');
    }
});