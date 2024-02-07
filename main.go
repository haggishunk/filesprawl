package main

import (
	"context"
	"log"
	"os"

	"github.com/haggishunk/filesprawl/internal/database"
	"github.com/haggishunk/filesprawl/internal/object"
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
		log.Printf("Connected to database with pid %s", conn.PgConn().PID())
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

	// sample code to retrieve a hash from a db repo
	h := object.Hash{
		Hash: "872f92f3",
		Type: "md5",
	}

	// demonstrates how a query result from rclone could be
	// referenced against and persisted in a database
	id, err := repo.GetHash(context.Background(), h)
	if err != nil {
		log.Printf("Error retrieving hash: %s", err)
	}
	log.Println("ID: ", id)
}
