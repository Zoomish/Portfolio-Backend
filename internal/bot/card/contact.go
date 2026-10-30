package card

import (
	"encoding/json"
	"strings"

	"github.com/Zoomish/portfolio-backend/internal/domain"
)

func Contact(raw json.RawMessage) (View, error) {
	var doc struct {
		Data struct {
			Email     *string         `json:"email"`
			BasicInfo json.RawMessage `json:"basic_info"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return View{}, err
	}
	var info struct {
		FullName   string `json:"full_name"`
		ProfileURL string `json:"profile_url"`
	}
	if len(doc.Data.BasicInfo) > 0 && string(doc.Data.BasicInfo) != "null" {
		if err := json.Unmarshal(doc.Data.BasicInfo, &info); err != nil {
			return View{}, err
		}
	}
	email := ""
	if doc.Data.Email != nil {
		email = strings.TrimSpace(*doc.Data.Email)
	}
	if info.FullName == "" && email == "" && info.ProfileURL == "" {
		return View{}, ErrIncomplete
	}
	var b strings.Builder
	b.WriteString("<b>Контакт</b>\n")
	if info.FullName != "" {
		b.WriteString(escape(info.FullName))
		b.WriteString("\n")
	}
	if email != "" {
		b.WriteString(escape(email))
	} else {
		b.WriteString("Почта не указана")
	}
	if info.ProfileURL != "" {
		b.WriteString("\n")
		b.WriteString(link("Профиль", info.ProfileURL))
	}
	return View{Text: b.String(), Keys: sectionKeys(ScreenContact), Event: domain.EventContact}, nil
}
