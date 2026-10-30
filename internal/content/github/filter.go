package github

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"
)

const maxRepos = 6

var blockedNameParts = []string{
	"portfolio",
	"zoomish",
	"mesto-ad",
	"alko-dubai-ecom",
	"bad",
}

func Select(repos []Candidate, readmeOf func(fullName string) (string, error)) ([]Repo, error) {
	sorted := append([]Candidate(nil), repos...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Stars < sorted[j].Stars
	})

	out := make([]Repo, 0, maxRepos)
	for _, repo := range sorted {
		if len(out) >= maxRepos {
			break
		}
		if blockedName(repo.Name) {
			continue
		}
		if !repo.Private {
			out = append(out, toRepo(repo))
			continue
		}
		if readmeOf == nil {
			return nil, fmt.Errorf("readme lookup is nil")
		}
		readme, err := readmeOf(repo.FullName)
		if err != nil {
			slog.Warn("skip repo readme", "repo", repo.FullName, "err", err)
			continue
		}
		text := strings.ToLower(readme)
		if strings.Contains(text, "project") || strings.Contains(text, "проект") {
			out = append(out, toRepo(repo))
		}
	}
	return out, nil
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
