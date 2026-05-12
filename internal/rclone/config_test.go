package rclone

import "testing"

func TestNewListConfig(t *testing.T) {
	opt := &ListOption{
		DirsOnly:      false,
		FilesOnly:     false,
		HashTypes:     []string{"dropbox", "md5"},
		NoMimeType:    true,
		NoModTime:     true,
		Recurse:       false,
		ShowEncrypted: true,
		ShowHash:      true,
		ShowOrigIDs:   true,
	}
	want := ListConfig{
		Fs:     "remote:",
		Remote: "path/to/list",
		Opt:    opt,
	}
	got := NewListConfig("remote:", "path/to/list", opt)
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
