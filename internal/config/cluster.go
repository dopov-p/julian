package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func DSN() string {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	pgDSN := os.Getenv("POSTGRES_PG_DSN")
	sslMode := os.Getenv("POSTGRES_SSL_MODE")

	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf("%s&sslmode=%s", pgDSN, sslMode)
}

type Cluster struct {
	Conn *pgxpool.Pool
}

func NewCluster(ctx context.Context) (*Cluster, error) {
	dsn := DSN()

	dbpool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}

	// Test the connection
	if err = dbpool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("dbpool.Ping: %w", err)
	}

	log.Printf("Successfully connected to database")

	return &Cluster{Conn: dbpool}, nil
}

func (c *Cluster) Close() {
	if c.Conn != nil {
		c.Conn.Close()
	}
}
