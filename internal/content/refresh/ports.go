package refresh

import (
	"context"
	"encoding/json"

	"github.com/Zoomish/portfolio-backend/internal/content"
	"github.com/Zoomish/portfolio-backend/internal/content/github"
	"github.com/Zoomish/portfolio-backend/internal/content/linkedin"
)

type RepoSource interface {
	List(ctx context.Context) ([]github.Candidate, error)
}

type ProfileSource interface {
	Fetch(ctx context.Context, part linkedin.Part) (json.RawMessage, error)
}

type Store interface {
	Save(ctx context.Context, snap content.Snapshot) error
}
