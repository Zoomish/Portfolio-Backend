package postgres

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/Zoomish/portfolio-backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, pg config.Postgres) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(connectionURL(pg))
	if err != nil {
		return nil, fmt.Errorf("postgres config: %w", err)
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
