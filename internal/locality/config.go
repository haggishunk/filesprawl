package locality

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ConfigManager loads and caches locality configuration
type ConfigManager struct {
	cache       *cachedConfig
	cacheTTL    time.Duration
	localPath   string // Path to local .filesprawl/locality.json override
	mu          sync.RWMutex
}

type cachedConfig struct {
	config    Config
	timestamp time.Time
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(cacheTTL time.Duration) *ConfigManager {
	return &ConfigManager{
		cacheTTL: cacheTTL,
	}
}

// SetLocalConfigPath sets the path to local configuration
func (cm *ConfigManager) SetLocalConfigPath(path string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.localPath = path
}

// LoadConfig loads configuration from local file and/or remote
// Priority: local override > remote > empty config
func (cm *ConfigManager) LoadConfig(ctx context.Context) (Config, error) {
	cm.mu.RLock()
	if cm.cache != nil && time.Since(cm.cache.timestamp) < cm.cacheTTL {
		defer cm.mu.RUnlock()
		return cm.cache.config, nil
	}
	cm.mu.RUnlock()

	config := Config{Rules: []MappingRule{}}

	// Try to load local config first (higher priority)
	if cm.localPath != "" {
		localConfig, err := LoadConfigFromFile(cm.localPath)
		if err == nil {
			config = localConfig
		}
	}

	// Cache the result
	cm.mu.Lock()
	cm.cache = &cachedConfig{
		config:    config,
		timestamp: time.Now(),
	}
	cm.mu.Unlock()

	return config, nil
}

// InvalidateCache clears the configuration cache
func (cm *ConfigManager) InvalidateCache() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cache = nil
}

// LoadConfigFromFile loads configuration from a JSON file
func LoadConfigFromFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	return config, nil
}

// SaveConfigToFile saves configuration to a JSON file
func SaveConfigToFile(path string, config Config) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
