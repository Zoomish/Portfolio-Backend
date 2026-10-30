package github

type Candidate struct {
	Name        string
	FullName    string
	Description string
	Language    string
	HTMLURL     string
	Homepage    *string
	Stars       int
	Forks       int
	Private     bool
}

type Repo struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Language    string  `json:"language"`
	URL         string  `json:"url"`
	Link        *string `json:"link"`
	Stars       int     `json:"stars"`
	Forks       int     `json:"forks"`
}
