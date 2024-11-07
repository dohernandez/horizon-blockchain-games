package storage_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dohernandez/horizon-blockchain-games/internal/storage"
)

func TestFile_Load(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	data, err := storage.NewFileSystem("../../resources", "sample_data.csv").Load(ctx)
	require.NoError(t, err)

	require.NotNil(t, data)
}

func TestFile_LoadStep(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	data, err := storage.NewFileSystem("testdata", "").LoadStep(ctx, "extraction.csv")
	require.NoError(t, err)

	require.NotNil(t, data)
}

func TestFile_SaveStep(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	data := []byte("sample data")

	f := storage.NewFileSystem("testdata", "")
	err := f.SaveStep(ctx, "calculation.csv", data)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := os.Remove("testdata/calculation.csv")
		require.NoError(t, err)
	})

	require.FileExists(t, "testdata/calculation.csv")
}
