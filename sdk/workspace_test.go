package tapd

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListWorkspaces(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/workspaces/user_participant_projects" {
			t.Errorf("unexpected path: %s, want /workspaces/user_participant_projects", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Workspace":{"id":"1","name":"Project1","status":"active","category":"project"}},{"Workspace":{"id":"2","name":"Org","status":"active","category":"organization"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	workspaces, err := c.ListWorkspaces()
	if err != nil {
		t.Fatalf("ListWorkspaces() unexpected error: %v", err)
	}
	if len(workspaces) != 1 {
		t.Fatalf("expected 1 workspace after filtering organization, got %d", len(workspaces))
	}
	if workspaces[0].ID != "1" {
		t.Errorf("workspace ID = %q, want %q", workspaces[0].ID, "1")
	}
	if workspaces[0].Name != "Project1" {
		t.Errorf("workspace Name = %q, want %q", workspaces[0].Name, "Project1")
	}
	if workspaces[0].Category != "project" {
		t.Errorf("workspace Category = %q, want %q", workspaces[0].Category, "project")
	}
}

func TestGetWorkspaceInfo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/workspaces/get_workspace_info" {
			t.Errorf("unexpected path: %s, want /workspaces/get_workspace_info", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Workspace":{"id":"10","name":"MyProject","status":"active","category":"project","creator":"admin"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	ws, err := c.GetWorkspaceInfo("10")
	if err != nil {
		t.Fatalf("GetWorkspaceInfo() unexpected error: %v", err)
	}
	if ws.ID != "10" {
		t.Errorf("workspace ID = %q, want %q", ws.ID, "10")
	}
	if ws.Name != "MyProject" {
		t.Errorf("workspace Name = %q, want %q", ws.Name, "MyProject")
	}
	if ws.Status != "active" {
		t.Errorf("workspace Status = %q, want %q", ws.Status, "active")
	}
	if ws.Creator != "admin" {
		t.Errorf("workspace Creator = %q, want %q", ws.Creator, "admin")
	}
}

func TestGetWorkspaceInfo_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	_, err := c.GetWorkspaceInfo("999")
	if err == nil {
		t.Fatal("GetWorkspaceInfo() expected error for empty data, got nil")
	}
	var tapdErr *TAPDError
	if !errors.As(err, &tapdErr) {
		t.Fatalf("expected *TAPDError, got %T", err)
	}
	if tapdErr.ExitCode != 2 {
		t.Errorf("ExitCode = %d, want 2", tapdErr.ExitCode)
	}
}
