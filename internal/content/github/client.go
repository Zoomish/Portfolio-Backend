package github

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const pageSize = 100
const maxPages = 30

type Client struct {
	token string
	http  *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
		http:  &http.Client{Timeout: 20 * time.Second},
	}
}

func (c *Client) List(ctx context.Context) ([]Candidate, error) {
	var all []Candidate
	for page := 1; page <= maxPages; page++ {
		query := url.Values{}
		query.Set("per_page", strconv.Itoa(pageSize))
		query.Set("page", strconv.Itoa(page))
		query.Set("type", "owner")
		query.Set("sort", "updated")
		query.Set("direction", "desc")
		body, status, err := c.get(ctx, "https://api.github.com/user/repos?"+query.Encode(), map[string]string{
			"Accept": "application/vnd.github+json",
		})
		if err != nil {
			return nil, fmt.Errorf("github repos page %d: %w", page, err)
		}
		if status != http.StatusOK {
			return nil, fmt.Errorf("github repos page %d: status %d", page, status)
		}
		batch, err := DecodeRepos(body)
		if err != nil {
			return nil, fmt.Errorf("github repos page %d: %w", page, err)
		}
		if len(batch) == 0 {
			break
		}
		all = append(all, batch...)
		if len(batch) < pageSize {
			break
		}
	}
	return all, nil
}

func (c *Client) Readme(ctx context.Context, fullName string) (string, error) {
	body, status, err := c.get(ctx, "https://api.github.com/repos/"+fullName+"/readme", map[string]string{
		"Accept": "application/vnd.github.v3.raw",
	})
	if err != nil {
		return "", fmt.Errorf("github readme %s: %w", fullName, err)
	}
	if status == http.StatusNotFound {
		return "", nil
	}
	if status != http.StatusOK {
		return "", fmt.Errorf("github readme %s: status %d", fullName, status)
	}
	return string(body), nil
}

func (c *Client) get(ctx context.Context, rawURL string, headers map[string]string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("request: %w", err)
	}
	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("User-Agent", "portfolio-backend")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read body: %w", err)
	}
	return body, resp.StatusCode, nil
}
