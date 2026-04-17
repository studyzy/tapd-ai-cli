package tapd

import (
	"encoding/json"

	"github.com/studyzy/tapd-sdk-go/model"
)

// GetWorkflowTransitions 获取工作流状态流转细则
func (c *Client) GetWorkflowTransitions(req *model.WorkflowRequest) (json.RawMessage, error) {
	return c.doGet("/workflows/all_transitions", req.ToParams())
}

// GetWorkflowStatusMap 获取工作流状态中英文映射
func (c *Client) GetWorkflowStatusMap(req *model.WorkflowRequest) (json.RawMessage, error) {
	return c.doGet("/workflows/status_map", req.ToParams())
}

// GetWorkflowLastSteps 获取工作流结束状态
func (c *Client) GetWorkflowLastSteps(req *model.WorkflowRequest) (json.RawMessage, error) {
	return c.doGet("/workflows/last_steps", req.ToParams())
}
