package app

import (
	"fmt"
	"net/http"

	"github.com/costynus/loyalty-system/config"
	"github.com/costynus/loyalty-system/internal/usecase"
	"github.com/costynus/loyalty-system/internal/usecase/repo"
	"github.com/costynus/loyalty-system/pkg/logger"
	"github.com/costynus/loyalty-system/pkg/postgres"
    "github.com/go-chi/chi/v5"
)


func Run(cfg *config.Config){
    l := logger.New(cfg.Log.Level)

    pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PollMax))
    if err != nil {
        l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
    }
    defer pg.Close()

    uc := usecase.New(
        repo.New(pg),
    )

    handler := chi.NewRouter()
    NewRouter(
        handler,
        uc,
        l,
    )

    l.Info("app - Run")

    l.Fatal(http.ListenAndServe(cfg.HTTP.Address, handler))
}
