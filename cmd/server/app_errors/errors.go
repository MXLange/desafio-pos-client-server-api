package apperrors

import "errors"

var (
	ErrNilDB error = errors.New("db cannot be nil")
	ErrNilPriceRepository error = errors.New("price repository cannot be nil")
)