package fsys

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/astronetes/sdk-go/internal/fsys"
)

func GetFileContent(path string) ([]byte, error) {
	dirPath, filename := filepath.Split(path)
	// TODO: The above line could not work as expected for all the filesystems
	return fsys.GetFileContent(dirPath, filename)
}

func GetAbsoluteFilePath(relativePath string) (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("file://%s/%s", path, relativePath), nil
}
