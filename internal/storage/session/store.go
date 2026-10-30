package session

import (
	"context"
	"fmt"

	botsession "github.com/Zoomish/portfolio-backend/internal/bot/session"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) Load(ctx context.Context, userID int64) (botsession.State, error) {
	var state botsession.State
	var last, errID *int64
	err := s.pool.QueryRow(ctx, `
		INSERT INTO telegram_session (user_id)
		VALUES ($1)
		ON CONFLICT (user_id) DO UPDATE
		SET user_id = EXCLUDED.user_id
		RETURNING waiting_for, last_bot_message_id, error_message_id`,
		userID).Scan(&state.WaitingFor, &last, &errID)
	if err != nil {
		return botsession.State{}, fmt.Errorf("load session: %w", err)
	}
	state.UserID = userID
	if last != nil {
		state.LastBotMessageID = *last
	}
	if errID != nil {
		state.ErrorMessageID = *errID
	}
	return state, nil
}

func (s *Store) Save(ctx context.Context, state botsession.State) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE telegram_session
		SET waiting_for = $2,
		    last_bot_message_id = $3,
		    error_message_id = $4
		WHERE user_id = $1`,
		state.UserID, state.WaitingFor, nullID(state.LastBotMessageID), nullID(state.ErrorMessageID))
	if err != nil {
		return fmt.Errorf("save session: %w", err)
	}
	return nil
}

func nullID(id int64) any {
	if id == 0 {
		return nil
	}
	return id
}
