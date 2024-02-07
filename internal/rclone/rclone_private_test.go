package rclone

import "testing"

func TestListThem(t *testing.T) {
	src := "remote1:/"
	ListThem([]string{src})
}
