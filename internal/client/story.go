package client

import (
	"encoding/json"
	"fmt"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/studyzy/tapd-ai-cli/internal/model"
)

// ListStories 查询需求列表，返回强类型 Story 切片，自动过滤 custom_field 等无用字段
func (c *Client) ListStories(req *model.ListStoriesRequest) ([]model.Story, error) {
	data, err := c.doGet("/stories", req.ToParams())
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse story list: %w", err)
	}

	var results []model.Story
	for _, item := range rawList {
		if raw, ok := item["Story"]; ok {
			var story model.Story
			if err := json.Unmarshal(raw, &story); err == nil {
				results = append(results, story)
			}
		}
	}
	return results, nil
}

// GetStory 获取单个需求详情，description 字段自动从 HTML 转换为 Markdown
func (c *Client) GetStory(workspaceID, id string) (*model.Story, error) {
	params := map[string]string{
		"workspace_id": workspaceID,
		"id":           id,
	}

	data, err := c.doGet("/stories", params)
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(rawList) == 0 {
		return nil, &TAPDError{ExitCode: 2, Message: fmt.Sprintf("story %s not found", id)}
	}

	raw, ok := rawList[0]["Story"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var story model.Story
	if err := json.Unmarshal(raw, &story); err != nil {
		return nil, fmt.Errorf("failed to parse story: %w", err)
	}

	// HTML 转 Markdown
	if story.Description != "" {
		md, err := htmltomarkdown.ConvertString(story.Description)
		if err == nil {
			story.Description = md
		}
	}

	story.URL = fmt.Sprintf("https://www.tapd.cn/%s/prong/stories/view/%s", workspaceID, id)

	return &story, nil
}

// CreateStory 创建需求
func (c *Client) CreateStory(req *model.CreateStoryRequest) (*model.SuccessResponse, error) {
	data, err := c.doPost("/stories", req.ToParams())
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse create response: %w", err)
	}

	raw, ok := wrapper["Story"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var created struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(raw, &created); err != nil {
		return nil, fmt.Errorf("failed to parse created story: %w", err)
	}

	return &model.SuccessResponse{
		Success: true,
		ID:      created.ID,
		URL:     fmt.Sprintf("https://www.tapd.cn/%s/prong/stories/view/%s", req.WorkspaceID, created.ID),
	}, nil
}

// UpdateStory 更新需求
func (c *Client) UpdateStory(req *model.UpdateStoryRequest) (*model.Story, error) {
	data, err := c.doPost("/stories", req.ToParams())
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse update response: %w", err)
	}

	raw, ok := wrapper["Story"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var story model.Story
	if err := json.Unmarshal(raw, &story); err != nil {
		return nil, fmt.Errorf("failed to parse updated story: %w", err)
	}

	return &story, nil
}

// CountStories 查询需求数量
func (c *Client) CountStories(req *model.CountStoriesRequest) (int, error) {
	data, err := c.doGet("/stories/count", req.ToParams())
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
