import React from 'react';

// 侧边栏分组条目
export interface TaskGroupSidebarItem {
  id: string;
  title: string;
  count: number;
}

// 侧边栏组件参数
interface TaskGroupSidebarProps {
  items: TaskGroupSidebarItem[];
  activeId: string | null;
  onSelect: (id: string) => void;
}

const TaskGroupSidebar: React.FC<TaskGroupSidebarProps> = ({ items, activeId, onSelect }) => {
  return (
    <nav className="task-group-sidebar" aria-label="分组定位器">
      <div className="task-group-sidebar-title">分组</div>
      <div className="task-group-sidebar-list">
        {items.map(item => (
          <button
            key={item.id}
            type="button"
            className={`task-group-sidebar-item ${activeId === item.id ? 'active' : ''}`}
            onClick={() => onSelect(item.id)}
            title={item.title}
          >
            <span className="task-group-sidebar-item-text">{item.title}</span>
            <span className="task-group-sidebar-item-count">{item.count}</span>
          </button>
        ))}
      </div>
    </nav>
  );
};

export default TaskGroupSidebar;

