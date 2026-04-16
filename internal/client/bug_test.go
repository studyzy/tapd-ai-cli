package client_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/studyzy/tapd-ai-cli/internal/client"
)

func TestListBugs(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bugs" {
			t.Errorf("unexpected path: %s, want /bugs", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Bug":{"id":"500","title":"Bug1","status":"new","priority":"high"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	params := map[string]string{
		"workspace_id": "1",
	}
	results, err := c.ListBugs(params)
	if err != nil {
		t.Fatalf("ListBugs() unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 bug, got %d", len(results))
	}
	if results[0]["id"] != "500" {
		t.Errorf("bug id = %v, want %q", results[0]["id"], "500")
	}
	if results[0]["title"] != "Bug1" {
		t.Errorf("bug title = %v, want %q", results[0]["title"], "Bug1")
	}
	if results[0]["status"] != "new" {
		t.Errorf("bug status = %v, want %q", results[0]["status"], "new")
	}
	if results[0]["priority"] != "high" {
		t.Errorf("bug priority = %v, want %q", results[0]["priority"], "high")
	}
}

func TestGetBug_HTMLToMarkdown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bugs" {
			t.Errorf("unexpected path: %s, want /bugs", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Bug":{"id":"500","title":"Bug1","description":"<p>Steps to <em>reproduce</em></p>"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	result, err := c.GetBug("1", "500")
	if err != nil {
		t.Fatalf("GetBug() unexpected error: %v", err)
	}

	desc, ok := result["description"].(string)
	if !ok {
		t.Fatal("description is not a string")
	}
	if !strings.Contains(desc, "*reproduce*") {
		t.Errorf("description = %q, want to contain %q", desc, "*reproduce*")
	}

	urlStr, ok := result["url"].(string)
	if !ok || urlStr == "" {
		t.Error("url field should be populated")
	}
	if !strings.Contains(urlStr, "/1/bugtrace/bugs/view/500") {
		t.Errorf("url = %q, want to contain %q", urlStr, "/1/bugtrace/bugs/view/500")
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

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	params := map[string]string{
		"workspace_id": "1",
		"title":        "New Bug",
	}
	resp, err := c.CreateBug(params)
	if err != nil {
		t.Fatalf("CreateBug() unexpected error: %v", err)
	}
	if !resp.Success {
		t.Error("expected Success = true")
	}
	if resp.ID != "600" {
		t.Errorf("ID = %q, want %q", resp.ID, "600")
	}
	if !strings.Contains(resp.URL, "/1/bugtrace/bugs/view/600") {
		t.Errorf("URL = %q, want to contain %q", resp.URL, "/1/bugtrace/bugs/view/600")
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

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	params := map[string]string{
		"workspace_id": "1",
	}
	count, err := c.CountBugs(params)
	if err != nil {
		t.Fatalf("CountBugs() unexpected error: %v", err)
	}
	if count != 17 {
		t.Errorf("count = %d, want 17", count)
	}
}
