package source

import (
	"context"
	"fmt"
	"os"
	"path"
)

// File is a source that loads data from a file.
type File struct {
	dir  string
	file string
}

// NewFile creates a new file source.
func NewFile(dir, file string) *File {
	return &File{dir: dir, file: file}
}

// Load loads the data from the file.
func (f *File) Load(_ context.Context) ([]byte, error) {
	data, err := os.ReadFile(path.Join(f.dir, f.file))
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	return data, nil
}
