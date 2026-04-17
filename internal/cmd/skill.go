package cmd

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/output"
)

//go:embed skill_template.md
var skillTemplate string

// aiTool 定义一个 AI Coding 工具的元信息
type aiTool struct {
	Name string // 显示名称
	Dir  string // 配置目录名，如 ".claude"
}

// supportedTools 支持的 AI Coding 工具列表
var supportedTools = []aiTool{
	{Name: "Claude Code", Dir: ".claude"},
	{Name: "CodeBuddy", Dir: ".codebuddy"},
	{Name: "Cursor", Dir: ".cursor"},
	{Name: "Windsurf", Dir: ".windsurf"},
	{Name: "Trae", Dir: ".trae"},
	{Name: "Codex", Dir: ".codex"},
	{Name: "Gemini CLI", Dir: ".gemini"},
	{Name: "Cline", Dir: ".cline"},
	{Name: "Roo Code", Dir: ".roo"},
	{Name: "Augment", Dir: ".augment"},
}

// skillCmd 是 skill 父命令
var skillCmd = &cobra.Command{
	Use:   "skill",
	Short: "AI Coding 工具 Skill 管理",
}

// skillInitCmd 是 skill init 子命令
var skillInitCmd = &cobra.Command{
	Use:   "init",
	Short: "为 AI Coding 工具生成 TAPD CLI 的 SKILL.md",
	Long:  "扫描当前目录，交互式选择 AI Coding 工具，为选中的工具生成 TAPD CLI 的 SKILL.md 指令文件。",
	RunE:  runSkillInit,
}

func init() {
	skillCmd.AddCommand(skillInitCmd)
	rootCmd.AddCommand(skillCmd)
}

// runSkillInit 执行 skill init 的主流程
func runSkillInit(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		output.PrintError(os.Stderr, "skill_init_error",
			"Failed to get current directory: "+err.Error(),
			"Check file system permissions.")
		os.Exit(output.ExitAPIError)
		return nil
	}

	selected := detectExistingTools(cwd)

	selected, cancelled := runInteractiveSelect(os.Stdin, os.Stderr, selected)
	if cancelled {
		fmt.Fprintln(os.Stderr, "已取消。")
		return nil
	}

	// 检查是否有选中的工具
	hasSelected := false
	for _, s := range selected {
		if s {
			hasSelected = true
			break
		}
	}
	if !hasSelected {
		fmt.Fprintln(os.Stderr, "未选择任何工具。")
		return nil
	}

	content := generateSkillContent(rootCmd)

	var generated []string
	for i, tool := range supportedTools {
		if !selected[i] {
			continue
		}
		path, writeErr := writeSkillFile(cwd, tool.Dir, []byte(content))
		if writeErr != nil {
			output.PrintError(os.Stderr, "skill_init_error",
				fmt.Sprintf("Failed to write %s: %s", path, writeErr.Error()),
				"Check file system permissions.")
			os.Exit(output.ExitAPIError)
			return nil
		}
		generated = append(generated, path)
	}

	// 输出结果
	fmt.Fprintf(os.Stderr, "已为 %d 个工具生成 TAPD skill：\n", len(generated))
	for _, p := range generated {
		fmt.Fprintf(os.Stderr, "  %s\n", p)
	}
	return nil
}

// detectExistingTools 检测当前目录下已存在的 AI Coding 工具配置目录
func detectExistingTools(cwd string) []bool {
	selected := make([]bool, len(supportedTools))
	for i, tool := range supportedTools {
		info, err := os.Stat(filepath.Join(cwd, tool.Dir))
		if err == nil && info.IsDir() {
			selected[i] = true
		}
	}
	return selected
}

// runInteractiveSelect 运行交互式多选菜单，返回选中状态和是否取消
func runInteractiveSelect(input *os.File, w *os.File, selected []bool) ([]bool, bool) {
	scanner := bufio.NewScanner(input)
	for {
		printMenu(w, selected)
		fmt.Fprint(w, "> ")
		if !scanner.Scan() {
			return selected, true
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		switch line {
		case "y", "Y":
			return selected, false
		case "q", "Q":
			return selected, true
		case "a", "A":
			for i := range selected {
				selected[i] = true
			}
		case "n", "N":
			for i := range selected {
				selected[i] = false
			}
		default:
			num, err := strconv.Atoi(line)
			if err != nil || num < 1 || num > len(supportedTools) {
				fmt.Fprintf(w, "无效输入：%s\n", line)
				continue
			}
			selected[num-1] = !selected[num-1]
		}
	}
}

// printMenu 打印交互式多选菜单
func printMenu(w *os.File, selected []bool) {
	fmt.Fprintln(w, "请选择要初始化 TAPD skill 的 AI 编程工具：")
	for i, tool := range supportedTools {
		mark := " "
		if selected[i] {
			mark = "x"
		}
		fmt.Fprintf(w, "  [%s] %d. %s (%s/)\n", mark, i+1, tool.Name, tool.Dir)
	}
	fmt.Fprintln(w, "切换选中：输入编号 | 全选：a | 全不选：n | 确认：y | 取消：q")
}

// generateSkillContent 生成 SKILL.md 的完整内容
func generateSkillContent(root *cobra.Command) string {
	cmdRef := buildCommandReferenceMarkdown(root)
	return strings.Replace(skillTemplate, "{{COMMAND_REFERENCE}}", cmdRef, 1)
}

// buildCommandReferenceMarkdown 从命令树动态生成 Markdown 格式的命令参考
func buildCommandReferenceMarkdown(root *cobra.Command) string {
	lines := buildSpecLines(root)
	var b strings.Builder
	lastGroup := ""
	for _, l := range lines {
		if l.group != lastGroup {
			if lastGroup != "" {
				b.WriteString("```\n\n")
			}
			short := getGroupShort(root, l.group)
			if short != "" {
				b.WriteString("### " + l.group + " — " + short + "\n\n")
			} else {
				b.WriteString("### " + l.group + "\n\n")
			}
			b.WriteString("```bash\n")
			lastGroup = l.group
		}
		b.WriteString(l.text + "\n")
	}
	if lastGroup != "" {
		b.WriteString("```\n")
	}
	return b.String()
}

// getGroupShort 获取指定分组命令的 Short 描述
func getGroupShort(root *cobra.Command, group string) string {
	for _, cmd := range root.Commands() {
		if cmd.Name() == group {
			return cmd.Short
		}
	}
	return ""
}

// writeSkillFile 将 SKILL.md 写入指定工具的 skills/tapd/ 目录
func writeSkillFile(cwd, toolDir string, content []byte) (string, error) {
	dir := filepath.Join(cwd, toolDir, "skills", "tapd")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return filepath.Join(dir, "SKILL.md"), err
	}
	path := filepath.Join(dir, "SKILL.md")
	return path, os.WriteFile(path, content, 0644)
}
