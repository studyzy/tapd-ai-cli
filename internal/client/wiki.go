// Package client 中的 wiki.go 实现了 Wiki 文档相关的 API 调用
package client

import (
	"encoding/json"
	"fmt"
)

// ListWikis 查询 Wiki 文档列表
func (c *Client) ListWikis(params map[string]string) ([]map[string]interface{}, error) {
	data, err := c.doGet("/tapd_wikis", params)
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse wiki list: %w", err)
	}

	var results []map[string]interface{}
	for _, item := range rawList {
		if raw, ok := item["Wiki"]; ok {
			var obj map[string]interface{}
			if err := json.Unmarshal(raw, &obj); err == nil {
				results = append(results, obj)
			}
		}
	}
	return results, nil
}

// GetWiki 获取单个 Wiki 文档详情，返回包含 markdown_description 字段的完整内容
func (c *Client) GetWiki(workspaceID, id string) (map[string]interface{}, error) {
	params := map[string]string{
		"workspace_id": workspaceID,
		"id":           id,
	}

	data, err := c.doGet("/tapd_wikis", params)
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(rawList) == 0 {
		return nil, &TAPDError{ExitCode: 2, Message: fmt.Sprintf("wiki %s not found", id)}
	}

	raw, ok := rawList[0]["Wiki"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("failed to parse wiki: %w", err)
	}

	return result, nil
}
