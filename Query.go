package db

import (
	"context"

	"github.com/uoul/go-async"
)

// Query executes a SQL query and returns the results as a slice of type T.
//
// The function performs the following operations:
// 1. Executes the provided SQL query with the given arguments using the database session
// 2. Parses the returned rows into a slice of the specified type T
// 3. Ensures proper resource cleanup by closing the rows
//
// Type parameter T must be a type that can be populated from database rows.
// The actual mapping from database rows to type T is handled by parseDbResult.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - conn: Database session (connection or transaction) to execute the query on
//   - query: SQL query string to execute
//   - args: Variadic arguments to be used as query parameters (prevents SQL injection)
//
// Returns:
//   - []T: Slice of results parsed from the query, empty slice if no rows match
//   - error: Non-nil if query execution or result parsing fails
func Query[T any](ctx context.Context, conn IDbSession, query string, args ...any) ([]T, error) {
	rows, err := conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result, err := parseDbResult[T](rows)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// QueryAsync executes a SQL query asynchronously and returns the results as a slice of type T.
//
// This function wraps the synchronous Query function in an asynchronous execution context,
// allowing the query to run concurrently without blocking the caller. The result is returned
// as an async.Result that can be awaited or processed later.
//
// The function leverages Go's concurrency model to execute the database query in a separate
// goroutine, making it suitable for scenarios where you want to:
// - Execute multiple independent queries in parallel
// - Avoid blocking the main execution flow while waiting for database results
// - Implement non-blocking data fetching patterns
func QueryAsync[T any](ctx context.Context, conn IDbSession, query string, args ...any) async.Result[[]T] {
	return async.Do(
		ctx,
		func(ctx context.Context) ([]T, error) {
			return Query[T](ctx, conn, query, args...)
		},
	)
}
