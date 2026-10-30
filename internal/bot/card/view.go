package card

import (
	"errors"
	"strings"

	"github.com/Zoomish/portfolio-backend/internal/telegram"
)

const (
	ScreenHome     = "home"
	ScreenAbout    = "about"
	ScreenProjects = "projects"
	ScreenContact  = "contact"
)

var ErrIncomplete = errors.New("snapshot incomplete")

type View struct {
	Text  string
	Keys  telegram.Keyboard
	Event string
}

func Kind(rest string) (string, bool) {
	switch rest {
	case ScreenHome, ScreenAbout, ScreenProjects, ScreenContact:
		return rest, true
	default:
		return "", false
	}
}

func Unavailable() View {
	return View{
		Text: "Данные ещё не готовы. Попробуй чуть позже.",
		Keys: HomeKeys(),
	}
}

func HomeKeys() telegram.Keyboard {
	return telegram.Keyboard{
		{{Text: "Обо мне", Data: "menu-about"}},
		{{Text: "Проекты", Data: "menu-projects"}},
		{{Text: "Контакт", Data: "menu-contact"}},
	}
}

func sectionKeys(current string) telegram.Keyboard {
	var row []telegram.Button
	if current != ScreenAbout {
		row = append(row, telegram.Button{Text: "Обо мне", Data: "menu-about"})
	}
	if current != ScreenProjects {
		row = append(row, telegram.Button{Text: "Проекты", Data: "menu-projects"})
	}
	if current != ScreenContact {
		row = append(row, telegram.Button{Text: "Контакт", Data: "menu-contact"})
	}
	return telegram.Keyboard{
		row,
		{{Text: "В начало", Data: "menu-home"}},
	}
}

func escape(s string) string {
	return strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
	).Replace(s)
}

func clip(s string, n int) string {
	runes := []rune(s)
	if n <= 0 || len(runes) <= n {
		return s
	}
	return string(runes[:n-1]) + "…"
}

func link(label, href string) string {
	if !strings.HasPrefix(href, "https://") && !strings.HasPrefix(href, "http://") {
		return escape(label)
	}
	return `<a href="` + escape(href) + `">` + escape(label) + `</a>`
}
