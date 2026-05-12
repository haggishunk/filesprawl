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
