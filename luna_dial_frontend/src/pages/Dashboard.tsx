import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Moon,
  LogOut,
  User,
  Plus,
  ChevronDown,
  ChevronRight,
  Calendar,
  Target,
  FileText
} from 'lucide-react';
import useAuthStore from '../store/auth';
import taskService from '../services/task';
import { Task, TaskStatus, PeriodType } from '../types';

const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const { user, logout } = useAuthStore();
  const [currentPeriod, setCurrentPeriod] = useState<PeriodType>('day');
  const [tasks, setTasks] = useState<Task[]>([]);
  const [expandedTasks, setExpandedTasks] = useState<Set<string>>(new Set());
  const [isLoading, setIsLoading] = useState(false);

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  const toggleTaskExpand = (taskId: string) => {
    setExpandedTasks(prev => {
      const newSet = new Set(prev);
      if (newSet.has(taskId)) {
        newSet.delete(taskId);
      } else {
        newSet.add(taskId);
      }
      return newSet;
    });
  };

  const getTaskStatusColor = (status: TaskStatus) => {
    switch (status) {
      case TaskStatus.NotStarted:
        return 'bg-gray-100 text-gray-700';
      case TaskStatus.InProgress:
        return 'bg-blue-100 text-blue-700';
      case TaskStatus.Completed:
        return 'bg-green-100 text-green-700';
      case TaskStatus.Cancelled:
        return 'bg-red-100 text-red-700';
      default:
        return 'bg-gray-100 text-gray-700';
    }
  };

  const getTaskStatusLabel = (status: TaskStatus) => {
    const labels = {
      [TaskStatus.NotStarted]: '未开始',
      [TaskStatus.InProgress]: '进行中',
      [TaskStatus.Completed]: '已完成',
      [TaskStatus.Cancelled]: '已取消',
    };
    return labels[status];
  };

  const renderTaskTree = (task: Task, level: number = 0) => {
    const isExpanded = expandedTasks.has(task.id);

    return (
      <div key={task.id} className="mb-1">
        <div
          className={`flex items-center p-3 hover:bg-gray-50 rounded-lg cursor-pointer`}
          style={{ marginLeft: `${level * 24}px` }}
        >
          {task.has_children && (
            <button
              onClick={() => toggleTaskExpand(task.id)}
              className="mr-2 text-gray-500 hover:text-gray-700"
            >
              {isExpanded ? <ChevronDown className="w-4 h-4" /> : <ChevronRight className="w-4 h-4" />}
            </button>
          )}
          {!task.has_children && <span className="w-6 mr-2" />}

          <span className="mr-3 text-xl">{task.icon || '📋'}</span>
          <span className="flex-1 font-medium text-gray-900">{task.title}</span>

          <span className={`px-2 py-1 text-xs font-medium rounded-full ${getTaskStatusColor(task.status)}`}>
            {getTaskStatusLabel(task.status)}
          </span>
        </div>

        {isExpanded && task.children && task.children.map(child => renderTaskTree(child, level + 1))}
      </div>
    );
  };

  useEffect(() => {
    // 加载初始数据
    const loadTasks = async () => {
      setIsLoading(true);
      try {
        const response = await taskService.getTaskTree();
        setTasks(response.items);
      } catch (error) {
        console.error('Failed to load tasks:', error);
      } finally {
        setIsLoading(false);
      }
    };

    loadTasks();
  }, []);

  return (
    <div className="min-h-screen bg-gray-50">
      {/* 顶部导航栏 */}
      <header className="bg-white shadow-sm border-b border-gray-200">
        <div className="px-6 py-4">
          <div className="flex items-center justify-between">
            {/* Logo和标题 */}
            <div className="flex items-center">
              <Moon className="w-8 h-8 text-indigo-600 mr-3" />
              <h1 className="text-2xl font-bold text-gray-900">Luna Dial</h1>
            </div>

            {/* 周期切换器 */}
            <div className="flex items-center space-x-2">
              {(['day', 'week', 'month', 'quarter', 'year'] as PeriodType[]).map(period => (
                <button
                  key={period}
                  onClick={() => setCurrentPeriod(period)}
                  className={`px-4 py-2 rounded-lg font-medium transition-colors ${
                    currentPeriod === period
                      ? 'bg-indigo-600 text-white'
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                  }`}
                >
                  {period === 'day' && '日'}
                  {period === 'week' && '周'}
                  {period === 'month' && '月'}
                  {period === 'quarter' && '季'}
                  {period === 'year' && '年'}
                </button>
              ))}
            </div>

            {/* 用户信息 */}
            <div className="flex items-center space-x-4">
              <span className="text-gray-700 font-medium">{user?.name || user?.username}</span>
              <button
                onClick={() => navigate('/profile')}
                className="p-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors"
              >
                <User className="w-5 h-5" />
              </button>
              <button
                onClick={handleLogout}
                className="p-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors"
              >
                <LogOut className="w-5 h-5" />
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* 主内容区域 */}
      <main className="container mx-auto px-6 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* 左侧：任务树 */}
          <div className="lg:col-span-2">
            <div className="bg-white rounded-xl shadow-sm border border-gray-200">
              <div className="p-6 border-b border-gray-200">
                <div className="flex items-center justify-between">
                  <h2 className="text-xl font-semibold text-gray-900 flex items-center">
                    <Target className="w-5 h-5 mr-2 text-indigo-600" />
                    任务树
                  </h2>
                  <button className="flex items-center px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors">
                    <Plus className="w-4 h-4 mr-2" />
                    新建任务
                  </button>
                </div>
              </div>

              <div className="p-6">
                {isLoading ? (
                  <div className="text-center py-8 text-gray-500">加载中...</div>
                ) : tasks.length > 0 ? (
                  <div className="space-y-1">
                    {tasks.map(task => renderTaskTree(task))}
                  </div>
                ) : (
                  <div className="text-center py-8 text-gray-500">
                    暂无任务，点击上方按钮创建第一个任务
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* 右侧：日志和统计 */}
          <div className="space-y-6">
            {/* 今日统计 */}
            <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                <Calendar className="w-5 h-5 mr-2 text-indigo-600" />
                今日统计
              </h3>
              <div className="space-y-3">
                <div className="flex justify-between items-center">
                  <span className="text-gray-600">总任务</span>
                  <span className="font-semibold text-gray-900">0</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-600">已完成</span>
                  <span className="font-semibold text-green-600">0</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-600">进行中</span>
                  <span className="font-semibold text-blue-600">0</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-gray-600">努力评分</span>
                  <span className="font-semibold text-gray-900">0</span>
                </div>
              </div>
            </div>

            {/* 最近日志 */}
            <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-semibold text-gray-900 flex items-center">
                  <FileText className="w-5 h-5 mr-2 text-indigo-600" />
                  最近日志
                </h3>
                <button className="text-indigo-600 hover:text-indigo-700 text-sm font-medium">
                  查看全部
                </button>
              </div>
              <div className="text-gray-500 text-center py-4">
                暂无日志
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
};

export default Dashboard;