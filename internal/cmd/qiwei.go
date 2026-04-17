// Package cmd 中的 qiwei.go 实现了企业微信消息发送命令
package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/output"
)

var (
	flagQiweiWebhook string
	flagQiweiMsg     string
)

// qiweiCmd 是 qiwei 父命令
var qiweiCmd = &cobra.Command{
	Use:   "qiwei",
	Short: "企业微信消息管理",
}

var qiweiSendCmd = &cobra.Command{
	Use:   "send",
	Short: "发送消息到企业微信群（通过机器人 Webhook）",
	RunE:  runQiweiSend,
}

func init() {
	qiweiSendCmd.Flags().StringVar(&flagQiweiWebhook, "webhook", "", "企业微信机器人 Webhook URL（必需，也可通过 BOT_URL 环境变量设置）")
	qiweiSendCmd.Flags().StringVar(&flagQiweiMsg, "msg", "", "消息内容（Markdown 格式，必需）")

	qiweiCmd.AddCommand(qiweiSendCmd)
	rootCmd.AddCommand(qiweiCmd)
}

func runQiweiSend(cmd *cobra.Command, args []string) error {
	// webhook 优先使用命令行参数，否则使用环境变量
	webhook := flagQiweiWebhook
	if webhook == "" {
		webhook = os.Getenv("BOT_URL")
	}
	if webhook == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"Webhook URL is required",
			"Usage: tapd qiwei send --webhook <url> --msg <content>, or set BOT_URL environment variable")
		os.Exit(output.ExitParamError)
		return nil
	}

	if flagQiweiMsg == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"--msg is required",
			"Usage: tapd qiwei send --webhook <url> --msg <content>")
		os.Exit(output.ExitParamError)
		return nil
	}

	err := apiClient.SendQiweiMessage(context.Background(), webhook, flagQiweiMsg)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}

	return output.PrintJSON(os.Stdout, map[string]interface{}{
		"success": true,
		"message": "Message sent successfully",
	}, !flagPretty)
}
