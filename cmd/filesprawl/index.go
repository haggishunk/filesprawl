package main

import (
	"fmt"
	"os"

	"github.com/haggishunk/filesprawl/internal/cli"
	"github.com/haggishunk/filesprawl/internal/operation"
	"github.com/haggishunk/filesprawl/internal/rclone"
	"github.com/haggishunk/filesprawl/internal/remote"
	"github.com/spf13/cobra"
)

var (
	indexRemote string
	indexPath   string
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index files from a remote storage location",
	Long: `Index all files from the specified remote, computing hashes and
persisting metadata to the database.`,
	Example: `  # Index entire remote
  filesprawl index --remote media

  # Index specific path within remote
  filesprawl index --remote media --path code/flux`,
	RunE: runIndexCmd,
}

func init() {
	indexCmd.Flags().StringVarP(&indexRemote, "remote", "r", "", "remote name (required)")
	indexCmd.Flags().StringVarP(&indexPath, "path", "p", "", "path within remote")
	_ = indexCmd.MarkFlagRequired("remote")

	rootCmd.AddCommand(indexCmd)
}

func runIndexCmd(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	repo, cleanup, err := cli.OpenRepository(ctx)
	if err != nil {
		return err
	}
	defer cleanup()

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to determine hostname: %w", err)
	}

	remoteName := normalizeRemoteName(indexRemote)
	if remoteName == "" {
		return fmt.Errorf("remote name is required")
	}

	rem := remote.NewRemote(hostname, remoteName)
	scn := operation.NewScanner(
		operation.WithRemote(&rem),
		operation.WithRepo(repo),
	)

	lo := rclone.NewListOption(rclone.ListOptionFilesOnly())
	lc := rclone.NewListConfig(toRcloneFS(rem.Name), indexPath, &lo)
	if err := operation.Scan(ctx, &scn, lc); err != nil {
		return err
	}

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "indexed remote %s from host %s\n", rem.Name, rem.Hostname)
	if err != nil {
		return fmt.Errorf("failed to write index result: %w", err)
	}

	return nil
}
