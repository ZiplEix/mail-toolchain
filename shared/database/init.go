package database

import (
	"context"
	"fmt"

	"github.com/ZiplEix/mail-toolchain/shared/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
