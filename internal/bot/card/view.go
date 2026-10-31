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
	ScreenPosts    = "posts"
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
	case ScreenHome, ScreenAbout, ScreenProjects, ScreenPosts, ScreenContact:
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
		{{Text: "👤 Обо мне", Data: "menu-about"}},
		{{Text: "💼 Проекты", Data: "menu-projects"}},
		{{Text: "📝 Посты", Data: "menu-posts"}},
		{{Text: "✉️ Контакт", Data: "menu-contact"}},
	}
}

func sectionKeys(current string) telegram.Keyboard {
	items := []struct {
		screen string
		text   string
		data   string
	}{
		{ScreenAbout, "👤 Обо мне", "menu-about"},
		{ScreenProjects, "💼 Проекты", "menu-projects"},
		{ScreenPosts, "📝 Посты", "menu-posts"},
		{ScreenContact, "✉️ Контакт", "menu-contact"},
	}
	keys := make(telegram.Keyboard, 0, len(items)+1)
	for _, item := range items {
		if item.screen == current {
			continue
		}
		keys = append(keys, []telegram.Button{{Text: item.text, Data: item.data}})
	}
	keys = append(keys, []telegram.Button{{Text: "🏠 В начало", Data: "menu-home"}})
	return keys
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
