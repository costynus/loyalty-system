package usecase

import "context"

type (
    Gophermart interface {
        PingRepo(context.Context) error
    }

    GophermartRepo interface {
        Ping(context.Context) error
    }
)