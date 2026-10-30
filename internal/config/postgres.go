package config

import (
	"fmt"
	"strconv"
	"strings"
)

type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

func loadPostgres(getenv func(string) string, env string) (Postgres, error) {
	host, err := required(getenv, "POSTGRES_HOST")
	if err != nil {
		return Postgres{}, err
	}
	user, err := required(getenv, "POSTGRES_USER")
	if err != nil {
		return Postgres{}, err
	}
	password, err := required(getenv, "POSTGRES_PASSWORD")
	if err != nil {
		return Postgres{}, err
	}
	database, err := required(getenv, "POSTGRES_DATABASE")
	if err != nil {
		return Postgres{}, err
	}

	portRaw := strings.TrimSpace(getenv("POSTGRES_PORT"))
	if portRaw == "" {
		portRaw = "5432"
	}
	port, err := strconv.Atoi(portRaw)
	if err != nil {
		return Postgres{}, fmt.Errorf("POSTGRES_PORT: %w", err)
	}

	return Postgres{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
		SSLMode:  postgresSSL(getenv, env),
	}, nil
}

func postgresSSL(getenv func(string) string, env string) string {
	if strings.EqualFold(strings.TrimSpace(getenv("POSTGRES_SSL")), "disable") {
		return "disable"
	}
	if env == "dev" || env == "development" {
		return "disable"
	}
	return "require"
}
