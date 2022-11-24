package repo

import (
	"context"
	"fmt"

	 sq "github.com/Masterminds/squirrel"
	"github.com/costynus/loyalty-system/internal/entity"
	"github.com/costynus/loyalty-system/pkg/postgres"
	"github.com/shopspring/decimal"
    "github.com/georgysavva/scany/pgxscan"
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


func (r *GophermartRepo) CreateUser(ctx context.Context, login, passwordHash string) (entity.User, error) {
    var user entity.User
    user.Login = login
    user.PasswordHash = passwordHash

    sql, args, err := r.Builder.
        Insert("public.user").
        Columns("login", "password_hash").
        Values(login, passwordHash).
        Suffix("RETURNING id").
        ToSql()
    if err != nil {
        return entity.User{}, fmt.Errorf("GophermartRepo - CreateUser - r.Builder: %w", err)
    }

    tx, err := r.Pool.Begin(ctx)
    if err != nil {
        return entity.User{}, fmt.Errorf("GophermartRepo - CreateUser - r.Pool.Begin: %w", err)
    }
    defer tx.Rollback(ctx)

    err = tx.QueryRow(ctx, sql, args...).Scan(&user.ID)

    if err != nil {
        return entity.User{}, fmt.Errorf("GophermartRepo - CreateUser - tx.QueryRow: %w", err)
    }

    tx.Commit(ctx)
    return user, nil
}

func (r *GophermartRepo) GetUserWithLogin(ctx context.Context, login string) (entity.User, error) {
    sql, args, err := r.Builder.
        Select("id", "login", "password_hash").
        From("public.user").
        Where(sq.Eq{"login": login}).
        ToSql()

    if err != nil {
        return entity.User{}, fmt.Errorf("GophermartRepo - GetUserWithLogin - r.Builder: %w", err)
    }

    dst := make([]entity.User, 0)
    if err = pgxscan.Select(ctx, r.Pool, &dst, sql, args...); err != nil {
        return entity.User{}, fmt.Errorf("GophermartRepo - GetUserWithLogin - pgxscan.Select: %w", err)
    }

    if len(dst) == 0 {
        return entity.User{}, ErrNotFound
    }

    return dst[0], nil
}

func (r *GophermartRepo) GetOrderList(ctx context.Context, userId int) ([]entity.Order, error) {
    sql, args, err := r.Builder.
        Select("order_number", "status", "accrual", "uploaded_at").
        From("public.order").
        Where(sq.Eq{"user_id": userId}).
        ToSql()
    if err != nil {
        return nil, fmt.Errorf("gophermartRepo - GetOrderList - r.Builder: %w", err)
    }

    dst := make([]entity.Order, 0)
    if err = pgxscan.Select(ctx, r.Pool, &dst, sql, args...); err != nil {
        return nil, fmt.Errorf("GophermartRepo - GetOrderList - pgxscan.Select: %w", err)
    }
    return dst, nil
}

func (r *GophermartRepo) GetCurrentBalance(ctx context.Context, userId int) (entity.Balance, error) {
    sql, args, err := r.Builder.
        Select("balance", "withdrawal").
        From("public.balance").
        Where(sq.Eq{"user_id": userId}).
        ToSql()
    if err != nil {
        return entity.Balance{}, fmt.Errorf("GophermartRepo - GetcurrentBalance - r.Builder: %w", err)
    }

    dst := make([]entity.Balance, 0)
    if err = pgxscan.Select(ctx, r.Pool, &dst, sql, args...); err != nil {
        return entity.Balance{}, fmt.Errorf("GophermartRepo - GetCurrentBalance - pgxsan.Select: %w", err)
    }

    if len(dst) == 0 {
        return entity.Balance{}, ErrNotFound
    }

    return dst[0], nil
}

func (r *GophermartRepo) UpdateBalance(ctx context.Context, userId int, value decimal.Decimal) error {
    // TODO: code me pls
    return nil
}

func (r *GophermartRepo) AddWithdrawal(ctx context.Context, userId int, orderNum string, value decimal.Decimal) error {
    // TODO: code me pls
    return nil
}

func (r *GophermartRepo) GetWithdrawalList(ctx context.Context, userId int, orderNum string) ([]entity.Withdraw, error) {
    sql, args, err := r.Builder.
        Select("order_num", "sum_number", "updated_at").
        From("public.withdrawal").
        Where(sq.Eq{"user_id": userId}).
        ToSql()
    if err != nil {
        return nil, fmt.Errorf("GophermartRepo - GetWithdrawalList - r.Builder: %w", err)
    }

    dst := make([]entity.Withdraw, 0)
    if err = pgxscan.Select(ctx, r.Pool, &dst, sql, args...); err != nil {
        return nil, fmt.Errorf("GophermartRepo - GetWithdrawalList - pgxscan.Select: %w", err)
    }

    return dst, nil
}
