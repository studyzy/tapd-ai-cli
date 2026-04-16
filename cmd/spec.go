// Package cmd 中的 spec.go 实现了 AI 自发现命令，输出 Tool Definition JSON
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/studyzy/tapd-ai-cli/internal/output"
)

// specCmd 输出 OpenAI/Anthropic 兼容的 Tool Definition JSON
var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "输出 Tool Definition JSON（AI 自发现）",
	Long:  "遍历 Cobra 命令树，为每个叶子命令生成 OpenAI/Anthropic 兼容的 Tool Definition，使 AI Agent 可自动发现所有命令。",
	RunE:  runSpec,
}

func init() {
	rootCmd.AddCommand(specCmd)
}

// toolDefinition 表示单个工具定义
type toolDefinition struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  toolParameters `json:"parameters"`
}

// toolParameters 表示工具参数的 JSON Schema
type toolParameters struct {
	Type       string                  `json:"type"`
	Properties map[string]toolProperty `json:"properties,omitempty"`
	Required   []string                `json:"required,omitempty"`
}

// toolProperty 表示单个参数属性
type toolProperty struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Default     interface{} `json:"default,omitempty"`
}

func runSpec(cmd *cobra.Command, args []string) error {
	tools := buildToolDefinitions(rootCmd)
	return output.PrintJSON(os.Stdout, tools, !flagPretty)
}

// buildToolDefinitions 遍历命令树，为每个叶子命令生成工具定义
func buildToolDefinitions(root *cobra.Command) []toolDefinition {
	var tools []toolDefinition
	walkCommands(root, "", &tools)
	return tools
}

// walkCommands 递归遍历命令树
func walkCommands(cmd *cobra.Command, prefix string, tools *[]toolDefinition) {
	for _, child := range cmd.Commands() {
		if child.Hidden || child.Name() == "help" || child.Name() == "completion" {
			continue
		}

		fullName := child.Name()
		if prefix != "" {
			fullName = prefix + "_" + child.Name()
		}

		if child.HasSubCommands() {
			walkCommands(child, fullName, tools)
		} else {
			tool := commandToTool(child, fullName)
			*tools = append(*tools, tool)
		}
	}
}

// commandToTool 将 Cobra 命令转换为工具定义
func commandToTool(cmd *cobra.Command, name string) toolDefinition {
	props := make(map[string]toolProperty)
	var required []string

	// 收集位置参数
	// 从 Use 字段解析位置参数（如 "show <story_id>"）
	if argName := extractArgName(cmd.Use); argName != "" {
		props[argName] = toolProperty{
			Type:        "string",
			Description: argName,
		}
		required = append(required, argName)
	}

	// 收集标志
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden || f.Name == "help" {
			return
		}
		prop := toolProperty{
			Type:        "string",
			Description: f.Usage,
		}
		if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" {
			prop.Default = f.DefValue
		}
		props[f.Name] = prop
	})

	return toolDefinition{
		Name:        "tapd_" + name,
		Description: cmd.Short,
		Parameters: toolParameters{
			Type:       "object",
			Properties: props,
			Required:   required,
		},
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
