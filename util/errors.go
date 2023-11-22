/*
MIT License

Copyright (c) 2023 Narayan Babu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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
