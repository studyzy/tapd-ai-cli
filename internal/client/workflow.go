// Package client 中的 workflow.go 封装了 TAPD 工作流相关 API
package client

import "encoding/json"

// GetWorkflowTransitions 获取工作流状态流转细则
func (c *Client) GetWorkflowTransitions(params map[string]string) (json.RawMessage, error) {
	return c.doGet("/workflows/all_transitions", params)
}

// GetWorkflowStatusMap 获取工作流状态中英文映射
func (c *Client) GetWorkflowStatusMap(params map[string]string) (json.RawMessage, error) {
	return c.doGet("/workflows/status_map", params)
}

// GetWorkflowLastSteps 获取工作流结束状态
func (c *Client) GetWorkflowLastSteps(params map[string]string) (json.RawMessage, error) {
	return c.doGet("/workflows/last_steps", params)
}
