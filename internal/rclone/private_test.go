package rclone

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListRemotes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/config/listremotes" {
			t.Fatalf("expected config/listremotes path, got %s", r.URL.Path)
		}
		fmt.Fprint(w, `{"remotes":["media:","backup:"]}`)
	}))
	defer server.Close()

	oldBaseURL := baseURL
	baseURL = server.URL
	t.Cleanup(func() {
		baseURL = oldBaseURL
	})

	got, err := ListRemotes(context.Background())
	if err != nil {
		t.Fatalf("ListRemotes returned error: %v", err)
	}

	if len(got.Remotes) != 2 {
		t.Fatalf("expected 2 remotes, got %d", len(got.Remotes))
	}
	if got.Remotes[0] != "media:" || got.Remotes[1] != "backup:" {
		t.Fatalf("unexpected remotes: %+v", got.Remotes)
	}
}

func TestListRemotesError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"error":"boom"}`)
	}))
	defer server.Close()

	oldBaseURL := baseURL
	baseURL = server.URL
	t.Cleanup(func() {
		baseURL = oldBaseURL
	})

	_, err := ListRemotes(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); got == "" || got == "boom" {
		t.Fatalf("expected wrapped error, got %q", got)
	}
}
