// Package cmd 中的 story.go 实现了需求管理命令
package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/output"
	"github.com/studyzy/tapd-sdk-go/model"
)

var (
	flagStatus      string
	flagOwner       string
	flagIterationID string
	flagLimit       int
	flagPage        int
	flagName        string
	flagDescription string
	flagDescFile    string
	flagPriority    string
	flagParentID    string
)

// storyCmd 是 story 父命令
var storyCmd = &cobra.Command{
	Use:   "story",
	Short: "需求管理",
}

var storyListCmd = &cobra.Command{
	Use:   "list",
	Short: "查询需求列表",
	RunE:  runStoryList,
}

var storyShowCmd = &cobra.Command{
	Use:   "show <story_id>",
	Short: "查看需求详情",
	Args:  cobra.ExactArgs(1),
	RunE:  runStoryShow,
}

var storyCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建需求",
	Long: `创建需求，描述支持三种输入方式：
  1. --description <text>  直接传入描述文本
  2. --file <path>         从本地文件读取描述内容
  3. echo "..." | tapd story create --name <title>  通过 stdin 管道输入`,
	RunE: runStoryCreate,
}

var storyUpdateCmd = &cobra.Command{
	Use:   "update <story_id>",
	Short: "更新需求",
	Long: `更新需求，描述支持三种输入方式：
  1. --description <text>  直接传入描述文本
  2. --file <path>         从本地文件读取描述内容
  3. echo "..." | tapd story update <story_id>  通过 stdin 管道输入`,
	Args: cobra.ExactArgs(1),
	RunE: runStoryUpdate,
}

var storyCountCmd = &cobra.Command{
	Use:   "count",
	Short: "查询需求数量",
	RunE:  runStoryCount,
}

var storyTodoCmd = &cobra.Command{
	Use:   "todo",
	Short: "查询当前用户待办需求",
	RunE:  runStoryTodo,
}

func init() {
	storyListCmd.Flags().StringVar(&flagStatus, "status", "", "按状态筛选（用 workflow status-map 查询可用值）")
	storyListCmd.Flags().StringVar(&flagOwner, "owner", "", "按处理人筛选")
	storyListCmd.Flags().StringVar(&flagIterationID, "iteration-id", "", "按迭代 ID 筛选")
	storyListCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	storyListCmd.Flags().IntVar(&flagPage, "page", 1, "页码")

	storyCreateCmd.Flags().StringVar(&flagName, "name", "", "需求标题（必需）")
	storyCreateCmd.Flags().StringVar(&flagDescription, "description", "", "描述")
	storyCreateCmd.Flags().StringVar(&flagDescFile, "file", "", "从本地文件读取描述内容")
	storyCreateCmd.Flags().StringVar(&flagOwner, "owner", "", "处理人")
	storyCreateCmd.Flags().StringVar(&flagPriority, "priority", "", "优先级（High/Middle/Low/Nice To Have）")
	storyCreateCmd.Flags().StringVar(&flagIterationID, "iteration-id", "", "关联迭代 ID")
	storyCreateCmd.Flags().StringVar(&flagParentID, "parent-id", "", "父需求 ID（创建子需求时使用）")

	storyUpdateCmd.Flags().StringVar(&flagName, "name", "", "新标题")
	storyUpdateCmd.Flags().StringVar(&flagDescription, "description", "", "新描述")
	storyUpdateCmd.Flags().StringVar(&flagDescFile, "file", "", "从本地文件读取描述内容")
	storyUpdateCmd.Flags().StringVar(&flagStatus, "status", "", "新状态（用 workflow status-map 查询可用值）")
	storyUpdateCmd.Flags().StringVar(&flagOwner, "owner", "", "新处理人")
	storyUpdateCmd.Flags().StringVar(&flagPriority, "priority", "", "新优先级（High/Middle/Low/Nice To Have）")

	storyCountCmd.Flags().StringVar(&flagStatus, "status", "", "按状态筛选（用 workflow status-map 查询可用值）")

	storyTodoCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	storyTodoCmd.Flags().IntVar(&flagPage, "page", 1, "页码")

	storyCmd.AddCommand(storyListCmd, storyShowCmd, storyCreateCmd, storyUpdateCmd, storyCountCmd, storyTodoCmd)
	rootCmd.AddCommand(storyCmd)
}

func runStoryList(cmd *cobra.Command, args []string) error {
	req := &model.ListStoriesRequest{
		WorkspaceID: flagWorkspaceID,
		Status:      flagStatus,
		Owner:       flagOwner,
		IterationID: flagIterationID,
		Fields:      "id,name,status,owner,modified",
		Limit:       flagLimit,
		Page:        flagPage,
	}
	stories, err := apiClient.ListStories(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}

	total, _ := apiClient.CountStories(context.Background(), &model.CountStoriesRequest{
		WorkspaceID: flagWorkspaceID,
		Status:      flagStatus,
	})

	resp := &model.ListResponse{
		Items:   stories,
		Total:   total,
		Page:    flagPage,
		Limit:   flagLimit,
		HasMore: total > flagPage*flagLimit,
	}
	return output.PrintJSON(os.Stdout, resp, !flagPretty)
}

func runStoryShow(cmd *cobra.Command, args []string) error {
	story, err := apiClient.GetStory(context.Background(), flagWorkspaceID, args[0])
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	story.Description = htmlToMarkdown(story.Description)
	if err := printDetail(story, "description"); err != nil {
		return err
	}
	printComments(flagWorkspaceID, "stories", args[0])
	return nil
}

func runStoryCreate(cmd *cobra.Command, args []string) error {
	if flagName == "" {
		output.PrintError(os.Stderr, "missing_parameter", "--name is required", "Usage: tapd story create --name <title>")
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
		Description:   description,
		Owner:         flagOwner,
		Creator:       apiClient.GetNick(),
		PriorityLabel: flagPriority,
		IterationID:   flagIterationID,
		ParentID:      flagParentID,
	}
	story, err := apiClient.CreateStory(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return printSuccessResponse(story.ID, story.URL, "")
}

func runStoryUpdate(cmd *cobra.Command, args []string) error {
	description, err := readDescription()
	if err != nil {
		output.PrintError(os.Stderr, "file_error", err.Error(), "Check that the file path is correct and readable")
		os.Exit(output.ExitParamError)
		return nil
	}

	req := &model.UpdateStoryRequest{
		WorkspaceID:   flagWorkspaceID,
		ID:            args[0],
		Name:          flagName,
		Description:   description,
		VStatus:       flagStatus,
		Owner:         flagOwner,
		PriorityLabel: flagPriority,
	}
	story, err := apiClient.UpdateStory(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return printSuccessResponse(story.ID, fmt.Sprintf("%s/%s/prong/stories/view/%s", apiClient.WebURL(), flagWorkspaceID, story.ID), "")
}

func runStoryCount(cmd *cobra.Command, args []string) error {
	req := &model.CountStoriesRequest{
		WorkspaceID: flagWorkspaceID,
		Status:      flagStatus,
	}
	count, err := apiClient.CountStories(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, &model.CountResponse{Count: count}, !flagPretty)
}

func runStoryTodo(cmd *cobra.Command, args []string) error {
	req := &model.GetTodoRequest{
		WorkspaceID: flagWorkspaceID,
		EntityType:  "story",
		Limit:       flagLimit,
		Page:        flagPage,
	}
	stories, err := apiClient.GetTodoStories(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}

	resp := &model.ListResponse{
		Items: stories,
		Page:  flagPage,
		Limit: flagLimit,
	}
	return output.PrintJSON(os.Stdout, resp, !flagPretty)
}

// readDescription 从 --description、--file 或 stdin 读取描述内容
// 优先级：--description > --file > stdin
// TAPD API 的 description 字段期望 HTML 格式，因此自动将 Markdown 转换为 HTML
func readDescription() (string, error) {
	var content string
	if flagDescription != "" {
		content = flagDescription
	} else if flagDescFile != "" {
		data, err := os.ReadFile(flagDescFile)
		if err != nil {
			return "", err
		}
		content = string(data)
	} else {
		// 尝试从 stdin 读取（仅当 stdin 不是终端时）
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return "", err
			}
			content = string(data)
		}
	}
	if content != "" {
		content = markdownToHTML(content)
	}
	return content, nil
}

// addOptionalParam 当值非空时添加到参数 map
func addOptionalParam(params map[string]string, key, value string) {
	if value != "" {
		params[key] = value
	}
}

// addPaginationParams 添加分页参数
func addPaginationParams(params map[string]string, limit, page int) {
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}
	if page > 0 {
		params["page"] = fmt.Sprintf("%d", page)
	}
}
