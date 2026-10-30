package location

import "testing"

func TestNewLocation(t *testing.T) {
	loc := NewLocation("archlinux", "local")
	
	if loc.Id != -1 {
		t.Errorf("expected id -1, got %d", loc.Id)
	}
	if loc.Name != "archlinux" {
		t.Errorf("expected name archlinux, got %s", loc.Name)
	}
	if loc.Type != "local" {
		t.Errorf("expected type local, got %s", loc.Type)
	}
	if loc.Persisted {
		t.Error("expected persisted false")
	}
}

func TestLocalLocation(t *testing.T) {
	loc := LocalLocation("host1")
	
	if loc.Name != "host1" {
		t.Errorf("expected name host1, got %s", loc.Name)
	}
	if loc.Type != "local" {
		t.Errorf("expected type local, got %s", loc.Type)
	}
	if !loc.IsLocal() {
		t.Error("expected IsLocal() to return true")
	}
	if loc.IsRemote() {
		t.Error("expected IsRemote() to return false")
	}
}

func TestRemoteLocation(t *testing.T) {
	loc := RemoteLocation("dropbox")
	
	if loc.Name != "dropbox" {
		t.Errorf("expected name dropbox, got %s", loc.Name)
	}
	if loc.Type != "remote" {
		t.Errorf("expected type remote, got %s", loc.Type)
	}
	if loc.IsLocal() {
		t.Error("expected IsLocal() to return false")
	}
	if !loc.IsRemote() {
		t.Error("expected IsRemote() to return true")
	}
}

func TestNewLocationWithOptions(t *testing.T) {
	loc := NewLocation("host1", "local", WithLocationID(42))
	
	if loc.Id != 42 {
		t.Errorf("expected id 42, got %d", loc.Id)
	}
}
