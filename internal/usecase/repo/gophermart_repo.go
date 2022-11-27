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

func (r *GophermartRepo) UpdateOrderAccrual(ctx context.Context, orderNumber string, accrual decimal.Decimal) error {
    sql, args, err := r.Builder.
        Update("public.order").
        Set("accrual", accrual).
        Where(sq.Eq{"order_number": orderNumber}).
        ToSql()
    if err != nil {
        return fmt.Errorf("GophermartRepo - UpdateOrderAccrual - r.Builder: %w", err)
    }

    _, err = r.Pool.Exec(ctx, sql, args...)
    if err != nil {
        return fmt.Errorf("GophermartRepo - UpdateOrderAccrual - r.Pool.Exec: %w", err)
    }
    return nil
}

func (r *GophermartRepo) UpdateOrderStatus(ctx context.Context, orderNumber, status string) error {
    sql, args, err := r.Builder.
        Update("public.order").
        Set("status", status).
        Where(sq.Eq{"order_number": orderNumber}).
        ToSql()
    if err != nil {
        return fmt.Errorf("GophermartRepo - UpdateOrderStatus - r.builder: %w", err)
    }
    _, err = r.Pool.Exec(ctx, sql, args...)
    if err != nil {
        return fmt.Errorf("GophermartRepo - UpdateOrderStatus - r.Pool.Exec: %w", err)
    }
    return nil
}

func (r *GophermartRepo) GetOrderByOrderNumber(ctx context.Context, orderNumber string) (entity.Order, error) {
    sql, args, err := r.Builder.
        Select("order_number", "status", "accrual", "uploaded_at", "user_id").
        From("public.order").
        Where(sq.Eq{"order_number": orderNumber}).
        OrderBy("uploaded_at").
        ToSql()
    if err != nil {
        return entity.Order{}, fmt.Errorf("GophermartRepo - GetOrderByOrderNumber - r.Builder: %w", err)
    }
    
    dst := make([]entity.Order, 0)
    if err = pgxscan.Select(ctx, r.Pool, &dst, sql, args...); err != nil {
        return entity.Order{}, fmt.Errorf("GophermartRepo - GetOrderByOrderNumber - pgxsan.Select: %w", err)
    }

    if len(dst) == 0 {
        return entity.Order{}, ErrNotFound
    }

    return dst[0], nil
}

func (r *GophermartRepo) CreateUserBalance(ctx context.Context, userID int) error {
    sql, args, err := r.Builder.
        Insert("public.balance").
        Columns("balance", "withdrawal", "user_id").
        Values(0, 0, userID).
        ToSql()
    if err != nil {
        return fmt.Errorf("GophermartRepo - CreateUserBalance - r.Builder: %w", err)
    }

    _, err = r.Pool.Exec(ctx, sql, args...)
    if err != nil {
        return fmt.Errorf("GophermartRepo - CreateUserBalance - r.Pool.Exec: %w", err)
    }
    return nil
}

func (r *GophermartRepo) CreateOrder(ctx context.Context, userID int, orderNumber string) error {
    sql, args, err := r.Builder.
        Insert("public.order").
        Columns("order_number", "user_id", "accrual").
        Values(orderNumber, userID, 0).
        ToSql()
    if err != nil {
        return fmt.Errorf("GophermartRepo - CreateOrder - r.Builder: %w", err)
    }

    _, err = r.Pool.Exec(ctx, sql, args...)
    if err != nil {
        return fmt.Errorf("GophermartRepo - CreateOrder - r.Pool.Exec: %w", err)
    }

    return nil
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

func (r *GophermartRepo) GetOrderList(ctx context.Context, userID int) ([]entity.Order, error) {
    sql, args, err := r.Builder.
        Select("order_number", "status", "accrual", "uploaded_at").
        From("public.order").
        Where(sq.Eq{"user_id": userID}).
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

func (r *GophermartRepo) GetCurrentBalance(ctx context.Context, userID int) (entity.Balance, error) {
    sql, args, err := r.Builder.
        Select("balance", "withdrawal").
        From("public.balance").
        Where(sq.Eq{"user_id": userID}).
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

func (r *GophermartRepo) UpdateBalance(ctx context.Context, userID int, balance, withdrawal decimal.Decimal) error {
    sql, args, err := r.Builder.
        Update("public.balance").
        Set("balance", balance).
        Set("withdrawal", withdrawal).
        Where(sq.Eq{"user_id": userID}).
        ToSql()
    if err != nil {
        return fmt.Errorf("GophermartRepo - UpdateBalance - r.Builder: %w", err)
    }

    _, err = r.Pool.Exec(ctx, sql, args...)
    if err != nil {
        return fmt.Errorf("GophermartRepo - UpdateBalance - r.Pool.Exec: %w", err)
    }
    return nil
}

func (r *GophermartRepo) AddWithdrawal(ctx context.Context, userID int, orderNum string, value decimal.Decimal) error {
    sql, args, err := r.Builder.
        Insert("public.withdrawal").
        Columns("order_number", "sum_number", "user_id").
        Values(orderNum, value, userID).
        ToSql()
    if err != nil {
        return fmt.Errorf("GophermartRepo - AddWithdrawal - r.Builder: %w", err)
    }

    _, err = r.Pool.Exec(ctx, sql, args...)
    if err != nil {
        return fmt.Errorf("GophermartRepo - AddWithdrawal - r.Pool.Exec: %w", err)
    }
    return nil
}

func (r *GophermartRepo) GetWithdrawalList(ctx context.Context, userID int, orderNum string) ([]entity.Withdraw, error) {
    sql, args, err := r.Builder.
        Select("order_number", "sum_number", "updated_at").
        From("public.withdrawal").
        Where(sq.Eq{"user_id": userID}).
        ToSql()
    if err != nil {
        return nil, fmt.Errorf("GophermartRepo - GetWithdrawalList - r.Builder: %w", err)
    }

    dst := make([]entity.Withdraw, 0)
    if err = pgxscan.Select(ctx, r.Pool, &dst, sql, args...); err != nil {
        return nil, fmt.Errorf("GophermartRepo - GetWithdrawalList - pgxscan.Select: %w", err)
    }
    if len(dst) == 0 {
        return nil, nil
    }

    return dst, nil
}
