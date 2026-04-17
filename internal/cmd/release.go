// Package cmd 中的 release.go 实现了发布计划管理命令
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/output"
	"github.com/studyzy/tapd-sdk-go/model"
)

var (
	flagReleaseStatus string
)

// releaseCmd 是 release 父命令
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "发布计划管理",
}

var releaseListCmd = &cobra.Command{
	Use:   "list",
	Short: "查询发布计划列表",
	RunE:  runReleaseList,
}

func init() {
	releaseListCmd.Flags().StringVar(&flagName, "name", "", "按名称筛选")
	releaseListCmd.Flags().StringVar(&flagReleaseStatus, "status", "", "按状态筛选（open/done）")
	releaseListCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	releaseListCmd.Flags().IntVar(&flagPage, "page", 1, "页码")

	releaseCmd.AddCommand(releaseListCmd)
	rootCmd.AddCommand(releaseCmd)
}

func runReleaseList(cmd *cobra.Command, args []string) error {
	req := &model.WorkspaceIDRequest{
		WorkspaceID: flagWorkspaceID,
	}

	releases, err := apiClient.ListReleases(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, releases, !flagPretty)
}
