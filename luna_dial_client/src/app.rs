use std::io;

use crossterm::event::{self, Event, KeyCode, KeyEvent, KeyEventKind};
use ratatui::{
    DefaultTerminal, Frame,
    buffer::Buffer,
    layout::{Constraint, Direction, Layout, Rect},
    style::{Style, Color, Stylize},
    symbols::border,
    text::{Line, Text, Span},
    widgets::{Block, Paragraph, Widget},
};

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
}

impl App {
    pub fn new() -> Self {
        App {
            view_mode: ViewMode::GlobalTree,
            input_mode: InputMode::Normal,
            session: Session::new(),
            running: true,
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
            }
     
        ];

        let nav_line = Line::from(nav_bar);
        let nav_paragraph = Paragraph::new(nav_line);
        frame.render_widget(nav_paragraph, chunks[0]);

        let main_block = Block::bordered().title(title);
        frame.render_widget(main_block, chunks[1]);

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
            _ => {}
        }
    }
}
