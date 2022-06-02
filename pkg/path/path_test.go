package path

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRootPath(t *testing.T) {
	t.Run("get_root_path", func(t *testing.T) {
		rootPath := RootPath()
		assert.NotEmpty(t, rootPath)
	})
}
