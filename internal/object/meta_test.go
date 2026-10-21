package object

import (
	"testing"
)

func TestBase(t *testing.T) {
	obj := NewMeta("file", "full/dir/path/file", "")
	want := "full/dir/path"
	if got := obj.BasePath; want != got {
		t.Errorf("Got %s, want %s", got, want)
	}
}

func TestMetaSymlinkAndScanRootOptions(t *testing.T) {
	obj := NewMeta("link", "dir/link", "application/octet-stream",
		WithMetaSymlink(true),
		WithMetaScanRoot("/photos"),
	)
	if !obj.IsSymlink {
		t.Errorf("expected IsSymlink to be true")
	}
	if obj.ScanRoot != "/photos" {
		t.Errorf("Got scan root %s, want /photos", obj.ScanRoot)
	}

	plain := NewMeta("file", "dir/file", "")
	if plain.IsSymlink {
		t.Errorf("expected IsSymlink to default to false")
	}
	if plain.ScanRoot != "" {
		t.Errorf("expected ScanRoot to default to empty, got %s", plain.ScanRoot)
	}
}
