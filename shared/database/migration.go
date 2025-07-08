package database

import (
	"context"
	"fmt"

	"github.com/ZiplEix/mail-toolchain/shared/logger"
)

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

func MigrateUsersTable() error {
	sql := `
	CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := Pool.Exec(context.Background(), sql)
	if err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}
	logger.Info("Users table migrated successfully")
	return nil
}
