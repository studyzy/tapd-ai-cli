package tapd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/studyzy/tapd-sdk-go/model"
)

func TestGetCustomFields(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stories/custom_fields_settings" {
			t.Errorf("unexpected path: %s, want /stories/custom_fields_settings", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"CustomFieldConfig":{"id":"cf1","workspace_id":"1","custom_field":"custom_field_one","name":"自定义字段1","type":"text","enabled":"1"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	req := &model.GetCustomFieldsRequest{
		WorkspaceID: "1",
		EntityType:  "stories",
	}
	configs, err := c.GetCustomFields(req)
	if err != nil {
		t.Fatalf("GetCustomFields() unexpected error: %v", err)
	}
	if len(configs) != 1 {
		t.Fatalf("expected 1 config, got %d", len(configs))
	}
	if configs[0].ID != "cf1" {
		t.Errorf("id = %q, want %q", configs[0].ID, "cf1")
	}
	if configs[0].CustomField != "custom_field_one" {
		t.Errorf("custom_field = %q, want %q", configs[0].CustomField, "custom_field_one")
	}
	if configs[0].Name != "自定义字段1" {
		t.Errorf("name = %q, want %q", configs[0].Name, "自定义字段1")
	}
	if configs[0].Type != "text" {
		t.Errorf("type = %q, want %q", configs[0].Type, "text")
	}
}

func TestGetStoryFieldsLabel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stories/get_fields_lable" {
			t.Errorf("unexpected path: %s, want /stories/get_fields_lable", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"id":"ID","name":"标题","status":"状态","owner":"处理人"},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	req := &model.WorkspaceIDRequest{
		WorkspaceID: "1",
	}
	labels, err := c.GetStoryFieldsLabel(req)
	if err != nil {
		t.Fatalf("GetStoryFieldsLabel() unexpected error: %v", err)
	}
	if len(labels) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(labels))
	}
	if labels["id"] != "ID" {
		t.Errorf("id = %q, want %q", labels["id"], "ID")
	}
	if labels["name"] != "标题" {
		t.Errorf("name = %q, want %q", labels["name"], "标题")
	}
	if labels["status"] != "状态" {
		t.Errorf("status = %q, want %q", labels["status"], "状态")
	}
	if labels["owner"] != "处理人" {
		t.Errorf("owner = %q, want %q", labels["owner"], "处理人")
	}
}

func TestGetStoryFieldsInfo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stories/get_fields_info" {
			t.Errorf("unexpected path: %s, want /stories/get_fields_info", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"status":{"name":"status","label":"状态","html_type":"select"}},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	req := &model.WorkspaceIDRequest{
		WorkspaceID: "1",
	}
	fields, err := c.GetStoryFieldsInfo(req)
	if err != nil {
		t.Fatalf("GetStoryFieldsInfo() unexpected error: %v", err)
	}
	info, ok := fields["status"]
	if !ok {
		t.Fatal("expected 'status' key in fields")
	}
	if info.Name != "status" {
		t.Errorf("name = %q, want %q", info.Name, "status")
	}
	if info.Label != "状态" {
		t.Errorf("label = %q, want %q", info.Label, "状态")
	}
	if info.HTMLType != "select" {
		t.Errorf("html_type = %q, want %q", info.HTMLType, "select")
	}
}

func TestGetWorkitemTypes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/workitem_types" {
			t.Errorf("unexpected path: %s, want /workitem_types", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"WorkitemType":{"id":"wt1","workspace_id":"1","name":"用户故事","entity_type":"story","status":"open"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	req := &model.WorkspaceIDRequest{
		WorkspaceID: "1",
	}
	types, err := c.GetWorkitemTypes(req)
	if err != nil {
		t.Fatalf("GetWorkitemTypes() unexpected error: %v", err)
	}
	if len(types) != 1 {
		t.Fatalf("expected 1 workitem type, got %d", len(types))
	}
	if types[0].ID != "wt1" {
		t.Errorf("id = %q, want %q", types[0].ID, "wt1")
	}
	if types[0].Name != "用户故事" {
		t.Errorf("name = %q, want %q", types[0].Name, "用户故事")
	}
	if types[0].EntityType != "story" {
		t.Errorf("entity_type = %q, want %q", types[0].EntityType, "story")
	}
}
