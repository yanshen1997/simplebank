package db

/**
 *	事务实现转账功能
 */

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	tQ := New(tx)
	err = fn(tQ)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rockback error: err: %v , rberr: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParam struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
}

func (store *Store) TransferTx(ctx context.Context, tp TransferTxParam) (TransferTxResult, error) {
	var res TransferTxResult

	store.execTx(ctx, func(q *Queries) error {
		var err error
		res.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: tp.FromAccountID,
			ToAccountID:   tp.ToAccountID,
			Amount:        tp.Amount,
		})

		if err != nil {
			return err
		}

		res.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: tp.FromAccountID,
			Amount:    -tp.Amount,
		})

		if err != nil {
			return err
		}

		res.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: tp.ToAccountID,
			Amount:    tp.Amount,
		})

		if err != nil {
			return err
		}

		res.FromAccount, err = store.AddAccountBalance(ctx, AddAccountBalanceParams{
			Amount: -tp.Amount,
			ID:     tp.FromAccountID,
		})
		if err != nil {
			return err
		}

		res.ToAccount, err = store.AddAccountBalance(ctx, AddAccountBalanceParams{
			Amount: tp.Amount,
			ID:     tp.ToAccountID,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return res, nil
}
