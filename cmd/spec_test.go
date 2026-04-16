package cmd

import (
	"testing"
)

// TestExtractArgName 测试从 Use 字段提取位置参数名
func TestExtractArgName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"show <story_id>", "story_id"},
		{"list", ""},
		{"update <bug_id>", "bug_id"},
		{"switch <workspace_id>", "workspace_id"},
		{"login", ""},
		{"show <id> extra", "id"},
		{"no-angle-brackets", ""},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := extractArgName(tc.input)
			if got != tc.want {
				t.Errorf("extractArgName(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

// TestBuildToolDefinitions 测试工具定义生成
func TestBuildToolDefinitions(t *testing.T) {
	tools := buildToolDefinitions(rootCmd)

	if len(tools) == 0 {
		t.Fatal("buildToolDefinitions should return non-empty slice")
	}

	// 构建名称索引方便查找
	toolMap := make(map[string]toolDefinition)
	for _, tool := range tools {
		toolMap[tool.Name] = tool
	}

	// 验证预期工具存在
	expectedTools := []string{
		"tapd_auth_login",
		"tapd_workspace_list",
		"tapd_story_list",
		"tapd_bug_list",
		"tapd_iteration_list",
		"tapd_spec",
	}
	for _, name := range expectedTools {
		if _, ok := toolMap[name]; !ok {
			t.Errorf("expected tool %q not found in definitions", name)
		}
	}

	// 验证每个工具都有非空的 Name 和 Description
	for _, tool := range tools {
		if tool.Name == "" {
			t.Error("tool Name should not be empty")
		}
		if tool.Description == "" {
			t.Errorf("tool %q Description should not be empty", tool.Name)
		}
		if tool.Parameters.Type != "object" {
			t.Errorf("tool %q Parameters.Type = %q, want %q", tool.Name, tool.Parameters.Type, "object")
		}
	}

	// 验证带位置参数的工具有 required 字段
	toolsWithArgs := []struct {
		name     string
		argName  string
	}{
		{"tapd_story_show", "story_id"},
		{"tapd_bug_show", "bug_id"},
		{"tapd_workspace_switch", "workspace_id"},
	}
	for _, tc := range toolsWithArgs {
		tool, ok := toolMap[tc.name]
		if !ok {
			// 工具可能不存在（取决于命令注册），跳过
			continue
		}
		if len(tool.Parameters.Required) == 0 {
			t.Errorf("tool %q should have required parameters", tc.name)
		}
		// 验证 required 中包含预期的参数名
		foundArg := false
		for _, r := range tool.Parameters.Required {
			if r == tc.argName {
				foundArg = true
				break
			}
		}
		if !foundArg {
			t.Errorf("tool %q required should contain %q, got %v", tc.name, tc.argName, tool.Parameters.Required)
		}
		// 验证 properties 中也有该参数
		if _, ok := tool.Parameters.Properties[tc.argName]; !ok {
			t.Errorf("tool %q properties should contain %q", tc.name, tc.argName)
		}
	}
}

// TestSpecCommand_Exists 验证 specCmd 已注册在 rootCmd 下
func TestSpecCommand_Exists(t *testing.T) {
	found := false
	for _, sub := range rootCmd.Commands() {
		if sub.Name() == "spec" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("rootCmd should have 'spec' subcommand")
	}

	if specCmd.Use != "spec" {
		t.Errorf("specCmd.Use = %q, want %q", specCmd.Use, "spec")
	}
	if specCmd.Short == "" {
		t.Error("specCmd.Short should not be empty")
	}
	if specCmd.RunE == nil {
		t.Error("specCmd.RunE should not be nil")
	}
}
