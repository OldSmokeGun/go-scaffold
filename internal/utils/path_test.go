package utils

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestIsDir(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	testFile, err := os.Executable()
	require.NoError(t, err)

	excepts := map[string]map[string]interface{}{
		"is_dir(rel)":     {"path": "../utils", "except": true},
		"is_not_dir(rel)": {"path": "../internal.go", "except": false},
		"is_dir(abs)":     {"path": wd, "except": true},
		"is_not_dir(abs)": {"path": testFile, "except": false},
	}

	for k, v := range excepts {
		t.Run(k, func(t *testing.T) {
			ok, err := IsDir(v["path"].(string))
			require.NoError(t, err)
			assert.Equal(t, v["except"], ok)
		})
	}
}

func TestPathExist(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	excepts := map[string]map[string]interface{}{
		"rel_is_exist":     {"path": "../utils", "except": true},
		"rel_is_not_exist": {"path": "../not_exist.go", "except": false},
		"abs_is_exist":     {"path": wd, "except": true},
		"abs_is_not_exist": {"path": "/a/b/c/d/e/f", "except": false},
	}

	for k, v := range excepts {
		t.Run(k, func(t *testing.T) {
			ok := PathExist(v["path"].(string))
			assert.Equal(t, v["except"], ok)
		})
	}
}
