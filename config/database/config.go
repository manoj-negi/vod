package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	DB *pgxpool.Pool
}

// func NewPostgres(ctx context.Context, connString string)
func NewPostgres(ctx context.Context, connString string) (*Postgres, error) {
	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	fmt.Println("--db connected successfully-")
	return &Postgres{DB: db}, nil
}

func (p *Postgres) Ping(ctx context.Context) error {
	return p.DB.Ping(ctx)
}

func (p *Postgres) Close() error {
	defer p.DB.Close()
	return nil
}
