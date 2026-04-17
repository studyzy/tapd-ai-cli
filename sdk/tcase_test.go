package tapd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/studyzy/tapd-sdk-go/model"
)

func TestListTCases(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tcases" {
			t.Errorf("unexpected path: %s, want /tcases", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Tcase":{"id":"1001","name":"TC1","workspace_id":"1","status":"normal","priority":"high","creator":"tester1","created":"2026-03-10"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	req := &model.ListTCasesRequest{
		WorkspaceID: "1",
	}
	tcases, err := c.ListTCases(req)
	if err != nil {
		t.Fatalf("ListTCases() unexpected error: %v", err)
	}
	if len(tcases) != 1 {
		t.Fatalf("expected 1 tcase, got %d", len(tcases))
	}
	if tcases[0].ID != "1001" {
		t.Errorf("tcase id = %q, want %q", tcases[0].ID, "1001")
	}
	if tcases[0].Name != "TC1" {
		t.Errorf("tcase name = %q, want %q", tcases[0].Name, "TC1")
	}
	if tcases[0].Status != "normal" {
		t.Errorf("tcase status = %q, want %q", tcases[0].Status, "normal")
	}
	if tcases[0].Priority != "high" {
		t.Errorf("tcase priority = %q, want %q", tcases[0].Priority, "high")
	}
	if tcases[0].Creator != "tester1" {
		t.Errorf("tcase creator = %q, want %q", tcases[0].Creator, "tester1")
	}
}

func TestCountTCases(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tcases/count" {
			t.Errorf("unexpected path: %s, want /tcases/count", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"count":25},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	req := &model.CountTCasesRequest{
		WorkspaceID: "1",
	}
	count, err := c.CountTCases(req)
	if err != nil {
		t.Fatalf("CountTCases() unexpected error: %v", err)
	}
	if count != 25 {
		t.Errorf("count = %d, want 25", count)
	}
}

func TestCreateTCase(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/tcases" {
			t.Errorf("unexpected path: %s, want /tcases", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"Tcase":{"id":"1002","name":"New TC","workspace_id":"1","status":"normal"}},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	req := &model.CreateTCaseRequest{
		WorkspaceID: "1",
		Name:        "New TC",
	}
	tc, err := c.CreateTCase(req)
	if err != nil {
		t.Fatalf("CreateTCase() unexpected error: %v", err)
	}
	if tc.ID != "1002" {
		t.Errorf("tcase id = %q, want %q", tc.ID, "1002")
	}
	if tc.Name != "New TC" {
		t.Errorf("tcase name = %q, want %q", tc.Name, "New TC")
	}
}

func TestBatchCreateTCases(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/tcases/batch_save" {
			t.Errorf("unexpected path: %s, want /tcases/batch_save", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"result":"ok"},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	req := &model.BatchCreateTCasesRequest{
		WorkspaceID: "1",
		Data:        `[{"name":"TC1"},{"name":"TC2"}]`,
	}
	result, err := c.BatchCreateTCases(req)
	if err != nil {
		t.Fatalf("BatchCreateTCases() unexpected error: %v", err)
	}
	if result == nil {
		t.Error("expected non-nil result")
	}
}
