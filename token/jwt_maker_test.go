package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"github.com/yanshen1997/simplebank/util"
)

func TestJwtMaker(t *testing.T) {
	maker, err := NewJwtMaker(util.RandomString(32))
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

func TestExpiredTocken(t *testing.T) {
	maker, err := NewJwtMaker(util.RandomString(32))
	require.NoError(t, err)
	username := util.GetRandomOwner()
	token, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidTocken(t *testing.T) {
	payload, err := NewPayload(util.GetRandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJwtMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
