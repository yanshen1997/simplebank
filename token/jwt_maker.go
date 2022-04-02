package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const MinJWTSecretKeyLen = 32

type JwtMaker struct {
	secretKey string
}

func NewJwtMaker(secretKey string) (Maker, error) {
	if len(secretKey) < MinJWTSecretKeyLen {
		return nil, fmt.Errorf("secret key length smaller than minmum length %d", MinJWTSecretKeyLen)
	}
	return &JwtMaker{
		secretKey: secretKey,
	}, nil
}

func (maker *JwtMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

// VerifyToken 验证token并返回payload
func (maker *JwtMaker) VerifyToken(token string) (*Payload, error) {
	var keyFunc func(*jwt.Token) (interface{}, error)
	keyFunc = func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		if verr, ok := err.(*jwt.ValidationError); ok {
			if errors.Is(verr.Inner, ErrExpiredToken) {
				return nil, ErrExpiredToken
			}
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil

}
