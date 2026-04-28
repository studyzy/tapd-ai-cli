// Package cmd 中的 coverage_test.go 补充 mock server 单元测试以提升覆盖率
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

// resetFlags 将所有包级 flag 变量重置为零值，避免测试之间互相干扰
func resetFlags() {
	flagStatus = ""
	flagOwner = ""
	flagIterationID = ""
	flagLimit = 10
	flagPage = 1
	flagName = ""
	flagDescription = ""
	flagDescFile = ""
	flagPriority = ""
	flagParentID = ""
	flagLabel = ""
	flagCC = ""
	flagDeveloper = ""
	flagCategoryID = ""
	flagBegin = ""
	flagDue = ""
	flagModule = ""
	flagResolution = ""
	flagEffort = ""
	flagCustomField = nil
	flagStoryID = ""
	flagTitle = ""
	flagSeverity = ""
	flagReporter = ""
	flagCurrentOwner = ""
	flagCurrentUser = ""
	flagOrder = ""
	flagJSON = false
	flagPretty = false
	flagNoComments = true // 默认不输出评论避免额外请求
}

// setupMockServer 创建 mock HTTP server，根据请求路径返回对应的 JSON
func setupMockServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, func()) {
	t.Helper()
	srv := httptest.NewServer(handler)
	oldClient := apiClient
	oldWsID := flagWorkspaceID
	oldStdout := os.Stdout

	apiClient = tapd.NewClientWithBaseURL(srv.URL, srv.URL, "test-token", "", "")
	flagWorkspaceID = "12345"

	cleanup := func() {
		srv.Close()
		apiClient = oldClient
		flagWorkspaceID = oldWsID
		os.Stdout = oldStdout
	}
	return srv, cleanup
}

// captureStdout 替换 os.Stdout 为 pipe，返回恢复函数和读取器
func captureStdout(t *testing.T) (restore func(), reader *os.File) {
	t.Helper()
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	os.Stdout = w
	return func() {
		w.Close()
		os.Stdout = oldStdout
	}, r
}

// drainReader 读取并丢弃 reader 中的内容
func drainReader(r *os.File) {
	io.ReadAll(r)
	r.Close()
}

// ===================== parseCustomFields 测试 =====================

func TestParseCustomFields_Nil(t *testing.T) {
	result := parseCustomFields(nil)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestParseCustomFields_Empty(t *testing.T) {
	result := parseCustomFields([]string{})
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestParseCustomFields_Single(t *testing.T) {
	result := parseCustomFields([]string{"custom_field_one=hello"})
	if len(result) != 1 || result["custom_field_one"] != "hello" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestParseCustomFields_Multiple(t *testing.T) {
	result := parseCustomFields([]string{"custom_field_1=v1", "custom_field_2=v2"})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result["custom_field_1"] != "v1" || result["custom_field_2"] != "v2" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestParseCustomFields_ValueWithEquals(t *testing.T) {
	result := parseCustomFields([]string{"key=a=b=c"})
	if result["key"] != "a=b=c" {
		t.Errorf("expected 'a=b=c', got %q", result["key"])
	}
}

func TestParseCustomFields_NoEquals(t *testing.T) {
	result := parseCustomFields([]string{"noequals"})
	if len(result) != 0 {
		t.Errorf("expected empty map for invalid entry, got %v", result)
	}
}

func TestParseCustomFields_EmptyKey(t *testing.T) {
	result := parseCustomFields([]string{"=value"})
	if len(result) != 0 {
		t.Errorf("expected empty map for empty key, got %v", result)
	}
}

func TestParseCustomFields_EmptyValue(t *testing.T) {
	result := parseCustomFields([]string{"key="})
	if result["key"] != "" {
		t.Errorf("expected empty value, got %q", result["key"])
	}
}

// ===================== Story 命令测试 =====================

// storyAPIHandler 为 story 相关 API 提供 mock 响应
func storyAPIHandler(t *testing.T, captured *url.Values) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.Method == http.MethodPost {
			*captured = r.PostForm
		} else {
			*captured = r.URL.Query()
		}
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/stories/count"):
			w.Write([]byte(`{"status":1,"data":{"count":5}}`))
		case strings.HasSuffix(path, "/stories"):
			if r.Method == http.MethodPost {
				w.Write([]byte(`{"status":1,"data":{"Story":{"id":"10001","name":"Test","url":"http://test/story/10001"}}}`))
			} else {
				w.Write([]byte(`{"status":1,"data":[{"Story":{"id":"10001","name":"Test","status":"open","owner":"alice"}}]}`))
			}
		case strings.Contains(path, "/comments/count"):
			w.Write([]byte(`{"status":1,"data":{"count":0}}`))
		case strings.Contains(path, "/comments"):
			w.Write([]byte(`{"status":1,"data":[]}`))
		default:
			w.Write([]byte(`{"status":1,"data":{}}`))
		}
	}
}

func TestRunStoryCreate_PassesNewFlags(t *testing.T) {
	resetFlags()
	var captured url.Values
	_, cleanup := setupMockServer(t, storyAPIHandler(t, &captured))
	defer cleanup()

	flagName = "新需求"
	flagDescription = "描述内容"
	flagDeveloper = "bob"
	flagCC = "carol"
	flagCategoryID = "cat001"
	flagBegin = "2026-01-01"
	flagDue = "2026-01-31"
	flagLabel = "标签A"
	flagCustomField = []string{"custom_field_1=val1", "custom_field_2=val2"}

	restore, reader := captureStdout(t)
	err := runStoryCreate(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runStoryCreate failed: %v", err)
	}
	if captured.Get("developer") != "bob" {
		t.Errorf("developer = %q, want %q", captured.Get("developer"), "bob")
	}
	if captured.Get("cc") != "carol" {
		t.Errorf("cc = %q, want %q", captured.Get("cc"), "carol")
	}
	if captured.Get("category_id") != "cat001" {
		t.Errorf("category_id = %q, want %q", captured.Get("category_id"), "cat001")
	}
	if captured.Get("begin") != "2026-01-01" {
		t.Errorf("begin = %q, want %q", captured.Get("begin"), "2026-01-01")
	}
	if captured.Get("due") != "2026-01-31" {
		t.Errorf("due = %q, want %q", captured.Get("due"), "2026-01-31")
	}
	if captured.Get("label") != "标签A" {
		t.Errorf("label = %q, want %q", captured.Get("label"), "标签A")
	}
	if captured.Get("custom_field_1") != "val1" {
		t.Errorf("custom_field_1 = %q, want %q", captured.Get("custom_field_1"), "val1")
	}
	if captured.Get("custom_field_2") != "val2" {
		t.Errorf("custom_field_2 = %q, want %q", captured.Get("custom_field_2"), "val2")
	}
}

func TestRunStoryUpdate_PassesNewFlags(t *testing.T) {
	resetFlags()
	var captured url.Values
	_, cleanup := setupMockServer(t, storyAPIHandler(t, &captured))
	defer cleanup()

	flagDeveloper = "dev1"
	flagCC = "cc1"
	flagCurrentUser = "user1"
	flagCategoryID = "cat1"
	flagBegin = "2026-02-01"
	flagDue = "2026-02-28"
	flagLabel = "L1"
	flagCustomField = []string{"custom_field_one=cv1"}

	restore, reader := captureStdout(t)
	err := runStoryUpdate(nil, []string{"10001"})
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runStoryUpdate failed: %v", err)
	}
	if captured.Get("developer") != "dev1" {
		t.Errorf("developer = %q, want %q", captured.Get("developer"), "dev1")
	}
	if captured.Get("cc") != "cc1" {
		t.Errorf("cc = %q, want %q", captured.Get("cc"), "cc1")
	}
	if captured.Get("current_user") != "user1" {
		t.Errorf("current_user = %q, want %q", captured.Get("current_user"), "user1")
	}
	if captured.Get("category_id") != "cat1" {
		t.Errorf("category_id = %q, want %q", captured.Get("category_id"), "cat1")
	}
	if captured.Get("begin") != "2026-02-01" {
		t.Errorf("begin = %q, want %q", captured.Get("begin"), "2026-02-01")
	}
	if captured.Get("label") != "L1" {
		t.Errorf("label = %q, want %q", captured.Get("label"), "L1")
	}
	if captured.Get("custom_field_one") != "cv1" {
		t.Errorf("custom_field_one = %q, want %q", captured.Get("custom_field_one"), "cv1")
	}
}

func TestRunStoryShow_Mock(t *testing.T) {
	resetFlags()
	handler := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.Contains(path, "/comments"):
			w.Write([]byte(`{"status":1,"data":[]}`))
		default:
			w.Write([]byte(`{"status":1,"data":[{"Story":{"id":"10001","name":"Test Story","description":"<p>Hello</p>","status":"open","owner":"alice"}}]}`))
		}
	}
	_, cleanup := setupMockServer(t, handler)
	defer cleanup()
	flagJSON = true

	restore, reader := captureStdout(t)
	err := runStoryShow(nil, []string{"10001"})
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runStoryShow failed: %v", err)
	}
}

func TestRunStoryList_PassesNewFlags(t *testing.T) {
	resetFlags()
	var listQuery url.Values
	handler := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/stories/count"):
			w.Write([]byte(`{"status":1,"data":{"count":5}}`))
		case strings.HasSuffix(path, "/stories"):
			if r.Method == http.MethodGet {
				listQuery = r.URL.Query()
			}
			w.Write([]byte(`{"status":1,"data":[{"Story":{"id":"10001","name":"Test","status":"open","owner":"alice"}}]}`))
		default:
			w.Write([]byte(`{"status":1,"data":{}}`))
		}
	}
	_, cleanup := setupMockServer(t, handler)
	defer cleanup()

	flagName = "测试"
	flagPriority = "High"
	flagCategoryID = "cat1"
	flagLabel = "tag1"
	flagOrder = "created desc"

	restore, reader := captureStdout(t)
	err := runStoryList(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runStoryList failed: %v", err)
	}
	if listQuery.Get("name") != "测试" {
		t.Errorf("name = %q, want %q", listQuery.Get("name"), "测试")
	}
	if listQuery.Get("priority_label") != "High" {
		t.Errorf("priority_label = %q, want %q", listQuery.Get("priority_label"), "High")
	}
	if listQuery.Get("category_id") != "cat1" {
		t.Errorf("category_id = %q, want %q", listQuery.Get("category_id"), "cat1")
	}
	if listQuery.Get("order") != "created desc" {
		t.Errorf("order = %q, want %q", listQuery.Get("order"), "created desc")
	}
}

// ===================== Bug 命令测试 =====================

func bugAPIHandler(t *testing.T, captured *url.Values) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.Method == http.MethodPost {
			*captured = r.PostForm
		} else {
			*captured = r.URL.Query()
		}
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/bugs/count"):
			w.Write([]byte(`{"status":1,"data":{"count":3}}`))
		case strings.HasSuffix(path, "/bugs"):
			if r.Method == http.MethodPost {
				w.Write([]byte(`{"status":1,"data":{"Bug":{"id":"20001","title":"Test Bug","url":"http://test/bug/20001"}}}`))
			} else {
				w.Write([]byte(`{"status":1,"data":[{"Bug":{"id":"20001","title":"Test Bug","status":"new","current_owner":"alice","severity":"normal"}}]}`))
			}
		case strings.Contains(path, "/comments/count"):
			w.Write([]byte(`{"status":1,"data":{"count":0}}`))
		case strings.Contains(path, "/comments"):
			w.Write([]byte(`{"status":1,"data":[]}`))
		default:
			w.Write([]byte(`{"status":1,"data":{}}`))
		}
	}
}

func TestRunBugCreate_PassesNewFlags(t *testing.T) {
	resetFlags()
	var captured url.Values
	_, cleanup := setupMockServer(t, bugAPIHandler(t, &captured))
	defer cleanup()

	flagTitle = "新缺陷"
	flagDescription = "缺陷描述"
	flagPriority = "high"
	flagSeverity = "serious"
	flagCurrentOwner = "bob"
	flagCC = "carol"
	flagIterationID = "iter001"
	flagModule = "模块A"
	flagLabel = "bug-tag"
	flagBegin = "2026-03-01"
	flagDue = "2026-03-15"
	flagCustomField = []string{"custom_field_3=v3"}

	restore, reader := captureStdout(t)
	err := runBugCreate(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runBugCreate failed: %v", err)
	}
	if captured.Get("current_owner") != "bob" {
		t.Errorf("current_owner = %q, want %q", captured.Get("current_owner"), "bob")
	}
	if captured.Get("cc") != "carol" {
		t.Errorf("cc = %q, want %q", captured.Get("cc"), "carol")
	}
	if captured.Get("iteration_id") != "iter001" {
		t.Errorf("iteration_id = %q, want %q", captured.Get("iteration_id"), "iter001")
	}
	if captured.Get("module") != "模块A" {
		t.Errorf("module = %q, want %q", captured.Get("module"), "模块A")
	}
	if captured.Get("label") != "bug-tag" {
		t.Errorf("label = %q, want %q", captured.Get("label"), "bug-tag")
	}
	if captured.Get("begin") != "2026-03-01" {
		t.Errorf("begin = %q, want %q", captured.Get("begin"), "2026-03-01")
	}
	if captured.Get("custom_field_3") != "v3" {
		t.Errorf("custom_field_3 = %q, want %q", captured.Get("custom_field_3"), "v3")
	}
}

func TestRunBugUpdate_PassesNewFlags(t *testing.T) {
	resetFlags()
	var captured url.Values
	_, cleanup := setupMockServer(t, bugAPIHandler(t, &captured))
	defer cleanup()

	flagCC = "cc_user"
	flagIterationID = "iter002"
	flagModule = "模块B"
	flagLabel = "fix"
	flagBegin = "2026-04-01"
	flagDue = "2026-04-15"
	flagCurrentUser = "updater"
	flagResolution = "fixed"
	flagCustomField = []string{"custom_field_5=v5"}

	restore, reader := captureStdout(t)
	err := runBugUpdate(nil, []string{"20001"})
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runBugUpdate failed: %v", err)
	}
	if captured.Get("cc") != "cc_user" {
		t.Errorf("cc = %q, want %q", captured.Get("cc"), "cc_user")
	}
	if captured.Get("iteration_id") != "iter002" {
		t.Errorf("iteration_id = %q, want %q", captured.Get("iteration_id"), "iter002")
	}
	if captured.Get("current_user") != "updater" {
		t.Errorf("current_user = %q, want %q", captured.Get("current_user"), "updater")
	}
	if captured.Get("resolution") != "fixed" {
		t.Errorf("resolution = %q, want %q", captured.Get("resolution"), "fixed")
	}
	if captured.Get("custom_field_5") != "v5" {
		t.Errorf("custom_field_5 = %q, want %q", captured.Get("custom_field_5"), "v5")
	}
}

func TestRunBugShow_Mock(t *testing.T) {
	resetFlags()
	handler := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.Contains(path, "/comments"):
			w.Write([]byte(`{"status":1,"data":[]}`))
		default:
			w.Write([]byte(`{"status":1,"data":[{"Bug":{"id":"20001","title":"Test Bug","description":"<p>Desc</p>","status":"new"}}]}`))
		}
	}
	_, cleanup := setupMockServer(t, handler)
	defer cleanup()
	flagJSON = true

	restore, reader := captureStdout(t)
	err := runBugShow(nil, []string{"20001"})
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runBugShow failed: %v", err)
	}
}

func TestRunBugList_PassesNewFlags(t *testing.T) {
	resetFlags()
	var listQuery url.Values
	handler := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/bugs/count"):
			w.Write([]byte(`{"status":1,"data":{"count":3}}`))
		case strings.HasSuffix(path, "/bugs"):
			if r.Method == http.MethodGet {
				listQuery = r.URL.Query()
			}
			w.Write([]byte(`{"status":1,"data":[{"Bug":{"id":"20001","title":"Test Bug","status":"new","current_owner":"alice","severity":"normal"}}]}`))
		default:
			w.Write([]byte(`{"status":1,"data":{}}`))
		}
	}
	_, cleanup := setupMockServer(t, handler)
	defer cleanup()

	flagTitle = "测试缺陷"
	flagIterationID = "it1"
	flagModule = "mod1"
	flagLabel = "lbl1"
	flagReporter = "reporter1"
	flagOrder = "modified desc"

	restore, reader := captureStdout(t)
	err := runBugList(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runBugList failed: %v", err)
	}
	if listQuery.Get("title") != "测试缺陷" {
		t.Errorf("title = %q, want %q", listQuery.Get("title"), "测试缺陷")
	}
	if listQuery.Get("iteration_id") != "it1" {
		t.Errorf("iteration_id = %q, want %q", listQuery.Get("iteration_id"), "it1")
	}
	if listQuery.Get("module") != "mod1" {
		t.Errorf("module = %q, want %q", listQuery.Get("module"), "mod1")
	}
	if listQuery.Get("reporter") != "reporter1" {
		t.Errorf("reporter = %q, want %q", listQuery.Get("reporter"), "reporter1")
	}
	if listQuery.Get("order") != "modified desc" {
		t.Errorf("order = %q, want %q", listQuery.Get("order"), "modified desc")
	}
}

// ===================== Task 命令测试 =====================

func taskAPIHandler(t *testing.T, captured *url.Values) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.Method == http.MethodPost {
			*captured = r.PostForm
		} else {
			*captured = r.URL.Query()
		}
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/tasks/count"):
			w.Write([]byte(`{"status":1,"data":{"count":2}}`))
		case strings.HasSuffix(path, "/tasks"):
			if r.Method == http.MethodPost {
				w.Write([]byte(`{"status":1,"data":{"Task":{"id":"30001","name":"Test Task","url":"http://test/task/30001"}}}`))
			} else {
				w.Write([]byte(`{"status":1,"data":[{"Task":{"id":"30001","name":"Test Task","status":"open","owner":"alice"}}]}`))
			}
		case strings.Contains(path, "/comments/count"):
			w.Write([]byte(`{"status":1,"data":{"count":0}}`))
		case strings.Contains(path, "/comments"):
			w.Write([]byte(`{"status":1,"data":[]}`))
		default:
			w.Write([]byte(`{"status":1,"data":{}}`))
		}
	}
}

func TestRunTaskCreate_PassesNewFlags(t *testing.T) {
	resetFlags()
	var captured url.Values
	_, cleanup := setupMockServer(t, taskAPIHandler(t, &captured))
	defer cleanup()

	flagName = "新任务"
	flagDescription = "任务描述"
	flagCC = "cc_task"
	flagBegin = "2026-05-01"
	flagDue = "2026-05-15"
	flagIterationID = "iter_task"
	flagEffort = "8"
	flagLabel = "task-tag"
	flagCustomField = []string{"custom_field_9=v9"}

	restore, reader := captureStdout(t)
	err := runTaskCreate(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runTaskCreate failed: %v", err)
	}
	if captured.Get("cc") != "cc_task" {
		t.Errorf("cc = %q, want %q", captured.Get("cc"), "cc_task")
	}
	if captured.Get("begin") != "2026-05-01" {
		t.Errorf("begin = %q, want %q", captured.Get("begin"), "2026-05-01")
	}
	if captured.Get("due") != "2026-05-15" {
		t.Errorf("due = %q, want %q", captured.Get("due"), "2026-05-15")
	}
	if captured.Get("iteration_id") != "iter_task" {
		t.Errorf("iteration_id = %q, want %q", captured.Get("iteration_id"), "iter_task")
	}
	if captured.Get("effort") != "8" {
		t.Errorf("effort = %q, want %q", captured.Get("effort"), "8")
	}
	if captured.Get("label") != "task-tag" {
		t.Errorf("label = %q, want %q", captured.Get("label"), "task-tag")
	}
	if captured.Get("custom_field_9") != "v9" {
		t.Errorf("custom_field_9 = %q, want %q", captured.Get("custom_field_9"), "v9")
	}
}

func TestRunTaskUpdate_PassesNewFlags(t *testing.T) {
	resetFlags()
	var captured url.Values
	_, cleanup := setupMockServer(t, taskAPIHandler(t, &captured))
	defer cleanup()

	flagCC = "cc_upd"
	flagBegin = "2026-06-01"
	flagDue = "2026-06-30"
	flagStoryID = "story999"
	flagIterationID = "iter999"
	flagPriority = "Middle"
	flagEffort = "16"
	flagLabel = "upd-tag"
	flagCurrentUser = "operator"
	flagCustomField = []string{"custom_field_7=v7"}

	restore, reader := captureStdout(t)
	err := runTaskUpdate(nil, []string{"30001"})
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runTaskUpdate failed: %v", err)
	}
	if captured.Get("cc") != "cc_upd" {
		t.Errorf("cc = %q, want %q", captured.Get("cc"), "cc_upd")
	}
	if captured.Get("story_id") != "story999" {
		t.Errorf("story_id = %q, want %q", captured.Get("story_id"), "story999")
	}
	if captured.Get("iteration_id") != "iter999" {
		t.Errorf("iteration_id = %q, want %q", captured.Get("iteration_id"), "iter999")
	}
	if captured.Get("priority_label") != "Middle" {
		t.Errorf("priority_label = %q, want %q", captured.Get("priority_label"), "Middle")
	}
	if captured.Get("effort") != "16" {
		t.Errorf("effort = %q, want %q", captured.Get("effort"), "16")
	}
	if captured.Get("current_user") != "operator" {
		t.Errorf("current_user = %q, want %q", captured.Get("current_user"), "operator")
	}
	if captured.Get("custom_field_7") != "v7" {
		t.Errorf("custom_field_7 = %q, want %q", captured.Get("custom_field_7"), "v7")
	}
}

func TestRunTaskShow_Mock(t *testing.T) {
	resetFlags()
	handler := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.Contains(path, "/comments"):
			w.Write([]byte(`{"status":1,"data":[]}`))
		default:
			w.Write([]byte(`{"status":1,"data":[{"Task":{"id":"30001","name":"Test Task","description":"<p>Task Desc</p>","status":"open"}}]}`))
		}
	}
	_, cleanup := setupMockServer(t, handler)
	defer cleanup()
	flagJSON = true

	restore, reader := captureStdout(t)
	err := runTaskShow(nil, []string{"30001"})
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runTaskShow failed: %v", err)
	}
}

func TestRunTaskList_PassesNewFlags(t *testing.T) {
	resetFlags()
	var listQuery url.Values
	handler := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/tasks/count"):
			w.Write([]byte(`{"status":1,"data":{"count":2}}`))
		case strings.HasSuffix(path, "/tasks"):
			if r.Method == http.MethodGet {
				listQuery = r.URL.Query()
			}
			w.Write([]byte(`{"status":1,"data":[{"Task":{"id":"30001","name":"Test Task","status":"open","owner":"alice"}}]}`))
		default:
			w.Write([]byte(`{"status":1,"data":{}}`))
		}
	}
	_, cleanup := setupMockServer(t, handler)
	defer cleanup()

	flagName = "任务搜索"
	flagStoryID = "s123"
	flagIterationID = "i456"
	flagPriority = "High"
	flagLabel = "t-label"
	flagOrder = "id desc"

	restore, reader := captureStdout(t)
	err := runTaskList(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runTaskList failed: %v", err)
	}
	if listQuery.Get("name") != "任务搜索" {
		t.Errorf("name = %q, want %q", listQuery.Get("name"), "任务搜索")
	}
	if listQuery.Get("story_id") != "s123" {
		t.Errorf("story_id = %q, want %q", listQuery.Get("story_id"), "s123")
	}
	if listQuery.Get("iteration_id") != "i456" {
		t.Errorf("iteration_id = %q, want %q", listQuery.Get("iteration_id"), "i456")
	}
	if listQuery.Get("priority_label") != "High" {
		t.Errorf("priority_label = %q, want %q", listQuery.Get("priority_label"), "High")
	}
	if listQuery.Get("order") != "id desc" {
		t.Errorf("order = %q, want %q", listQuery.Get("order"), "id desc")
	}
}

// ===================== Iteration 命令测试 =====================

func iterationAPIHandler(t *testing.T, captured *url.Values) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.Method == http.MethodPost {
			*captured = r.PostForm
		} else {
			*captured = r.URL.Query()
		}
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/iterations/count"):
			w.Write([]byte(`{"status":1,"data":{"count":1}}`))
		case strings.HasSuffix(path, "/iterations"):
			if r.Method == http.MethodPost {
				w.Write([]byte(`{"status":1,"data":{"Iteration":{"id":"40001","name":"Sprint 1","workspace_id":"12345"}}}`))
			} else {
				w.Write([]byte(`{"status":1,"data":[{"Iteration":{"id":"40001","name":"Sprint 1","status":"open"}}]}`))
			}
		default:
			w.Write([]byte(`{"status":1,"data":{}}`))
		}
	}
}

func TestRunIterationCreate_PassesNewFlags(t *testing.T) {
	resetFlags()
	var captured url.Values
	_, cleanup := setupMockServer(t, iterationAPIHandler(t, &captured))
	defer cleanup()

	flagName = "Sprint X"
	flagStartDate = "2026-07-01"
	flagEndDate = "2026-07-14"
	flagCreator = "pm_user"
	flagLabel = "sprint-tag"
	flagParentID = "parent001"

	restore, reader := captureStdout(t)
	err := runIterationCreate(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runIterationCreate failed: %v", err)
	}
	if captured.Get("label") != "sprint-tag" {
		t.Errorf("label = %q, want %q", captured.Get("label"), "sprint-tag")
	}
	if captured.Get("parent_id") != "parent001" {
		t.Errorf("parent_id = %q, want %q", captured.Get("parent_id"), "parent001")
	}
}

func TestRunIterationUpdate_Mock(t *testing.T) {
	resetFlags()
	var captured url.Values
	_, cleanup := setupMockServer(t, iterationAPIHandler(t, &captured))
	defer cleanup()

	flagCurrentUser = "pm_user"
	flagName = "Sprint Y"
	flagStatus = "done"

	restore, reader := captureStdout(t)
	err := runIterationUpdate(nil, []string{"40001"})
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runIterationUpdate failed: %v", err)
	}
	if captured.Get("current_user") != "pm_user" {
		t.Errorf("current_user = %q, want %q", captured.Get("current_user"), "pm_user")
	}
	if captured.Get("name") != "Sprint Y" {
		t.Errorf("name = %q, want %q", captured.Get("name"), "Sprint Y")
	}
}

func TestRunIterationCount_Mock(t *testing.T) {
	resetFlags()
	var captured url.Values
	_, cleanup := setupMockServer(t, iterationAPIHandler(t, &captured))
	defer cleanup()

	restore, reader := captureStdout(t)
	err := runIterationCount(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runIterationCount failed: %v", err)
	}
}

func TestRunIterationList_PassesNewFlags(t *testing.T) {
	resetFlags()
	var listQuery url.Values
	handler := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/iterations/count"):
			w.Write([]byte(`{"status":1,"data":{"count":1}}`))
		case strings.HasSuffix(path, "/iterations"):
			if r.Method == http.MethodGet {
				listQuery = r.URL.Query()
			}
			w.Write([]byte(`{"status":1,"data":[{"Iteration":{"id":"40001","name":"Sprint 1","status":"open"}}]}`))
		default:
			w.Write([]byte(`{"status":1,"data":{}}`))
		}
	}
	_, cleanup := setupMockServer(t, handler)
	defer cleanup()

	flagName = "Sprint"
	flagCreator = "pm1"
	flagOrder = "startdate asc"

	restore, reader := captureStdout(t)
	err := runIterationList(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runIterationList failed: %v", err)
	}
	if listQuery.Get("name") != "Sprint" {
		t.Errorf("name = %q, want %q", listQuery.Get("name"), "Sprint")
	}
	if listQuery.Get("creator") != "pm1" {
		t.Errorf("creator = %q, want %q", listQuery.Get("creator"), "pm1")
	}
	if listQuery.Get("order") != "startdate asc" {
		t.Errorf("order = %q, want %q", listQuery.Get("order"), "startdate asc")
	}
}

// ===================== printComments 测试 =====================

func TestPrintComments_WithComments(t *testing.T) {
	resetFlags()
	flagNoComments = false
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":1,"data":[{"Comment":{"id":"c1","author":"alice","created":"2026-01-01","description":"<p>评论内容</p>"}}]}`))
	}
	_, cleanup := setupMockServer(t, handler)
	defer cleanup()

	restore, reader := captureStdout(t)
	printComments("12345", "stories", "10001")
	restore()

	data, _ := io.ReadAll(reader)
	reader.Close()
	output := string(data)
	if !strings.Contains(output, "评论") {
		t.Errorf("expected comment output, got: %s", output)
	}
}

func TestPrintComments_NoComments(t *testing.T) {
	resetFlags()
	flagNoComments = false
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":1,"data":[]}`))
	}
	_, cleanup := setupMockServer(t, handler)
	defer cleanup()

	restore, reader := captureStdout(t)
	printComments("12345", "stories", "10001")
	restore()

	data, _ := io.ReadAll(reader)
	reader.Close()
	if len(data) > 0 {
		t.Errorf("expected no output for empty comments, got: %s", string(data))
	}
}

func TestPrintComments_Skipped(t *testing.T) {
	resetFlags()
	flagNoComments = true
	// 不需要 mock server，--no-comments 直接跳过
	oldClient := apiClient
	defer func() { apiClient = oldClient }()
	apiClient = tapd.NewClient("test", "", "")

	restore, reader := captureStdout(t)
	printComments("12345", "stories", "10001")
	restore()

	data, _ := io.ReadAll(reader)
	reader.Close()
	if len(data) > 0 {
		t.Errorf("expected no output when --no-comments, got: %s", string(data))
	}
}

func TestPrintComments_JSONOutput(t *testing.T) {
	resetFlags()
	flagNoComments = false
	flagJSON = true
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":1,"data":[{"Comment":{"id":"c1","author":"alice","created":"2026-01-01","description":"<p>JSON评论</p>"}}]}`))
	}
	_, cleanup := setupMockServer(t, handler)
	defer cleanup()

	restore, reader := captureStdout(t)
	printComments("12345", "stories", "10001")
	restore()

	data, _ := io.ReadAll(reader)
	reader.Close()
	output := string(data)
	if !strings.Contains(output, "comments") {
		t.Errorf("expected JSON comment output, got: %s", output)
	}
}

// ===================== Flag 注册验证测试 =====================

func TestNewFlagsRegistered_StoryList(t *testing.T) {
	for _, name := range []string{"name", "priority", "category-id", "label", "order"} {
		if storyListCmd.Flags().Lookup(name) == nil {
			t.Errorf("story list should register --%s", name)
		}
	}
}

func TestNewFlagsRegistered_StoryCreate(t *testing.T) {
	for _, name := range []string{"developer", "cc", "category-id", "begin", "due", "label", "custom-field"} {
		if storyCreateCmd.Flags().Lookup(name) == nil {
			t.Errorf("story create should register --%s", name)
		}
	}
}

func TestNewFlagsRegistered_StoryUpdate(t *testing.T) {
	for _, name := range []string{"developer", "cc", "current-user", "category-id", "begin", "due", "label", "custom-field"} {
		if storyUpdateCmd.Flags().Lookup(name) == nil {
			t.Errorf("story update should register --%s", name)
		}
	}
}

func TestNewFlagsRegistered_TaskList(t *testing.T) {
	for _, name := range []string{"name", "story-id", "iteration-id", "priority", "label", "order"} {
		if taskListCmd.Flags().Lookup(name) == nil {
			t.Errorf("task list should register --%s", name)
		}
	}
}

func TestNewFlagsRegistered_TaskCreate(t *testing.T) {
	for _, name := range []string{"cc", "begin", "due", "iteration-id", "effort", "label", "custom-field"} {
		if taskCreateCmd.Flags().Lookup(name) == nil {
			t.Errorf("task create should register --%s", name)
		}
	}
}

func TestNewFlagsRegistered_TaskUpdate(t *testing.T) {
	for _, name := range []string{"cc", "begin", "due", "story-id", "iteration-id", "priority", "effort", "label", "current-user", "custom-field"} {
		if taskUpdateCmd.Flags().Lookup(name) == nil {
			t.Errorf("task update should register --%s", name)
		}
	}
}

func TestNewFlagsRegistered_BugList(t *testing.T) {
	for _, name := range []string{"title", "iteration-id", "module", "label", "reporter", "order"} {
		if bugListCmd.Flags().Lookup(name) == nil {
			t.Errorf("bug list should register --%s", name)
		}
	}
}

func TestNewFlagsRegistered_BugCreate(t *testing.T) {
	for _, name := range []string{"current-owner", "cc", "iteration-id", "module", "label", "begin", "due", "custom-field"} {
		if bugCreateCmd.Flags().Lookup(name) == nil {
			t.Errorf("bug create should register --%s", name)
		}
	}
}

func TestNewFlagsRegistered_BugUpdate(t *testing.T) {
	for _, name := range []string{"cc", "iteration-id", "module", "label", "begin", "due", "current-user", "resolution", "custom-field"} {
		if bugUpdateCmd.Flags().Lookup(name) == nil {
			t.Errorf("bug update should register --%s", name)
		}
	}
}

func TestNewFlagsRegistered_IterationList(t *testing.T) {
	for _, name := range []string{"name", "creator", "order"} {
		if iterationListCmd.Flags().Lookup(name) == nil {
			t.Errorf("iteration list should register --%s", name)
		}
	}
}

func TestNewFlagsRegistered_IterationCreate(t *testing.T) {
	for _, name := range []string{"label", "parent-id"} {
		if iterationCreateCmd.Flags().Lookup(name) == nil {
			t.Errorf("iteration create should register --%s", name)
		}
	}
}

// ===================== 错误路径测试（提升 create/update 覆盖率） =====================

// testExitCode 在子进程中运行会 os.Exit 的函数，返回是否退出了
// 由于 runXxxCreate 在参数缺失时调用 os.Exit，直接调用会终止测试进程，
// 所以这里只测试正常路径已够。但我们可以测试 readDescription 从文件读取的分支，
// 以及 iteration 的参数校验（它们也是返回而非 os.Exit 前的逻辑）。

func TestReadDescription_FromFile_NewFlags(t *testing.T) {
	resetFlags()
	tmpFile := t.TempDir() + "/desc.md"
	os.WriteFile(tmpFile, []byte("# Hello\nWorld"), 0644)
	flagDescription = ""
	flagDescFile = tmpFile

	content, err := readDescription()
	if err != nil {
		t.Fatalf("readDescription failed: %v", err)
	}
	if content == "" {
		t.Fatal("expected non-empty content from file")
	}
	// 应该被转换为 HTML
	if !strings.Contains(content, "Hello") {
		t.Errorf("content should contain Hello, got: %s", content)
	}
}

func TestReadDescription_FromDescription(t *testing.T) {
	resetFlags()
	flagDescription = "直接描述"
	flagDescFile = ""

	content, err := readDescription()
	if err != nil {
		t.Fatalf("readDescription failed: %v", err)
	}
	if content == "" {
		t.Fatal("expected non-empty content")
	}
}

func TestReadDescription_FileNotFound_NewFlags(t *testing.T) {
	resetFlags()
	flagDescription = ""
	flagDescFile = "/nonexistent/path/desc.md"

	_, err := readDescription()
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestReadDescription_Empty(t *testing.T) {
	resetFlags()
	flagDescription = ""
	flagDescFile = ""

	content, err := readDescription()
	if err != nil {
		t.Fatalf("readDescription failed: %v", err)
	}
	if content != "" {
		t.Errorf("expected empty content, got: %s", content)
	}
}

func TestPrintDetail_JSON(t *testing.T) {
	resetFlags()
	flagJSON = true
	flagPretty = false

	restore, reader := captureStdout(t)
	err := printDetail(map[string]string{"name": "test"}, "")
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("printDetail failed: %v", err)
	}
}

func TestPrintDetail_Markdown(t *testing.T) {
	resetFlags()
	flagJSON = false
	flagPretty = false

	type sample struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	data := sample{Name: "test", Description: "hello world"}

	restore, reader := captureStdout(t)
	err := printDetail(data, "description")
	restore()

	out, _ := io.ReadAll(reader)
	reader.Close()

	if err != nil {
		t.Fatalf("printDetail failed: %v", err)
	}
	if !strings.Contains(string(out), "hello world") {
		t.Errorf("expected markdown body, got: %s", string(out))
	}
}

func TestUseJSONOutput(t *testing.T) {
	flagJSON = false
	flagPretty = false
	if useJSONOutput() {
		t.Error("expected false when both flags are false")
	}
	flagJSON = true
	if !useJSONOutput() {
		t.Error("expected true when --json is set")
	}
	flagJSON = false
	flagPretty = true
	if !useJSONOutput() {
		t.Error("expected true when --pretty is set")
	}
}

// ===================== 其他 0% 函数的 mock 测试 =====================

// genericAPIHandler 返回通用的成功 JSON 响应
func genericAPIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.Contains(path, "/all_transitions"):
			w.Write([]byte(`{"status":1,"data":[{"Name":"open-done","StepPrevious":"open","StepNext":"done"}]}`))
		case strings.Contains(path, "/status_map"):
			w.Write([]byte(`{"status":1,"data":{"open":"打开","done":"完成"}}`))
		case strings.Contains(path, "/last_steps"):
			w.Write([]byte(`{"status":1,"data":{"done":"完成"}}`))
		case strings.Contains(path, "/custom_fields_settings"):
			w.Write([]byte(`{"status":1,"data":[{"CustomFieldConfig":{"custom_field":"custom_field_one","name":"测试字段","type":"text"}}]}`))
		case strings.Contains(path, "/stories/get_fields_label"):
			w.Write([]byte(`{"status":1,"data":{"Story":{"id":"ID","name":"标题"}}}`))
		case strings.Contains(path, "/stories/get_fields_info"):
			w.Write([]byte(`{"status":1,"data":{"Story":{"id":{"label":"ID","type":"text"}}}}`))
		case strings.Contains(path, "/workitem_types"):
			w.Write([]byte(`{"status":1,"data":[{"WorkitemType":{"id":"1","name":"默认类别"}}]}`))
		default:
			w.Write([]byte(`{"status":1,"data":{}}`))
		}
	}
}

func TestRunWorkflowStatusMap_Mock(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, genericAPIHandler())
	defer cleanup()
	flagSystem = "story"
	flagWorkitemTypeID = "1"

	restore, reader := captureStdout(t)
	err := runWorkflowStatusMap(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runWorkflowStatusMap failed: %v", err)
	}
}

func TestRunWorkflowTransitions_Mock(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, genericAPIHandler())
	defer cleanup()
	flagSystem = "story"
	flagWorkitemTypeID = "1"

	restore, reader := captureStdout(t)
	err := runWorkflowTransitions(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runWorkflowTransitions failed: %v", err)
	}
}

func TestRunWorkflowLastSteps_Mock(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, genericAPIHandler())
	defer cleanup()
	flagSystem = "bug"

	restore, reader := captureStdout(t)
	err := runWorkflowLastSteps(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runWorkflowLastSteps failed: %v", err)
	}
}

func TestRunCustomFieldList_Mock(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, genericAPIHandler())
	defer cleanup()
	flagEntityType = "stories"

	restore, reader := captureStdout(t)
	err := runCustomFieldList(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runCustomFieldList failed: %v", err)
	}
}

func TestRunStoryFieldLabel_Mock(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, genericAPIHandler())
	defer cleanup()

	restore, reader := captureStdout(t)
	err := runStoryFieldLabel(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runStoryFieldLabel failed: %v", err)
	}
}

func TestRunStoryFieldInfo_Mock(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, genericAPIHandler())
	defer cleanup()

	restore, reader := captureStdout(t)
	err := runStoryFieldInfo(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runStoryFieldInfo failed: %v", err)
	}
}

func TestRunWorkitemTypeList_Mock(t *testing.T) {
	resetFlags()
	_, cleanup := setupMockServer(t, genericAPIHandler())
	defer cleanup()

	restore, reader := captureStdout(t)
	err := runWorkitemTypeList(nil, nil)
	restore()
	drainReader(reader)

	if err != nil {
		t.Fatalf("runWorkitemTypeList failed: %v", err)
	}
}
