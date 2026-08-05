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

type Options struct {
	User     string
	Password string
	Database string
	Host     string
	Port     uint16
	SSLMode  string
	MaxConns int32
}

func (c *Options) Validate() error {
	var errs []error

	if strings.TrimSpace(c.User) == "" {
		errs = append(errs, errors.New("user is empty"))
	}
	if c.Password == "" {
		errs = append(errs, errors.New("password is empty"))
	}
	if strings.TrimSpace(c.Database) == "" {
		errs = append(errs, errors.New("database is empty"))
	}
	if strings.TrimSpace(c.Host) == "" {
		errs = append(errs, errors.New("host is empty"))
	}
	if c.Port == 0 {
		errs = append(errs, errors.New("port is empty"))
	}
	if strings.TrimSpace(c.SSLMode) == "" {
		errs = append(errs, errors.New("ssl mode is empty"))
	}

	return errors.Join(errs...)
}

type Postgres struct {
	Pool       *pgxpool.Pool
	Transactor *Transactor
}

func New(ctx context.Context, o Options, options ...PoolOption) (*Postgres, error) {
	if err := o.Validate(); err != nil {
		return nil, fmt.Errorf("validate params: %w", err)
	}

	str := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		o.Host, o.Port, o.User, o.Password, o.Database, o.SSLMode)

	config, err := pgxpool.ParseConfig(str)
	if err != nil {
		return nil, fmt.Errorf("pool config parse: %w", err)
	}

	config.MaxConnIdleTime = 3 * time.Minute
	config.MaxConnLifetime = 15 * time.Minute
	config.HealthCheckPeriod = 30 * time.Second
	config.MaxConns = 30
	if o.MaxConns > 0 {
		config.MaxConns = o.MaxConns
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

	return &Postgres{
		Pool:       pool,
		Transactor: newTransactor(pool),
	}, nil
}
