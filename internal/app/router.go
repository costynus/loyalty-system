package app

import (
	"encoding/json"
	"net/http"

	"github.com/costynus/loyalty-system/internal/entity"
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

    // auth
    handler.Route("/api/user", func(r chi.Router) {
        r.Post("/register", registrationUser(uc, l))
        r.Post("/login", loginUser(uc, l))
    })
}

func registrationUser(uc usecase.Gophermart, l logger.Interface) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var userAuth entity.UserAuth

        if err := json.NewDecoder(r.Body).Decode(&userAuth); err != nil {
            http.Error(w, "bad request", http.StatusBadRequest)
            return
        }

        if err := uc.CreateNewUser(r.Context(), userAuth); err != nil {
            errorHandler(w, err)
            return
        }

        // TODO: generate jwtToken
        // TODO: set jwtToken

        w.WriteHeader(http.StatusOK)
    }
}

func loginUser(uc usecase.Gophermart, l logger.Interface) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        var userAuth entity.UserAuth

        if err := json.NewDecoder(r.Body).Decode(&userAuth); err != nil {
            http.Error(w, "bad request", http.StatusBadRequest)
            return
        }

        if err := uc.CheckUser(r.Context(), userAuth); err != nil {
            errorHandler(w, err)
            return
        }

        // TODO: generate jwtToken
        // TODO: set jwtToken

        w.WriteHeader(http.StatusOK)
    }
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
