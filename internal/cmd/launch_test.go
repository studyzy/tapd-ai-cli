package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	tapd "github.com/studyzy/tapd-sdk-go"
)

func TestLaunchCommand_HasSubcommands(t *testing.T) {
	want := map[string]bool{
		"list":      false,
		"count":     false,
		"create":    false,
		"update":    false,
		"templates": false,
		"fields":    false,
	}
	for _, sub := range launchCmd.Commands() {
		if _, ok := want[sub.Name()]; ok {
			want[sub.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("launchCmd should have %q subcommand", name)
		}
	}
}

func TestLaunchCreate_SendsExpectedParams(t *testing.T) {
	withLaunchTestClient(t, func(r *http.Request) string {
		if r.URL.Path != "/launch_forms" {
			t.Errorf("unexpected path: %s, want /launch_forms", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm failed: %v", err)
		}
		checks := map[string]string{
			"workspace_id":     "10104801",
			"creator":          "tapd_api",
			"template_id":      "tpl-1",
			"title":            "上线评审",
			"release_type":     "正常发布",
			"custom_field_one": "value1",
			"cus_自定义字段的名称":     "alias-value",
		}
		for k, want := range checks {
			if got := r.PostFormValue(k); got != want {
				t.Errorf("%s = %q, want %q", k, got, want)
			}
		}
		return `{"status":1,"data":{"LaunchForm":{"id":"lf-1","workspace_id":"10104801","title":"上线评审"}},"info":"success"}`
	}, func() {
		flagLaunchTemplateID = "tpl-1"
		flagCreator = "tapd_api"
		flagLaunchTitle = "上线评审"
		flagLaunchReleaseType = "正常发布"
		flagLaunchCustomFields = map[string]string{"custom_field_one": "value1"}
		flagLaunchAliasFields = map[string]string{"自定义字段的名称": "alias-value"}
		if err := runLaunchCreate(launchCreateCmd, nil); err != nil {
			t.Fatalf("runLaunchCreate() unexpected error: %v", err)
		}
	})
}

func TestLaunchUpdate_SendsExpectedParams(t *testing.T) {
	withLaunchTestClient(t, func(r *http.Request) string {
		if r.URL.Path != "/launch_forms" {
			t.Errorf("unexpected path: %s, want /launch_forms", r.URL.Path)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm failed: %v", err)
		}
		checks := map[string]string{
			"id":               "lf-1",
			"workspace_id":     "10104801",
			"status":           "finished",
			"release_result":   "release_success",
			"release_comment":  "done",
			"remark":           "ok",
			"custom_field_two": "value2",
		}
		for k, want := range checks {
			if got := r.PostFormValue(k); got != want {
				t.Errorf("%s = %q, want %q", k, got, want)
			}
		}
		return `{"status":1,"data":{"LaunchForm":{"id":"lf-1","workspace_id":"10104801","status":"finished"}},"info":"success"}`
	}, func() {
		flagStatus = "finished"
		flagLaunchReleaseResult = "release_success"
		flagLaunchReleaseComment = "done"
		flagLaunchRemark = "ok"
		flagLaunchCustomFields = map[string]string{"custom_field_two": "value2"}
		if err := runLaunchUpdate(launchUpdateCmd, []string{"lf-1"}); err != nil {
			t.Fatalf("runLaunchUpdate() unexpected error: %v", err)
		}
	})
}

func TestLaunchListAndCount_SendsExpectedQuery(t *testing.T) {
	calls := 0
	withLaunchTestClient(t, func(r *http.Request) string {
		calls++
		if r.URL.Query().Get("workspace_id") != "10104801" {
			t.Errorf("workspace_id = %q, want 10104801", r.URL.Query().Get("workspace_id"))
		}
		if r.URL.Query().Get("release_type") != "正常发布" {
			t.Errorf("release_type = %q, want 正常发布", r.URL.Query().Get("release_type"))
		}
		if r.URL.Query().Get("limit") != "5" {
			t.Errorf("limit = %q, want 5", r.URL.Query().Get("limit"))
		}
		if r.URL.Path == "/launch_forms/count" {
			return `{"status":1,"data":{"count":1},"info":"success"}`
		}
		if r.URL.Path != "/launch_forms" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		return `{"status":1,"data":[{"LaunchForm":{"id":"lf-1","workspace_id":"10104801","release_type":"正常发布"}}],"info":"success"}`
	}, func() {
		flagLaunchReleaseType = "正常发布"
		flagLimit = 5
		flagPage = 1
		if err := runLaunchList(launchListCmd, nil); err != nil {
			t.Fatalf("runLaunchList() unexpected error: %v", err)
		}
	})
	if calls != 2 {
		t.Fatalf("expected list and count calls, got %d", calls)
	}
}

func TestLaunchTemplatesAndFields(t *testing.T) {
	tests := []struct {
		name string
		run  func() error
		path string
		body string
	}{
		{
			name: "templates",
			run:  func() error { return runLaunchTemplates(launchTemplatesCmd, nil) },
			path: "/launch_forms/templates",
			body: `{"status":1,"data":[{"template":{"id":"tpl-1","name":"默认流程"}}],"info":"success"}`,
		},
		{
			name: "fields",
			run:  func() error { return runLaunchFields(launchFieldsCmd, nil) },
			path: "/launch_forms/custom_fields_settings",
			body: `{"status":1,"data":[{"CustomFieldConfig":{"id":"cf-1","workspace_id":"10104801","entry_type":"launchform","custom_field":"custom_field_one","name":"DB变更"}}],"info":"success"}`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			withLaunchTestClient(t, func(r *http.Request) string {
				if r.URL.Path != tc.path {
					t.Errorf("unexpected path: %s, want %s", r.URL.Path, tc.path)
				}
				return tc.body
			}, func() {
				if err := tc.run(); err != nil {
					t.Fatalf("run unexpected error: %v", err)
				}
			})
		})
	}
}

func withLaunchTestClient(t *testing.T, handler func(*http.Request) string, run func()) {
	t.Helper()
	resetLaunchFlags()
	oldClient := apiClient
	oldWorkspaceID := flagWorkspaceID
	oldPretty := flagPretty
	t.Cleanup(func() {
		apiClient = oldClient
		flagWorkspaceID = oldWorkspaceID
		flagPretty = oldPretty
		resetLaunchFlags()
	})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(handler(r)))
	}))
	t.Cleanup(srv.Close)

	apiClient = tapd.NewClientWithBaseURL(srv.URL, srv.URL, "test-token", "", "")
	flagWorkspaceID = "10104801"
	flagPretty = false
	run()
}

func resetLaunchFlags() {
	flagLaunchID = ""
	flagCreator = ""
	flagStatus = ""
	flagLaunchTitle = ""
	flagLaunchTemplateID = ""
	flagLaunchVersionType = ""
	flagLaunchBaseline = ""
	flagLaunchReleaseModel = ""
	flagLaunchRoadmapVersion = ""
	flagLaunchReleaseType = ""
	flagLaunchChangeType = ""
	flagLaunchSignedBy = ""
	flagLaunchArchivedBy = ""
	flagLaunchCC = ""
	flagLaunchChangeNotifier = ""
	flagLaunchSignerComment = ""
	flagLaunchReleaseResult = ""
	flagLaunchReleaseComment = ""
	flagLaunchRemark = ""
	flagLaunchFields = ""
	flagLimit = 10
	flagPage = 1
	flagLaunchCustomFields = nil
	flagLaunchAliasFields = nil
}
