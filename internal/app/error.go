package app

import (
	"errors"
	"net/http"

	"github.com/costynus/loyalty-system/internal/usecase"
)

func errorHandler(w http.ResponseWriter, err error) {
	if errors.Is(err, usecase.ErrNotImplemented) {
		http.Error(w, err.Error(), http.StatusNotImplemented)
	} else if errors.Is(err, usecase.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
    } else if errors.Is(err, usecase.ErrConflict) {
        http.Error(w, err.Error(), http.StatusConflict)
    } else if errors.Is(err, usecase.ErrUnauthorized) {
        http.Error(w, err.Error(), http.StatusUnauthorized)
	} else {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
