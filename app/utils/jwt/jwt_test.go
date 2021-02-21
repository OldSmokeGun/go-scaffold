package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
	"time"
)

var TestKey = "123456"

func TestNewToken(t *testing.T) {
	excepts := map[string]map[string]string{
		"without_key": {"key": ""},
		"with_key":    {"key": TestKey},
	}

	var (
		token *Token
		err   error
	)
	for k, v := range excepts {
		t.Run(k, func(t *testing.T) {
			token, err = NewToken(WithKey(v["key"]))
			if k == "without_key" {
				if assert.ErrorIs(t, err, ErrMissingKey) {
					assert.Nil(t, token)
				}
			} else if k == "with_key" {
				if assert.NoError(t, err) {
					assert.NotNil(t, token)
				}
			}
		})
	}
}

func TestTokenCreate(t *testing.T) {
	exceptToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2R5Ijp7Im5hbWUiOiJUb20iLCJzZXgiOiJtYWxlIn0sImF1ZCI6InNlcnZlciIsImV4cCI6NTA4Mzg4NTY5NSwianRpIjoiYWJjZGUiLCJpYXQiOjc4ODkxODQwMCwiaXNzIjoiYXBwIiwibmJmIjo3ODg5MTg0MDAsInN1YiI6ImF1dGhlbnRpY2F0aW9uIHRva2VuIn0.KFq7Nea4aJEPuBr_sEGSHZKbvJheqb8uUp4ilPwdP7o"

	now, err := time.Parse("2006-01-02 15:04:05", "1995-01-01 00:00:00")
	require.NoError(t, err)

	token, err := NewToken(
		WithKey(TestKey),
		WithAlg(DefaultAlg),
		WithBody(map[string]interface{}{
			"name": "Tom",
			"sex":  "male",
		}),
		WithClaims(jwt.StandardClaims{
			Audience:  "server",
			ExpiresAt: now.Add(time.Second * time.Duration(math.MaxUint32)).Unix(),
			Id:        "abcde",
			IssuedAt:  now.Unix(),
			Issuer:    "app",
			NotBefore: now.Unix(),
			Subject:   "authentication token",
		}),
	)
	require.NoError(t, err)

	signedToken, err := token.Create()
	require.NoError(t, err)

	assert.Equal(t, exceptToken, signedToken)
}

func TestTokenParse(t *testing.T) {
	signedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2R5Ijp7Im5hbWUiOiJUb20iLCJzZXgiOiJtYWxlIn0sImF1ZCI6InNlcnZlciIsImV4cCI6NTA4Mzg4NTY5NSwianRpIjoiYWJjZGUiLCJpYXQiOjc4ODkxODQwMCwiaXNzIjoiYXBwIiwibmJmIjo3ODg5MTg0MDAsInN1YiI6ImF1dGhlbnRpY2F0aW9uIHRva2VuIn0.KFq7Nea4aJEPuBr_sEGSHZKbvJheqb8uUp4ilPwdP7o"

	token, err := NewToken(
		WithKey(TestKey),
	)
	require.NoError(t, err)

	parsedToken, claims, err := token.Parse(signedToken)
	require.NoError(t, err)

	if assert.NotNil(t, parsedToken) {
		if assert.True(t, parsedToken.Valid) {
			assert.Equal(t, "Tom", claims.Body["name"])
			assert.Equal(t, "male", claims.Body["sex"])
		}
	}
}

func TestTokenExpire(t *testing.T) {
	now, err := time.Parse("2006-01-02 15:04:05", "1995-01-01 00:00:00")
	require.NoError(t, err)

	token, err := NewToken(
		WithKey(TestKey),
		WithAlg(DefaultAlg),
		WithBody(map[string]interface{}{
			"name": "Tom",
			"sex":  "male",
		}),
		WithClaims(jwt.StandardClaims{
			Audience:  "server",
			ExpiresAt: now.Unix(),
			Id:        "abcde",
			IssuedAt:  now.Unix(),
			Issuer:    "app",
			NotBefore: now.Unix(),
			Subject:   "authentication token",
		}),
	)
	require.NoError(t, err)

	signedToken, err := token.Create()
	require.NoError(t, err)

	token, err = NewToken(
		WithKey(TestKey),
	)
	require.NoError(t, err)

	parsedToken, _, err := token.Parse(signedToken)

	assert.ErrorIs(t, err, ErrExpired)
	if assert.NotNil(t, parsedToken) {
		assert.False(t, parsedToken.Valid)
	}
}

func TestTokenNotValidYet(t *testing.T) {
	token, err := NewToken(
		WithKey(TestKey),
		WithAlg(DefaultAlg),
		WithBody(map[string]interface{}{
			"name": "Tom",
			"sex":  "male",
		}),
		WithClaims(jwt.StandardClaims{
			Audience:  "server",
			ExpiresAt: time.Now().Add(time.Second * time.Duration(math.MaxUint32)).Unix(),
			Id:        "abcde",
			IssuedAt:  time.Now().Unix(),
			Issuer:    "app",
			NotBefore: time.Now().Add(time.Second * 3600).Unix(),
			Subject:   "authentication token",
		}),
	)
	require.NoError(t, err)

	signedToken, err := token.Create()
	require.NoError(t, err)

	token, err = NewToken(
		WithKey(TestKey),
	)
	require.NoError(t, err)

	parsedToken, _, err := token.Parse(signedToken)

	assert.ErrorIs(t, err, ErrNotValidYet)
	if assert.NotNil(t, parsedToken) {
		assert.False(t, parsedToken.Valid)
	}
}
