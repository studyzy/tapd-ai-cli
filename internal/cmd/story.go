// Package cmd 中的 story.go 实现了需求管理命令
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/model"
	"github.com/studyzy/tapd-ai-cli/internal/output"
)

var (
	flagStatus      string
	flagOwner       string
	flagIterationID string
	flagLimit       int
	flagPage        int
	flagName        string
	flagDescription string
	flagPriority    string
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
	RunE:  runStoryCreate,
}

var storyUpdateCmd = &cobra.Command{
	Use:   "update <story_id>",
	Short: "更新需求",
	Args:  cobra.ExactArgs(1),
	RunE:  runStoryUpdate,
}

var storyCountCmd = &cobra.Command{
	Use:   "count",
	Short: "查询需求数量",
	RunE:  runStoryCount,
}

func init() {
	storyListCmd.Flags().StringVar(&flagStatus, "status", "", "按状态筛选")
	storyListCmd.Flags().StringVar(&flagOwner, "owner", "", "按处理人筛选")
	storyListCmd.Flags().StringVar(&flagIterationID, "iteration-id", "", "按迭代 ID 筛选")
	storyListCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	storyListCmd.Flags().IntVar(&flagPage, "page", 1, "页码")

	storyCreateCmd.Flags().StringVar(&flagName, "name", "", "需求标题（必需）")
	storyCreateCmd.Flags().StringVar(&flagDescription, "description", "", "描述")
	storyCreateCmd.Flags().StringVar(&flagOwner, "owner", "", "处理人")
	storyCreateCmd.Flags().StringVar(&flagPriority, "priority", "", "优先级（High/Middle/Low/Nice To Have）")
	storyCreateCmd.Flags().StringVar(&flagIterationID, "iteration-id", "", "关联迭代 ID")

	storyUpdateCmd.Flags().StringVar(&flagName, "name", "", "新标题")
	storyUpdateCmd.Flags().StringVar(&flagStatus, "status", "", "新状态")
	storyUpdateCmd.Flags().StringVar(&flagOwner, "owner", "", "新处理人")
	storyUpdateCmd.Flags().StringVar(&flagPriority, "priority", "", "新优先级")

	storyCountCmd.Flags().StringVar(&flagStatus, "status", "", "按状态筛选")

	storyCmd.AddCommand(storyListCmd, storyShowCmd, storyCreateCmd, storyUpdateCmd, storyCountCmd)
	rootCmd.AddCommand(storyCmd)
}

func runStoryList(cmd *cobra.Command, args []string) error {
	params := map[string]string{
		"workspace_id": flagWorkspaceID,
		"entity_type":  "stories",
	}
	addOptionalParam(params, "status", flagStatus)
	addOptionalParam(params, "owner", flagOwner)
	addOptionalParam(params, "iteration_id", flagIterationID)
	addPaginationParams(params, flagLimit, flagPage)
	params["fields"] = "id,name,status,owner,modified"

	stories, err := apiClient.ListStories(params)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}

	total, _ := apiClient.CountStories(map[string]string{
		"workspace_id": flagWorkspaceID,
		"entity_type":  "stories",
		"status":       flagStatus,
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
	args[0] = expandShortID(args[0], flagWorkspaceID)
	story, err := apiClient.GetStory(flagWorkspaceID, args[0], "stories")
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
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

	params := map[string]string{
		"workspace_id": flagWorkspaceID,
		"name":         flagName,
	}
	if flagDescription != "" {
		flagDescription = markdownToHTML(flagDescription)
	}
	addOptionalParam(params, "description", flagDescription)
	addOptionalParam(params, "owner", flagOwner)
	addOptionalParam(params, "priority_label", flagPriority)
	addOptionalParam(params, "iteration_id", flagIterationID)

	result, err := apiClient.CreateStory(params, "stories")
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, result, !flagPretty)
}

func runStoryUpdate(cmd *cobra.Command, args []string) error {
	args[0] = expandShortID(args[0], flagWorkspaceID)
	params := map[string]string{
		"workspace_id": flagWorkspaceID,
		"id":           args[0],
	}
	addOptionalParam(params, "name", flagName)
	addOptionalParam(params, "v_status", flagStatus)
	addOptionalParam(params, "owner", flagOwner)
	addOptionalParam(params, "priority_label", flagPriority)

	result, err := apiClient.UpdateStory(params, "stories")
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, result, !flagPretty)
}

func runStoryCount(cmd *cobra.Command, args []string) error {
	params := map[string]string{
		"workspace_id": flagWorkspaceID,
		"entity_type":  "stories",
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
