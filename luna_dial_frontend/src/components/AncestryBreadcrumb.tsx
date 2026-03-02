import React from 'react';

// 面包屑条目
export interface BreadcrumbItem {
  id: string;
  title: string;
}

// 面包屑展示状态
export type BreadcrumbStatus = 'loading' | 'error' | 'ready' | 'empty';

// 面包屑组件参数
interface AncestryBreadcrumbProps {
  items: BreadcrumbItem[];
  status: BreadcrumbStatus;
  onRetry?: () => void;
}

const AncestryBreadcrumb: React.FC<AncestryBreadcrumbProps> = ({ items, status, onRetry }) => {
  if (status === 'loading') {
    return <div className="breadcrumb-status">路径加载中…</div>;
  }

  if (status === 'error') {
    return (
      <div className="breadcrumb-status">
        <span>路径加载失败</span>
        {onRetry && (
          <button type="button" className="breadcrumb-retry" onClick={onRetry}>
            重试
          </button>
        )}
      </div>
    );
  }

  if (status === 'empty' || items.length === 0) {
    return <div className="breadcrumb-status">无路径</div>;
  }

  return (
    <nav className="breadcrumb" aria-label="面包屑路径">
      {items.map((item, index) => {
        const isLast = index === items.length - 1;
        return (
          <React.Fragment key={item.id}>
            <span className={`breadcrumb-item ${isLast ? 'is-current' : ''}`}>{item.title}</span>
            {!isLast && <span className="breadcrumb-sep" aria-hidden="true">›</span>}
          </React.Fragment>
        );
      })}
    </nav>
  );
};

export default AncestryBreadcrumb;

