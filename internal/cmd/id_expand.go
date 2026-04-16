// Package cmd 中的 id_expand.go 提供短 ID 自动展开为 TAPD 长 ID 的辅助函数
package cmd

import (
	"fmt"
	"unicode"
)

// expandShortID 将短 ID 自动展开为 TAPD 长 ID。
// 规则：如果 id 全部是数字且长度 <= 9，则左补零到 9 位，前缀 "11" + workspaceID。
// 如果 id 已经超过 9 位或包含非数字字符，则原样返回。
func expandShortID(id, workspaceID string) string {
	if id == "" || workspaceID == "" {
		return id
	}
	// 检查是否全部为数字
	for _, c := range id {
		if !unicode.IsDigit(c) {
			return id
		}
	}
	// 长度超过 9 位视为已经是完整 ID
	if len(id) > 9 {
		return id
	}
	// 左补零到 9 位，前缀 "11" + workspaceID
	padded := fmt.Sprintf("%09s", id)
	return "11" + workspaceID + padded
}
