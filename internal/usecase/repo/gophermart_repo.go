package repo

import (
	"context"

	"github.com/costynus/loyalty-system/internal/entity"
	"github.com/costynus/loyalty-system/pkg/postgres"
	"github.com/shopspring/decimal"
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


func (r *GophermartRepo) CreateUser(ctx context.Context, login, password_hash string) (entity.User, error) {
    return entity.User{}, nil
}

func (r *GophermartRepo) GetUserWithLogin(ctx context.Context, login string) (entity.User, error) {
    return entity.User{}, nil
}

func (r *GophermartRepo) GetOrderList(ctx context.Context, user_id int) ([]entity.Order, error) {
    return nil, nil
}

func (r *GophermartRepo) GetCurrentBalance(ctx context.Context, user_id int) (entity.Balance, error) {
    return entity.Balance{}, nil
}

func (r *GophermartRepo) UpdateBalance(ctx context.Context, user_id int, value decimal.Decimal) error {
    return nil
}

func (r *GophermartRepo) AddWithdrawal(ctx context.Context, user_id int, order_num string, value decimal.Decimal) error {
    return nil
}

func (r *GophermartRepo) GetWithdrawalList(ctx context.Context, user_id int, order_num string) ([]entity.Withdraw, error) {
    return nil, nil
}
