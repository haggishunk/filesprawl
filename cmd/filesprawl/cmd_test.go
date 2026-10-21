package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

// TestResolveRemoteWithRemoteFlag exercises the resolve-remote CLI path with a
// remote-scoped rule selected via the --remote flag value.
func TestResolveRemoteWithRemoteFlag(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "locality.json")
	configJSON := `{
	  "rules": [
	    {"remotePattern": "/photos/*", "localPattern": "/backup/generic/*", "priority": 1},
	    {"remotePattern": "/photos/*", "localPattern": "/backup/s3-photos/*", "remoteName": "s3-prod", "priority": 10}
	  ]
	}`
	if err := os.WriteFile(configPath, []byte(configJSON), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Setenv("FILESPRAWL_LOCALITY_CONFIG", configPath)

	localityRemoteName = "s3-prod"
	defer func() { localityRemoteName = "" }()

	out := captureStdout(t, func() {
		if err := runResolveRemote(resolvRemoteCmd, []string{"/photos/cat.jpg"}); err != nil {
			t.Fatalf("runResolveRemote failed: %v", err)
		}
	})

	if !strings.Contains(out, "Remote name: s3-prod") {
		t.Fatalf("expected remote name in output, got: %s", out)
	}
	if !strings.Contains(out, "Local:  /backup/s3-photos/cat.jpg") {
		t.Fatalf("expected remote-scoped mapping in output, got: %s", out)
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = old }()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("close pipe writer: %v", err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read pipe: %v", err)
	}
	return string(out)
}
