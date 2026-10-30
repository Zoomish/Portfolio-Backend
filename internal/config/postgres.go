package config

import (
	"fmt"
	"os"
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
	CAFile   string
	CACert   string
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

	sslMode, caFile, caCert, err := postgresSSL(getenv, env)
	if err != nil {
		return Postgres{}, err
	}

	return Postgres{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
		SSLMode:  sslMode,
		CAFile:   caFile,
		CACert:   caCert,
	}, nil
}

func postgresSSL(getenv func(string) string, env string) (string, string, string, error) {
	if strings.EqualFold(strings.TrimSpace(getenv("POSTGRES_SSL")), "disable") {
		return "disable", "", "", nil
	}

	caCert := strings.TrimSpace(getenv("POSTGRES_SSL_CA"))
	caFile := strings.TrimSpace(getenv("POSTGRES_SSL_CA_FILE"))
	if caCert == "" && caFile == "" {
		if env == "dev" || env == "development" {
			return "disable", "", "", nil
		}
		if _, statErr := os.Stat("ca.pem"); statErr == nil {
			return "verify-full", "ca.pem", "", nil
		}
		return "", "", "", fmt.Errorf("postgres ssl ca is required (POSTGRES_SSL_CA or POSTGRES_SSL_CA_FILE)")
	}
	return "verify-full", caFile, caCert, nil
}
