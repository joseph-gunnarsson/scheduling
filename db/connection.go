package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func GetDBConnection() *pgx.Conn {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		log.Fatalf("Failed to ping database: %v\n", err)
		conn.Close(ctx)
	}

	return conn
}
