package snapshot

import (
	"context"
	"errors"
	"fmt"

	"github.com/Zoomish/portfolio-backend/internal/content"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) Get(ctx context.Context, source content.Source) (content.Snapshot, error) {
	var snap content.Snapshot
	err := s.pool.QueryRow(ctx, `
		SELECT payload, fetched_at
		FROM content_snapshot
		WHERE source = $1`, string(source)).Scan(&snap.Payload, &snap.FetchedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return content.Snapshot{}, fmt.Errorf("get snapshot %s: %w", source, content.ErrNotFound)
	}
	if err != nil {
		return content.Snapshot{}, fmt.Errorf("get snapshot %s: %w", source, err)
	}
	snap.Source = source
	return snap, nil
}

func (s *Store) Save(ctx context.Context, snap content.Snapshot) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO content_snapshot (source, payload, fetched_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (source) DO UPDATE
		SET payload = EXCLUDED.payload,
		    fetched_at = EXCLUDED.fetched_at,
		    updated_at = now()`,
		string(snap.Source), snap.Payload, snap.FetchedAt)
	if err != nil {
		return fmt.Errorf("save snapshot %s: %w", snap.Source, err)
	}
	return nil
}
