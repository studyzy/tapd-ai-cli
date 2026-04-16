// Package cmd 中的 bug.go 实现了缺陷管理命令
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/model"
	"github.com/studyzy/tapd-ai-cli/internal/output"
)

var (
	flagSeverity string
	flagTitle    string
)

// bugCmd 是 bug 父命令
var bugCmd = &cobra.Command{
	Use:   "bug",
	Short: "缺陷管理",
}

var bugListCmd = &cobra.Command{
	Use:   "list",
	Short: "查询缺陷列表",
	RunE:  runBugList,
}

var bugShowCmd = &cobra.Command{
	Use:   "show <bug_id>",
	Short: "查看缺陷详情",
	Args:  cobra.ExactArgs(1),
	RunE:  runBugShow,
}

var bugCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建缺陷",
	Long: `创建缺陷，描述支持三种输入方式：
  1. --description <text>  直接传入描述文本
  2. --file <path>         从本地文件读取描述内容
  3. echo "..." | tapd bug create --title <title>  通过 stdin 管道输入`,
	RunE: runBugCreate,
}

var bugUpdateCmd = &cobra.Command{
	Use:   "update <bug_id>",
	Short: "更新缺陷",
	Long: `更新缺陷，描述支持三种输入方式：
  1. --description <text>  直接传入描述文本
  2. --file <path>         从本地文件读取描述内容
  3. echo "..." | tapd bug update <bug_id>  通过 stdin 管道输入`,
	Args: cobra.ExactArgs(1),
	RunE: runBugUpdate,
}

var bugCountCmd = &cobra.Command{
	Use:   "count",
	Short: "查询缺陷数量",
	RunE:  runBugCount,
}

var bugTodoCmd = &cobra.Command{
	Use:   "todo",
	Short: "查询当前用户待办缺陷",
	RunE:  runBugTodo,
}

func init() {
	bugListCmd.Flags().StringVar(&flagStatus, "status", "", "按状态筛选")
	bugListCmd.Flags().StringVar(&flagPriority, "priority", "", "按优先级筛选（urgent/high/medium/low/insignificant）")
	bugListCmd.Flags().StringVar(&flagSeverity, "severity", "", "按严重程度筛选（fatal/serious/normal/prompt/advice）")
	bugListCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	bugListCmd.Flags().IntVar(&flagPage, "page", 1, "页码")

	bugCreateCmd.Flags().StringVar(&flagTitle, "title", "", "缺陷标题（必需）")
	bugCreateCmd.Flags().StringVar(&flagDescription, "description", "", "描述")
	bugCreateCmd.Flags().StringVar(&flagDescFile, "file", "", "从本地文件读取描述内容")
	bugCreateCmd.Flags().StringVar(&flagPriority, "priority", "", "优先级")
	bugCreateCmd.Flags().StringVar(&flagSeverity, "severity", "", "严重程度")

	bugUpdateCmd.Flags().StringVar(&flagTitle, "title", "", "新标题")
	bugUpdateCmd.Flags().StringVar(&flagDescription, "description", "", "新描述")
	bugUpdateCmd.Flags().StringVar(&flagDescFile, "file", "", "从本地文件读取描述内容")
	bugUpdateCmd.Flags().StringVar(&flagStatus, "status", "", "新状态")
	bugUpdateCmd.Flags().StringVar(&flagPriority, "priority", "", "新优先级")
	bugUpdateCmd.Flags().StringVar(&flagSeverity, "severity", "", "新严重程度")

	bugCountCmd.Flags().StringVar(&flagStatus, "status", "", "按状态筛选")

	bugTodoCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	bugTodoCmd.Flags().IntVar(&flagPage, "page", 1, "页码")

	bugCmd.AddCommand(bugListCmd, bugShowCmd, bugCreateCmd, bugUpdateCmd, bugCountCmd, bugTodoCmd)
	rootCmd.AddCommand(bugCmd)
}

func runBugList(cmd *cobra.Command, args []string) error {
	req := &model.ListBugsRequest{
		WorkspaceID:   flagWorkspaceID,
		PriorityLabel: flagPriority,
		Severity:      flagSeverity,
		Status:        flagStatus,
		Limit:         fmt.Sprintf("%d", flagLimit),
		Page:          fmt.Sprintf("%d", flagPage),
	}
	bugs, err := apiClient.ListBugs(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}

	total, _ := apiClient.CountBugs(&model.CountBugsRequest{
		WorkspaceID: flagWorkspaceID,
		Status:      flagStatus,
	})

	resp := &model.ListResponse{
		Items:   bugs,
		Total:   total,
		Page:    flagPage,
		Limit:   flagLimit,
		HasMore: total > flagPage*flagLimit,
	}
	return output.PrintJSON(os.Stdout, resp, !flagPretty)
}

func runBugShow(cmd *cobra.Command, args []string) error {
	bug, err := apiClient.GetBug(flagWorkspaceID, args[0])
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	if err := printDetail(bug, "description"); err != nil {
		return err
	}
	printComments(flagWorkspaceID, "bug", args[0])
	return nil
}

func runBugCreate(cmd *cobra.Command, args []string) error {
	if flagTitle == "" {
		output.PrintError(os.Stderr, "missing_parameter", "--title is required", "Usage: tapd bug create --title <title>")
		os.Exit(output.ExitParamError)
		return nil
	}

	description, err := readDescription()
	if err != nil {
		output.PrintError(os.Stderr, "file_error", err.Error(), "Check that the file path is correct and readable")
		os.Exit(output.ExitParamError)
		return nil
	}

	req := &model.CreateBugRequest{
		WorkspaceID:   flagWorkspaceID,
		Title:         flagTitle,
		Description:   description,
		Reporter:      apiClient.Nick,
		PriorityLabel: flagPriority,
		Severity:      flagSeverity,
	}
	result, err := apiClient.CreateBug(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, result, !flagPretty)
}

func runBugUpdate(cmd *cobra.Command, args []string) error {
	description, err := readDescription()
	if err != nil {
		output.PrintError(os.Stderr, "file_error", err.Error(), "Check that the file path is correct and readable")
		os.Exit(output.ExitParamError)
		return nil
	}

	req := &model.UpdateBugRequest{
		WorkspaceID:   flagWorkspaceID,
		ID:            args[0],
		Title:         flagTitle,
		Description:   description,
		VStatus:       flagStatus,
		PriorityLabel: flagPriority,
		Severity:      flagSeverity,
	}
	result, err := apiClient.UpdateBug(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, result, !flagPretty)
}

func runBugCount(cmd *cobra.Command, args []string) error {
	req := &model.CountBugsRequest{
		WorkspaceID: flagWorkspaceID,
		Status:      flagStatus,
	}
	count, err := apiClient.CountBugs(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, &model.CountResponse{Count: count}, !flagPretty)
}

func runBugTodo(cmd *cobra.Command, args []string) error {
	req := &model.GetTodoRequest{
		WorkspaceID: flagWorkspaceID,
		EntityType:  "bug",
		Limit:       fmt.Sprintf("%d", flagLimit),
		Page:        fmt.Sprintf("%d", flagPage),
	}
	bugs, err := apiClient.GetTodoBugs(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}

	resp := &model.ListResponse{
		Items: bugs,
		Page:  flagPage,
		Limit: flagLimit,
	}
	return output.PrintJSON(os.Stdout, resp, !flagPretty)
}
