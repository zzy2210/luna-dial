"""命令行接口"""

import os
import sys
from datetime import datetime, date, timedelta
from typing import Optional

import click
from rich.console import Console
from rich.table import Table
from rich.prompt import Prompt, Confirm
from rich.tree import Tree
from rich.panel import Panel
from dateutil.parser import parse as parse_date

from .client import OKRClient, OKRClientError, PlanViewError, ScoreTrendError, TaskCreationError
from .models import TaskType, TaskStatus, TimeScale, EntryType, TaskRequest, JournalRequest, TaskTree, PlanResponse, ScoreTrendResponse

console = Console()


def get_client() -> OKRClient:
    """获取 API 客户端实例"""
    return OKRClient()


def handle_error(func):
    """错误处理装饰器"""
    import functools
    
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        try:
            return func(*args, **kwargs)
        except (OKRClientError, PlanViewError, ScoreTrendError, TaskCreationError) as e:
            console.print(f"[red]错误: {e}[/red]")
            sys.exit(1)
        except Exception as e:
            console.print(f"[red]未知错误: {e}[/red]")
            sys.exit(1)
    return wrapper


def parse_date_input(date_str: str) -> datetime:
    """解析日期输入"""
    if not date_str:
        return datetime.now()
    
    try:
        # 尝试解析各种日期格式
        return parse_date(date_str)
    except Exception:
        console.print(f"[red]无效的日期格式: {date_str}[/red]")
        sys.exit(1)


def display_task_tree(task_tree: TaskTree, tree: Tree = None, is_root: bool = True) -> Tree:
    """显示任务树结构"""
    task = task_tree  # 兼容扁平结构
    
    # 状态图标
    status_icons = {
        TaskStatus.PENDING: "⏳",
        TaskStatus.IN_PROGRESS: "🔄", 
        TaskStatus.COMPLETED: "✅"
    }
    
    # 任务标题
    task_label = f"{status_icons.get(task.status, '❓')} {task.title}"
    
    if is_root:
        tree = Tree(task_label)
        current_node = tree
    else:
        current_node = tree.add(task_label)
    
    # 递归添加子任务
    for child in task_tree.children:
        display_task_tree(child, current_node, False)
    
    return tree


def display_plan_view(plan: PlanResponse):
    """显示计划视图"""
    time_range = plan.time_range
    stats = plan.stats
    
    # 标题
    title = f"📋 {time_range.start.strftime('%Y-%m-%d')} ~ {time_range.end.strftime('%Y-%m-%d')} 计划视图"
    console.print(f"\n[bold cyan]{title}[/bold cyan]")
    console.print("━" * len(title))
    
    # 统计概览
    console.print(f"\n📊 [bold]统计概览:[/bold]")
    console.print(f"• 总任务数: [cyan]{stats.total_tasks}[/cyan]")
    if stats.total_tasks > 0:
        completed_pct = (stats.completed_tasks / stats.total_tasks) * 100
        in_progress_pct = (stats.in_progress_tasks / stats.total_tasks) * 100
        pending_pct = (stats.pending_tasks / stats.total_tasks) * 100
        
        console.print(f"• 已完成: [green]{stats.completed_tasks}[/green] ({completed_pct:.1f}%)")
        console.print(f"• 进行中: [blue]{stats.in_progress_tasks}[/blue] ({in_progress_pct:.1f}%)")
        console.print(f"• 待开始: [yellow]{stats.pending_tasks}[/yellow] ({pending_pct:.1f}%)")
    
    console.print(f"• 总分: [magenta]{stats.total_score}[/magenta] / 完成分数: [green]{stats.completed_score}[/green]")
    
    # 任务树
    if plan.tasks:
        console.print(f"\n🌳 [bold]任务树:[/bold]")
        for task_tree in plan.tasks:
            tree = display_task_tree(task_tree)
            console.print(tree)
    else:
        console.print("\n[yellow]该时间段没有任务[/yellow]")
    
    # 相关日志
    if plan.journals:
        console.print(f"\n📝 [bold]相关日志 ({len(plan.journals)}条):[/bold]")
        for journal in plan.journals[:5]:  # 只显示前5条
            content_preview = journal.content[:50] + "..." if len(journal.content) > 50 else journal.content
            console.print(f"• {journal.created_at.strftime('%Y-%m-%d')}: {content_preview}")
        
        if len(plan.journals) > 5:
            console.print(f"... 还有 {len(plan.journals) - 5} 条日志")
    else:
        console.print("\n[yellow]该时间段没有相关日志[/yellow]")


def display_score_trend(trend: ScoreTrendResponse):
    """显示分数趋势"""
    time_range = trend.time_range
    summary = trend.summary
    # 标题
    title = f"📈 {time_range.start.strftime('%Y-%m-%d')} ~ {time_range.end.strftime('%Y-%m-%d')} 分数趋势"
    console.print(f"\n[bold cyan]{title}[/bold cyan]")
    console.print("━" * len(title))
    # 趋势摘要
    console.print(f"\n📊 [bold]趋势摘要:[/bold]")
    if summary is not None:
        console.print(f"• 总分: [magenta]{summary.total_score}[/magenta]")
        console.print(f"• 总任务: [cyan]{summary.total_tasks}[/cyan]")
        console.print(f"• 平均分: [blue]{summary.average_score:.2f}[/blue]")
        console.print(f"• 平均任务数: [blue]{summary.average_task_count:.2f}[/blue]")
        console.print(f"• 最高分: [green]{summary.max_score}[/green]")
        console.print(f"• 最低分: [red]{summary.min_score}[/red]")
    else:
        console.print("[yellow]无趋势摘要数据[/yellow]")
    # 趋势图
    if trend.labels and trend.scores:
        console.print(f"\n📈 [bold]趋势图:[/bold]")
        max_score = max(trend.scores) if trend.scores else 1
        for i, (label, score, count) in enumerate(zip(trend.labels, trend.scores, trend.counts)):
            bar_length = int((score / max_score) * 20) if max_score > 0 else 0
            bar = "▓" * bar_length + "░" * (20 - bar_length)
            console.print(f"{label} {bar} {score}分 ({count}任务)")


@click.group()
def cli():
    """OKR 管理系统命令行工具"""
    pass


# 认证相关命令
@cli.command()
@handle_error
def login():
    """用户登录"""
    username = Prompt.ask("用户名")
    password = Prompt.ask("密码", password=True)
    # 如果用户名或密码为空，设置为默认值
    if not username:
        username = "admin"
    if not password:
        password = "your-password-word"
    
    client = get_client()
    auth_response = client.login(username, password)
    
    console.print(f"[green]登录成功！欢迎，{auth_response.user.username}[/green]")


@cli.command()
@handle_error
def logout():
    """用户登出"""
    client = get_client()
    client.logout()
    console.print("[green]已登出[/green]")


@cli.command()
@handle_error
def me():
    """显示当前用户信息"""
    client = get_client()
    user = client.get_current_user()
    
    table = Table(title="用户信息")
    table.add_column("字段", style="cyan")
    table.add_column("值", style="magenta")
    
    table.add_row("ID", user.id)
    table.add_row("用户名", user.username)
    table.add_row("邮箱", user.email)
    table.add_row("创建时间", user.created_at.strftime("%Y-%m-%d %H:%M:%S"))
    
    console.print(table)


# 计划视图命令组
@cli.group()
def plan():
    """计划视图管理"""
    pass


@plan.command("view")
@click.option("--scale", type=click.Choice([s.value for s in TimeScale]), required=True, help="时间尺度")
@click.option("--time-ref", required=True, help="时间参考")
@handle_error
def plan_view(scale: str, time_ref: str):
    """查看计划视图"""
    client = get_client()
    plan_response = client.get_plan_view(TimeScale(scale), time_ref)
    display_plan_view(plan_response)


# 新增：快捷计划视图命令
@plan.command("today")
@handle_error
def plan_today():
    """查看今日计划"""
    client = get_client()
    today = datetime.now().strftime('%Y-%m-%d')
    plan_response = client.get_plan_view(TimeScale.DAY, today)
    display_plan_view(plan_response)


@plan.command("week")
@handle_error
def plan_this_week():
    """查看本周计划（ISO周编号）"""
    client = get_client()
    from datetime import datetime
    now = datetime.now()
    year, week, _ = now.isocalendar()
    time_ref = f"{year}-W{week:02d}"
    plan_response = client.get_plan_view(TimeScale.WEEK, time_ref)
    display_plan_view(plan_response)


@plan.command("month")
@handle_error
def plan_this_month():
    """查看本月计划"""
    client = get_client()
    this_month = datetime.now().strftime('%Y-%m')
    plan_response = client.get_plan_view(TimeScale.MONTH, this_month)
    display_plan_view(plan_response)


@plan.command("quarter")
@handle_error
def plan_this_quarter():
    """查看本季度计划"""
    client = get_client()
    now = datetime.now()
    quarter = (now.month - 1) // 3 + 1
    time_ref = f"{now.year}-Q{quarter}"
    plan_response = client.get_plan_view(TimeScale.QUARTER, time_ref)
    display_plan_view(plan_response)


@plan.command("year")
@handle_error
def plan_this_year():
    """查看本年计划"""
    client = get_client()
    this_year = str(datetime.now().year)
    plan_response = client.get_plan_view(TimeScale.YEAR, this_year)
    display_plan_view(plan_response)


@plan.command("quarterly")
@click.argument("year", type=int)
@click.argument("quarter", type=int)
@handle_error
def plan_quarterly(year: int, quarter: int):
    """查看指定季度计划（便捷命令）"""
    client = get_client()
    plan_response = client.get_plan_view_for_quarter(year, quarter)
    display_plan_view(plan_response)


@plan.command("monthly")
@click.argument("year", type=int)
@click.argument("month", type=int)
@handle_error
def plan_monthly(year: int, month: int):
    """查看指定月份计划（便捷命令）"""
    client = get_client()
    plan_response = client.get_plan_view_for_month(year, month)
    display_plan_view(plan_response)


# 统计命令组
@cli.group()
def stats():
    """统计分析"""
    pass


@stats.command("trend")
@click.option("--scale", type=click.Choice([s.value for s in TimeScale]), required=True, help="统计尺度")
@click.option("--time-ref", required=True, help="时间参考")
@handle_error
def stats_trend(scale: str, time_ref: str):
    """查看分数趋势"""
    client = get_client()
    trend_response = client.get_score_trend(TimeScale(scale), time_ref)
    display_score_trend(trend_response)


# 新增：快捷分数趋势命令
@stats.command("today")
@handle_error
def stats_today():
    """查看今日分数趋势"""
    client = get_client()
    today = datetime.now().strftime('%Y-%m-%d')
    trend_response = client.get_score_trend(TimeScale.DAY, today)
    display_score_trend(trend_response)


@stats.command("week")
@handle_error
def stats_this_week():
    """查看本周分数趋势"""
    client = get_client()
    # 计算当前周数（ISO周）
    now = datetime.now()
    year, week, _ = now.isocalendar()
    time_ref = f"{year}-W{week:02d}"
    trend_response = client.get_score_trend(TimeScale.WEEK, time_ref)
    display_score_trend(trend_response)


@stats.command("month")
@handle_error
def stats_this_month():
    """查看本月分数趋势"""
    client = get_client()
    this_month = datetime.now().strftime('%Y-%m')
    trend_response = client.get_score_trend(TimeScale.MONTH, this_month)
    display_score_trend(trend_response)


@stats.command("quarter")
@handle_error
def stats_this_quarter():
    """查看本季度分数趋势"""
    client = get_client()
    now = datetime.now()
    quarter = (now.month - 1) // 3 + 1
    time_ref = f"{now.year}-Q{quarter}"
    trend_response = client.get_score_trend(TimeScale.QUARTER, time_ref)
    display_score_trend(trend_response)


@stats.command("year")
@handle_error
def stats_this_year():
    """查看本年分数趋势"""
    client = get_client()
    this_year = str(datetime.now().year)
    trend_response = client.get_score_trend(TimeScale.YEAR, this_year)
    display_score_trend(trend_response)


@stats.command("monthly-trend")
@click.argument("year", type=int)
@click.argument("month", type=int)
@handle_error
def stats_monthly_trend(year: int, month: int):
    """查看月度分数趋势（便捷命令）"""
    client = get_client()
    trend_response = client.get_monthly_score_trend(year, month)
    display_score_trend(trend_response)


@stats.command("quarterly-trend")
@click.argument("year", type=int)
@click.argument("quarter", type=int)
@handle_error
def stats_quarterly_trend(year: int, quarter: int):
    """查看季度分数趋势（便捷命令）"""
    client = get_client()
    trend_response = client.get_quarterly_score_trend(year, quarter)
    display_score_trend(trend_response)


# 任务相关命令组
@cli.group()
def task():
    """任务管理"""
    pass


@task.command("list")
@click.option("--type", "task_type", type=click.Choice([t.value for t in TaskType]), help="任务类型")
@click.option("--date", help="日期 (YYYY-MM-DD 或其他格式)")
@click.option("--status", type=click.Choice([s.value for s in TaskStatus]), help="任务状态")
@handle_error
def list_tasks(task_type: Optional[str], date: Optional[str], status: Optional[str]):
    """查看任务列表"""
    client = get_client()
    
    # 构建查询参数
    kwargs = {}
    if task_type:
        kwargs["task_type"] = TaskType(task_type)
    if status:
        kwargs["status"] = TaskStatus(status)
    
    # 如果指定了日期，设置日期范围
    if date:
        target_date = parse_date_input(date).date()
        kwargs["start_date"] = datetime.combine(target_date, datetime.min.time())
        kwargs["end_date"] = datetime.combine(target_date, datetime.max.time())
    
    tasks = client.get_tasks(**kwargs)
    
    if not tasks:
        console.print("[yellow]没有找到任务[/yellow]")
        return
    
    table = Table(title="任务列表")
    table.add_column("ID", style="dim")
    table.add_column("标题", style="cyan")
    table.add_column("类型", style="green")
    table.add_column("状态", style="yellow")
    table.add_column("分数", style="magenta")
    table.add_column("开始时间", style="blue")
    table.add_column("结束时间", style="blue")
    
    for task in tasks:
        status_color = {
            "pending": "yellow",
            "in-progress": "blue",
            "completed": "green"
        }.get(task.status.value, "white")
        
        table.add_row(
            task.id[:8] + "...",
            task.title,
            task.type.value,
            f"[{status_color}]{task.status.value}[/{status_color}]",
            str(task.score) if task.score else "-",
            task.start_date.strftime("%m-%d"),
            task.end_date.strftime("%m-%d")
        )
    
    console.print(table)


@task.command("create")
@click.option("--title", required=True, help="任务标题")
@click.option("--desc", help="任务描述")
@click.option("--type", "task_type", type=click.Choice([t.value for t in TaskType]), default="day", help="任务类型")
@click.option("--start-date", help="开始日期")
@click.option("--end-date", help="结束日期")
@click.option("--score", type=int, help="分数 (1-10)")
@click.option("--quick-month", is_flag=True, help="自动设置为本月任务")
@click.option("--quick-year", is_flag=True, help="自动设置为本年任务")
@click.option("--quick-quarter", is_flag=True, help="自动设置为本季度任务") 
@click.option("--quick-week", is_flag=True, help="自动设置为本周任务")
@handle_error
def create_task(title: str, desc: Optional[str], task_type: str, start_date: Optional[str], 
              end_date: Optional[str], score: Optional[int], quick_month: bool, quick_year: bool,
              quick_quarter: bool, quick_week: bool):
    """创建任务"""
    client = get_client()
    
    # 检查快速创建选项
    if quick_month:
        task = client.create_this_month_task(title, desc, score)
        console.print(f"[green]本月任务创建成功！ID: {task.id}[/green]")
        console.print(f"时间范围: {task.start_date.strftime('%Y-%m-%d')} ~ {task.end_date.strftime('%Y-%m-%d')}")
        return
    elif quick_year:
        task = client.create_this_year_task(title, desc, score)
        console.print(f"[green]本年任务创建成功！ID: {task.id}[/green]")
        console.print(f"时间范围: {task.start_date.strftime('%Y-%m-%d')} ~ {task.end_date.strftime('%Y-%m-%d')}")
        return
    elif quick_quarter:
        task = client.create_this_quarter_task(title, desc, score)
        console.print(f"[green]本季度任务创建成功！ID: {task.id}[/green]")
        console.print(f"时间范围: {task.start_date.strftime('%Y-%m-%d')} ~ {task.end_date.strftime('%Y-%m-%d')}")
        return
    elif quick_week:
        task = client.create_this_week_task(title, desc, score)
        console.print(f"[green]本周任务创建成功！ID: {task.id}[/green]")
        console.print(f"时间范围: {task.start_date.strftime('%Y-%m-%d')} ~ {task.end_date.strftime('%Y-%m-%d')}")
        return
    
    # 常规创建流程
    print("[调试] CLI create_task 被调用")
    print("[调试] get_client 返回:", client)
    # 解析日期
    if start_date:
        start_dt = parse_date_input(start_date)
    else:
        start_dt = datetime.now()
    
    if end_date:
        end_dt = parse_date_input(end_date)
    else:
        # 根据任务类型设置默认结束时间
        if task_type == "day":
            end_dt = start_dt.replace(hour=23, minute=59, second=59)
        elif task_type == "week":
            end_dt = start_dt.replace(hour=23, minute=59, second=59) + timedelta(days=7)
        else:
            end_dt = start_dt.replace(hour=23, minute=59, second=59)
    
    # 转为 Go 端能识别的 RFC3339 格式（不带微秒，带 Z）
    def to_rfc3339(dt):
        return dt.replace(microsecond=0).isoformat() + 'Z'

    # 只在 score 不为 None 时传递，否则为 None
    task_request = TaskRequest(
        title=title,
        description=desc,
        type=TaskType(task_type),
        start_date=to_rfc3339(start_dt),
        end_date=to_rfc3339(end_dt),
        score=score if score is not None else None,
        status=TaskStatus.PENDING
    )
    print("[调试] CLI 发送的 task_request:", task_request.model_dump(mode="json", exclude_none=True))
    task = client.create_task(task_request)
    console.print(f"[green]任务创建成功！ID: {task.id}[/green]")


# 便捷任务创建命令
@task.command("today")
@click.argument("title")
@click.option("--desc", help="任务描述")
@click.option("--score", type=int, help="分数 (1-10)")
@handle_error
def create_today_task(title: str, desc: Optional[str], score: Optional[int]):
    """创建今日任务"""
    client = get_client()
    task = client.create_today_task(title, desc, score)
    console.print(f"[green]今日任务创建成功！ID: {task.id}[/green]")
    console.print(f"时间范围: {task.start_date.strftime('%Y-%m-%d')} ~ {task.end_date.strftime('%Y-%m-%d')}")


@task.command("week")
@click.argument("title")
@click.option("--desc", help="任务描述")
@click.option("--score", type=int, help="分数 (1-10)")
@handle_error
def create_week_task(title: str, desc: Optional[str], score: Optional[int]):
    """创建本周任务"""
    client = get_client()
    task = client.create_this_week_task(title, desc, score)
    console.print(f"[green]本周任务创建成功！ID: {task.id}[/green]")
    console.print(f"时间范围: {task.start_date.strftime('%Y-%m-%d')} ~ {task.end_date.strftime('%Y-%m-%d')}")


@task.command("month")
@click.argument("title")
@click.option("--desc", help="任务描述")
@click.option("--score", type=int, help="分数 (1-10)")
@handle_error
def create_month_task(title: str, desc: Optional[str], score: Optional[int]):
    """创建本月任务"""
    client = get_client()
    task = client.create_this_month_task(title, desc, score)
    console.print(f"[green]本月任务创建成功！ID: {task.id}[/green]")
    console.print(f"时间范围: {task.start_date.strftime('%Y-%m-%d')} ~ {task.end_date.strftime('%Y-%m-%d')}")


@task.command("quarter")
@click.argument("title")
@click.option("--desc", help="任务描述")
@click.option("--score", type=int, help="分数 (1-10)")
@click.option("--year", type=int, help="指定年份")
@click.option("--q", type=int, help="指定季度 (1-4)")
@handle_error
def create_quarter_task(title: str, desc: Optional[str], score: Optional[int], 
                       year: Optional[int], q: Optional[int]):
    """创建季度任务"""
    client = get_client()
    
    if year and q:
        # 创建指定季度任务
        task = client.create_quarter_task(title, year, q, desc, score)
        console.print(f"[green]{year}年第{q}季度任务创建成功！ID: {task.id}[/green]")
    else:
        # 创建本季度任务
        task = client.create_this_quarter_task(title, desc, score)
        console.print(f"[green]本季度任务创建成功！ID: {task.id}[/green]")
    
    console.print(f"时间范围: {task.start_date.strftime('%Y-%m-%d')} ~ {task.end_date.strftime('%Y-%m-%d')}")


@task.command("year")
@click.argument("title")
@click.option("--desc", help="任务描述")
@click.option("--score", type=int, help="分数 (1-10)")
@handle_error
def create_year_task(title: str, desc: Optional[str], score: Optional[int]):
    """创建本年任务"""
    client = get_client()
    task = client.create_this_year_task(title, desc, score)
    console.print(f"[green]本年任务创建成功！ID: {task.id}[/green]")
    console.print(f"时间范围: {task.start_date.strftime('%Y-%m-%d')} ~ {task.end_date.strftime('%Y-%m-%d')}")


@task.command("update")
@click.argument("task_id")
@click.option("--title", help="任务标题")
@click.option("--desc", help="任务描述")
@click.option("--status", type=click.Choice([s.value for s in TaskStatus]), help="任务状态")
@click.option("--score", type=int, help="分数 (1-10)")
@handle_error
def update_task(task_id: str, title: Optional[str], desc: Optional[str], status: Optional[str], score: Optional[int]):
    """更新任务"""
    client = get_client()
    
    # 获取现有任务
    existing_task = client.get_task(task_id)
    
    # 更新字段
    task_request = TaskRequest(
        title=title or existing_task.title,
        description=desc if desc is not None else existing_task.description,
        type=existing_task.type,
        start_date=existing_task.start_date,
        end_date=existing_task.end_date,
        status=TaskStatus(status) if status else existing_task.status,
        score=score if score is not None else existing_task.score,
        parent_id=existing_task.parent_id,
        tags=existing_task.tags
    )
    
    task = client.update_task(task_id, task_request)
    console.print(f"[green]任务更新成功！[/green]")


@task.command("done")
@click.argument("task_id")
@handle_error
def complete_task(task_id: str):
    """完成任务"""
    client = get_client()
    task = client.complete_task(task_id)
    console.print(f"[green]任务 '{task.title}' 已完成！[/green]")


# 日志相关命令组
@cli.group()
def journal():
    """日志管理"""
    pass


@journal.command("list")
@click.option("--scale", type=click.Choice([s.value for s in TimeScale]), help="时间尺度")
@click.option("--date", help="日期")
@handle_error
def list_journals(scale: Optional[str], date: Optional[str]):
    """查看日志列表"""
    client = get_client()
    
    kwargs = {}
    if scale:
        kwargs["time_scale"] = TimeScale(scale)
    if date:
        target_date = parse_date_input(date)
        kwargs["start_time"] = target_date
        kwargs["end_time"] = target_date.replace(hour=23, minute=59, second=59)
    
    journals = client.get_journals(**kwargs)
    
    if not journals:
        console.print("[yellow]没有找到日志[/yellow]")
        return
    
    table = Table(title="日志列表")
    table.add_column("ID", style="dim")
    table.add_column("内容", style="cyan")
    table.add_column("时间尺度", style="green")
    table.add_column("类型", style="yellow")
    table.add_column("创建时间", style="blue")
    
    for journal in journals:
        content_preview = journal.content[:50] + "..." if len(journal.content) > 50 else journal.content
        table.add_row(
            journal.id[:8] + "...",
            content_preview,
            journal.time_scale.value,
            journal.entry_type.value,
            journal.created_at.strftime("%m-%d %H:%M")
        )
    
    console.print(table)


@journal.command("create")
@click.option("--content", required=True, help="日志内容")
@click.option("--scale", type=click.Choice([s.value for s in TimeScale]), default="day", help="时间尺度")
@click.option("--type", "entry_type", type=click.Choice([e.value for e in EntryType]), default="reflection", help="日志类型")
@handle_error
def create_journal(content: str, scale: str, entry_type: str):
    """创建日志"""
    client = get_client()
    
    journal_request = JournalRequest(
        content=content,
        time_scale=TimeScale(scale),
        entry_type=EntryType(entry_type),
        time_reference=datetime.now().strftime("%Y-%m-%d")
    )
    
    journal = client.create_journal(journal_request)
    console.print(f"[green]日志创建成功！ID: {journal.id}[/green]")


@journal.command("edit")
@click.argument("journal_id")
@click.option("--content", required=True, help="新的日志内容")
@handle_error
def edit_journal(journal_id: str, content: str):
    """编辑日志"""
    client = get_client()
    
    # 获取现有日志
    existing_journal = client.get_journal(journal_id)
    
    journal_request = JournalRequest(
        content=content,
        time_scale=existing_journal.time_scale,
        entry_type=existing_journal.entry_type,
        time_reference=existing_journal.time_reference
    )
    
    journal = client.update_journal(journal_id, journal_request)
    console.print(f"[green]日志更新成功！[/green]")


@journal.command("delete")
@click.argument("journal_id")
@handle_error
def delete_journal(journal_id: str):
    """删除日志"""
    if not Confirm.ask("确定要删除这个日志吗？"):
        return
    
    client = get_client()
    client.delete_journal(journal_id)
    console.print(f"[green]日志删除成功！[/green]")


def main():
    """主函数，用于打包后的入口点"""
    cli()


if __name__ == "__main__":
    cli()
