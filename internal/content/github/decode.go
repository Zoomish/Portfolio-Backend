package github

import "encoding/json"

type apiRepo struct {
	Name        string  `json:"name"`
	FullName    string  `json:"full_name"`
	Private     bool    `json:"private"`
	HTMLURL     string  `json:"html_url"`
	Homepage    *string `json:"homepage"`
	Description *string `json:"description"`
	Language    *string `json:"language"`
	Stars       int     `json:"stargazers_count"`
	Forks       int     `json:"forks_count"`
}

func DecodeRepos(body []byte) ([]Candidate, error) {
	var raw []apiRepo
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	out := make([]Candidate, 0, len(raw))
	for _, repo := range raw {
		out = append(out, Candidate{
			Name:        repo.Name,
			FullName:    repo.FullName,
			Description: deref(repo.Description),
			Language:    deref(repo.Language),
			HTMLURL:     repo.HTMLURL,
			Homepage:    repo.Homepage,
			Stars:       repo.Stars,
			Forks:       repo.Forks,
			Private:     repo.Private,
		})
	}
	return out, nil
}

func deref(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
