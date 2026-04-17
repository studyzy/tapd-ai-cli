package tapd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/studyzy/tapd-sdk-go/model"
)

func TestGetCommitMsg(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/svn_commits/get_scm_copy_keywords" {
			t.Errorf("unexpected path: %s, want /svn_commits/get_scm_copy_keywords", r.URL.Path)
		}
		if r.URL.Query().Get("workspace_id") != "1" {
			t.Errorf("workspace_id = %q, want %q", r.URL.Query().Get("workspace_id"), "1")
		}
		if r.URL.Query().Get("object_id") != "100" {
			t.Errorf("object_id = %q, want %q", r.URL.Query().Get("object_id"), "100")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"commit_keyword":"--story=100"},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "test-token", "", "")
	result, err := c.GetCommitMsg(&model.GetCommitMsgRequest{
		WorkspaceID: "1",
		ObjectID:    "100",
		Type:        "story",
	})
	if err != nil {
		t.Fatalf("GetCommitMsg() unexpected error: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(result, &data); err != nil {
		t.Fatalf("failed to parse result: %v", err)
	}
	if data["commit_keyword"] != "--story=100" {
		t.Errorf("commit_keyword = %v, want %q", data["commit_keyword"], "--story=100")
	}
}

func TestListReleases(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/releases" {
			t.Errorf("unexpected path: %s, want /releases", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Release":{"id":"10","name":"v1.0","status":"open","workspace_id":"1","creator":"admin","startdate":"2026-01-01","enddate":"2026-06-30"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "test-token", "", "")
	releases, err := c.ListReleases(&model.WorkspaceIDRequest{WorkspaceID: "1"})
	if err != nil {
		t.Fatalf("ListReleases() unexpected error: %v", err)
	}
	if len(releases) != 1 {
		t.Fatalf("expected 1 release, got %d", len(releases))
	}
	if releases[0].ID != "10" {
		t.Errorf("release id = %q, want %q", releases[0].ID, "10")
	}
	if releases[0].Name != "v1.0" {
		t.Errorf("release name = %q, want %q", releases[0].Name, "v1.0")
	}
	if releases[0].Status != "open" {
		t.Errorf("release status = %q, want %q", releases[0].Status, "open")
	}
	if releases[0].Creator != "admin" {
		t.Errorf("release creator = %q, want %q", releases[0].Creator, "admin")
	}
}

func TestListReleases_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "test-token", "", "")
	releases, err := c.ListReleases(&model.WorkspaceIDRequest{WorkspaceID: "1"})
	if err != nil {
		t.Fatalf("ListReleases() unexpected error: %v", err)
	}
	if len(releases) != 0 {
		t.Errorf("expected 0 releases, got %d", len(releases))
	}
}

func TestGetTodoStories(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/user_oauth/get_user_todo_story" {
			t.Errorf("unexpected path: %s, want /user_oauth/get_user_todo_story", r.URL.Path)
		}
		if r.URL.Query().Get("workspace_id") != "1" {
			t.Errorf("workspace_id = %q, want %q", r.URL.Query().Get("workspace_id"), "1")
		}
		if r.URL.Query().Get("limit") != "5" {
			t.Errorf("limit = %q, want %q", r.URL.Query().Get("limit"), "5")
		}
		if r.URL.Query().Get("page") != "2" {
			t.Errorf("page = %q, want %q", r.URL.Query().Get("page"), "2")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Story":{"id":"100","name":"待办需求","status":"open","owner":"user1"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "test-token", "", "")
	stories, err := c.GetTodoStories(&model.GetTodoRequest{
		WorkspaceID: "1",
		EntityType:  "story",
		Limit:       "5",
		Page:        "2",
	})
	if err != nil {
		t.Fatalf("GetTodoStories() unexpected error: %v", err)
	}
	if len(stories) != 1 {
		t.Fatalf("expected 1 story, got %d", len(stories))
	}
	if stories[0].ID != "100" {
		t.Errorf("story id = %q, want %q", stories[0].ID, "100")
	}
	if stories[0].Name != "待办需求" {
		t.Errorf("story name = %q, want %q", stories[0].Name, "待办需求")
	}
	if stories[0].Owner != "user1" {
		t.Errorf("story owner = %q, want %q", stories[0].Owner, "user1")
	}
}

func TestGetTodoStories_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "test-token", "", "")
	stories, err := c.GetTodoStories(&model.GetTodoRequest{
		WorkspaceID: "1",
		EntityType:  "story",
	})
	if err != nil {
		t.Fatalf("GetTodoStories() unexpected error: %v", err)
	}
	if len(stories) != 0 {
		t.Errorf("expected 0 stories, got %d", len(stories))
	}
}

func TestGetTodoTasks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/user_oauth/get_user_todo_task" {
			t.Errorf("unexpected path: %s, want /user_oauth/get_user_todo_task", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Task":{"id":"200","name":"待办任务","status":"open","owner":"user2"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "test-token", "", "")
	tasks, err := c.GetTodoTasks(&model.GetTodoRequest{
		WorkspaceID: "1",
		EntityType:  "task",
	})
	if err != nil {
		t.Fatalf("GetTodoTasks() unexpected error: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].ID != "200" {
		t.Errorf("task id = %q, want %q", tasks[0].ID, "200")
	}
	if tasks[0].Name != "待办任务" {
		t.Errorf("task name = %q, want %q", tasks[0].Name, "待办任务")
	}
}

func TestGetTodoBugs(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/user_oauth/get_user_todo_bug" {
			t.Errorf("unexpected path: %s, want /user_oauth/get_user_todo_bug", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Bug":{"id":"300","title":"待办缺陷","status":"new","current_owner":"user3","severity":"normal"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "test-token", "", "")
	bugs, err := c.GetTodoBugs(&model.GetTodoRequest{
		WorkspaceID: "1",
		EntityType:  "bug",
	})
	if err != nil {
		t.Fatalf("GetTodoBugs() unexpected error: %v", err)
	}
	if len(bugs) != 1 {
		t.Fatalf("expected 1 bug, got %d", len(bugs))
	}
	if bugs[0].ID != "300" {
		t.Errorf("bug id = %q, want %q", bugs[0].ID, "300")
	}
	if bugs[0].Title != "待办缺陷" {
		t.Errorf("bug title = %q, want %q", bugs[0].Title, "待办缺陷")
	}
	if bugs[0].Severity != "normal" {
		t.Errorf("bug severity = %q, want %q", bugs[0].Severity, "normal")
	}
}

func TestSendQiweiMessage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("Content-Type = %q, want %q", ct, "application/json")
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to parse request body: %v", err)
		}
		if payload["msgtype"] != "markdown" {
			t.Errorf("msgtype = %v, want %q", payload["msgtype"], "markdown")
		}
		md, ok := payload["markdown"].(map[string]interface{})
		if !ok {
			t.Fatal("expected markdown field in payload")
		}
		if md["content"] != "hello world" {
			t.Errorf("content = %v, want %q", md["content"], "hello world")
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewClientWithBaseURL("https://api.tapd.cn", "test-token", "", "")
	err := c.SendQiweiMessage(srv.URL, "hello world")
	if err != nil {
		t.Fatalf("SendQiweiMessage() unexpected error: %v", err)
	}
}

func TestSendQiweiMessage_EmptyWebhook(t *testing.T) {
	c := NewClientWithBaseURL("https://api.tapd.cn", "test-token", "", "")
	err := c.SendQiweiMessage("", "hello")
	if err == nil {
		t.Fatal("expected error for empty webhook URL")
	}
}

func TestSendQiweiMessage_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL("https://api.tapd.cn", "test-token", "", "")
	err := c.SendQiweiMessage(srv.URL, "hello")
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
}
