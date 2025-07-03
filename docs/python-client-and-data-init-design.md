### **Python 客户端与后端数据初始化 - 设计与执行计划**

本文档旨在为 OKR 管理系统的 Python 客户端提供设计规范，并规划后端添加初始测试数据的执行步骤。

#### **第一部分：后端数据初始化**

**目标**：创建一个可用于登录和测试的 `admin` 用户。

**设计**：

1.  **用户凭证**：
    *   用户名：`admin`
    *   密码：与 `config.ini` 文件中定义的数据库密码一致。
2.  **实现方式**：
    *   为了保证安全和可维护性，我们将通过在 `main.go` 中添加一个一次性的启动任务来创建这个用户。
    *   该任务会在应用启动时检查 `admin` 用户是否存在。如果不存在，它将读取数据库配置，获取密码，对密码进行哈希处理，然后创建用户。
    *   这种方式避免了将明文密码写入代码或迁移脚本，并且只在需要时执行一次。

**执行计划**：

1.  **分析密码哈希逻辑** ✅：检查 `internal/service/user_service.go`，确认现有的用户注册流程是如何对密码进行哈希处理的。
2.  **修改 `main.go`** ✅：
    *   在 `main.go` 的 `main` 函数中，初始化数据库连接之后，添加一个新的函数调用，如 `config.CreateDefaultUserIfNotExist(client, cfg)`。
    *   这个函数将封装创建用户的逻辑。
3.  **实现 `CreateDefaultUserIfNotExist`** ✅：
    *   在 `internal/config/migration.go` 文件中实现该函数。
    *   函数逻辑：
        a.  查询数据库中是否存在用户名为 `admin` 的用户。
        b.  如果存在，则直接返回。
        c.  如果不存在，读取 `config.Database.Password` 作为 `admin` 用户的密码。
        d.  使用与用户注册时相同的哈希算法，对密码进行加密。
        e.  创建一个新的 `ent.User` 实例，并将其保存到数据库。
        f.  打印一条日志，提示 `admin` 用户已成功创建。

---

#### **第二部分：Python 客户端设计**

**目标**：创建一个功能完善、易于使用的 Python 客户端，支持用户进行核心的 OKR 管理操作。客户端将以命令行工具（CLI）的形式提供。

**1. 项目结构**

```
okr-python-client/
├── okr_client/
│   ├── __init__.py
│   ├── client.py       # 核心API客户端，处理HTTP请求和认证
│   ├── models.py       # Pydantic模型，用于数据验证和序列化
│   └── cli.py          # Click命令行接口定义
├── tests/
│   ├── test_client.py
│   └── test_cli.py
├── .env.example        # 环境变量示例文件
├── requirements.txt    # Python依赖
└── README.md           # 项目说明
```

**2. 核心功能与命令设计**

我们将使用 `click` 库来构建一个结构清晰的 CLI。

*   **基础命令**: `okr`

*   **认证**:
    *   `okr login`: 提示输入用户名和密码，成功后将 JWT Token 保存到本地配置文件（例如 `~/.okr/config`）。
    *   `okr me`: 显示当前登录的用户信息。

*   **任务 (Tasks)**:
    *   `okr task list [--type <TYPE>] [--date <DATE>]`: 查看任务。
        *   `--type`: 可选，按类型过滤 (year, quarter, month, week, day)。
        *   `--date`: 可选，查看指定日期/周/月/季/年的计划。例如 `2025-07-11`, `2025-W28`, `2025-07`, `2025-Q3`, `2025`。默认为今天。
    *   `okr task create --title "<TITLE>" [--desc "<DESC>"] [--type <TYPE>] ...`: 创建新任务。
    *   `okr task update <TASK_ID> [--title "<TITLE>"] [--status <STATUS>] ...`: 更新任务。
    *   `okr task done <TASK_ID>`: 将任务状态标记为 `completed`。

*   **日志 (Journals)**:
    *   `okr journal list [--scale <SCALE>] [--date <DATE>]`: 查看日志。
        *   `--scale`: 可选，按时间尺度过滤 (year, quarter, month, week, day)。
        *   `--date`: 可选，指定时间范围。默认为今天。
    *   `okr journal create --content "<CONTENT>"`: 为今天创建一篇日志。
    *   `okr journal edit <JOURNAL_ID> --content "<CONTENT>"`: 编辑日志。
    *   `okr journal delete <JOURNAL_ID>`: 删除日志。

**3. 技术选型**

*   **HTTP 客户端**: `requests` - 简单易用，功能强大。
*   **CLI 框架**: `click` - 声明式，易于构建复杂的命令行工具。
*   **数据模型**: `pydantic` - 用于请求和响应数据的验证和类型提示。
*   **配置管理**: 从环境变量或本地配置文件 (`~/.okr/config`) 读取后端 API 地址和认证 Token。

**4. 执行计划**

1.  **创建项目结构** ✅：按照上述设计创建 `okr-python-client` 目录和文件。
2.  **实现 `client.py`** ✅：
    *   创建 `OKRClient` 类。
    *   实现 `login` 方法，处理认证并保存 token。
    *   为 `Task` 和 `Journal` 资源分别实现 `create`, `get`, `list`, `update`, `delete` 等方法。
    *   所有请求方法都应自动附加认证头。
3.  **实现 `models.py`** ✅：
    *   根据 `ai-design-specification.md` 中的 API 规范，使用 Pydantic 定义 `Task`, `Journal`, `User` 等数据模型。
4.  **实现 `cli.py`** ✅：
    *   使用 `click` 创建命令组和子命令。
    *   将 CLI 命令映射到 `OKRClient` 的相应方法。
    *   处理用户输入和命令行参数。
    *   格式化输出，使其在终端中清晰易读。
5.  **编写 `README.md`** ✅：
    *   提供清晰的安装和使用说明。
    *   列出所有可用的命令及其选项。
6.  **编写测试** ✅：为 `client` 和 `cli` 的核心功能编写单元测试。

---

## 执行状态总结

### ✅ 第一部分：后端数据初始化 - 已完成

- [x] 分析了密码哈希逻辑（Argon2ID）
- [x] 在 `internal/config/migration.go` 中实现了 `CreateDefaultUserIfNotExist` 函数
- [x] 修改了 `cmd/server/main.go` 在数据库迁移后调用创建默认用户
- [x] 添加了必要的导入和依赖
- [x] 实现了密码哈希和 UUID 生成功能

**结果**: 系统启动时会自动检查并创建 admin 用户（用户名: admin，密码: 与数据库密码相同）

### ✅ 第二部分：Python 客户端 - 已完成

- [x] 创建了完整的项目结构
- [x] 实现了所有必需的 Python 文件：
  - `okr_client/client.py` - 核心 API 客户端
  - `okr_client/models.py` - Pydantic 数据模型
  - `okr_client/cli.py` - Click 命令行接口
  - `okr_client/__init__.py` - 模块初始化
- [x] 创建了配置文件和文档：
  - `requirements.txt` - Python 依赖
  - `.env.example` - 环境变量示例
  - `README.md` - 详细使用说明
- [x] 编写了基础测试文件：
  - `tests/test_client.py` - 客户端测试
  - `tests/test_cli.py` - CLI 测试

**功能包括**:
- 用户认证 (`login`, `logout`, `me`)
- 任务管理 (`list`, `create`, `update`, `done`)
- 日志管理 (`list`, `create`, `edit`, `delete`)
- 本地配置和 Token 管理
- 错误处理和用户友好的输出

### 使用说明

1. **启动后端服务**（这将自动创建 admin 用户）：
   ```bash
   cd /home/y1nhui/work/github_owner/okr-web
   go run cmd/server/main.go
   ```

2. **安装并使用 Python 客户端**：
   ```bash
   cd okr-python-client
   pip install -r requirements.txt
   python -m okr_client.cli login
   # 输入用户名: admin
   # 输入密码: your-password-word (数据库密码)
   ```

3. **开始使用**：
   ```bash
   # 查看当前用户
   python -m okr_client.cli me
   
   # 创建今日任务
   python -m okr_client.cli task create --title "学习Python客户端使用"
   
   # 查看任务
   python -m okr_client.cli task list
   
   # 创建今日日志
   python -m okr_client.cli journal create --content "今天完成了客户端开发和测试"
   ```

**任务已全部完成！** 🎉
