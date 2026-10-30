package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zoomish/portfolio-backend/internal/content"
	"github.com/Zoomish/portfolio-backend/internal/content/serve"
)

func TestGithubReposJSON(t *testing.T) {
	handler := New(
		stubGithub{body: serve.GithubBody{Data: json.RawMessage(`[{"name":"api","link":null}]`), Cached: true}},
		stubLinkedIn{},
		stubReady{},
	)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/github/repos", nil)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("missing cors")
	}
	var body serve.GithubBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if !body.Cached || string(body.Data) != `[{"name":"api","link":null}]` {
		t.Fatalf("body %s", rec.Body.Bytes())
	}
}

func TestLinkedInAllJSON(t *testing.T) {
	part := serve.PartBody{Data: json.RawMessage(`{"success":true}`), Cached: true}
	handler := New(stubGithub{}, stubLinkedIn{all: serve.AllBody{Posts: part, Details: part, Contact: part}}, stubReady{})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/linkedin/all", nil)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.Bytes())
	}
	var body serve.AllBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if !body.Posts.Cached || !body.Details.Cached || !body.Contact.Cached {
		t.Fatalf("body %+v", body)
	}
}

func TestContentNotReady(t *testing.T) {
	handler := New(stubGithub{err: content.ErrNotReady}, stubLinkedIn{}, stubReady{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/github/repos", nil)
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status %d", rec.Code)
	}
}

func TestAppLiveness(t *testing.T) {
	handler := New(stubGithub{}, stubLinkedIn{}, stubReady{err: errors.New("db down")})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/app", nil)
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("liveness %d", rec.Code)
	}

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/ready", nil)
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("ready %d", rec.Code)
	}
}

func TestOptionsCORS(t *testing.T) {
	handler := New(stubGithub{}, stubLinkedIn{}, stubReady{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/api/github/repos", nil)
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("missing cors")
	}
}

type stubGithub struct {
	body serve.GithubBody
	err  error
}

func (s stubGithub) Github(context.Context) (serve.GithubBody, error) {
	return s.body, s.err
}

type stubLinkedIn struct {
	all serve.AllBody
	err error
}

func (s stubLinkedIn) LinkedInAll(context.Context) (serve.AllBody, error) {
	return s.all, s.err
}

func (s stubLinkedIn) LinkedInPart(context.Context, content.Source) (serve.PartBody, error) {
	return serve.PartBody{}, s.err
}

type stubReady struct {
	err error
}

func (s stubReady) Ready(context.Context) error { return s.err }
