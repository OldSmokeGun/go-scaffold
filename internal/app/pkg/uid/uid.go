package uid

import (
	"github.com/bwmarrin/snowflake"
	"math/rand"
	"time"
)

const nodeLimit = 1024

type Uid struct {
	node int64
}

func New(options ...Option) *Uid {
	uid := &Uid{}

	for _, option := range options {
		option(uid)
	}

	return uid
}

type Option func(uid *Uid)

func WithNode(node int64) Option {
	return func(uid *Uid) {
		uid.node = node
	}
}

func (u *Uid) Generate() (int64, error) {
	var node int64

	if u.node == 0 {
		node = int64(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(nodeLimit))
	} else {
		node = u.node
	}

	sn, err := snowflake.NewNode(node)
	if err != nil {
		return 0, err
	}

	return sn.Generate().Int64(), nil
}
