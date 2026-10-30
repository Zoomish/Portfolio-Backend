package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Zoomish/portfolio-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) Upsert(ctx context.Context, telegramID int64, name, username string, admin bool) (domain.User, error) {
	if name == "" {
		name = "—"
	}
	var uname *string
	if username != "" {
		uname = &username
	}
	row := s.pool.QueryRow(ctx, `
		INSERT INTO users (telegram_id, admin, name, username)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (telegram_id) DO UPDATE
		SET admin = EXCLUDED.admin,
		    name = EXCLUDED.name,
		    username = EXCLUDED.username,
		    updated_at = now()
		RETURNING id, telegram_id, admin, name, username, created_at, updated_at`,
		telegramID, admin, name, uname)
	user, err := scanUser(row)
	if err != nil {
		return domain.User{}, fmt.Errorf("upsert user: %w", err)
	}
	return user, nil
}

func (s *Store) Count(ctx context.Context) (int, error) {
	var n int
	if err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&n); err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}
	return n, nil
}

func (s *Store) CountSince(ctx context.Context, since time.Time) (int, error) {
	var n int
	if err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE updated_at >= $1`, since).Scan(&n); err != nil {
		return 0, fmt.Errorf("count users since: %w", err)
	}
	return n, nil
}

func (s *Store) Recent(ctx context.Context, limit int) ([]domain.User, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, telegram_id, admin, name, username, created_at, updated_at
		FROM users
		ORDER BY updated_at DESC
		LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("recent users: %w", err)
	}
	defer rows.Close()

	var out []domain.User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("recent users: %w", err)
		}
		out = append(out, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recent users: %w", err)
	}
	return out, nil
}

func scanUser(row pgx.Row) (domain.User, error) {
	var user domain.User
	var username *string
	err := row.Scan(&user.ID, &user.TelegramID, &user.Admin, &user.Name, &username, &user.CreatedAt, &user.UpdatedAt)
	if username != nil {
		user.Username = *username
	}
	return user, err
}
