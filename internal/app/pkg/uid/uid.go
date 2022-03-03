package uid

import (
	"github.com/bwmarrin/snowflake"
	"math/rand"
	"time"
)

var defaultRand = rand.New(rand.NewSource(time.Now().UnixNano()))

const maxNode = 1<<10 - 1

type uid struct {
	node int64
	rand *rand.Rand
}

type Option func(uid *uid)

func WithNode(node int64) Option {
	return func(uid *uid) {
		uid.node = node
	}
}

func WithRand(rand *rand.Rand) Option {
	return func(uid *uid) {
		uid.rand = rand
	}
}

func Generate(options ...Option) (int64, error) {
	u := &uid{rand: defaultRand}

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
		return 0, err
	}

	return sn.Generate().Int64(), nil
}
