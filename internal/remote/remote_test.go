package remote

import "testing"

func TestNewRemoteDefaults(t *testing.T) {
	got := NewRemote("host1", "media:")

	if got.Id != -1 {
		t.Fatalf("expected default id -1, got %d", got.Id)
	}
	if got.Hostname != "host1" {
		t.Fatalf("expected hostname host1, got %s", got.Hostname)
	}
	if got.Name != "media:" {
		t.Fatalf("expected name media:, got %s", got.Name)
	}
	if got.Type != "rclone" {
		t.Fatalf("expected type rclone, got %s", got.Type)
	}
	if got.Persisted {
		t.Fatal("expected remote to start unpersisted")
	}
}
