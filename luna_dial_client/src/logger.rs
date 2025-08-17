use chrono::Local;
use std::fs;
use std::path::Path;
use tracing_subscriber::prelude::*;

pub fn init_logger() -> Result<(), Box<dyn std::error::Error>> {
    let logs_dir = "logs";
    if !Path::new(logs_dir).exists() {
        fs::create_dir_all(logs_dir)?;
        println!("创建日志目录: {}", logs_dir);
    }

    let current_date = Local::now().format("%Y-%m-%d").to_string();
    let log_file_path = format!("{}/luna_dial_client_{}.log", logs_dir, current_date);

    // 创建文件追加器
    let file_appender = tracing_appender::rolling::never(
        logs_dir,
        format!("luna_dial_client_{}.log", current_date),
    );
    let (non_blocking, _guard) = tracing_appender::non_blocking(file_appender);

    // 设置 tracing subscriber
    tracing_subscriber::registry()
        .with(
            tracing_subscriber::fmt::layer()
                .with_writer(non_blocking)
                .with_ansi(false) // 文件中不使用颜色
                .with_target(true)
                .with_level(true)
                .with_thread_ids(false)
                .with_file(false)
                .with_line_number(false),
        )
        .init();

    tracing::info!("日志系统初始化完成，日志文件: {}", log_file_path);
    Ok(())
}
