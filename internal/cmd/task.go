// Package cmd 中的 task.go 实现了任务管理命令
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/model"
	"github.com/studyzy/tapd-ai-cli/internal/output"
)

var flagStoryID string

// taskCmd 是 task 父命令
var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "任务管理",
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "查询任务列表",
	RunE:  runTaskList,
}

var taskShowCmd = &cobra.Command{
	Use:   "show <task_id>",
	Short: "查看任务详情",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskShow,
}

var taskCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建任务",
	RunE:  runTaskCreate,
}

var taskUpdateCmd = &cobra.Command{
	Use:   "update <task_id>",
	Short: "更新任务",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskUpdate,
}

var taskCountCmd = &cobra.Command{
	Use:   "count",
	Short: "查询任务数量",
	RunE:  runTaskCount,
}

func init() {
	taskListCmd.Flags().StringVar(&flagStatus, "status", "", "按状态筛选（open/progressing/done）")
	taskListCmd.Flags().StringVar(&flagOwner, "owner", "", "按处理人筛选")
	taskListCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	taskListCmd.Flags().IntVar(&flagPage, "page", 1, "页码")

	taskCreateCmd.Flags().StringVar(&flagName, "name", "", "任务标题（必需）")
	taskCreateCmd.Flags().StringVar(&flagDescription, "description", "", "描述")
	taskCreateCmd.Flags().StringVar(&flagOwner, "owner", "", "处理人")
	taskCreateCmd.Flags().StringVar(&flagPriority, "priority", "", "优先级")
	taskCreateCmd.Flags().StringVar(&flagStoryID, "story-id", "", "关联需求 ID")

	taskUpdateCmd.Flags().StringVar(&flagName, "name", "", "新标题")
	taskUpdateCmd.Flags().StringVar(&flagStatus, "status", "", "新状态（open/progressing/done）")
	taskUpdateCmd.Flags().StringVar(&flagOwner, "owner", "", "新处理人")

	taskCountCmd.Flags().StringVar(&flagStatus, "status", "", "按状态筛选")

	taskCmd.AddCommand(taskListCmd, taskShowCmd, taskCreateCmd, taskUpdateCmd, taskCountCmd)
	rootCmd.AddCommand(taskCmd)
}

func runTaskList(cmd *cobra.Command, args []string) error {
	params := map[string]string{
		"workspace_id": flagWorkspaceID,
		"entity_type":  "tasks",
	}
	addOptionalParam(params, "status", flagStatus)
	addOptionalParam(params, "owner", flagOwner)
	addPaginationParams(params, flagLimit, flagPage)
	params["fields"] = "id,name,status,owner,modified"

	tasks, err := apiClient.ListStories(params)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}

	total, _ := apiClient.CountStories(map[string]string{
		"workspace_id": flagWorkspaceID,
		"entity_type":  "tasks",
		"status":       flagStatus,
	})

	resp := &model.ListResponse{
		Items:   tasks,
		Total:   total,
		Page:    flagPage,
		Limit:   flagLimit,
		HasMore: total > flagPage*flagLimit,
	}
	return output.PrintJSON(os.Stdout, resp, !flagPretty)
}

func runTaskShow(cmd *cobra.Command, args []string) error {
	task, err := apiClient.GetStory(flagWorkspaceID, args[0], "tasks")
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	if err := printDetail(task, "description"); err != nil {
		return err
	}
	printComments(flagWorkspaceID, "tasks", args[0])
	return nil
}

func runTaskCreate(cmd *cobra.Command, args []string) error {
	if flagName == "" {
		output.PrintError(os.Stderr, "missing_parameter", "--name is required", "Usage: tapd task create --name <title>")
		os.Exit(output.ExitParamError)
		return nil
	}

	params := map[string]string{
		"workspace_id": flagWorkspaceID,
		"name":         flagName,
	}
	addOptionalParam(params, "description", flagDescription)
	addOptionalParam(params, "owner", flagOwner)
	addOptionalParam(params, "priority_label", flagPriority)
	addOptionalParam(params, "story_id", flagStoryID)

	result, err := apiClient.CreateStory(params, "tasks")
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, result, !flagPretty)
}

func runTaskUpdate(cmd *cobra.Command, args []string) error {
	params := map[string]string{
		"workspace_id": flagWorkspaceID,
		"id":           args[0],
	}
	addOptionalParam(params, "name", flagName)
	addOptionalParam(params, "status", flagStatus)
	addOptionalParam(params, "owner", flagOwner)

	result, err := apiClient.UpdateStory(params, "tasks")
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, result, !flagPretty)
}

func runTaskCount(cmd *cobra.Command, args []string) error {
	params := map[string]string{
		"workspace_id": flagWorkspaceID,
		"entity_type":  "tasks",
	}
	addOptionalParam(params, "status", flagStatus)

	count, err := apiClient.CountStories(params)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, &model.CountResponse{Count: count}, !flagPretty)
}
