package database

import (
	"context"
	"fmt"
	"time"

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
		for i := range 5 {
			logger.Warnf("Failed to connect to PostgreSQL, retrying in 2 seconds... (%d/5)", i+1)
			if err := Pool.Ping(context.Background()); err == nil {
				logger.Info("Connected to PostgreSQL after retry")
				return nil
			}
			time.Sleep(2 * time.Second)
		}
	}
	logger.Info("Connected to PostgreSQL")
	return nil
}
