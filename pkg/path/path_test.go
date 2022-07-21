package path

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootPath(t *testing.T) {
	t.Run("get_root_path", func(t *testing.T) {
		rootPath := RootPath()
		assert.NotEmpty(t, rootPath)
	})
}
