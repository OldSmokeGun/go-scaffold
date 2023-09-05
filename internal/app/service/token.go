package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type AccountTokenData struct {
	UserID int64 `json:"userID"`
}

type AccountTokenClaims struct {
	jwt.RegisteredClaims
	Data *AccountTokenData `json:"data"`
}

type AccountTokenService struct {
	Key    string
	Expire time.Duration
}

func NewAccountTokenService(key string) *AccountTokenService {
	return &AccountTokenService{
		Key: key,
	}
}

func (s *AccountTokenService) Generate(expire time.Duration, data AccountTokenData) (string, error) {
	claims := AccountTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
		},
		Data: &data,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sign, err := token.SignedString([]byte(s.Key))
	if err != nil {
		return "", errors.WithStack(err)
	}

	return sign, nil
}

func (s *AccountTokenService) ValidateToken(token string) (*AccountTokenClaims, error) {
	t, err := jwt.ParseWithClaims(token, &AccountTokenClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.Key), nil
	})

	claim, ok := t.Claims.(*AccountTokenClaims)
	if !t.Valid || !ok {
		return nil, errors.WithStack(err)
	}

	return claim, nil
}

func ParseAccountTokenUnverified(token string) (*AccountTokenClaims, error) {
	t, _, err := jwt.NewParser().ParseUnverified(token, &AccountTokenClaims{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return t.Claims.(*AccountTokenClaims), nil
}
