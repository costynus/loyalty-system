package app

import (
	"errors"
	"net/http"

	"github.com/costynus/loyalty-system/internal/usecase"
)

func errorHandler(w http.ResponseWriter, err error) {
	if errors.Is(err, usecase.ErrNotImplemented) {
		http.Error(w, err.Error(), http.StatusNotImplemented) // 501
	} else if errors.Is(err, usecase.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound) // 404
    } else if errors.Is(err, usecase.ErrConflict) {
        http.Error(w, err.Error(), http.StatusConflict) // 409
    } else if errors.Is(err, usecase.ErrUnauthorized) {
        http.Error(w, err.Error(), http.StatusUnauthorized) // 401
    } else if errors.Is(err, usecase.ErrPaymentRequired) {
        http.Error(w, err.Error(), http.StatusPaymentRequired) // 402
    } else if errors.Is(err, usecase.ErrUnprocessableEntity) {
        http.Error(w, err.Error(), http.StatusUnprocessableEntity) // 422
	} else {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
