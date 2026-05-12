package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/haggishunk/filesprawl/internal/rclone"
	"github.com/spf13/cobra"
)

// version is set at build time via -ldflags "-X main.version=<value>". The
// Makefile and release workflow strip any leading "v" so the embedded value is
// bare semver matching the published artifact.
var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "filesprawl",
	Short: "File indexing and duplicate detection across remote storage",
	Long: `Filesprawl indexes and analyzes files across multiple remote
storage locations, identifying duplicates and optimizing storage usage.`,
	Version:       version,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var listRemotesCmd = &cobra.Command{
	Use:   "list-remotes",
	Short: "List configured rclone remotes",
	Long: `List all rclone remote names as configured in the active local
rcd session running on http://localhost:5572.`,
	Example: `  filesprawl list-remotes`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runListRemotes(cmd.Context(), cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(listRemotesCmd)
}

// Execute runs the root command using the supplied context. It is the single
// entry point used by main and by integration tests that need to invoke the
// CLI programmatically.
func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}

func runListRemotes(ctx context.Context, stdout io.Writer) error {
	remotes, err := rclone.ListRemotes(ctx)
	if err != nil {
		return err
	}

	for _, name := range remotes.Remotes {
		if _, err := fmt.Fprintln(stdout, name); err != nil {
			return fmt.Errorf("failed to write remote list: %w", err)
		}
	}

	return nil
}

// normalizeRemoteName trims whitespace and the rclone trailing colon from a
// remote name so it can be compared and stored consistently.
func normalizeRemoteName(name string) string {
	return strings.TrimSuffix(strings.TrimSpace(name), ":")
}

// toRcloneFS converts a normalized remote name back to the colon-suffixed form
// expected by rclone fs arguments. An empty name yields an empty string.
func toRcloneFS(name string) string {
	trimmed := normalizeRemoteName(name)
	if trimmed == "" {
		return ""
	}
	return trimmed + ":"
}
