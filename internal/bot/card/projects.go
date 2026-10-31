package card

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Zoomish/portfolio-backend/internal/content/github"
	"github.com/Zoomish/portfolio-backend/internal/domain"
)

func Projects(raw json.RawMessage) (View, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return View{}, ErrIncomplete
	}
	var repos []github.Repo
	if err := json.Unmarshal(raw, &repos); err != nil {
		return View{}, err
	}
	var b strings.Builder
	b.WriteString("<b>💼 Проекты</b>")
	if len(repos) == 0 {
		b.WriteString("\n\nПубличных проектов пока нет.")
		return View{Text: b.String(), Keys: sectionKeys(ScreenProjects), Event: domain.EventProjects}, nil
	}
	for _, repo := range repos {
		b.WriteString("\n\n<b>")
		b.WriteString(escape(repo.Name))
		b.WriteString("</b>")
		if repo.Description != "" {
			b.WriteString("\n")
			b.WriteString(escape(clip(repo.Description, 300)))
		}
		b.WriteString("\n")
		b.WriteString(meta(repo))
		if repo.URL != "" {
			b.WriteString("\n")
			b.WriteString(link(repo.URL, repo.URL))
		}
	}
	return View{Text: b.String(), Keys: sectionKeys(ScreenProjects), Event: domain.EventProjects}, nil
}

func meta(repo github.Repo) string {
	parts := make([]string, 0, 2)
	if repo.Language != "" {
		parts = append(parts, escape(repo.Language))
	}
	parts = append(parts, fmt.Sprintf("★ %d", repo.Stars))
	return strings.Join(parts, " · ")
}
