package card

import (
	"encoding/json"
	"strings"

	"github.com/Zoomish/portfolio-backend/internal/domain"
)

type profile struct {
	Fullname   string
	Headline   string
	About      string
	ProfileURL string
	Location   string
}

func Home(raw json.RawMessage) (View, error) {
	item, err := parseProfile(raw)
	if err != nil {
		return View{}, err
	}
	var b strings.Builder
	b.WriteString("Здравствуйте. Я <b>")
	b.WriteString(escape(item.Fullname))
	b.WriteString("</b>.\n")
	if item.Headline != "" {
		b.WriteString(escape(item.Headline))
		b.WriteString("\n")
	}
	b.WriteString("\nВыберите раздел.")
	return View{Text: b.String(), Keys: HomeKeys()}, nil
}

func About(raw json.RawMessage) (View, error) {
	item, err := parseProfile(raw)
	if err != nil {
		return View{}, err
	}
	var b strings.Builder
	b.WriteString("<b>")
	b.WriteString(escape(item.Fullname))
	b.WriteString("</b>\n")
	if item.Headline != "" {
		b.WriteString(escape(item.Headline))
		b.WriteString("\n")
	}
	if item.About != "" {
		b.WriteString("\n")
		b.WriteString(escape(clip(item.About, 3000)))
		b.WriteString("\n")
	}
	if item.Location != "" {
		b.WriteString("\n")
		b.WriteString(escape(item.Location))
		b.WriteString("\n")
	}
	if item.ProfileURL != "" {
		b.WriteString("\n")
		b.WriteString(link("Профиль", item.ProfileURL))
	}
	return View{Text: b.String(), Keys: sectionKeys(ScreenAbout), Event: domain.EventAbout}, nil
}

func parseProfile(raw json.RawMessage) (profile, error) {
	var doc struct {
		Data struct {
			BasicInfo json.RawMessage `json:"basic_info"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return profile{}, err
	}
	if len(doc.Data.BasicInfo) == 0 || string(doc.Data.BasicInfo) == "null" {
		return profile{}, ErrIncomplete
	}
	var info struct {
		Fullname   string `json:"fullname"`
		Headline   string `json:"headline"`
		About      string `json:"about"`
		ProfileURL string `json:"profile_url"`
		Location   struct {
			Full string `json:"full"`
			City string `json:"city"`
		} `json:"location"`
	}
	if err := json.Unmarshal(doc.Data.BasicInfo, &info); err != nil {
		return profile{}, err
	}
	if strings.TrimSpace(info.Fullname) == "" && strings.TrimSpace(info.Headline) == "" {
		return profile{}, ErrIncomplete
	}
	location := info.Location.Full
	if location == "" {
		location = info.Location.City
	}
	return profile{
		Fullname:   info.Fullname,
		Headline:   info.Headline,
		About:      info.About,
		ProfileURL: info.ProfileURL,
		Location:   location,
	}, nil
}
