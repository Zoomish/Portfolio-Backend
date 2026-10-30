package github

import "testing"

func TestSelect(t *testing.T) {
	home := "https://example.com"
	tests := []struct {
		name  string
		repos []Candidate
		want  []string
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
			want: []string{"low", "mid", "a", "b"},
		},
		{
			name: "private is excluded before the cap",
			repos: []Candidate{
				priv("secret", "owner/secret", 0),
				pub("a", 1),
				pub("b", 2),
				pub("c", 3),
				pub("d", 4),
				pub("e", 5),
			},
			want: []string{"a", "b", "c", "d"},
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
			got := Select(tt.repos)
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
