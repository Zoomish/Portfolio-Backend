package callback

import "testing"

func TestSplit(t *testing.T) {
	tests := []struct {
		in, head, rest string
	}{
		{"menu-about", "menu", "about"},
		{"menu-projects", "menu", "projects"},
		{"menu-contact", "menu", "contact"},
		{"menu-posts", "menu", "posts"},
		{"menu-home", "menu", "home"},
		{"stats-recent", "stats", "recent"},
		{"stats-summary", "stats", "summary"},
	}
	for _, tt := range tests {
		head, rest := Split(tt.in)
		if head != tt.head || rest != tt.rest {
			t.Fatalf("%s -> %s %s", tt.in, head, rest)
		}
	}
}
