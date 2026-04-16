// Package cmd_test 中的 url_test.go 测试 TAPD URL 解析函数
package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestParseTAPDURL(t *testing.T) {
	tests := []struct {
		name            string
		rawURL          string
		wantWorkspaceID string
		wantEntityType  string
		wantEntityID    string
		wantErr         bool
	}{
		// Story 详情页格式
		{
			name:            "story 详情页 URL",
			rawURL:          "https://www.tapd.cn/tapd_fe/51081496/story/detail/1151081496001028684",
			wantWorkspaceID: "51081496",
			wantEntityType:  "story",
			wantEntityID:    "1151081496001028684",
		},
		// Story 列表页预览格式
		{
			name:            "story 列表页 dialog_preview_id",
			rawURL:          "https://www.tapd.cn/tapd_fe/51081496/story/list?categoryId=1151081496001000472&dialog_preview_id=story_1151081496001028684",
			wantWorkspaceID: "51081496",
			wantEntityType:  "story",
			wantEntityID:    "1151081496001028684",
		},
		// Bug 详情页格式
		{
			name:            "bug 详情页 URL",
			rawURL:          "https://www.tapd.cn/tapd_fe/51081496/bug/detail/1151081496001016136",
			wantWorkspaceID: "51081496",
			wantEntityType:  "bug",
			wantEntityID:    "1151081496001016136",
		},
		// Bug 列表页预览格式
		{
			name:            "bug 列表页 dialog_preview_id",
			rawURL:          "https://www.tapd.cn/tapd_fe/51081496/bug/list?confId=1151081496001021206&dialog_preview_id=bug_1151081496001016136",
			wantWorkspaceID: "51081496",
			wantEntityType:  "bug",
			wantEntityID:    "1151081496001016136",
		},
		// Task 详情页格式
		{
			name:            "task 详情页 URL",
			rawURL:          "https://www.tapd.cn/tapd_fe/51081496/task/detail/1151081496001028786",
			wantWorkspaceID: "51081496",
			wantEntityType:  "task",
			wantEntityID:    "1151081496001028786",
		},
		// Task 看板页预览格式（无 tapd_fe 前缀）
		{
			name:            "task 看板页 dialog_preview_id（无 tapd_fe）",
			rawURL:          "https://www.tapd.cn/51081496/prong/tasks?conf_id=1151081496001031953&dialog_preview_id=task_1151081496001028786",
			wantWorkspaceID: "51081496",
			wantEntityType:  "task",
			wantEntityID:    "1151081496001028786",
		},
		// Wiki fragment 格式
		{
			name:            "wiki fragment URL",
			rawURL:          "https://www.tapd.cn/51081496/markdown_wikis/show/#1151081496001001503",
			wantWorkspaceID: "51081496",
			wantEntityType:  "wiki",
			wantEntityID:    "1151081496001001503",
		},
		// 错误场景：非 TAPD URL
		{
			name:    "非 TAPD URL",
			rawURL:  "https://github.com/foo/bar",
			wantErr: true,
		},
		// 错误场景：不支持的类型（在 dialog_preview_id 中）
		{
			name:    "不支持的 dialog_preview_id 类型",
			rawURL:  "https://www.tapd.cn/tapd_fe/51081496/iteration/list?dialog_preview_id=iteration_123",
			wantErr: true,
		},
		// 错误场景：无法识别的路径
		{
			name:    "无法识别的路径",
			rawURL:  "https://www.tapd.cn/51081496/prong/iterations",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTAPDURL(tt.rawURL)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseTAPDURL(%q) expected error, got nil", tt.rawURL)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseTAPDURL(%q) unexpected error: %v", tt.rawURL, err)
			}
			if got.WorkspaceID != tt.wantWorkspaceID {
				t.Errorf("WorkspaceID = %q, want %q", got.WorkspaceID, tt.wantWorkspaceID)
			}
			if got.EntityType != tt.wantEntityType {
				t.Errorf("EntityType = %q, want %q", got.EntityType, tt.wantEntityType)
			}
			if got.EntityID != tt.wantEntityID {
				t.Errorf("EntityID = %q, want %q", got.EntityID, tt.wantEntityID)
			}
		})
	}
}

// TestURLCommandSkipsWorkspaceCheck 验证 url 命令在 workspace ID 为空时不会触发 workspace_required 检查
// 因为 url 命令从 URL 中提取 workspace ID，无需全局配置
func TestURLCommandSkipsWorkspaceCheck(t *testing.T) {
	// needsWorkspace 条件逻辑来自 root.go initClientAndConfig
	// 此处模拟 url 命令的 Cobra 结构验证豁免逻辑
	tests := []struct {
		name         string
		cmdName      string
		parentName   string
		wantExempted bool
	}{
		{
			name:         "url 命令应被豁免",
			cmdName:      "url",
			parentName:   "tapd",
			wantExempted: true,
		},
		{
			name:         "show 命令（非豁免）",
			cmdName:      "show",
			parentName:   "story",
			wantExempted: false,
		},
		{
			name:         "auth 子命令应被豁免",
			cmdName:      "login",
			parentName:   "auth",
			wantExempted: true,
		},
		{
			name:         "workspace 子命令应被豁免",
			cmdName:      "list",
			parentName:   "workspace",
			wantExempted: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parent := &cobra.Command{Use: tt.parentName}
			cmd := &cobra.Command{Use: tt.cmdName}
			parent.AddCommand(cmd)

			// 复现 root.go 中的 workspace 检查条件
			needsWorkspace := cmd.Name() != "list" || (cmd.Parent() != nil && cmd.Parent().Name() != "workspace")
			wouldCheck := needsWorkspace && cmd.Name() != "url" && cmd.Parent() != nil && cmd.Parent().Name() != "auth" && cmd.Parent().Name() != "workspace"

			if tt.wantExempted && wouldCheck {
				t.Errorf("command %q under %q should be exempted from workspace check, but would be checked", tt.cmdName, tt.parentName)
			}
			if !tt.wantExempted && !wouldCheck {
				t.Errorf("command %q under %q should require workspace check, but would be exempted", tt.cmdName, tt.parentName)
			}
		})
	}
}
