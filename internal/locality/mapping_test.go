package locality

import (
	"os"
	"testing"
)

func TestSimplePathMapping(t *testing.T) {
	config := Config{
		Rules: []MappingRule{
			{
				RemotePattern: "/documents",
				LocalPattern:  "/home/user/Documents",
				Priority:      1,
			},
		},
	}

	mapper := NewMapper(config)

	// Test remote to local
	local, err := mapper.RemoteToLocal("/documents")
	if err != nil {
		t.Fatalf("RemoteToLocal failed: %v", err)
	}
	if local != "/home/user/Documents" {
		t.Fatalf("expected /home/user/Documents, got %s", local)
	}
}

func TestWildcardPatternMapping(t *testing.T) {
	config := Config{
		Rules: []MappingRule{
			{
				RemotePattern: "/data/*",
				LocalPattern:  "/home/user/data/*",
				Priority:      1,
			},
		},
	}

	mapper := NewMapper(config)

	// Test with wildcard
	local, err := mapper.RemoteToLocal("/data/files/document.txt")
	if err != nil {
		t.Fatalf("RemoteToLocal failed: %v", err)
	}
	if local != "/home/user/data/files/document.txt" {
		t.Fatalf("expected /home/user/data/files/document.txt, got %s", local)
	}
}

func TestVariableSubstitution(t *testing.T) {
	config := Config{
		Rules: []MappingRule{
			{
				RemotePattern: "/{category}/files",
				LocalPattern:  "/home/user/{category}",
				Priority:      1,
			},
		},
	}

	mapper := NewMapper(config)

	// Test variable matching
	local, err := mapper.RemoteToLocal("/documents/files")
	if err != nil {
		t.Fatalf("RemoteToLocal failed: %v", err)
	}
	if local != "/home/user/documents" {
		t.Fatalf("expected /home/user/documents, got %s", local)
	}
}

func TestRemoteToLocalNotFound(t *testing.T) {
	config := Config{
		Rules: []MappingRule{
			{
				RemotePattern: "/documents",
				LocalPattern:  "/home/user/Documents",
				Priority:      1,
			},
		},
	}

	mapper := NewMapper(config)

	_, err := mapper.RemoteToLocal("/nonexistent/path")
	if err == nil {
		t.Fatalf("expected error for unmapped path, got nil")
	}
}

func TestLocalToRemote(t *testing.T) {
	config := Config{
		Rules: []MappingRule{
			{
				RemotePattern: "/documents",
				LocalPattern:  "/home/user/Documents",
				Priority:      1,
			},
		},
	}

	mapper := NewMapper(config)

	// Test local to remote
	remotes, err := mapper.LocalToRemote("/home/user/Documents")
	if err != nil {
		t.Fatalf("LocalToRemote failed: %v", err)
	}
	if len(remotes) != 1 {
		t.Fatalf("expected 1 remote, got %d", len(remotes))
	}
	if remotes[0] != "/documents" {
		t.Fatalf("expected /documents, got %s", remotes[0])
	}
}

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		pattern      string
		path         string
		shouldMatch  bool
		expectedVars map[string]string
	}{
		{"/documents", "/documents", true, map[string]string{}},
		{"/documents", "/data", false, nil},
		{"/data/*", "/data/file.txt", true, map[string]string{"_rest": "file.txt"}},
		{"/{category}/files", "/docs/files", true, map[string]string{"category": "docs"}},
	}

	for _, tt := range tests {
		matches, _ := matchPattern(tt.pattern, tt.path)
		if matches != tt.shouldMatch {
			t.Fatalf("pattern %s against path %s: expected %v, got %v",
				tt.pattern, tt.path, tt.shouldMatch, matches)
		}
	}
}

func TestGetDefaultConfigPath(t *testing.T) {
	// Save original env vars
	origEnv := os.Getenv("FILESPRAWL_LOCALITY_CONFIG")
	origXDG := os.Getenv("XDG_CONFIG_HOME")
	origHome := os.Getenv("HOME")
	defer func() {
		if origEnv != "" {
			os.Setenv("FILESPRAWL_LOCALITY_CONFIG", origEnv)
		} else {
			os.Unsetenv("FILESPRAWL_LOCALITY_CONFIG")
		}
		if origXDG != "" {
			os.Setenv("XDG_CONFIG_HOME", origXDG)
		} else {
			os.Unsetenv("XDG_CONFIG_HOME")
		}
		if origHome != "" {
			os.Setenv("HOME", origHome)
		}
	}()

	// Test 1: Explicit FILESPRAWL_LOCALITY_CONFIG env var takes priority
	testPath := "/custom/path/locality.json"
	os.Setenv("FILESPRAWL_LOCALITY_CONFIG", testPath)
	os.Unsetenv("XDG_CONFIG_HOME")
	if path := GetDefaultConfigPath(); path != testPath {
		t.Fatalf("expected custom path %s, got %s", testPath, path)
	}

	// Test 2: XDG_CONFIG_HOME is checked (but file may not exist)
	os.Unsetenv("FILESPRAWL_LOCALITY_CONFIG")
	os.Setenv("XDG_CONFIG_HOME", "/tmp")
	path := GetDefaultConfigPath()
	// Should include filesprawl/locality.json
	if path == "" {
		t.Fatalf("expected path, got empty string")
	}

	// Test 3: Fall back to $HOME/.config/filesprawl/locality.json
	os.Unsetenv("FILESPRAWL_LOCALITY_CONFIG")
	os.Unsetenv("XDG_CONFIG_HOME")
	os.Setenv("HOME", "/home/testuser")
	path = GetDefaultConfigPath()
	if path == "" {
		t.Fatalf("expected default path, got empty string")
	}
	// Should end with .filesprawl/locality.json or .config/filesprawl/locality.json
	if !containsSubstring(path, "filesprawl") {
		t.Fatalf("expected path to contain 'filesprawl', got %s", path)
	}
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
