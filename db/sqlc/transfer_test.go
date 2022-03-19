package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	account1, account2 := createRandomAccount(t), createRandomAccount(t)
	createTransferByAccounts(t, account1, account2)
}

func createTransferByAccounts(t *testing.T, account1, account2 Account) Transfer {
	transfer, err := testQueries.CreateTransfer(context.Background(), CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        account1.Balance,
	})

	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, transfer.Amount, account1.Balance)
	require.Equal(t, transfer.FromAccountID, account1.ID)
	require.Equal(t, transfer.ToAccountID, account2.ID)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	return transfer
}

func TestGetTransfer(t *testing.T) {
	account1, account2 := createRandomAccount(t), createRandomAccount(t)
	transfer := createTransferByAccounts(t, account1, account2)
	queryTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, queryTransfer)
	require.Equal(t, transfer.Amount, queryTransfer.Amount)
	require.Equal(t, transfer.FromAccountID, queryTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, queryTransfer.ToAccountID)
	require.Equal(t, transfer.ID, queryTransfer.ID)
	require.WithinDuration(t, transfer.CreatedAt, queryTransfer.CreatedAt, time.Second)
}

func TestListTransfer(t *testing.T) {
	account1, account2 := createRandomAccount(t), createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createTransferByAccounts(t, account1, account2)
	}
	transfers, err := testQueries.ListTransfer(context.Background(), ListTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Limit:         5,
		Offset:        5,
	})
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, v := range transfers {
		require.NotEmpty(t, v)
		require.True(t, account1.ID == v.FromAccountID && account2.ID == v.ToAccountID)
	}

}
