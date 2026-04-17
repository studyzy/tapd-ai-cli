package tapd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/studyzy/tapd-sdk-go/model"
)

func TestGetWorkflowTransitions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/workflows/all_transitions" {
			t.Errorf("unexpected path: %s, want /workflows/all_transitions", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Name":"开始处理","StepPrevious":"open","StepNext":"progressing","Inform":[{"InformType":"owner","InformId":""}],"Appendfield":[{"DBModel":"Story","FieldName":"owner","Notnull":"1","Sort":"0","DefaultValue":[]}]}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	req := &model.WorkflowRequest{
		WorkspaceID: "1",
		System:      "story",
	}
	transitions, err := c.GetWorkflowTransitions(req)
	if err != nil {
		t.Fatalf("GetWorkflowTransitions() unexpected error: %v", err)
	}
	if len(transitions) != 1 {
		t.Fatalf("expected 1 transition, got %d", len(transitions))
	}
	if transitions[0].Name != "开始处理" {
		t.Errorf("transition name = %q, want %q", transitions[0].Name, "开始处理")
	}
	if transitions[0].StepPrevious != "open" {
		t.Errorf("step previous = %q, want %q", transitions[0].StepPrevious, "open")
	}
	if transitions[0].StepNext != "progressing" {
		t.Errorf("step next = %q, want %q", transitions[0].StepNext, "progressing")
	}
}

func TestGetWorkflowStatusMap(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/workflows/status_map" {
			t.Errorf("unexpected path: %s, want /workflows/status_map", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"open":"未开始","progressing":"进行中","done":"已完成"},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	req := &model.WorkflowRequest{
		WorkspaceID: "1",
		System:      "story",
	}
	statusMap, err := c.GetWorkflowStatusMap(req)
	if err != nil {
		t.Fatalf("GetWorkflowStatusMap() unexpected error: %v", err)
	}
	if len(statusMap) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(statusMap))
	}
	if statusMap["open"] != "未开始" {
		t.Errorf("open = %q, want %q", statusMap["open"], "未开始")
	}
	if statusMap["progressing"] != "进行中" {
		t.Errorf("progressing = %q, want %q", statusMap["progressing"], "进行中")
	}
	if statusMap["done"] != "已完成" {
		t.Errorf("done = %q, want %q", statusMap["done"], "已完成")
	}
}

func TestGetWorkflowLastSteps(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/workflows/last_steps" {
			t.Errorf("unexpected path: %s, want /workflows/last_steps", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"done":"已完成","rejected":"已拒绝"},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "test-token", "", "")
	req := &model.WorkflowRequest{
		WorkspaceID: "1",
		System:      "story",
	}
	lastSteps, err := c.GetWorkflowLastSteps(req)
	if err != nil {
		t.Fatalf("GetWorkflowLastSteps() unexpected error: %v", err)
	}
	if len(lastSteps) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(lastSteps))
	}
	if lastSteps["done"] != "已完成" {
		t.Errorf("done = %q, want %q", lastSteps["done"], "已完成")
	}
	if lastSteps["rejected"] != "已拒绝" {
		t.Errorf("rejected = %q, want %q", lastSteps["rejected"], "已拒绝")
	}
}
