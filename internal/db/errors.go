package db

import (
	"errors"
)

var (
	ErrNotFound = errors.New("report not found")
	ErrConflict = errors.New("report conflict")
)
