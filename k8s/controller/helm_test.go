package controller

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func Test_values_WithEntry(t *testing.T) {
	type fields struct {
		name string
		vars map[string]interface{}
	}
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Values
	}{
		{
			name: "with empty vars",
			fields: fields{
				name: "values",
				vars: map[string]interface{}{},
			},
			args: args{
				key:   "var1",
				value: 120,
			},
			want: &values{
				name: "values",
				vars: map[string]interface{}{
					"var1": 120,
				},
			},
		},
		{
			name: "override an existing var",
			fields: fields{
				name: "values",
				vars: map[string]interface{}{
					"var1": 120,
				},
			},
			args: args{
				key:   "var1",
				value: 121,
			},
			want: &values{
				name: "values",
				vars: map[string]interface{}{
					"var1": 121,
				},
			},
		},
		{
			name: "override an existing var and change the type of the var",
			fields: fields{
				name: "values",
				vars: map[string]interface{}{
					"var1": 120,
				},
			},
			args: args{
				key:   "var1",
				value: "home",
			},
			want: &values{
				name: "values",
				vars: map[string]interface{}{
					"var1": "home",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &values{
				name: tt.fields.name,
				vars: tt.fields.vars,
			}
			if got := v.WithEntry(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_values_WithEntries(t *testing.T) {
	type fields struct {
		name string
		vars map[string]interface{}
	}
	type args struct {
		entries map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Values
	}{
		{
			name: "vars is empty",
			fields: fields{
				vars: map[string]interface{}{},
			},
			args: args{
				entries: map[string]interface{}{
					"var1": true,
					"var2": 120,
				},
			},
			want: &values{
				vars: map[string]interface{}{
					"var1": true,
					"var2": 120,
				},
			},
		},
		{
			name: "override some of the provided vars",
			fields: fields{
				vars: map[string]interface{}{
					"var1": false,
				},
			},
			args: args{
				entries: map[string]interface{}{
					"var1": true,
					"var2": 120,
				},
			},
			want: &values{
				vars: map[string]interface{}{
					"var1": true,
					"var2": 120,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &values{
				name: tt.fields.name,
				vars: tt.fields.vars,
			}
			if got := v.WithEntries(tt.args.entries); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithEntries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_values_Build(t *testing.T) {
	pwd, _ := os.Getwd()
	type fields struct {
		name string
		vars map[string]interface{}
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "simple template",
			fields: fields{
				vars: map[string]interface{}{
					"var1": 120,
					"var2": "hello",
				},
			},
			args: args{
				path: fmt.Sprintf("file://%s/testdata/sample-values.yml", pwd),
			},
			want: "id: 120\n" +
				"person: hello\n" +
				"option:\n" +
				"  - 120\n" +
				"  - hello\n" +
				"children: 120-hello",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &values{
				name: tt.fields.name,
				vars: tt.fields.vars,
			}
			got, err := v.Build(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Build() got = %v, want %v", got, tt.want)
			}
		})
	}
}
