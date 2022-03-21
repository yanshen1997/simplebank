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

		if tp.FromAccountID < tp.ToAccountID {
			res.FromAccount, res.ToAccount, err = addMoney(ctx, store.Queries, tp.FromAccountID, -tp.Amount, tp.ToAccountID, tp.Amount)
		} else {
			res.ToAccount, res.FromAccount, err = addMoney(ctx, store.Queries, tp.ToAccountID, tp.Amount, tp.FromAccountID, -tp.Amount)
		}
		if err != nil {
			return err
		}
		return nil
	})

	return res, nil
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
