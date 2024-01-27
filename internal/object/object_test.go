package object

import (
	"testing"
)

func TestBase(t *testing.T) {
	obj := Object{
		Name: "file",
		Path: "full/dir/path/file",
	}
	want := "file"
	if got := obj.Base(); want != got {
		t.Errorf("Got %s, want %s", got, want)
	}
}
