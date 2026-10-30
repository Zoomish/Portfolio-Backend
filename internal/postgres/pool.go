package postgres

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/Zoomish/portfolio-backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, pg config.Postgres) (*pgxpool.Pool, error) {
	rawURL := connectionURL(pg)
	poolCfg, err := pgxpool.ParseConfig(rawURL)
	if err != nil {
		return nil, fmt.Errorf("postgres config: %w", err)
	}
	if pg.SSLMode != "disable" {
		tlsCfg, err := tlsConfig(pg)
		if err != nil {
			return nil, err
		}
		poolCfg.ConnConfig.TLSConfig = tlsCfg
	}
	poolCfg.MaxConns = 4

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("postgres pool: %w", err)
	}
	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres ping: %w", err)
	}
	return pool, nil
}

func connectionURL(pg config.Postgres) string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(pg.User, pg.Password),
		Host:   net.JoinHostPort(pg.Host, strconv.Itoa(pg.Port)),
		Path:   "/" + pg.Database,
	}
	query := u.Query()
	query.Set("sslmode", pg.SSLMode)
	u.RawQuery = query.Encode()
	return u.String()
}

func tlsConfig(pg config.Postgres) (*tls.Config, error) {
	pem, err := caPEM(pg)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(pem) {
		return nil, fmt.Errorf("postgres ca pem has no certificates")
	}
	return &tls.Config{
		RootCAs:    pool,
		ServerName: pg.Host,
		MinVersion: tls.VersionTLS12,
	}, nil
}

func caPEM(pg config.Postgres) ([]byte, error) {
	if pg.CACert != "" {
		return []byte(pg.CACert), nil
	}
	if pg.CAFile == "" {
		return nil, fmt.Errorf("postgres ca is empty")
	}
	body, err := os.ReadFile(pg.CAFile)
	if err != nil {
		return nil, fmt.Errorf("read postgres ca: %w", err)
	}
	return body, nil
}
