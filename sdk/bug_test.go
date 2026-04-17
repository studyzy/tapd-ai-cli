package tapd

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/studyzy/tapd-sdk-go/model"
)

func TestListBugs(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bugs" {
			t.Errorf("unexpected path: %s, want /bugs", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Bug":{"id":"500","title":"Bug1","status":"new","priority":"high","current_owner":"test","created":"2026-03-06"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	results, err := c.ListBugs(&model.ListBugsRequest{WorkspaceID: "1"})
	if err != nil {
		t.Fatalf("ListBugs() unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 bug, got %d", len(results))
	}
	if results[0].ID != "500" {
		t.Errorf("bug id = %v, want %q", results[0].ID, "500")
	}
	if results[0].Title != "Bug1" {
		t.Errorf("bug title = %v, want %q", results[0].Title, "Bug1")
	}
	if results[0].Status != "new" {
		t.Errorf("bug status = %v, want %q", results[0].Status, "new")
	}
	if results[0].Priority != "high" {
		t.Errorf("bug priority = %v, want %q", results[0].Priority, "high")
	}
	if results[0].CurrentOwner != "test" {
		t.Errorf("bug current_owner = %v, want %q", results[0].CurrentOwner, "test")
	}
	if results[0].Created != "2026-03-06" {
		t.Errorf("bug created = %v, want %q", results[0].Created, "2026-03-06")
	}
}

func TestListBugs_FiltersCustomFields(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// 返回包含大量 custom_field 的数据，强类型反序列化应自动忽略
		w.Write([]byte(`{"status":1,"data":[{"Bug":{"id":"501","title":"Bug2","custom_field_1":"val1","custom_field_50":"val50","custom_field_100":"val100"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	results, err := c.ListBugs(&model.ListBugsRequest{WorkspaceID: "1"})
	if err != nil {
		t.Fatalf("ListBugs() unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 bug, got %d", len(results))
	}
	if results[0].ID != "501" {
		t.Errorf("bug id = %v, want %q", results[0].ID, "501")
	}
	if results[0].Title != "Bug2" {
		t.Errorf("bug title = %v, want %q", results[0].Title, "Bug2")
	}
}

func TestGetBug_PreservesHTML(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bugs" {
			t.Errorf("unexpected path: %s, want /bugs", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Bug":{"id":"500","title":"Bug1","description":"<p>Steps to <em>reproduce</em></p>","current_owner":"DevinZeng","reporter":"tester","created":"2026-03-06 15:35:27","resolved":"2026-03-09 18:36:10","fixer":"DevinZeng","severity":"high","resolution":"fixed","iteration_id":"123"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	result, err := c.GetBug("1", "500")
	if err != nil {
		t.Fatalf("GetBug() unexpected error: %v", err)
	}

	// SDK 保留原始 HTML，不做转换
	if !strings.Contains(result.Description, "<em>reproduce</em>") {
		t.Errorf("description = %q, want to contain %q", result.Description, "<em>reproduce</em>")
	}

	if result.URL == "" {
		t.Error("url field should be populated")
	}
	if !strings.Contains(result.URL, "/1/bugtrace/bugs/view/500") {
		t.Errorf("url = %q, want to contain %q", result.URL, "/1/bugtrace/bugs/view/500")
	}

	// 验证新增字段正确映射
	if result.CurrentOwner != "DevinZeng" {
		t.Errorf("current_owner = %q, want %q", result.CurrentOwner, "DevinZeng")
	}
	if result.Reporter != "tester" {
		t.Errorf("reporter = %q, want %q", result.Reporter, "tester")
	}
	if result.Created != "2026-03-06 15:35:27" {
		t.Errorf("created = %q, want %q", result.Created, "2026-03-06 15:35:27")
	}
	if result.Resolved != "2026-03-09 18:36:10" {
		t.Errorf("resolved = %q, want %q", result.Resolved, "2026-03-09 18:36:10")
	}
	if result.Fixer != "DevinZeng" {
		t.Errorf("fixer = %q, want %q", result.Fixer, "DevinZeng")
	}
	if result.Resolution != "fixed" {
		t.Errorf("resolution = %q, want %q", result.Resolution, "fixed")
	}
	if result.IterationID != "123" {
		t.Errorf("iteration_id = %q, want %q", result.IterationID, "123")
	}
}

func TestCreateBug(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/bugs" {
			t.Errorf("unexpected path: %s, want /bugs", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"Bug":{"id":"600","title":"New Bug"}},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	bug, err := c.CreateBug(&model.CreateBugRequest{
		WorkspaceID: "1",
		Title:       "New Bug",
	})
	if err != nil {
		t.Fatalf("CreateBug() unexpected error: %v", err)
	}
	if bug.ID != "600" {
		t.Errorf("ID = %q, want %q", bug.ID, "600")
	}
	if bug.Title != "New Bug" {
		t.Errorf("Title = %q, want %q", bug.Title, "New Bug")
	}
	if !strings.Contains(bug.URL, "/1/bugtrace/bugs/view/600") {
		t.Errorf("URL = %q, want to contain %q", bug.URL, "/1/bugtrace/bugs/view/600")
	}
}

func TestCountBugs(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bugs/count" {
			t.Errorf("unexpected path: %s, want /bugs/count", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"count":17},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	count, err := c.CountBugs(&model.CountBugsRequest{WorkspaceID: "1"})
	if err != nil {
		t.Fatalf("CountBugs() unexpected error: %v", err)
	}
	if count != 17 {
		t.Errorf("count = %d, want 17", count)
	}
}

func TestUpdateBug(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/bugs" {
			t.Errorf("unexpected path: %s, want /bugs", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"Bug":{"id":"500","title":"Fixed","status":"resolved"}},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	result, err := c.UpdateBug(&model.UpdateBugRequest{
		WorkspaceID: "1",
		ID:          "500",
		Title:       "Fixed",
		Status:      "resolved",
	})
	if err != nil {
		t.Fatalf("UpdateBug() unexpected error: %v", err)
	}
	if result.ID != "500" {
		t.Errorf("id = %v, want %q", result.ID, "500")
	}
	if result.Title != "Fixed" {
		t.Errorf("title = %v, want %q", result.Title, "Fixed")
	}
	if result.Status != "resolved" {
		t.Errorf("status = %v, want %q", result.Status, "resolved")
	}
}
