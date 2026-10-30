package domain

import "time"

const (
	EventStart    = "start"
	EventAbout    = "about"
	EventProjects = "projects"
	EventContact  = "contact"
)

type User struct {
	ID         int64
	TelegramID int64
	Admin      bool
	Name       string
	Username   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
