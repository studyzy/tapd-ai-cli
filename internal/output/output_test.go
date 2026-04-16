package output_test

import (
	"bytes"
	"encoding/json"
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
