// Package cmd 中的 comment.go 实现了评论管理命令
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/output"
	"github.com/studyzy/tapd-sdk-go/model"
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
	Long: `添加评论，评论内容支持三种输入方式：
  1. --description <text>  直接传入评论内容
  2. --file <path>         从本地文件读取评论内容
  3. echo "..." | tapd comment add --entry-type <type> --entry-id <id>  通过 stdin 管道输入`,
	RunE: runCommentAdd,
}

var commentUpdateCmd = &cobra.Command{
	Use:   "update <comment_id>",
	Short: "更新评论",
	Long: `更新评论，评论内容支持三种输入方式：
  1. --description <text>  直接传入评论内容
  2. --file <path>         从本地文件读取评论内容
  3. echo "..." | tapd comment update <comment_id>  通过 stdin 管道输入`,
	Args: cobra.ExactArgs(1),
	RunE: runCommentUpdate,
}

var commentCountCmd = &cobra.Command{
	Use:   "count",
	Short: "查询评论数量",
	RunE:  runCommentCount,
}

func init() {
	// list 子命令
	commentListCmd.Flags().StringVar(&flagEntryType, "entry-type", "", "评论类型（stories|bug|bug_remark|tasks）")
	commentListCmd.Flags().StringVar(&flagEntryID, "entry-id", "", "条目 ID")
	commentListCmd.Flags().StringVar(&flagCommentAuthor, "author", "", "按评论人筛选")
	commentListCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	commentListCmd.Flags().IntVar(&flagPage, "page", 1, "页码")
	commentListCmd.Flags().StringVar(&flagOrder, "order", "", "排序方式（如 created desc）")

	// add 子命令
	commentAddCmd.Flags().StringVar(&flagEntryType, "entry-type", "", "评论类型（stories|bug|bug_remark|tasks，必需）")
	commentAddCmd.Flags().StringVar(&flagEntryID, "entry-id", "", "条目 ID（必需）")
	commentAddCmd.Flags().StringVar(&flagDescription, "description", "", "评论内容")
	commentAddCmd.Flags().StringVar(&flagDescFile, "file", "", "从本地文件读取评论内容")
	commentAddCmd.Flags().StringVar(&flagCommentAuthor, "author", "", "评论人（可选，默认使用当前登录用户）")
	commentAddCmd.Flags().StringVar(&flagReplyID, "reply-id", "", "回复的评论 ID（可选）")

	// update 子命令
	commentUpdateCmd.Flags().StringVar(&flagDescription, "description", "", "新的评论内容")
	commentUpdateCmd.Flags().StringVar(&flagDescFile, "file", "", "从本地文件读取评论内容")

	// count 子命令
	commentCountCmd.Flags().StringVar(&flagEntryType, "entry-type", "", "评论类型（stories|bug|bug_remark|tasks）")
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
		Limit:       flagLimit,
		Page:        flagPage,
	}
	comments, err := apiClient.ListComments(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}

	// 尝试获取总数用于分页信息
	total, _ := apiClient.CountComments(context.Background(), &model.CountCommentsRequest{
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
	description, err := readDescription()
	if err != nil {
		output.PrintError(os.Stderr, "file_error", err.Error(), "Check that the file path is correct and readable")
		os.Exit(output.ExitParamError)
		return nil
	}

	if flagEntryType == "" || flagEntryID == "" || description == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"--entry-type, --entry-id and description are required",
			"Usage: tapd comment add --entry-type <stories|bug|bug_remark|tasks> --entry-id <id> --description <content>")
		os.Exit(output.ExitParamError)
		return nil
	}

	// author 优先使用命令行参数，否则使用当前登录用户昵称
	author := flagCommentAuthor
	if author == "" {
		author = ensureNick()
	}
	req := &model.AddCommentRequest{
		WorkspaceID: flagWorkspaceID,
		EntryType:   flagEntryType,
		EntryID:     flagEntryID,
		Description: description,
		Author:      author,
		ReplyID:     flagReplyID,
	}
	result, err := apiClient.AddComment(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return printSuccessResponse(result.ID, "", "")
}

func runCommentUpdate(cmd *cobra.Command, args []string) error {
	description, err := readDescription()
	if err != nil {
		output.PrintError(os.Stderr, "file_error", err.Error(), "Check that the file path is correct and readable")
		os.Exit(output.ExitParamError)
		return nil
	}

	if description == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"description is required",
			fmt.Sprintf("Usage: tapd comment update %s --description <content>", args[0]))
		os.Exit(output.ExitParamError)
		return nil
	}

	req := &model.UpdateCommentRequest{
		WorkspaceID:   flagWorkspaceID,
		ID:            args[0],
		Description:   description,
		ChangeCreator: ensureNick(),
	}
	result, err := apiClient.UpdateComment(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return printSuccessResponse(result.ID, "", "")
}

func runCommentCount(cmd *cobra.Command, args []string) error {
	req := &model.CountCommentsRequest{
		WorkspaceID: flagWorkspaceID,
		EntryType:   flagEntryType,
		EntryID:     flagEntryID,
	}
	count, err := apiClient.CountComments(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, &model.CountResponse{Count: count}, !flagPretty)
}
