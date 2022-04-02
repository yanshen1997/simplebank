package token

import "time"

type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)
	// VerifyToken 验证token并返回payload
	VerifyToken(token string) (*Payload, error)
}
