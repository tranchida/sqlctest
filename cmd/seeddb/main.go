package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	ctx := context.Background()

	var url string
	if url = os.Getenv("POSTGRESQL_URL"); url == "" {
		url = "postgres://user:password@localhost:5432/sqlctest?sslmode=disable"
	}

	conn, err := pgxpool.New(ctx, url)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	schema, err := os.ReadFile("build/sqlc/schema.sql")
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(ctx, string(schema))
	if err != nil {
		panic(err)
	}

}
