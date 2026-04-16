// Package cmd 中的 auth.go 实现了认证相关命令
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/client"
	"github.com/studyzy/tapd-ai-cli/internal/config"
	"github.com/studyzy/tapd-ai-cli/internal/model"
	"github.com/studyzy/tapd-ai-cli/internal/output"
)

var flagLocal bool

// authCmd 是 auth 父命令
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "认证管理",
}

// loginCmd 是 auth login 子命令
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "登录并持久化凭据到配置文件",
	Long:  "支持 --access-token 或 --api-user/--api-password 两种认证方式。默认写入 ~/.tapd.json，使用 --local 写入当前目录。",
	RunE:  runLogin,
}

func init() {
	loginCmd.Flags().BoolVar(&flagLocal, "local", false, "将凭据写入当前目录的 .tapd.json")
	authCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(authCmd)
}

// runLogin 执行登录逻辑
func runLogin(cmd *cobra.Command, args []string) error {
	accessToken := flagAccessToken
	apiUser := flagAPIUser
	apiPassword := flagAPIPassword

	// 验证参数
	if accessToken == "" && (apiUser == "" || apiPassword == "") {
		output.PrintError(os.Stderr, "missing_parameter",
			"Provide --access-token or --api-user with --api-password",
			"Usage: tapd auth login --access-token <token> OR --api-user <user> --api-password <password>")
		os.Exit(output.ExitParamError)
		return nil
	}

	// 验证凭据有效性
	c := client.NewClient(accessToken, apiUser, apiPassword)
	if err := c.TestAuth(); err != nil {
		output.PrintError(os.Stderr, "authentication_failed",
			"Invalid credentials: "+err.Error(),
			"Please check your access token or API user/password and try again.")
		os.Exit(output.ExitAuthError)
		return nil
	}

	// 构建配置
	cfg := &model.Config{}
	if accessToken != "" {
		cfg.AccessToken = accessToken
	} else {
		cfg.APIUser = apiUser
		cfg.APIPassword = apiPassword
	}

	// 确定写入路径
	configPath, err := config.GetConfigPath(flagLocal)
	if err != nil {
		output.PrintError(os.Stderr, "config_error",
			"Failed to determine config path: "+err.Error(),
			"Check file system permissions.")
		os.Exit(output.ExitAPIError)
		return nil
	}

	// 保存配置
	if err := config.SaveConfig(cfg, configPath); err != nil {
		output.PrintError(os.Stderr, "config_error",
			"Failed to save config: "+err.Error(),
			"Check file system permissions for "+configPath)
		os.Exit(output.ExitAPIError)
		return nil
	}

	return output.PrintJSON(os.Stdout, &model.SuccessResponse{Success: true}, !flagPretty)
}
