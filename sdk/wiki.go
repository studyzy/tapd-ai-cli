package tapd

import (
	"encoding/json"
	"fmt"

	"github.com/studyzy/tapd-sdk-go/model"
)

// ListWikis 查询 Wiki 文档列表，返回强类型 []model.Wiki
func (c *Client) ListWikis(req *model.ListWikisRequest) ([]model.Wiki, error) {
	data, err := c.doGet("/tapd_wikis", req.ToParams())
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

// GetWiki 获取单个 Wiki 文档详情，description 字段保留原始 HTML
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

	wiki.URL = fmt.Sprintf("%s/%s/markdown_wikis/view/%s", c.webURL, workspaceID, id)

	return &wiki, nil
}

// CreateWiki 创建 Wiki 文档，返回创建后的完整 Wiki 对象
func (c *Client) CreateWiki(req *model.CreateWikiRequest) (*model.Wiki, error) {
	data, err := c.doPost("/tapd_wikis", req.ToParams())
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse create wiki response: %w", err)
	}

	raw, ok := wrapper["Wiki"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var wiki model.Wiki
	if err := json.Unmarshal(raw, &wiki); err != nil {
		return nil, fmt.Errorf("failed to parse created wiki: %w", err)
	}

	wiki.URL = fmt.Sprintf("%s/%s/markdown_wikis/view/%s", c.webURL, req.WorkspaceID, wiki.ID)

	return &wiki, nil
}

// UpdateWiki 更新 Wiki 文档
func (c *Client) UpdateWiki(req *model.UpdateWikiRequest) (*model.Wiki, error) {
	data, err := c.doPost("/tapd_wikis", req.ToParams())
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse update wiki response: %w", err)
	}

	raw, ok := wrapper["Wiki"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var wiki model.Wiki
	if err := json.Unmarshal(raw, &wiki); err != nil {
		return nil, fmt.Errorf("failed to parse updated wiki: %w", err)
	}

	return &wiki, nil
}
