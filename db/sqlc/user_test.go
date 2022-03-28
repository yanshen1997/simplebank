package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yanshen1997/simplebank/util"
)

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func createRandomUser(t *testing.T) User {
	createArgs := CreateUserParams{
		Username:       util.GetRandomOwner(),
		HashedPassword: "to be improved",
		FullName:       util.GetRandomOwner(),
		Email:          util.GetRandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), createArgs)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user.Username, createArgs.Username)
	require.Equal(t, user.HashedPassword, createArgs.HashedPassword)
	require.Equal(t, user.FullName, createArgs.FullName)
	require.Equal(t, user.Email, createArgs.Email)

	require.NotZero(t, user.PasswordChangedAt)
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	queryRes, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, queryRes)
	require.Equal(t, user.Username, queryRes.Username)
	require.Equal(t, user.HashedPassword, queryRes.HashedPassword)
	require.Equal(t, user.FullName, queryRes.FullName)
	require.Equal(t, user.Email, queryRes.Email)
	require.WithinDuration(t, user.CreatedAt, queryRes.CreatedAt, time.Second)
	require.WithinDuration(t, user.PasswordChangedAt, queryRes.PasswordChangedAt, time.Second)
}
