package util

import "errors"

// Error constants
const (
	ErrStrSourceNotFound = "source not found"
	ErrStrInvalidType    = "invalid type"
	ErrStrUserNotFound   = "user not found"
	// Add other error constants here
)

var (
	ErrCategoryNotFound   = errors.New("category not found")
	ErrInvalidInput       = errors.New("invalid input data")
	ErrDatabase           = errors.New("database error")
	ErrCategoryNameLength = errors.New("category name length exceeds limit")
	ErrCategoryDescLength = errors.New("category description length exceeds limit")
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrInvalidFilter       = errors.New("invalid filter provided")
)

var (
	ErrSourceNotFound = errors.New("source not found")
	ErrInvalidType    = errors.New("invalid source type; only 'CREDIT' and 'SAVINGS' are allowed")
)
