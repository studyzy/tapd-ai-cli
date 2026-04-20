package cmd

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	tapd "github.com/studyzy/tapd-sdk-go"
)

func TestBuildSpecLines_StoryUpdateIncludesIterationID(t *testing.T) {
	lines := buildSpecLines(rootCmd)
	for _, line := range lines {
		if extractCommandPath(line.text) != "tapd story update" {
			continue
		}
		if !strings.Contains(line.text, "[--iteration-id]") {
			t.Fatalf("story update spec line should contain iteration flag, got: %s", line.text)
		}
		return
	}

	t.Fatal("tapd story update spec line not found")
}

func TestRunStoryUpdate_PassesIterationID(t *testing.T) {
	t.Helper()

	oldClient := apiClient
	oldWorkspaceID := flagWorkspaceID
	oldPretty := flagPretty
	oldName := flagName
	oldDescription := flagDescription
	oldStatus := flagStatus
	oldOwner := flagOwner
	oldPriority := flagPriority
	oldIterationID := flagIterationID
	oldStdout := os.Stdout

	t.Cleanup(func() {
		apiClient = oldClient
		flagWorkspaceID = oldWorkspaceID
		flagPretty = oldPretty
		flagName = oldName
		flagDescription = oldDescription
		flagStatus = oldStatus
		flagOwner = oldOwner
		flagPriority = oldPriority
		flagIterationID = oldIterationID
		os.Stdout = oldStdout
	})

	var captured url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/stories" {
			t.Fatalf("expected path /stories, got %s", r.URL.Path)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm failed: %v", err)
		}
		captured = r.PostForm
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":1,"data":{"Story":{"id":"10001","name":"Updated"}},"info":"success"}`))
	}))
	defer srv.Close()

	apiClient = tapd.NewClientWithBaseURL(srv.URL, srv.URL, "test-token", "", "")
	flagWorkspaceID = "51081496"
	flagName = "更新后的需求"
	flagDescription = "新的描述"
	flagStatus = "进行中"
	flagOwner = "alice"
	flagPriority = "High"
	flagIterationID = "123"
	flagPretty = false

	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	os.Stdout = stdoutWriter

	runErr := runStoryUpdate(nil, []string{"10001"})
	_ = stdoutWriter.Close()
	os.Stdout = oldStdout

	if runErr != nil {
		t.Fatalf("runStoryUpdate returned error: %v", runErr)
	}
	if _, err := io.ReadAll(stdoutReader); err != nil {
		t.Fatalf("failed to read stdout: %v", err)
	}
	_ = stdoutReader.Close()

	if got := captured.Get("workspace_id"); got != "51081496" {
		t.Fatalf("workspace_id = %q, want %q", got, "51081496")
	}
	if got := captured.Get("id"); got != "10001" {
		t.Fatalf("id = %q, want %q", got, "10001")
	}
	if got := captured.Get("iteration_id"); got != "123" {
		t.Fatalf("iteration_id = %q, want %q", got, "123")
	}
	if got := captured.Get("v_status"); got != "进行中" {
		t.Fatalf("v_status = %q, want %q", got, "进行中")
	}
	if got := captured.Get("description"); got == "" {
		t.Fatal("description should be converted and sent when provided")
	}
}

func TestStoryUpdateFlagRegistered(t *testing.T) {
	flag := storyUpdateCmd.Flags().Lookup("iteration-id")
	if flag == nil {
		t.Fatal("story update should register --iteration-id")
	}
	if flag.Usage != "新迭代 ID" {
		t.Fatalf("iteration-id usage = %q, want %q", flag.Usage, "新迭代 ID")
	}
}
