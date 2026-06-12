package app

import (
	"net/http"

	"github.com/costynus/loyalty-system/internal/usecase"
	"github.com/costynus/loyalty-system/pkg/logger"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(handler *chi.Mux, uc usecase.Gophermart, l logger.Interface) {
    handler.Use(middleware.Logger)

    // checker
    handler.Get("/healthz", healthzHandler())
    handler.Get("/ping", pingHandler(uc, l))
}


func healthzHandler() http.HandlerFunc{
    return func(w http.ResponseWriter, r *http.Request) {w.WriteHeader(http.StatusOK) }
}

func pingHandler(uc usecase.Gophermart, l logger.Interface) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := uc.PingRepo(r.Context()); err != nil {
            //l.Error(err),
            http.Error(w, "repo error", http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusOK)
    }
}
