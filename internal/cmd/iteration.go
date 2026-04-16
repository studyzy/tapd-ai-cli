// Package cmd 中的 iteration.go 实现了迭代管理命令
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/model"
	"github.com/studyzy/tapd-ai-cli/internal/output"
)

var (
	flagStartDate   string
	flagEndDate     string
	flagCreator     string
	flagCurrentUser string
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

var iterationCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建迭代",
	RunE:  runIterationCreate,
}

var iterationUpdateCmd = &cobra.Command{
	Use:   "update <iteration_id>",
	Short: "更新迭代",
	Args:  cobra.ExactArgs(1),
	RunE:  runIterationUpdate,
}

func init() {
	iterationListCmd.Flags().StringVar(&flagStatus, "status", "", "按状态筛选（open/done）")

	iterationCreateCmd.Flags().StringVar(&flagName, "name", "", "迭代标题（必需）")
	iterationCreateCmd.Flags().StringVar(&flagStartDate, "startdate", "", "开始日期（必需，格式：2006-01-02）")
	iterationCreateCmd.Flags().StringVar(&flagEndDate, "enddate", "", "结束日期（必需，格式：2006-01-02）")
	iterationCreateCmd.Flags().StringVar(&flagCreator, "creator", "", "创建人（必需）")
	iterationCreateCmd.Flags().StringVar(&flagDescription, "description", "", "详细描述")
	iterationCreateCmd.Flags().StringVar(&flagStatus, "status", "", "状态（open/done，默认 open）")

	iterationUpdateCmd.Flags().StringVar(&flagCurrentUser, "current-user", "", "变更人（必需）")
	iterationUpdateCmd.Flags().StringVar(&flagName, "name", "", "新标题")
	iterationUpdateCmd.Flags().StringVar(&flagStartDate, "startdate", "", "新开始日期（格式：2006-01-02）")
	iterationUpdateCmd.Flags().StringVar(&flagEndDate, "enddate", "", "新结束日期（格式：2006-01-02）")
	iterationUpdateCmd.Flags().StringVar(&flagDescription, "description", "", "新描述")
	iterationUpdateCmd.Flags().StringVar(&flagStatus, "status", "", "新状态（open/done）")

	iterationCmd.AddCommand(iterationListCmd, iterationCreateCmd, iterationUpdateCmd)
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

func runIterationCreate(cmd *cobra.Command, args []string) error {
	if flagName == "" {
		output.PrintError(os.Stderr, "missing_parameter", "--name is required", "Usage: tapd iteration create --name <title> --startdate <date> --enddate <date> --creator <user>")
		os.Exit(output.ExitParamError)
		return nil
	}
	if flagStartDate == "" {
		output.PrintError(os.Stderr, "missing_parameter", "--startdate is required", "Usage: tapd iteration create --name <title> --startdate <date> --enddate <date> --creator <user>")
		os.Exit(output.ExitParamError)
		return nil
	}
	if flagEndDate == "" {
		output.PrintError(os.Stderr, "missing_parameter", "--enddate is required", "Usage: tapd iteration create --name <title> --startdate <date> --enddate <date> --creator <user>")
		os.Exit(output.ExitParamError)
		return nil
	}
	if flagCreator == "" {
		output.PrintError(os.Stderr, "missing_parameter", "--creator is required", "Usage: tapd iteration create --name <title> --startdate <date> --enddate <date> --creator <user>")
		os.Exit(output.ExitParamError)
		return nil
	}

	params := map[string]string{
		"workspace_id": flagWorkspaceID,
		"name":         flagName,
		"startdate":    flagStartDate,
		"enddate":      flagEndDate,
		"creator":      flagCreator,
	}
	if flagDescription != "" {
		flagDescription = markdownToHTML(flagDescription)
	}
	addOptionalParam(params, "description", flagDescription)
	addOptionalParam(params, "status", flagStatus)

	result, err := apiClient.CreateIteration(params)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, result, !flagPretty)
}

func runIterationUpdate(cmd *cobra.Command, args []string) error {
	if flagCurrentUser == "" {
		output.PrintError(os.Stderr, "missing_parameter", "--current-user is required", "Usage: tapd iteration update <iteration_id> --current-user <user>")
		os.Exit(output.ExitParamError)
		return nil
	}

	args[0] = expandShortID(args[0], flagWorkspaceID)
	params := map[string]string{
		"workspace_id": flagWorkspaceID,
		"id":           args[0],
		"current_user": flagCurrentUser,
	}
	addOptionalParam(params, "name", flagName)
	addOptionalParam(params, "startdate", flagStartDate)
	addOptionalParam(params, "enddate", flagEndDate)
	if flagDescription != "" {
		flagDescription = markdownToHTML(flagDescription)
	}
	addOptionalParam(params, "description", flagDescription)
	addOptionalParam(params, "status", flagStatus)

	result, err := apiClient.UpdateIteration(params)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, result, !flagPretty)
}
