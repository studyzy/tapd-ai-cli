// Package client 中的 misc.go 封装了 TAPD 其他杂项 API（提交关键字、发布计划、待办、企业微信）
package client

import (
	"encoding/json"
	"fmt"

	"github.com/studyzy/tapd-ai-cli/internal/model"
)

// GetCommitMsg 获取源码提交关键字
func (c *Client) GetCommitMsg(params map[string]string) (json.RawMessage, error) {
	return c.doGet("/svn_commits/get_scm_copy_keywords", params)
}

// ListReleases 查询发布计划列表
func (c *Client) ListReleases(params map[string]string) ([]model.Release, error) {
	data, err := c.doGet("/releases", params)
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse release list: %w", err)
	}

	var results []model.Release
	for _, item := range rawList {
		if raw, ok := item["Release"]; ok {
			var r model.Release
			if err := json.Unmarshal(raw, &r); err == nil {
				results = append(results, r)
			}
		}
	}
	return results, nil
}

// GetTodo 获取用户待办事项
// entityType 取值：story, bug, task
func (c *Client) GetTodo(params map[string]string) (json.RawMessage, error) {
	entityType := params["entity_type"]
	delete(params, "entity_type")

	endpoint := "/user_oauth/get_user_todo_story"
	switch entityType {
	case "bug":
		endpoint = "/user_oauth/get_user_todo_bug"
	case "task":
		endpoint = "/user_oauth/get_user_todo_task"
	default:
		// 默认使用 story
	}
	return c.doGet(endpoint, params)
}

// SendQiweiMessage 发送消息到企业微信群
// 注意：此功能需要配置企业微信机器人 webhook URL
func (c *Client) SendQiweiMessage(webhookURL, msg string) error {
	if webhookURL == "" {
		return fmt.Errorf("qiwei webhook URL is not configured")
	}

	// 构造请求体
	msgType := "markdown"
	payload := map[string]interface{}{
		"msgtype": msgType,
		msgType: map[string]string{
			"content": msg,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return c.doPostJSON(webhookURL, payloadBytes)
}
