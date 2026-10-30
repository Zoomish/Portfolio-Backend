package serve

import (
	"context"
	"encoding/json"

	"github.com/Zoomish/portfolio-backend/internal/content"
)

type GithubBody struct {
	Data   json.RawMessage `json:"data"`
	Cached bool            `json:"cached"`
}

func (s *Service) Github(ctx context.Context) (GithubBody, error) {
	if err := s.ensure(ctx, content.SourceGithubRepos); err != nil {
		return GithubBody{}, err
	}
	snap, err := s.read(ctx, content.SourceGithubRepos)
	if err != nil {
		return GithubBody{}, err
	}
	return GithubBody{Data: snap.Payload, Cached: true}, nil
}
