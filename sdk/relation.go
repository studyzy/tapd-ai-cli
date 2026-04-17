package tapd

import (
	"encoding/json"

	"github.com/studyzy/tapd-sdk-go/model"
)

// GetRelatedBugs 获取需求关联的缺陷 ID 列表
func (c *Client) GetRelatedBugs(req *model.GetRelatedBugsRequest) (json.RawMessage, error) {
	return c.doGet("/stories/get_related_bugs", req.ToParams())
}

// CreateRelation 创建实体关联关系
func (c *Client) CreateRelation(req *model.CreateRelationRequest) (json.RawMessage, error) {
	return c.doPost("/relations", req.ToParams())
}
