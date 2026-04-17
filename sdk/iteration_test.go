package tapd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/studyzy/tapd-sdk-go/model"
)

func TestListIterations(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/iterations" {
			t.Errorf("unexpected path: %s, want /iterations", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":1,"data":[{"Iteration":{"id":"3001","name":"Sprint 1","status":"open","startdate":"2026-04-01","enddate":"2026-04-15"}}],"info":"success"}`)
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	iterations, err := c.ListIterations(&model.ListIterationsRequest{
		WorkspaceID: "1",
	})
	if err != nil {
		t.Fatalf("ListIterations() unexpected error: %v", err)
	}
	if len(iterations) != 1 {
		t.Fatalf("expected 1 iteration, got %d", len(iterations))
	}
	iter := iterations[0]
	if iter.ID != "3001" {
		t.Errorf("iteration ID = %q, want %q", iter.ID, "3001")
	}
	if iter.Name != "Sprint 1" {
		t.Errorf("iteration Name = %q, want %q", iter.Name, "Sprint 1")
	}
	if iter.Status != "open" {
		t.Errorf("iteration Status = %q, want %q", iter.Status, "open")
	}
	if iter.StartDate != "2026-04-01" {
		t.Errorf("iteration StartDate = %q, want %q", iter.StartDate, "2026-04-01")
	}
	if iter.EndDate != "2026-04-15" {
		t.Errorf("iteration EndDate = %q, want %q", iter.EndDate, "2026-04-15")
	}
}

func TestCreateIteration(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/iterations" {
			t.Errorf("unexpected path: %s, want /iterations", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":1,"data":{"Iteration":{"id":"3002","name":"Sprint 2","workspace_id":"1","startdate":"2026-04-16","enddate":"2026-04-30","status":"open","creator":"testuser","created":"2026-04-16 10:00:00"}},"info":"success"}`)
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	iteration, err := c.CreateIteration(&model.CreateIterationRequest{
		WorkspaceID: "1",
		Name:        "Sprint 2",
		StartDate:   "2026-04-16",
		EndDate:     "2026-04-30",
		Creator:     "testuser",
	})
	if err != nil {
		t.Fatalf("CreateIteration() unexpected error: %v", err)
	}
	if iteration.ID != "3002" {
		t.Errorf("ID = %q, want %q", iteration.ID, "3002")
	}
	if iteration.Name != "Sprint 2" {
		t.Errorf("Name = %q, want %q", iteration.Name, "Sprint 2")
	}
	if iteration.WorkspaceID != "1" {
		t.Errorf("WorkspaceID = %q, want %q", iteration.WorkspaceID, "1")
	}
}

func TestUpdateIteration(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/iterations" {
			t.Errorf("unexpected path: %s, want /iterations", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":1,"data":{"Iteration":{"id":"3001","name":"Sprint 1 Updated","workspace_id":"1","startdate":"2026-04-01","enddate":"2026-04-15","status":"done","description":"updated desc","creator":"testuser","modified":"2026-04-16 12:00:00"}},"info":"success"}`)
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	result, err := c.UpdateIteration(&model.UpdateIterationRequest{
		WorkspaceID: "1",
		ID:          "3001",
		CurrentUser: "testuser",
		Name:        "Sprint 1 Updated",
		Status:      "done",
		Description: "updated desc",
	})
	if err != nil {
		t.Fatalf("UpdateIteration() unexpected error: %v", err)
	}
	if result.ID != "3001" {
		t.Errorf("ID = %q, want %q", result.ID, "3001")
	}
	if result.Name != "Sprint 1 Updated" {
		t.Errorf("Name = %q, want %q", result.Name, "Sprint 1 Updated")
	}
	if result.Status != "done" {
		t.Errorf("Status = %q, want %q", result.Status, "done")
	}
	if result.Description != "updated desc" {
		t.Errorf("Description = %q, want %q", result.Description, "updated desc")
	}
}
