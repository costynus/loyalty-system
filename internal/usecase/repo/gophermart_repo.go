package repo

import (
	"context"

	"github.com/costynus/loyalty-system/pkg/postgres"
)

type GophermartRepo struct {
    *postgres.Postgres
}

func New(pg *postgres.Postgres) *GophermartRepo{
    return &GophermartRepo{pg}
}

func (r *GophermartRepo) Ping(ctx context.Context) error {
    return r.Pool.Ping(ctx)
}
