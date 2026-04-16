package client_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/studyzy/tapd-ai-cli/internal/client"
)

func TestListStories(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stories" {
			t.Errorf("unexpected path: %s, want /stories", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Story":{"id":"100","name":"Test Story","status":"open","owner":"user1","modified":"2026-01-01"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	params := map[string]string{
		"workspace_id": "1",
	}
	results, err := c.ListStories(params)
	if err != nil {
		t.Fatalf("ListStories() unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 story, got %d", len(results))
	}
	if results[0]["id"] != "100" {
		t.Errorf("story id = %v, want %q", results[0]["id"], "100")
	}
	if results[0]["name"] != "Test Story" {
		t.Errorf("story name = %v, want %q", results[0]["name"], "Test Story")
	}
	if results[0]["status"] != "open" {
		t.Errorf("story status = %v, want %q", results[0]["status"], "open")
	}
	if results[0]["owner"] != "user1" {
		t.Errorf("story owner = %v, want %q", results[0]["owner"], "user1")
	}
}

func TestGetStory_HTMLToMarkdown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stories" {
			t.Errorf("unexpected path: %s, want /stories", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Story":{"id":"100","name":"Test Story","description":"<p>Hello <strong>World</strong></p>"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	result, err := c.GetStory("1", "100", "stories")
	if err != nil {
		t.Fatalf("GetStory() unexpected error: %v", err)
	}

	desc, ok := result["description"].(string)
	if !ok {
		t.Fatal("description is not a string")
	}
	if !strings.Contains(desc, "**World**") {
		t.Errorf("description = %q, want to contain %q", desc, "**World**")
	}

	urlStr, ok := result["url"].(string)
	if !ok || urlStr == "" {
		t.Error("url field should be populated")
	}
	if !strings.Contains(urlStr, "/1/prong/stories/view/100") {
		t.Errorf("url = %q, want to contain %q", urlStr, "/1/prong/stories/view/100")
	}
}

func TestCreateStory(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/stories" {
			t.Errorf("unexpected path: %s, want /stories", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"Story":{"id":"200","name":"New"}},"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	params := map[string]string{
		"workspace_id": "1",
		"name":         "New",
	}
	resp, err := c.CreateStory(params, "stories")
	if err != nil {
		t.Fatalf("CreateStory() unexpected error: %v", err)
	}
	if !resp.Success {
		t.Error("expected Success = true")
	}
	if resp.ID != "200" {
		t.Errorf("ID = %q, want %q", resp.ID, "200")
	}
	if !strings.Contains(resp.URL, "/1/prong/stories/view/200") {
		t.Errorf("URL = %q, want to contain %q", resp.URL, "/1/prong/stories/view/200")
	}
}

func TestCountStories(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stories/count" {
			t.Errorf("unexpected path: %s, want /stories/count", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"count":42},"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	params := map[string]string{
		"workspace_id": "1",
	}
	count, err := c.CountStories(params)
	if err != nil {
		t.Fatalf("CountStories() unexpected error: %v", err)
	}
	if count != 42 {
		t.Errorf("count = %d, want 42", count)
	}
}

func TestListStories_Tasks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks" {
			t.Errorf("unexpected path: %s, want /tasks", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Task":{"id":"300","name":"Test Task"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	params := map[string]string{
		"workspace_id": "1",
		"entity_type":  "tasks",
	}
	results, err := c.ListStories(params)
	if err != nil {
		t.Fatalf("ListStories(tasks) unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 task, got %d", len(results))
	}
	if results[0]["id"] != "300" {
		t.Errorf("task id = %v, want %q", results[0]["id"], "300")
	}
	if results[0]["name"] != "Test Task" {
		t.Errorf("task name = %v, want %q", results[0]["name"], "Test Task")
	}
}
