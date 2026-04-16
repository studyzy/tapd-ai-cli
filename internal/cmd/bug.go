// Package cmd 中的 bug.go 实现了缺陷管理命令
package cmd

import (
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
	RunE:  runBugCreate,
}

var bugUpdateCmd = &cobra.Command{
	Use:   "update <bug_id>",
	Short: "更新缺陷",
	Args:  cobra.ExactArgs(1),
	RunE:  runBugUpdate,
}

var bugCountCmd = &cobra.Command{
	Use:   "count",
	Short: "查询缺陷数量",
	RunE:  runBugCount,
}

func init() {
	bugListCmd.Flags().StringVar(&flagStatus, "status", "", "按状态筛选")
	bugListCmd.Flags().StringVar(&flagPriority, "priority", "", "按优先级筛选（urgent/high/medium/low/insignificant）")
	bugListCmd.Flags().StringVar(&flagSeverity, "severity", "", "按严重程度筛选（fatal/serious/normal/prompt/advice）")
	bugListCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	bugListCmd.Flags().IntVar(&flagPage, "page", 1, "页码")

	bugCreateCmd.Flags().StringVar(&flagTitle, "title", "", "缺陷标题（必需）")
	bugCreateCmd.Flags().StringVar(&flagDescription, "description", "", "描述")
	bugCreateCmd.Flags().StringVar(&flagPriority, "priority", "", "优先级")
	bugCreateCmd.Flags().StringVar(&flagSeverity, "severity", "", "严重程度")

	bugUpdateCmd.Flags().StringVar(&flagTitle, "title", "", "新标题")
	bugUpdateCmd.Flags().StringVar(&flagStatus, "status", "", "新状态")
	bugUpdateCmd.Flags().StringVar(&flagPriority, "priority", "", "新优先级")
	bugUpdateCmd.Flags().StringVar(&flagSeverity, "severity", "", "新严重程度")

	bugCountCmd.Flags().StringVar(&flagStatus, "status", "", "按状态筛选")

	bugCmd.AddCommand(bugListCmd, bugShowCmd, bugCreateCmd, bugUpdateCmd, bugCountCmd)
	rootCmd.AddCommand(bugCmd)
}

func runBugList(cmd *cobra.Command, args []string) error {
	params := map[string]string{
		"workspace_id": flagWorkspaceID,
	}
	addOptionalParam(params, "status", flagStatus)
	addOptionalParam(params, "priority_label", flagPriority)
	addOptionalParam(params, "severity", flagSeverity)
	addPaginationParams(params, flagLimit, flagPage)

	bugs, err := apiClient.ListBugs(params)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}

	total, _ := apiClient.CountBugs(map[string]string{
		"workspace_id": flagWorkspaceID,
		"status":       flagStatus,
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
	args[0] = expandShortID(args[0], flagWorkspaceID)
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

	params := map[string]string{
		"workspace_id": flagWorkspaceID,
		"title":        flagTitle,
	}
	if flagDescription != "" {
		flagDescription = markdownToHTML(flagDescription)
	}
	addOptionalParam(params, "description", flagDescription)
	addOptionalParam(params, "priority_label", flagPriority)
	addOptionalParam(params, "severity", flagSeverity)

	result, err := apiClient.CreateBug(params)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, result, !flagPretty)
}

func runBugUpdate(cmd *cobra.Command, args []string) error {
	args[0] = expandShortID(args[0], flagWorkspaceID)
	params := map[string]string{
		"workspace_id": flagWorkspaceID,
		"id":           args[0],
	}
	addOptionalParam(params, "title", flagTitle)
	addOptionalParam(params, "v_status", flagStatus)
	addOptionalParam(params, "priority_label", flagPriority)
	addOptionalParam(params, "severity", flagSeverity)

	result, err := apiClient.UpdateBug(params)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, result, !flagPretty)
}

func runBugCount(cmd *cobra.Command, args []string) error {
	params := map[string]string{
		"workspace_id": flagWorkspaceID,
	}
	addOptionalParam(params, "status", flagStatus)

	count, err := apiClient.CountBugs(params)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, &model.CountResponse{Count: count}, !flagPretty)
}
