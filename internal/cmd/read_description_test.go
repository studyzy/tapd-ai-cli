// Package cmd 中的 read_description_test.go 测试 readDescription 函数
package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestReadDescription_FromFlag 测试通过 --description flag 传入纯文本描述会被转换为 HTML
func TestReadDescription_FromFlag(t *testing.T) {
	origDesc := flagDescription
	origFile := flagDescFile
	defer func() {
		flagDescription = origDesc
		flagDescFile = origFile
	}()

	flagDescription = "简单纯文本描述"
	flagDescFile = ""

	result, err := readDescription()
	if err != nil {
		t.Fatalf("readDescription() error: %v", err)
	}
	// 纯文本被 markdownToHTML 包裹为 <p> 标签
	if !strings.Contains(result, "简单纯文本描述") {
		t.Errorf("readDescription() should contain original text, got %q", result)
	}
}

// TestReadDescription_FromFlag_Markdown 测试通过 --description flag 传入 Markdown 内容会被转换为 HTML
func TestReadDescription_FromFlag_Markdown(t *testing.T) {
	origDesc := flagDescription
	origFile := flagDescFile
	defer func() {
		flagDescription = origDesc
		flagDescFile = origFile
	}()

	flagDescription = "# 标题\n\n段落内容"
	flagDescFile = ""

	result, err := readDescription()
	if err != nil {
		t.Fatalf("readDescription() error: %v", err)
	}
	if !strings.Contains(result, "<h1") {
		t.Errorf("readDescription() should convert markdown heading to HTML, got %q", result)
	}
	if !strings.Contains(result, "<p>") {
		t.Errorf("readDescription() should convert markdown paragraph to HTML, got %q", result)
	}
}

// TestReadDescription_FromFile 测试通过 --file flag 从文件读取 Markdown 内容并转换为 HTML
func TestReadDescription_FromFile(t *testing.T) {
	origDesc := flagDescription
	origFile := flagDescFile
	defer func() {
		flagDescription = origDesc
		flagDescFile = origFile
	}()

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "desc.md")
	content := "# 从文件读取的描述\n\n这是一段详细描述。"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	flagDescription = ""
	flagDescFile = tmpFile

	result, err := readDescription()
	if err != nil {
		t.Fatalf("readDescription() error: %v", err)
	}
	// Markdown 文件应被转换为 HTML
	if !strings.Contains(result, "<h1") {
		t.Errorf("readDescription() should convert file markdown to HTML, got %q", result)
	}
	if !strings.Contains(result, "<p>") {
		t.Errorf("readDescription() should contain <p> tag, got %q", result)
	}
}

// TestReadDescription_FlagPriorityOverFile 测试 --description 优先于 --file
func TestReadDescription_FlagPriorityOverFile(t *testing.T) {
	origDesc := flagDescription
	origFile := flagDescFile
	defer func() {
		flagDescription = origDesc
		flagDescFile = origFile
	}()

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "desc.md")
	if err := os.WriteFile(tmpFile, []byte("# 文件内容"), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	flagDescription = "flag 内容"
	flagDescFile = tmpFile

	result, err := readDescription()
	if err != nil {
		t.Fatalf("readDescription() error: %v", err)
	}
	// flag 优先，内容经过 markdownToHTML 转换
	if !strings.Contains(result, "flag 内容") {
		t.Errorf("readDescription() should contain 'flag 内容', got %q", result)
	}
	// 不应包含文件内容
	if strings.Contains(result, "文件内容") {
		t.Errorf("readDescription() should not contain file content, got %q", result)
	}
}

// TestReadDescription_FileNotFound 测试文件不存在时返回错误
func TestReadDescription_FileNotFound(t *testing.T) {
	origDesc := flagDescription
	origFile := flagDescFile
	defer func() {
		flagDescription = origDesc
		flagDescFile = origFile
	}()

	flagDescription = ""
	flagDescFile = "/nonexistent/path/desc.md"

	_, err := readDescription()
	if err == nil {
		t.Error("readDescription() should return error for nonexistent file")
	}
}

// TestReadDescription_EmptyFlagsNoStdin 测试所有来源为空时返回空字符串
func TestReadDescription_EmptyFlagsNoStdin(t *testing.T) {
	origDesc := flagDescription
	origFile := flagDescFile
	defer func() {
		flagDescription = origDesc
		flagDescFile = origFile
	}()

	flagDescription = ""
	flagDescFile = ""

	result, err := readDescription()
	if err != nil {
		t.Fatalf("readDescription() error: %v", err)
	}
	if result != "" {
		t.Errorf("readDescription() = %q, want empty string", result)
	}
}
