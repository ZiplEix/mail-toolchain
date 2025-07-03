package database

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

func Init(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	var err error
	once.Do(func() {
		if pool, err = pgxpool.New(ctx, dsn); err != nil {
			return
		}
		err = pool.Ping(ctx)
	})
	return pool, err
}

func Pool() *pgxpool.Pool {
	return pool
}
