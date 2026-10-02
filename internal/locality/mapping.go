package locality

import (
	"fmt"
	"strings"
)

// MappingRule represents a single path mapping rule
type MappingRule struct {
	RemotePattern string // e.g., "/documents/*"
	LocalPattern  string // e.g., "/home/user/documents/*"
	Priority      int    // Higher priority = more specific, applied first
}

// Config represents locality mapping configuration
type Config struct {
	Rules []MappingRule
}

// Mapper handles bidirectional path resolution
type Mapper struct {
	config Config
}

// NewMapper creates a new locality mapper with the given configuration
func NewMapper(config Config) *Mapper {
	return &Mapper{
		config: config,
	}
}

// RemoteToLocal resolves a remote path to a local path
func (m *Mapper) RemoteToLocal(remotePath string) (string, error) {
	for _, rule := range m.config.Rules {
		if matches, vars := matchPattern(rule.RemotePattern, remotePath); matches {
			return applyPattern(rule.LocalPattern, vars), nil
		}
	}
	return "", fmt.Errorf("no mapping found for remote path: %s", remotePath)
}

// LocalToRemote resolves a local path to remote paths
func (m *Mapper) LocalToRemote(localPath string) ([]string, error) {
	var results []string
	for _, rule := range m.config.Rules {
		if matches, vars := matchPattern(rule.LocalPattern, localPath); matches {
			results = append(results, applyPattern(rule.RemotePattern, vars))
		}
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no mapping found for local path: %s", localPath)
	}
	return results, nil
}

// matchPattern checks if a path matches a pattern and extracts variables
func matchPattern(pattern, path string) (bool, map[string]string) {
	// Convert pattern to regex
	// {var} -> capture group
	parts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(parts) > len(pathParts) {
		return false, nil
	}

	vars := make(map[string]string)
	for i, part := range parts {
		if i >= len(pathParts) {
			return false, nil
		}

		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			// Variable
			varName := part[1 : len(part)-1]
			vars[varName] = pathParts[i]
		} else if part == "*" {
			// Wildcard matches rest of path
			vars["_rest"] = strings.Join(pathParts[i:], "/")
			return true, vars
		} else if part != pathParts[i] {
			// Literal mismatch
			return false, nil
		}
	}

	return true, vars
}

// applyPattern applies variable substitution to a pattern
func applyPattern(pattern string, vars map[string]string) string {
	result := pattern
	for key, val := range vars {
		if key == "_rest" {
			result = strings.TrimSuffix(result, "*")
			result = strings.TrimSuffix(result, "/")
			result += "/" + val
		} else {
			result = strings.ReplaceAll(result, "{"+key+"}", val)
		}
	}
	return result
}
