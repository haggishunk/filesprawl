package locality

import (
	"os"
	"path/filepath"
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
	local, err := mapper.RemoteToLocal("", "/documents")
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
	local, err := mapper.RemoteToLocal("", "/data/files/document.txt")
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
	local, err := mapper.RemoteToLocal("", "/documents/files")
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

	_, err := mapper.RemoteToLocal("", "/nonexistent/path")
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
	remotes, err := mapper.LocalToRemote("", "/home/user/Documents")
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

func TestRemoteSpecificRules(t *testing.T) {
	config := Config{
		Rules: []MappingRule{
			{
				RemotePattern: "/photos/*",
				LocalPattern:  "/backup/generic/*",
				Priority:      1,
			},
			{
				RemotePattern: "/photos/*",
				LocalPattern:  "/backup/s3-photos/*",
				RemoteName:    "s3-prod",
				Priority:      10,
			},
		},
	}

	mapper := NewMapper(config)

	// Remote-specific rule wins through higher priority
	local, err := mapper.RemoteToLocal("s3-prod", "/photos/cat.jpg")
	if err != nil {
		t.Fatalf("RemoteToLocal failed: %v", err)
	}
	if local != "/backup/s3-photos/cat.jpg" {
		t.Fatalf("expected /backup/s3-photos/cat.jpg, got %s", local)
	}

	// Other remotes fall back to the unscoped rule
	local, err = mapper.RemoteToLocal("google-drive", "/photos/cat.jpg")
	if err != nil {
		t.Fatalf("RemoteToLocal failed: %v", err)
	}
	if local != "/backup/generic/cat.jpg" {
		t.Fatalf("expected /backup/generic/cat.jpg, got %s", local)
	}

	// No remote name given skips remote-scoped rules
	local, err = mapper.RemoteToLocal("", "/photos/cat.jpg")
	if err != nil {
		t.Fatalf("RemoteToLocal failed: %v", err)
	}
	if local != "/backup/generic/cat.jpg" {
		t.Fatalf("expected /backup/generic/cat.jpg, got %s", local)
	}
}

func TestRemoteRegexRules(t *testing.T) {
	config := Config{
		Rules: []MappingRule{
			{
				RemotePattern: "/data/*",
				LocalPattern:  "/backup/prod/*",
				RemoteRegex:   "s3-.*",
				Priority:      1,
			},
		},
	}

	mapper := NewMapper(config)

	local, err := mapper.RemoteToLocal("s3-eu", "/data/file.txt")
	if err != nil {
		t.Fatalf("RemoteToLocal failed: %v", err)
	}
	if local != "/backup/prod/file.txt" {
		t.Fatalf("expected /backup/prod/file.txt, got %s", local)
	}

	if _, err := mapper.RemoteToLocal("gcs-eu", "/data/file.txt"); err == nil {
		t.Fatalf("expected no mapping for non-matching remote")
	}
}

func TestRegexPathPattern(t *testing.T) {
	config := Config{
		Rules: []MappingRule{
			{
				RemotePattern: `/archive/\d{4}/\d{2}/.*`,
				LocalPattern:  "/backup/archive",
				IsRegex:       true,
				Priority:      1,
			},
		},
	}

	mapper := NewMapper(config)

	local, err := mapper.RemoteToLocal("", "/archive/2026/10/report.pdf")
	if err != nil {
		t.Fatalf("RemoteToLocal failed: %v", err)
	}
	if local != "/backup/archive" {
		t.Fatalf("expected /backup/archive, got %s", local)
	}

	if _, err := mapper.RemoteToLocal("", "/archive/report.pdf"); err == nil {
		t.Fatalf("expected no mapping for non-matching path")
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		expectErr bool
	}{
		{
			name:   "valid simple rule",
			config: Config{Rules: []MappingRule{{RemotePattern: "/a", LocalPattern: "/b"}}},
		},
		{
			name:      "missing local pattern",
			config:    Config{Rules: []MappingRule{{RemotePattern: "/a"}}},
			expectErr: true,
		},
		{
			name:      "remoteName and remoteRegex together",
			config:    Config{Rules: []MappingRule{{RemotePattern: "/a", LocalPattern: "/b", RemoteName: "x", RemoteRegex: "y"}}},
			expectErr: true,
		},
		{
			name:      "invalid remote regex",
			config:    Config{Rules: []MappingRule{{RemotePattern: "/a", LocalPattern: "/b", RemoteRegex: "["}}},
			expectErr: true,
		},
		{
			name:      "invalid regex path pattern",
			config:    Config{Rules: []MappingRule{{RemotePattern: "(", LocalPattern: "/b", IsRegex: true}}},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		err := ValidateConfig(tt.config)
		if (err != nil) != tt.expectErr {
			t.Fatalf("%s: expected error=%v, got %v", tt.name, tt.expectErr, err)
		}
	}
}

func TestPriorityOrdering(t *testing.T) {
	config := Config{
		Rules: []MappingRule{
			{RemotePattern: "/data/*", LocalPattern: "/low/*", Priority: 1},
			{RemotePattern: "/data/*", LocalPattern: "/high/*", Priority: 10},
		},
	}
	mapper := NewMapper(config)

	// Higher priority wins even though it is listed later in the config
	local, err := mapper.RemoteToLocal("", "/data/file.txt")
	if err != nil {
		t.Fatalf("RemoteToLocal failed: %v", err)
	}
	if local != "/high/file.txt" {
		t.Fatalf("expected /high/file.txt, got %s", local)
	}

	// Equal priority keeps the configured order
	config.Rules[1].Priority = 1
	local, err = mapper.RemoteToLocal("", "/data/file.txt")
	if err != nil {
		t.Fatalf("RemoteToLocal failed: %v", err)
	}
	if local != "/low/file.txt" {
		t.Fatalf("expected /low/file.txt, got %s", local)
	}
}

func TestMixedConfigurations(t *testing.T) {
	config := Config{
		Rules: []MappingRule{
			{RemotePattern: `/archive/\d{4}/.*`, LocalPattern: "/backup/archive", IsRegex: true, Priority: 10},
			{RemotePattern: "/documents/*", LocalPattern: "/home/user/documents/*", Priority: 1},
		},
	}
	mapper := NewMapper(config)

	// Regex rule matches date-structured paths
	local, err := mapper.RemoteToLocal("", "/archive/2026/report.pdf")
	if err != nil {
		t.Fatalf("RemoteToLocal failed: %v", err)
	}
	if local != "/backup/archive" {
		t.Fatalf("expected /backup/archive, got %s", local)
	}

	// Simple rule with wildcard still works alongside regex rules
	local, err = mapper.RemoteToLocal("", "/documents/notes.txt")
	if err != nil {
		t.Fatalf("RemoteToLocal failed: %v", err)
	}
	if local != "/home/user/documents/notes.txt" {
		t.Fatalf("expected /home/user/documents/notes.txt, got %s", local)
	}
}

func TestLoadConfigFromFileValidation(t *testing.T) {
	dir := t.TempDir()

	valid := `{"rules": [{"remotePattern": "/a/*", "localPattern": "/b/*"}]}`
	validPath := filepath.Join(dir, "valid.json")
	if err := os.WriteFile(validPath, []byte(valid), 0644); err != nil {
		t.Fatalf("write valid config: %v", err)
	}
	config, err := LoadConfigFromFile(validPath)
	if err != nil {
		t.Fatalf("expected valid config to load, got %v", err)
	}
	if len(config.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(config.Rules))
	}

	invalid := `{"rules": [{"remotePattern": "/a/*", "localPattern": "/b/*", "remoteRegex": "["}]}`
	invalidPath := filepath.Join(dir, "invalid.json")
	if err := os.WriteFile(invalidPath, []byte(invalid), 0644); err != nil {
		t.Fatalf("write invalid config: %v", err)
	}
	if _, err := LoadConfigFromFile(invalidPath); err == nil {
		t.Fatalf("expected invalid config to be rejected")
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
