package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ZiplEix/mail-toolchain/shared/logger"
	"github.com/jackc/pgx/v5"
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
	logger.Info("Connected to PostgreSQL")
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
	if err != nil {
		return fmt.Errorf("failed to create mails table: %v", err)
	}
	logger.Info("Mails table migrated successfully")
	return nil
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

func GetAllMails() ([]Mail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := Pool.Query(ctx, `
		SELECT id, sender, recipients, raw_data, received_at
		FROM mails
		ORDER BY received_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query mails: %v", err)
	}
	defer rows.Close()

	var mails []Mail
	for rows.Next() {
		var m Mail
		err := rows.Scan(&m.ID, &m.Sender, &m.Recipients, &m.RawData, &m.ReceivedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan mail: %v", err)
		}
		mails = append(mails, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return mails, nil
}

func GetMailsInUIDRange(start, end int) ([]Mail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string
	var rows pgx.Rows
	var err error

	if end == -1 {
		query = `SELECT id, sender, recipients, raw_data, received_at FROM mails WHERE id >= $1 ORDER BY id ASC`
		rows, err = Pool.Query(ctx, query, start)
	} else {
		query = `SELECT id, sender, recipients, raw_data, received_at FROM mails WHERE id BETWEEN $1 AND $2 ORDER BY id ASC`
		rows, err = Pool.Query(ctx, query, start, end)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query mails in UID range: %v", err)
	}
	defer rows.Close()

	var mails []Mail
	for rows.Next() {
		var m Mail
		err := rows.Scan(&m.ID, &m.Sender, &m.Recipients, &m.RawData, &m.ReceivedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan mail: %v", err)
		}
		mails = append(mails, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return mails, nil
}
