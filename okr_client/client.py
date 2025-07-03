class OKRClient:
    """OKR API 客户端"""

    def create_task(self, task_request: TaskRequest) -> Task:
        """创建任务"""
        payload = task_request.model_dump(mode="json", exclude_none=True)
        print("[调试] create_task 请求体:", payload)
        response = self._request("POST", "/tasks", json=payload)
        return Task(**response["data"])

    def update_task(self, task_id: str, task_request: TaskRequest) -> Task:
        """更新任务"""
        payload = task_request.model_dump(mode="json", exclude_none=True)
        print("[调试] update_task 请求体:", payload)
        response = self._request("PUT", f"/tasks/{task_id}", json=payload)
        return Task(**response["data"])

    def create_journal(self, journal_request: JournalRequest) -> JournalEntry:
        """创建日志条目"""
        payload = journal_request.model_dump(mode="json", exclude_none=True)
        print("[调试] create_journal 请求体:", payload)
        response = self._request("POST", "/journals", json=payload)
        return JournalEntry(**response["data"])

    def update_journal(self, journal_id: str, journal_request: JournalRequest) -> JournalEntry:
        """更新日志条目"""
        payload = journal_request.model_dump(mode="json", exclude_none=True)
        print("[调试] update_journal 请求体:", payload)
        response = self._request("PUT", f"/journals/{journal_id}", json=payload)
        return JournalEntry(**response["data"])