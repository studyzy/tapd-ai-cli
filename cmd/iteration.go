// Package cmd 中的 iteration.go 实现了迭代查询命令
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/model"
	"github.com/studyzy/tapd-ai-cli/internal/output"
)

// iterationCmd 是 iteration 父命令
var iterationCmd = &cobra.Command{
	Use:   "iteration",
	Short: "迭代管理",
}

var iterationListCmd = &cobra.Command{
	Use:   "list",
	Short: "查询迭代列表",
	RunE:  runIterationList,
}

func init() {
	iterationListCmd.Flags().StringVar(&flagStatus, "status", "", "按状态筛选（open/done）")
	iterationCmd.AddCommand(iterationListCmd)
	rootCmd.AddCommand(iterationCmd)
}

func runIterationList(cmd *cobra.Command, args []string) error {
	params := map[string]string{
		"workspace_id": flagWorkspaceID,
	}
	addOptionalParam(params, "status", flagStatus)

	iterations, err := apiClient.ListIterations(params)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}

	resp := &model.ListResponse{
		Items: iterations,
		Total: len(iterations),
	}
	return output.PrintJSON(os.Stdout, resp, !flagPretty)
}
