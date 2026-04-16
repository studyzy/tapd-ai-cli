// Package client 中的 custom_field.go 封装了 TAPD 自定义字段和需求字段相关 API
package client

import (
	"encoding/json"

	"github.com/studyzy/tapd-ai-cli/internal/model"
)

// GetCustomFields 获取自定义字段配置
// entityType 取值：stories, tasks, iterations, tcases
func (c *Client) GetCustomFields(req *model.GetCustomFieldsRequest) (json.RawMessage, error) {
	return c.doGet("/"+req.EntityType+"/custom_fields_settings", req.ToParams())
}

// GetStoryFieldsLabel 获取需求所有字段的中英文名
func (c *Client) GetStoryFieldsLabel(req *model.WorkspaceIDRequest) (json.RawMessage, error) {
	return c.doGet("/stories/get_fields_lable", req.ToParams())
}

// GetStoryFieldsInfo 获取需求所有字段及候选值
func (c *Client) GetStoryFieldsInfo(req *model.WorkspaceIDRequest) (json.RawMessage, error) {
	return c.doGet("/stories/get_fields_info", req.ToParams())
}

// GetWorkitemTypes 获取需求类别列表
func (c *Client) GetWorkitemTypes(req *model.WorkspaceIDRequest) (json.RawMessage, error) {
	return c.doGet("/workitem_types", req.ToParams())
}
