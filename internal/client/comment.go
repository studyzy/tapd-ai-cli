package client

import (
	"encoding/json"
	"fmt"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/studyzy/tapd-ai-cli/internal/model"
)

// ListComments 查询评论列表，返回强类型 Comment 切片
func (c *Client) ListComments(req *model.ListCommentsRequest) ([]model.Comment, error) {
	data, err := c.doGet("/comments", req.ToParams())
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse comment list: %w", err)
	}

	var results []model.Comment
	for _, item := range rawList {
		if raw, ok := item["Comment"]; ok {
			var comment model.Comment
			if err := json.Unmarshal(raw, &comment); err == nil {
				// HTML 转 Markdown
				if comment.Description != "" {
					md, err := htmltomarkdown.ConvertString(comment.Description)
					if err == nil {
						comment.Description = md
					}
				}
				results = append(results, comment)
			}
		}
	}
	return results, nil
}

// AddComment 添加评论，返回新建的评论对象
func (c *Client) AddComment(req *model.AddCommentRequest) (*model.Comment, error) {
	data, err := c.doPost("/comments", req.ToParams())
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse create comment response: %w", err)
	}

	raw, ok := wrapper["Comment"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var comment model.Comment
	if err := json.Unmarshal(raw, &comment); err != nil {
		return nil, fmt.Errorf("failed to parse created comment: %w", err)
	}

	return &comment, nil
}

// UpdateComment 更新评论，返回更新后的评论对象
func (c *Client) UpdateComment(req *model.UpdateCommentRequest) (*model.Comment, error) {
	data, err := c.doPost("/comments", req.ToParams())
	if err != nil {
		return nil, err
	}

	// 尝试直接解析（文档示例格式：data 直接是评论对象）
	var comment model.Comment
	if err := json.Unmarshal(data, &comment); err == nil && comment.ID != "" {
		return &comment, nil
	}

	// 兼容：可能有 "Comment" 包裹层
	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse update comment response: %w", err)
	}

	raw, ok := wrapper["Comment"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	if err := json.Unmarshal(raw, &comment); err != nil {
		return nil, fmt.Errorf("failed to parse updated comment: %w", err)
	}

	return &comment, nil
}

// CountComments 查询评论数量
func (c *Client) CountComments(req *model.CountCommentsRequest) (int, error) {
	data, err := c.doGet("/comments/count", req.ToParams())
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
