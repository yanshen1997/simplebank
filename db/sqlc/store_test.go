package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	account1, account2 := createRandomAccount(t), createRandomAccount(t)
	store := NewStore(testDb)
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParam{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.CreatedAt)
		require.NotZero(t, transfer.ID)
		_, err = testQueries.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry, toEntry := result.FromEntry, result.ToEntry
		require.NotEmpty(t, fromEntry)
		require.NotEmpty(t, toEntry)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.Equal(t, toEntry.Amount, amount)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//account test here
		fromAccount, toAccount := result.FromAccount, result.ToAccount
		require.NotEmpty(t, fromAccount)
		require.NotEmpty(t, toAccount)
		require.NotZero(t, fromAccount.CreatedAt)
		require.NotZero(t, toAccount.CreatedAt)
		require.Equal(t, fromAccount.ID, account1.ID)
		require.Equal(t, toAccount.ID, account2.ID)
		diff1, diff2 := account1.Balance-fromAccount.Balance, toAccount.Balance-account2.Balance
		require.True(t, diff1 == diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

	}

	updatedFromAccount, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedFromAccount)
	require.Equal(t, updatedFromAccount.ID, account1.ID)
	updatedToAccount, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedToAccount)
	require.Equal(t, updatedToAccount.ID, account2.ID)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedFromAccount.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedToAccount.Balance)

}
func TestDLTransferTx(t *testing.T) {
	account1, account2 := createRandomAccount(t), createRandomAccount(t)
	store := NewStore(testDb)
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccount, toAccount := account1, account2
		if i%2 == 0 {
			fromAccount = account2
			toAccount = account1
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParam{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedFromAccount, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedFromAccount)
	require.Equal(t, updatedFromAccount.ID, account1.ID)
	updatedToAccount, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedToAccount)
	require.Equal(t, updatedToAccount.ID, account2.ID)

	require.Equal(t, account1.Balance, updatedFromAccount.Balance)
	require.Equal(t, account2.Balance, updatedToAccount.Balance)
}
