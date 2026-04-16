// Package cmd 中的 commit_msg.go 实现了源码提交关键字查询命令
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/model"
	"github.com/studyzy/tapd-ai-cli/internal/output"
)

var (
	flagCommitMsgObjectID string
	flagCommitMsgType     string
)

// commitMsgCmd 是 commit-msg 父命令
var commitMsgCmd = &cobra.Command{
	Use:   "commit-msg",
	Short: "源码提交关键字管理",
}

var commitMsgGetCmd = &cobra.Command{
	Use:   "get",
	Short: "获取源码提交关键字（用于关联 git commit 到 TAPD 条目）",
	RunE:  runCommitMsgGet,
}

func init() {
	commitMsgGetCmd.Flags().StringVar(&flagCommitMsgObjectID, "object-id", "", "条目 ID（必需）")
	commitMsgGetCmd.Flags().StringVar(&flagCommitMsgType, "type", "story", "条目类型（story|task|bug）")

	commitMsgCmd.AddCommand(commitMsgGetCmd)
	rootCmd.AddCommand(commitMsgCmd)
}

func runCommitMsgGet(cmd *cobra.Command, args []string) error {
	if flagCommitMsgObjectID == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"--object-id is required",
			"Usage: tapd commit-msg get --object-id <id> [--type story|task|bug]")
		os.Exit(output.ExitParamError)
		return nil
	}

	req := &model.GetCommitMsgRequest{
		WorkspaceID: flagWorkspaceID,
		ObjectID:    flagCommitMsgObjectID,
		Type:        flagCommitMsgType,
	}

	data, err := apiClient.GetCommitMsg(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, data, !flagPretty)
}
