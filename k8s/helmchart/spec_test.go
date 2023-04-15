package helmchart

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/astronetes/sdk-go/internal/testfuncs"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
)

func Test_spec_chart(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	chartsPath := filepath.Join(pwd, "testdata", "charts-temp")
	if err := os.MkdirAll(chartsPath, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	if err := testfuncs.PullPackagedChart("bitnami", "https://charts.bitnami.com/bitnami", "nginx", "13.2.33", chartsPath); err != nil {
		t.Fatal(err)
	}

	type fields struct {
		dist               chartDist
		chartPath          string
		valuesTemplatePath string
		vars               map[string]interface{}
		_chart             *chart.Chart
		_values            map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		want    *chart.Chart
		wantErr bool
	}{
		{
			name: "chart path is not provided ",
			fields: fields{
				dist:      packaged,
				chartPath: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "the disttribution of the chart is not supported ",
			fields: fields{
				dist:      packaged + 1,
				chartPath: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "the path to the chart doesn't exist",
			fields: fields{
				dist:      packaged,
				chartPath: "file://unknown_path",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "the file system is not supported",
			fields: fields{
				dist:      packaged,
				chartPath: "unsupported://unknown_path",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "the chart is loaded successfully",
			fields: fields{
				dist:      packaged,
				chartPath: "file://" + chartsPath + "/nginx-13.2.33.tgz",
			},
			want: func() *chart.Chart {
				c, err := loader.Load(chartsPath + "/nginx-13.2.33.tgz")
				if err != nil {
					t.Fatal(err)
				}
				return c
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &spec{
				dist:               tt.fields.dist,
				chartPath:          tt.fields.chartPath,
				valuesTemplatePath: tt.fields.valuesTemplatePath,
				vars:               tt.fields.vars,
				_chart:             tt.fields._chart,
				_values:            tt.fields._values,
			}
			got, err := s.chart()
			if (err != nil) != tt.wantErr {
				t.Errorf("chart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("chart() got = %v, want %v", got, tt.want)
			}
		})
	}
}
