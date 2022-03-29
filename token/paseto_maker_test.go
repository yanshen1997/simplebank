package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yanshen1997/simplebank/util"
	"golang.org/x/crypto/chacha20poly1305"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(chacha20poly1305.KeySize))
	require.NoError(t, err)
	username := util.GetRandomOwner()
	duration := time.Now()
	token, err := maker.CreateToken(username, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, duration, payload.IssuedAt, time.Second)
	require.WithinDuration(t, duration.Add(time.Minute), payload.ExpiredAt, time.Second)
}
