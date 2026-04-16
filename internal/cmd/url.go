// Package cmd 中的 url.go 实现了根据 TAPD URL 查询对应条目详情的通用命令
package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/output"
)

// urlCmd 是根据 TAPD URL 查询条目详情的命令
var urlCmd = &cobra.Command{
	Use:   "url <url>",
	Short: "根据 TAPD URL 查询对应条目详情（支持需求、缺陷、任务、Wiki）",
	Long: `根据 TAPD URL 自动识别条目类型并查询详情。

支持以下 URL 格式：
  需求详情页:  https://www.tapd.cn/tapd_fe/{workspace_id}/story/detail/{id}
  需求列表页:  https://www.tapd.cn/tapd_fe/{workspace_id}/story/list?...&dialog_preview_id=story_{id}
  缺陷详情页:  https://www.tapd.cn/tapd_fe/{workspace_id}/bug/detail/{id}
  缺陷列表页:  https://www.tapd.cn/tapd_fe/{workspace_id}/bug/list?...&dialog_preview_id=bug_{id}
  任务详情页:  https://www.tapd.cn/tapd_fe/{workspace_id}/task/detail/{id}
  任务看板页:  https://www.tapd.cn/{workspace_id}/prong/tasks?...&dialog_preview_id=task_{id}
  Wiki 文档:   https://www.tapd.cn/{workspace_id}/markdown_wikis/show/#{id}`,
	Args: cobra.ExactArgs(1),
	RunE: runURLQuery,
}

func init() {
	rootCmd.AddCommand(urlCmd)
}

// parsedTAPDURL 保存从 TAPD URL 中解析出的关键信息
type parsedTAPDURL struct {
	WorkspaceID string
	EntityType  string
	EntityID    string
}

// parseTAPDURL 解析 TAPD URL，提取工作区 ID、条目类型和条目 ID
// 支持三种格式：
//  1. dialog_preview_id 查询参数格式（story/bug/task）
//  2. /tapd_fe/{ws}/{type}/detail/{id} 详情页路径格式
//  3. /{ws}/markdown_wikis/show/#{id} Wiki fragment 格式
func parseTAPDURL(rawURL string) (*parsedTAPDURL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	host := u.Hostname()
	if host != "www.tapd.cn" && host != "tapd.cn" {
		return nil, fmt.Errorf("not a valid TAPD URL: host must be www.tapd.cn or tapd.cn, got %q", host)
	}

	segments := splitPath(u.Path)

	// 格式一：检查 dialog_preview_id 查询参数（优先级最高）
	if previewID := u.Query().Get("dialog_preview_id"); previewID != "" {
		entityType, entityID, err := parsePreviewID(previewID)
		if err != nil {
			return nil, err
		}
		workspaceID, err := extractWorkspaceID(segments)
		if err != nil {
			return nil, err
		}
		return &parsedTAPDURL{WorkspaceID: workspaceID, EntityType: entityType, EntityID: entityID}, nil
	}

	// 格式三：Wiki fragment 格式（/{ws}/markdown_wikis/show/#{id}）
	if containsSegment(segments, "markdown_wikis") {
		workspaceID, err := extractWorkspaceID(segments)
		if err != nil {
			return nil, err
		}
		wikiID := strings.TrimSpace(u.Fragment)
		if wikiID == "" {
			return nil, fmt.Errorf("wiki URL missing ID in fragment (#)")
		}
		return &parsedTAPDURL{WorkspaceID: workspaceID, EntityType: "wiki", EntityID: wikiID}, nil
	}

	// 格式二：/tapd_fe/{ws}/{type}/detail/{id} 详情页路径
	entityType, entityID, err := parseDetailPath(segments)
	if err != nil {
		return nil, err
	}
	workspaceID, err := extractWorkspaceID(segments)
	if err != nil {
		return nil, err
	}
	return &parsedTAPDURL{WorkspaceID: workspaceID, EntityType: entityType, EntityID: entityID}, nil
}

// splitPath 将 URL 路径按 "/" 分割并过滤空字符串
func splitPath(path string) []string {
	var segments []string
	for _, s := range strings.Split(path, "/") {
		if s != "" {
			segments = append(segments, s)
		}
	}
	return segments
}

// containsSegment 检查路径段中是否包含指定值
func containsSegment(segments []string, target string) bool {
	for _, s := range segments {
		if s == target {
			return true
		}
	}
	return false
}

// extractWorkspaceID 从路径段中提取工作区 ID
// tapd_fe 路径：segments[0]="tapd_fe", segments[1]=workspaceID
// 直接路径：segments[0]=workspaceID
func extractWorkspaceID(segments []string) (string, error) {
	if len(segments) == 0 {
		return "", fmt.Errorf("URL path is empty, cannot extract workspace ID")
	}
	if segments[0] == "tapd_fe" {
		if len(segments) < 2 {
			return "", fmt.Errorf("URL path too short to extract workspace ID")
		}
		return segments[1], nil
	}
	return segments[0], nil
}

// parsePreviewID 解析 dialog_preview_id 参数值（格式："{type}_{id}"）
func parsePreviewID(previewID string) (entityType, entityID string, err error) {
	knownTypes := []string{"story", "bug", "task"}
	for _, t := range knownTypes {
		prefix := t + "_"
		if strings.HasPrefix(previewID, prefix) {
			return t, strings.TrimPrefix(previewID, prefix), nil
		}
	}
	return "", "", fmt.Errorf("unsupported entity type in dialog_preview_id: %q, supported types: story, bug, task, wiki", previewID)
}

// parseDetailPath 从详情页路径中解析条目类型和 ID
// 路径格式：/tapd_fe/{ws}/{type}/detail/{id} 或 /{ws}/{type}/detail/{id}
func parseDetailPath(segments []string) (entityType, entityID string, err error) {
	// 查找 "detail" 段的位置
	for i, s := range segments {
		if s == "detail" && i > 0 && i+1 < len(segments) {
			typeSegment := segments[i-1]
			id := segments[i+1]
			switch typeSegment {
			case "story":
				return "story", id, nil
			case "bug":
				return "bug", id, nil
			case "task":
				return "task", id, nil
			default:
				return "", "", fmt.Errorf("unsupported entity type: %q, supported types: story, bug, task, wiki", typeSegment)
			}
		}
	}
	return "", "", fmt.Errorf("cannot identify TAPD entity type from URL path")
}

// runURLQuery 是 url 命令的执行函数
func runURLQuery(cmd *cobra.Command, args []string) error {
	parsed, err := parseTAPDURL(args[0])
	if err != nil {
		output.PrintError(os.Stderr, "invalid_tapd_url", err.Error(),
			"provide a TAPD URL like https://www.tapd.cn/tapd_fe/{workspace_id}/story/detail/{id}")
		os.Exit(output.ExitParamError)
		return nil
	}

	// URL 中的 workspaceID 覆盖全局配置
	workspaceID := parsed.WorkspaceID

	switch parsed.EntityType {
	case "story":
		result, err := apiClient.GetStory(workspaceID, parsed.EntityID)
		if err != nil {
			handleAPIError(err)
			return nil
		}
		if err := printDetail(result, "description"); err != nil {
			return err
		}
		printComments(workspaceID, "stories", parsed.EntityID)
		return nil

	case "bug":
		result, err := apiClient.GetBug(workspaceID, parsed.EntityID)
		if err != nil {
			handleAPIError(err)
			return nil
		}
		if err := printDetail(result, "description"); err != nil {
			return err
		}
		printComments(workspaceID, "bug", parsed.EntityID)
		return nil

	case "task":
		result, err := apiClient.GetTask(workspaceID, parsed.EntityID)
		if err != nil {
			handleAPIError(err)
			return nil
		}
		if err := printDetail(result, "description"); err != nil {
			return err
		}
		printComments(workspaceID, "tasks", parsed.EntityID)
		return nil

	case "wiki":
		result, err := apiClient.GetWiki(workspaceID, parsed.EntityID)
		if err != nil {
			handleAPIError(err)
			return nil
		}
		if err := printDetail(result, "markdown_description"); err != nil {
			return err
		}
		printComments(workspaceID, "wiki", parsed.EntityID)
		return nil

	default:
		output.PrintError(os.Stderr, "unsupported_entity_type",
			fmt.Sprintf("unsupported TAPD entity type: %q", parsed.EntityType),
			"supported types: story, bug, task, wiki")
		os.Exit(output.ExitParamError)
		return nil
	}
}

// handleAPIError 统一处理 API 调用错误
func handleAPIError(err error) {
	output.PrintError(os.Stderr, "api_error", err.Error(), "")
	os.Exit(output.ExitAPIError)
}
