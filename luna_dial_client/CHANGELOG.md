# 更新日志

所有重要的项目更改都将记录在此文件中。

该项目遵循[语义化版本控制](https://semver.org/lang/zh-CN/)规范。

## [1.1.0] - 2025-07-11

### 新增功能
- ✨ **计划视图功能**: 新增 `/plan` API 支持，可获取指定时间周期的任务综合视图
  - 支持任务树形结构显示
  - 提供统计信息汇总（总任务数、完成数、进行中、待开始等）
  - 显示相关时间周期的日志条目
  - 支持年/季/月/周/日等多种时间尺度

- 📈 **分数趋势分析**: 新增 `/stats/score-trend-ref` API 支持
  - 按指定时间粒度聚合任务分数和数量
  - 返回时间序列数据和统计摘要
  - 支持图形化趋势显示

- ⚡ **便捷任务创建**: 大幅简化任务创建流程
  - 基于当前时间的快速创建：`create_today_task()`, `create_this_week_task()` 等
  - 基于指定时间的创建：`create_quarter_task()`, `create_month_task()` 等
  - 自动计算合适的时间范围
  - 智能参数验证

- 🎯 **时间计算工具**: 新增 `utils.py` 模块
  - 支持多种时间范围计算：今日、本周、本月、本季度、本年
  - 支持指定时间周期计算：指定季度、月份、周
  - 完整的时间格式验证
  - 时间参考字符串解析

### CLI 命令扩展
- 🆕 **计划视图命令组**: `okr plan`
  - `okr plan view --scale --time-ref`: 查看计划视图
  - `okr plan quarterly <year> <quarter>`: 查看季度计划
  - `okr plan monthly <year> <month>`: 查看月度计划

- 📊 **统计命令组扩展**: `okr stats`
  - `okr stats trend --scale --time-ref`: 查看分数趋势
  - `okr stats monthly-trend <year> <month>`: 月度趋势
  - `okr stats quarterly-trend <year> <quarter>`: 季度趋势

- ⚡ **便捷任务创建命令**:
  - `okr task today <title>`: 创建今日任务
  - `okr task week <title>`: 创建本周任务
  - `okr task month <title>`: 创建本月任务
  - `okr task quarter <title>`: 创建本季度任务
  - `okr task year <title>`: 创建本年任务
  - 支持 `--year`, `--q`, `--month`, `--week` 参数指定时间
  - 为 `task create` 添加快捷选项：`--quick-month`, `--quick-year` 等

### 数据模型扩展
- 📦 **新增模型类**:
  - `TimeRange`: 时间范围模型
  - `TaskTree`: 任务树模型（支持递归结构）
  - `PlanRequest`/`PlanResponse`: 计划视图请求和响应
  - `ScoreTrendRequest`/`ScoreTrendResponse`: 分数趋势请求和响应
  - `PlanStats`: 计划统计信息
  - `TrendSummary`: 趋势摘要统计

### 用户体验改进
- 🎨 **增强的输出格式**:
  - 任务树形结构显示（使用 Rich Tree）
  - 彩色状态指示器和图标
  - 分数趋势图形化显示
  - 统计信息表格化展示
  - 进度百分比和可视化

- ⌚ **智能时间显示**:
  - 自动计算并显示时间范围
  - 支持多种时间格式解析
  - 用户友好的时间提示

### 错误处理和测试
- 🛡️ **新增异常类型**:
  - `PlanViewError`: 计划视图相关错误
  - `ScoreTrendError`: 分数趋势相关错误
  - `TaskCreationError`: 任务创建相关错误

- 🧪 **全面的测试覆盖**:
  - `tests/test_utils.py`: 时间计算函数测试
  - `tests/test_client_extensions.py`: 客户端扩展功能测试
  - 时间计算测试覆盖率 > 90%
  - 边界条件和异常情况测试

### 文档更新
- 📚 **完善的文档**:
  - 更新 README.md，新增扩展功能说明
  - 添加详细的使用示例
  - 更新 API 使用指南
  - 添加版本更新日志

### 技术改进
- 🔧 **代码质量提升**:
  - 遵循 PEP8 代码规范
  - 完整的类型注解
  - 丰富的文档字符串
  - 模块化设计

- 📦 **包结构优化**:
  - 新增 `utils.py` 工具模块
  - 更新 `__init__.py` 导出配置
  - 版本号更新为 1.1.0

## [1.0.0] - 2025-07-10

### 初始版本
- 🎉 **基础功能实现**:
  - 用户认证（登录、登出、用户信息）
  - 任务管理（创建、查询、更新、删除、完成）
  - 日志管理（创建、查询、更新、删除）
  - CLI 命令行接口

- 🔧 **核心特性**:
  - RESTful API 客户端
  - 命令行工具 `okr_client.cli`
  - 数据模型定义
  - 错误处理机制
  - 配置管理

- 📱 **CLI 命令**:
  - `okr login/logout/me`: 用户认证
  - `okr task list/create/update/done`: 任务管理
  - `okr journal list/create/edit/delete`: 日志管理

### 已知问题
- 无

---

## 更新说明

- **[新增]**: 新功能
- **[改进]**: 对现有功能的改进
- **[修复]**: 错误修复
- **[移除]**: 移除的功能
- **[安全]**: 安全相关更新

## 升级指南

### 从 1.0.0 到 1.1.0

1. **依赖更新**: 无需额外依赖
2. **API 兼容性**: 完全向后兼容
3. **新功能**: 可选使用新的计划视图和便捷创建功能
4. **配置**: 无需修改现有配置

### 迁移建议

- 现有的 CLI 命令继续工作
- 推荐使用新的便捷创建命令提高效率
- 可以尝试使用计划视图功能获得更好的任务总览
