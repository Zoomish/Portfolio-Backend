package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoadRequiresGitHubToken(t *testing.T) {
	env := map[string]string{
		"POSTGRES_HOST":     "localhost",
		"POSTGRES_PORT":     "5432",
		"POSTGRES_USER":     "user",
		"POSTGRES_PASSWORD": "secret",
		"POSTGRES_DATABASE": "portfolio",
		"POSTGRES_SSL":      "disable",
		"RAPIDAPI_KEY":      "key",
		"RAPIDAPI_HOST":     "example.rapidapi.com",
	}

	_, err := load(func(key string) string { return env[key] })
	if err == nil || !strings.Contains(err.Error(), "GITHUB_TOKEN") {
		t.Fatalf("expected missing token error, got %v", err)
	}
}

func TestLoadDefaults(t *testing.T) {
	env := map[string]string{
		"POSTGRES_HOST":     "localhost",
		"POSTGRES_USER":     "user",
		"POSTGRES_PASSWORD": "secret",
		"POSTGRES_DATABASE": "portfolio",
		"POSTGRES_SSL":      "disable",
		"GITHUB_TOKEN":      "gh",
		"RAPIDAPI_KEY":      "key",
		"RAPIDAPI_HOST":     "example.rapidapi.com",
		"APP_ENV":           "dev",
	}

	cfg, err := load(func(key string) string { return env[key] })
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != "3000" {
		t.Fatalf("port %s", cfg.Port)
	}
	if cfg.RefreshEvery != 6*time.Hour {
		t.Fatalf("refresh %s", cfg.RefreshEvery)
	}
	if cfg.KeepaliveEvery != 15*time.Minute {
		t.Fatalf("keepalive %s", cfg.KeepaliveEvery)
	}
	if cfg.LinkedInUser != "zoomish" {
		t.Fatalf("linkedin user %s", cfg.LinkedInUser)
	}
	if !cfg.KeepaliveEnabled {
		t.Fatal("keepalive should be enabled")
	}
	if cfg.Postgres.SSLMode != "disable" {
		t.Fatalf("ssl %s", cfg.Postgres.SSLMode)
	}
	if cfg.AlertAfter != 3 {
		t.Fatalf("alert after %d", cfg.AlertAfter)
	}
}
