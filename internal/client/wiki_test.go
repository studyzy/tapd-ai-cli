package client_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/studyzy/tapd-ai-cli/internal/client"
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
	params := map[string]string{
		"workspace_id": "51081496",
	}
	results, err := c.ListWikis(params)
	if err != nil {
		t.Fatalf("ListWikis() unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 wiki, got %d", len(results))
	}
	if results[0]["id"] != "1151081496001001503" {
		t.Errorf("wiki id = %v, want %q", results[0]["id"], "1151081496001001503")
	}
	if results[0]["name"] != "技术架构" {
		t.Errorf("wiki name = %v, want %q", results[0]["name"], "技术架构")
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
	if result["id"] != "1151081496001001503" {
		t.Errorf("id = %v, want %q", result["id"], "1151081496001001503")
	}
	if result["name"] != "技术架构" {
		t.Errorf("name = %v, want %q", result["name"], "技术架构")
	}
	if result["markdown_description"] == nil {
		t.Error("markdown_description should not be nil")
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
