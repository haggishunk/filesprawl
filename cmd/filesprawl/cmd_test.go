package main

import "testing"

func TestNormalizeRemoteName(t *testing.T) {
	if got := normalizeRemoteName("media:"); got != "media" {
		t.Fatalf("expected media, got %s", got)
	}
	if got := normalizeRemoteName(" media "); got != "media" {
		t.Fatalf("expected media, got %s", got)
	}
}

func TestToRcloneFS(t *testing.T) {
	if got := toRcloneFS("media"); got != "media:" {
		t.Fatalf("expected media:, got %s", got)
	}
	if got := toRcloneFS("media:"); got != "media:" {
		t.Fatalf("expected media:, got %s", got)
	}
}

// TestRootCmdVersionWired guards against regressions where rootCmd.Version
// stops reflecting the package-level `version` symbol that ldflags injects.
// Without it, --version would silently report an empty string in release
// builds even though the linker substitution succeeded.
func TestRootCmdVersionWired(t *testing.T) {
	if rootCmd.Version == "" {
		t.Fatal("rootCmd.Version is empty; ldflags wiring is broken")
	}
	if rootCmd.Version != version {
		t.Fatalf("rootCmd.Version (%q) does not match version (%q); declaration order may have changed", rootCmd.Version, version)
	}
}
