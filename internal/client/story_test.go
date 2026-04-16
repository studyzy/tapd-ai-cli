package client_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/studyzy/tapd-ai-cli/internal/client"
	"github.com/studyzy/tapd-ai-cli/internal/model"
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

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	req := &model.ListStoriesRequest{
		WorkspaceID: "1",
	}
	result, err := c.ListStories(req)
	if err != nil {
		t.Fatalf("ListStories() unexpected error: %v", err)
	}
	stories, ok := result.([]model.Story)
	if !ok {
		t.Fatalf("expected []model.Story, got %T", result)
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

func TestGetStory_HTMLToMarkdown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stories" {
			t.Errorf("unexpected path: %s, want /stories", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Story":{"id":"100","name":"Test Story","description":"<p>Hello <strong>World</strong></p>","creator":"admin","created":"2026-01-01 10:00:00"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	result, err := c.GetStory("1", "100", "stories")
	if err != nil {
		t.Fatalf("GetStory() unexpected error: %v", err)
	}

	story, ok := result.(*model.Story)
	if !ok {
		t.Fatalf("expected *model.Story, got %T", result)
	}

	if !strings.Contains(story.Description, "**World**") {
		t.Errorf("description = %q, want to contain %q", story.Description, "**World**")
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

func TestGetStory_Task(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks" {
			t.Errorf("unexpected path: %s, want /tasks", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Task":{"id":"300","name":"Test Task","description":"<p>Task desc</p>","story_id":"100","creator":"dev1","iteration_id":"50"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	result, err := c.GetStory("1", "300", "tasks")
	if err != nil {
		t.Fatalf("GetStory(tasks) unexpected error: %v", err)
	}

	task, ok := result.(*model.Task)
	if !ok {
		t.Fatalf("expected *model.Task, got %T", result)
	}
	if task.ID != "300" {
		t.Errorf("task id = %q, want %q", task.ID, "300")
	}
	if task.StoryID != "100" {
		t.Errorf("task story_id = %q, want %q", task.StoryID, "100")
	}
	if !strings.Contains(task.URL, "/1/prong/tasks/view/300") {
		t.Errorf("url = %q, want to contain %q", task.URL, "/1/prong/tasks/view/300")
	}
	if task.Creator != "dev1" {
		t.Errorf("creator = %q, want %q", task.Creator, "dev1")
	}
	if task.IterationID != "50" {
		t.Errorf("iteration_id = %q, want %q", task.IterationID, "50")
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
	req := &model.CreateStoryRequest{
		WorkspaceID: "1",
		Name:        "New",
		EntityType:  "stories",
	}
	resp, err := c.CreateStory(req)
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

func TestListStories_Tasks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks" {
			t.Errorf("unexpected path: %s, want /tasks", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Task":{"id":"300","name":"Test Task","story_id":"100"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	req := &model.ListStoriesRequest{
		WorkspaceID: "1",
		EntityType:  "tasks",
	}
	result, err := c.ListStories(req)
	if err != nil {
		t.Fatalf("ListStories(tasks) unexpected error: %v", err)
	}
	tasks, ok := result.([]model.Task)
	if !ok {
		t.Fatalf("expected []model.Task, got %T", result)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].ID != "300" {
		t.Errorf("task id = %v, want %q", tasks[0].ID, "300")
	}
	if tasks[0].Name != "Test Task" {
		t.Errorf("task name = %v, want %q", tasks[0].Name, "Test Task")
	}
	if tasks[0].StoryID != "100" {
		t.Errorf("task story_id = %v, want %q", tasks[0].StoryID, "100")
	}
}
