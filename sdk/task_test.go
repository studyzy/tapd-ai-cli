package tapd

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/studyzy/tapd-sdk-go/model"
)

func TestListTasks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks" {
			t.Errorf("unexpected path: %s, want /tasks", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Task":{"id":"300","name":"Test Task","status":"open","owner":"dev1","story_id":"100"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	tasks, err := c.ListTasks(&model.ListTasksRequest{
		WorkspaceID: "1",
	})
	if err != nil {
		t.Fatalf("ListTasks() unexpected error: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].ID != "300" {
		t.Errorf("task id = %q, want %q", tasks[0].ID, "300")
	}
	if tasks[0].Name != "Test Task" {
		t.Errorf("task name = %q, want %q", tasks[0].Name, "Test Task")
	}
	if tasks[0].StoryID != "100" {
		t.Errorf("task story_id = %q, want %q", tasks[0].StoryID, "100")
	}
}

func TestGetTask_PreservesHTML(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks" {
			t.Errorf("unexpected path: %s, want /tasks", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Task":{"id":"300","name":"Test Task","description":"<p>Task <strong>desc</strong></p>","story_id":"100","creator":"dev1","iteration_id":"50"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	task, err := c.GetTask("1", "300")
	if err != nil {
		t.Fatalf("GetTask() unexpected error: %v", err)
	}
	if task.ID != "300" {
		t.Errorf("task id = %q, want %q", task.ID, "300")
	}
	// SDK 保留原始 HTML，不做转换
	if !strings.Contains(task.Description, "<strong>desc</strong>") {
		t.Errorf("description = %q, want to contain %q", task.Description, "<strong>desc</strong>")
	}
	if !strings.Contains(task.URL, "/1/prong/tasks/view/300") {
		t.Errorf("url = %q, want to contain %q", task.URL, "/1/prong/tasks/view/300")
	}
	if task.StoryID != "100" {
		t.Errorf("task story_id = %q, want %q", task.StoryID, "100")
	}
	if task.Creator != "dev1" {
		t.Errorf("creator = %q, want %q", task.Creator, "dev1")
	}
	if task.IterationID != "50" {
		t.Errorf("iteration_id = %q, want %q", task.IterationID, "50")
	}
}

func TestListTasks_PreservesCustomFields(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Task":{"id":"300","name":"Task","custom_field_one":"t1","custom_plan_field_5":"p5"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	tasks, err := c.ListTasks(&model.ListTasksRequest{WorkspaceID: "1"})
	if err != nil {
		t.Fatalf("ListTasks() unexpected error: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].CustomFields["custom_field_one"] != "t1" {
		t.Errorf("custom_field_one = %q, want %q", tasks[0].CustomFields["custom_field_one"], "t1")
	}
	if tasks[0].CustomFields["custom_plan_field_5"] != "p5" {
		t.Errorf("custom_plan_field_5 = %q, want %q", tasks[0].CustomFields["custom_plan_field_5"], "p5")
	}
}

func TestGetTask_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	_, err := c.GetTask("1", "999")
	if err == nil {
		t.Fatal("expected error for not found task")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error = %q, want to contain %q", err.Error(), "not found")
	}
}

func TestCreateTask(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/tasks" {
			t.Errorf("unexpected path: %s, want /tasks", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"Task":{"id":"400","name":"New Task"}},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	task, err := c.CreateTask(&model.CreateTaskRequest{
		WorkspaceID: "1",
		Name:        "New Task",
		StoryID:     "100",
	})
	if err != nil {
		t.Fatalf("CreateTask() unexpected error: %v", err)
	}
	if task.ID != "400" {
		t.Errorf("ID = %q, want %q", task.ID, "400")
	}
	if task.Name != "New Task" {
		t.Errorf("Name = %q, want %q", task.Name, "New Task")
	}
	if !strings.Contains(task.URL, "/1/prong/tasks/view/400") {
		t.Errorf("URL = %q, want to contain %q", task.URL, "/1/prong/tasks/view/400")
	}
}

func TestUpdateTask(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/tasks" {
			t.Errorf("unexpected path: %s, want /tasks", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"Task":{"id":"300","name":"Updated Task","status":"done"}},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	task, err := c.UpdateTask(&model.UpdateTaskRequest{
		WorkspaceID: "1",
		ID:          "300",
		Name:        "Updated Task",
		Status:      "done",
	})
	if err != nil {
		t.Fatalf("UpdateTask() unexpected error: %v", err)
	}
	if task.ID != "300" {
		t.Errorf("task id = %q, want %q", task.ID, "300")
	}
	if task.Name != "Updated Task" {
		t.Errorf("task name = %q, want %q", task.Name, "Updated Task")
	}
	if task.Status != "done" {
		t.Errorf("task status = %q, want %q", task.Status, "done")
	}
}

func TestCountTasks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/count" {
			t.Errorf("unexpected path: %s, want /tasks/count", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"count":15},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	count, err := c.CountTasks(&model.CountTasksRequest{
		WorkspaceID: "1",
	})
	if err != nil {
		t.Fatalf("CountTasks() unexpected error: %v", err)
	}
	if count != 15 {
		t.Errorf("count = %d, want 15", count)
	}
}
