package turn

import (
	"github.com/Zoomish/portfolio-backend/internal/bot/session"
	"github.com/Zoomish/portfolio-backend/internal/domain"
)

type Turn struct {
	UpdateID   int64
	ChatID     int64
	Actor      domain.User
	Session    session.State
	IncomingID int64
	CallbackID string
}
