package uid

import (
	"fmt"
	"strconv"

	"github.com/sony/sonyflake"
)

type Generator interface {
	Generate() (string, error)
}

type UID struct {
	sf *sonyflake.Sonyflake // sonyflake 实例
}

func New() (*UID, error) {
	sf, err := sonyflake.New(sonyflake.Settings{})
	if err != nil {
		return nil, fmt.Errorf("new sony flake error: %w", err)
	}

	return &UID{sf: sf}, nil
}

func (u *UID) Generate() (string, error) {
	id, err := u.sf.NextID()
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(id, 10), nil
}
