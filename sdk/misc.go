package tapd

import (
	"encoding/json"
	"fmt"

	"github.com/studyzy/tapd-sdk-go/model"
)

// GetCommitMsg 获取源码提交关键字
func (c *Client) GetCommitMsg(req *model.GetCommitMsgRequest) (json.RawMessage, error) {
	return c.doGet("/svn_commits/get_scm_copy_keywords", req.ToParams())
}

// ListReleases 查询发布计划列表
func (c *Client) ListReleases(req *model.WorkspaceIDRequest) ([]model.Release, error) {
	data, err := c.doGet("/releases", req.ToParams())
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

// GetTodoStories 获取用户待办需求，返回强类型 Story 切片
func (c *Client) GetTodoStories(req *model.GetTodoRequest) ([]model.Story, error) {
	data, err := c.doGet("/user_oauth/get_user_todo_story", req.ToParams())
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse todo story list: %w", err)
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

// GetTodoTasks 获取用户待办任务，返回强类型 Task 切片
func (c *Client) GetTodoTasks(req *model.GetTodoRequest) ([]model.Task, error) {
	data, err := c.doGet("/user_oauth/get_user_todo_task", req.ToParams())
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse todo task list: %w", err)
	}

	var results []model.Task
	for _, item := range rawList {
		if raw, ok := item["Task"]; ok {
			var task model.Task
			if err := json.Unmarshal(raw, &task); err == nil {
				results = append(results, task)
			}
		}
	}
	return results, nil
}

// GetTodoBugs 获取用户待办缺陷，返回强类型 Bug 切片
func (c *Client) GetTodoBugs(req *model.GetTodoRequest) ([]model.Bug, error) {
	data, err := c.doGet("/user_oauth/get_user_todo_bug", req.ToParams())
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse todo bug list: %w", err)
	}

	var results []model.Bug
	for _, item := range rawList {
		if raw, ok := item["Bug"]; ok {
			var bug model.Bug
			if err := json.Unmarshal(raw, &bug); err == nil {
				results = append(results, bug)
			}
		}
	}
	return results, nil
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
