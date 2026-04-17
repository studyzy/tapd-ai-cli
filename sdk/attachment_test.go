package tapd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/studyzy/tapd-sdk-go/model"
)

func TestGetImage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/files/get_image" {
			t.Errorf("unexpected path: %s, want /files/get_image", r.URL.Path)
		}
		if r.URL.Query().Get("image_path") != "/tfl/pictures/abc.png" {
			t.Errorf("image_path = %q, want %q", r.URL.Query().Get("image_path"), "/tfl/pictures/abc.png")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"Attachment":{"type":"image","value":"abc.png","workspace_id":"1","filename":"abc.png","download_url":"https://example.com/abc.png"}},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "test-token", "", "")
	result, err := c.GetImage(&model.GetImageRequest{
		WorkspaceID: "1",
		ImagePath:   "/tfl/pictures/abc.png",
	})
	if err != nil {
		t.Fatalf("GetImage() unexpected error: %v", err)
	}
	if result.Filename != "abc.png" {
		t.Errorf("filename = %q, want %q", result.Filename, "abc.png")
	}
	if result.DownloadURL != "https://example.com/abc.png" {
		t.Errorf("download_url = %q, want %q", result.DownloadURL, "https://example.com/abc.png")
	}
	if result.Type != "image" {
		t.Errorf("type = %q, want %q", result.Type, "image")
	}
}

func TestGetImage_MissingAttachment(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// 返回不含 Attachment key 的响应
		w.Write([]byte(`{"status":1,"data":{"other":"value"},"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "test-token", "", "")
	_, err := c.GetImage(&model.GetImageRequest{
		WorkspaceID: "1",
		ImagePath:   "/tfl/pictures/abc.png",
	})
	if err == nil {
		t.Fatal("expected error for missing Attachment key")
	}
}

func TestGetAttachments(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/attachments" {
			t.Errorf("unexpected path: %s, want /attachments", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[{"Attachment":{"id":"501","type":"story","entry_id":"200","filename":"design.pdf","content_type":"application/pdf","owner":"uploader","workspace_id":"1","download_url":"https://example.com/design.pdf"}}],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "test-token", "", "")
	results, err := c.GetAttachments(&model.GetAttachmentsRequest{
		WorkspaceID: "1",
		EntryID:     "200",
	})
	if err != nil {
		t.Fatalf("GetAttachments() unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(results))
	}
	if results[0].ID != "501" {
		t.Errorf("attachment id = %q, want %q", results[0].ID, "501")
	}
	if results[0].Filename != "design.pdf" {
		t.Errorf("filename = %q, want %q", results[0].Filename, "design.pdf")
	}
	if results[0].DownloadURL != "https://example.com/design.pdf" {
		t.Errorf("download_url = %q, want %q", results[0].DownloadURL, "https://example.com/design.pdf")
	}
	if results[0].Owner != "uploader" {
		t.Errorf("owner = %q, want %q", results[0].Owner, "uploader")
	}
}

func TestGetAttachments_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[],"info":"success"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "test-token", "", "")
	results, err := c.GetAttachments(&model.GetAttachmentsRequest{
		WorkspaceID: "1",
		EntryID:     "200",
	})
	if err != nil {
		t.Fatalf("GetAttachments() unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 attachments, got %d", len(results))
	}
}

func TestGetAttachments_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":0,"data":null,"info":"permission denied"}`))
	}))
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "test-token", "", "")
	_, err := c.GetAttachments(&model.GetAttachmentsRequest{
		WorkspaceID: "1",
		EntryID:     "200",
	})
	if err == nil {
		t.Fatal("expected error for status=0")
	}
}
