import React from 'react';
import AncestryBreadcrumb, { BreadcrumbItem, BreadcrumbStatus } from './AncestryBreadcrumb';

// 分组卡片参数
interface TaskGroupCardProps {
  title: string;
  breadcrumbItems: BreadcrumbItem[];
  breadcrumbStatus: BreadcrumbStatus;
  metaText?: string;
  onRetryBreadcrumb?: () => void;
  children: React.ReactNode;
}

const TaskGroupCard: React.FC<TaskGroupCardProps> = ({
  title,
  breadcrumbItems,
  breadcrumbStatus,
  metaText,
  onRetryBreadcrumb,
  children
}) => {
  return (
    <section className="task-group-card" aria-label={`任务分组：${title}`}>
      <div className="task-group-header">
        <div className="task-group-header-main">
          <div className="task-group-title">{title}</div>
          <AncestryBreadcrumb
            items={breadcrumbItems}
            status={breadcrumbStatus}
            onRetry={onRetryBreadcrumb}
          />
        </div>
        {metaText && <div className="task-group-meta">{metaText}</div>}
      </div>
      <div className="task-group-body">{children}</div>
    </section>
  );
};

export default TaskGroupCard;

