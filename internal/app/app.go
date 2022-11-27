package app

import (
	"fmt"
	"net/http"

	"github.com/costynus/loyalty-system/config"
	"github.com/costynus/loyalty-system/internal/usecase"
	"github.com/costynus/loyalty-system/internal/usecase/repo"
	"github.com/costynus/loyalty-system/internal/usecase/webapi"
	"github.com/costynus/loyalty-system/pkg/logger"
	"github.com/costynus/loyalty-system/pkg/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	decimal.MarshalJSONWithoutQuotes = cfg.App.MarshalJSONWithoutQuotes

	err := migration(cfg.PG.URL, cfg.PG.MigDir)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - migration: %w", err))
	}
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PollMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	client := resty.New().SetBaseURL(cfg.App.AccrualSystemAddress)
	webAPI := webapi.New(client)

	uc := usecase.New(
		repo.New(pg),
		webAPI,
		cfg.App.WorkersCount,
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
