package main

import (
	"fmt"
	"os"

	"github.com/haggishunk/filesprawl/internal/locality"
	"github.com/spf13/cobra"
)

var localityCmd = &cobra.Command{
	Use:   "locality",
	Short: "Manage locality mappings between remote and local paths",
	Long: `Manage path mappings that define how remote storage paths correspond to
local filesystem paths.`,
}

var resolvRemoteCmd = &cobra.Command{
	Use:   "resolve-remote",
	Short: "Resolve a remote path to local path",
	Args:  cobra.ExactArgs(1),
	RunE:  runResolveRemote,
}

var resolveLocalCmd = &cobra.Command{
	Use:   "resolve-local",
	Short: "Resolve a local path to remote paths",
	Args:  cobra.ExactArgs(1),
	RunE:  runResolveLocal,
}

var showConfigCmd = &cobra.Command{
	Use:   "show-config",
	Short: "Show current locality configuration",
	RunE:  runShowConfig,
}

func runResolveRemote(cmd *cobra.Command, args []string) error {
	remotePath := args[0]

	// Load configuration from default location
	configPath := os.ExpandEnv("$HOME/.filesprawl/locality.json")
	config, err := loadConfigFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	mapper := locality.NewMapper(config)
	localPath, err := mapper.RemoteToLocal(remotePath)
	if err != nil {
		return err
	}

	fmt.Printf("Remote: %s\n", remotePath)
	fmt.Printf("Local:  %s\n", localPath)
	return nil
}

func runResolveLocal(cmd *cobra.Command, args []string) error {
	localPath := args[0]

	// Load configuration from default location
	configPath := os.ExpandEnv("$HOME/.filesprawl/locality.json")
	config, err := loadConfigFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	mapper := locality.NewMapper(config)
	remotePaths, err := mapper.LocalToRemote(localPath)
	if err != nil {
		return err
	}

	fmt.Printf("Local: %s\n", localPath)
	fmt.Println("Remote paths:")
	for _, rp := range remotePaths {
		fmt.Printf("  - %s\n", rp)
	}
	return nil
}

func runShowConfig(cmd *cobra.Command, args []string) error {
	configPath := os.ExpandEnv("$HOME/.filesprawl/locality.json")
	config, err := loadConfigFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Printf("Locality Configuration: %s\n", configPath)
	fmt.Printf("Rules: %d\n", len(config.Rules))
	for i, rule := range config.Rules {
		fmt.Printf("  [%d] %s -> %s (priority: %d)\n",
			i+1, rule.RemotePattern, rule.LocalPattern, rule.Priority)
	}
	return nil
}

func loadConfigFile(path string) (locality.Config, error) {
	if _, err := os.Stat(path); err != nil {
		// File doesn't exist, return empty config
		return locality.Config{Rules: []locality.MappingRule{}}, nil
	}
	// File exists, load it
	return locality.LoadConfigFromFile(path)
}

func init() {
	rootCmd.AddCommand(localityCmd)
	localityCmd.AddCommand(resolvRemoteCmd, resolveLocalCmd, showConfigCmd)
}
