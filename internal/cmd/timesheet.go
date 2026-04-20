// Package cmd 中的 timesheet.go 实现了花费工时管理命令
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
	flagTimesheetEntityType string
	flagTimesheetEntityID   string
	flagTimesheetSpent      string
	flagTimesheetRemain     string
	flagTimesheetDate       string
	flagTimesheetOwner      string
	flagTimesheetMemo       string
)

// timesheetCmd 是 timesheet 父命令
var timesheetCmd = &cobra.Command{
	Use:   "timesheet",
	Short: "花费工时管理",
}

var timesheetListCmd = &cobra.Command{
	Use:   "list",
	Short: "查询花费工时列表",
	RunE:  runTimesheetList,
}

var timesheetAddCmd = &cobra.Command{
	Use:   "add",
	Short: "填写花费工时",
	RunE:  runTimesheetAdd,
}

var timesheetUpdateCmd = &cobra.Command{
	Use:   "update <timesheet_id>",
	Short: "更新花费工时",
	Args:  cobra.ExactArgs(1),
	RunE:  runTimesheetUpdate,
}

func init() {
	// list 子命令
	timesheetListCmd.Flags().StringVar(&flagTimesheetEntityType, "entity-type", "", "实体类型（story|task|bug）")
	timesheetListCmd.Flags().StringVar(&flagTimesheetEntityID, "entity-id", "", "实体 ID")
	timesheetListCmd.Flags().StringVar(&flagTimesheetOwner, "owner", "", "按工时填写人筛选")
	timesheetListCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	timesheetListCmd.Flags().IntVar(&flagPage, "page", 1, "页码")

	// add 子命令
	timesheetAddCmd.Flags().StringVar(&flagTimesheetEntityType, "entity-type", "", "实体类型（story|task|bug，必需）")
	timesheetAddCmd.Flags().StringVar(&flagTimesheetEntityID, "entity-id", "", "实体 ID（必需）")
	timesheetAddCmd.Flags().StringVar(&flagTimesheetSpent, "timespent", "", "花费工时，如 2h 或 0.5d（必需）")
	timesheetAddCmd.Flags().StringVar(&flagTimesheetOwner, "owner", "", "工时填写人（可选，默认当前用户）")
	timesheetAddCmd.Flags().StringVar(&flagTimesheetDate, "spentdate", "", "花费日期，如 2025-01-01（可选）")
	timesheetAddCmd.Flags().StringVar(&flagTimesheetMemo, "memo", "", "备注（可选）")
	timesheetAddCmd.Flags().StringVar(&flagTimesheetRemain, "timeremain", "", "剩余工时（可选）")

	// update 子命令
	timesheetUpdateCmd.Flags().StringVar(&flagTimesheetSpent, "timespent", "", "花费工时")
	timesheetUpdateCmd.Flags().StringVar(&flagTimesheetRemain, "timeremain", "", "剩余工时")
	timesheetUpdateCmd.Flags().StringVar(&flagTimesheetMemo, "memo", "", "备注")

	timesheetCmd.AddCommand(timesheetListCmd, timesheetAddCmd, timesheetUpdateCmd)
	rootCmd.AddCommand(timesheetCmd)
}

func runTimesheetList(cmd *cobra.Command, args []string) error {
	req := &model.ListTimesheetsRequest{
		WorkspaceID: flagWorkspaceID,
		EntityType:  flagTimesheetEntityType,
		EntityID:    flagTimesheetEntityID,
		Owner:       flagTimesheetOwner,
		Limit:       flagLimit,
		Page:        flagPage,
	}
	timesheets, err := apiClient.ListTimesheets(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, timesheets, !flagPretty)
}

func runTimesheetAdd(cmd *cobra.Command, args []string) error {
	if flagTimesheetEntityType == "" || flagTimesheetEntityID == "" || flagTimesheetSpent == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"--entity-type, --entity-id and --timespent are required",
			"Usage: tapd timesheet add --entity-type <story|task|bug> --entity-id <id> --timespent <time>")
		os.Exit(output.ExitParamError)
		return nil
	}

	// owner 优先使用命令行参数，否则使用当前登录用户昵称
	owner := flagTimesheetOwner
	if owner == "" {
		owner = ensureNick()
	}
	req := &model.AddTimesheetRequest{
		WorkspaceID: flagWorkspaceID,
		EntityType:  flagTimesheetEntityType,
		EntityID:    flagTimesheetEntityID,
		Timespent:   flagTimesheetSpent,
		Owner:       owner,
		Timeremain:  flagTimesheetRemain,
		Spentdate:   flagTimesheetDate,
		Memo:        flagTimesheetMemo,
	}
	result, err := apiClient.AddTimesheet(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, result, !flagPretty)
}

func runTimesheetUpdate(cmd *cobra.Command, args []string) error {
	if flagTimesheetSpent == "" && flagTimesheetRemain == "" && flagTimesheetMemo == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"At least one of --timespent, --timeremain, --memo is required",
			fmt.Sprintf("Usage: tapd timesheet update %s --timespent <time>", args[0]))
		os.Exit(output.ExitParamError)
		return nil
	}

	req := &model.UpdateTimesheetRequest{
		WorkspaceID: flagWorkspaceID,
		ID:          args[0],
		Timespent:   flagTimesheetSpent,
		Timeremain:  flagTimesheetRemain,
		Memo:        flagTimesheetMemo,
	}
	result, err := apiClient.UpdateTimesheet(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, result, !flagPretty)
}
