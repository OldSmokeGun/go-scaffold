package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomString(t *testing.T) {
	excepts := map[string]map[string]int{
		"length_0":  {"except": 0},
		"length_32": {"except": 32},
		"length_64": {"except": 64},
	}

	for k, v := range excepts {
		t.Run(k, func(t *testing.T) {
			s := RandomString(v["except"])
			assert.Len(t, s, v["except"])
		})
	}
}
