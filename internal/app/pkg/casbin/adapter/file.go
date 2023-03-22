package adapter

import (
	"go-scaffold/internal/config"

	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

// NewFileAdapter build casin file adapter
func NewFileAdapter(conf config.CasbinFileAdapter) *fileadapter.FilteredAdapter {
	return fileadapter.NewFilteredAdapter(conf.Path)
}
