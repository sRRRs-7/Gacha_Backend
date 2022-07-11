package db

import (
	"database/sql"
)

type Store interface {
	Querier
	// TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore {
		Queries: New(db),
		db: db,
	}
}

