// Package cmd 中的 attachment.go 实现了附件和图片管理命令
package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/output"
	"github.com/studyzy/tapd-sdk-go/model"
)

var (
	flagAttachmentEntryID string
	flagAttachmentType    string
	flagImagePath         string
)

// attachmentCmd 是 attachment 父命令
var attachmentCmd = &cobra.Command{
	Use:   "attachment",
	Short: "附件管理",
}

var attachmentListCmd = &cobra.Command{
	Use:   "list",
	Short: "查询条目的附件列表",
	RunE:  runAttachmentList,
}

// imageCmd 是 image 父命令
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "图片管理",
}

var imageGetCmd = &cobra.Command{
	Use:   "get",
	Short: "获取图片下载链接",
	RunE:  runImageGet,
}

func init() {
	// attachment list 子命令
	attachmentListCmd.Flags().StringVar(&flagAttachmentEntryID, "entry-id", "", "条目 ID（必需）")
	attachmentListCmd.Flags().StringVar(&flagAttachmentType, "type", "", "条目类型（story|bug|task）")
	attachmentListCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	attachmentListCmd.Flags().IntVar(&flagPage, "page", 1, "页码")

	// image get 子命令
	imageGetCmd.Flags().StringVar(&flagImagePath, "image-path", "", "图片路径（必需，从条目描述中获取）")

	attachmentCmd.AddCommand(attachmentListCmd)
	rootCmd.AddCommand(attachmentCmd)

	imageCmd.AddCommand(imageGetCmd)
	rootCmd.AddCommand(imageCmd)
}

func runAttachmentList(cmd *cobra.Command, args []string) error {
	if flagAttachmentEntryID == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"--entry-id is required",
			"Usage: tapd attachment list --entry-id <id>")
		os.Exit(output.ExitParamError)
		return nil
	}

	req := &model.GetAttachmentsRequest{
		WorkspaceID: flagWorkspaceID,
		EntryID:     flagAttachmentEntryID,
		Type:        flagAttachmentType,
		Limit:       flagLimit,
		Page:        flagPage,
	}

	attachments, err := apiClient.GetAttachments(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, attachments, !flagPretty)
}

func runImageGet(cmd *cobra.Command, args []string) error {
	if flagImagePath == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"--image-path is required",
			"Usage: tapd image get --image-path <path>")
		os.Exit(output.ExitParamError)
		return nil
	}

	req := &model.GetImageRequest{
		WorkspaceID: flagWorkspaceID,
		ImagePath:   flagImagePath,
	}

	img, err := apiClient.GetImage(context.Background(), req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, img, !flagPretty)
}
