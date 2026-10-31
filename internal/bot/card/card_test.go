package card

import (
	"strings"
	"testing"
)

func TestHomeAndAboutFromProfile(t *testing.T) {
	raw := []byte(`{"data":{"basic_info":{
		"fullname":"Egor <Rezvyj>",
		"headline":"Go Backend",
		"about":"Builds & ships",
		"profile_url":"https://linkedin.com/in/zoomish",
		"location":{"full":"Archangel, Russia","city":"Archangel"}
	}}}`)
	home, err := Home(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(home.Text, "Egor &lt;Rezvyj&gt;") || !strings.Contains(home.Text, "Go Backend") {
		t.Fatalf("home %s", home.Text)
	}
	if !strings.Contains(home.Text, "<b>Кнопки:</b>") || !strings.Contains(home.Text, "📝 Посты —") {
		t.Fatalf("legend %s", home.Text)
	}
	if len(home.Keys) != 4 {
		t.Fatalf("keys %d", len(home.Keys))
	}
	for _, row := range home.Keys {
		if len(row) != 1 {
			t.Fatalf("row %v", row)
		}
	}
	wantButtons := []string{"👤 Обо мне", "💼 Проекты", "📝 Посты", "✉️ Контакт"}
	for i, label := range wantButtons {
		if home.Keys[i][0].Text != label {
			t.Fatalf("button %d %s", i, home.Keys[i][0].Text)
		}
	}
	if home.Event != "" {
		t.Fatal("home records no screen event")
	}
	about, err := About(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(about.Text, "👤 ") || !strings.Contains(about.Text, "Builds &amp; ships") || !strings.Contains(about.Text, "Archangel, Russia") {
		t.Fatalf("about %s", about.Text)
	}
	if !strings.Contains(about.Text, `href="https://linkedin.com/in/zoomish"`) {
		t.Fatalf("link %s", about.Text)
	}
}

func TestProjectsAndContact(t *testing.T) {
	repos := []byte(`[{"name":"api","description":"edge","language":"Go","url":"https://github.com/zoomish/api","stars":2}]`)
	view, err := Projects(repos)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(view.Text, "<b>💼 Проекты</b>") || !strings.Contains(view.Text, "<b>api</b>") || !strings.Contains(view.Text, "Go · ★ 2") {
		t.Fatalf("projects %s", view.Text)
	}
	contactRaw := []byte(`{"data":{"email":null,"basic_info":{"full_name":"Egor","profile_url":"https://www.linkedin.com/in/zoomish"}}}`)
	contact, err := Contact(contactRaw)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(contact.Text, "Egor") || !strings.Contains(contact.Text, "Почта не указана") {
		t.Fatalf("contact %s", contact.Text)
	}
}

func TestPostsKeepsFiveNewest(t *testing.T) {
	raw := []byte(`{"data":{"posts":[
		{"text":"oldest","url":"https://linkedin.com/p/1","posted_at":{"timestamp":1,"relative":"6d"}},
		{"text":"second","url":"https://linkedin.com/p/2","posted_at":{"timestamp":2,"relative":"5d"}},
		{"text":"third","url":"https://linkedin.com/p/3","posted_at":{"timestamp":3,"relative":"4d"}},
		{"text":"fourth","url":"https://linkedin.com/p/4","posted_at":{"timestamp":4,"relative":"3d"}},
		{"text":"fifth","url":"https://linkedin.com/p/5","posted_at":{"timestamp":5,"relative":"2d"}},
		{"text":"newest","url":"https://linkedin.com/p/6","posted_at":{"timestamp":6,"relative":"1d"}}
	]}}`)
	view, err := Posts(raw)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(view.Text, "oldest") {
		t.Fatalf("oldest kept: %s", view.Text)
	}
	if !strings.Contains(view.Text, "newest") || !strings.Contains(view.Text, "second") {
		t.Fatalf("posts %s", view.Text)
	}
	if !strings.Contains(view.Text, `href="https://linkedin.com/p/6"`) || !strings.Contains(view.Text, "1d") {
		t.Fatalf("link %s", view.Text)
	}
	if view.Event != "posts" {
		t.Fatalf("event %s", view.Event)
	}
}

func TestIncompleteProfile(t *testing.T) {
	if _, err := Home([]byte(`{"data":{}}`)); err != ErrIncomplete {
		t.Fatal(err)
	}
}
