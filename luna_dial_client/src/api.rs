use crate::models::{ApiResponse, PlanData, PlanRequest, Task};
use chrono::{Local, Utc};
use reqwest::Client;
#[derive(Debug, Default)]
pub struct ApiClient {
    client: Client,
    base_url: String,
}

impl ApiClient {
    pub fn new(base_url: String) -> Self {
        let client = Client::new();
        ApiClient { client, base_url }
    }

    // 只写了今日任务，但是后续应该是 task + document
    pub async fn get_today_plan(&self) -> Result<Vec<Task>, ApiError> {
        let today = chrono::Utc::now().date_naive();
        let start_time = today
            .and_hms_opt(0, 0, 0)
            .unwrap()
            .and_local_timezone(Local)
            .unwrap()
            .with_timezone(&Utc);
        let end_time = start_time + chrono::Duration::days(1);
        let request_body = PlanRequest {
            //TODO 这里注意 感觉不对劲 可能是服务端也有要修改的地方
            period_type: "day".to_string(),
            start_date: start_time,
            end_date: end_time,
        };
        let url = format!("{}/api/v1/plans", self.base_url);
        let response = self
            .client
            .post(url)
            .header("Authorization", "Bearer your_session_id") // 需要处理认证
            .json(&request_body) // 自动序列化为JSON
            .send()
            .await?;

        // 4. 检查HTTP状态码
        if !response.status().is_success() {
            return Err(ApiError::Business {
                code: response.status().as_u16() as u32,
                message: format!("HTTP错误: {}", response.status()),
            });
        }

        // 5. 解析响应体
        let api_response: ApiResponse<PlanData> = response.json().await?;

        // 6. 检查业务状态码
        if !api_response.success {
            return Err(ApiError::Business {
                code: api_response.code,
                message: api_response.message,
            });
        }

        // 7. 返回任务列表
        Ok(api_response.data.tasks)
    }
}

#[derive(Debug)]
pub enum ApiError {
    Network(reqwest::Error),
    Parse(serde_json::Error),
    Business { code: u32, message: String },
    Auth,
}

impl From<reqwest::Error> for ApiError {
    fn from(err: reqwest::Error) -> Self {
        ApiError::Network(err)
    }
}

impl From<serde_json::Error> for ApiError {
    fn from(err: serde_json::Error) -> Self {
        ApiError::Parse(err)
    }
}

impl std::fmt::Display for ApiError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            ApiError::Network(err) => write!(f, "网络错误: {}", err),
            ApiError::Parse(err) => write!(f, "解析错误: {}", err),
            ApiError::Business { code, message } => write!(f, "API错误 {}: {}", code, message),
            ApiError::Auth => write!(f, "认证错误"),
        }
    }
}

impl std::error::Error for ApiError {}
