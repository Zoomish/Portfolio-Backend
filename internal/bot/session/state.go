package session

import "context"

const WaitErrorReport = "error_report"

type State struct {
	UserID           int64
	WaitingFor       string
	LastBotMessageID int64
	ErrorMessageID   int64
}

func (s State) Waiting() bool {
	return s.WaitingFor != ""
}

func Hint(waiting string) string {
	switch waiting {
	case WaitErrorReport:
		return "Жду описание проблемы одним сообщением."
	default:
		return "Жду сообщение."
	}
}

type Store interface {
	Load(ctx context.Context, userID int64) (State, error)
	Save(ctx context.Context, state State) error
}
