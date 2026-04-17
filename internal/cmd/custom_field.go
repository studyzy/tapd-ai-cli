// Package cmd 中的 custom_field.go 实现了自定义字段和需求字段信息查询命令
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/output"
	"github.com/studyzy/tapd-sdk-go/model"
)

var flagEntityType string

// customFieldCmd 是 custom-field 命令
var customFieldCmd = &cobra.Command{
	Use:   "custom-field",
	Short: "自定义字段管理",
}

var customFieldListCmd = &cobra.Command{
	Use:   "list",
	Short: "获取自定义字段配置",
	RunE:  runCustomFieldList,
}

// storyFieldCmd 是 story-field 命令
var storyFieldCmd = &cobra.Command{
	Use:   "story-field",
	Short: "需求字段信息",
}

var storyFieldLabelCmd = &cobra.Command{
	Use:   "label",
	Short: "获取需求字段中英文名",
	RunE:  runStoryFieldLabel,
}

var storyFieldInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "获取需求字段及候选值",
	RunE:  runStoryFieldInfo,
}

// workitemTypeCmd 是 workitem-type 命令
var workitemTypeCmd = &cobra.Command{
	Use:   "workitem-type",
	Short: "需求类别管理",
}

var workitemTypeListCmd = &cobra.Command{
	Use:   "list",
	Short: "获取需求类别列表",
	RunE:  runWorkitemTypeList,
}

func init() {
	customFieldListCmd.Flags().StringVar(&flagEntityType, "entity-type", "", "类型（stories|tasks|iterations|tcases，必需）")
	customFieldCmd.AddCommand(customFieldListCmd)
	rootCmd.AddCommand(customFieldCmd)

	storyFieldCmd.AddCommand(storyFieldLabelCmd, storyFieldInfoCmd)
	rootCmd.AddCommand(storyFieldCmd)

	workitemTypeListCmd.Flags().StringVar(&flagName, "name", "", "按名称筛选")
	workitemTypeListCmd.Flags().IntVar(&flagLimit, "limit", 30, "返回数量限制")
	workitemTypeListCmd.Flags().IntVar(&flagPage, "page", 1, "页码")
	workitemTypeCmd.AddCommand(workitemTypeListCmd)
	rootCmd.AddCommand(workitemTypeCmd)
}

func runCustomFieldList(cmd *cobra.Command, args []string) error {
	if flagEntityType == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"--entity-type is required",
			"Usage: tapd custom-field list --entity-type <stories|tasks|iterations|tcases>")
		os.Exit(output.ExitParamError)
		return nil
	}

	req := &model.GetCustomFieldsRequest{
		WorkspaceID: flagWorkspaceID,
		EntityType:  flagEntityType,
	}

	data, err := apiClient.GetCustomFields(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, data, !flagPretty)
}

func runStoryFieldLabel(cmd *cobra.Command, args []string) error {
	req := &model.WorkspaceIDRequest{
		WorkspaceID: flagWorkspaceID,
	}

	data, err := apiClient.GetStoryFieldsLabel(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, data, !flagPretty)
}

func runStoryFieldInfo(cmd *cobra.Command, args []string) error {
	req := &model.WorkspaceIDRequest{
		WorkspaceID: flagWorkspaceID,
	}

	data, err := apiClient.GetStoryFieldsInfo(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, data, !flagPretty)
}

func runWorkitemTypeList(cmd *cobra.Command, args []string) error {
	req := &model.WorkspaceIDRequest{
		WorkspaceID: flagWorkspaceID,
	}

	data, err := apiClient.GetWorkitemTypes(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, data, !flagPretty)
}
