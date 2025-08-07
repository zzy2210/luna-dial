# todo-list

不确定做没做
梳理api能力
- [ ] 分页查看任务
- [ ] 查看指定任务的任务树
- [ ] 查看任务的父任务（遍历）
- [ ] 分页查看全部任务（按任务树形式排列）
- [ ] 分页查看日志

## 实现方案：混合方案（数据库冗余字段 + 业务层树结构）

### 数据库层面改动
1. **添加树优化字段**：
   - `has_children BOOLEAN DEFAULT FALSE` - 是否有子任务
   - `children_count INT DEFAULT 0` - 直接子任务数量
   - `root_task_id VARCHAR(36)` - 指向根任务ID
   - `tree_depth INT DEFAULT 0` - 在树中的深度

2. **索引优化**：
   - `idx_tasks_root_task_id` - 根任务ID索引
   - `idx_tasks_no_parent` - 根任务查询索引

### 业务层面改动
1. **Task结构扩展**：
   - `Children []*Task` - 子任务列表（内存构建，不存数据库）
   - 利用新增的数据库字段提高查询效率

### 核心查询策略
1. **获取任意任务的完整子任务树**：
   - 一次查询：`WHERE (id = ? OR root_task_id = ?) AND user_id = ?`
   - 可选状态过滤：`AND status != 'cancelled'` (默认不显示已取消的任务)
   - 内存构建树结构

2. **全局任务树视图（分页）**：
   - 步骤1：分页获取根任务 `WHERE parent_id IS NULL OR parent_id = '' AND status != 'cancelled'`
   - 步骤2：批量获取相关子任务 `WHERE root_task_id IN (...) AND status != 'cancelled'`
   - 步骤3：内存组装完整树

3. **父任务链查询**：
   - 利用现有递归逻辑或通过root_task_id优化
   - 包含所有状态的任务（便于理解完整的层级关系）

**状态过滤策略**：
- **默认视图**：隐藏已取消的任务，保持界面整洁
- **完整视图**：显示所有状态的任务，用于回顾和分析
- **灵活过滤**：支持按状态组合查询（如只看进行中的任务）

### API设计要点
- 支持任意深度的任务树查询
- 全局视图按根任务分页
- 保持树结构的完整性
- 优化长期大量任务的查询性能

#### API响应格式（嵌套树结构）
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": "task1",
      "title": "2024年目标",
      "task_type": 4,
      "tree_depth": 0,
      "has_children": true,
      "children_count": 2,
      "children": [
        {
          "id": "task2",
          "title": "Q1目标",
          "parent_id": "task1",
          "tree_depth": 1,
          "has_children": true,
          "children_count": 1,
          "children": [
            {
              "id": "task3",
              "title": "1月任务",
              "parent_id": "task2",
              "tree_depth": 2,
              "has_children": false,
              "children_count": 0,
              "children": []
            }
          ]
        }
      ]
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 50,
    "total_pages": 5,
    "has_next": true
  }
}
```

#### 新增API接口
1. **移动任务到新父任务**：
   ```
   PUT /api/v1/tasks/:task_id/move
   Body: { "new_parent_id": "parent_task_id" }  // 空字符串表示变成根任务
   ```

2. **获取可选的父任务列表**（避免循环引用）：
   ```
   GET /api/v1/tasks/:task_id/possible-parents
   // 返回可以作为该任务父任务的任务列表（排除自己和自己的子任务）
   ```

### 数据一致性维护

#### 1. 新建任务
```go
func CreateTask(task *Task) error {
    if task.ParentID != "" {
        // 查询父任务信息
        parent := GetTask(task.ParentID)
        task.TreeDepth = parent.TreeDepth + 1
        task.RootTaskID = parent.RootTaskID
        
        // 更新父任务的子任务计数
        UpdateParentChildrenCount(task.ParentID, +1)
    } else {
        // 根任务
        task.TreeDepth = 0
        task.RootTaskID = task.ID  // 根任务的root_task_id指向自己
    }
    
    task.HasChildren = false
    task.ChildrenCount = 0
    return SaveTask(task)
}
```

#### 2. 删除任务处理策略
**采用级联删除 + 状态管理的设计理念**：
- **设计思路**：尽量不要隐藏自己的过往
- **推荐做法**：任务做不完时设置为"进行中"或"取消"状态，而不是删除
- **删除场景**：仅在确实需要物理删除时使用（如测试数据、错误创建等）

```go
func DeleteTask(taskID string) error {
    // 级联删除策略
    // 1. 递归删除所有子任务
    childTasks := GetChildTasks(taskID)
    for _, child := range childTasks {
        DeleteTask(child.ID)  // 递归删除
    }
    
    // 2. 更新父任务的计数
    task := GetTask(taskID)
    if task.ParentID != "" {
        UpdateParentChildrenCount(task.ParentID, -1)
        UpdateParentHasChildren(task.ParentID)  // 重新计算是否有子任务
    }
    
    // 3. 删除当前任务
    return DeleteTaskFromDB(taskID)
}

// 推荐的状态管理方式
func CancelTask(taskID string) error {
    task := GetTask(taskID)
    task.Status = TaskStatusCancelled  // 设置为取消状态
    return UpdateTask(task)
}
```

**设计优势**：
- ✅ **保留历史**：通过状态管理保留任务历史，便于回顾和分析
- ✅ **数据完整性**：避免因删除导致的数据关系破坏
- ✅ **简化逻辑**：减少复杂的孤儿任务处理逻辑

#### 3. 移动任务（改变父任务）
**常见场景：将独立任务加入现有任务树**
例如：一个独立的"月任务"加入到某个"年任务"下

```go
func MoveTask(taskID, newParentID string) error {
    task := GetTask(taskID)
    oldParentID := task.ParentID
    
    // 1. 更新旧父任务的计数（如果有的话）
    if oldParentID != "" {
        UpdateParentChildrenCount(oldParentID, -1)
    }
    
    // 2. 计算新的树信息
    var newTreeDepth int
    var newRootTaskID string
    
    if newParentID != "" {
        // 加入到现有任务树
        newParent := GetTask(newParentID)
        newTreeDepth = newParent.TreeDepth + 1
        newRootTaskID = newParent.RootTaskID
        
        // 更新新父任务的计数
        UpdateParentChildrenCount(newParentID, +1)
    } else {
        // 变成独立的根任务
        newTreeDepth = 0
        newRootTaskID = taskID
    }
    
    // 3. 递归更新整个子树的树信息
    UpdateTaskTreeInfo(taskID, newTreeDepth, newRootTaskID)
    
    return nil
}

func UpdateTaskTreeInfo(taskID string, newDepth int, newRootID string) {
    task := GetTask(taskID)
    oldRootID := task.RootTaskID
    
    // 更新当前任务
    task.TreeDepth = newDepth
    task.RootTaskID = newRootID
    task.ParentID = newParentID  // 更新父任务ID
    UpdateTask(task)
    
    // 递归更新所有子任务
    childTasks := GetChildTasks(taskID)
    for _, child := range childTasks {
        UpdateTaskTreeInfo(child.ID, newDepth+1, newRootID)
    }
    
    // 如果原来是根任务且有子任务，需要特殊处理
    if oldRootID == taskID && len(childTasks) > 0 {
        // 原来这个任务是根任务，现在变成子任务了
        // 所有子任务的root_task_id都需要更新
    }
}
```

**使用场景示例**：
1. **独立任务 → 加入任务树**：
   - 原来：`月任务A` (独立任务，parent_id="", root_task_id=自己的ID)
   - 现在：`年任务B > 月任务A` (parent_id=年任务B的ID, root_task_id=年任务B的ID)

2. **任务在不同树间移动**：
   - 原来：`年任务X > 季任务Y > 月任务A`
   - 现在：`年任务Z > 月任务A` (从一个树移动到另一个树)

#### 4. 错误恢复机制
- 定期检查数据一致性
- 提供数据修复工具
- 在关键操作时验证树结构完整性