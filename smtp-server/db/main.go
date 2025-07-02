package db

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func InitDB() error {
	dsn := os.Getenv("POSTGRES_URL")
	if dsn == "" {
		return fmt.Errorf("POSTGRES_URL not set")
	}

	var err error
	dbPool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	if err := dbPool.Ping(context.Background()); err != nil {
		return fmt.Errorf("unable to ping database: %v", err)
	}
	fmt.Println("Connected to database")

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS mails (
		id SERIAL PRIMARY KEY,
		sender VARCHAR(255) NOT NULL,
		recipients TEXT[] NOT NULL,
		raw_data TEXT NOT NULL,
		received_at TIMESTAMP NOT NULL DEFAULT NOW()
	);
	`

	_, err = dbPool.Exec(context.Background(), createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create mails table: %w", err)
	}

	return nil
}

func SaveMailToDB(sender string, recipients []string, rawData []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	raw := strings.Join(rawData, "\r\n")

	_, err := dbPool.Exec(ctx,
		"INSERT INTO mails(sender, recipients, raw_data, received_at) VALUES($1, $2, $3, NOW())",
		sender, recipients, raw)

	return err
}
