#[derive(Debug, Default)]
pub struct Session {
    pub auth_token:Option<String>,
}


impl Session {
    pub fn new() -> Self {
        Session {
            auth_token: None,
        }
    }

    pub fn set_auth_token(&mut self, token: String) {
        self.auth_token = Some(token);
    }

    pub fn clear_auth_token(&mut self) {
        self.auth_token = None;
    }
}