package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Postgres -.
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool
}

// New -.
func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return pg, nil
}

type txKey struct{}

func injectTx(ctx context.Context, tx *pgx.Tx) context.Context {
    return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) *pgx.Tx {
    if tx, ok := ctx.Value(txKey{}).(*pgx.Tx); ok {
        return tx
    }
    return nil
}

// Close -.
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

func (p *Postgres) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error{
    tx, err := p.Pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("postgres - WithinTransaction - p.Pool.Begin: %w", err)
    }

    err = tFunc(injectTx(ctx, &tx))
    if err != nil {
        tx.Rollback(ctx)
        return fmt.Errorf("postgres - WithinTransaction - tFunc: %w", err)
    }
    tx.Commit(ctx)
    return nil
}
