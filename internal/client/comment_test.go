package client_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/studyzy/tapd-ai-cli/internal/client"
	"github.com/studyzy/tapd-ai-cli/internal/model"
)

func TestListComments(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/comments" {
			t.Errorf("unexpected path: %s, want /comments", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Comment":{"id":"100","title":"评论标题","description":"评论内容","author":"tester","entry_type":"stories","entry_id":"200","workspace_id":"1","created":"2026-04-16 10:00:00"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	results, err := c.ListComments(&model.ListCommentsRequest{
		WorkspaceID: "1",
	})
	if err != nil {
		t.Fatalf("ListComments() unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(results))
	}
	if results[0].ID != "100" {
		t.Errorf("comment id = %v, want %q", results[0].ID, "100")
	}
	if results[0].Author != "tester" {
		t.Errorf("comment author = %v, want %q", results[0].Author, "tester")
	}
	if results[0].EntryType != "stories" {
		t.Errorf("comment entry_type = %v, want %q", results[0].EntryType, "stories")
	}
	if results[0].EntryID != "200" {
		t.Errorf("comment entry_id = %v, want %q", results[0].EntryID, "200")
	}
}

func TestListComments_HTMLToMarkdown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Comment":{"id":"101","description":"<p>This is <strong>bold</strong> and <em>italic</em></p>","author":"tester","entry_type":"bug","entry_id":"300","workspace_id":"1"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	results, err := c.ListComments(&model.ListCommentsRequest{WorkspaceID: "1"})
	if err != nil {
		t.Fatalf("ListComments() unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(results))
	}
	if !strings.Contains(results[0].Description, "**bold**") {
		t.Errorf("description = %q, want to contain %q", results[0].Description, "**bold**")
	}
	if !strings.Contains(results[0].Description, "*italic*") {
		t.Errorf("description = %q, want to contain %q", results[0].Description, "*italic*")
	}
}

func TestAddComment(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/comments" {
			t.Errorf("unexpected path: %s, want /comments", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"Comment":{"id":"102","title":"在状态 [新] 添加","description":"test comment","author":"tester","entry_type":"stories","entry_id":"200","workspace_id":"1","created":"2026-04-16 11:00:00"}},"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	result, err := c.AddComment(&model.AddCommentRequest{
		WorkspaceID: "1",
		EntryType:   "stories",
		EntryID:     "200",
		Description: "test comment",
		Author:      "tester",
	})
	if err != nil {
		t.Fatalf("AddComment() unexpected error: %v", err)
	}
	if result.ID != "102" {
		t.Errorf("comment id = %q, want %q", result.ID, "102")
	}
	if result.Description != "test comment" {
		t.Errorf("comment description = %q, want %q", result.Description, "test comment")
	}
	if result.Author != "tester" {
		t.Errorf("comment author = %q, want %q", result.Author, "tester")
	}
}

func TestUpdateComment(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/comments" {
			t.Errorf("unexpected path: %s, want /comments", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		// 更新评论接口的 data 直接是评论对象，无 "Comment" 包裹
		w.Write([]byte(`{"status":1,"data":{"id":"102","title":"在状态 [新] 添加","description":"updated content","author":"tester","entry_type":"stories","entry_id":"200","workspace_id":"1","created":"2026-04-16 11:00:00","modified":"2026-04-16 12:00:00"},"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	result, err := c.UpdateComment(&model.UpdateCommentRequest{
		WorkspaceID: "1",
		ID:          "102",
		Description: "updated content",
	})
	if err != nil {
		t.Fatalf("UpdateComment() unexpected error: %v", err)
	}
	if result.ID != "102" {
		t.Errorf("comment id = %q, want %q", result.ID, "102")
	}
	if result.Description != "updated content" {
		t.Errorf("comment description = %q, want %q", result.Description, "updated content")
	}
	if result.Modified != "2026-04-16 12:00:00" {
		t.Errorf("comment modified = %q, want %q", result.Modified, "2026-04-16 12:00:00")
	}
}

func TestCountComments(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/comments/count" {
			t.Errorf("unexpected path: %s, want /comments/count", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"count":61},"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	count, err := c.CountComments(&model.CountCommentsRequest{
		WorkspaceID: "1",
	})
	if err != nil {
		t.Fatalf("CountComments() unexpected error: %v", err)
	}
	if count != 61 {
		t.Errorf("count = %d, want 61", count)
	}
}
