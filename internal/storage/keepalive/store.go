package keepalive

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const rewriteEvery = 20

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) Touch(ctx context.Context) error {
	var seq int64
	err := s.pool.QueryRow(ctx, `
		UPDATE keepalive
		SET seq = seq + 1, touched_at = now()
		WHERE id = 1
		RETURNING seq`).Scan(&seq)
	if errors.Is(err, pgx.ErrNoRows) {
		if _, err = s.pool.Exec(ctx, `INSERT INTO keepalive (id, seq, touched_at) VALUES (1, 1, now())`); err != nil {
			return fmt.Errorf("keepalive insert: %w", err)
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("keepalive update: %w", err)
	}
	if seq%rewriteEvery != 0 {
		return nil
	}
	return s.rewrite(ctx, seq)
}

func (s *Store) rewrite(ctx context.Context, seq int64) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("keepalive tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx, `DELETE FROM keepalive WHERE id = 1`); err != nil {
		return fmt.Errorf("keepalive delete: %w", err)
	}
	if _, err = tx.Exec(ctx, `INSERT INTO keepalive (id, touched_at, seq) VALUES (1, now(), $1)`, seq); err != nil {
		return fmt.Errorf("keepalive reinsert: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("keepalive commit: %w", err)
	}
	return nil
}

func (s *Store) Ready(ctx context.Context) error {
	var one int
	if err := s.pool.QueryRow(ctx, `SELECT 1 FROM keepalive WHERE id = 1`).Scan(&one); err != nil {
		return fmt.Errorf("ready: %w", err)
	}
	return nil
}
