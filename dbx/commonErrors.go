package dbx

import "errors"

// ErrNoRows is returned by Query.First() or similar methods
// when no rows are found for the given query.
var ErrNoRows = errors.New("no rows in result set")

// Bạn có thể thêm các lỗi chung khác ở đây tùy theo nhu cầu của package dbx.
// Ví dụ:

// ErrDuplicateEntry is returned when an insert or update operation
// violates a unique constraint.
var ErrDuplicateEntry = errors.New("duplicate entry")

// ErrInvalidInput is returned when an operation receives invalid input data.
var ErrInvalidInput = errors.New("invalid input")

// ErrNotFound is a general error indicating a resource was not found.
// Could be used as an alternative to ErrNoRows for broader context.
var ErrNotFound = errors.New("resource not found")

// ErrConnectionFailed is returned when there's an issue establishing or
// maintaining a database connection.
var ErrConnectionFailed = errors.New("database connection failed")

// ErrTransactionFailed is returned when a database transaction fails to commit or rollback.
var ErrTransactionFailed = errors.New("database transaction failed")
