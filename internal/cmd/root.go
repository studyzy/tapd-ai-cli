// Package cmd 定义了 tapd-ai-cli 的所有 Cobra 命令
package cmd

import (
	"context"
	"fmt"
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
		return initClientAndConfig(cmd)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
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
		lines := buildSpecLines(rootCmd)
		printSpecOutput(os.Stdout, rootCmd, lines)
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
