package usecase

import (
	"context"
	"fmt"

	"github.com/costynus/loyalty-system/internal/entity"
	"github.com/costynus/loyalty-system/internal/usecase/repo"
	"github.com/costynus/loyalty-system/internal/usecase/webapi"
	"github.com/costynus/loyalty-system/pkg/password"
)

type GophermartUseCase struct {
    repo GophermartRepo
    webAPI GophermartWebAPI
    orderCh chan <- string
}

func New(r GophermartRepo, w GophermartWebAPI, workersCount int) *GophermartUseCase{
    orderCh := make(chan string)
    uc := &GophermartUseCase{
        repo: r,
        webAPI: w,
        orderCh: orderCh,
    }

    for i := 0; i < workersCount; i++ {
        go func() {
            for orderNumber := range orderCh {
                err := uc.ProcessOrder(orderNumber)
                switch err {
                case webapi.ErrTooManyRequests:
                    orderCh <- orderNumber
                }
            }
        }()
    }

    return uc
}

func (uc *GophermartUseCase) ProcessOrder(orderNumber string) error {
    order, err := uc.webAPI.GetOrderInfo(orderNumber)
    switch err {
    case webapi.ErrInternalServerError:
        err = uc.repo.UpdateOrderStatus(context.TODO(), orderNumber, "INVALID")
        if err != nil {
            return fmt.Errorf("GophermartUseCase - ProcessOrder - uc.repo.UpdateOrderStatus: %w", err)
        }
        return err
    case webapi.ErrTooManyRequests:
        return err
    default:
        if err != nil {
            return err
        }
        addOrder, err := uc.repo.GetOrderByOrderNumber(context.TODO(), orderNumber)
        if err != nil {
            return fmt.Errorf("GophermartUseCase - ProcessOrder - uc.repo.GetOrderByOrderNumber: %w", err)
        }
        order.UserID = addOrder.UserID
        err = uc.repo.UpdateOrderStatus(context.TODO(), orderNumber, order.Status)
        if err != nil {
            return fmt.Errorf("GophermartUseCase - ProcessOrder - uc.repo.UpdateOrderStatus: %w", err)
        }
        switch order.Status {
            case "PROCESSED": 
                err = uc.repo.UpdateOrderAccrual(context.TODO(), orderNumber, order.Accrual)
                if err != nil {
                    return fmt.Errorf("GophermartUseCase - ProcessOrder - uc.repo.UpdateOrderAccrual: %w", err)
                }
                balance, err := uc.repo.GetCurrentBalance(context.TODO(), order.UserID)
                if err != nil {
                    return fmt.Errorf("GophermartUseCase - ProcessOrder - uc.repo.GetCurrentBalance: %w", err)
                }
                err = uc.repo.UpdateBalance(context.TODO(), order.UserID, balance.Current.Add(order.Accrual), balance.Withdraw)
                if err != nil {
                    return fmt.Errorf("GophermartUseCase - ProcessOrder - uc.repo.UpdateBalance: %w", err)
                }
            case "INVALID":
                return nil
            default:
                uc.orderCh <- orderNumber
        }
    }
    return nil
}

func (uc *GophermartUseCase) PingRepo(ctx context.Context) error {
    return uc.repo.Ping(ctx)
}


func (uc *GophermartUseCase) CreateNewUser(ctx context.Context, userAuth entity.UserAuth) (entity.User, error) {
    passwordHash, err := password.HashPassword(userAuth.Password)
    if err != nil {
        return entity.User{}, err
    }

    _, err = uc.repo.GetUserWithLogin(ctx, userAuth.Login)
    if err == nil {
        return entity.User{}, ErrConflict
    }

    user, err := uc.repo.CreateUser(ctx, userAuth.Login, passwordHash) 
    if err != nil {
        return entity.User{}, err
    }

    return user, nil
}

func (uc *GophermartUseCase) CheckUser(ctx context.Context, userAuth entity.UserAuth) (entity.User, error){
    passwordHash, err := password.HashPassword(userAuth.Password)
    if err != nil {
        return entity.User{}, err
    }

    user, err := uc.repo.GetUserWithLogin(ctx, userAuth.Login)
    if err != nil {
        return entity.User{}, ErrUnauthorized
    }

    isValidPassword := password.CheckPasswordHash(userAuth.Password, passwordHash)
    if !isValidPassword {
        return entity.User{}, ErrUnauthorized
    }


    return user, nil
}

func (uc *GophermartUseCase) UploadOrder(ctx context.Context, userID int, orderNum string) (bool, error) {
    order, err := uc.repo.GetOrderByOrderNumber(ctx, orderNum)
    switch err {
    case nil:
        if order.UserID == userID {
            return true, nil
        }
        return false, ErrConflict
    case repo.ErrNotFound:
        err = uc.repo.CreateOrder(ctx, userID, orderNum)
        if err != nil {
            return false, fmt.Errorf("GophermartUseCase - UploadOrder - uc.repo.CreateOrder: %w", err)
        }
        uc.orderCh <- orderNum
    default:
        return false, fmt.Errorf("GophermartUseCase - UploadOrder - uc.repo.GetOrderByOrderNumber: %w", err)

    }

    return false, nil
}


func (uc *GophermartUseCase) GetOrderList(ctx context.Context, userID int) ([]entity.Order, error) {
    orderList, err := uc.repo.GetOrderList(ctx, userID)
    return orderList, err
}


func (uc *GophermartUseCase) GetCurrentBalance(ctx context.Context, userID int) (entity.Balance, error) {
    balance, err := uc.repo.GetCurrentBalance(ctx, userID)
    return balance, err
}


func (uc *GophermartUseCase) Withdraw(ctx context.Context, userID int, withdrawal entity.Withdrawal) error {
    balance, err := uc.repo.GetCurrentBalance(ctx, userID)
    if err != nil {
        return err
    }

    if withdrawal.Sum.GreaterThan(balance.Current) {
        return ErrPaymentRequired
    }

    withdrawalList, err := uc.repo.GetWithdrawalList(ctx, userID, withdrawal.Order)
    if err != nil {
        return err
    }
    if withdrawalList != nil {
        return ErrUnprocessableEntity
    }

    err = uc.repo.UpdateBalance(ctx, userID, balance.Current.Sub(withdrawal.Sum), balance.Withdraw.Add(withdrawal.Sum))
    if err != nil {
        return err
    }

    err = uc.repo.AddWithdrawal(ctx, userID, withdrawal.Order, withdrawal.Sum)
    if err != nil {
        return err
    }

    return nil
}


func (uc *GophermartUseCase) GetWithdrawList(ctx context.Context, userID int) ([]entity.Withdraw, error) {
    withdrawalList, err := uc.repo.GetWithdrawalList(ctx, userID, "")
    return withdrawalList, err
}