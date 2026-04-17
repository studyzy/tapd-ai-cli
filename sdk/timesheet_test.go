package tapd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/studyzy/tapd-sdk-go/model"
)

func TestListTimesheets(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/timesheets" {
			t.Errorf("unexpected path: %s, want /timesheets", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Timesheet":{"id":"1001","entity_type":"story","entity_id":"200","timespent":"2h","owner":"tester","spentdate":"2026-04-16","memo":"开发功能"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	results, err := c.ListTimesheets(&model.ListTimesheetsRequest{
		WorkspaceID: "1",
	})
	if err != nil {
		t.Fatalf("ListTimesheets() unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 timesheet, got %d", len(results))
	}
	if results[0].ID != "1001" {
		t.Errorf("timesheet id = %q, want %q", results[0].ID, "1001")
	}
	if results[0].EntityType != "story" {
		t.Errorf("entity_type = %q, want %q", results[0].EntityType, "story")
	}
	if results[0].EntityID != "200" {
		t.Errorf("entity_id = %q, want %q", results[0].EntityID, "200")
	}
	if results[0].Timespent != "2h" {
		t.Errorf("timespent = %q, want %q", results[0].Timespent, "2h")
	}
	if results[0].Owner != "tester" {
		t.Errorf("owner = %q, want %q", results[0].Owner, "tester")
	}
	if results[0].Memo != "开发功能" {
		t.Errorf("memo = %q, want %q", results[0].Memo, "开发功能")
	}
}

func TestListTimesheets_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	results, err := c.ListTimesheets(&model.ListTimesheetsRequest{
		WorkspaceID: "1",
	})
	if err != nil {
		t.Fatalf("ListTimesheets() unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 timesheets, got %d", len(results))
	}
}

func TestAddTimesheet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/timesheets" {
			t.Errorf("unexpected path: %s, want /timesheets", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"Timesheet":{"id":"1002","entity_type":"story","entity_id":"200","timespent":"3h","owner":"tester","spentdate":"2026-04-16","created":"2026-04-16 10:00:00"}},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	result, err := c.AddTimesheet(&model.AddTimesheetRequest{
		WorkspaceID: "1",
		EntityType:  "story",
		EntityID:    "200",
		Timespent:   "3h",
		Owner:       "tester",
	})
	if err != nil {
		t.Fatalf("AddTimesheet() unexpected error: %v", err)
	}
	if result.ID != "1002" {
		t.Errorf("timesheet id = %q, want %q", result.ID, "1002")
	}
	if result.Timespent != "3h" {
		t.Errorf("timespent = %q, want %q", result.Timespent, "3h")
	}
	if result.Owner != "tester" {
		t.Errorf("owner = %q, want %q", result.Owner, "tester")
	}
}

func TestUpdateTimesheet_Wrapped(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"Timesheet":{"id":"1002","timespent":"5h","memo":"updated","modified":"2026-04-16 12:00:00"}},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	result, err := c.UpdateTimesheet(&model.UpdateTimesheetRequest{
		WorkspaceID: "1",
		ID:          "1002",
		Timespent:   "5h",
		Memo:        "updated",
	})
	if err != nil {
		t.Fatalf("UpdateTimesheet() unexpected error: %v", err)
	}
	if result.ID != "1002" {
		t.Errorf("timesheet id = %q, want %q", result.ID, "1002")
	}
	if result.Timespent != "5h" {
		t.Errorf("timespent = %q, want %q", result.Timespent, "5h")
	}
	if result.Memo != "updated" {
		t.Errorf("memo = %q, want %q", result.Memo, "updated")
	}
}

func TestUpdateTimesheet_Direct(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// 直接返回 Timesheet 对象，不包裹在 "Timesheet" key 中
		w.Write([]byte(`{"status":1,"data":{"id":"1002","timespent":"5h","memo":"updated","modified":"2026-04-16 12:00:00"},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	result, err := c.UpdateTimesheet(&model.UpdateTimesheetRequest{
		WorkspaceID: "1",
		ID:          "1002",
		Timespent:   "5h",
	})
	if err != nil {
		t.Fatalf("UpdateTimesheet() unexpected error: %v", err)
	}
	if result.ID != "1002" {
		t.Errorf("timesheet id = %q, want %q", result.ID, "1002")
	}
}
