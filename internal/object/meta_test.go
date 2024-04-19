package object

import (
	"testing"
)

func TestBase(t *testing.T) {
	obj := Meta{
		Name: "file",
		Path: "full/dir/path/file",
	}
	want := "file/dir/path"
	if got := obj.BasePath; want != got {
		t.Errorf("Got %s, want %s", got, want)
	}
}
