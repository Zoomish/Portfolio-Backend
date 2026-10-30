package refresh

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Zoomish/portfolio-backend/internal/content"
	"github.com/Zoomish/portfolio-backend/internal/content/github"
	"github.com/Zoomish/portfolio-backend/internal/content/linkedin"
)

type Cycle struct {
	repos   RepoSource
	profile ProfileSource
	store   Store
}

func NewCycle(repos RepoSource, profile ProfileSource, store Store) *Cycle {
	return &Cycle{repos: repos, profile: profile, store: store}
}

func (c *Cycle) Run(ctx context.Context) error {
	errs := make([]error, 4)
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		errs[0] = c.github(ctx)
	}()
	go func() {
		defer wg.Done()
		errs[1] = c.profilePart(ctx, linkedin.Posts, content.SourceLinkedInPosts)
	}()
	go func() {
		defer wg.Done()
		errs[2] = c.profilePart(ctx, linkedin.Profile, content.SourceLinkedInProfile)
	}()
	go func() {
		defer wg.Done()
		errs[3] = c.profilePart(ctx, linkedin.Contact, content.SourceLinkedInContact)
	}()
	wg.Wait()
	return errors.Join(errs...)
}

func (c *Cycle) github(ctx context.Context) error {
	list, err := c.repos.List(ctx)
	if err != nil {
		return fmt.Errorf("github list: %w", err)
	}
	selected, err := github.Select(list, func(fullName string) (string, error) {
		return c.repos.Readme(ctx, fullName)
	})
	if err != nil {
		return err
	}
	payload, err := json.Marshal(selected)
	if err != nil {
		return fmt.Errorf("github encode: %w", err)
	}
	return c.store.Save(ctx, content.Snapshot{
		Source:    content.SourceGithubRepos,
		Payload:   payload,
		FetchedAt: time.Now().UTC(),
	})
}

func (c *Cycle) profilePart(ctx context.Context, part linkedin.Part, source content.Source) error {
	raw, err := c.profile.Fetch(ctx, part)
	if err != nil {
		return fmt.Errorf("linkedin %s: %w", part, err)
	}
	return c.store.Save(ctx, content.Snapshot{
		Source:    source,
		Payload:   raw,
		FetchedAt: time.Now().UTC(),
	})
}
