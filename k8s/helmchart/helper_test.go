package helmchart

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/astronetes/sdk-go/internal/testfuncs"
)

func Test_readFile(t *testing.T) {
	fileContent := `
	Hello my friend,
		How're you?
	Regards'
	`
	filePath, err := testfuncs.CreateTemporalFile(fileContent)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filePath)

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "read a local file",
			args: args{
				path: fmt.Sprintf("file://%s", filePath),
			},
			want:    []byte(fileContent),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("readFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
