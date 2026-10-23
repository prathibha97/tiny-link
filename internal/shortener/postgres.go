package shortener

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresPool(ctx context.Context, dsn string) (*PostgresStore, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	return &PostgresStore{
		pool: pool,
	}, nil
}

func (s *PostgresStore) Save(
	ctx context.Context,
	mapping URLMapping,
) error {
	_, err := s.pool.Exec(ctx,
		`
	INSERT INTO urls (code, original_url)
	VALUES ($1, $2)
	`,

		mapping.Code,
		mapping.OriginalURL,
	)
	if err != nil {
		return fmt.Errorf("save url mapping: %w", err)
	}
	return nil
}

func (s *PostgresStore) Get(ctx context.Context, code string) (URLMapping, error) {
	var mapping URLMapping
	err := s.pool.QueryRow(ctx,
		`SELECT code, original_url
	FROM urls
	WHERE code = $1
	`,
		code,
	).Scan(&mapping.Code, &mapping.OriginalURL)

	if errors.Is(err, pgx.ErrNoRows) {
		return URLMapping{}, ErrNotFound
	}

	if err != nil {
		return URLMapping{}, fmt.Errorf("get URL mapping: %w", err)
	}
	return mapping, nil
}

func (s *PostgresStore) Close() {
	s.pool.Close()
}
