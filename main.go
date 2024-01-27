package main

import (
	"context"
	"fmt"

	"github.com/haggishunk/filesprawl/internal/clients"
)

func main() {
	// out, err := list()
	ctx := context.Background()
	err := clients.ListJSON(ctx, "dbox:", "rollbar/macbook")
	if err != nil {
		fmt.Printf("Error: %q", err)
	}
}
