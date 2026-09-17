package authtoken

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Generate issue token
func Generate(userID int64, salt string, ts int64) string {
	idStr := strconv.FormatInt(userID, 10)
	tsStr := strconv.FormatInt(ts, 10)
	sum := md5.Sum([]byte(idStr + tsStr + salt))
	sign := hex.EncodeToString(sum[:])

	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s|%s|%s", idStr, tsStr, sign)))
}

// Parse token
func Parse(token string) (userID int64, ts int64, sign string, err error) {
	raw, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return 0, 0, "", err
	}
	parts := strings.Split(string(raw), "|")
	if len(parts) != 3 {
		return 0, 0, "", errors.New("invalid auth token")
	}
	userID, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, "", err
	}
	ts, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, "", err
	}

	return userID, ts, parts[2], nil
}

// VerifySign verify signature with salt
func VerifySign(userID, ts int64, salt, sign string) bool {
	idStr := strconv.FormatInt(userID, 10)
	tsStr := strconv.FormatInt(ts, 10)
	sum := md5.Sum([]byte(idStr + tsStr + salt))

	return hex.EncodeToString(sum[:]) == sign
}

// InExpireWindow determine whether the issuance time is within the valid window
func InExpireWindow(ts int64, now time.Time, expire time.Duration) bool {
	n := now.Unix()
	if n < ts {
		return false
	}
	if n > ts+int64(expire.Seconds()) {
		return false
	}

	return true
}
