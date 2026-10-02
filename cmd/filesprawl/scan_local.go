package main

import (
	"fmt"

	"github.com/haggishunk/filesprawl/internal/cli"
	"github.com/haggishunk/filesprawl/internal/operation"
	"github.com/spf13/cobra"
)

var scanLocalPath string

var scanLocalCmd = &cobra.Command{
	Use:   "scan-local",
	Short: "Index files from a local filesystem path",
	Long: `Index all files from the specified local filesystem path, computing hashes and
persisting metadata to the database.`,
	Example: `  # Index entire local directory
  filesprawl scan-local /home/user/documents

  # Index from environment variable
  LOCAL_PATH=/media filesprawl scan-local $LOCAL_PATH`,
	Args: cobra.ExactArgs(1),
	RunE: runScanLocalCmd,
}

func runScanLocalCmd(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	if len(args) == 0 {
		return fmt.Errorf("local path is required")
	}

	localPath := args[0]

	repo, cleanup, err := cli.OpenRepository(ctx)
	if err != nil {
		return err
	}
	defer cleanup()

	ls := operation.NewLocalScanner(
		operation.WithLocalRoot(localPath),
		operation.WithLocalRepo(repo),
	)

	if err := operation.ScanLocal(ctx, &ls); err != nil {
		return err
	}

	fmt.Printf("Successfully indexed local path: %s\n", localPath)
	return nil
}

func init() {
	rootCmd.AddCommand(scanLocalCmd)
}
