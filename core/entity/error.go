package entity

import "errors"

var (
	ErrorBind         = errors.New("bind error")
	ErrorBadRequest   = errors.New("bad request")
	ErrorNotFound     = errors.New("not found")
	ErrorInternal     = errors.New("internal error")
	ErrorForbidden    = errors.New("forbidden")
	ErrorUnauthorized = errors.New("unauthorized")
	ErrorConflict     = errors.New("conflict")
)
