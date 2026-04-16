package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/studyzy/tapd-ai-cli/internal/model"
	"github.com/studyzy/tapd-ai-cli/internal/output"
)

func TestPrintJSON_Compact(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]string{"name": "test", "value": "123"}
	err := output.PrintJSON(&buf, data, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	// 紧凑模式不应包含缩进
	want := `{"name":"test","value":"123"}` + "\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrintJSON_Indent(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]string{"key": "val"}
	err := output.PrintJSON(&buf, data, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	want := "{\n  \"key\": \"val\"\n}\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPrintJSON_OmitEmpty(t *testing.T) {
	var buf bytes.Buffer
	// Story 结构体中空字符串字段应被 omitempty 省略
	data := model.Story{ID: "123", Name: "test story"}
	err := output.PrintJSON(&buf, data, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	// 空字符串字段不应出现在输出中
	for _, field := range []string{"status", "owner", "priority", "description", "iteration_id", "modified", "url"} {
		if _, ok := m[field]; ok {
			t.Errorf("expected field %q to be omitted, but it was present", field)
		}
	}
	// 非空字段应存在
	if m["id"] != "123" {
		t.Errorf("expected id=123, got %v", m["id"])
	}
	if m["name"] != "test story" {
		t.Errorf("expected name=test story, got %v", m["name"])
	}
}

func TestPrintError(t *testing.T) {
	var buf bytes.Buffer
	output.PrintError(&buf, "auth_error", "invalid token", "run tapd auth login")
	got := buf.String()
	var resp model.ErrorResponse
	if err := json.Unmarshal([]byte(got), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if resp.Error != "auth_error" {
		t.Errorf("error: got %q, want %q", resp.Error, "auth_error")
	}
	if resp.Message != "invalid token" {
		t.Errorf("message: got %q, want %q", resp.Message, "invalid token")
	}
	if resp.Hint != "run tapd auth login" {
		t.Errorf("hint: got %q, want %q", resp.Hint, "run tapd auth login")
	}
}

func TestPrintError_NoHTMLEscape(t *testing.T) {
	var buf bytes.Buffer
	output.PrintError(&buf, "auth_error", "invalid", "run 'tapd auth login --api-user <user>'")
	got := buf.String()
	// 确保尖括号不被转义为 \u003c / \u003e
	if strings.Contains(got, `\u003c`) || strings.Contains(got, `\u003e`) {
		t.Errorf("angle brackets should not be HTML-escaped, got: %s", got)
	}
	if !strings.Contains(got, "<user>") {
		t.Errorf("expected literal <user> in output, got: %s", got)
	}
}

func TestPrintError_EmptyHint(t *testing.T) {
	var buf bytes.Buffer
	output.PrintError(&buf, "not_found", "resource not found", "")
	got := buf.String()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	// hint 为空时应被 omitempty 省略
	if _, ok := m["hint"]; ok {
		t.Error("expected hint field to be omitted when empty")
	}
	if m["error"] != "not_found" {
		t.Errorf("error: got %v, want not_found", m["error"])
	}
	if m["message"] != "resource not found" {
		t.Errorf("message: got %v, want resource not found", m["message"])
	}
}

func TestPrintSuccess(t *testing.T) {
	var buf bytes.Buffer
	data := model.SuccessResponse{Success: true, ID: "456"}
	err := output.PrintSuccess(&buf, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	var resp model.SuccessResponse
	if err := json.Unmarshal([]byte(got), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.ID != "456" {
		t.Errorf("id: got %q, want %q", resp.ID, "456")
	}
}

func TestExitCodes(t *testing.T) {
	tests := []struct {
		name string
		code int
		want int
	}{
		{"ExitSuccess", output.ExitSuccess, 0},
		{"ExitAuthError", output.ExitAuthError, 1},
		{"ExitNotFound", output.ExitNotFound, 2},
		{"ExitParamError", output.ExitParamError, 3},
		{"ExitAPIError", output.ExitAPIError, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code != tt.want {
				t.Errorf("%s = %d, want %d", tt.name, tt.code, tt.want)
			}
		})
	}
}

func TestPrintMarkdown_StoryWithDescription(t *testing.T) {
	var buf bytes.Buffer
	story := model.Story{
		ID:          "123",
		Name:        "登录优化",
		Status:      "open",
		Owner:       "张三",
		Description: "# 需求背景\n\n这是一个登录优化需求。",
	}
	err := output.PrintMarkdown(&buf, story, "description")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()

	// 应以 YAML frontmatter 开头
	if !strings.HasPrefix(got, "---\n") {
		t.Errorf("expected output to start with ---\\n, got %q", got[:20])
	}

	// frontmatter 中应包含元数据字段
	if !strings.Contains(got, "id: 123") {
		t.Errorf("expected frontmatter to contain id: 123, got:\n%s", got)
	}
	if !strings.Contains(got, "name: 登录优化") {
		t.Errorf("expected frontmatter to contain name: 登录优化, got:\n%s", got)
	}
	if !strings.Contains(got, "status: open") {
		t.Errorf("expected frontmatter to contain status: open, got:\n%s", got)
	}
	if !strings.Contains(got, "owner: 张三") {
		t.Errorf("expected frontmatter to contain owner: 张三, got:\n%s", got)
	}

	// description 不应出现在 frontmatter 中
	parts := strings.SplitN(got, "---\n", 3)
	if len(parts) < 3 {
		t.Fatalf("expected output to have frontmatter and body, got:\n%s", got)
	}
	frontmatter := parts[1]
	if strings.Contains(frontmatter, "description:") {
		t.Errorf("description should not appear in frontmatter, got:\n%s", frontmatter)
	}

	// body 部分应包含描述内容
	body := parts[2]
	if !strings.Contains(body, "# 需求背景") {
		t.Errorf("expected body to contain description content, got:\n%s", body)
	}
}

func TestPrintMarkdown_OmitEmpty(t *testing.T) {
	var buf bytes.Buffer
	story := model.Story{ID: "456", Name: "测试需求"}
	err := output.PrintMarkdown(&buf, story, "description")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()

	// 空字段不应出现在 frontmatter 中
	if strings.Contains(got, "status:") {
		t.Errorf("empty status should be omitted, got:\n%s", got)
	}
	if strings.Contains(got, "owner:") {
		t.Errorf("empty owner should be omitted, got:\n%s", got)
	}

	// 非空字段应存在
	if !strings.Contains(got, "id: 456") {
		t.Errorf("expected id: 456 in output, got:\n%s", got)
	}
}

func TestPrintMarkdown_NoDescription(t *testing.T) {
	var buf bytes.Buffer
	story := model.Story{ID: "789", Name: "无描述需求", Status: "done"}
	err := output.PrintMarkdown(&buf, story, "description")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()

	// 应有 frontmatter，但没有 body
	if !strings.HasPrefix(got, "---\n") {
		t.Errorf("expected output to start with ---\\n")
	}
	// frontmatter 结束后不应有多余内容（除了结尾换行）
	parts := strings.SplitN(got, "---\n", 3)
	if len(parts) < 3 {
		t.Fatalf("expected frontmatter format, got:\n%s", got)
	}
	body := strings.TrimSpace(parts[2])
	if body != "" {
		t.Errorf("expected empty body when no description, got: %q", body)
	}
}

func TestPrintMarkdown_BugWithDescription(t *testing.T) {
	var buf bytes.Buffer
	bug := model.Bug{
		ID:          "100",
		Title:       "登录崩溃",
		Status:      "new",
		Severity:    "fatal",
		Description: "## 复现步骤\n\n1. 打开登录页\n2. 点击提交",
	}
	err := output.PrintMarkdown(&buf, bug, "description")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()

	if !strings.Contains(got, "title: 登录崩溃") {
		t.Errorf("expected title in frontmatter, got:\n%s", got)
	}
	if !strings.Contains(got, "## 复现步骤") {
		t.Errorf("expected description content in body, got:\n%s", got)
	}
}

func TestPrintMarkdown_Pointer(t *testing.T) {
	var buf bytes.Buffer
	story := &model.Story{ID: "111", Name: "指针测试"}
	err := output.PrintMarkdown(&buf, story, "description")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "id: 111") {
		t.Errorf("expected id in output for pointer input, got:\n%s", got)
	}
}

func TestPrintMarkdown_NonStruct(t *testing.T) {
	var buf bytes.Buffer
	err := output.PrintMarkdown(&buf, "not a struct", "description")
	if err == nil {
		t.Fatal("expected error for non-struct input, got nil")
	}
}
