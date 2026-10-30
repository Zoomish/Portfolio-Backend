package event

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) Add(ctx context.Context, userID int64, kind string) error {
	if userID == 0 || kind == "" {
		return fmt.Errorf("add event: user and kind are required")
	}
	_, err := s.pool.Exec(ctx, `
		INSERT INTO user_event (user_id, kind)
		VALUES ($1, $2)`, userID, kind)
	if err != nil {
		return fmt.Errorf("add event: %w", err)
	}
	return nil
}

func (s *Store) Counts(ctx context.Context) (map[string]int, error) {
	rows, err := s.pool.Query(ctx, `SELECT kind, COUNT(*) FROM user_event GROUP BY kind`)
	if err != nil {
		return nil, fmt.Errorf("count events: %w", err)
	}
	defer rows.Close()

	out := map[string]int{}
	for rows.Next() {
		var kind string
		var n int
		if err := rows.Scan(&kind, &n); err != nil {
			return nil, fmt.Errorf("count events: %w", err)
		}
		out[kind] = n
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("count events: %w", err)
	}
	return out, nil
}
