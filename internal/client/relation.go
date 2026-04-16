// Package client 中的 relation.go 封装了 TAPD 关联关系相关 API
package client

import (
	"encoding/json"

	"github.com/studyzy/tapd-ai-cli/internal/model"
)

// GetRelatedBugs 获取需求关联的缺陷 ID 列表
func (c *Client) GetRelatedBugs(req *model.GetRelatedBugsRequest) (json.RawMessage, error) {
	return c.doGet("/stories/get_related_bugs", req.ToParams())
}

// CreateRelation 创建实体关联关系
func (c *Client) CreateRelation(req *model.CreateRelationRequest) (json.RawMessage, error) {
	return c.doPost("/relations", req.ToParams())
}
