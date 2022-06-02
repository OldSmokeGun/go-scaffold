package file

import (
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

type Config struct {
	Path string
}

// New casin file adapter
func New(config *Config) *fileadapter.FilteredAdapter {
	return fileadapter.NewFilteredAdapter(config.Path)
}
