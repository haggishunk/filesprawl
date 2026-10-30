package remote

import (
	"testing"

	"github.com/haggishunk/filesprawl/internal/location"
)

func TestNewRcloneRemote(t *testing.T) {
	got := NewRcloneRemote("media", "")

	if got.Id != -1 {
		t.Fatalf("expected default id -1, got %d", got.Id)
	}
	if got.Name != "" {
		t.Fatalf("expected name empty, got %s", got.Name)
	}
	if !got.Location.IsRemote() {
		t.Fatal("expected location to be remote type")
	}
	if got.Location.Name != "media" {
		t.Fatalf("expected location name media, got %s", got.Location.Name)
	}
	if got.Persisted {
		t.Fatal("expected remote to start unpersisted")
	}
}

func TestNewLocalRemote(t *testing.T) {
	got := NewLocalRemote("archlinux", "/media/photos")

	if got.Id != -1 {
		t.Fatalf("expected default id -1, got %d", got.Id)
	}
	if got.Name != "/media/photos" {
		t.Fatalf("expected name /media/photos, got %s", got.Name)
	}
	if !got.Location.IsLocal() {
		t.Fatal("expected location to be local type")
	}
	if got.Location.Name != "archlinux" {
		t.Fatalf("expected location name archlinux, got %s", got.Location.Name)
	}
	if got.Persisted {
		t.Fatal("expected remote to start unpersisted")
	}
}

func TestNewRemoteWithLocation(t *testing.T) {
	loc := location.RemoteLocation("dropbox")
	got := NewRemote(loc, "media:")

	if got.Name != "media:" {
		t.Fatalf("expected name media:, got %s", got.Name)
	}
	if got.Location.Name != "dropbox" {
		t.Fatalf("expected location name dropbox, got %s", got.Location.Name)
	}
}
