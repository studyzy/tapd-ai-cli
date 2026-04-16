package client_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/studyzy/tapd-ai-cli/internal/client"
	"github.com/studyzy/tapd-ai-cli/internal/model"
)

func TestUpdateStory(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/stories" {
			t.Errorf("unexpected path: %s, want /stories", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":1,"data":{"Story":{"id":"100","name":"Updated","status":"done"}},"info":"success"}`)
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	req := &model.UpdateStoryRequest{
		WorkspaceID: "1",
		ID:          "100",
		EntityType:  "stories",
		Name:        "Updated",
		Status:      "done",
	}
	result, err := c.UpdateStory(req)
	if err != nil {
		t.Fatalf("UpdateStory() unexpected error: %v", err)
	}
	story, ok := result.(*model.Story)
	if !ok {
		t.Fatalf("expected *model.Story, got %T", result)
	}
	if story.ID != "100" {
		t.Errorf("id = %v, want %q", story.ID, "100")
	}
	if story.Name != "Updated" {
		t.Errorf("name = %v, want %q", story.Name, "Updated")
	}
	if story.Status != "done" {
		t.Errorf("status = %v, want %q", story.Status, "done")
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
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":1,"data":{"Bug":{"id":"500","title":"Fixed","status":"resolved"}},"info":"success"}`)
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
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

func TestNewClient(t *testing.T) {
	c := client.NewClient("test-token", "", "")
	if c == nil {
		t.Fatal("NewClient() returned nil")
	}
}

func TestTAPDError_Error(t *testing.T) {
	e := &client.TAPDError{
		HTTPStatus: 400,
		ExitCode:   3,
		Message:    "bad request",
	}
	got := e.Error()
	if got != "bad request" {
		t.Errorf("Error() = %q, want %q", got, "bad request")
	}
}
