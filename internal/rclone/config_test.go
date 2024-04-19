package rclone

import "testing"

func TestNewListConfig(t *testing.T) {
	want := ListConfig{
		Fs:     "remote:",
		Remote: "path/to/list",
		Opt: &ListOption{
			DirsOnly:      false,
			FilesOnly:     false,
			HashTypes:     []string{"dropbox", "md5"},
			NoMimeType:    true,
			NoModTime:     true,
			Recurse:       false,
			ShowEncrypted: true,
			ShowHash:      true,
			ShowOrigIDs:   true,
		},
	}
	got := NewListConfig("remote:", "path/to/list")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
