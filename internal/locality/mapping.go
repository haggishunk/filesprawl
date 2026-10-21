package locality

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// MappingRule represents a single path mapping rule
type MappingRule struct {
	RemotePattern string `json:"remotePattern"` // e.g., "/documents/*"
	LocalPattern  string `json:"localPattern"`  // e.g., "/home/user/documents/*"
	RemoteName    string `json:"remoteName"`    // Optional: rule only applies to this remote
	RemoteRegex   string `json:"remoteRegex"`   // Optional: rule only applies to remotes matching this regex
	IsRegex       bool   `json:"isRegex"`       // Treat the path patterns as regular expressions
	Priority      int    `json:"priority"`      // Higher priority = more specific, applied first
}

// Config represents locality mapping configuration
type Config struct {
	Rules []MappingRule `json:"rules"`
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

// RemoteToLocal resolves a remote path to a local path. An empty remoteName
// restricts matching to rules that are not scoped to a particular remote.
func (m *Mapper) RemoteToLocal(remoteName, remotePath string) (string, error) {
	for _, rule := range m.rulesForRemote(remoteName) {
		if rule.IsRegex {
			if matchRegexPattern(rule.RemotePattern, remotePath) {
				return rule.LocalPattern, nil
			}
			continue
		}
		if matches, vars := matchPattern(rule.RemotePattern, remotePath); matches {
			return applyPattern(rule.LocalPattern, vars), nil
		}
	}
	return "", fmt.Errorf("no mapping found for remote path: %s", remotePath)
}

// LocalToRemote resolves a local path to remote paths. An empty remoteName
// restricts matching to rules that are not scoped to a particular remote.
func (m *Mapper) LocalToRemote(remoteName, localPath string) ([]string, error) {
	var results []string
	for _, rule := range m.rulesForRemote(remoteName) {
		if rule.IsRegex {
			if matchRegexPattern(rule.LocalPattern, localPath) {
				results = append(results, rule.RemotePattern)
			}
			continue
		}
		if matches, vars := matchPattern(rule.LocalPattern, localPath); matches {
			results = append(results, applyPattern(rule.RemotePattern, vars))
		}
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no mapping found for local path: %s", localPath)
	}
	return results, nil
}

// rulesForRemote returns the rules applicable to remoteName, ordered by
// descending priority. Rules with equal priority keep their configured order.
func (m *Mapper) rulesForRemote(remoteName string) []MappingRule {
	rules := make([]MappingRule, 0, len(m.config.Rules))
	for _, rule := range m.config.Rules {
		if ruleAppliesToRemote(rule, remoteName) {
			rules = append(rules, rule)
		}
	}
	sort.SliceStable(rules, func(i, j int) bool {
		return rules[i].Priority > rules[j].Priority
	})
	return rules
}

// ruleAppliesToRemote reports whether a rule is in scope for the given remote.
// Rules without remote constraints apply to every remote.
func ruleAppliesToRemote(rule MappingRule, remoteName string) bool {
	if rule.RemoteName == "" && rule.RemoteRegex == "" {
		return true
	}
	if remoteName == "" {
		return false
	}
	if rule.RemoteName != "" {
		return rule.RemoteName == remoteName
	}
	return matchRegexPattern(rule.RemoteRegex, remoteName)
}

// matchRegexPattern reports whether value matches the given regular expression.
// Invalid expressions never match; they are reported by ValidateConfig.
func matchRegexPattern(pattern, value string) bool {
	re, err := compileRegex(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(value)
}

var (
	regexCacheMu sync.RWMutex
	regexCache   = map[string]*regexp.Regexp{}
)

// compileRegex compiles and caches a regular expression, anchoring it so that
// patterns must match the whole value.
func compileRegex(pattern string) (*regexp.Regexp, error) {
	regexCacheMu.RLock()
	re, ok := regexCache[pattern]
	regexCacheMu.RUnlock()
	if ok {
		return re, nil
	}

	re, err := regexp.Compile("^(?:" + pattern + ")$")
	if err != nil {
		return nil, err
	}

	regexCacheMu.Lock()
	regexCache[pattern] = re
	regexCacheMu.Unlock()
	return re, nil
}

// ValidateConfig checks that every rule is well formed
func ValidateConfig(config Config) error {
	for i, rule := range config.Rules {
		if rule.RemotePattern == "" || rule.LocalPattern == "" {
			return fmt.Errorf("rule %d: remotePattern and localPattern are required", i+1)
		}
		if rule.RemoteName != "" && rule.RemoteRegex != "" {
			return fmt.Errorf("rule %d: remoteName and remoteRegex are mutually exclusive", i+1)
		}
		if rule.RemoteRegex != "" {
			if _, err := compileRegex(rule.RemoteRegex); err != nil {
				return fmt.Errorf("rule %d: invalid remoteRegex %q: %w", i+1, rule.RemoteRegex, err)
			}
		}
		if rule.IsRegex {
			if _, err := compileRegex(rule.RemotePattern); err != nil {
				return fmt.Errorf("rule %d: invalid remotePattern %q: %w", i+1, rule.RemotePattern, err)
			}
			if _, err := compileRegex(rule.LocalPattern); err != nil {
				return fmt.Errorf("rule %d: invalid localPattern %q: %w", i+1, rule.LocalPattern, err)
			}
		}
	}
	return nil
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
