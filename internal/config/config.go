package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Env              string
	Port             string
	Postgres         Postgres
	GitHubToken      string
	RapidAPIKey      string
	RapidAPIHost     string
	LinkedInUser     string
	RefreshEvery     time.Duration
	KeepaliveEvery   time.Duration
	KeepaliveEnabled bool
	AlertAfter       int
	TelegramToken    string
	AdminTelegramID  int64
}

func Load() (Config, error) {
	if err := loadDotEnv(".env"); err != nil {
		return Config{}, err
	}
	return load(os.Getenv)
}

func load(getenv func(string) string) (Config, error) {
	env := firstNonEmpty(getenv("APP_ENV"), getenv("NODE_ENV"), "dev")

	pg, err := loadPostgres(getenv, env)
	if err != nil {
		return Config{}, err
	}

	token, err := required(getenv, "GITHUB_TOKEN")
	if err != nil {
		return Config{}, err
	}
	apiKey, err := required(getenv, "RAPIDAPI_KEY")
	if err != nil {
		return Config{}, err
	}
	apiHost, err := required(getenv, "RAPIDAPI_HOST")
	if err != nil {
		return Config{}, err
	}

	port := strings.TrimSpace(getenv("PORT"))
	if port == "" {
		port = "3000"
	}
	if _, err := strconv.Atoi(port); err != nil {
		return Config{}, fmt.Errorf("PORT: %w", err)
	}

	refreshEvery, err := duration(getenv, "CONTENT_REFRESH_INTERVAL", "6h")
	if err != nil {
		return Config{}, err
	}
	keepaliveEvery, err := duration(getenv, "DB_KEEPALIVE_INTERVAL", "15m")
	if err != nil {
		return Config{}, err
	}

	alertAfter := 3
	if raw := strings.TrimSpace(getenv("ALERT_AFTER_FAILURES")); raw != "" {
		alertAfter, err = strconv.Atoi(raw)
		if err != nil || alertAfter < 1 {
			return Config{}, fmt.Errorf("ALERT_AFTER_FAILURES must be a positive integer")
		}
	}

	var adminID int64
	if raw := strings.TrimSpace(getenv("ADMIN_TELEGRAM_ID")); raw != "" {
		adminID, err = strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return Config{}, fmt.Errorf("ADMIN_TELEGRAM_ID: %w", err)
		}
	}

	user := strings.TrimSpace(getenv("LINKEDIN_USERNAME"))
	if user == "" {
		user = "zoomish"
	}

	return Config{
		Env:              env,
		Port:             port,
		Postgres:         pg,
		GitHubToken:      token,
		RapidAPIKey:      apiKey,
		RapidAPIHost:     apiHost,
		LinkedInUser:     user,
		RefreshEvery:     refreshEvery,
		KeepaliveEvery:   keepaliveEvery,
		KeepaliveEnabled: !strings.EqualFold(strings.TrimSpace(getenv("DB_KEEPALIVE_ENABLED")), "false"),
		AlertAfter:       alertAfter,
		TelegramToken:    strings.TrimSpace(getenv("TELEGRAM_TOKEN")),
		AdminTelegramID:  adminID,
	}, nil
}

func required(getenv func(string) string, key string) (string, error) {
	value := strings.TrimSpace(getenv(key))
	if value == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return value, nil
}

func duration(getenv func(string) string, key, fallback string) (time.Duration, error) {
	raw := strings.TrimSpace(getenv(key))
	if raw == "" {
		raw = fallback
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", key, err)
	}
	if d <= 0 {
		return 0, fmt.Errorf("%s must be positive", key)
	}
	return d, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
