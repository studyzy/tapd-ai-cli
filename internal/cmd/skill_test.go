package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildCommandReferenceMarkdown(t *testing.T) {
	result := buildCommandReferenceMarkdown(rootCmd)

	// 应包含已知的命令分组
	expectedGroups := []string{"story", "bug", "task", "iteration", "wiki"}
	for _, group := range expectedGroups {
		if !strings.Contains(result, "### "+group) {
			t.Errorf("command reference should contain group %q", group)
		}
	}

	// 应包含 bash 代码块
	if !strings.Contains(result, "```bash") {
		t.Error("command reference should contain ```bash code blocks")
	}

	// 应包含具体命令
	if !strings.Contains(result, "tapd story list") {
		t.Error("command reference should contain 'tapd story list'")
	}
	if !strings.Contains(result, "tapd bug show") {
		t.Error("command reference should contain 'tapd bug show'")
	}
}

func TestGetGroupShort(t *testing.T) {
	tests := []struct {
		group    string
		wantNon  bool // 期望返回非空
	}{
		{group: "story", wantNon: true},
		{group: "bug", wantNon: true},
		{group: "nonexistent", wantNon: false},
	}
	for _, tt := range tests {
		got := getGroupShort(rootCmd, tt.group)
		if tt.wantNon && got == "" {
			t.Errorf("getGroupShort(%q) = empty, want non-empty", tt.group)
		}
		if !tt.wantNon && got != "" {
			t.Errorf("getGroupShort(%q) = %q, want empty", tt.group, got)
		}
	}
}

func TestDetectExistingTools(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建部分工具目录
	os.Mkdir(filepath.Join(tmpDir, ".claude"), 0755)
	os.Mkdir(filepath.Join(tmpDir, ".codebuddy"), 0755)

	selected := detectExistingTools(tmpDir)

	// .claude 应被选中（索引 0）
	if !selected[0] {
		t.Error("expected .claude to be detected")
	}
	// .codebuddy 应被选中（索引 1）
	if !selected[1] {
		t.Error("expected .codebuddy to be detected")
	}
	// .cursor 不应被选中（索引 2）
	if selected[2] {
		t.Error("expected .cursor to NOT be detected")
	}
}

func TestDetectExistingToolsEmpty(t *testing.T) {
	tmpDir := t.TempDir()

	selected := detectExistingTools(tmpDir)

	for i, s := range selected {
		if s {
			t.Errorf("expected tool %d (%s) to NOT be detected in empty dir",
				i, supportedTools[i].Name)
		}
	}
}

func TestGenerateSkillContent(t *testing.T) {
	content := generateSkillContent(rootCmd)

	// 不应包含未替换的占位符
	if strings.Contains(content, "{{COMMAND_REFERENCE}}") {
		t.Error("generated content should not contain {{COMMAND_REFERENCE}} placeholder")
	}

	// 应包含固定部分
	if !strings.Contains(content, "name: tapd") {
		t.Error("generated content should contain skill name")
	}
	if !strings.Contains(content, "go install") {
		t.Error("generated content should contain install instructions")
	}
	if !strings.Contains(content, "TAPD_ACCESS_TOKEN") {
		t.Error("generated content should contain auth instructions")
	}
	// 应包含动态生成的命令参考
	if !strings.Contains(content, "## 命令参考") {
		t.Error("generated content should contain command reference header")
	}
	if !strings.Contains(content, "tapd story") {
		t.Error("generated content should contain dynamically generated story commands")
	}
}

func TestWriteSkillFile(t *testing.T) {
	tmpDir := t.TempDir()
	content := []byte("test skill content")

	path, err := writeSkillFile(tmpDir, ".testool", content)
	if err != nil {
		t.Fatalf("writeSkillFile failed: %v", err)
	}

	expectedPath := filepath.Join(tmpDir, ".testool", "skills", "tapd", "SKILL.md")
	if path != expectedPath {
		t.Errorf("path = %q, want %q", path, expectedPath)
	}

	// 验证文件内容
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	if string(got) != "test skill content" {
		t.Errorf("file content = %q, want %q", string(got), "test skill content")
	}
}

func TestWriteSkillFileOverwrite(t *testing.T) {
	tmpDir := t.TempDir()

	// 写入第一次
	_, err := writeSkillFile(tmpDir, ".testool", []byte("old content"))
	if err != nil {
		t.Fatalf("first write failed: %v", err)
	}

	// 覆盖写入
	path, err := writeSkillFile(tmpDir, ".testool", []byte("new content"))
	if err != nil {
		t.Fatalf("overwrite failed: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read overwritten file: %v", err)
	}
	if string(got) != "new content" {
		t.Errorf("file content = %q, want %q", string(got), "new content")
	}
}

func TestRunInteractiveSelectConfirm(t *testing.T) {
	// 创建模拟输入：选择第 3 项，然后确认
	r, w, _ := os.Pipe()
	w.WriteString("3\ny\n")
	w.Close()

	devNull, _ := os.Open(os.DevNull)
	defer devNull.Close()

	selected := []bool{true, true, false, false, false, false, false, false, false, false}
	result, cancelled := runInteractiveSelect(r, devNull, selected)

	if cancelled {
		t.Error("expected not cancelled")
	}
	if !result[2] {
		t.Error("expected item 3 to be toggled on")
	}
	if !result[0] || !result[1] {
		t.Error("expected items 1,2 to remain selected")
	}
}

func TestRunInteractiveSelectCancel(t *testing.T) {
	r, w, _ := os.Pipe()
	w.WriteString("q\n")
	w.Close()

	devNull, _ := os.Open(os.DevNull)
	defer devNull.Close()

	selected := []bool{true, false, false, false, false, false, false, false, false, false}
	_, cancelled := runInteractiveSelect(r, devNull, selected)

	if !cancelled {
		t.Error("expected cancelled")
	}
}

func TestRunInteractiveSelectAll(t *testing.T) {
	r, w, _ := os.Pipe()
	w.WriteString("a\ny\n")
	w.Close()

	devNull, _ := os.Open(os.DevNull)
	defer devNull.Close()

	selected := make([]bool, len(supportedTools))
	result, cancelled := runInteractiveSelect(r, devNull, selected)

	if cancelled {
		t.Error("expected not cancelled")
	}
	for i, s := range result {
		if !s {
			t.Errorf("expected tool %d to be selected after 'a'", i)
		}
	}
}
