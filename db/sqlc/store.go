package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Store interface {
	Querier
	ExchangeTx(ctx context.Context, arg ExchangeTxParams) (ExchangeTxResult, error)
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

// execTx executes a function within a database transaction
func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("txErr: %v, rbErr: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type ExchangeTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	ItemID1       int64 `json:"item_id_1"`
	ItemID2       int64 `json:"item_id_2"`
}

type ExchangeTxResult struct {
	Exchange1    Exchange `json:"exchange_1"`
	Exchange2    Exchange `json:"exchange_2"`
	Gallery1	 Gallery  `json:"gallery_1"`
	Gallery2	 Gallery  `json:"gallery_2"`
}

func (s *SQLStore) ExchangeTx(ctx context.Context, arg ExchangeTxParams) (ExchangeTxResult, error) {
	var result ExchangeTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		result.Exchange1, err = q.CreateExchange(ctx, CreateExchangeParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			ItemID:        arg.ItemID1,
		})
		if err != nil {
			return err
		}

		result.Exchange2, err = q.CreateExchange(ctx, CreateExchangeParams{
			FromAccountID: arg.ToAccountID,
			ToAccountID:   arg.FromAccountID,
			ItemID:        arg.ItemID2,
		})
		if err != nil {
			return err
		}

		result.Gallery1, err = q.UpdateGallery(ctx, UpdateGalleryParams{
			OwnerID: arg.FromAccountID,
			ItemID: arg.ItemID1,
			OwnerID_2: arg.ToAccountID,
			ExchangeAt: time.Now(),
		})
		if err != nil {
			return err
		}

		result.Gallery2, err = q.UpdateGallery(ctx, UpdateGalleryParams{
			OwnerID: arg.ToAccountID,
			ItemID: arg.ItemID2,
			OwnerID_2: arg.FromAccountID,
			ExchangeAt: time.Now(),
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

