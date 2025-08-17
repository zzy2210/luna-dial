use std::io;

mod api;
mod app;
mod logger;
mod models;
mod session;

use app::App;
use crossterm::event::{self, Event, KeyCode, KeyEvent, KeyEventKind};
use ratatui::{
    DefaultTerminal, Frame,
    buffer::Buffer,
    layout::Rect,
    style::Stylize,
    symbols::border,
    text::{Line, Text},
    widgets::{Block, Paragraph, Widget},
};

#[tokio::main]
async fn main() -> io::Result<()> {
    logger::init_logger().expect("日志系统初始化失败");
    let mut terminal = ratatui::init();
    let mut app = App::new();
    let result = app.run(&mut terminal).await;
    ratatui::restore();
    result
}
