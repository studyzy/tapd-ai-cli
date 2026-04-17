package tapd

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/studyzy/tapd-sdk-go/model"
)

func TestListStories(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stories" {
			t.Errorf("unexpected path: %s, want /stories", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Story":{"id":"100","name":"Test Story","status":"open","owner":"user1","modified":"2026-01-01","creator":"admin","iteration_id":"50"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	req := &model.ListStoriesRequest{
		WorkspaceID: "1",
	}
	stories, err := c.ListStories(req)
	if err != nil {
		t.Fatalf("ListStories() unexpected error: %v", err)
	}
	if len(stories) != 1 {
		t.Fatalf("expected 1 story, got %d", len(stories))
	}
	if stories[0].ID != "100" {
		t.Errorf("story id = %v, want %q", stories[0].ID, "100")
	}
	if stories[0].Name != "Test Story" {
		t.Errorf("story name = %v, want %q", stories[0].Name, "Test Story")
	}
	if stories[0].Status != "open" {
		t.Errorf("story status = %v, want %q", stories[0].Status, "open")
	}
	if stories[0].Owner != "user1" {
		t.Errorf("story owner = %v, want %q", stories[0].Owner, "user1")
	}
	if stories[0].Creator != "admin" {
		t.Errorf("story creator = %v, want %q", stories[0].Creator, "admin")
	}
	if stories[0].IterationID != "50" {
		t.Errorf("story iteration_id = %v, want %q", stories[0].IterationID, "50")
	}
}

func TestGetStory_PreservesHTML(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stories" {
			t.Errorf("unexpected path: %s, want /stories", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Story":{"id":"100","name":"Test Story","description":"<p>Hello <strong>World</strong></p>","creator":"admin","created":"2026-01-01 10:00:00"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	story, err := c.GetStory("1", "100")
	if err != nil {
		t.Fatalf("GetStory() unexpected error: %v", err)
	}

	// SDK 保留原始 HTML，不做转换
	if !strings.Contains(story.Description, "<strong>World</strong>") {
		t.Errorf("description = %q, want to contain %q", story.Description, "<strong>World</strong>")
	}

	if story.URL == "" {
		t.Error("url field should be populated")
	}
	if !strings.Contains(story.URL, "/1/prong/stories/view/100") {
		t.Errorf("url = %q, want to contain %q", story.URL, "/1/prong/stories/view/100")
	}
	if story.Creator != "admin" {
		t.Errorf("creator = %q, want %q", story.Creator, "admin")
	}
	if story.Created != "2026-01-01 10:00:00" {
		t.Errorf("created = %q, want %q", story.Created, "2026-01-01 10:00:00")
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

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	req := &model.CreateStoryRequest{
		WorkspaceID: "1",
		Name:        "New",
	}
	story, err := c.CreateStory(req)
	if err != nil {
		t.Fatalf("CreateStory() unexpected error: %v", err)
	}
	if story.ID != "200" {
		t.Errorf("ID = %q, want %q", story.ID, "200")
	}
	if story.Name != "New" {
		t.Errorf("Name = %q, want %q", story.Name, "New")
	}
	if !strings.Contains(story.URL, "/1/prong/stories/view/200") {
		t.Errorf("URL = %q, want to contain %q", story.URL, "/1/prong/stories/view/200")
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

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	req := &model.CountStoriesRequest{
		WorkspaceID: "1",
	}
	count, err := c.CountStories(req)
	if err != nil {
		t.Fatalf("CountStories() unexpected error: %v", err)
	}
	if count != 42 {
		t.Errorf("count = %d, want 42", count)
	}
}

func TestUpdateStory(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/stories" {
			t.Errorf("unexpected path: %s, want /stories", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"Story":{"id":"100","name":"Updated","status":"done"}},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	req := &model.UpdateStoryRequest{
		WorkspaceID: "1",
		ID:          "100",
		Name:        "Updated",
		Status:      "done",
	}
	story, err := c.UpdateStory(req)
	if err != nil {
		t.Fatalf("UpdateStory() unexpected error: %v", err)
	}
	if story.ID != "100" {
		t.Errorf("story id = %q, want %q", story.ID, "100")
	}
	if story.Name != "Updated" {
		t.Errorf("story name = %q, want %q", story.Name, "Updated")
	}
}
