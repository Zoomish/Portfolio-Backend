package serve

import (
	"context"
	"encoding/json"

	"github.com/Zoomish/portfolio-backend/internal/content"
)

type PartBody struct {
	Data   json.RawMessage `json:"data"`
	Cached bool            `json:"cached"`
}

type AllBody struct {
	Posts   PartBody `json:"posts"`
	Details PartBody `json:"details"`
	Contact PartBody `json:"contact"`
}

func (s *Service) LinkedInAll(ctx context.Context) (AllBody, error) {
	sources := []content.Source{
		content.SourceLinkedInPosts,
		content.SourceLinkedInProfile,
		content.SourceLinkedInContact,
	}
	if err := s.ensure(ctx, sources...); err != nil {
		return AllBody{}, err
	}
	posts, err := s.part(ctx, content.SourceLinkedInPosts)
	if err != nil {
		return AllBody{}, err
	}
	details, err := s.part(ctx, content.SourceLinkedInProfile)
	if err != nil {
		return AllBody{}, err
	}
	contact, err := s.part(ctx, content.SourceLinkedInContact)
	if err != nil {
		return AllBody{}, err
	}
	return AllBody{Posts: posts, Details: details, Contact: contact}, nil
}

func (s *Service) LinkedInPart(ctx context.Context, source content.Source) (PartBody, error) {
	if err := s.ensure(ctx, source); err != nil {
		return PartBody{}, err
	}
	return s.part(ctx, source)
}

func (s *Service) part(ctx context.Context, source content.Source) (PartBody, error) {
	snap, err := s.read(ctx, source)
	if err != nil {
		return PartBody{}, err
	}
	return PartBody{Data: snap.Payload, Cached: true}, nil
}
