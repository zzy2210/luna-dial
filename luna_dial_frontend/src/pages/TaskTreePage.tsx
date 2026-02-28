import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import TaskTree from '../components/TaskTree';
import TaskViewDialog from '../components/TaskViewDialog';
import TaskEditDialog from '../components/TaskEditDialog';
import taskService from '../services/task';
import { Task, TaskStatus } from '../types';
import { TASK_STATUS_TO_API } from '../constants/task';
import '../styles/task-tree-page.css';

const TaskTreePage: React.FC = () => {
  const navigate = useNavigate();
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const pageSize = 10;

  // 对话框状态
  const [showViewTaskDialog, setShowViewTaskDialog] = useState(false);
  const [showTaskDialog, setShowTaskDialog] = useState(false);
  const [viewingTask, setViewingTask] = useState<Task | null>(null);
  const [editingTask, setEditingTask] = useState<Task | null>(null);
  const [parentTaskIdForCreation, setParentTaskIdForCreation] = useState<string | undefined>(undefined);

  useEffect(() => {
    loadTaskTree();
  }, [page]);

  const loadTaskTree = async () => {
    setLoading(true);
    try {
      const response = await taskService.getTaskTree(page, pageSize, true);
      setTasks(response.items || []);
      setTotalPages(response.pagination?.total_pages || 1);
    } catch (error) {
      console.error('Failed to load task tree:', error);
      setTasks([]);
    } finally {
      setLoading(false);
    }
  };

  const handleTaskStatusChange = async (taskId: string, status: TaskStatus) => {
    try {
      await taskService.updateTask(taskId, {
        status: TASK_STATUS_TO_API[status]
      });

      // 刷新任务树
      loadTaskTree();
    } catch (error) {
      console.error('Failed to update task status:', error);
    }
  };

  const handleTaskClick = (task: Task) => {
    setViewingTask(task);
    setShowViewTaskDialog(true);
  };

  const handleCreateTask = () => {
    setEditingTask(null);
    setParentTaskIdForCreation(undefined);
    setShowTaskDialog(true);
  };

  const handleCreateSubtask = (parentTaskId: string) => {
    setParentTaskIdForCreation(parentTaskId);
    setEditingTask(null);
    setShowViewTaskDialog(false);
    setShowTaskDialog(true);
  };

  const handleEditTask = (task: Task) => {
    setEditingTask(task);
    setShowViewTaskDialog(false);
    setShowTaskDialog(true);
  };

  const handleDeleteTask = async (taskId: string) => {
    try {
      await taskService.deleteTask(taskId);
      setShowViewTaskDialog(false);
      loadTaskTree();
    } catch (error) {
      console.error('Failed to delete task:', error);
      alert('删除任务失败，请重试');
    }
  };

  const handleTaskScoreChange = async (taskId: string, score: number) => {
    try {
      await taskService.updateScore(taskId, score);
      loadTaskTree();
    } catch (error) {
      console.error('Failed to update task score:', error);
    }
  };

  const handlePrevPage = () => {
    if (page > 1) {
      setPage(page - 1);
    }
  };

  const handleNextPage = () => {
    if (page < totalPages) {
      setPage(page + 1);
    }
  };

  return (
    <div className="task-tree-page">
      {/* 顶部导航栏 */}
      <header className="task-tree-navbar">
        <div className="navbar-left">
          <button className="btn-back" onClick={() => navigate('/dashboard')}>
            ← 返回
          </button>
          <h1>任务树</h1>
        </div>
        <div className="navbar-right">
          <button className="btn-create-root-task" onClick={handleCreateTask}>
            + 新建根任务
          </button>
        </div>
      </header>

      {/* 主内容区 */}
      <main className="task-tree-content">
        <div className="task-tree-container">
          {loading ? (
            <div className="loading-state">加载中...</div>
          ) : tasks.length > 0 ? (
            <>
              <TaskTree
                tasks={tasks}
                onTaskStatusChange={handleTaskStatusChange}
                onTaskClick={handleTaskClick}
              />

              {/* 分页控制 */}
              {totalPages > 1 && (
                <div className="pagination">
                  <button
                    className="btn-page"
                    onClick={handlePrevPage}
                    disabled={page === 1}
                  >
                    上一页
                  </button>
                  <span className="page-info">
                    第 {page} / {totalPages} 页
                  </span>
                  <button
                    className="btn-page"
                    onClick={handleNextPage}
                    disabled={page === totalPages}
                  >
                    下一页
                  </button>
                </div>
              )}
            </>
          ) : (
            <div className="empty-state">
              <div className="empty-icon">📋</div>
              <div className="empty-text">暂无任务</div>
              <button className="btn-create-first" onClick={handleCreateTask}>
                创建第一个任务
              </button>
            </div>
          )}
        </div>
      </main>

      {/* 任务查看对话框 */}
      {showViewTaskDialog && viewingTask && (
        <TaskViewDialog
          task={viewingTask}
          onClose={() => setShowViewTaskDialog(false)}
          onEdit={(task) => handleEditTask(task)}
          onDelete={(taskId) => handleDeleteTask(taskId)}
          onScoreUpdate={handleTaskScoreChange}
          onCreateSubtask={handleCreateSubtask}
        />
      )}

      {/* 任务编辑对话框 */}
      {showTaskDialog && (
        <TaskEditDialog
          task={editingTask}
          onClose={() => {
            setShowTaskDialog(false);
            setParentTaskIdForCreation(undefined);
          }}
          onSuccess={() => {
            setShowTaskDialog(false);
            setParentTaskIdForCreation(undefined);
            loadTaskTree();
          }}
          currentPeriod="day"
          parentTaskId={parentTaskIdForCreation}
        />
      )}
    </div>
  );
};

export default TaskTreePage;
