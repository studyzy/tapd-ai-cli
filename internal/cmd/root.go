// Package cmd 定义了 tapd-ai-cli 的所有 Cobra 命令
package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/config"
	"github.com/studyzy/tapd-ai-cli/internal/output"
	tapd "github.com/studyzy/tapd-sdk-go"
	"github.com/studyzy/tapd-sdk-go/model"
)

var (
	// 全局标志
	flagWorkspaceID string
	flagPretty      bool
	flagJSON        bool
	flagNoComments  bool
	flagAccessToken string
	flagAPIUser     string
	flagAPIPassword string

	// 全局共享的客户端和配置
	apiClient *tapd.Client
	appConfig *config.Config

	// list 命令共享的 filter 标志
	flagFilter []string

	// filterFlagDesc 是 --filter 标志的描述文本，统一常量避免 13 处重复
	filterFlagDesc = "高级过滤条件（可重复，格式：field=OP<value>，支持 LIKE/EQ/CONTAINS 等 OpenAPI 特殊查询语法）"
)

// rootCmd 是 CLI 的根命令
var rootCmd = &cobra.Command{
	Use:   "tapd",
	Short: "面向 AI Agent 的 TAPD 命令行工具",
	Long:  "tapd-ai-cli 是一个面向 AI Agent 的 TAPD 命令行工具，通过 TAPD Open API 实现项目管理核心操作。",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// auth login 命令不需要预加载配置和客户端
		if cmd.Name() == "login" || cmd.Name() == "init" {
			return nil
		}
		// --version 不需要认证
		if v, _ := cmd.Flags().GetBool("version"); v {
			return nil
		}
		// 兜底的父命令（如 tapd bug）只输出 help 或报未知子命令，不需要认证
		if _, ok := cmd.Annotations[annotationGuardedParent]; ok {
			return nil
		}
		return initClientAndConfig(cmd)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute 执行根命令
func Execute() {
	guardUnknownSubcommands(rootCmd)
	if cmd, err := rootCmd.ExecuteC(); err != nil {
		printCommandError(os.Stderr, cmd, err)
		os.Exit(1)
	}
}

// printCommandError 输出未被子命令自行处理的错误。
// rootCmd 开启了 SilenceErrors，cobra 层错误（未知参数/子命令、位置参数个数不对等）
// 和 RunE 直接返回的错误都不会被打印；不在这里输出的话，调用方只能看到退出码 1、零输出。
func printCommandError(w io.Writer, cmd *cobra.Command, err error) {
	path := "tapd <command>"
	if cmd != nil {
		path = cmd.CommandPath()
	}
	output.PrintError(w, "command_error", err.Error(), fmt.Sprintf("Run '%s --help' for usage", path))
}

// annotationGuardedParent 标记由 guardUnknownSubcommands 补了兜底 RunE 的父命令
const annotationGuardedParent = "tapd/guarded-parent"

// guardUnknownSubcommands 给没有 Run 的父命令补兜底 RunE。
// cobra 只在根命令上校验子命令名：`tapd bug comment add` 这类写错的子命令
// 会打印 bug 的 help 后以退出码 0 返回，调用方误以为执行成功。
func guardUnknownSubcommands(parent *cobra.Command) {
	for _, c := range parent.Commands() {
		if !c.HasSubCommands() {
			continue
		}
		if !c.Runnable() {
			if c.Annotations == nil {
				c.Annotations = map[string]string{}
			}
			c.Annotations[annotationGuardedParent] = ""
			c.RunE = func(cmd *cobra.Command, args []string) error {
				if len(args) > 0 {
					return fmt.Errorf("unknown command %q for %q", args[0], cmd.CommandPath())
				}
				return cmd.Help()
			}
		}
		guardUnknownSubcommands(c)
	}
}

func init() {
	// 根命令自定义 help：输出紧凑参考卡（原 spec 子命令的功能）
	defaultHelp := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd != rootCmd {
			defaultHelp(cmd, args)
			return
		}
		coreLines, advancedLines := buildSpecLines(rootCmd)
		printSpecOutput(os.Stdout, rootCmd, coreLines, advancedLines)
	})

	rootCmd.PersistentFlags().StringVar(&flagWorkspaceID, "workspace-id", "", "指定工作区 ID（覆盖本地配置）")
	rootCmd.PersistentFlags().BoolVar(&flagPretty, "pretty", false, "输出带缩进的 JSON，仅供人类阅读，AI Agent 不应使用（浪费 token）")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "强制 JSON 输出（列表默认已是 JSON，详情默认 Markdown 更省 token；仅在需要从详情提取字段时使用）")
	rootCmd.PersistentFlags().StringVar(&flagAccessToken, "access-token", "", "TAPD Access Token")
	rootCmd.PersistentFlags().StringVar(&flagAPIUser, "api-user", "", "TAPD API 用户名")
	rootCmd.PersistentFlags().StringVar(&flagAPIPassword, "api-password", "", "TAPD API 密码")
	rootCmd.PersistentFlags().BoolVar(&flagNoComments, "no-comments", false, "不展示评论")
}

// initClientAndConfig 初始化配置和 API 客户端
func initClientAndConfig(cmd *cobra.Command) error {
	// 加载配置文件
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	appConfig = cfg

	// 命令行标志覆盖配置
	accessToken := flagAccessToken
	apiUser := flagAPIUser
	apiPassword := flagAPIPassword

	if accessToken == "" {
		accessToken = cfg.AccessToken
	}
	if apiUser == "" {
		apiUser = cfg.APIUser
	}
	if apiPassword == "" {
		apiPassword = cfg.APIPassword
	}

	// 检查是否有有效凭据
	if accessToken == "" && (apiUser == "" || apiPassword == "") {
		output.PrintError(os.Stderr, "authentication_required",
			"No valid credentials found",
			"Run 'tapd auth login --access-token <token>' or 'tapd auth login --api-user <user> --api-password <password>'. "+
				"You can also set TAPD_ACCESS_TOKEN or TAPD_API_USER/TAPD_API_PASSWORD environment variables.")
		os.Exit(output.ExitAuthError)
	}

	apiClient = tapd.NewClientWithBaseURL(cfg.APIBaseURL, cfg.BaseURL, accessToken, apiUser, apiPassword)

	// workspace-id 标志覆盖配置
	if flagWorkspaceID == "" {
		flagWorkspaceID = cfg.WorkspaceID
	}

	// 需要 workspace_id 的命令检查
	// 以下命令不需要 workspace_id：
	// - workspace list：列出所有工作区
	// - auth 子命令：认证操作
	// - url 命令：从 URL 中提取 workspace ID
	skipWorkspace := map[string]bool{"auth": true, "workspace": true}
	parentName := ""
	if cmd.Parent() != nil {
		parentName = cmd.Parent().Name()
	}
	needsWorkspace := !skipWorkspace[parentName] &&
		cmd.Name() != "url" &&
		!(cmd.Name() == "list" && parentName == "workspace")
	if needsWorkspace && flagWorkspaceID == "" {
		output.PrintError(os.Stderr, "workspace_required",
			"No workspace ID configured",
			"Run 'tapd workspace switch <id>' or use --workspace-id flag.")
		os.Exit(output.ExitParamError)
	}

	return nil
}

// cmdContext 安全获取命令的 context，当 cmd 为 nil 或 context 未设置时回退到 context.Background()
func cmdContext(cmd *cobra.Command) context.Context {
	if cmd != nil {
		defer func() { recover() }()
		if ctx := cmd.Context(); ctx != nil {
			return ctx
		}
	}
	return context.Background()
}

// ensureNick 按需获取当前用户昵称，仅在首次调用时发起 HTTP 请求
func ensureNick() string {
	if apiClient.GetNick() == "" {
		apiClient.FetchNick(context.Background())
	}
	return apiClient.GetNick()
}

// useJSONOutput 判断是否应使用 JSON 格式输出，--pretty 隐含 --json
func useJSONOutput() bool {
	return flagJSON || flagPretty
}

// printDetail 输出单条详情，默认 Markdown 格式，--json/--pretty 时输出 JSON
// bodyField 指定作为 Markdown body 的字段名（JSON tag 名称）
func printDetail(data interface{}, bodyField string) error {
	if useJSONOutput() {
		return output.PrintJSON(os.Stdout, data, !flagPretty)
	}
	return output.PrintMarkdown(os.Stdout, data, bodyField)
}

// printSuccessResponse 输出创建/更新操作的精简成功响应，节省 AI Agent token 消耗
func printSuccessResponse(id, url, workspaceID string) error {
	resp := &model.SuccessResponse{
		Success:     true,
		ID:          id,
		URL:         url,
		WorkspaceID: workspaceID,
	}
	return output.PrintJSON(os.Stdout, resp, !flagPretty)
}

// printComments 获取并输出条目的评论列表
// entryType 取值：stories|bug|tasks，entryID 为条目 ID
// 当 --no-comments 标志启用或获取失败时静默跳过
func printComments(workspaceID, entryType, entryID string) {
	if flagNoComments {
		return
	}
	comments, err := apiClient.ListComments(context.Background(), &model.ListCommentsRequest{
		WorkspaceID: workspaceID,
		EntryType:   entryType,
		EntryID:     entryID,
	})
	if err != nil || len(comments) == 0 {
		return
	}
	if useJSONOutput() {
		fmt.Fprintln(os.Stdout)
		output.PrintJSON(os.Stdout, map[string]interface{}{
			"comments": comments,
			"count":    len(comments),
		}, !flagPretty)
		return
	}
	fmt.Fprintf(os.Stdout, "\n## 评论 (%d)\n\n", len(comments))
	for i := range comments {
		comments[i].Description = htmlToMarkdown(comments[i].Description)
	}
	for _, c := range comments {
		fmt.Fprintf(os.Stdout, "**%s** (%s):\n%s\n\n", c.Author, c.Created, c.Description)
	}
}

// parseCustomFields 将 ["key1=value1","key2=value2"] 形式的切片解析为 map[string]string
// 用于支持 --custom-field key=value 可重复 flag
func parseCustomFields(fields []string) map[string]string {
	if len(fields) == 0 {
		return nil
	}
	m := make(map[string]string, len(fields))
	for _, f := range fields {
		k, v, ok := strings.Cut(f, "=")
		if !ok || k == "" {
			continue
		}
		m[k] = v
	}
	return m
}

// parseFilters 解析 --filter 参数为 map，格式无效时打印警告到 stderr
func parseFilters(fields []string) map[string]string {
	if len(fields) == 0 {
		return nil
	}
	m := make(map[string]string, len(fields))
	for _, f := range fields {
		k, v, ok := strings.Cut(f, "=")
		if !ok || k == "" {
			fmt.Fprintf(os.Stderr, "warning: skipping invalid filter %q (expected format: key=value)\n", f)
			continue
		}
		m[k] = v
	}
	return m
}

// listWithFilters 是一个泛型辅助函数，用于 list 命令支持 --filter 参数。
// 它将 request struct 的 ToParams() 结果与 --filter 解析出的额外参数合并到新 map，
// 然后通过 SDK 的 DoGet 方法直接发送请求，再用 ParseList 解析 TAPD 包装响应。
// 注意：不会修改传入的 params map；filter 值会覆盖同名 params 键。
func listWithFilters[T any](ctx context.Context, client *tapd.Client, endpoint string, params map[string]string, filters []string, wrapperKey string) ([]T, error) {
	merged := make(map[string]string, len(params)+len(filters))
	for k, v := range params {
		merged[k] = v
	}
	for k, v := range parseFilters(filters) {
		merged[k] = v
	}
	data, err := client.DoGet(ctx, endpoint, merged)
	if err != nil {
		return nil, err
	}
	return tapd.ParseList[T](data, wrapperKey)
}
