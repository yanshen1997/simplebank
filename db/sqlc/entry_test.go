package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createEntryByAccount(t, account)
}

func createEntryByAccount(t *testing.T, account Account) Entry {
	entry, err := testQueries.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: account.ID,
		Amount:    account.Balance,
	})

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, account.ID, entry.AccountID)
	require.Equal(t, account.Balance, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry := createEntryByAccount(t, account)
	queryEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, queryEntry)
	require.Equal(t, queryEntry.AccountID, entry.AccountID)
	require.Equal(t, queryEntry.ID, entry.ID)
	require.Equal(t, queryEntry.Amount, entry.Amount)
	require.WithinDuration(t, entry.CreatedAt, queryEntry.CreatedAt, time.Second)
}

func TestListEntry(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createEntryByAccount(t, account)
	}
	entries, err := testQueries.ListEntry(context.Background(), ListEntryParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	})

	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, v := range entries {
		require.NotEmpty(t, v)
		require.Equal(t, v.AccountID, account.ID)
	}

}
