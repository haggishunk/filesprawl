package main

import "testing"

func TestParseIndexArgs(t *testing.T) {
	got, err := parseIndexArgs([]string{"--remote", "media:", "--path", "code/flux"})
	if err != nil {
		t.Fatalf("parseIndexArgs returned error: %v", err)
	}

	if got.Remote != "media:" {
		t.Fatalf("expected remote media:, got %s", got.Remote)
	}
	if got.Path != "code/flux" {
		t.Fatalf("expected path code/flux, got %s", got.Path)
	}
}

func TestParseIndexArgsRequiresRemote(t *testing.T) {
	_, err := parseIndexArgs([]string{"--path", "code/flux"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNormalizeRemoteName(t *testing.T) {
	if got := normalizeRemoteName("media:"); got != "media" {
		t.Fatalf("expected media, got %s", got)
	}
	if got := normalizeRemoteName(" media "); got != "media" {
		t.Fatalf("expected media, got %s", got)
	}
}

func TestToRcloneFS(t *testing.T) {
	if got := toRcloneFS("media"); got != "media:" {
		t.Fatalf("expected media:, got %s", got)
	}
	if got := toRcloneFS("media:"); got != "media:" {
		t.Fatalf("expected media:, got %s", got)
	}
}
