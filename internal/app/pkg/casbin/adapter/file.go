package adapter

import (
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

// NewFileAdapter build casin file adapter
func NewFileAdapter(fp string) *fileadapter.FilteredAdapter {
	return fileadapter.NewFilteredAdapter(fp)
}
