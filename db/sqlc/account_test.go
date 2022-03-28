package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yanshen1997/simplebank/util"
)

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	createArgs := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.GetRandomBalance(),
		Currency: util.GetRandomCurrancy(),
	}
	account, err := testQueries.CreateAccount(context.Background(), createArgs)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, account.Balance, createArgs.Balance)
	require.Equal(t, account.Owner, createArgs.Owner)
	require.Equal(t, account.Currency, createArgs.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func deleteAfterTest(t *testing.T, id int64) {
	err := testQueries.DeleteAccount(context.Background(), id)
	require.NoError(t, err)
	accountAgain, err := testQueries.GetAccount(context.Background(), id)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accountAgain)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	accountAgain, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accountAgain)

}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	queryRes, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, queryRes)
	require.Equal(t, account.ID, queryRes.ID)
	require.Equal(t, account.Owner, queryRes.Owner)
	require.Equal(t, account.Balance, queryRes.Balance)
	require.Equal(t, account.Currency, queryRes.Currency)
	require.WithinDuration(t, account.CreatedAt, queryRes.CreatedAt, time.Second)
	deleteAfterTest(t, account.ID)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	args := ListAccountParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccount(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, accounts, int(args.Limit))
	for _, v := range accounts {
		require.NotEmpty(t, v)
	}
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)
	args := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.GetRandomBalance(),
	}
	accountAfterUpdate, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, accountAfterUpdate)
	require.Equal(t, account.ID, accountAfterUpdate.ID)
	require.Equal(t, account.Owner, accountAfterUpdate.Owner)
	require.Equal(t, args.Balance, accountAfterUpdate.Balance)
	require.Equal(t, account.Currency, accountAfterUpdate.Currency)
	require.WithinDuration(t, account.CreatedAt, accountAfterUpdate.CreatedAt, time.Second)
}
