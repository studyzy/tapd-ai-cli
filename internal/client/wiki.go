// Package client 中的 wiki.go 实现了 Wiki 文档相关的 API 调用
package client

import (
	"encoding/json"
	"fmt"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/studyzy/tapd-ai-cli/internal/model"
)

// ListWikis 查询 Wiki 文档列表，返回强类型 []model.Wiki
func (c *Client) ListWikis(params map[string]string) ([]model.Wiki, error) {
	data, err := c.doGet("/tapd_wikis", params)
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse wiki list: %w", err)
	}

	var results []model.Wiki
	for _, item := range rawList {
		if raw, ok := item["Wiki"]; ok {
			var wiki model.Wiki
			if err := json.Unmarshal(raw, &wiki); err == nil {
				results = append(results, wiki)
			}
		}
	}
	return results, nil
}

// GetWiki 获取单个 Wiki 文档详情，description 字段自动从 HTML 转换为 Markdown
// 返回强类型 *model.Wiki，自动过滤无用字段
func (c *Client) GetWiki(workspaceID, id string) (*model.Wiki, error) {
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

	var wiki model.Wiki
	if err := json.Unmarshal(raw, &wiki); err != nil {
		return nil, fmt.Errorf("failed to parse wiki: %w", err)
	}

	if wiki.Description != "" {
		md, err := htmltomarkdown.ConvertString(wiki.Description)
		if err == nil {
			wiki.Description = md
		}
	}
	wiki.URL = fmt.Sprintf("https://www.tapd.cn/%s/markdown_wikis/view/%s", workspaceID, id)

	return &wiki, nil
}
