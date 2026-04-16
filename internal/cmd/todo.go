// Package cmd 中的 todo.go 实现了用户待办事项查询命令
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/model"
	"github.com/studyzy/tapd-ai-cli/internal/output"
)

var (
	flagTodoEntityType string
)

// todoCmd 是 todo 父命令
var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "待办事项管理",
}

var todoListCmd = &cobra.Command{
	Use:   "list",
	Short: "查询当前用户的待办事项",
	RunE:  runTodoList,
}

func init() {
	todoListCmd.Flags().StringVar(&flagTodoEntityType, "entity-type", "story", "待办类型（story|bug|task）")
	todoListCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	todoListCmd.Flags().IntVar(&flagPage, "page", 1, "页码")

	todoCmd.AddCommand(todoListCmd)
	rootCmd.AddCommand(todoCmd)
}

func runTodoList(cmd *cobra.Command, args []string) error {
	req := &model.GetTodoRequest{
		WorkspaceID: flagWorkspaceID,
		EntityType:  flagTodoEntityType,
	}

	data, err := apiClient.GetTodo(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, data, !flagPretty)
}
