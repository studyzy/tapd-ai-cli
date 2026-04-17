package tapd

import (
	"encoding/json"
	"fmt"

	"github.com/studyzy/tapd-sdk-go/model"
)

// ListWorkspaces 获取当前用户参与的项目列表，自动过滤 category 为 organization 的条目
func (c *Client) ListWorkspaces() ([]model.Workspace, error) {
	params := map[string]string{}
	if c.Nick != "" {
		params["nick"] = c.Nick
	}
	data, err := c.doGet("/workspaces/user_participant_projects", params)
	if err != nil {
		return nil, err
	}

	// TAPD 返回格式: [{"Workspace": {...}}, ...]
	var rawList []map[string]model.Workspace
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse workspace list: %w", err)
	}

	var workspaces []model.Workspace
	for _, item := range rawList {
		if ws, ok := item["Workspace"]; ok {
			if ws.Category != "organization" {
				workspaces = append(workspaces, ws)
			}
		}
	}
	return workspaces, nil
}

// GetWorkspaceInfo 获取指定工作区的详细信息
func (c *Client) GetWorkspaceInfo(workspaceID string) (*model.Workspace, error) {
	params := map[string]string{
		"workspace_id": workspaceID,
	}
	data, err := c.doGet("/workspaces/get_workspace_info", params)
	if err != nil {
		return nil, err
	}

	// TAPD 返回格式: {"Workspace": {...}} (单个对象)
	var wrapper map[string]model.Workspace
	if err := json.Unmarshal(data, &wrapper); err != nil {
		// 尝试数组格式: [{"Workspace": {...}}]
		var rawList []map[string]model.Workspace
		if err2 := json.Unmarshal(data, &rawList); err2 != nil {
			return nil, fmt.Errorf("failed to parse workspace info: %w", err)
		}
		if len(rawList) == 0 {
			return nil, &TAPDError{ExitCode: 2, Message: "workspace not found"}
		}
		if ws, ok := rawList[0]["Workspace"]; ok {
			return &ws, nil
		}
		return nil, fmt.Errorf("unexpected response format")
	}

	if ws, ok := wrapper["Workspace"]; ok {
		return &ws, nil
	}
	return nil, &TAPDError{ExitCode: 2, Message: "workspace not found"}
}
