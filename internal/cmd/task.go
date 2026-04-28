// Package cmd 中的 task.go 实现了任务管理命令
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/output"
	"github.com/studyzy/tapd-sdk-go/model"
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
	taskListCmd.Flags().StringVar(&flagName, "name", "", "按标题模糊匹配")
	taskListCmd.Flags().StringVar(&flagStoryID, "story-id", "", "按关联需求 ID 筛选")
	taskListCmd.Flags().StringVar(&flagIterationID, "iteration-id", "", "按迭代 ID 筛选")
	taskListCmd.Flags().StringVar(&flagPriority, "priority", "", "按优先级筛选")
	taskListCmd.Flags().StringVar(&flagLabel, "label", "", "按标签筛选")
	taskListCmd.Flags().StringVar(&flagOrder, "order", "", "排序规则（如 \"created desc\"）")
	taskListCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	taskListCmd.Flags().IntVar(&flagPage, "page", 1, "页码")

	taskCreateCmd.Flags().StringVar(&flagName, "name", "", "任务标题（必需）")
	taskCreateCmd.Flags().StringVar(&flagDescription, "description", "", "描述")
	taskCreateCmd.Flags().StringVar(&flagDescFile, "file", "", "从本地文件读取描述内容")
	taskCreateCmd.Flags().StringVar(&flagOwner, "owner", "", "处理人")
	taskCreateCmd.Flags().StringVar(&flagPriority, "priority", "", "优先级（High/Middle/Low/Nice To Have）")
	taskCreateCmd.Flags().StringVar(&flagStoryID, "story-id", "", "关联需求 ID")
	taskCreateCmd.Flags().StringVar(&flagCC, "cc", "", "抄送人")
	taskCreateCmd.Flags().StringVar(&flagBegin, "begin", "", "预计开始日期（格式：2006-01-02）")
	taskCreateCmd.Flags().StringVar(&flagDue, "due", "", "预计结束日期（格式：2006-01-02）")
	taskCreateCmd.Flags().StringVar(&flagIterationID, "iteration-id", "", "关联迭代 ID")
	taskCreateCmd.Flags().StringVar(&flagEffort, "effort", "", "预估工时")
	taskCreateCmd.Flags().StringVar(&flagLabel, "label", "", "标签（多个以竖线分隔）")
	taskCreateCmd.Flags().StringArrayVar(&flagCustomField, "custom-field", nil, "自定义字段（可重复，格式：key=value）")

	taskUpdateCmd.Flags().StringVar(&flagName, "name", "", "新标题")
	taskUpdateCmd.Flags().StringVar(&flagDescription, "description", "", "新描述")
	taskUpdateCmd.Flags().StringVar(&flagDescFile, "file", "", "从本地文件读取描述内容")
	taskUpdateCmd.Flags().StringVar(&flagStatus, "status", "", "新状态（open/progressing/done）")
	taskUpdateCmd.Flags().StringVar(&flagOwner, "owner", "", "新处理人")
	taskUpdateCmd.Flags().StringVar(&flagCC, "cc", "", "新抄送人")
	taskUpdateCmd.Flags().StringVar(&flagBegin, "begin", "", "新预计开始日期（格式：2006-01-02）")
	taskUpdateCmd.Flags().StringVar(&flagDue, "due", "", "新预计结束日期（格式：2006-01-02）")
	taskUpdateCmd.Flags().StringVar(&flagStoryID, "story-id", "", "新关联需求 ID")
	taskUpdateCmd.Flags().StringVar(&flagIterationID, "iteration-id", "", "新迭代 ID")
	taskUpdateCmd.Flags().StringVar(&flagPriority, "priority", "", "新优先级（High/Middle/Low/Nice To Have）")
	taskUpdateCmd.Flags().StringVar(&flagEffort, "effort", "", "新预估工时")
	taskUpdateCmd.Flags().StringVar(&flagLabel, "label", "", "新标签（多个以竖线分隔）")
	taskUpdateCmd.Flags().StringVar(&flagCurrentUser, "current-user", "", "操作人")
	taskUpdateCmd.Flags().StringArrayVar(&flagCustomField, "custom-field", nil, "自定义字段（可重复，格式：key=value）")

	taskCountCmd.Flags().StringVar(&flagStatus, "status", "", "按状态筛选（open/progressing/done）")

	taskTodoCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	taskTodoCmd.Flags().IntVar(&flagPage, "page", 1, "页码")

	taskCmd.AddCommand(taskListCmd, taskShowCmd, taskCreateCmd, taskUpdateCmd, taskCountCmd, taskTodoCmd)
	rootCmd.AddCommand(taskCmd)
}

func runTaskList(cmd *cobra.Command, args []string) error {
	req := &model.ListTasksRequest{
		WorkspaceID:   flagWorkspaceID,
		Name:          flagName,
		Status:        flagStatus,
		Owner:         flagOwner,
		StoryID:       flagStoryID,
		IterationID:   flagIterationID,
		PriorityLabel: flagPriority,
		Label:         flagLabel,
		Order:         flagOrder,
		Fields:        "id,name,status,owner,modified",
		Limit:         flagLimit,
		Page:          flagPage,
	}
	tasks, err := apiClient.ListTasks(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}

	total, _ := apiClient.CountTasks(context.Background(), &model.CountTasksRequest{
		WorkspaceID: flagWorkspaceID,
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
	task, err := apiClient.GetTask(context.Background(), flagWorkspaceID, args[0])
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	task.Description = htmlToMarkdown(task.Description)
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

	req := &model.CreateTaskRequest{
		WorkspaceID:   flagWorkspaceID,
		Name:          flagName,
		Description:   description,
		Owner:         flagOwner,
		Creator:       ensureNick(),
		PriorityLabel: flagPriority,
		StoryID:       flagStoryID,
		CC:            flagCC,
		Begin:         flagBegin,
		Due:           flagDue,
		IterationID:   flagIterationID,
		Effort:        flagEffort,
		Label:         flagLabel,
		CustomFields:  parseCustomFields(flagCustomField),
	}
	task, err := apiClient.CreateTask(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return printSuccessResponse(task.ID, task.URL, "")
}

func runTaskUpdate(cmd *cobra.Command, args []string) error {
	description, err := readDescription()
	if err != nil {
		output.PrintError(os.Stderr, "file_error", err.Error(), "Check that the file path is correct and readable")
		os.Exit(output.ExitParamError)
		return nil
	}

	req := &model.UpdateTaskRequest{
		WorkspaceID:   flagWorkspaceID,
		ID:            args[0],
		Name:          flagName,
		Description:   description,
		Status:        flagStatus,
		Owner:         flagOwner,
		CC:            flagCC,
		Begin:         flagBegin,
		Due:           flagDue,
		StoryID:       flagStoryID,
		IterationID:   flagIterationID,
		PriorityLabel: flagPriority,
		Effort:        flagEffort,
		Label:         flagLabel,
		CurrentUser:   flagCurrentUser,
		CustomFields:  parseCustomFields(flagCustomField),
	}
	task, err := apiClient.UpdateTask(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return printSuccessResponse(task.ID, fmt.Sprintf("%s/%s/prong/tasks/view/%s", apiClient.WebURL(), flagWorkspaceID, task.ID), "")
}

func runTaskCount(cmd *cobra.Command, args []string) error {
	req := &model.CountTasksRequest{
		WorkspaceID: flagWorkspaceID,
		Status:      flagStatus,
	}
	count, err := apiClient.CountTasks(context.Background(), req)
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
		Limit:       flagLimit,
		Page:        flagPage,
	}
	tasks, err := apiClient.GetTodoTasks(context.Background(), req)
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
