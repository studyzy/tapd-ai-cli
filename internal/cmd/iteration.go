// Package cmd 中的 iteration.go 实现了迭代管理命令
package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/output"
	"github.com/studyzy/tapd-sdk-go/model"
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

var iterationCountCmd = &cobra.Command{
	Use:   "count",
	Short: "查询迭代数量",
	RunE:  runIterationCount,
}

func init() {
	iterationListCmd.Flags().StringVar(&flagStatus, "status", "", "按状态筛选（open/done）")
	iterationListCmd.Flags().StringVar(&flagName, "name", "", "按标题模糊匹配")
	iterationListCmd.Flags().StringVar(&flagCreator, "creator", "", "按创建人筛选")
	iterationListCmd.Flags().StringVar(&flagOrder, "order", "", "排序规则（如 \"created desc\"）")
	iterationListCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	iterationListCmd.Flags().IntVar(&flagPage, "page", 1, "页码")

	iterationCreateCmd.Flags().StringVar(&flagName, "name", "", "迭代标题（必需）")
	iterationCreateCmd.Flags().StringVar(&flagStartDate, "startdate", "", "开始日期（必需，格式：2006-01-02）")
	iterationCreateCmd.Flags().StringVar(&flagEndDate, "enddate", "", "结束日期（必需，格式：2006-01-02）")
	iterationCreateCmd.Flags().StringVar(&flagCreator, "creator", "", "创建人（必需）")
	iterationCreateCmd.Flags().StringVar(&flagDescription, "description", "", "详细描述")
	iterationCreateCmd.Flags().StringVar(&flagStatus, "status", "", "状态（open/done，默认 open）")
	iterationCreateCmd.Flags().StringVar(&flagLabel, "label", "", "标签（多个以竖线分隔）")
	iterationCreateCmd.Flags().StringVar(&flagParentID, "parent-id", "", "上层计划 ID")

	iterationUpdateCmd.Flags().StringVar(&flagCurrentUser, "current-user", "", "变更人（必需）")
	iterationUpdateCmd.Flags().StringVar(&flagName, "name", "", "新标题")
	iterationUpdateCmd.Flags().StringVar(&flagStartDate, "startdate", "", "新开始日期（格式：2006-01-02）")
	iterationUpdateCmd.Flags().StringVar(&flagEndDate, "enddate", "", "新结束日期（格式：2006-01-02）")
	iterationUpdateCmd.Flags().StringVar(&flagDescription, "description", "", "新描述")
	iterationUpdateCmd.Flags().StringVar(&flagStatus, "status", "", "新状态（open/done）")

	iterationCountCmd.Flags().StringVar(&flagStatus, "status", "", "按状态筛选（open/done）")

	iterationCmd.AddCommand(iterationListCmd, iterationCreateCmd, iterationUpdateCmd, iterationCountCmd)
	rootCmd.AddCommand(iterationCmd)
}

func runIterationList(cmd *cobra.Command, args []string) error {
	req := &model.ListIterationsRequest{
		WorkspaceID: flagWorkspaceID,
		Name:        flagName,
		Status:      flagStatus,
		Creator:     flagCreator,
		Order:       flagOrder,
		Fields:      "id,name,status,startdate,enddate,modified",
		Limit:       flagLimit,
		Page:        flagPage,
	}
	iterations, err := apiClient.ListIterations(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}

	total, _ := apiClient.CountIterations(context.Background(), &model.CountIterationsRequest{
		WorkspaceID: flagWorkspaceID,
		Status:      flagStatus,
	})

	resp := &model.ListResponse{
		Items:   iterations,
		Total:   total,
		Page:    flagPage,
		Limit:   flagLimit,
		HasMore: total > flagPage*flagLimit,
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

	req := &model.CreateIterationRequest{
		WorkspaceID: flagWorkspaceID,
		Name:        flagName,
		StartDate:   flagStartDate,
		EndDate:     flagEndDate,
		Creator:     flagCreator,
		Description: flagDescription,
		Status:      flagStatus,
		Label:       flagLabel,
		ParentID:    flagParentID,
	}
	iteration, err := apiClient.CreateIteration(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return printSuccessResponse(iteration.ID, "", iteration.WorkspaceID)
}

func runIterationUpdate(cmd *cobra.Command, args []string) error {
	if flagCurrentUser == "" {
		output.PrintError(os.Stderr, "missing_parameter", "--current-user is required", "Usage: tapd iteration update <iteration_id> --current-user <user>")
		os.Exit(output.ExitParamError)
		return nil
	}

	req := &model.UpdateIterationRequest{
		WorkspaceID: flagWorkspaceID,
		ID:          args[0],
		CurrentUser: flagCurrentUser,
		Name:        flagName,
		StartDate:   flagStartDate,
		EndDate:     flagEndDate,
		Description: flagDescription,
		Status:      flagStatus,
	}
	iteration, err := apiClient.UpdateIteration(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return printSuccessResponse(iteration.ID, "", iteration.WorkspaceID)
}

func runIterationCount(cmd *cobra.Command, args []string) error {
	req := &model.CountIterationsRequest{
		WorkspaceID: flagWorkspaceID,
		Status:      flagStatus,
	}
	count, err := apiClient.CountIterations(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, &model.CountResponse{Count: count}, !flagPretty)
}
