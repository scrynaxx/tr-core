package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Params struct {
	User     string
	Password string
	Database string
	Host     string
	Port     uint16
	SSLMode  string
	MaxConns int32
}

func (p *Params) Validate() error {
	var errs []error

	if strings.TrimSpace(p.User) == "" {
		errs = append(errs, errors.New("user is empty"))
	}
	if p.Password == "" {
		errs = append(errs, errors.New("password is empty"))
	}
	if strings.TrimSpace(p.Database) == "" {
		errs = append(errs, errors.New("database is empty"))
	}
	if strings.TrimSpace(p.Host) == "" {
		errs = append(errs, errors.New("host is empty"))
	}
	if p.Port == 0 {
		errs = append(errs, errors.New("port is empty"))
	}
	if strings.TrimSpace(p.SSLMode) == "" {
		errs = append(errs, errors.New("ssl mode is empty"))
	}

	return errors.Join(errs...)
}

func NewPool(ctx context.Context, params Params, options ...PoolOption) (*pgxpool.Pool, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("validate params: %w", err)
	}

	config, err := pgxpool.ParseConfig(params.connectionString())
	if err != nil {
		return nil, fmt.Errorf("pool config parse: %w", err)
	}

	config.MaxConnIdleTime = 3 * time.Minute
	config.MaxConnLifetime = 15 * time.Minute
	config.HealthCheckPeriod = 30 * time.Second
	config.MaxConns = 30
	if params.MaxConns > 0 {
		config.MaxConns = params.MaxConns
	}
	config.MinConns = 5
	config.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("pool create: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pool ping: %w", err)
	}

	if err = otelpgx.RecordStats(pool); err != nil {
		return nil, fmt.Errorf("record pg stats: %w", err)
	}

	for _, option := range options {
		if err = option(ctx, pool); err != nil {
			pool.Close()
			return nil, fmt.Errorf("apply pool option: %w", err)
		}
	}

	return pool, nil
}

func (p *Params) connectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host,
		p.Port,
		p.User,
		p.Password,
		p.Database,
		p.SSLMode,
	)
}
