package vfs

import (
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
)

//go:generate go-assets-builder ../../app/templates -o file_system.go -p vfs -v filesystem

func LoadTemplatesFromFilesystem() (*template.Template, error) {
	t := template.New("")

	includeExtensions := map[string]struct{}{
		".tmpl": {}, ".tpl": {}, ".html": {},
	}

	for name, file := range filesystem.Files {
		if file.IsDir() {
			continue
		}

		if _, ok := includeExtensions[filepath.Ext(name)]; !ok {
			continue
		}

		h, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		t, err = t.New(strings.SplitAfter(name, "templates/")[1]).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}
