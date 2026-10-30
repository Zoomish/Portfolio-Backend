package linkedin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	key      string
	host     string
	username string
	http     *http.Client
}

func NewClient(key, host, username string) *Client {
	return &Client{
		key:      key,
		host:     host,
		username: username,
		http:     &http.Client{Timeout: 20 * time.Second},
	}
}

func (c *Client) Fetch(ctx context.Context, part Part) (json.RawMessage, error) {
	path := part.Path()
	if path == "" {
		return nil, fmt.Errorf("unknown linkedin part %q", part)
	}
	endpoint := "https://" + c.host + "/" + path + "?" + part.Query(c.username).Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("linkedin %s request: %w", part, err)
	}
	req.Header.Set("x-rapidapi-key", c.key)
	req.Header.Set("x-rapidapi-host", c.host)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("linkedin %s: %w", part, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, fmt.Errorf("linkedin %s body: %w", part, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("linkedin %s: status %d", part, resp.StatusCode)
	}
	if !json.Valid(body) {
		return nil, fmt.Errorf("linkedin %s: invalid json", part)
	}
	return json.RawMessage(body), nil
}
