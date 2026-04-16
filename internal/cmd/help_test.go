package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/pflag"
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

// TestBuildSpecLines 测试参考行生成
func TestBuildSpecLines(t *testing.T) {
	lines := buildSpecLines(rootCmd)

	if len(lines) == 0 {
		t.Fatal("buildSpecLines should return non-empty slice")
	}

	// 构建文本索引方便查找
	lineMap := make(map[string]specLine)
	for _, l := range lines {
		lineMap[extractCommandPath(l.text)] = l
	}

	// 验证预期命令存在
	expectedCmds := []string{
		"tapd auth login",
		"tapd workspace list",
		"tapd story list",
		"tapd bug list",
		"tapd iteration list",
	}
	for _, cmd := range expectedCmds {
		if _, ok := lineMap[cmd]; !ok {
			t.Errorf("expected command %q not found in spec lines", cmd)
		}
	}

	// 验证每行都以 "tapd " 开头且有分组
	for _, l := range lines {
		if !strings.HasPrefix(l.text, "tapd ") {
			t.Errorf("spec line should start with 'tapd ', got: %s", l.text)
		}
		if l.group == "" {
			t.Errorf("spec line should have a group, got empty for: %s", l.text)
		}
	}

	// 验证带位置参数的命令包含 <arg>
	argsTests := []struct {
		path    string
		argName string
	}{
		{"tapd story show", "<story_id>"},
		{"tapd bug show", "<bug_id>"},
		{"tapd workspace switch", "<workspace_id>"},
	}
	for _, tc := range argsTests {
		l, ok := lineMap[tc.path]
		if !ok {
			continue
		}
		if !strings.Contains(l.text, tc.argName) {
			t.Errorf("command %q should contain %q, got: %s", tc.path, tc.argName, l.text)
		}
	}
}

// TestIsFlagRequired 测试必填标志检测
func TestIsFlagRequired(t *testing.T) {
	tests := []struct {
		usage string
		want  bool
	}{
		{"需求标题（必需）", true},
		{"迭代标题（必需）", true},
		{"按状态筛选", false},
		{"描述", false},
		{"工时填写人（必填）", true},
	}

	for _, tc := range tests {
		t.Run(tc.usage, func(t *testing.T) {
			f := &pflag.Flag{Usage: tc.usage}
			got := isFlagRequired(f)
			if got != tc.want {
				t.Errorf("isFlagRequired(usage=%q) = %v, want %v", tc.usage, got, tc.want)
			}
		})
	}
}

// TestIsGlobalAuthFlag 测试全局认证标志检测
func TestIsGlobalAuthFlag(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"access-token", true},
		{"api-user", true},
		{"api-password", true},
		{"workspace-id", false},
		{"status", false},
	}
	for _, tc := range tests {
		if got := isGlobalAuthFlag(tc.name); got != tc.want {
			t.Errorf("isGlobalAuthFlag(%q) = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestIsGlobalDisplayFlag 测试全局展示标志检测
func TestIsGlobalDisplayFlag(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"workspace-id", true},
		{"pretty", true},
		{"json", true},
		{"no-comments", true},
		{"status", false},
		{"name", false},
	}
	for _, tc := range tests {
		if got := isGlobalDisplayFlag(tc.name); got != tc.want {
			t.Errorf("isGlobalDisplayFlag(%q) = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// extractCommandPath 从参考行文本中提取命令路径（到第一个 flag、[flag]、<arg> 或 # 之前的部分）
func extractCommandPath(text string) string {
	for i, c := range text {
		if c == '-' && i+1 < len(text) && text[i+1] == '-' {
			return strings.TrimSpace(text[:i])
		}
		if c == '[' {
			return strings.TrimSpace(text[:i])
		}
		if c == '<' {
			return strings.TrimSpace(text[:i])
		}
		if c == '#' {
			return strings.TrimSpace(text[:i])
		}
	}
	return strings.TrimSpace(text)
}
