package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/haggishunk/filesprawl/internal/database"
	"github.com/haggishunk/filesprawl/internal/operation"
	"github.com/haggishunk/filesprawl/internal/rclone"
	"github.com/haggishunk/filesprawl/internal/remote"
	"github.com/haggishunk/filesprawl/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type indexArgs struct {
	Remote string
	Path   string
}

func main() {
	os.Exit(run(context.Background(), os.Args[1:], os.Stdout, os.Stderr))
}

func run(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer) int {
	err := runCommand(ctx, args, stdout)
	if err == nil {
		return 0
	}

	fmt.Fprintln(stderr, err)
	return 1
}

func runCommand(ctx context.Context, args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: filesprawl <list-remotes|index>")
	}

	switch args[0] {
	case "list-remotes":
		return runListRemotes(ctx, stdout)
	case "index":
		parsed, err := parseIndexArgs(args[1:])
		if err != nil {
			return err
		}
		return runIndex(ctx, parsed, stdout)
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func parseIndexArgs(args []string) (indexArgs, error) {
	fs := flag.NewFlagSet("index", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	remoteName := fs.String("remote", "", "local rclone remote name")
	remotePath := fs.String("path", "", "path within the selected remote")

	err := fs.Parse(args)
	if err != nil {
		return indexArgs{}, fmt.Errorf("failed to parse index arguments: %w", err)
	}
	if fs.NArg() != 0 {
		return indexArgs{}, fmt.Errorf("unexpected positional arguments: %v", fs.Args())
	}
	*remoteName = strings.TrimSpace(*remoteName)
	*remotePath = strings.TrimSpace(*remotePath)
	if *remoteName == "" {
		return indexArgs{}, fmt.Errorf("remote name is required")
	}

	return indexArgs{Remote: *remoteName, Path: *remotePath}, nil
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

func runIndex(ctx context.Context, args indexArgs, stdout io.Writer) error {
	repo, cleanup, err := openRepository(ctx)
	if err != nil {
		return err
	}
	defer cleanup()

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to determine hostname: %w", err)
	}

	remoteName := normalizeRemoteName(args.Remote)
	rem := remote.NewRemote(hostname, remoteName)
	scn := operation.NewScanner(
		operation.WithRemote(&rem),
		operation.WithRepo(repo),
	)

	lo := rclone.NewListOption(rclone.ListOptionFilesOnly())
	lc := rclone.NewListConfig(toRcloneFS(rem.Name), args.Path, &lo)
	if err := operation.Scan(ctx, &scn, lc); err != nil {
		return err
	}

	_, err = fmt.Fprintf(stdout, "indexed remote %s from host %s\n", rem.Name, rem.Hostname)
	if err != nil {
		return fmt.Errorf("failed to write index result: %w", err)
	}

	return nil
}

func normalizeRemoteName(name string) string {
	return strings.TrimSuffix(strings.TrimSpace(name), ":")
}

func toRcloneFS(name string) string {
	trimmed := normalizeRemoteName(name)
	if trimmed == "" {
		return ""
	}
	return trimmed + ":"
}

func openRepository(ctx context.Context) (*repository.ObjectRepository, func(), error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, nil, fmt.Errorf("DATABASE_URL is required")
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing database url: %w", err)
	}
	config.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		log.Printf("Connected to database with pid %d", conn.PgConn().PID())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect: %w", err)
	}

	cleanup := func() {
		pool.Close()
	}

	return repository.NewObjectRepository(database.NewPgxDatabase(pool)), cleanup, nil
}
