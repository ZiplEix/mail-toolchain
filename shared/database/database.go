package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Init(dsn string) error {
	var err error
	Pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	if err := Pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("unable to ping database: %v", err)
	}
	fmt.Println("✅ Connected to PostgreSQL")
	return nil
}

func MigrateMailsTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS mails (
		id SERIAL PRIMARY KEY,
		sender VARCHAR(255) NOT NULL,
		recipients TEXT[] NOT NULL,
		raw_data TEXT NOT NULL,
		received_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`
	_, err := Pool.Exec(context.Background(), sql)
	return err
}

func SaveMail(sender string, recipients []string, rawData []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	raw := strings.Join(rawData, "\r\n")

	_, err := Pool.Exec(ctx,
		"INSERT INTO mails(sender, recipients, raw_data, received_at) VALUES($1, $2, $3, NOW())",
		sender, recipients, raw)

	return err
}
