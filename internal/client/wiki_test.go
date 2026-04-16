package client_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/studyzy/tapd-ai-cli/internal/client"
	"github.com/studyzy/tapd-ai-cli/internal/model"
)

func TestListWikis(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tapd_wikis" {
			t.Errorf("unexpected path: %s, want /tapd_wikis", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Wiki":{"id":"1151081496001001503","name":"技术架构","workspace_id":"51081496","creator":"张三"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	results, err := c.ListWikis(&model.ListWikisRequest{WorkspaceID: "51081496"})
	if err != nil {
		t.Fatalf("ListWikis() unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 wiki, got %d", len(results))
	}
	if results[0].ID != "1151081496001001503" {
		t.Errorf("wiki id = %v, want %q", results[0].ID, "1151081496001001503")
	}
	if results[0].Name != "技术架构" {
		t.Errorf("wiki name = %v, want %q", results[0].Name, "技术架构")
	}
}

func TestGetWiki_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tapd_wikis" {
			t.Errorf("unexpected path: %s, want /tapd_wikis", r.URL.Path)
		}
		if r.URL.Query().Get("id") != "1151081496001001503" {
			t.Errorf("expected id query param %q, got %q", "1151081496001001503", r.URL.Query().Get("id"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Wiki":{"id":"1151081496001001503","name":"技术架构","markdown_description":"# 技术架构\n\n本项目采用微服务架构。","creator":"张三","modified":"2026-04-01 10:00:00"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	result, err := c.GetWiki("51081496", "1151081496001001503")
	if err != nil {
		t.Fatalf("GetWiki() unexpected error: %v", err)
	}
	if result.ID != "1151081496001001503" {
		t.Errorf("id = %v, want %q", result.ID, "1151081496001001503")
	}
	if result.Name != "技术架构" {
		t.Errorf("name = %v, want %q", result.Name, "技术架构")
	}
	if result.MarkdownDescription == "" {
		t.Error("markdown_description should not be empty")
	}
}

func TestGetWiki_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[],"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	_, err := c.GetWiki("51081496", "9999999999")
	if err == nil {
		t.Fatal("GetWiki() expected error for not found, got nil")
	}
}

func TestCreateWiki(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/tapd_wikis" {
			t.Errorf("unexpected path: %s, want /tapd_wikis", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"Wiki":{"id":"1151081496001001600","name":"新文档","workspace_id":"51081496","creator":"testuser","markdown_description":"# Hello"}},"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	req := &model.CreateWikiRequest{
		WorkspaceID:         "51081496",
		Name:                "新文档",
		Creator:             "testuser",
		MarkdownDescription: "# Hello",
	}
	resp, err := c.CreateWiki(req)
	if err != nil {
		t.Fatalf("CreateWiki() unexpected error: %v", err)
	}
	if !resp.Success {
		t.Error("expected Success = true")
	}
	if resp.ID != "1151081496001001600" {
		t.Errorf("ID = %q, want %q", resp.ID, "1151081496001001600")
	}
	if resp.URL == "" {
		t.Error("URL should not be empty")
	}
}

func TestUpdateWiki(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/tapd_wikis" {
			t.Errorf("unexpected path: %s, want /tapd_wikis", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"Wiki":{"id":"1151081496001001503","name":"更新后的标题","workspace_id":"51081496","markdown_description":"# Updated","modifier":"testuser","modified":"2026-04-16 12:00:00"}},"info":"success"}`))
	}))
	defer srv.Close()

	c := client.NewClientWithBaseURL(srv.URL, "test-token", "", "")
	req := &model.UpdateWikiRequest{
		WorkspaceID:         "51081496",
		ID:                  "1151081496001001503",
		Name:                "更新后的标题",
		MarkdownDescription: "# Updated",
	}
	result, err := c.UpdateWiki(req)
	if err != nil {
		t.Fatalf("UpdateWiki() unexpected error: %v", err)
	}
	if result.ID != "1151081496001001503" {
		t.Errorf("ID = %q, want %q", result.ID, "1151081496001001503")
	}
	if result.Name != "更新后的标题" {
		t.Errorf("Name = %q, want %q", result.Name, "更新后的标题")
	}
	if result.MarkdownDescription != "# Updated" {
		t.Errorf("MarkdownDescription = %q, want %q", result.MarkdownDescription, "# Updated")
	}
}
