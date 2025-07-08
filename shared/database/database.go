package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

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
