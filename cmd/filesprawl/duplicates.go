package main

import (
	"fmt"

	"github.com/haggishunk/filesprawl/internal/analysis"
	"github.com/haggishunk/filesprawl/internal/cli"
	"github.com/spf13/cobra"
)

var (
	duplicatesRemote        string
	duplicatesAcrossRemotes bool
	duplicatesHashType      string
	duplicatesMinSize       int64
	duplicatesLimit         int
	duplicatesOffset        int
)

var duplicatesCmd = &cobra.Command{
	Use:   "duplicates",
	Short: "Find duplicate files",
	Long: `Find files with identical content hashes, either within a single
remote or across all indexed remotes.`,
	Example: `  # Find duplicates within a remote
  filesprawl duplicates --remote media

  # Find duplicates across all remotes
  filesprawl duplicates --across-remotes

  # Find large duplicates (>10MB)
  filesprawl duplicates --remote media --min-size 10485760

  # Filter by hash type and limit results
  filesprawl duplicates --across-remotes --hash-type sha256 --limit 20`,
	PreRunE: validateDuplicatesArgs,
	RunE:    runDuplicatesCmd,
}

func init() {
	duplicatesCmd.Flags().StringVarP(&duplicatesRemote, "remote", "r", "", "find duplicates within this remote")
	duplicatesCmd.Flags().BoolVar(&duplicatesAcrossRemotes, "across-remotes", false, "find duplicates across all remotes")
	duplicatesCmd.Flags().StringVar(&duplicatesHashType, "hash-type", "", "filter by hash type (md5, sha1, sha256, dropbox)")
	duplicatesCmd.Flags().Int64Var(&duplicatesMinSize, "min-size", 0, "minimum file size in bytes")
	duplicatesCmd.Flags().IntVar(&duplicatesLimit, "limit", 0, "maximum number of duplicate groups to return")
	duplicatesCmd.Flags().IntVar(&duplicatesOffset, "offset", 0, "skip this many duplicate groups")

	rootCmd.AddCommand(duplicatesCmd)
}

func validateDuplicatesArgs(_ *cobra.Command, _ []string) error {
	duplicatesRemote = normalizeRemoteName(duplicatesRemote)

	if duplicatesRemote != "" && duplicatesAcrossRemotes {
		return fmt.Errorf("--remote and --across-remotes are mutually exclusive")
	}
	if duplicatesRemote == "" && !duplicatesAcrossRemotes {
		return fmt.Errorf("must specify either --remote or --across-remotes")
	}

	if duplicatesHashType != "" {
		validHashTypes := map[string]bool{"md5": true, "sha1": true, "sha256": true, "dropbox": true}
		if !validHashTypes[duplicatesHashType] {
			return fmt.Errorf("invalid hash type %q, must be one of: md5, sha1, sha256, dropbox", duplicatesHashType)
		}
	}

	if duplicatesMinSize < 0 {
		return fmt.Errorf("min-size must be non-negative, got %d", duplicatesMinSize)
	}
	if duplicatesLimit < 0 {
		return fmt.Errorf("limit must be non-negative, got %d", duplicatesLimit)
	}
	if duplicatesOffset < 0 {
		return fmt.Errorf("offset must be non-negative, got %d", duplicatesOffset)
	}

	return nil
}

func runDuplicatesCmd(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	repo, cleanup, err := cli.OpenRepository(ctx)
	if err != nil {
		return err
	}
	defer cleanup()

	detector := analysis.NewDuplicateDetector(repo)

	filterOpts := analysis.FilterOptions{
		HashType: duplicatesHashType,
		MinSize:  duplicatesMinSize,
		Limit:    duplicatesLimit,
		Offset:   duplicatesOffset,
	}

	var groups []analysis.DuplicateGroup
	if duplicatesAcrossRemotes {
		groups, err = detector.FindAcrossRemotes(ctx, filterOpts)
		if err != nil {
			return fmt.Errorf("failed to find duplicates across remotes: %w", err)
		}
	} else {
		groups, err = detector.FindWithinRemote(ctx, duplicatesRemote, filterOpts)
		if err != nil {
			return fmt.Errorf("failed to find duplicates within remote %q: %w", duplicatesRemote, err)
		}
	}

	report := analysis.FormatReport(groups)
	if _, err := fmt.Fprint(cmd.OutOrStdout(), report); err != nil {
		return fmt.Errorf("failed to write duplicate report: %w", err)
	}

	return nil
}
