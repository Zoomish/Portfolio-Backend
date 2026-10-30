package content

import (
	"encoding/json"
	"errors"
	"time"
)

type Source string

const (
	SourceGithubRepos     Source = "github_repos"
	SourceLinkedInPosts   Source = "linkedin_posts"
	SourceLinkedInProfile Source = "linkedin_profile"
	SourceLinkedInContact Source = "linkedin_contact"
)

var (
	ErrNotFound = errors.New("snapshot not found")
	ErrNotReady = errors.New("snapshot not ready")
)

type Snapshot struct {
	Source    Source
	Payload   json.RawMessage
	FetchedAt time.Time
}
