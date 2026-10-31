package card

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/Zoomish/portfolio-backend/internal/domain"
)

const maxPosts = 5

type linkedInPost struct {
	Text     string `json:"text"`
	URL      string `json:"url"`
	PostedAt struct {
		Date      string `json:"date"`
		Relative  string `json:"relative"`
		Timestamp int64  `json:"timestamp"`
	} `json:"posted_at"`
}

func Posts(raw json.RawMessage) (View, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return View{}, ErrIncomplete
	}
	var doc struct {
		Data struct {
			Posts []linkedInPost `json:"posts"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return View{}, err
	}
	posts := doc.Data.Posts
	sort.SliceStable(posts, func(i, j int) bool {
		return posts[i].PostedAt.Timestamp > posts[j].PostedAt.Timestamp
	})
	if len(posts) > maxPosts {
		posts = posts[:maxPosts]
	}
	var b strings.Builder
	b.WriteString("<b>📝 Посты</b>")
	if len(posts) == 0 {
		b.WriteString("\n\nПостов пока нет.")
		return View{Text: b.String(), Keys: sectionKeys(ScreenPosts), Event: domain.EventPosts}, nil
	}
	for _, post := range posts {
		b.WriteString("\n\n")
		text := strings.TrimSpace(post.Text)
		if text == "" {
			text = "Без текста"
		}
		b.WriteString(escape(clip(text, 400)))
		if when := postWhen(post); when != "" {
			b.WriteString("\n")
			b.WriteString(escape(when))
		}
		if post.URL != "" {
			b.WriteString("\n")
			b.WriteString(link("Открыть", post.URL))
		}
	}
	return View{Text: b.String(), Keys: sectionKeys(ScreenPosts), Event: domain.EventPosts}, nil
}

func postWhen(post linkedInPost) string {
	if strings.TrimSpace(post.PostedAt.Relative) != "" {
		return post.PostedAt.Relative
	}
	return strings.TrimSpace(post.PostedAt.Date)
}
