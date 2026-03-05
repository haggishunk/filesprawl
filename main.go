package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/haggishunk/filesprawl/internal/analysis"
	"github.com/haggishunk/filesprawl/internal/database"
	"github.com/haggishunk/filesprawl/internal/object"
	"github.com/haggishunk/filesprawl/internal/operation"
	"github.com/haggishunk/filesprawl/internal/rclone"
	"github.com/haggishunk/filesprawl/internal/remote"
	"github.com/haggishunk/filesprawl/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// out, err := list()

	// err := clients.ListJSON(ctx, "dbox:", "rollbar/macbook")
	// if err != nil {
	// 	fmt.Printf("Error: %q", err)
	// }

	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error parsing database url: %w", err)
	}
	config.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		log.Printf("Connected to database with pid %d", conn.PgConn().PID())
		return nil
	}

	// opening the pool and closing the pool should span the lifetime
	// of the main thread
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("Failed to connect: %w", err)
	}
	defer pool.Close()

	repo := repository.NewObjectRepository(database.NewPgxDatabase(pool))

	// BEGIN SAMPLE CODE
	// retrieve a hash from a db repo
	h := object.NewHash("jkdjfw", "md5")

	// demonstrates how a query result from rclone could be
	// referenced against and persisted in a database
	err = repo.ReadHash(context.Background(), &h)
	if err != nil {
		log.Printf("Error retrieving hash: %s\n", err)
	}
	if h.Persisted {
		log.Printf("Found persisted hash: %s\n", h)
	} else {
		log.Printf("Hash not found in db: %s", h.Hash)
	}
	// END SAMPLE CODE

	rem := remote.Remote{}
	rem.Hostname = "guru"
	rem.Name = "dbox:"

	scn := operation.NewScanner(
		operation.WithRemote(&rem),
		operation.WithRepo(repo),
	)
	// persist scan start

	lo := rclone.NewListOption(rclone.ListOptionFilesOnly())
	lc := rclone.NewListConfig(rem.Name, "code/flux", &lo)

	err = operation.Scan(context.Background(), &scn, lc)
	if err != nil {
		log.Panic("Failed %w", err)
	}

	// Duplicate detection across all remotes
	detector := analysis.NewDuplicateDetector(repo)
	groups, err := detector.FindAcrossRemotes(context.Background(), analysis.FilterOptions{})
	if err != nil {
		log.Printf("Error finding duplicates: %v", err)
	} else {
		fmt.Print(analysis.FormatReport(groups))
	}
}
