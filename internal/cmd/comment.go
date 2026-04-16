// Package cmd 中的 comment.go 实现了评论管理命令
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/model"
	"github.com/studyzy/tapd-ai-cli/internal/output"
)

var (
	flagEntryType     string
	flagEntryID       string
	flagCommentAuthor string
	flagReplyID       string
	flagOrder         string
)

// commentCmd 是 comment 父命令
var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "评论管理",
}

var commentListCmd = &cobra.Command{
	Use:   "list",
	Short: "查询评论列表",
	RunE:  runCommentList,
}

var commentAddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加评论",
	RunE:  runCommentAdd,
}

var commentUpdateCmd = &cobra.Command{
	Use:   "update <comment_id>",
	Short: "更新评论",
	Args:  cobra.ExactArgs(1),
	RunE:  runCommentUpdate,
}

var commentCountCmd = &cobra.Command{
	Use:   "count",
	Short: "查询评论数量",
	RunE:  runCommentCount,
}

func init() {
	// list 子命令
	commentListCmd.Flags().StringVar(&flagEntryType, "entry-type", "", "评论类型（stories|bug|tasks）")
	commentListCmd.Flags().StringVar(&flagEntryID, "entry-id", "", "条目 ID")
	commentListCmd.Flags().StringVar(&flagCommentAuthor, "author", "", "按评论人筛选")
	commentListCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	commentListCmd.Flags().IntVar(&flagPage, "page", 1, "页码")
	commentListCmd.Flags().StringVar(&flagOrder, "order", "", "排序方式（如 created desc）")

	// add 子命令
	commentAddCmd.Flags().StringVar(&flagEntryType, "entry-type", "", "评论类型（stories|bug|tasks，必需）")
	commentAddCmd.Flags().StringVar(&flagEntryID, "entry-id", "", "条目 ID（必需）")
	commentAddCmd.Flags().StringVar(&flagDescription, "description", "", "评论内容（必需）")
	commentAddCmd.Flags().StringVar(&flagCommentAuthor, "author", "", "评论人（可选，默认使用当前登录用户）")
	commentAddCmd.Flags().StringVar(&flagReplyID, "reply-id", "", "回复的评论 ID（可选）")

	// update 子命令
	commentUpdateCmd.Flags().StringVar(&flagDescription, "description", "", "新的评论内容（必需）")

	// count 子命令
	commentCountCmd.Flags().StringVar(&flagEntryType, "entry-type", "", "评论类型（stories|bug|tasks）")
	commentCountCmd.Flags().StringVar(&flagEntryID, "entry-id", "", "条目 ID")

	commentCmd.AddCommand(commentListCmd, commentAddCmd, commentUpdateCmd, commentCountCmd)
	rootCmd.AddCommand(commentCmd)
}

func runCommentList(cmd *cobra.Command, args []string) error {
	req := &model.ListCommentsRequest{
		WorkspaceID: flagWorkspaceID,
		EntryType:   flagEntryType,
		EntryID:     flagEntryID,
		Author:      flagCommentAuthor,
		Order:       flagOrder,
		Limit:       fmt.Sprintf("%d", flagLimit),
		Page:        fmt.Sprintf("%d", flagPage),
	}
	comments, err := apiClient.ListComments(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}

	// 尝试获取总数用于分页信息
	total, _ := apiClient.CountComments(&model.CountCommentsRequest{
		WorkspaceID: flagWorkspaceID,
		EntryType:   flagEntryType,
		EntryID:     flagEntryID,
		Author:      flagCommentAuthor,
	})

	resp := &model.ListResponse{
		Items:   comments,
		Total:   total,
		Page:    flagPage,
		Limit:   flagLimit,
		HasMore: total > flagPage*flagLimit,
	}
	return output.PrintJSON(os.Stdout, resp, !flagPretty)
}

func runCommentAdd(cmd *cobra.Command, args []string) error {
	if flagEntryType == "" || flagEntryID == "" || flagDescription == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"--entry-type, --entry-id and --description are required",
			"Usage: tapd comment add --entry-type <stories|bug|tasks> --entry-id <id> --description <content>")
		os.Exit(output.ExitParamError)
		return nil
	}

	// author 优先使用命令行参数，否则使用当前登录用户昵称
	author := flagCommentAuthor
	if author == "" {
		author = apiClient.Nick
	}
	req := &model.AddCommentRequest{
		WorkspaceID: flagWorkspaceID,
		EntryType:   flagEntryType,
		EntryID:     flagEntryID,
		Description: flagDescription,
		Author:      author,
		ReplyID:     flagReplyID,
	}
	result, err := apiClient.AddComment(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, result, !flagPretty)
}

func runCommentUpdate(cmd *cobra.Command, args []string) error {
	if flagDescription == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"--description is required",
			fmt.Sprintf("Usage: tapd comment update %s --description <content>", args[0]))
		os.Exit(output.ExitParamError)
		return nil
	}

	req := &model.UpdateCommentRequest{
		WorkspaceID:   flagWorkspaceID,
		ID:            args[0],
		Description:   flagDescription,
		ChangeCreator: apiClient.Nick,
	}
	result, err := apiClient.UpdateComment(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, result, !flagPretty)
}

func runCommentCount(cmd *cobra.Command, args []string) error {
	req := &model.CountCommentsRequest{
		WorkspaceID: flagWorkspaceID,
		EntryType:   flagEntryType,
		EntryID:     flagEntryID,
	}
	count, err := apiClient.CountComments(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, &model.CountResponse{Count: count}, !flagPretty)
}
