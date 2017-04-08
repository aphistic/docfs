package db

import "errors"

var (
	ErrExists    = errors.New("Entry already exists")
	ErrNotExists = errors.New("Entry does not exist")
)
