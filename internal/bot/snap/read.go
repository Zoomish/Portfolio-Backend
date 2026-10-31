package snap

import (
	"context"
	"encoding/json"

	"github.com/Zoomish/portfolio-backend/internal/content"
	"github.com/Zoomish/portfolio-backend/internal/content/serve"
)

type Reader struct {
	Snaps *serve.Service
}

func (r Reader) Profile(ctx context.Context) (json.RawMessage, error) {
	body, err := r.Snaps.LinkedInPart(ctx, content.SourceLinkedInProfile)
	if err != nil {
		return nil, err
	}
	return body.Data, nil
}

func (r Reader) Contact(ctx context.Context) (json.RawMessage, error) {
	body, err := r.Snaps.LinkedInPart(ctx, content.SourceLinkedInContact)
	if err != nil {
		return nil, err
	}
	return body.Data, nil
}

func (r Reader) Projects(ctx context.Context) (json.RawMessage, error) {
	body, err := r.Snaps.Github(ctx)
	if err != nil {
		return nil, err
	}
	return body.Data, nil
}

func (r Reader) Posts(ctx context.Context) (json.RawMessage, error) {
	body, err := r.Snaps.LinkedInPart(ctx, content.SourceLinkedInPosts)
	if err != nil {
		return nil, err
	}
	return body.Data, nil
}
