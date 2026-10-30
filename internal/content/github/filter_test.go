package github

import (
	"errors"
	"testing"
)

func TestSelect(t *testing.T) {
	home := "https://example.com"
	tests := []struct {
		name      string
		repos     []Candidate
		readmes   map[string]string
		readmeErr map[string]error
		want      []string
	}{
		{
			name: "public repos sorted by stars and capped",
			repos: []Candidate{
				pub("high", 10),
				pub("low", 1),
				pub("mid", 3),
				pub("a", 4),
				pub("b", 5),
				pub("c", 6),
				pub("d", 7),
			},
			want: []string{"low", "mid", "a", "b", "c", "d"},
		},
		{
			name: "skips blocked names",
			repos: []Candidate{
				pub("my-portfolio", 1),
				pub("Zoomish-site", 2),
				pub("mesto-ad-shop", 3),
				pub("alko-dubai-ecom-api", 4),
				pub("bad-demo", 5),
				pub("kept", 6),
			},
			want: []string{"kept"},
		},
		{
			name: "private readme must mention a project",
			repos: []Candidate{
				priv("hidden", "owner/hidden", 1),
				priv("shown", "owner/shown", 2),
				priv("ru", "owner/ru", 3),
			},
			readmes: map[string]string{
				"owner/hidden": "just a note",
				"owner/shown":  "This Project ships",
				"owner/ru":     "Учебный проект",
			},
			want: []string{"shown", "ru"},
		},
		{
			name: "empty readme skips private repo",
			repos: []Candidate{
				priv("empty", "owner/empty", 1),
				pub("public", 2),
			},
			readmes: map[string]string{},
			want:    []string{"public"},
		},
		{
			name: "readme error skips that repo",
			repos: []Candidate{
				priv("broken", "owner/broken", 1),
				pub("public", 2),
			},
			readmeErr: map[string]error{"owner/broken": errors.New("boom")},
			want:      []string{"public"},
		},
		{
			name: "homepage stays nil",
			repos: []Candidate{
				pub("plain", 1),
			},
			want: []string{"plain"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "homepage stays nil" {
				tt.repos[0].Homepage = nil
			}
			if tt.name == "public repos sorted by stars and capped" {
				tt.repos[0].Homepage = &home
			}
			got, err := Select(tt.repos, func(fullName string) (string, error) {
				if err := tt.readmeErr[fullName]; err != nil {
					return "", err
				}
				return tt.readmes[fullName], nil
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("len %d, want %d (%v)", len(got), len(tt.want), names(got))
			}
			for i, name := range tt.want {
				if got[i].Name != name {
					t.Fatalf("index %d got %s", i, got[i].Name)
				}
			}
			if tt.name == "homepage stays nil" && got[0].Link != nil {
				t.Fatal("expected null link")
			}
		})
	}
}

func pub(name string, stars int) Candidate {
	return Candidate{Name: name, FullName: "owner/" + name, HTMLURL: "https://github.com/owner/" + name, Stars: stars}
}

func priv(name, full string, stars int) Candidate {
	repo := pub(name, stars)
	repo.FullName = full
	repo.Private = true
	return repo
}

func names(repos []Repo) []string {
	out := make([]string, 0, len(repos))
	for _, repo := range repos {
		out = append(out, repo.Name)
	}
	return out
}
