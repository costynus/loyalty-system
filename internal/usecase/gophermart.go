package usecase

import "context"

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