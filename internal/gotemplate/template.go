package gotemplate

import (
	"bytes"
	"fmt"
	"text/template"
)

func ApplyTemplateWithVariables(name string, text string, vars map[string]interface{}) (string, error) {
	t, err := template.New(name).Parse(string(text))
	if err != nil {
		return "", fmt.Errorf("error parsing the content: '%w'", err)
	}
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, vars); err != nil {
		return "", fmt.Errorf("error processing the content: %w", err)
	}
	return buf.String(), nil
}
