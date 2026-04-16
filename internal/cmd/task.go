// Package cmd 中的 task.go 实现了任务管理命令
package cmd

import (
	"fmt"
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
	Long: `创建任务，描述支持三种输入方式：
  1. --description <text>  直接传入描述文本
  2. --file <path>         从本地文件读取描述内容
  3. echo "..." | tapd task create --name <title>  通过 stdin 管道输入`,
	RunE: runTaskCreate,
}

var taskUpdateCmd = &cobra.Command{
	Use:   "update <task_id>",
	Short: "更新任务",
	Long: `更新任务，描述支持三种输入方式：
  1. --description <text>  直接传入描述文本
  2. --file <path>         从本地文件读取描述内容
  3. echo "..." | tapd task update <task_id>  通过 stdin 管道输入`,
	Args: cobra.ExactArgs(1),
	RunE: runTaskUpdate,
}

var taskCountCmd = &cobra.Command{
	Use:   "count",
	Short: "查询任务数量",
	RunE:  runTaskCount,
}

var taskTodoCmd = &cobra.Command{
	Use:   "todo",
	Short: "查询当前用户待办任务",
	RunE:  runTaskTodo,
}

func init() {
	taskListCmd.Flags().StringVar(&flagStatus, "status", "", "按状态筛选（open/progressing/done）")
	taskListCmd.Flags().StringVar(&flagOwner, "owner", "", "按处理人筛选")
	taskListCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	taskListCmd.Flags().IntVar(&flagPage, "page", 1, "页码")

	taskCreateCmd.Flags().StringVar(&flagName, "name", "", "任务标题（必需）")
	taskCreateCmd.Flags().StringVar(&flagDescription, "description", "", "描述")
	taskCreateCmd.Flags().StringVar(&flagDescFile, "file", "", "从本地文件读取描述内容")
	taskCreateCmd.Flags().StringVar(&flagOwner, "owner", "", "处理人")
	taskCreateCmd.Flags().StringVar(&flagPriority, "priority", "", "优先级")
	taskCreateCmd.Flags().StringVar(&flagStoryID, "story-id", "", "关联需求 ID")

	taskUpdateCmd.Flags().StringVar(&flagName, "name", "", "新标题")
	taskUpdateCmd.Flags().StringVar(&flagDescription, "description", "", "新描述")
	taskUpdateCmd.Flags().StringVar(&flagDescFile, "file", "", "从本地文件读取描述内容")
	taskUpdateCmd.Flags().StringVar(&flagStatus, "status", "", "新状态（open/progressing/done）")
	taskUpdateCmd.Flags().StringVar(&flagOwner, "owner", "", "新处理人")

	taskCountCmd.Flags().StringVar(&flagStatus, "status", "", "按状态筛选")

	taskTodoCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	taskTodoCmd.Flags().IntVar(&flagPage, "page", 1, "页码")

	taskCmd.AddCommand(taskListCmd, taskShowCmd, taskCreateCmd, taskUpdateCmd, taskCountCmd, taskTodoCmd)
	rootCmd.AddCommand(taskCmd)
}

func runTaskList(cmd *cobra.Command, args []string) error {
	req := &model.ListStoriesRequest{
		WorkspaceID: flagWorkspaceID,
		EntityType:  "tasks",
		Status:      flagStatus,
		Owner:       flagOwner,
		Fields:      "id,name,status,owner,modified",
		Limit:       fmt.Sprintf("%d", flagLimit),
		Page:        fmt.Sprintf("%d", flagPage),
	}
	tasks, err := apiClient.ListStories(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}

	total, _ := apiClient.CountStories(&model.CountStoriesRequest{
		WorkspaceID: flagWorkspaceID,
		EntityType:  "tasks",
		Status:      flagStatus,
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

	description, err := readDescription()
	if err != nil {
		output.PrintError(os.Stderr, "file_error", err.Error(), "Check that the file path is correct and readable")
		os.Exit(output.ExitParamError)
		return nil
	}

	req := &model.CreateStoryRequest{
		WorkspaceID:   flagWorkspaceID,
		Name:          flagName,
		EntityType:    "tasks",
		Description:   description,
		Owner:         flagOwner,
		Creator:       apiClient.Nick,
		PriorityLabel: flagPriority,
		StoryID:       flagStoryID,
	}
	result, err := apiClient.CreateStory(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, result, !flagPretty)
}

func runTaskUpdate(cmd *cobra.Command, args []string) error {
	description, err := readDescription()
	if err != nil {
		output.PrintError(os.Stderr, "file_error", err.Error(), "Check that the file path is correct and readable")
		os.Exit(output.ExitParamError)
		return nil
	}

	req := &model.UpdateStoryRequest{
		WorkspaceID: flagWorkspaceID,
		ID:          args[0],
		EntityType:  "tasks",
		Name:        flagName,
		Description: description,
		Status:      flagStatus,
		Owner:       flagOwner,
	}
	result, err := apiClient.UpdateStory(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, result, !flagPretty)
}

func runTaskCount(cmd *cobra.Command, args []string) error {
	req := &model.CountStoriesRequest{
		WorkspaceID: flagWorkspaceID,
		EntityType:  "tasks",
		Status:      flagStatus,
	}
	count, err := apiClient.CountStories(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, &model.CountResponse{Count: count}, !flagPretty)
}

func runTaskTodo(cmd *cobra.Command, args []string) error {
	req := &model.GetTodoRequest{
		WorkspaceID: flagWorkspaceID,
		EntityType:  "task",
		Limit:       fmt.Sprintf("%d", flagLimit),
		Page:        fmt.Sprintf("%d", flagPage),
	}
	tasks, err := apiClient.GetTodoTasks(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}

	resp := &model.ListResponse{
		Items: tasks,
		Page:  flagPage,
		Limit: flagLimit,
	}
	return output.PrintJSON(os.Stdout, resp, !flagPretty)
}
