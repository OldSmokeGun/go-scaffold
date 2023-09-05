package domain

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/google/uuid"
)

type User struct {
	ID       int64    `json:"id"`
	Username string   `json:"username"`
	Password Password `json:"password"`
	Nickname string   `json:"nickname"`
	Phone    string   `json:"phone"`
	Salt     string   `json:"salt"`
}

func (u *User) RefreshSalt() {
	u.Salt = uuid.New().String()
}

func (u *User) ToProfile() *UserProfile {
	return &UserProfile{
		ID:       u.ID,
		Username: u.Username,
		Nickname: u.Nickname,
		Phone:    u.Phone,
	}
}

type Plaintext string

func (p Plaintext) Encrypt() Password {
	sum := md5.Sum([]byte(p))
	return Password(hex.EncodeToString(sum[:]))
}

type Password string

func (p Password) Verify(s Plaintext) bool {
	return p == s.Encrypt()
}

type UserProfile struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}
