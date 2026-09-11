package token

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type Claims struct {
	UserID    int64
	Timestamp int64
	Sign      string
}

func Generate(userID, ts int64, salt string) string {
	var (
		id  = strconv.FormatInt(userID, 10)
		tss = strconv.FormatInt(ts, 10)
	)

	sum := md5.Sum([]byte(fmt.Sprintf("%s%s%s", id, tss, salt))) //nolint:gosec // 与旧系统 token 签名兼容。
	sign := hex.EncodeToString(sum[:])

	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s|%s|%s", id, tss, sign)))
}

func Parse(token string) (*Claims, error) {
	raw, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(raw), "|")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	userID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, err
	}
	timestamp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, err
	}
	sign := parts[2]

	return &Claims{
		UserID:    userID,
		Timestamp: timestamp,
		Sign:      sign,
	}, nil
}

func Verify(token, sign string) (bool, error) {
	claims, err := Parse(token)
	if err != nil {
		return false, err
	}

	vs := Generate(claims.UserID, claims.Timestamp, claims.Sign)

	return vs == sign, nil
}
