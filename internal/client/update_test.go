package client_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/studyzy/tapd-ai-cli/internal/client"
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
	params := map[string]string{
		"workspace_id": "1",
		"id":           "100",
		"name":         "Updated",
		"status":       "done",
	}
	result, err := c.UpdateStory(params, "stories")
	if err != nil {
		t.Fatalf("UpdateStory() unexpected error: %v", err)
	}
	if result["id"] != "100" {
		t.Errorf("id = %v, want %q", result["id"], "100")
	}
	if result["name"] != "Updated" {
		t.Errorf("name = %v, want %q", result["name"], "Updated")
	}
	if result["status"] != "done" {
		t.Errorf("status = %v, want %q", result["status"], "done")
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
	params := map[string]string{
		"workspace_id": "1",
		"id":           "500",
		"title":        "Fixed",
		"status":       "resolved",
	}
	result, err := c.UpdateBug(params)
	if err != nil {
		t.Fatalf("UpdateBug() unexpected error: %v", err)
	}
	if result["id"] != "500" {
		t.Errorf("id = %v, want %q", result["id"], "500")
	}
	if result["title"] != "Fixed" {
		t.Errorf("title = %v, want %q", result["title"], "Fixed")
	}
	if result["status"] != "resolved" {
		t.Errorf("status = %v, want %q", result["status"], "resolved")
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
