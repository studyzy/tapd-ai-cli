// Package cmd 中的 launch.go 实现发布评审管理命令
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
	flagLaunchID             string
	flagLaunchTitle          string
	flagLaunchTemplateID     string
	flagLaunchVersionType    string
	flagLaunchBaseline       string
	flagLaunchReleaseModel   string
	flagLaunchRoadmapVersion string
	flagLaunchReleaseType    string
	flagLaunchChangeType     string
	flagLaunchSignedBy       string
	flagLaunchArchivedBy     string
	flagLaunchCC             string
	flagLaunchChangeNotifier string
	flagLaunchSignerComment  string
	flagLaunchReleaseResult  string
	flagLaunchReleaseComment string
	flagLaunchRemark         string
	flagLaunchFields         string
	flagLaunchCustomFields   map[string]string
	flagLaunchAliasFields    map[string]string
)

// launchCmd 是发布评审父命令。
var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "发布评审管理",
}

var launchListCmd = &cobra.Command{
	Use:   "list",
	Short: "查询发布评审列表",
	RunE:  runLaunchList,
}

var launchCountCmd = &cobra.Command{
	Use:   "count",
	Short: "查询发布评审数量",
	RunE:  runLaunchCount,
}

var launchCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建发布评审",
	RunE:  runLaunchCreate,
}

var launchUpdateCmd = &cobra.Command{
	Use:   "update <launch_form_id>",
	Short: "更新发布评审",
	Args:  cobra.ExactArgs(1),
	RunE:  runLaunchUpdate,
}

var launchTemplatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "查询发布评审模板",
	RunE:  runLaunchTemplates,
}

var launchFieldsCmd = &cobra.Command{
	Use:   "fields",
	Short: "查询发布评审自定义字段配置",
	RunE:  runLaunchFields,
}

func init() {
	bindLaunchListFlags(launchListCmd)
	bindLaunchListFlags(launchCountCmd)
	bindLaunchCreateFlags(launchCreateCmd)
	bindLaunchUpdateFlags(launchUpdateCmd)

	launchCmd.AddCommand(launchListCmd, launchCountCmd, launchCreateCmd, launchUpdateCmd, launchTemplatesCmd, launchFieldsCmd)
	rootCmd.AddCommand(launchCmd)
}

func bindLaunchListFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&flagLaunchID, "id", "", "按发布评审 ID 筛选")
	cmd.Flags().StringVar(&flagCreator, "creator", "", "按创建人筛选")
	cmd.Flags().StringVar(&flagStatus, "status", "", "按状态筛选（initial/auditing/signing/sign_completed/finished；初始化/评审中/待签发/签发结束/发布结束）")
	cmd.Flags().StringVar(&flagLaunchTitle, "title", "", "按标题筛选")
	cmd.Flags().StringVar(&flagLaunchReleaseType, "release-type", "", "按发布类型筛选")
	cmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	cmd.Flags().IntVar(&flagPage, "page", 1, "页码")
	cmd.Flags().StringVar(&flagLaunchFields, "fields", "", "返回字段列表（逗号分隔）")
}

func bindLaunchCreateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&flagLaunchTemplateID, "template-id", "", "模板 ID（必需）")
	cmd.Flags().StringVar(&flagCreator, "creator", "", "创建人（默认当前用户）")
	cmd.Flags().StringVar(&flagLaunchTitle, "title", "", "标题")
	cmd.Flags().StringVar(&flagLaunchVersionType, "version-type", "", "版本类型")
	cmd.Flags().StringVar(&flagLaunchBaseline, "baseline", "", "基线")
	cmd.Flags().StringVar(&flagLaunchReleaseModel, "release-model", "", "发布模块")
	cmd.Flags().StringVar(&flagLaunchRoadmapVersion, "roadmap-version", "", "路标版本")
	cmd.Flags().StringVar(&flagLaunchReleaseType, "release-type", "", "发布类型")
	cmd.Flags().StringVar(&flagLaunchSignedBy, "signed-by", "", "签发人，多个用户用 ; 分隔")
	cmd.Flags().StringVar(&flagLaunchArchivedBy, "archived-by", "", "发布确认人，多个用户用 ; 分隔")
	cmd.Flags().StringVar(&flagLaunchCC, "cc", "", "抄送人，多个用户用 ; 分隔")
	cmd.Flags().StringToStringVar(&flagLaunchCustomFields, "custom-field", nil, "自定义字段，格式 custom_field_one=value，可重复或逗号分隔")
	cmd.Flags().StringToStringVar(&flagLaunchAliasFields, "alias-field", nil, "别名自定义字段，格式 字段名=value，会提交为 cus_字段名")
}

func bindLaunchUpdateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&flagLaunchTitle, "title", "", "标题")
	cmd.Flags().StringVar(&flagStatus, "status", "", "发布评审状态（initial/auditing/signing/sign_completed/finished；初始化/评审中/待签发/签发结束/发布结束）")
	cmd.Flags().StringVar(&flagLaunchVersionType, "version-type", "", "版本类型")
	cmd.Flags().StringVar(&flagLaunchBaseline, "baseline", "", "基线")
	cmd.Flags().StringVar(&flagLaunchReleaseModel, "release-model", "", "发布模块")
	cmd.Flags().StringVar(&flagLaunchRoadmapVersion, "roadmap-version", "", "路标版本")
	cmd.Flags().StringVar(&flagLaunchReleaseType, "release-type", "", "发布类型")
	cmd.Flags().StringVar(&flagLaunchChangeType, "change-type", "", "变更类型")
	cmd.Flags().StringVar(&flagLaunchSignedBy, "signed-by", "", "签发人，多个用户用 ; 分隔")
	cmd.Flags().StringVar(&flagLaunchArchivedBy, "archived-by", "", "发布确认人，多个用户用 ; 分隔")
	cmd.Flags().StringVar(&flagLaunchCC, "cc", "", "抄送人，多个用户用 ; 分隔")
	cmd.Flags().StringVar(&flagLaunchChangeNotifier, "change-notifier", "", "变更通知人，多个用户用 ; 分隔")
	cmd.Flags().StringVar(&flagLaunchSignerComment, "signer-comment", "", "签发意见")
	cmd.Flags().StringVar(&flagLaunchReleaseResult, "release-result", "", "发布结果（release_success/release_fail）")
	cmd.Flags().StringVar(&flagLaunchReleaseComment, "release-comment", "", "发布意见")
	cmd.Flags().StringVar(&flagLaunchRemark, "remark", "", "备注")
	cmd.Flags().StringToStringVar(&flagLaunchCustomFields, "custom-field", nil, "自定义字段，格式 custom_field_one=value，可重复或逗号分隔")
	cmd.Flags().StringToStringVar(&flagLaunchAliasFields, "alias-field", nil, "别名自定义字段，格式 字段名=value，会提交为 cus_字段名")
}

func newLaunchListRequest() *model.ListLaunchFormsRequest {
	return &model.ListLaunchFormsRequest{
		WorkspaceID: flagWorkspaceID,
		ID:          flagLaunchID,
		Creator:     flagCreator,
		Status:      flagStatus,
		Title:       flagLaunchTitle,
		ReleaseType: flagLaunchReleaseType,
		Limit:       flagLimit,
		Page:        flagPage,
		Fields:      flagLaunchFields,
	}
}

func runLaunchList(cmd *cobra.Command, args []string) error {
	forms, err := apiClient.ListLaunchForms(context.Background(), newLaunchListRequest())
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}

	total, _ := apiClient.CountLaunchForms(context.Background(), newLaunchListRequest())
	resp := &model.ListResponse{
		Items:   forms,
		Total:   total,
		Page:    flagPage,
		Limit:   flagLimit,
		HasMore: total > flagPage*flagLimit,
	}
	return output.PrintJSON(os.Stdout, resp, !flagPretty)
}

func runLaunchCount(cmd *cobra.Command, args []string) error {
	count, err := apiClient.CountLaunchForms(context.Background(), newLaunchListRequest())
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, &model.CountResponse{Count: count}, !flagPretty)
}

func runLaunchCreate(cmd *cobra.Command, args []string) error {
	if flagLaunchTemplateID == "" {
		output.PrintError(os.Stderr, "missing_parameter", "--template-id is required", "Usage: tapd launch create --template-id <template_id>")
		os.Exit(output.ExitParamError)
		return nil
	}

	creator := flagCreator
	if creator == "" {
		creator = ensureNick()
	}
	req := &model.CreateLaunchFormRequest{
		WorkspaceID:    flagWorkspaceID,
		Creator:        creator,
		TemplateID:     flagLaunchTemplateID,
		Title:          flagLaunchTitle,
		VersionType:    flagLaunchVersionType,
		Baseline:       flagLaunchBaseline,
		ReleaseModel:   flagLaunchReleaseModel,
		RoadmapVersion: flagLaunchRoadmapVersion,
		ReleaseType:    flagLaunchReleaseType,
		SignedBy:       flagLaunchSignedBy,
		ArchivedBy:     flagLaunchArchivedBy,
		CC:             flagLaunchCC,
		CustomFields:   flagLaunchCustomFields,
		AliasFields:    flagLaunchAliasFields,
	}
	form, err := apiClient.CreateLaunchForm(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return printSuccessResponse(form.ID, launchFormURL(form.ID), form.WorkspaceID)
}

func runLaunchUpdate(cmd *cobra.Command, args []string) error {
	req := &model.UpdateLaunchFormRequest{
		WorkspaceID:    flagWorkspaceID,
		ID:             args[0],
		Title:          flagLaunchTitle,
		Status:         flagStatus,
		VersionType:    flagLaunchVersionType,
		Baseline:       flagLaunchBaseline,
		ReleaseModel:   flagLaunchReleaseModel,
		RoadmapVersion: flagLaunchRoadmapVersion,
		ReleaseType:    flagLaunchReleaseType,
		ChangeType:     flagLaunchChangeType,
		SignedBy:       flagLaunchSignedBy,
		ArchivedBy:     flagLaunchArchivedBy,
		CC:             flagLaunchCC,
		ChangeNotifier: flagLaunchChangeNotifier,
		SignerComment:  flagLaunchSignerComment,
		ReleaseResult:  flagLaunchReleaseResult,
		ReleaseComment: flagLaunchReleaseComment,
		Remark:         flagLaunchRemark,
		CustomFields:   flagLaunchCustomFields,
		AliasFields:    flagLaunchAliasFields,
	}
	form, err := apiClient.UpdateLaunchForm(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return printSuccessResponse(form.ID, launchFormURL(form.ID), form.WorkspaceID)
}

func runLaunchTemplates(cmd *cobra.Command, args []string) error {
	data, err := apiClient.ListLaunchFormTemplates(context.Background(), &model.WorkspaceIDRequest{WorkspaceID: flagWorkspaceID})
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, data, !flagPretty)
}

func runLaunchFields(cmd *cobra.Command, args []string) error {
	data, err := apiClient.GetLaunchFormCustomFields(context.Background(), &model.WorkspaceIDRequest{WorkspaceID: flagWorkspaceID})
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, data, !flagPretty)
}

func launchFormURL(id string) string {
	if id == "" || flagWorkspaceID == "" || apiClient == nil {
		return ""
	}
	return fmt.Sprintf("%s/%s/launch_forms/view/%s", apiClient.WebURL(), flagWorkspaceID, id)
}
