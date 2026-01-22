package db

import (
	"context"
	"database/sql"

	"github.com/uoul/go-async"
)

// TransactionScopeFunction is a function type that executes database operations within a transaction context.
//
// This function type is designed to encapsulate business logic that requires transactional guarantees.
// It receives a context for cancellation/timeout control and a transaction object for executing
// database operations, and returns a result of type T along with any error that occurred.
//
// Type parameter T represents the return type of the transaction function, allowing for
// flexible return values based on the specific use case.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control during transaction execution
//   - tx: Active database transaction to use for executing queries and commands
//
// Returns:
//   - T: Result of the transaction operations (type determined by caller)
//   - error: Non-nil if any operation within the transaction fails
type TransactionScopeFunction[T any] func(ctx context.Context, tx *sql.Tx) (T, error)

// ExecuteInTransaction executes the provided function within a database transaction.
//
// This function creates a new transaction using the provided database connection and
// executes the given TransactionScopeFunction within that transaction context. If the
// function completes successfully, the transaction is committed; otherwise, it is rolled back.
// The transaction is also rolled back if a panic occurs during execution (via deferred rollback).
//
// Type parameter T represents the return type of the transaction function, allowing for
// flexible return values based on the specific business logic requirements.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control, propagated to the transaction
//   - db: Database connection to use for creating the transaction
//   - tsf: Function to execute within the transaction scope
//   - opts: Optional transaction options (isolation level, read-only mode, etc.).
//     If not provided, default transaction options are used.
//
// Returns:
//   - T: The result returned by the transaction function
//   - error: Non-nil if transaction creation, execution, or commit fails
func ExecuteInTransaction[T any](ctx context.Context, db IDbConnection, tsf TransactionScopeFunction[T], opts ...sql.TxOptions) (T, error) {
	var txOpts *sql.TxOptions = nil
	if len(opts) > 0 {
		txOpts = &opts[0]
	}
	// Create transaction
	tx, err := db.BeginTx(ctx, txOpts)
	if err != nil {
		return *new(T), err
	}
	defer tx.Rollback()
	// Execute TransactionScopeFunction
	r, err := tsf(ctx, tx)
	if err != nil {
		return *new(T), err
	}
	// Commit changes
	if err := tx.Commit(); err != nil {
		return *new(T), err
	}
	// Return result
	return r, nil
}

// ExecuteInTransactionAsync executes a database transaction asynchronously and returns the result.
//
// This function wraps the synchronous ExecuteInTransaction function in an asynchronous execution
// context, allowing the transaction to run concurrently without blocking the caller. The result
// is returned as an async.Result that can be awaited or processed later.
//
// The function leverages Go's concurrency model to execute the entire transaction lifecycle
// (begin, execute, commit/rollback) in a separate goroutine, making it suitable for scenarios where:
// - You want to execute multiple independent transactions in parallel
// - You need to avoid blocking the main execution flow while waiting for transaction completion
// - You're implementing non-blocking data persistence patterns
// - You want to improve throughput by overlapping transaction execution with other work
//
// Type parameter T represents the return type of the transaction function, allowing for
// flexible return values based on the specific business logic requirements.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control, propagated to the underlying transaction
//   - db: Database connection to use for creating the transaction
//   - tsf: Function to execute within the transaction scope
//   - opts: Optional transaction options (isolation level, read-only mode, etc.).
//     If not provided, default transaction options are used.
//
// Returns:
//   - async.Result[T]: An async result object containing either:
//   - The result returned by the transaction function
//   - An error if transaction creation, execution, or commit fails
func ExecuteInTransactionAsync[T any](ctx context.Context, db IDbConnection, tsf TransactionScopeFunction[T], opts ...sql.TxOptions) async.Result[T] {
	return async.Do(
		ctx,
		func(ctx context.Context) (T, error) {
			return ExecuteInTransaction(ctx, db, tsf, opts...)
		},
	)
}
