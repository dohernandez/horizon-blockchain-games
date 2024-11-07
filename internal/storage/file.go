package storage

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// File is a storage that save/loads data to/from a file.
type File struct {
	dir string
}

// NewFile creates a new file storage.
func NewFile(dir string) *File {
	return &File{dir: dir}
}

// Load loads the data from the file.
func (f *File) Load(_ context.Context, file string) ([]byte, error) {
	filePath := filepath.Clean(path.Join(f.dir, file))

	if !strings.HasPrefix(filePath, f.dir) {
		return nil, fmt.Errorf("invalid path")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	return data, nil
}

// Save saves the data to the file.
func (f *File) Save(_ context.Context, file string, data []byte) error {
	filePath := filepath.Clean(path.Join(f.dir, file))

	if !strings.HasPrefix(filePath, f.dir) {
		return fmt.Errorf("invalid path")
	}

	st, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer st.Close() //nolint:errcheck

	_, err = st.Write(data)
	if err != nil {
		return err
	}

	// Flush and ensure data is written to disk.
	return st.Sync()
}
