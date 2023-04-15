package testfuncs

import "os"

const defTmpFilePrefix = "astronetes-sdk-go"

func CreateTemporalFile(content string) (string, error) {
	file, err := os.CreateTemp(os.TempDir(), defTmpFilePrefix)
	if err != nil {
		return "", err
	}
	if _, err := file.WriteString(content); err != nil {
		return "", err
	}
	return file.Name(), nil
}
