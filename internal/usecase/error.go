package usecase

import "errors"

var ErrNotImplemented = errors.New("not implemented")
var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")
var ErrUnauthorized = errors.New("unauthorized")
