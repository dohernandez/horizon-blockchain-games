package storage

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// FileSystem is a storage that save/loads data to/from a file.
type FileSystem struct {
	dir string

	file string
}

// NewFileSystem creates a new file storage.
func NewFileSystem(dir, file string) *FileSystem {
	return &FileSystem{
		dir:  dir,
		file: file,
	}
}

// Load loads the data from the file.
func (f *FileSystem) Load(_ context.Context) ([]byte, error) {
	data, err := os.ReadFile(path.Join(f.dir, f.file))
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	return data, nil
}

// LoadStep loads the data from the file.
func (f *FileSystem) LoadStep(_ context.Context, file string) ([]byte, error) {
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

// SaveStep saves the data to the file.
func (f *FileSystem) SaveStep(_ context.Context, file string, data []byte) error {
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
