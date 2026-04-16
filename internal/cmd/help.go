// Package cmd 中的 spec.go 实现了紧凑参考卡生成逻辑，供根命令 --help 使用
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// specLine 表示一条命令参考行
type specLine struct {
	group string // 命令所属分组（第一级子命令名）
	text  string // 完整的命令参考文本
}

// buildSpecLines 遍历命令树，为每个叶子命令生成参考行
func buildSpecLines(root *cobra.Command) []specLine {
	var lines []specLine
	walkSpecCommands(root, "", "", &lines)
	return lines
}

// walkSpecCommands 递归遍历命令树，收集叶子命令的参考行
func walkSpecCommands(cmd *cobra.Command, prefix string, group string, lines *[]specLine) {
	for _, child := range cmd.Commands() {
		if child.Hidden || child.Name() == "help" || child.Name() == "completion" {
			continue
		}

		fullPath := child.Name()
		if prefix != "" {
			fullPath = prefix + " " + child.Name()
		}

		// 确定分组名：取第一级子命令名
		currentGroup := group
		if currentGroup == "" {
			currentGroup = child.Name()
		}

		if child.HasSubCommands() {
			walkSpecCommands(child, fullPath, currentGroup, lines)
		} else {
			line := commandToLine(child, fullPath)
			*lines = append(*lines, specLine{group: currentGroup, text: line})
		}
	}
}

// commandToLine 将 Cobra 命令转换为一行紧凑参考文本
func commandToLine(cmd *cobra.Command, path string) string {
	var b strings.Builder
	b.WriteString("tapd ")
	b.WriteString(path)

	// 添加位置参数
	if argName := extractArgName(cmd.Use); argName != "" {
		b.WriteString(" <")
		b.WriteString(argName)
		b.WriteString(">")
	}

	// 收集标志，区分必填和可选
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden || f.Name == "help" {
			return
		}
		// 跳过全局认证标志（spec 不需要展示）
		if isGlobalAuthFlag(f.Name) {
			return
		}
		// 跳过全局展示标志（已在 header 中展示）
		if isGlobalDisplayFlag(f.Name) {
			return
		}
		b.WriteString(" ")
		b.WriteString(formatFlag(f))
	})

	// 添加描述注释
	if cmd.Short != "" {
		b.WriteString("  # ")
		b.WriteString(cmd.Short)
	}

	return b.String()
}

// formatFlag 将一个 flag 格式化为紧凑文本
// 必填标志：--flag=<val>
// 可选带默认值：[--flag=default]
// 可选无默认值：[--flag]
func formatFlag(f *pflag.Flag) string {
	if isFlagRequired(f) {
		return "--" + f.Name + "=<" + f.Name + ">"
	}
	if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" {
		return "[--" + f.Name + "=" + f.DefValue + "]"
	}
	return "[--" + f.Name + "]"
}

// isFlagRequired 判断标志是否为必填（通过检测 Usage 文本中的关键字）
func isFlagRequired(f *pflag.Flag) bool {
	usage := f.Usage
	return strings.Contains(usage, "必需") || strings.Contains(usage, "必填")
}

// isGlobalAuthFlag 判断是否为全局认证标志
func isGlobalAuthFlag(name string) bool {
	switch name {
	case "access-token", "api-user", "api-password":
		return true
	default:
		return false
	}
}

// isGlobalDisplayFlag 判断是否为全局展示标志（header 中已展示）
func isGlobalDisplayFlag(name string) bool {
	switch name {
	case "workspace-id", "pretty", "json", "no-comments":
		return true
	default:
		return false
	}
}

// printSpecOutput 输出完整的参考卡文本
func printSpecOutput(w *os.File, root *cobra.Command, lines []specLine) {
	// 标题行
	fmt.Fprintf(w, "tapd - %s\n", root.Short)
	fmt.Fprintln(w, "Global: [--workspace-id=<id>] [--json] [--pretty] [--no-comments]")

	// 按分组输出
	lastGroup := ""
	for _, l := range lines {
		if l.group != lastGroup {
			fmt.Fprintln(w)
			fmt.Fprintf(w, "# %s\n", l.group)
			lastGroup = l.group
		}
		fmt.Fprintln(w, l.text)
	}
}

// extractArgName 从 Use 字段提取位置参数名（如 "show <story_id>" -> "story_id"）
func extractArgName(use string) string {
	start := -1
	for i, c := range use {
		if c == '<' {
			start = i + 1
		} else if c == '>' && start >= 0 {
			return use[start:i]
		}
	}
	return ""
}
