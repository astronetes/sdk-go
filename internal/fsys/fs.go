package fsys

import "fmt"

func GetFileContent(dirPath, filename string) ([]byte, error) {
	fs, err := getInstance().FS(dirPath)
	if err != nil {
		return nil, fmt.Errorf("unsupported file system for path '%s': '%v'", dirPath, err)
	}
	file, err := fs.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening the file '%s': '%v'", filename, err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("unexpected error reading the file info for file '%s': '%v'", filename, err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("the provided path is a directory not a file")
	}
	var b = make([]byte, info.Size())

	if _, err := file.Read(b); err != nil {
		return nil, fmt.Errorf("error reading the file content for '%s': '%v'", filename, err)
	}
	return b, nil
}
