use std::{collections::HashSet, io};

use chrono::Utc;
use crossterm::event::{self, Event, KeyCode, KeyEvent, KeyEventKind};
use ratatui::{
    DefaultTerminal, Frame,
    layout::{Constraint, Direction, Layout},
    style::{Color, Style},
    text::{Line, Span, Text},
    widgets::{Block, List, ListItem, Paragraph},
};

use crate::models::{Period, PeriodType, Task, TaskPriority, TaskStatus};
use crate::session::Session;

#[derive(Clone, Copy, Debug, Default)]
pub enum ViewMode {
    #[default]
    Today, // 今日
    ThisWeek,       // 本周
    ThisMonth,      // 本月
    ThisQuarter,    // 本季度
    ThisYear,       // 本年
    ExecutionStats, // 执行情况
    CustomTime,     // 自定义时间
    GlobalTree,     // 全局任务树
}

#[derive(Debug, Default)]
#[allow(dead_code)]
pub enum InputMode {
    #[default]
    Normal,
    EditingTaskName,
}

#[derive(Debug, Default)]
pub struct App {
    pub view_mode: ViewMode,
    pub input_mode: InputMode,
    pub session: Session,

    pub running: bool,

    pub tasks: Vec<Task>,                // 所有任务列表
    pub selected_task_index: usize,      // 当前选中的任务索引
    pub expanded_tasks: HashSet<String>, // 记录展开的任务ID
}

impl App {
    pub fn new() -> Self {
        let now = Utc::now();

        // 构建嵌套的任务树结构，模拟后端返回的数据
        let test_tasks = vec![
            // 根任务1：2025年度规划
            Task {
                id: "1".to_string(),
                title: "2025年度规划".to_string(),
                task_type: PeriodType::Yearly,
                time_period: Period {
                    start: now,
                    end: now + chrono::Duration::days(365),
                },
                status: TaskStatus::InProgress,
                tags: vec!["重要".to_string(), "年度".to_string()],
                icon: "🎯".to_string(),
                score: 0,
                priority: TaskPriority::High,
                parent_id: None,
                user_id: "user1".to_string(),
                created_at: now - chrono::Duration::days(30),
                updated_at: now - chrono::Duration::days(1),

                // 树结构字段
                has_children: true,
                children_count: 2,
                tree_depth: 0,
                root_task_id: None,
                children: vec![
                    // 子任务1.1：Q1季度目标
                    Task {
                        id: "2".to_string(),
                        title: "Q1季度目标".to_string(),
                        task_type: PeriodType::Quarterly,
                        time_period: Period {
                            start: now,
                            end: now + chrono::Duration::days(90),
                        },
                        status: TaskStatus::InProgress,
                        tags: vec!["重要".to_string(), "季度".to_string()],
                        icon: "�".to_string(),
                        score: 0,
                        priority: TaskPriority::High,
                        parent_id: Some("1".to_string()),
                        user_id: "user1".to_string(),
                        created_at: now - chrono::Duration::days(25),
                        updated_at: now - chrono::Duration::days(2),

                        // 树结构字段
                        has_children: true,
                        children_count: 2,
                        tree_depth: 1,
                        root_task_id: Some("1".to_string()),
                        children: vec![
                            // 子任务1.1.1：1月学习计划
                            Task {
                                id: "3".to_string(),
                                title: "1月学习Rust".to_string(),
                                task_type: PeriodType::Monthly,
                                time_period: Period {
                                    start: now,
                                    end: now + chrono::Duration::days(30),
                                },
                                status: TaskStatus::InProgress,
                                tags: vec!["学习".to_string(), "技术".to_string()],
                                icon: "🦀".to_string(),
                                score: 0,
                                priority: TaskPriority::Urgent,
                                parent_id: Some("2".to_string()),
                                user_id: "user1".to_string(),
                                created_at: now - chrono::Duration::days(20),
                                updated_at: now,

                                // 树结构字段
                                has_children: true,
                                children_count: 1,
                                tree_depth: 2,
                                root_task_id: Some("1".to_string()),
                                children: vec![
                                    // 子任务1.1.1.1：每日TUI练习
                                    Task {
                                        id: "4".to_string(),
                                        title: "完成TUI客户端开发".to_string(),
                                        task_type: PeriodType::Daily,
                                        time_period: Period {
                                            start: now,
                                            end: now + chrono::Duration::hours(8),
                                        },
                                        status: TaskStatus::InProgress,
                                        tags: vec!["学习".to_string(), "项目".to_string()],
                                        icon: "�".to_string(),
                                        score: 8,
                                        priority: TaskPriority::High,
                                        parent_id: Some("3".to_string()),
                                        user_id: "user1".to_string(),
                                        created_at: now - chrono::Duration::days(1),
                                        updated_at: now,

                                        // 树结构字段
                                        has_children: false,
                                        children_count: 0,
                                        tree_depth: 3,
                                        root_task_id: Some("1".to_string()),
                                        children: vec![],
                                    },
                                ],
                            },
                            // 子任务1.1.2：2月项目实战
                            Task {
                                id: "5".to_string(),
                                title: "2月项目实战".to_string(),
                                task_type: PeriodType::Monthly,
                                time_period: Period {
                                    start: now + chrono::Duration::days(30),
                                    end: now + chrono::Duration::days(60),
                                },
                                status: TaskStatus::NotStarted,
                                tags: vec!["项目".to_string(), "实战".to_string()],
                                icon: "🚀".to_string(),
                                score: 0,
                                priority: TaskPriority::Medium,
                                parent_id: Some("2".to_string()),
                                user_id: "user1".to_string(),
                                created_at: now - chrono::Duration::days(15),
                                updated_at: now - chrono::Duration::days(5),

                                // 树结构字段
                                has_children: false,
                                children_count: 0,
                                tree_depth: 2,
                                root_task_id: Some("1".to_string()),
                                children: vec![],
                            },
                        ],
                    },
                    // 子任务1.2：Q2季度目标
                    Task {
                        id: "6".to_string(),
                        title: "Q2季度目标".to_string(),
                        task_type: PeriodType::Quarterly,
                        time_period: Period {
                            start: now + chrono::Duration::days(90),
                            end: now + chrono::Duration::days(180),
                        },
                        status: TaskStatus::NotStarted,
                        tags: vec!["重要".to_string(), "季度".to_string()],
                        icon: "�".to_string(),
                        score: 0,
                        priority: TaskPriority::Medium,
                        parent_id: Some("1".to_string()),
                        user_id: "user1".to_string(),
                        created_at: now - chrono::Duration::days(10),
                        updated_at: now - chrono::Duration::days(3),

                        // 树结构字段
                        has_children: false,
                        children_count: 0,
                        tree_depth: 1,
                        root_task_id: Some("1".to_string()),
                        children: vec![],
                    },
                ],
            },
            // 根任务2：健康计划（独立的根任务）
            Task {
                id: "7".to_string(),
                title: "健康管理计划".to_string(),
                task_type: PeriodType::Yearly,
                time_period: Period {
                    start: now,
                    end: now + chrono::Duration::days(365),
                },
                status: TaskStatus::NotStarted,
                tags: vec!["健康".to_string(), "生活".to_string()],
                icon: "💪".to_string(),
                score: 0,
                priority: TaskPriority::Medium,
                parent_id: None,
                user_id: "user1".to_string(),
                created_at: now - chrono::Duration::days(5),
                updated_at: now - chrono::Duration::days(1),

                // 树结构字段
                has_children: false,
                children_count: 0,
                tree_depth: 0,
                root_task_id: None,
                children: vec![],
            },
        ];

        App {
            view_mode: ViewMode::GlobalTree,
            input_mode: InputMode::Normal,
            session: Session::new(),
            running: true,
            tasks: test_tasks,
            selected_task_index: 0,
            expanded_tasks: HashSet::new(), // 初始所有任务都是折叠的
        }
    }
    pub fn run(&mut self, terminal: &mut DefaultTerminal) -> io::Result<()> {
        while self.running {
            terminal.draw(|frame| self.draw(frame))?;
            self.handle_events()?;
        }
        Ok(())
    }

    fn draw(&self, frame: &mut Frame) {
        let chunks = Layout::default()
            .direction(Direction::Vertical)
            .constraints([
                Constraint::Length(1), // 顶部导航栏
                Constraint::Fill(1),
                Constraint::Length(1), // 底部提示栏
            ])
            .split(frame.area());

        let title = match self.view_mode {
            ViewMode::GlobalTree => "全局任务树",
            ViewMode::Today => "今日",
            ViewMode::ThisWeek => "本周",
            ViewMode::ThisMonth => "本月",
            ViewMode::ThisQuarter => "本季度",
            ViewMode::ThisYear => "本年",
            ViewMode::ExecutionStats => "执行情况",
            ViewMode::CustomTime => "自定义时间",
        };
        // 顶部导航栏
        let nav_bar = vec![
            if matches!(self.view_mode, ViewMode::Today) {
                Span::styled("[今日]", Style::default().fg(Color::Yellow))
            } else {
                Span::styled("[今日]", Style::default().fg(Color::White))
            },
            Span::raw(" "),
            if matches!(self.view_mode, ViewMode::ThisWeek) {
                Span::styled("[本周]", Style::default().fg(Color::Yellow))
            } else {
                Span::styled("[本周]", Style::default().fg(Color::White))
            },
            Span::raw(" "),
            if matches!(self.view_mode, ViewMode::ThisMonth) {
                Span::styled("[本月]", Style::default().fg(Color::Yellow))
            } else {
                Span::styled("[本月]", Style::default().fg(Color::White))
            },
            Span::raw(" "),
            if matches!(self.view_mode, ViewMode::ThisQuarter) {
                Span::styled("[本季度]", Style::default().fg(Color::Yellow))
            } else {
                Span::styled("[本季度]", Style::default().fg(Color::White))
            },
            Span::raw(" "),
            if matches!(self.view_mode, ViewMode::ThisYear) {
                Span::styled("[本年]", Style::default().fg(Color::Yellow))
            } else {
                Span::styled("[本年]", Style::default().fg(Color::White))
            },
            Span::raw(" "),
            if matches!(self.view_mode, ViewMode::ExecutionStats) {
                Span::styled("[执行情况]", Style::default().fg(Color::Yellow))
            } else {
                Span::styled("[执行情况]", Style::default().fg(Color::White))
            },
            Span::raw(" "),
            if matches!(self.view_mode, ViewMode::CustomTime) {
                Span::styled("[自定义时间]", Style::default().fg(Color::Yellow))
            } else {
                Span::styled("[自定义时间]", Style::default().fg(Color::White))
            },
            Span::raw(" "),
            if matches!(self.view_mode, ViewMode::GlobalTree) {
                Span::styled("[全局任务]", Style::default().fg(Color::Yellow))
            } else {
                Span::styled("[全局任务]", Style::default().fg(Color::White))
            },
        ];

        let nav_line = Line::from(nav_bar);
        let nav_paragraph = Paragraph::new(nav_line);
        frame.render_widget(nav_paragraph, chunks[0]);

        let main_block = Block::bordered().title(title);
        if matches!(self.view_mode, ViewMode::GlobalTree) {
            // TODO 更完善
            let visible_tasks = self.get_visible_tasks();
            let task_items: Vec<ListItem> = visible_tasks
                .iter()
                .enumerate()
                .map(|(i, (task, depth))| {
                    let indicator = if i == self.selected_task_index {
                        "► "
                    } else {
                        "  "
                    };
                    let expand_icon = if task.has_children {
                        if self.expanded_tasks.contains(&task.id) {
                            "▼ "
                        } else {
                            "▶ "
                        }
                    } else {
                        "  "
                    };
                    let indent = " ".repeat(depth * 2);
                    let content = format!(
                        "{}{}{}{} {}",
                        indicator,
                        indent,
                        expand_icon,
                        task.status.icon(),
                        task.title
                    );
                    ListItem::new(content)
                })
                .collect();

            let task_list = List::new(task_items).block(main_block);
            frame.render_widget(task_list, chunks[1]);
        } else {
            let text = Text::from("当前视图未实现");
            let paragraph = Paragraph::new(text).block(main_block);
            frame.render_widget(paragraph, chunks[1]);
        }

        let help_text = "Tab: 切换视图 | q: 退出";
        let help_paragraph = Paragraph::new(help_text);

        frame.render_widget(help_paragraph, chunks[2]);
    }

    fn handle_events(&mut self) -> io::Result<()> {
        match event::read()? {
            Event::Key(key_event) if key_event.kind == KeyEventKind::Press => {
                self.handle_key_event(key_event)
            }
            _ => {}
        }
        Ok(())
    }

    fn handle_key_event(&mut self, key_event: KeyEvent) {
        match key_event.code {
            KeyCode::Char('q') => {
                self.running = false;
            }
            KeyCode::Tab => {
                self.view_mode = match self.view_mode {
                    ViewMode::Today => ViewMode::ThisWeek,
                    ViewMode::ThisWeek => ViewMode::ThisMonth,
                    ViewMode::ThisMonth => ViewMode::ThisQuarter,
                    ViewMode::ThisQuarter => ViewMode::ThisYear,
                    ViewMode::ThisYear => ViewMode::ExecutionStats,
                    ViewMode::ExecutionStats => ViewMode::CustomTime,
                    ViewMode::CustomTime => ViewMode::GlobalTree,
                    ViewMode::GlobalTree => ViewMode::Today,
                };
            }
            KeyCode::Up => {
                if self.selected_task_index > 0 {
                    self.selected_task_index -= 1;
                }
            }
            KeyCode::Down => {
                let visible_tasks = self.get_visible_tasks();
                if self.selected_task_index < visible_tasks.len() - 1 {
                    self.selected_task_index += 1;
                }
            }
            KeyCode::Char(' ') => {
                // 拿id
                if let Some(current_task) = self.get_current_selected_task() {
                    let task_id = current_task.id.clone();
                    // 切换状态
                    if let Some(task) = self.find_task_mut(&task_id) {
                        task.status = task.status.next();
                    }
                }
            }
            // 收起
            KeyCode::Left => {
                if let Some(current_task) = self.get_current_selected_task() {
                    let task_id = current_task.id.clone();
                    self.expanded_tasks.remove(&task_id);
                }
            }
            // 展开
            KeyCode::Right => {
                if let Some(current_task) = self.get_current_selected_task() {
                    self.expanded_tasks.insert(current_task.id.clone());
                }
            }
            _ => {}
        }
    }

    // 获取所有可见任务  任务-层级
    fn get_visible_tasks(&self) -> Vec<(&Task, usize)> {
        let mut result = Vec::new();
        for task in &self.tasks {
            // 根目录统一给到 add
            self.add_task_if_visible(task, 0, &mut result);
        }
        result
    }

    // 为任务-层级序列添加任务
    fn add_task_if_visible<'a>(
        &self,
        task: &'a Task,
        depth: usize,
        result: &mut Vec<(&'a Task, usize)>,
    ) {
        // 直接添加传入数据
        result.push((task, depth));
        // 如果传入的task 是展开的
        if self.expanded_tasks.contains(&task.id) {
            for child in &task.children {
                // 递归添加子任务
                self.add_task_if_visible(child, depth + 1, result);
            }
        }
    }

    // 获取当前选择的任务
    fn get_current_selected_task(&self) -> Option<&Task> {
        let visible_tasks = self.get_visible_tasks();
        visible_tasks
            .get(self.selected_task_index)
            .map(|(task, _)| *task)
    }

    // 递归查找任务
    fn find_task_mut(&mut self, task_id: &str) -> Option<&mut Task> {
        for task in &mut self.tasks {
            if let Some(found) = Self::find_in_tree(task, task_id) {
                return Some(found);
            }
        }
        None
    }

    fn find_in_tree<'a>(task: &'a mut Task, task_id: &str) -> Option<&'a mut Task> {
        if task.id == task_id {
            return Some(task);
        }
        for child in &mut task.children {
            if let Some(found) = Self::find_in_tree(child, task_id) {
                return Some(found);
            }
        }
        None
    }
}
