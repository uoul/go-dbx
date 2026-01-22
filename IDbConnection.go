package db

import (
	"context"
	"database/sql"
)

type IDbConnection interface {
	IDbSession
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}
