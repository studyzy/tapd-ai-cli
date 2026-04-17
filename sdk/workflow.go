package tapd

import (
	"encoding/json"
	"fmt"

	"github.com/studyzy/tapd-sdk-go/model"
)

// GetWorkflowTransitions 获取工作流状态流转细则，返回强类型 []model.WorkflowTransition
// API 文档：https://open.tapd.cn/document/api-doc/API文档/api_reference/workflow/get_workflow_all_transitions.html
func (c *Client) GetWorkflowTransitions(req *model.WorkflowRequest) ([]model.WorkflowTransition, error) {
	data, err := c.doGet("/workflows/all_transitions", req.ToParams())
	if err != nil {
		return nil, err
	}

	var transitions []model.WorkflowTransition
	if err := json.Unmarshal(data, &transitions); err != nil {
		return nil, fmt.Errorf("failed to parse workflow transitions: %w", err)
	}
	return transitions, nil
}

// GetWorkflowStatusMap 获取工作流状态中英文映射，返回 map[string]string（英文状态名→中文状态名）
// API 文档：https://open.tapd.cn/document/api-doc/API文档/api_reference/workflow/get_workflow_status_map.html
func (c *Client) GetWorkflowStatusMap(req *model.WorkflowRequest) (map[string]string, error) {
	data, err := c.doGet("/workflows/status_map", req.ToParams())
	if err != nil {
		return nil, err
	}

	var statusMap map[string]string
	if err := json.Unmarshal(data, &statusMap); err != nil {
		return nil, fmt.Errorf("failed to parse workflow status map: %w", err)
	}
	return statusMap, nil
}

// GetWorkflowLastSteps 获取工作流结束状态，返回 map[string]string（英文状态名→中文状态名）
// API 文档：https://open.tapd.cn/document/api-doc/API文档/api_reference/workflow/get_workflow_last_steps.html
func (c *Client) GetWorkflowLastSteps(req *model.WorkflowRequest) (map[string]string, error) {
	data, err := c.doGet("/workflows/last_steps", req.ToParams())
	if err != nil {
		return nil, err
	}

	var lastSteps map[string]string
	if err := json.Unmarshal(data, &lastSteps); err != nil {
		return nil, fmt.Errorf("failed to parse workflow last steps: %w", err)
	}
	return lastSteps, nil
}
