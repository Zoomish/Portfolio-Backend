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
	if home.Event != "" {
		t.Fatal("home records no screen event")
	}
	about, err := About(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(about.Text, "Builds &amp; ships") || !strings.Contains(about.Text, "Archangel, Russia") {
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
	if !strings.Contains(view.Text, "<b>api</b>") || !strings.Contains(view.Text, "Go · ★ 2") {
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

func TestIncompleteProfile(t *testing.T) {
	if _, err := Home([]byte(`{"data":{}}`)); err != ErrIncomplete {
		t.Fatal(err)
	}
}
