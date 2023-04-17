package helmchart

import (
	"bytes"
	"fmt"

	"github.com/astronetes/sdk-go/internal/gotemplate"
	"github.com/ghodss/yaml"
	"golang.org/x/exp/maps"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
)

type chartDist int

const (
	defaultValuesTemplateName           = "values"
	packaged                  chartDist = iota
)

type Spec interface {
	WithValuesFileTemplate(valuesTemplatePath string) Spec
	With(key string, value interface{}) Spec
	WithValues(entries map[string]interface{}) Spec
	values() (map[string]interface{}, error)
	chart() (*chart.Chart, error)
	chartAndValues() (*chart.Chart, map[string]interface{}, error)
}

type spec struct {
	dist                  chartDist
	chartPath             string
	valuesTemplatePath    string
	valuesTemplateContent string
	vars                  map[string]interface{}
	_chart                *chart.Chart
	_values               map[string]interface{}
}

func (s *spec) hasValuesTemplate() bool {
	return s.valuesTemplatePath != "" || s.valuesTemplateContent != ""
}

func (s *spec) chart() (*chart.Chart, error) {
	if s._chart != nil {
		return s._chart, nil
	}

	content, err := readFile(s.chartPath)
	if err != nil {
		return nil, fmt.Errorf("error reading content of packaged chart: ''%w", err)
	}

	if s.dist == packaged {
		s._chart, err = loader.LoadArchive(bytes.NewReader(content))
		if err != nil {
			return nil, fmt.Errorf("error loading packaged chart: ''%w", err)
		}
		return s._chart, nil
	}
	return nil, fmt.Errorf("unsupported distribution type of chart")
}

func (s *spec) values() (map[string]interface{}, error) {
	if s._values != nil {
		return s._values, nil
	}

	if !s.hasValuesTemplate() {
		return s.vars, nil
	}
	var content = s.valuesTemplateContent
	if s.valuesTemplatePath != "" {
		contentBytes, err := readFile(s.valuesTemplatePath)
		if err != nil {
			return nil, fmt.Errorf("error reading content of values template: ''%w", err)
		}
		content = string(contentBytes)
	}

	valuesContent, err := gotemplate.ApplyTemplateWithVariables(defaultValuesTemplateName, content, s.vars)
	if err != nil {
		return nil, err
	}
	s._values = map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(valuesContent), &s._values); err != nil {
		return nil, fmt.Errorf("error unmarshaling the values map: '%w'", err)
	}
	return s._values, nil
}

func (s *spec) chartAndValues() (*chart.Chart, map[string]interface{}, error) {
	chart, err := s.chart()
	if err != nil {
		return nil, nil, err
	}
	values, err := s.values()
	if err != nil {
		return nil, nil, err
	}
	return chart, values, nil
}

// LoadPackagedChart initializes a Spec struct from a packaged chart spec.
// The path to the chart can reference to any of the supported filesystem.
// It returns a pointer to a Spec.
func LoadPackagedChart(chartPath string) Spec {
	return &spec{
		dist:                  packaged,
		chartPath:             chartPath,
		valuesTemplatePath:    "",
		valuesTemplateContent: "",
		vars:                  make(map[string]interface{}, 0),
		_chart:                nil,
		_values:               nil,
	}
}
func (s *spec) WithValuesTextTemplate(text string) Spec {
	s.valuesTemplateContent = text
	s.valuesTemplatePath = ""

	return s
}

func (s *spec) WithValuesFileTemplate(valuesTemplatePath string) Spec {
	s.valuesTemplatePath = valuesTemplatePath
	s.valuesTemplateContent = ""

	return s
}

func (s *spec) With(key string, value interface{}) Spec {
	s.vars[key] = value

	return s
}

func (s *spec) WithValues(entries map[string]interface{}) Spec {
	maps.Copy(s.vars, entries)

	return s
}
