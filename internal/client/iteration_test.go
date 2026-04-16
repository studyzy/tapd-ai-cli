package client_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/studyzy/tapd-ai-cli/internal/client"
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

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	params := map[string]string{
		"workspace_id": "1",
	}
	iterations, err := c.ListIterations(params)
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
