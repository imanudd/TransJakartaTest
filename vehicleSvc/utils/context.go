package utils

import (
	"context"
	"database/sql"
)

var txKey = "tx"

func SetTxCtx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func GetTxContext(ctx context.Context) *sql.Tx {
	raw, ok := ctx.Value("tx").(*sql.Tx)
	if ok {
		return raw
	}
	return nil
}
