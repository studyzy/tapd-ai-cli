// Package cmd 中的 read_description_test.go 测试 readDescription 函数
package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// TestReadDescription_FromFlag 测试通过 --description flag 传入描述
func TestReadDescription_FromFlag(t *testing.T) {
	// 保存并恢复全局变量
	origDesc := flagDescription
	origFile := flagDescFile
	defer func() {
		flagDescription = origDesc
		flagDescFile = origFile
	}()

	flagDescription = "来自 flag 的描述"
	flagDescFile = ""

	result, err := readDescription()
	if err != nil {
		t.Fatalf("readDescription() error: %v", err)
	}
	if result != "来自 flag 的描述" {
		t.Errorf("readDescription() = %q, want %q", result, "来自 flag 的描述")
	}
}

// TestReadDescription_FromFile 测试通过 --file flag 从文件读取描述
func TestReadDescription_FromFile(t *testing.T) {
	origDesc := flagDescription
	origFile := flagDescFile
	defer func() {
		flagDescription = origDesc
		flagDescFile = origFile
	}()

	// 创建临时文件
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
	if result != content {
		t.Errorf("readDescription() = %q, want %q", result, content)
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

	// 创建临时文件
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "desc.md")
	if err := os.WriteFile(tmpFile, []byte("文件内容"), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	flagDescription = "flag 内容"
	flagDescFile = tmpFile

	result, err := readDescription()
	if err != nil {
		t.Fatalf("readDescription() error: %v", err)
	}
	if result != "flag 内容" {
		t.Errorf("readDescription() = %q, want %q", result, "flag 内容")
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

	// 注意：在测试环境中 stdin 通常是终端设备，所以 readDescription 应返回空字符串
	result, err := readDescription()
	if err != nil {
		t.Fatalf("readDescription() error: %v", err)
	}
	if result != "" {
		t.Errorf("readDescription() = %q, want empty string", result)
	}
}
