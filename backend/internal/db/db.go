package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// OpenDB establishes a connection pool to PostgreSQL and verifies it with a ping.
func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Configure the connection pool settings for high concurrency
	db.SetMaxOpenConns(25)                  // Max concurrent connections to Postgres
	db.SetMaxIdleConns(25)                  // Max idle connections to keep alive
	db.SetConnMaxIdleTime(15 * time.Minute) // How long a connection can sit idle before closing

	// Create a context with a 5-second timeout to ensure the ping doesn't hang forever
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// PingContext actually establishes the connection to verify credentials and reachability
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
