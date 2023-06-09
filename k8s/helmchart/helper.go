package helmchart

import (
	"path/filepath"

	"github.com/astronetes/sdk-go/internal/fsys"
)

func readFile(path string) ([]byte, error) {
	dirPath, filename := filepath.Split(path)
	// TODO: The above line could not work as expected for all the filesystems
	return fsys.GetFileContent(dirPath, filename)
}
