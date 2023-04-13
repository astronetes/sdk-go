package controller

import (
	"path/filepath"

	"github.com/astronetes/sdk-go/internal/fsys"
	"github.com/astronetes/sdk-go/internal/gotemplate"
	"golang.org/x/exp/maps"
)

type Values interface {
	WithEntry(key string, value interface{}) Values
	WithEntries(entries map[string]interface{}) Values
	Build(path string) (string, error)
}

func NewHelmValues(name string) Values {
	return &values{
		name: name,
		vars: make(map[string]interface{}, 0),
	}
}

type values struct {
	name string
	vars map[string]interface{}
}

func (v *values) WithEntry(key string, value interface{}) Values {
	v.vars[key] = value
	return v
}

func (v *values) WithEntries(entries map[string]interface{}) Values {
	maps.Copy(v.vars, entries)
	return v
}

func (v *values) Build(path string) (string, error) {
	dirPath, filename := filepath.Split(path)
	//TODO: The above line could not work as expected for all the filesystems
	content, err := fsys.GetFileContent(dirPath, filename)
	if err != nil {
		return "", err
	}
	return gotemplate.ApplyTemplateWithVariables(v.name, string(content), v.vars)
}
