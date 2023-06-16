package models

import "errors"

var (
	ErrEmailTaken = errors.New("models: email address already in use")
	ErrNoData     = errors.New("create sql: no rows in result set")
)
