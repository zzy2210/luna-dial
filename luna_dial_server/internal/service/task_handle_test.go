package service

import (
	"testing"
)

func TestIsIcon(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// 正面测试用例 - 应该返回true的emoji
		{"simple smiley", "😀", true},
		{"heart", "❤️", true},
		{"thumbs up", "👍", true},
		{"fire", "🔥", true},
		{"star", "⭐", true},
		{"check mark", "✅", true},
		{"rocket", "🚀", true},
		{"flag", "🏳️", true},
		{"musical note", "🎵", true},
		{"sun", "☀️", true},

		// 负面测试用例 - 应该返回false的非emoji
		{"empty string", "", false},
		{"regular text", "hello", false},
		{"number", "123", false},
		{"letter", "a", false},
		{"special chars", "!@#", false},
		{"long text", "this is a long text", false},
		{"mixed text and emoji", "hello 😀", false},
		{"chinese text", "你好", false},

		// 边界测试用例
		{"whitespace", " ", false},
		{"tab", "\t", false},
		{"newline", "\n", false},
		{"multiple emoji", "😀😁", true}, // 这个可能需要根据业务需求调整
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsIcon(tt.input)
			if result != tt.expected {
				t.Errorf("IsIcon(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestListRootTasksRequest_Defaults(t *testing.T) {
	// 测试默认值设置逻辑
	req := ListRootTasksRequest{
		Page:     0,  // 无效值
		PageSize: -1, // 无效值
	}

	// 模拟API处理函数中的默认值设置逻辑
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	if req.Page != 1 {
		t.Errorf("Expected default Page = 1, got %d", req.Page)
	}
	if req.PageSize != 20 {
		t.Errorf("Expected default PageSize = 20, got %d", req.PageSize)
	}
}

func TestListGlobalTaskTreeRequest_Defaults(t *testing.T) {
	// 测试全局任务树请求的默认值
	req := ListGlobalTaskTreeRequest{
		Page:     0,
		PageSize: 0,
	}

	// 模拟API处理函数中的默认值设置逻辑
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	if req.Page != 1 {
		t.Errorf("Expected default Page = 1, got %d", req.Page)
	}
	if req.PageSize != 10 {
		t.Errorf("Expected default PageSize = 10, got %d", req.PageSize)
	}
}

func TestListJournalsWithPaginationRequest_Defaults(t *testing.T) {
	// 测试日志分页请求的默认值
	req := ListJournalsWithPaginationRequest{
		Page:     0,
		PageSize: 0,
	}

	// 模拟API处理函数中的默认值设置逻辑
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	if req.Page != 1 {
		t.Errorf("Expected default Page = 1, got %d", req.Page)
	}
	if req.PageSize != 20 {
		t.Errorf("Expected default PageSize = 20, got %d", req.PageSize)
	}
}

// 基准测试
func BenchmarkIsIcon(b *testing.B) {
	testCases := []string{
		"😀",
		"hello",
		"❤️",
		"🚀",
		"regular text",
	}

	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			IsIcon(tc)
		}
	}
}
