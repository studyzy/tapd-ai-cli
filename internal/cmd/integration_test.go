// Package cmd 中的 integration_test.go 使用真实 TAPD API 进行集成测试。
// 需要设置环境变量才会执行：
//   - TAPD_ACCESS_TOKEN 或 TAPD_API_USER + TAPD_API_PASSWORD
//   - TAPD_WORKSPACE_ID
package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/studyzy/tapd-ai-cli/internal/client"
	"github.com/studyzy/tapd-ai-cli/internal/model"
)

// skipIfNoCredentials 检查环境变量，若无凭据则跳过测试
func skipIfNoCredentials(t *testing.T) {
	t.Helper()
	token := os.Getenv("TAPD_ACCESS_TOKEN")
	user := os.Getenv("TAPD_API_USER")
	pass := os.Getenv("TAPD_API_PASSWORD")
	if token == "" && (user == "" || pass == "") {
		t.Skip("Skipping integration test: no TAPD credentials in environment (set TAPD_ACCESS_TOKEN or TAPD_API_USER/TAPD_API_PASSWORD)")
	}
}

// skipIfNoWorkspace 检查 TAPD_WORKSPACE_ID 环境变量
func skipIfNoWorkspace(t *testing.T) {
	t.Helper()
	skipIfNoCredentials(t)
	if os.Getenv("TAPD_WORKSPACE_ID") == "" {
		t.Skip("Skipping integration test: TAPD_WORKSPACE_ID not set")
	}
}

// setupIntegrationClient 初始化真实 API 客户端
func setupIntegrationClient(t *testing.T) *client.Client {
	t.Helper()
	token := os.Getenv("TAPD_ACCESS_TOKEN")
	user := os.Getenv("TAPD_API_USER")
	pass := os.Getenv("TAPD_API_PASSWORD")
	return client.NewClient(token, user, pass)
}

// setupIntegrationCmd 初始化 cmd 包的全局变量用于集成测试
func setupIntegrationCmd(t *testing.T) {
	t.Helper()
	apiClient = setupIntegrationClient(t)
	flagWorkspaceID = os.Getenv("TAPD_WORKSPACE_ID")
	flagPretty = false
}

func TestIntegration_AuthTestAuth(t *testing.T) {
	skipIfNoCredentials(t)
	c := setupIntegrationClient(t)
	if err := c.TestAuth(); err != nil {
		t.Fatalf("TestAuth failed: %v", err)
	}
}

func TestIntegration_WorkspaceList(t *testing.T) {
	skipIfNoCredentials(t)
	c := setupIntegrationClient(t)

	workspaces, err := c.ListWorkspaces()
	if err != nil {
		t.Fatalf("ListWorkspaces failed: %v", err)
	}
	if len(workspaces) == 0 {
		t.Fatal("Expected at least one workspace")
	}
	// 验证没有 organization 类型
	for _, ws := range workspaces {
		if ws.Category == "organization" {
			t.Errorf("ListWorkspaces should filter organization entries, got: %+v", ws)
		}
	}
	t.Logf("Found %d workspaces", len(workspaces))
}

func TestIntegration_WorkspaceInfo(t *testing.T) {
	skipIfNoWorkspace(t)
	c := setupIntegrationClient(t)

	ws, err := c.GetWorkspaceInfo(os.Getenv("TAPD_WORKSPACE_ID"))
	if err != nil {
		t.Fatalf("GetWorkspaceInfo failed: %v", err)
	}
	if ws.ID == "" || ws.Name == "" {
		t.Errorf("Workspace missing fields: %+v", ws)
	}
	t.Logf("Workspace: id=%s name=%s", ws.ID, ws.Name)
}

func TestIntegration_StoryList(t *testing.T) {
	skipIfNoWorkspace(t)
	c := setupIntegrationClient(t)
	wsID := os.Getenv("TAPD_WORKSPACE_ID")

	params := map[string]string{
		"workspace_id": wsID,
		"entity_type":  "stories",
		"limit":        "3",
		"fields":       "id,name,status,owner,modified",
	}
	result, err := c.ListStories(params)
	if err != nil {
		t.Fatalf("ListStories failed: %v", err)
	}
	stories, ok := result.([]model.Story)
	if !ok {
		t.Fatalf("expected []model.Story, got %T", result)
	}
	t.Logf("Found %d stories", len(stories))
	for _, s := range stories {
		t.Logf("  Story: id=%v name=%v", s.ID, s.Name)
	}
}

func TestIntegration_StoryCount(t *testing.T) {
	skipIfNoWorkspace(t)
	c := setupIntegrationClient(t)
	wsID := os.Getenv("TAPD_WORKSPACE_ID")

	count, err := c.CountStories(map[string]string{
		"workspace_id": wsID,
		"entity_type":  "stories",
	})
	if err != nil {
		t.Fatalf("CountStories failed: %v", err)
	}
	t.Logf("Story count: %d", count)
}

func TestIntegration_BugList(t *testing.T) {
	skipIfNoWorkspace(t)
	c := setupIntegrationClient(t)
	wsID := os.Getenv("TAPD_WORKSPACE_ID")

	params := map[string]string{
		"workspace_id": wsID,
		"limit":        "3",
	}
	bugs, err := c.ListBugs(params)
	if err != nil {
		t.Fatalf("ListBugs failed: %v", err)
	}
	t.Logf("Found %d bugs", len(bugs))
}

func TestIntegration_BugCount(t *testing.T) {
	skipIfNoWorkspace(t)
	c := setupIntegrationClient(t)
	wsID := os.Getenv("TAPD_WORKSPACE_ID")

	count, err := c.CountBugs(map[string]string{
		"workspace_id": wsID,
	})
	if err != nil {
		t.Fatalf("CountBugs failed: %v", err)
	}
	t.Logf("Bug count: %d", count)
}

func TestIntegration_IterationList(t *testing.T) {
	skipIfNoWorkspace(t)
	c := setupIntegrationClient(t)
	wsID := os.Getenv("TAPD_WORKSPACE_ID")

	params := map[string]string{
		"workspace_id": wsID,
	}
	iterations, err := c.ListIterations(params)
	if err != nil {
		t.Fatalf("ListIterations failed: %v", err)
	}
	t.Logf("Found %d iterations", len(iterations))
}

func TestIntegration_RunWorkspaceList(t *testing.T) {
	skipIfNoCredentials(t)
	setupIntegrationCmd(t)

	var buf bytes.Buffer
	err := runWorkspaceList(nil, nil)
	if err != nil {
		t.Fatalf("runWorkspaceList failed: %v", err)
	}
	_ = buf
}

func TestIntegration_RunWorkspaceInfo(t *testing.T) {
	skipIfNoWorkspace(t)
	setupIntegrationCmd(t)

	err := runWorkspaceInfo(nil, nil)
	if err != nil {
		t.Fatalf("runWorkspaceInfo failed: %v", err)
	}
}

func TestIntegration_RunStoryList(t *testing.T) {
	skipIfNoWorkspace(t)
	setupIntegrationCmd(t)
	flagStatus = ""
	flagOwner = ""
	flagLimit = 3
	flagPage = 1

	err := runStoryList(nil, nil)
	if err != nil {
		t.Fatalf("runStoryList failed: %v", err)
	}
}

func TestIntegration_RunStoryCount(t *testing.T) {
	skipIfNoWorkspace(t)
	setupIntegrationCmd(t)
	flagStatus = ""

	err := runStoryCount(nil, nil)
	if err != nil {
		t.Fatalf("runStoryCount failed: %v", err)
	}
}

func TestIntegration_RunBugList(t *testing.T) {
	skipIfNoWorkspace(t)
	setupIntegrationCmd(t)
	flagStatus = ""
	flagPriority = ""
	flagSeverity = ""
	flagLimit = 3
	flagPage = 1

	err := runBugList(nil, nil)
	if err != nil {
		t.Fatalf("runBugList failed: %v", err)
	}
}

func TestIntegration_RunBugCount(t *testing.T) {
	skipIfNoWorkspace(t)
	setupIntegrationCmd(t)
	flagStatus = ""

	err := runBugCount(nil, nil)
	if err != nil {
		t.Fatalf("runBugCount failed: %v", err)
	}
}

func TestIntegration_RunIterationList(t *testing.T) {
	skipIfNoWorkspace(t)
	setupIntegrationCmd(t)
	flagStatus = ""

	err := runIterationList(nil, nil)
	if err != nil {
		t.Fatalf("runIterationList failed: %v", err)
	}
}

func TestIntegration_RunTaskList(t *testing.T) {
	skipIfNoWorkspace(t)
	setupIntegrationCmd(t)
	flagStatus = ""
	flagOwner = ""
	flagLimit = 3
	flagPage = 1

	err := runTaskList(nil, nil)
	if err != nil {
		t.Fatalf("runTaskList failed: %v", err)
	}
}

func TestIntegration_RunTaskCount(t *testing.T) {
	skipIfNoWorkspace(t)
	setupIntegrationCmd(t)
	flagStatus = ""

	err := runTaskCount(nil, nil)
	if err != nil {
		t.Fatalf("runTaskCount failed: %v", err)
	}
}

func TestIntegration_RunSpec(t *testing.T) {
	// spec 不需要凭据
	flagPretty = false
	err := runSpec(nil, nil)
	if err != nil {
		t.Fatalf("runSpec failed: %v", err)
	}
}

func TestIntegration_SpecOutputValid(t *testing.T) {
	flagPretty = false
	// 捕获 stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runSpec(nil, nil)
	if err != nil {
		t.Fatalf("runSpec failed: %v", err)
	}

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)

	var tools []json.RawMessage
	if err := json.Unmarshal(buf.Bytes(), &tools); err != nil {
		t.Fatalf("spec output is not valid JSON array: %v\nOutput: %s", err, buf.String())
	}
	if len(tools) == 0 {
		t.Fatal("spec output has no tools")
	}
	t.Logf("spec output contains %d tool definitions", len(tools))
}

// TestIntegration_E2E_CreateAndShowStory 创建一个需求然后查看详情（端到端）
// 注意：TAPD API 没有删除接口，创建后通过 t.Cleanup 更新标题标记为废弃
func TestIntegration_E2E_CreateAndShowStory(t *testing.T) {
	skipIfNoWorkspace(t)
	c := setupIntegrationClient(t)
	wsID := os.Getenv("TAPD_WORKSPACE_ID")

	// 创建
	result, err := c.CreateStory(map[string]string{
		"workspace_id": wsID,
		"name":         "[tapd-ai-cli integration test] 自动化测试需求",
	}, "stories")
	if err != nil {
		t.Fatalf("CreateStory failed: %v", err)
	}
	if !result.Success || result.ID == "" {
		t.Fatalf("Expected success with ID, got: %+v", result)
	}
	t.Logf("Created story: id=%s url=%s", result.ID, result.URL)

	// 清理：标记为废弃
	t.Cleanup(func() {
		_, err := c.UpdateStory(map[string]string{
			"workspace_id": wsID,
			"id":           result.ID,
			"name":         "[已废弃-自动化测试] 请忽略此需求",
		}, "stories")
		if err != nil {
			t.Logf("Cleanup: failed to mark story %s as deprecated: %v", result.ID, err)
		} else {
			t.Logf("Cleanup: marked story %s as deprecated", result.ID)
		}
	})

	// 查看详情
	detail, err := c.GetStory(wsID, result.ID, "stories")
	if err != nil {
		t.Fatalf("GetStory failed: %v", err)
	}
	story, ok := detail.(*model.Story)
	if !ok {
		t.Fatalf("expected *model.Story, got %T", detail)
	}
	if story.Name == "" {
		t.Errorf("Story name is empty: %+v", story)
	}
	t.Logf("Story detail: name=%s status=%v", story.Name, story.Status)
}

// TestIntegration_WorkspaceSwitch 测试 workspace switch 写入当前目录
func TestIntegration_WorkspaceSwitch(t *testing.T) {
	skipIfNoWorkspace(t)

	// 切换到临时目录
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	setupIntegrationCmd(t)
	wsID := os.Getenv("TAPD_WORKSPACE_ID")

	err := runWorkspaceSwitch(nil, []string{wsID})
	if err != nil {
		t.Fatalf("runWorkspaceSwitch failed: %v", err)
	}

	// 验证 .tapd.json 被创建
	data, err := os.ReadFile(".tapd.json")
	if err != nil {
		t.Fatalf("Failed to read .tapd.json: %v", err)
	}

	var cfg model.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("Failed to parse .tapd.json: %v", err)
	}
	if cfg.WorkspaceID != wsID {
		t.Errorf("Expected workspace_id=%s, got=%s", wsID, cfg.WorkspaceID)
	}
	t.Logf("workspace switch wrote .tapd.json with workspace_id=%s", cfg.WorkspaceID)
}

func TestIntegration_RunWikiList(t *testing.T) {
	skipIfNoWorkspace(t)
	setupIntegrationCmd(t)
	flagLimit = 3
	flagPage = 1
	flagWikiName = ""

	err := runWikiList(nil, nil)
	if err != nil {
		t.Fatalf("runWikiList failed: %v", err)
	}
}

func TestIntegration_WikiList_Client(t *testing.T) {
	skipIfNoWorkspace(t)
	c := setupIntegrationClient(t)
	wsID := os.Getenv("TAPD_WORKSPACE_ID")

	wikis, err := c.ListWikis(map[string]string{
		"workspace_id": wsID,
		"limit":        "3",
		"fields":       "id,name,creator,modified",
	})
	if err != nil {
		t.Fatalf("ListWikis failed: %v", err)
	}
	t.Logf("Found %d wikis", len(wikis))
	for _, w := range wikis {
		t.Logf("  Wiki: id=%v name=%v", w.ID, w.Name)
	}
}

func TestIntegration_RunWikiShow(t *testing.T) {
	skipIfNoWorkspace(t)
	c := setupIntegrationClient(t)
	wsID := os.Getenv("TAPD_WORKSPACE_ID")

	// 先获取一个真实 wiki id
	wikis, err := c.ListWikis(map[string]string{
		"workspace_id": wsID,
		"limit":        "1",
		"fields":       "id,name",
	})
	if err != nil {
		t.Fatalf("ListWikis failed: %v", err)
	}
	if len(wikis) == 0 {
		t.Skip("No wikis in workspace, skipping show test")
	}

	wikiID := wikis[0].ID
	t.Logf("Testing wiki show with id=%s", wikiID)

	setupIntegrationCmd(t)
	err = runWikiShow(nil, []string{wikiID})
	if err != nil {
		t.Fatalf("runWikiShow failed: %v", err)
	}
}

func TestIntegration_URLCommand_StoryURL(t *testing.T) {
	skipIfNoWorkspace(t)
	c := setupIntegrationClient(t)
	wsID := os.Getenv("TAPD_WORKSPACE_ID")

	// 获取一个真实 story id
	listResult, err := c.ListStories(map[string]string{
		"workspace_id": wsID,
		"entity_type":  "stories",
		"limit":        "1",
		"fields":       "id,name",
	})
	if err != nil {
		t.Skip("No stories available for URL test")
	}
	storyList, ok := listResult.([]model.Story)
	if !ok || len(storyList) == 0 {
		t.Skip("No stories available for URL test")
	}
	storyID := storyList[0].ID
	storyURL := "https://www.tapd.cn/tapd_fe/" + wsID + "/story/detail/" + storyID

	// 验证 URL 解析
	parsed, err := parseTAPDURL(storyURL)
	if err != nil {
		t.Fatalf("parseTAPDURL(%q) failed: %v", storyURL, err)
	}
	if parsed.EntityType != "story" {
		t.Errorf("EntityType = %q, want %q", parsed.EntityType, "story")
	}
	if parsed.EntityID != storyID {
		t.Errorf("EntityID = %q, want %q", parsed.EntityID, storyID)
	}
	if parsed.WorkspaceID != wsID {
		t.Errorf("WorkspaceID = %q, want %q", parsed.WorkspaceID, wsID)
	}

	// 验证实际 API 调用
	urlResult, err := c.GetStory(wsID, storyID, "stories")
	if err != nil {
		t.Fatalf("GetStory via URL failed: %v", err)
	}
	urlStory, ok := urlResult.(*model.Story)
	if !ok {
		t.Fatalf("expected *model.Story, got %T", urlResult)
	}
	t.Logf("URL→Story: id=%v name=%v", urlStory.ID, urlStory.Name)
}

// TestIntegration_E2E_StoryCommentFlow 创建需求后，对其进行添加评论、查询评论、更新评论、查询评论数量的端到端测试
func TestIntegration_E2E_StoryCommentFlow(t *testing.T) {
	skipIfNoWorkspace(t)
	c := setupIntegrationClient(t)
	wsID := os.Getenv("TAPD_WORKSPACE_ID")

	// 步骤 1：创建需求
	storyResult, err := c.CreateStory(map[string]string{
		"workspace_id": wsID,
		"name":         "[tapd-ai-cli integration test] 评论功能测试需求",
	}, "stories")
	if err != nil {
		t.Fatalf("CreateStory failed: %v", err)
	}
	if !storyResult.Success || storyResult.ID == "" {
		t.Fatalf("Expected success with ID, got: %+v", storyResult)
	}
	storyID := storyResult.ID
	t.Logf("Step 1: Created story id=%s", storyID)

	// 清理：标记需求为废弃
	t.Cleanup(func() {
		_, err := c.UpdateStory(map[string]string{
			"workspace_id": wsID,
			"id":           storyID,
			"name":         "[已废弃-自动化测试] 评论功能测试需求-请忽略",
		}, "stories")
		if err != nil {
			t.Logf("Cleanup: failed to mark story %s as deprecated: %v", storyID, err)
		} else {
			t.Logf("Cleanup: marked story %s as deprecated", storyID)
		}
	})

	// 步骤 2：添加评论（API 客户端层）
	c.FetchNick()
	author := c.Nick
	if author == "" {
		author = os.Getenv("TAPD_API_USER")
	}
	t.Logf("Step 2: using author=%q for comment", author)
	comment, err := c.AddComment(map[string]string{
		"workspace_id": wsID,
		"entry_type":   "stories",
		"entry_id":     storyID,
		"description":  "这是一条自动化测试评论",
		"author":       author,
	})
	if err != nil {
		t.Fatalf("AddComment failed: %v", err)
	}
	if comment.ID == "" {
		t.Fatalf("Expected comment ID, got empty")
	}
	commentID := comment.ID
	t.Logf("Step 2: Added comment id=%s author=%s", commentID, comment.Author)

	// 步骤 3：查询评论列表
	comments, err := c.ListComments(map[string]string{
		"workspace_id": wsID,
		"entry_type":   "stories",
		"entry_id":     storyID,
	})
	if err != nil {
		t.Fatalf("ListComments failed: %v", err)
	}
	if len(comments) == 0 {
		t.Fatal("Expected at least 1 comment, got 0")
	}
	found := false
	for _, cm := range comments {
		if cm.ID == commentID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Comment %s not found in list results", commentID)
	}
	t.Logf("Step 3: Listed %d comments, found target comment=%s", len(comments), commentID)

	// 步骤 4：更新评论
	updated, err := c.UpdateComment(map[string]string{
		"workspace_id":   wsID,
		"id":             commentID,
		"description":    "这是更新后的自动化测试评论",
		"change_creator": c.Nick,
	})
	if err != nil {
		t.Fatalf("UpdateComment failed: %v", err)
	}
	if updated.ID != commentID {
		t.Errorf("Updated comment id=%q, want %q", updated.ID, commentID)
	}
	if updated.Description == "" {
		t.Error("Updated comment description should not be empty")
	}
	t.Logf("Step 4: Updated comment id=%s description=%s", updated.ID, updated.Description)

	// 步骤 5：查询评论数量
	count, err := c.CountComments(map[string]string{
		"workspace_id": wsID,
		"entry_type":   "stories",
		"entry_id":     storyID,
	})
	if err != nil {
		t.Fatalf("CountComments failed: %v", err)
	}
	if count < 1 {
		t.Errorf("Expected count >= 1, got %d", count)
	}
	t.Logf("Step 5: Comment count=%d", count)
}

// TestIntegration_E2E_StoryCommentFlow_Cmd 使用命令层函数完成评论的端到端测试
func TestIntegration_E2E_StoryCommentFlow_Cmd(t *testing.T) {
	skipIfNoWorkspace(t)
	c := setupIntegrationClient(t)
	wsID := os.Getenv("TAPD_WORKSPACE_ID")

	// 创建需求
	storyResult, err := c.CreateStory(map[string]string{
		"workspace_id": wsID,
		"name":         "[tapd-ai-cli integration test] 评论命令层测试需求",
	}, "stories")
	if err != nil {
		t.Fatalf("CreateStory failed: %v", err)
	}
	storyID := storyResult.ID
	t.Logf("Created story id=%s for cmd-level comment test", storyID)

	t.Cleanup(func() {
		_, err := c.UpdateStory(map[string]string{
			"workspace_id": wsID,
			"id":           storyID,
			"name":         "[已废弃-自动化测试] 评论命令层测试需求-请忽略",
		}, "stories")
		if err != nil {
			t.Logf("Cleanup: failed to mark story %s as deprecated: %v", storyID, err)
		}
	})

	setupIntegrationCmd(t)
	if apiClient.Nick == "" {
		apiClient.FetchNick()
	}

	// 测试 runCommentAdd
	flagEntryType = "stories"
	flagEntryID = storyID
	flagDescription = "命令层集成测试评论"
	flagCommentAuthor = ""
	flagReplyID = ""
	err = runCommentAdd(nil, nil)
	if err != nil {
		t.Fatalf("runCommentAdd failed: %v", err)
	}
	t.Log("runCommentAdd succeeded")

	// 测试 runCommentList
	flagEntryType = "stories"
	flagEntryID = storyID
	flagCommentAuthor = ""
	flagOrder = ""
	flagLimit = 10
	flagPage = 1
	err = runCommentList(nil, nil)
	if err != nil {
		t.Fatalf("runCommentList failed: %v", err)
	}
	t.Log("runCommentList succeeded")

	// 测试 runCommentCount
	flagEntryType = "stories"
	flagEntryID = storyID
	err = runCommentCount(nil, nil)
	if err != nil {
		t.Fatalf("runCommentCount failed: %v", err)
	}
	t.Log("runCommentCount succeeded")
}

func TestIntegration_URLCommand_WikiURL(t *testing.T) {
	skipIfNoWorkspace(t)
	c := setupIntegrationClient(t)
	wsID := os.Getenv("TAPD_WORKSPACE_ID")

	// 获取一个真实 wiki id
	wikis, err := c.ListWikis(map[string]string{
		"workspace_id": wsID,
		"limit":        "1",
		"fields":       "id,name",
	})
	if err != nil || len(wikis) == 0 {
		t.Skip("No wikis available for URL test")
	}
	wikiID := wikis[0].ID
	wikiURL := "https://www.tapd.cn/" + wsID + "/markdown_wikis/show/#" + wikiID

	// 验证 URL 解析
	parsed, err := parseTAPDURL(wikiURL)
	if err != nil {
		t.Fatalf("parseTAPDURL(%q) failed: %v", wikiURL, err)
	}
	if parsed.EntityType != "wiki" {
		t.Errorf("EntityType = %q, want %q", parsed.EntityType, "wiki")
	}
	if parsed.EntityID != wikiID {
		t.Errorf("EntityID = %q, want %q", parsed.EntityID, wikiID)
	}

	// 验证实际 API 调用
	result, err := c.GetWiki(wsID, wikiID)
	if err != nil {
		t.Fatalf("GetWiki via URL failed: %v", err)
	}
	t.Logf("URL→Wiki: id=%v name=%v", result.ID, result.Name)
}

// === 以下为新增命令的集成测试 ===

func TestIntegration_RunTodoList(t *testing.T) {
	skipIfNoWorkspace(t)
	setupIntegrationCmd(t)
	flagTodoEntityType = "story"
	flagLimit = 3
	flagPage = 1

	err := runTodoList(nil, nil)
	if err != nil {
		t.Fatalf("runTodoList failed: %v", err)
	}
}

func TestIntegration_RunTodoList_Bug(t *testing.T) {
	skipIfNoWorkspace(t)
	setupIntegrationCmd(t)
	flagTodoEntityType = "bug"
	flagLimit = 3
	flagPage = 1

	err := runTodoList(nil, nil)
	if err != nil {
		t.Fatalf("runTodoList(bug) failed: %v", err)
	}
}

func TestIntegration_RunTodoList_Task(t *testing.T) {
	skipIfNoWorkspace(t)
	setupIntegrationCmd(t)
	flagTodoEntityType = "task"
	flagLimit = 3
	flagPage = 1

	err := runTodoList(nil, nil)
	if err != nil {
		t.Fatalf("runTodoList(task) failed: %v", err)
	}
}

func TestIntegration_RunTimesheetList(t *testing.T) {
	skipIfNoWorkspace(t)
	setupIntegrationCmd(t)
	flagTimesheetEntityType = ""
	flagTimesheetEntityID = ""
	flagTimesheetOwner = ""
	flagLimit = 3
	flagPage = 1

	err := runTimesheetList(nil, nil)
	if err != nil {
		t.Fatalf("runTimesheetList failed: %v", err)
	}
}

func TestIntegration_RunReleaseList(t *testing.T) {
	skipIfNoWorkspace(t)
	setupIntegrationCmd(t)
	flagName = ""
	flagReleaseStatus = ""
	flagLimit = 3
	flagPage = 1

	err := runReleaseList(nil, nil)
	if err != nil {
		t.Fatalf("runReleaseList failed: %v", err)
	}
}

func TestIntegration_ReleaseList_Client(t *testing.T) {
	skipIfNoWorkspace(t)
	c := setupIntegrationClient(t)
	wsID := os.Getenv("TAPD_WORKSPACE_ID")

	releases, err := c.ListReleases(map[string]string{
		"workspace_id": wsID,
		"limit":        "3",
	})
	if err != nil {
		t.Fatalf("ListReleases failed: %v", err)
	}
	t.Logf("Found %d releases", len(releases))
	for _, r := range releases {
		t.Logf("  Release: id=%s name=%s status=%s", r.ID, r.Name, r.Status)
	}
}

func TestIntegration_TimesheetList_Client(t *testing.T) {
	skipIfNoWorkspace(t)
	c := setupIntegrationClient(t)
	wsID := os.Getenv("TAPD_WORKSPACE_ID")

	timesheets, err := c.ListTimesheets(map[string]string{
		"workspace_id": wsID,
		"limit":        "3",
	})
	if err != nil {
		t.Fatalf("ListTimesheets failed: %v", err)
	}
	t.Logf("Found %d timesheets", len(timesheets))
	for _, ts := range timesheets {
		t.Logf("  Timesheet: id=%s entity_type=%s timespent=%s owner=%s", ts.ID, ts.EntityType, ts.Timespent, ts.Owner)
	}
}

func TestIntegration_TodoList_Client(t *testing.T) {
	skipIfNoWorkspace(t)
	c := setupIntegrationClient(t)
	wsID := os.Getenv("TAPD_WORKSPACE_ID")

	data, err := c.GetTodo(map[string]string{
		"workspace_id": wsID,
		"entity_type":  "story",
		"limit":        "3",
	})
	if err != nil {
		t.Fatalf("GetTodo failed: %v", err)
	}
	t.Logf("Todo data length: %d bytes", len(data))
}

func TestIntegration_CommitMsg_Client(t *testing.T) {
	skipIfNoWorkspace(t)
	c := setupIntegrationClient(t)
	wsID := os.Getenv("TAPD_WORKSPACE_ID")

	// 先获取一个真实 story ID
	listResult, err := c.ListStories(map[string]string{
		"workspace_id": wsID,
		"entity_type":  "stories",
		"limit":        "1",
		"fields":       "id",
	})
	if err != nil {
		t.Skip("No stories available for commit-msg test")
	}
	stories, ok := listResult.([]model.Story)
	if !ok || len(stories) == 0 {
		t.Skip("No stories available for commit-msg test")
	}
	storyID := stories[0].ID

	data, err := c.GetCommitMsg(map[string]string{
		"workspace_id": wsID,
		"object_id":    storyID,
		"type":         "story",
	})
	if err != nil {
		t.Fatalf("GetCommitMsg failed: %v", err)
	}
	t.Logf("CommitMsg data: %s", string(data))
}

func TestIntegration_RunCommitMsgGet(t *testing.T) {
	skipIfNoWorkspace(t)
	c := setupIntegrationClient(t)
	wsID := os.Getenv("TAPD_WORKSPACE_ID")

	// 先获取一个真实 story ID
	listResult, err := c.ListStories(map[string]string{
		"workspace_id": wsID,
		"entity_type":  "stories",
		"limit":        "1",
		"fields":       "id",
	})
	if err != nil {
		t.Skip("No stories available for commit-msg test")
	}
	stories, ok := listResult.([]model.Story)
	if !ok || len(stories) == 0 {
		t.Skip("No stories available for commit-msg test")
	}

	setupIntegrationCmd(t)
	flagCommitMsgObjectID = stories[0].ID
	flagCommitMsgType = "story"

	err = runCommitMsgGet(nil, nil)
	if err != nil {
		t.Fatalf("runCommitMsgGet failed: %v", err)
	}
}
