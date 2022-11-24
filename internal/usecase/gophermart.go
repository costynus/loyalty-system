package usecase

import (
	"context"
	"fmt"

	"github.com/costynus/loyalty-system/internal/entity"
	"github.com/costynus/loyalty-system/pkg/password"
)

type GophermartUseCase struct {
    repo GophermartRepo
}

func New(r GophermartRepo) *GophermartUseCase{
    return &GophermartUseCase{
        repo: r,
    }
}

func (uc *GophermartUseCase) PingRepo(ctx context.Context) error {
    return uc.repo.Ping(ctx)
}


func (uc *GophermartUseCase) CreateNewUser(ctx context.Context, userAuth entity.UserAuth) (entity.User, error) {
    passwordHash, err := password.HashPassword(userAuth.Password)
    if err != nil {
        return entity.User{}, err
    }

    userNew, err := uc.repo.GetUserWithLogin(ctx, userAuth.Login)
    fmt.Println(userNew)
    fmt.Println(err)
    
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
    // TODO: code me pls
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
    return uc.repo.WithinTransaction(ctx, func(txCtx context.Context) error {
        balance, err := uc.repo.GetCurrentBalance(txCtx, userID)
        if err != nil {
            return err
        }

        if withdrawal.Sum.GreaterThan(balance.Current) {
            return ErrPaymentRequired
        }

        withdrawalList, err := uc.repo.GetWithdrawalList(txCtx, userID, withdrawal.Order)
        if err != nil {
            return err
        }
        if withdrawalList != nil {
            return ErrUnprocessableEntity
        }

        err = uc.repo.UpdateBalance(txCtx, userID, withdrawal.Sum)
        if err != nil {
            return err
        }

        err = uc.repo.AddWithdrawal(txCtx, userID, withdrawal.Order, withdrawal.Sum)
        if err != nil {
            return err
        }

        return nil
    })
}


func (uc *GophermartUseCase) GetWithdrawList(ctx context.Context, userID int) ([]entity.Withdraw, error) {
    withdrawalList, err := uc.repo.GetWithdrawalList(ctx, userID, "")
    return withdrawalList, err
}