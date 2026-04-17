package tapd

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListCategories_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/story_categories" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("workspace_id") != "12345" {
			t.Errorf("unexpected workspace_id: %s", r.URL.Query().Get("workspace_id"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":1,"data":[{"Category":{"id":"100","name":"前端需求","description":"前端相关","parent_id":"0","workspace_id":"12345"}},{"Category":{"id":"101","name":"后端需求","description":"后端相关","parent_id":"0","workspace_id":"12345"}}]}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "token123", "", "")
	categories, err := c.ListCategories(map[string]string{"workspace_id": "12345"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(categories) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(categories))
	}
	if categories[0].ID != "100" || categories[0].Name != "前端需求" {
		t.Errorf("unexpected first category: %+v", categories[0])
	}
	if categories[1].ID != "101" || categories[1].Name != "后端需求" {
		t.Errorf("unexpected second category: %+v", categories[1])
	}
}

func TestListCategories_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":1,"data":[]}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "token123", "", "")
	categories, err := c.ListCategories(map[string]string{"workspace_id": "12345"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(categories) != 0 {
		t.Fatalf("expected 0 categories, got %d", len(categories))
	}
}

func TestListCategories_WithNameFilter(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("name") != "前端" {
			t.Errorf("expected name filter '前端', got '%s'", r.URL.Query().Get("name"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":1,"data":[{"Category":{"id":"100","name":"前端需求"}}]}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "token123", "", "")
	categories, err := c.ListCategories(map[string]string{"workspace_id": "12345", "name": "前端"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(categories) != 1 {
		t.Fatalf("expected 1 category, got %d", len(categories))
	}
}

func TestListCategories_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":0,"info":"invalid workspace"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "token123", "", "")
	_, err := c.ListCategories(map[string]string{"workspace_id": "invalid"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
