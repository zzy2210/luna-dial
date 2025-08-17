use std::{collections::HashSet, io, time::Duration, vec};

use chrono::{Local, Utc};
use crossterm::event::{self, Event, KeyCode, KeyEvent, KeyEventKind, poll};
use ratatui::{
    DefaultTerminal, Frame,
    layout::{Constraint, Direction, Layout},
    prelude::Rect,
    style::{Color, Style},
    text::{Line, Span, Text},
    widgets::{Block, List, ListItem, Paragraph},
    widgets::{Widget, Wrap},
};

use crate::api::ApiClient;
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

    pub api_client: ApiClient, // API 客户端
    pub loading: bool,         // 显示加载状态

    pub tasks: Vec<Task>,                // 所有任务列表
    pub selected_task_index: usize,      // 当前选中的任务索引
    pub expanded_tasks: HashSet<String>, // 记录展开的任务ID

    // 错误处理 感觉不需要pub
    error_msg: Option<String>, // 错误信息
    show_error: bool,          // 是否显示错误提示
}

impl App {
    pub fn new() -> Self {
        //TODO  空数据 后续或许可以删除？
        let test_tasks = vec![];
        let api_client = ApiClient::new("http://localhost:8081".to_string());
        App {
            view_mode: ViewMode::GlobalTree,
            input_mode: InputMode::Normal,
            session: Session::new(),
            running: true,
            tasks: test_tasks,
            selected_task_index: 0,
            api_client: api_client,
            loading: false,
            expanded_tasks: HashSet::new(), // 初始所有任务都是折叠的
            show_error: false,
            error_msg: None, // 初始没有错误信息
        }
    }
    pub async fn run(&mut self, terminal: &mut DefaultTerminal) -> io::Result<()> {
        while self.running {
            terminal.draw(|frame| self.draw(frame))?;
            self.handle_events().await?;
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
        match self.view_mode {
            ViewMode::GlobalTree => {
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
            }
            ViewMode::Today => {
                // 确保加载今日任务

                let horizontal_chunks = Layout::default()
                    .direction(Direction::Horizontal)
                    .constraints([Constraint::Percentage(60), Constraint::Percentage(40)])
                    .split(chunks[1]);

                let today_tasks = self.get_visible_tasks();
                let task_items: Vec<ListItem> = today_tasks
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
                let task_list = List::new(task_items).block(Block::bordered().title("📋 今日任务"));
                frame.render_widget(task_list, horizontal_chunks[0]);

                // 右侧：相关文档列表（暂时显示占位符）
                let doc_items = vec![
                    ListItem::new("📝 设计文档v1.0"),
                    ListItem::new("📊 API接口文档"),
                    ListItem::new("💭 需求分析"),
                ];

                let doc_list = List::new(doc_items).block(Block::bordered().title("📚 相关文档"));
                frame.render_widget(doc_list, horizontal_chunks[1]);
            }
            _ => {
                let text = Text::from("当前视图未实现");
                let paragraph = Paragraph::new(text).block(main_block);
                frame.render_widget(paragraph, chunks[1]);
            }
        }

        let help_text = "Tab: 切换视图 | q: 退出";
        let help_paragraph = Paragraph::new(help_text);

        frame.render_widget(help_paragraph, chunks[2]);
        if self.show_error {
            if let Some(error_msg) = &self.error_msg {
                // 创建一个居中的错误弹窗
                let error_area = self.centered_rect(60, 20, frame.area());
                let error_block = Block::bordered()
                    .title("错误")
                    .style(Style::default().fg(Color::Red));
                let error_paragraph = Paragraph::new(error_msg.as_str())
                    .block(error_block)
                    .wrap(Wrap { trim: true });
                frame.render_widget(error_paragraph, error_area);
            }
        }
    }

    // 辅助函数：创建居中的矩形区域
    fn centered_rect(&self, percent_x: u16, percent_y: u16, r: Rect) -> Rect {
        let popup_layout = Layout::default()
            .direction(Direction::Vertical)
            .constraints([
                Constraint::Percentage((100 - percent_y) / 2),
                Constraint::Percentage(percent_y),
                Constraint::Percentage((100 - percent_y) / 2),
            ])
            .split(r);

        Layout::default()
            .direction(Direction::Horizontal)
            .constraints([
                Constraint::Percentage((100 - percent_x) / 2),
                Constraint::Percentage(percent_x),
                Constraint::Percentage((100 - percent_x) / 2),
            ])
            .split(popup_layout[1])[1]
    }

    async fn handle_events(&mut self) -> io::Result<()> {
        if event::poll(Duration::from_millis(250))? {
            if let Event::Key(key_event) = event::read()? {
                if key_event.kind == KeyEventKind::Press {
                    // 处理按键事件
                    self.handle_key_event(key_event).await;
                }
            }
        }

        Ok(())
    }

    async fn handle_key_event(&mut self, key_event: KeyEvent) {
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

                if matches!(self.view_mode, ViewMode::Today) {
                    // 切换到今日视图时，加载今日任务
                    if let Err(e) = self.load_today_tasks().await {
                        // eprintln!("加载今日任务失败: {}", e);
                        // TODO
                        /*
                        1. 记录log error
                        2. 打印错误
                         */

                        self.error_msg = Some(format!("加载今日任务失败: {}", e));
                        self.show_error = true;
                    } else {
                        self.show_error = false; // 成功加载任务后隐藏错误
                        self.error_msg = None; // 清除错误信息
                    }
                }
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

            //关闭错误提示
            KeyCode::Esc => {
                if self.show_error {
                    self.show_error = false;
                    self.error_msg = None; // 清除错误信息
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

    // 请求获取今日任务
    pub async fn load_today_tasks(&mut self) -> Result<(), crate::api::ApiError> {
        self.loading = true;
        match self.api_client.get_today_plan().await {
            Ok(tasks) => {
                self.tasks = tasks;
                self.loading = false;
                Ok(()) // 成功时返回 Ok(())
            }
            Err(e) => {
                self.loading = false;
                Err(e) // 失败时返回错误，让调用者处理
            }
        }
    }
}
