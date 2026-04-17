package tapd

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newMockServer 创建一个模拟 TAPD API 的测试服务器
func newMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func TestNewClient_BearerAuth(t *testing.T) {
	var gotAuth string
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[],"info":"success"}`))
	})
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "my-token", "", "")
	err := c.TestAuth()
	if err != nil {
		t.Fatalf("TestAuth() unexpected error: %v", err)
	}
	expected := "Bearer my-token"
	if gotAuth != expected {
		t.Errorf("Authorization header = %q, want %q", gotAuth, expected)
	}
}

func TestNewClient_BasicAuth(t *testing.T) {
	var gotAuth string
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[],"info":"success"}`))
	})
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "", "user1", "pass1")
	err := c.TestAuth()
	if err != nil {
		t.Fatalf("TestAuth() unexpected error: %v", err)
	}
	expectedEncoded := base64.StdEncoding.EncodeToString([]byte("user1:pass1"))
	expected := "Basic " + expectedEncoded
	if gotAuth != expected {
		t.Errorf("Authorization header = %q, want %q", gotAuth, expected)
	}
}

func TestNewClient_AccessTokenPreferred(t *testing.T) {
	var gotAuth string
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[],"info":"success"}`))
	})
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "my-token", "user1", "pass1")
	err := c.TestAuth()
	if err != nil {
		t.Fatalf("TestAuth() unexpected error: %v", err)
	}
	expected := "Bearer my-token"
	if gotAuth != expected {
		t.Errorf("Authorization header = %q, want %q; access_token should take priority", gotAuth, expected)
	}
}

func TestTestAuth_Success(t *testing.T) {
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/quickstart/testauth" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":[],"info":"success"}`))
	})
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "token", "", "")
	err := c.TestAuth()
	if err != nil {
		t.Fatalf("TestAuth() expected nil error, got: %v", err)
	}
}

func TestTestAuth_Failure(t *testing.T) {
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":0,"data":[],"info":"auth failed"}`))
	})
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "bad-token", "", "")
	err := c.TestAuth()
	if err == nil {
		t.Fatal("TestAuth() expected error for status=0, got nil")
	}
	var tapdErr *TAPDError
	if !errors.As(err, &tapdErr) {
		t.Fatalf("expected *TAPDError, got %T", err)
	}
	if tapdErr.ExitCode != 4 {
		t.Errorf("ExitCode = %d, want 4", tapdErr.ExitCode)
	}
}

func TestHTTP401_ExitCode1(t *testing.T) {
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`Unauthorized`))
	})
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "bad-token", "", "")
	err := c.TestAuth()
	if err == nil {
		t.Fatal("expected error for HTTP 401")
	}
	var tapdErr *TAPDError
	if !errors.As(err, &tapdErr) {
		t.Fatalf("expected *TAPDError, got %T", err)
	}
	if tapdErr.ExitCode != 1 {
		t.Errorf("ExitCode = %d, want 1", tapdErr.ExitCode)
	}
}

func TestHTTP404_ExitCode2(t *testing.T) {
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`Not Found`))
	})
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "token", "", "")
	err := c.TestAuth()
	if err == nil {
		t.Fatal("expected error for HTTP 404")
	}
	var tapdErr *TAPDError
	if !errors.As(err, &tapdErr) {
		t.Fatalf("expected *TAPDError, got %T", err)
	}
	if tapdErr.ExitCode != 2 {
		t.Errorf("ExitCode = %d, want 2", tapdErr.ExitCode)
	}
}

func TestHTTP429_ExitCode4(t *testing.T) {
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`Rate limited`))
	})
	defer srv.Close()

	c := NewClientWithBaseURL(srv.URL, "token", "", "")
	err := c.TestAuth()
	if err == nil {
		t.Fatal("expected error for HTTP 429")
	}
	var tapdErr *TAPDError
	if !errors.As(err, &tapdErr) {
		t.Fatalf("expected *TAPDError, got %T", err)
	}
	if tapdErr.ExitCode != 4 {
		t.Errorf("ExitCode = %d, want 4", tapdErr.ExitCode)
	}
}

func TestDoGet_ParsesDataCorrectly(t *testing.T) {
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":1,"data":{"id":"12345","name":"test"},"info":"success"}`))
	})
	defer srv.Close()

	// 通过 TestAuth 验证 GET 请求能正常解析，TestAuth 内部调用 doGet
	// 这里只验证不会报错，因为 TestAuth 不返回 data
	c := NewClientWithBaseURL(srv.URL, "token", "", "")
	err := c.TestAuth()
	if err != nil {
		t.Fatalf("expected successful GET request, got error: %v", err)
	}
}

func TestTAPDError_Error(t *testing.T) {
	e := &TAPDError{
		HTTPStatus: 400,
		ExitCode:   3,
		Message:    "bad request",
	}
	got := e.Error()
	if got != "bad request" {
		t.Errorf("Error() = %q, want %q", got, "bad request")
	}
}
