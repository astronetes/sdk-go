package fsys

import "fmt"

func GetFileContent(dirPath, filename string) ([]byte, error) {
	fs, err := getInstance().FS(dirPath)
	if err != nil {
		return nil, fmt.Errorf("unsupported file system for path '%s': '%w'", dirPath, err)
	}

	file, err := fs.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening the file '%s': '%w'", filename, err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("unexpected error reading the file info for file '%s': '%w'", filename, err)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("the provided path is a directory not a file")
	}

	bytes := make([]byte, info.Size())

	if _, err := file.Read(bytes); err != nil {
		return nil, fmt.Errorf("error reading the file content for '%s': '%w'", filename, err)
	}

	return bytes, nil
}
