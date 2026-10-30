package github

import (
	"sort"
	"strings"
)

const maxRepos = 4

var blockedNameParts = []string{
	"portfolio",
	"zoomish",
	"mesto-ad",
	"alko-dubai-ecom",
	"bad",
}

func Select(repos []Candidate) []Repo {
	public := make([]Candidate, 0, len(repos))
	for _, repo := range repos {
		if repo.Private {
			continue
		}
		public = append(public, repo)
	}
	sort.SliceStable(public, func(i, j int) bool {
		return public[i].Stars < public[j].Stars
	})

	out := make([]Repo, 0, maxRepos)
	for _, repo := range public {
		if len(out) >= maxRepos {
			break
		}
		if blockedName(repo.Name) {
			continue
		}
		out = append(out, toRepo(repo))
	}
	return out
}

func blockedName(name string) bool {
	lower := strings.ToLower(name)
	for _, part := range blockedNameParts {
		if strings.Contains(lower, part) {
			return true
		}
	}
	return false
}

func toRepo(repo Candidate) Repo {
	return Repo{
		Name:        repo.Name,
		Description: repo.Description,
		Language:    repo.Language,
		URL:         repo.HTMLURL,
		Link:        repo.Homepage,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
	}
}
