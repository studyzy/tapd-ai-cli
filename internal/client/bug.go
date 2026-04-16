package client

import (
	"encoding/json"
	"fmt"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/studyzy/tapd-ai-cli/internal/model"
)

// ListBugs 查询缺陷列表
func (c *Client) ListBugs(params map[string]string) ([]map[string]interface{}, error) {
	data, err := c.doGet("/bugs", params)
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse bug list: %w", err)
	}

	var results []map[string]interface{}
	for _, item := range rawList {
		if raw, ok := item["Bug"]; ok {
			var obj map[string]interface{}
			if err := json.Unmarshal(raw, &obj); err == nil {
				results = append(results, obj)
			}
		}
	}
	return results, nil
}

// GetBug 获取单个缺陷详情，description 字段自动从 HTML 转换为 Markdown
func (c *Client) GetBug(workspaceID, id string) (map[string]interface{}, error) {
	params := map[string]string{
		"workspace_id": workspaceID,
		"id":           id,
	}

	data, err := c.doGet("/bugs", params)
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(rawList) == 0 {
		return nil, &TAPDError{ExitCode: 2, Message: fmt.Sprintf("bug %s not found", id)}
	}

	raw, ok := rawList[0]["Bug"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("failed to parse bug: %w", err)
	}

	// HTML 转 Markdown
	if desc, ok := result["description"].(string); ok && desc != "" {
		md, err := htmltomarkdown.ConvertString(desc)
		if err == nil {
			result["description"] = md
		}
	}

	result["url"] = fmt.Sprintf("https://www.tapd.cn/%s/bugtrace/bugs/view/%s", workspaceID, id)

	return result, nil
}

// CreateBug 创建缺陷
func (c *Client) CreateBug(params map[string]string) (*model.SuccessResponse, error) {
	data, err := c.doPost("/bugs", params)
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse create response: %w", err)
	}

	raw, ok := wrapper["Bug"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var created map[string]interface{}
	if err := json.Unmarshal(raw, &created); err != nil {
		return nil, fmt.Errorf("failed to parse created bug: %w", err)
	}

	id := fmt.Sprintf("%v", created["id"])
	wsID := params["workspace_id"]

	return &model.SuccessResponse{
		Success: true,
		ID:      id,
		URL:     fmt.Sprintf("https://www.tapd.cn/%s/bugtrace/bugs/view/%s", wsID, id),
	}, nil
}

// UpdateBug 更新缺陷
func (c *Client) UpdateBug(params map[string]string) (map[string]interface{}, error) {
	data, err := c.doPost("/bugs", params)
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse update response: %w", err)
	}

	raw, ok := wrapper["Bug"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("failed to parse updated bug: %w", err)
	}

	return result, nil
}

// CountBugs 查询缺陷数量
func (c *Client) CountBugs(params map[string]string) (int, error) {
	data, err := c.doGet("/bugs/count", params)
	if err != nil {
		return 0, err
	}

	var result map[string]int
	if err := json.Unmarshal(data, &result); err != nil {
		return 0, fmt.Errorf("failed to parse count response: %w", err)
	}

	if count, ok := result["count"]; ok {
		return count, nil
	}
	return 0, nil
}
