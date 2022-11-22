package usecase

import (
	"context"

	"github.com/costynus/loyalty-system/internal/entity"
)

type (
    Gophermart interface {
        PingRepo(context.Context) error
        CreateNewUser(context.Context, entity.UserAuth) error
        CheckUser(context.Context, entity.UserAuth) error
    }

    GophermartRepo interface {
        Ping(context.Context) error
    }
)