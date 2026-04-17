package model

import "testing"

// TestSetOptional 测试 setOptional 辅助函数
func TestSetOptional(t *testing.T) {
	params := map[string]string{}
	setOptional(params, "key1", "val1")
	setOptional(params, "key2", "")
	if params["key1"] != "val1" {
		t.Errorf("expected key1=val1, got %q", params["key1"])
	}
	if _, ok := params["key2"]; ok {
		t.Error("empty value should not be set")
	}
}

// === story.go ===

func TestListStoriesRequest_ToParams(t *testing.T) {
	req := &ListStoriesRequest{
		WorkspaceID:   "100",
		Name:          "test",
		PriorityLabel: "high",
	}
	params := req.ToParams()
	if params["workspace_id"] != "100" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["name"] != "test" {
		t.Errorf("name: got %q", params["name"])
	}
	if params["priority_label"] != "high" {
		t.Errorf("priority_label: got %q", params["priority_label"])
	}
	if _, ok := params["status"]; ok {
		t.Error("empty status should not be in params")
	}
	if _, ok := params["id"]; ok {
		t.Error("empty id should not be in params")
	}
}

func TestCreateStoryRequest_ToParams(t *testing.T) {
	req := &CreateStoryRequest{
		WorkspaceID:   "1",
		Name:          "Test Story",
		PriorityLabel: "high",
		CustomFields:  map[string]string{"custom_field_one": "cf1"},
	}
	params := req.ToParams()
	if params["workspace_id"] != "1" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["name"] != "Test Story" {
		t.Errorf("name: got %q", params["name"])
	}
	if params["priority_label"] != "high" {
		t.Errorf("priority_label: got %q", params["priority_label"])
	}
	if params["custom_field_one"] != "cf1" {
		t.Errorf("custom_field_one: got %q", params["custom_field_one"])
	}
	if _, ok := params["description"]; ok {
		t.Error("empty description should not be in params")
	}
}

func TestUpdateStoryRequest_ToParams(t *testing.T) {
	req := &UpdateStoryRequest{
		WorkspaceID:  "1",
		ID:           "99",
		Status:       "done",
		CustomFields: map[string]string{"custom_field_two": "cf2"},
	}
	params := req.ToParams()
	if params["workspace_id"] != "1" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["id"] != "99" {
		t.Errorf("id: got %q", params["id"])
	}
	if params["status"] != "done" {
		t.Errorf("status: got %q", params["status"])
	}
	if params["custom_field_two"] != "cf2" {
		t.Errorf("custom_field_two: got %q", params["custom_field_two"])
	}
	if _, ok := params["name"]; ok {
		t.Error("empty name should not be in params")
	}
}

func TestCountStoriesRequest_ToParams(t *testing.T) {
	req := &CountStoriesRequest{
		WorkspaceID: "1",
		Status:      "open",
	}
	params := req.ToParams()
	if params["workspace_id"] != "1" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["status"] != "open" {
		t.Errorf("status: got %q", params["status"])
	}

	// 无可选字段
	req2 := &CountStoriesRequest{WorkspaceID: "2"}
	params2 := req2.ToParams()
	if _, ok := params2["status"]; ok {
		t.Error("empty status should not be in params")
	}
}

// === bug.go ===

func TestListBugsRequest_ToParams(t *testing.T) {
	req := &ListBugsRequest{
		WorkspaceID:   "10",
		Title:         "crash",
		PriorityLabel: "urgent",
		Severity:      "fatal",
		Limit:         "50",
	}
	params := req.ToParams()
	if params["workspace_id"] != "10" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["title"] != "crash" {
		t.Errorf("title: got %q", params["title"])
	}
	if params["priority_label"] != "urgent" {
		t.Errorf("priority_label: got %q", params["priority_label"])
	}
	if params["severity"] != "fatal" {
		t.Errorf("severity: got %q", params["severity"])
	}
	if params["limit"] != "50" {
		t.Errorf("limit: got %q", params["limit"])
	}
	if _, ok := params["status"]; ok {
		t.Error("empty status should not be in params")
	}
}

func TestCreateBugRequest_ToParams(t *testing.T) {
	req := &CreateBugRequest{
		WorkspaceID:  "10",
		Title:        "NPE",
		Severity:     "fatal",
		CurrentOwner: "alice",
		CustomFields: map[string]string{"custom_field_three": "cf3"},
	}
	params := req.ToParams()
	if params["workspace_id"] != "10" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["title"] != "NPE" {
		t.Errorf("title: got %q", params["title"])
	}
	if params["severity"] != "fatal" {
		t.Errorf("severity: got %q", params["severity"])
	}
	if params["current_owner"] != "alice" {
		t.Errorf("current_owner: got %q", params["current_owner"])
	}
	if params["custom_field_three"] != "cf3" {
		t.Errorf("custom_field_three: got %q", params["custom_field_three"])
	}
	if _, ok := params["description"]; ok {
		t.Error("empty description should not be in params")
	}
}

func TestUpdateBugRequest_ToParams(t *testing.T) {
	req := &UpdateBugRequest{
		WorkspaceID:  "10",
		ID:           "55",
		Status:       "resolved",
		CustomFields: map[string]string{"custom_field_four": "cf4"},
	}
	params := req.ToParams()
	if params["workspace_id"] != "10" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["id"] != "55" {
		t.Errorf("id: got %q", params["id"])
	}
	if params["status"] != "resolved" {
		t.Errorf("status: got %q", params["status"])
	}
	if params["custom_field_four"] != "cf4" {
		t.Errorf("custom_field_four: got %q", params["custom_field_four"])
	}
	if _, ok := params["title"]; ok {
		t.Error("empty title should not be in params")
	}
}

func TestCountBugsRequest_ToParams(t *testing.T) {
	req := &CountBugsRequest{
		WorkspaceID: "10",
		Status:      "new",
	}
	params := req.ToParams()
	if params["workspace_id"] != "10" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["status"] != "new" {
		t.Errorf("status: got %q", params["status"])
	}
	if _, ok := params["title"]; ok {
		t.Error("empty title should not be in params")
	}
}

// === task.go ===

func TestListTasksRequest_ToParams(t *testing.T) {
	req := &ListTasksRequest{
		WorkspaceID: "20",
		Status:      "open",
		Owner:       "bob",
		Limit:       "100",
	}
	params := req.ToParams()
	if params["workspace_id"] != "20" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["status"] != "open" {
		t.Errorf("status: got %q", params["status"])
	}
	if params["owner"] != "bob" {
		t.Errorf("owner: got %q", params["owner"])
	}
	if params["limit"] != "100" {
		t.Errorf("limit: got %q", params["limit"])
	}
	if _, ok := params["name"]; ok {
		t.Error("empty name should not be in params")
	}
}

func TestCreateTaskRequest_ToParams(t *testing.T) {
	req := &CreateTaskRequest{
		WorkspaceID:  "20",
		Name:         "Impl feature",
		Owner:        "bob",
		CustomFields: map[string]string{"custom_field_five": "cf5"},
	}
	params := req.ToParams()
	if params["workspace_id"] != "20" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["name"] != "Impl feature" {
		t.Errorf("name: got %q", params["name"])
	}
	if params["owner"] != "bob" {
		t.Errorf("owner: got %q", params["owner"])
	}
	if params["custom_field_five"] != "cf5" {
		t.Errorf("custom_field_five: got %q", params["custom_field_five"])
	}
	if _, ok := params["description"]; ok {
		t.Error("empty description should not be in params")
	}
}

func TestUpdateTaskRequest_ToParams(t *testing.T) {
	req := &UpdateTaskRequest{
		WorkspaceID:  "20",
		ID:           "77",
		Status:       "done",
		CustomFields: map[string]string{"custom_field_six": "cf6"},
	}
	params := req.ToParams()
	if params["workspace_id"] != "20" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["id"] != "77" {
		t.Errorf("id: got %q", params["id"])
	}
	if params["status"] != "done" {
		t.Errorf("status: got %q", params["status"])
	}
	if params["custom_field_six"] != "cf6" {
		t.Errorf("custom_field_six: got %q", params["custom_field_six"])
	}
	if _, ok := params["name"]; ok {
		t.Error("empty name should not be in params")
	}
}

func TestCountTasksRequest_ToParams(t *testing.T) {
	req := &CountTasksRequest{
		WorkspaceID: "20",
		Status:      "progressing",
	}
	params := req.ToParams()
	if params["workspace_id"] != "20" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["status"] != "progressing" {
		t.Errorf("status: got %q", params["status"])
	}
	if _, ok := params["name"]; ok {
		t.Error("empty name should not be in params")
	}
}

// === iteration.go ===

func TestListIterationsRequest_ToParams(t *testing.T) {
	req := &ListIterationsRequest{
		WorkspaceID: "30",
		Status:      "open",
		Name:        "Sprint 1",
	}
	params := req.ToParams()
	if params["workspace_id"] != "30" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["status"] != "open" {
		t.Errorf("status: got %q", params["status"])
	}
	if params["name"] != "Sprint 1" {
		t.Errorf("name: got %q", params["name"])
	}
	if _, ok := params["description"]; ok {
		t.Error("empty description should not be in params")
	}
}

func TestCreateIterationRequest_ToParams(t *testing.T) {
	req := &CreateIterationRequest{
		WorkspaceID: "30",
		Name:        "Sprint 2",
		StartDate:   "2025-01-01",
		EndDate:     "2025-01-14",
		Creator:     "admin",
		Description: "test sprint",
	}
	params := req.ToParams()
	if params["workspace_id"] != "30" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["name"] != "Sprint 2" {
		t.Errorf("name: got %q", params["name"])
	}
	if params["startdate"] != "2025-01-01" {
		t.Errorf("startdate: got %q", params["startdate"])
	}
	if params["enddate"] != "2025-01-14" {
		t.Errorf("enddate: got %q", params["enddate"])
	}
	if params["creator"] != "admin" {
		t.Errorf("creator: got %q", params["creator"])
	}
	if params["description"] != "test sprint" {
		t.Errorf("description: got %q", params["description"])
	}
	if _, ok := params["status"]; ok {
		t.Error("empty status should not be in params")
	}
}

func TestUpdateIterationRequest_ToParams(t *testing.T) {
	req := &UpdateIterationRequest{
		WorkspaceID: "30",
		ID:          "88",
		CurrentUser: "admin",
		Name:        "Sprint 2 updated",
	}
	params := req.ToParams()
	if params["workspace_id"] != "30" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["id"] != "88" {
		t.Errorf("id: got %q", params["id"])
	}
	if params["current_user"] != "admin" {
		t.Errorf("current_user: got %q", params["current_user"])
	}
	if params["name"] != "Sprint 2 updated" {
		t.Errorf("name: got %q", params["name"])
	}
	if _, ok := params["description"]; ok {
		t.Error("empty description should not be in params")
	}
}

func TestCountIterationsRequest_ToParams(t *testing.T) {
	req := &CountIterationsRequest{
		WorkspaceID: "30",
		Status:      "done",
	}
	params := req.ToParams()
	if params["workspace_id"] != "30" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["status"] != "done" {
		t.Errorf("status: got %q", params["status"])
	}
	if _, ok := params["name"]; ok {
		t.Error("empty name should not be in params")
	}
}

// === tcase.go ===

func TestListTCasesRequest_ToParams(t *testing.T) {
	req := &ListTCasesRequest{
		WorkspaceID: "40",
		Name:        "login test",
		Priority:    "1",
	}
	params := req.ToParams()
	if params["workspace_id"] != "40" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["name"] != "login test" {
		t.Errorf("name: got %q", params["name"])
	}
	if params["priority"] != "1" {
		t.Errorf("priority: got %q", params["priority"])
	}
	if _, ok := params["status"]; ok {
		t.Error("empty status should not be in params")
	}
}

func TestCountTCasesRequest_ToParams(t *testing.T) {
	req := &CountTCasesRequest{
		WorkspaceID: "40",
		Status:      "normal",
	}
	params := req.ToParams()
	if params["workspace_id"] != "40" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["status"] != "normal" {
		t.Errorf("status: got %q", params["status"])
	}
	if _, ok := params["name"]; ok {
		t.Error("empty name should not be in params")
	}
}

func TestCreateTCaseRequest_ToParams(t *testing.T) {
	req := &CreateTCaseRequest{
		WorkspaceID: "40",
		Name:        "new case",
		Priority:    "2",
		Steps:       "step1;step2",
	}
	params := req.ToParams()
	if params["workspace_id"] != "40" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["name"] != "new case" {
		t.Errorf("name: got %q", params["name"])
	}
	if params["priority"] != "2" {
		t.Errorf("priority: got %q", params["priority"])
	}
	if params["steps"] != "step1;step2" {
		t.Errorf("steps: got %q", params["steps"])
	}
	if _, ok := params["status"]; ok {
		t.Error("empty status should not be in params")
	}
}

func TestBatchCreateTCasesRequest_ToParams(t *testing.T) {
	req := &BatchCreateTCasesRequest{
		WorkspaceID: "40",
		Data:        `[{"name":"case1"}]`,
	}
	params := req.ToParams()
	if params["workspace_id"] != "40" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["data"] != `[{"name":"case1"}]` {
		t.Errorf("data: got %q", params["data"])
	}

	// 无 data
	req2 := &BatchCreateTCasesRequest{WorkspaceID: "40"}
	params2 := req2.ToParams()
	if _, ok := params2["data"]; ok {
		t.Error("empty data should not be in params")
	}
}

// === comment.go ===

func TestListCommentsRequest_ToParams(t *testing.T) {
	req := &ListCommentsRequest{
		WorkspaceID: "50",
		EntryType:   "stories",
		Author:      "charlie",
	}
	params := req.ToParams()
	if params["workspace_id"] != "50" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["entry_type"] != "stories" {
		t.Errorf("entry_type: got %q", params["entry_type"])
	}
	if params["author"] != "charlie" {
		t.Errorf("author: got %q", params["author"])
	}
	if _, ok := params["title"]; ok {
		t.Error("empty title should not be in params")
	}
}

func TestAddCommentRequest_ToParams(t *testing.T) {
	req := &AddCommentRequest{
		WorkspaceID: "50",
		Description: "looks good",
		Author:      "charlie",
		EntryType:   "stories",
		EntryID:     "1001",
		RootID:      "500",
	}
	params := req.ToParams()
	if params["workspace_id"] != "50" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["description"] != "looks good" {
		t.Errorf("description: got %q", params["description"])
	}
	if params["author"] != "charlie" {
		t.Errorf("author: got %q", params["author"])
	}
	if params["entry_type"] != "stories" {
		t.Errorf("entry_type: got %q", params["entry_type"])
	}
	if params["entry_id"] != "1001" {
		t.Errorf("entry_id: got %q", params["entry_id"])
	}
	if params["root_id"] != "500" {
		t.Errorf("root_id: got %q", params["root_id"])
	}
	if _, ok := params["reply_id"]; ok {
		t.Error("empty reply_id should not be in params")
	}
}

func TestUpdateCommentRequest_ToParams(t *testing.T) {
	req := &UpdateCommentRequest{
		WorkspaceID: "50",
		ID:          "600",
		Description: "updated comment",
	}
	params := req.ToParams()
	if params["workspace_id"] != "50" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["id"] != "600" {
		t.Errorf("id: got %q", params["id"])
	}
	if params["description"] != "updated comment" {
		t.Errorf("description: got %q", params["description"])
	}
	if _, ok := params["change_creator"]; ok {
		t.Error("empty change_creator should not be in params")
	}
}

func TestCountCommentsRequest_ToParams(t *testing.T) {
	req := &CountCommentsRequest{
		WorkspaceID: "50",
		EntryType:   "bug",
	}
	params := req.ToParams()
	if params["workspace_id"] != "50" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["entry_type"] != "bug" {
		t.Errorf("entry_type: got %q", params["entry_type"])
	}
	if _, ok := params["author"]; ok {
		t.Error("empty author should not be in params")
	}
}

// === wiki.go ===

func TestListWikisRequest_ToParams(t *testing.T) {
	req := &ListWikisRequest{
		WorkspaceID: "60",
		Creator:     "dave",
		Limit:       "20",
	}
	params := req.ToParams()
	if params["workspace_id"] != "60" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["creator"] != "dave" {
		t.Errorf("creator: got %q", params["creator"])
	}
	if params["limit"] != "20" {
		t.Errorf("limit: got %q", params["limit"])
	}
	if _, ok := params["name"]; ok {
		t.Error("empty name should not be in params")
	}
}

func TestCreateWikiRequest_ToParams(t *testing.T) {
	req := &CreateWikiRequest{
		WorkspaceID:         "60",
		Name:                "API Guide",
		Creator:             "dave",
		MarkdownDescription: "# Hello",
	}
	params := req.ToParams()
	if params["workspace_id"] != "60" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["name"] != "API Guide" {
		t.Errorf("name: got %q", params["name"])
	}
	if params["creator"] != "dave" {
		t.Errorf("creator: got %q", params["creator"])
	}
	if params["markdown_description"] != "# Hello" {
		t.Errorf("markdown_description: got %q", params["markdown_description"])
	}
	if _, ok := params["description"]; ok {
		t.Error("empty description should not be in params")
	}
}

func TestUpdateWikiRequest_ToParams(t *testing.T) {
	req := &UpdateWikiRequest{
		WorkspaceID: "60",
		ID:          "700",
		Name:        "Updated Guide",
	}
	params := req.ToParams()
	if params["workspace_id"] != "60" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["id"] != "700" {
		t.Errorf("id: got %q", params["id"])
	}
	if params["name"] != "Updated Guide" {
		t.Errorf("name: got %q", params["name"])
	}
	if _, ok := params["description"]; ok {
		t.Error("empty description should not be in params")
	}
}

// === timesheet.go ===

func TestListTimesheetsRequest_ToParams(t *testing.T) {
	req := &ListTimesheetsRequest{
		WorkspaceID: "70",
		EntityType:  "story",
		Owner:       "eve",
	}
	params := req.ToParams()
	if params["workspace_id"] != "70" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["entity_type"] != "story" {
		t.Errorf("entity_type: got %q", params["entity_type"])
	}
	if params["owner"] != "eve" {
		t.Errorf("owner: got %q", params["owner"])
	}
	if _, ok := params["entity_id"]; ok {
		t.Error("empty entity_id should not be in params")
	}
}

func TestAddTimesheetRequest_ToParams(t *testing.T) {
	req := &AddTimesheetRequest{
		WorkspaceID: "70",
		EntityType:  "task",
		EntityID:    "800",
		Timespent:   "3h",
		Owner:       "eve",
		Memo:        "coding",
	}
	params := req.ToParams()
	if params["workspace_id"] != "70" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["entity_type"] != "task" {
		t.Errorf("entity_type: got %q", params["entity_type"])
	}
	if params["entity_id"] != "800" {
		t.Errorf("entity_id: got %q", params["entity_id"])
	}
	if params["timespent"] != "3h" {
		t.Errorf("timespent: got %q", params["timespent"])
	}
	if params["owner"] != "eve" {
		t.Errorf("owner: got %q", params["owner"])
	}
	if params["memo"] != "coding" {
		t.Errorf("memo: got %q", params["memo"])
	}
	if _, ok := params["timeremain"]; ok {
		t.Error("empty timeremain should not be in params")
	}
}

func TestUpdateTimesheetRequest_ToParams(t *testing.T) {
	req := &UpdateTimesheetRequest{
		WorkspaceID: "70",
		ID:          "900",
		Timespent:   "5h",
	}
	params := req.ToParams()
	if params["workspace_id"] != "70" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["id"] != "900" {
		t.Errorf("id: got %q", params["id"])
	}
	if params["timespent"] != "5h" {
		t.Errorf("timespent: got %q", params["timespent"])
	}
	if _, ok := params["memo"]; ok {
		t.Error("empty memo should not be in params")
	}
}

// === attachment.go ===

func TestGetImageRequest_ToParams(t *testing.T) {
	req := &GetImageRequest{
		WorkspaceID: "80",
		ImagePath:   "/img/test.png",
	}
	params := req.ToParams()
	if params["workspace_id"] != "80" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["image_path"] != "/img/test.png" {
		t.Errorf("image_path: got %q", params["image_path"])
	}

	// 无 image_path
	req2 := &GetImageRequest{WorkspaceID: "80"}
	params2 := req2.ToParams()
	if _, ok := params2["image_path"]; ok {
		t.Error("empty image_path should not be in params")
	}
}

func TestGetAttachmentsRequest_ToParams(t *testing.T) {
	req := &GetAttachmentsRequest{
		WorkspaceID: "80",
		Type:        "story",
		EntryID:     "1001",
		Limit:       "10",
	}
	params := req.ToParams()
	if params["workspace_id"] != "80" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["type"] != "story" {
		t.Errorf("type: got %q", params["type"])
	}
	if params["entry_id"] != "1001" {
		t.Errorf("entry_id: got %q", params["entry_id"])
	}
	if params["limit"] != "10" {
		t.Errorf("limit: got %q", params["limit"])
	}
	if _, ok := params["page"]; ok {
		t.Error("empty page should not be in params")
	}
}

// === request.go ===

func TestGetCustomFieldsRequest_ToParams(t *testing.T) {
	req := &GetCustomFieldsRequest{
		WorkspaceID: "90",
		EntityType:  "stories",
	}
	params := req.ToParams()
	if params["workspace_id"] != "90" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if len(params) != 1 {
		t.Errorf("expected 1 param, got %d", len(params))
	}
}

func TestWorkspaceIDRequest_ToParams(t *testing.T) {
	req := &WorkspaceIDRequest{
		WorkspaceID: "90",
	}
	params := req.ToParams()
	if params["workspace_id"] != "90" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if len(params) != 1 {
		t.Errorf("expected 1 param, got %d", len(params))
	}
}

func TestWorkflowRequest_ToParams(t *testing.T) {
	req := &WorkflowRequest{
		WorkspaceID:    "90",
		System:         "story",
		WorkitemTypeID: "wt1",
	}
	params := req.ToParams()
	if params["workspace_id"] != "90" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["system"] != "story" {
		t.Errorf("system: got %q", params["system"])
	}
	if params["workitem_type_id"] != "wt1" {
		t.Errorf("workitem_type_id: got %q", params["workitem_type_id"])
	}

	// 无可选字段
	req2 := &WorkflowRequest{WorkspaceID: "90"}
	params2 := req2.ToParams()
	if _, ok := params2["system"]; ok {
		t.Error("empty system should not be in params")
	}
}

func TestGetCommitMsgRequest_ToParams(t *testing.T) {
	req := &GetCommitMsgRequest{
		WorkspaceID: "90",
		ObjectID:    "obj1",
		Type:        "story",
	}
	params := req.ToParams()
	if params["workspace_id"] != "90" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["object_id"] != "obj1" {
		t.Errorf("object_id: got %q", params["object_id"])
	}
	if params["type"] != "story" {
		t.Errorf("type: got %q", params["type"])
	}
	if len(params) != 3 {
		t.Errorf("expected 3 params, got %d", len(params))
	}
}

func TestGetTodoRequest_ToParams(t *testing.T) {
	req := &GetTodoRequest{
		WorkspaceID: "90",
		EntityType:  "story",
		Limit:       "50",
		Page:        "2",
	}
	params := req.ToParams()
	if params["workspace_id"] != "90" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["limit"] != "50" {
		t.Errorf("limit: got %q", params["limit"])
	}
	if params["page"] != "2" {
		t.Errorf("page: got %q", params["page"])
	}
	if _, ok := params["order"]; ok {
		t.Error("empty order should not be in params")
	}
}

func TestGetRelatedBugsRequest_ToParams(t *testing.T) {
	req := &GetRelatedBugsRequest{
		WorkspaceID: "90",
		StoryID:     "s1",
	}
	params := req.ToParams()
	if params["workspace_id"] != "90" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["story_id"] != "s1" {
		t.Errorf("story_id: got %q", params["story_id"])
	}
	if len(params) != 2 {
		t.Errorf("expected 2 params, got %d", len(params))
	}
}

func TestCreateRelationRequest_ToParams(t *testing.T) {
	req := &CreateRelationRequest{
		WorkspaceID: "90",
		SourceType:  "story",
		TargetType:  "bug",
		SourceID:    "s1",
		TargetID:    "b1",
	}
	params := req.ToParams()
	if params["workspace_id"] != "90" {
		t.Errorf("workspace_id: got %q", params["workspace_id"])
	}
	if params["source_type"] != "story" {
		t.Errorf("source_type: got %q", params["source_type"])
	}
	if params["target_type"] != "bug" {
		t.Errorf("target_type: got %q", params["target_type"])
	}
	if params["source_id"] != "s1" {
		t.Errorf("source_id: got %q", params["source_id"])
	}
	if params["target_id"] != "b1" {
		t.Errorf("target_id: got %q", params["target_id"])
	}
	if len(params) != 5 {
		t.Errorf("expected 5 params, got %d", len(params))
	}
}
