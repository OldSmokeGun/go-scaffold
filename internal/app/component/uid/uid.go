package uid

//go:generate mockgen -source=uid.go -destination=uid_mock.go -package=uid -mock_names=Generator=MockUid

import (
	"github.com/bwmarrin/snowflake"
	"math/rand"
	"time"
)

var defaultRand = rand.New(rand.NewSource(time.Now().UnixNano()))

const maxNode = 1<<10 - 1

type Generator interface {
	Generate(options ...Option) (string, error)
}

type Uid struct {
	node int64
	rand *rand.Rand
}

func New() *Uid {
	return &Uid{rand: defaultRand}
}

type Option func(uid *Uid)

func WithNode(node int64) Option {
	return func(uid *Uid) {
		uid.node = node
	}
}

func WithRand(rand *rand.Rand) Option {
	return func(uid *Uid) {
		uid.rand = rand
	}
}

func (u *Uid) Generate(options ...Option) (string, error) {
	for _, option := range options {
		option(u)
	}

	var node int64

	if u.node == 0 {
		node = int64(u.rand.Intn(maxNode + 1))
	} else {
		node = u.node
	}

	sn, err := snowflake.NewNode(node)
	if err != nil {
		return "", err
	}

	return sn.Generate().String(), nil
}
