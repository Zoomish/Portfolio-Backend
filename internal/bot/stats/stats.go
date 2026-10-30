package stats

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Zoomish/portfolio-backend/internal/domain"
	"github.com/Zoomish/portfolio-backend/internal/telegram"
)

const RecentLimit = 20

type Summary struct {
	Total    int
	Week     int
	Start    int
	About    int
	Projects int
	Contact  int
}

type Users interface {
	Count(ctx context.Context) (int, error)
	CountSince(ctx context.Context, since time.Time) (int, error)
	Recent(ctx context.Context, limit int) ([]domain.User, error)
}

type Events interface {
	Counts(ctx context.Context) (map[string]int, error)
}

type Query interface {
	Summary(ctx context.Context, now time.Time) (Summary, error)
	Recent(ctx context.Context, limit int) ([]domain.User, error)
}

type query struct {
	users  Users
	events Events
}

func NewQuery(users Users, events Events) Query {
	return query{users: users, events: events}
}

func (q query) Summary(ctx context.Context, now time.Time) (Summary, error) {
	total, err := q.users.Count(ctx)
	if err != nil {
		return Summary{}, fmt.Errorf("count users: %w", err)
	}
	week, err := q.users.CountSince(ctx, now.Add(-7*24*time.Hour))
	if err != nil {
		return Summary{}, fmt.Errorf("count recent users: %w", err)
	}
	counts, err := q.events.Counts(ctx)
	if err != nil {
		return Summary{}, fmt.Errorf("count events: %w", err)
	}
	return Summary{
		Total:    total,
		Week:     week,
		Start:    counts[domain.EventStart],
		About:    counts[domain.EventAbout],
		Projects: counts[domain.EventProjects],
		Contact:  counts[domain.EventContact],
	}, nil
}

func (q query) Recent(ctx context.Context, limit int) ([]domain.User, error) {
	users, err := q.users.Recent(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("recent users: %w", err)
	}
	return users, nil
}

func Allowed(actor, admin int64) bool {
	return admin != 0 && actor == admin
}

func FormatSummary(s Summary) string {
	return fmt.Sprintf("<b>Статистика</b>\n\nЛюдей: %d\nЗа 7 дней: %d\n\n/start: %d\nОбо мне: %d\nПроекты: %d\nКонтакт: %d",
		s.Total, s.Week, s.Start, s.About, s.Projects, s.Contact)
}

func FormatRecent(users []domain.User) string {
	if len(users) == 0 {
		return "<b>Последние</b>\n\nПока никого нет."
	}
	var b strings.Builder
	b.WriteString("<b>Последние</b>")
	for _, user := range users {
		b.WriteString("\n\n")
		b.WriteString(escape(user.Name))
		if user.Username != "" {
			b.WriteString(" · @")
			b.WriteString(escape(user.Username))
		}
		b.WriteString("\n<code>")
		b.WriteString(fmt.Sprintf("%d", user.TelegramID))
		b.WriteString("</code>\n")
		b.WriteString(user.UpdatedAt.UTC().Format("02.01.2006 15:04"))
		b.WriteString(" UTC")
	}
	return b.String()
}

func SummaryKeys() telegram.Keyboard {
	return telegram.Keyboard{{{Text: "Последние 20", Data: "stats-recent"}}}
}

func RecentKeys() telegram.Keyboard {
	return telegram.Keyboard{{{Text: "К сводке", Data: "stats-summary"}}}
}

func escape(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(s)
}
