# OKR Python 客户端

OKR 管理系统的 Python 命令行客户端，提供完整的任务和日志管理功能，包括计划视图、分数趋势分析和便捷任务创建等扩展功能。

## 版本更新

### v1.1.0 新功能
- 🆕 **计划视图**: 查看指定时间周期的任务综合视图和统计信息
- 📈 **分数趋势**: 获取任务分数和数量的时间序列分析  
- ⚡ **便捷任务创建**: 快速创建基于时间周期的任务（今日、本周、本月、本季度、本年）
- 🚀 **快捷视图命令**: 一键查看今日/本周/本月/本季度/本年的计划视图和分数趋势
- 🎯 **智能时间计算**: 自动计算各种时间周期的开始和结束时间
- 🎨 **增强的CLI体验**: 丰富的颜色输出、图标和树形结构显示

## 安装

1. 克隆项目并进入目录：
```bash
cd okr-python-client
```

2. 创建并激活虚拟环境：
```bash
# 创建虚拟环境
python3 -m venv venv

# 激活虚拟环境
source venv/bin/activate
```

3. 安装依赖：
```bash
pip install -r requirements.txt
```

4. 设置环境变量（可选）：
```bash
cp .env.example .env
# 编辑 .env 文件中的配置
```

## 使用方法

### 认证

首先需要登录：
```bash
python -m okr_client.cli login
```

查看当前用户信息：
```bash
python -m okr_client.cli me
```

登出：
```bash
python -m okr_client.cli logout
```

### 🆕 计划视图

查看计划视图可以获得指定时间周期内的任务综合视图，包括任务树结构、统计信息和相关日志。

#### ⚡ 快捷计划视图命令
```bash
# 查看今日计划
python -m okr_client.cli plan today

# 查看本周计划
python -m okr_client.cli plan week

# 查看本月计划
python -m okr_client.cli plan month

# 查看本季度计划
python -m okr_client.cli plan quarter

# 查看本年计划
python -m okr_client.cli plan year
```

#### 基本计划视图命令
```bash
# 查看2024年第4季度计划
python -m okr_client.cli plan view --scale quarter --time-ref 2024-Q4

# 查看2025年7月计划
python -m okr_client.cli plan view --scale month --time-ref 2025-07

# 查看2025年第15周计划
python -m okr_client.cli plan view --scale week --time-ref 2025-W15
```

#### 指定时间计划视图命令
```bash
# 查看指定季度计划
python -m okr_client.cli plan quarterly 2024 4

# 查看指定月份计划
python -m okr_client.cli plan monthly 2025 7
```

### 📈 分数趋势分析

分数趋势功能可以分析指定时间周期内的任务分数变化趋势和统计摘要。

#### ⚡ 快捷分数趋势命令
```bash
# 查看今日分数趋势
python -m okr_client.cli stats today

# 查看本周分数趋势
python -m okr_client.cli stats week

# 查看本月分数趋势
python -m okr_client.cli stats month

# 查看本季度分数趋势
python -m okr_client.cli stats quarter

# 查看本年分数趋势
python -m okr_client.cli stats year
```

#### 基本趋势命令
```bash
# 查看2025年7月的分数趋势
python -m okr_client.cli stats trend --scale month --time-ref 2025-07

# 查看2024年第4季度的分数趋势
python -m okr_client.cli stats trend --scale quarter --time-ref 2024-Q4
```

#### 指定时间趋势命令
```bash
# 查看月度分数趋势
python -m okr_client.cli stats monthly-trend 2025 7

# 查看季度分数趋势
python -m okr_client.cli stats quarterly-trend 2024 4
```

### ⚡ 便捷任务创建

新版本提供了多种便捷的任务创建方式，自动计算时间范围，大大简化了任务创建过程。

#### 基于当前时间的快速创建
```bash
# 创建今日任务
python -m okr_client.cli task today "完成代码审查"

# 创建本周任务  
python -m okr_client.cli task week "完成项目架构设计"

# 创建本月任务
python -m okr_client.cli task month "学习Go语言" --desc "深入学习Go并完成一个项目" --score 8

# 创建本季度任务
python -m okr_client.cli task quarter "提升编程技能" --score 9

# 创建本年任务
python -m okr_client.cli task year "成为全栈工程师"
```

#### 基于指定时间的创建
```bash
# 创建指定季度任务
python -m okr_client.cli task quarter "Q4目标" --year 2024 --q 4

# 创建指定月份任务  
python -m okr_client.cli task month "7月计划" --year 2025 --month 7

# 创建指定周任务
python -m okr_client.cli task week "第15周计划" --year 2025 --week 15
```

#### 扩展的create命令快捷选项
```bash
# 使用快捷选项创建本月任务
python -m okr_client.cli task create "学习计划" --quick-month

# 使用快捷选项创建本年任务
python -m okr_client.cli task create "年度目标" --quick-year

# 使用快捷选项创建本季度任务
python -m okr_client.cli task create "季度OKR" --quick-quarter
```

### 任务管理

查看任务列表：
```bash
# 查看所有任务
python -m okr_client.cli task list

# 按类型筛选
python -m okr_client.cli task list --type day

# 按日期筛选
python -m okr_client.cli task list --date 2025-07-11

# 按状态筛选
python -m okr_client.cli task list --status completed
```

创建任务：
```bash
# 创建简单任务
python -m okr_client.cli task create --title "完成项目文档"

# 创建详细任务
python -m okr_client.cli task create \
  --title "完成项目文档" \
  --desc "编写用户手册和API文档" \
  --type week \
  --score 8
```

更新任务：
```bash
# 更新任务状态
python -m okr_client.cli task update TASK_ID --status in-progress

# 更新任务分数
python -m okr_client.cli task update TASK_ID --score 9
```

完成任务：
```bash
python -m okr_client.cli task done TASK_ID
```

### 日志管理

查看日志列表：
```bash
# 查看所有日志
python -m okr_client.cli journal list

# 按时间尺度筛选
python -m okr_client.cli journal list --scale day

# 按日期筛选
python -m okr_client.cli journal list --date 2025-07-11
```

创建日志：
```bash
# 创建今日日志
python -m okr_client.cli journal create --content "今天完成了客户端开发"

# 创建周日志
python -m okr_client.cli journal create \
  --content "本周完成了后端API和Python客户端" \
  --scale week \
  --type summary
```

编辑日志：
```bash
python -m okr_client.cli journal edit JOURNAL_ID --content "更新的日志内容"
```

删除日志：
```bash
python -m okr_client.cli journal delete JOURNAL_ID
```

## 程序化API使用

除了CLI工具，也可以在Python代码中直接使用客户端：

```python
from okr_client import OKRClient, TimeScale

# 创建客户端
client = OKRClient()

# 登录
client.login("username", "password")

# 便捷任务创建
task = client.create_today_task("完成用户故事", "实现用户注册功能", 7)
month_task = client.create_this_month_task("学习新技术")
quarter_task = client.create_quarter_task("Q4目标", 2024, 4, "完成年度OKR", 9)

# 获取计划视图
plan = client.get_plan_view(TimeScale.QUARTER, "2024-Q4")
month_plan = client.get_plan_view_for_month(2025, 7)

# 获取分数趋势
trend = client.get_score_trend(TimeScale.MONTH, "2025-07")
quarterly_trend = client.get_quarterly_score_trend(2024, 4)
```

## 时间格式说明

支持多种时间格式：
- **年格式**: `2024`
- **季度格式**: `2024-Q4` (第4季度)
- **月格式**: `2025-07` (7月)
- **周格式**: `2025-W15` (第15周，ISO周标准)
- **日格式**: `2025-07-11` (ISO日期格式)

## 配置

客户端会在 `~/.okr/config` 保存认证信息。

环境变量：
- `OKR_API_BASE_URL`: API 服务器地址 (默认: http://localhost:8081/api)
- `OKR_CONFIG_PATH`: 配置文件路径 (默认: ~/.okr/config)

## 错误处理

如果遇到认证错误，请重新登录：
```bash
python -m okr_client.cli login
```

如果 API 服务器无法访问，请检查：
1. 服务器是否运行
2. `OKR_API_BASE_URL` 是否正确
3. 网络连接是否正常

## 开发

运行测试：
```bash
python -m pytest tests/
```

运行特定测试：
```bash
# 测试时间工具函数
python -m pytest tests/test_utils.py

# 测试客户端扩展功能
python -m pytest tests/test_client_extensions.py
```

## 更新日志

### v1.1.0 (2025-07-11)
- ✨ 新增计划视图功能，支持任务综合视图展示
- ✨ 新增分数趋势分析功能，支持时间序列统计
- ✨ 新增便捷任务创建功能，支持多种时间周期
- 🚀 新增快捷视图命令：plan/stats today/week/month/quarter/year
- ✨ 新增丰富的CLI命令和便捷选项
- 🎨 改进输出格式，支持树形结构、颜色和图标
- 🧪 添加全面的单元测试和集成测试
- 📚 更新文档和使用示例

### v1.0.0
- 🎉 初始版本
- ✅ 基础任务管理功能
- ✅ 日志管理功能
- ✅ 用户认证功能

## 许可证

MIT License
