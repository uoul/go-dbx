package db

import (
	"context"
	"database/sql"
)

type IDbSession interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}
