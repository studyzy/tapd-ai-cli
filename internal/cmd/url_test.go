// Package cmd_test 中的 url_test.go 测试 TAPD URL 解析函数
package cmd

import (
	"testing"
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
