// Package cmd 中的 tcase.go 实现了测试用例管理命令
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/model"
	"github.com/studyzy/tapd-ai-cli/internal/output"
)

var (
	flagTCaseID           string
	flagTCasePrecondition string
	flagTCaseSteps        string
	flagTCaseExpectation  string
	flagTCaseType         string
	flagTCaseCreator      string
	flagTCasesJSON        string
)

// tcaseCmd 是 tcase 父命令
var tcaseCmd = &cobra.Command{
	Use:   "tcase",
	Short: "测试用例管理",
}

var tcaseListCmd = &cobra.Command{
	Use:   "list",
	Short: "查询测试用例列表",
	RunE:  runTCaseList,
}

var tcaseCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建或更新测试用例",
	RunE:  runTCaseCreate,
}

var tcaseBatchCreateCmd = &cobra.Command{
	Use:   "batch-create",
	Short: "批量创建测试用例",
	RunE:  runTCaseBatchCreate,
}

func init() {
	tcaseListCmd.Flags().StringVar(&flagStatus, "status", "", "按状态筛选（updating|abandon|normal）")
	tcaseListCmd.Flags().StringVar(&flagPriority, "priority", "", "按等级筛选")
	tcaseListCmd.Flags().StringVar(&flagTCaseCreator, "creator", "", "按创建人筛选")
	tcaseListCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	tcaseListCmd.Flags().IntVar(&flagPage, "page", 1, "页码")

	tcaseCreateCmd.Flags().StringVar(&flagTCaseID, "id", "", "测试用例 ID（有值时为更新，无值时为创建）")
	tcaseCreateCmd.Flags().StringVar(&flagName, "name", "", "用例名称（创建时必需）")
	tcaseCreateCmd.Flags().StringVar(&flagStatus, "status", "", "状态（updating|abandon|normal）")
	tcaseCreateCmd.Flags().StringVar(&flagTCasePrecondition, "precondition", "", "前置条件")
	tcaseCreateCmd.Flags().StringVar(&flagTCaseSteps, "steps", "", "用例步骤")
	tcaseCreateCmd.Flags().StringVar(&flagTCaseExpectation, "expectation", "", "预期结果")
	tcaseCreateCmd.Flags().StringVar(&flagTCaseType, "type", "", "用例类型")
	tcaseCreateCmd.Flags().StringVar(&flagPriority, "priority", "", "用例等级")
	tcaseCreateCmd.Flags().StringVar(&flagTCaseCreator, "creator", "", "创建人")

	tcaseBatchCreateCmd.Flags().StringVar(&flagTCasesJSON, "tcases", "", "测试用例 JSON 数组（必需）")

	tcaseCmd.AddCommand(tcaseListCmd, tcaseCreateCmd, tcaseBatchCreateCmd)
	rootCmd.AddCommand(tcaseCmd)
}

func runTCaseList(cmd *cobra.Command, args []string) error {
	req := &model.ListTCasesRequest{
		WorkspaceID: flagWorkspaceID,
		Status:      flagStatus,
		Priority:    flagPriority,
		Limit:       fmt.Sprintf("%d", flagLimit),
		Page:        fmt.Sprintf("%d", flagPage),
	}

	tcases, err := apiClient.ListTCases(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}

	countReq := &model.CountTCasesRequest{
		WorkspaceID: flagWorkspaceID,
		Status:      flagStatus,
	}
	total, _ := apiClient.CountTCases(countReq)

	resp := &model.ListResponse{
		Items:   tcases,
		Total:   total,
		Page:    flagPage,
		Limit:   flagLimit,
		HasMore: total > flagPage*flagLimit,
	}
	return output.PrintJSON(os.Stdout, resp, !flagPretty)
}

func runTCaseCreate(cmd *cobra.Command, args []string) error {
	if flagTCaseID == "" && flagName == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"--name is required for creating a test case",
			"Usage: tapd tcase create --name <name>")
		os.Exit(output.ExitParamError)
		return nil
	}

	req := &model.CreateTCaseRequest{
		WorkspaceID:  flagWorkspaceID,
		Name:         flagName,
		Precondition: flagTCasePrecondition,
		Steps:        flagTCaseSteps,
		Expectation:  flagTCaseExpectation,
		Type:         flagTCaseType,
		Priority:     flagPriority,
		Creator:      flagTCaseCreator,
	}

	result, err := apiClient.CreateTCase(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, result, !flagPretty)
}

func runTCaseBatchCreate(cmd *cobra.Command, args []string) error {
	if flagTCasesJSON == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"--tcases is required",
			"Usage: tapd tcase batch-create --tcases '[{\"name\":\"case1\"},{\"name\":\"case2\"}]'")
		os.Exit(output.ExitParamError)
		return nil
	}

	// 验证 JSON 格式
	var tcases []map[string]interface{}
	if err := json.Unmarshal([]byte(flagTCasesJSON), &tcases); err != nil {
		output.PrintError(os.Stderr, "invalid_parameter",
			fmt.Sprintf("invalid JSON for --tcases: %v", err), "")
		os.Exit(output.ExitParamError)
		return nil
	}

	// 为每个用例添加 workspace_id
	for i := range tcases {
		if _, ok := tcases[i]["workspace_id"]; !ok {
			tcases[i]["workspace_id"] = flagWorkspaceID
		}
	}

	tcasesBytes, _ := json.Marshal(tcases)
	req := &model.BatchCreateTCasesRequest{
		WorkspaceID: flagWorkspaceID,
		Data:        string(tcasesBytes),
	}

	data, err := apiClient.BatchCreateTCases(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, data, !flagPretty)
}
