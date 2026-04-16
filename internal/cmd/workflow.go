// Package cmd 中的 workflow.go 实现了工作流状态管理命令
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/output"
)

var (
	flagSystem         string
	flagWorkitemTypeID string
)

// workflowCmd 是 workflow 父命令
var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "工作流状态管理",
}

var workflowTransitionsCmd = &cobra.Command{
	Use:   "transitions",
	Short: "获取状态流转规则",
	RunE:  runWorkflowTransitions,
}

var workflowStatusMapCmd = &cobra.Command{
	Use:   "status-map",
	Short: "获取状态中英文映射",
	RunE:  runWorkflowStatusMap,
}

var workflowLastStepsCmd = &cobra.Command{
	Use:   "last-steps",
	Short: "获取结束状态",
	RunE:  runWorkflowLastSteps,
}

func init() {
	workflowTransitionsCmd.Flags().StringVar(&flagSystem, "system", "", "系统名（story|bug，必需）")
	workflowTransitionsCmd.Flags().StringVar(&flagWorkitemTypeID, "workitem-type-id", "", "需求类别 ID（必需）")

	workflowStatusMapCmd.Flags().StringVar(&flagSystem, "system", "", "系统名（story|bug，必需）")
	workflowStatusMapCmd.Flags().StringVar(&flagWorkitemTypeID, "workitem-type-id", "", "需求类别 ID（必需）")

	workflowLastStepsCmd.Flags().StringVar(&flagSystem, "system", "", "系统名（story|bug，必需）")
	workflowLastStepsCmd.Flags().StringVar(&flagWorkitemTypeID, "workitem-type-id", "", "需求类别 ID")

	workflowCmd.AddCommand(workflowTransitionsCmd, workflowStatusMapCmd, workflowLastStepsCmd)
	rootCmd.AddCommand(workflowCmd)
}

func runWorkflowTransitions(cmd *cobra.Command, args []string) error {
	if flagSystem == "" || flagWorkitemTypeID == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"--system and --workitem-type-id are required",
			"Usage: tapd workflow transitions --system <story|bug> --workitem-type-id <id>")
		os.Exit(output.ExitParamError)
		return nil
	}

	params := map[string]string{
		"workspace_id":     flagWorkspaceID,
		"system":           flagSystem,
		"workitem_type_id": flagWorkitemTypeID,
	}

	data, err := apiClient.GetWorkflowTransitions(params)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, data, !flagPretty)
}

func runWorkflowStatusMap(cmd *cobra.Command, args []string) error {
	if flagSystem == "" || flagWorkitemTypeID == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"--system and --workitem-type-id are required",
			"Usage: tapd workflow status-map --system <story|bug> --workitem-type-id <id>")
		os.Exit(output.ExitParamError)
		return nil
	}

	params := map[string]string{
		"workspace_id":     flagWorkspaceID,
		"system":           flagSystem,
		"workitem_type_id": flagWorkitemTypeID,
	}

	data, err := apiClient.GetWorkflowStatusMap(params)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, data, !flagPretty)
}

func runWorkflowLastSteps(cmd *cobra.Command, args []string) error {
	if flagSystem == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"--system is required",
			"Usage: tapd workflow last-steps --system <story|bug>")
		os.Exit(output.ExitParamError)
		return nil
	}

	params := map[string]string{
		"workspace_id": flagWorkspaceID,
		"system":       flagSystem,
	}
	addOptionalParam(params, "workitem_type_id", flagWorkitemTypeID)

	data, err := apiClient.GetWorkflowLastSteps(params)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, data, !flagPretty)
}
