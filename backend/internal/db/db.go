package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/paulakelo/finetic/backend/internal/parser"
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

func InsertTransaction(ctx context.Context, db *sql.DB, userID string, tx *parser.MpesaTransaction, rawSMS string) error {
	query := `
		INSERT INTO transactions 
		(user_id, amount, type, mpesa_receipt_number, sender_or_recipient, raw_sms_data, transaction_date)
		VALUES (KES1, KES2, KES3, KES4, KES5, KES6, KES7)
	`

	// ExecContext executes the query without returning any rows.
	_, err := db.ExecContext(
		ctx,
		query,
		userID,
		tx.Amount,
		tx.Type,
		tx.Counterparty,
		rawSMS,
		tx.Date,
	)

	return err
}
