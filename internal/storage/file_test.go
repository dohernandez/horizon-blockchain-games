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

	data, err := storage.NewFile("testdata").Load(ctx, "extraction.csv")
	require.NoError(t, err)

	require.NotNil(t, data)
}

func TestFile_Save(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	data := []byte("sample data")

	f := storage.NewFile("testdata")
	err := f.Save(ctx, "calculation.csv", data)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := os.Remove("testdata/calculation.csv")
		require.NoError(t, err)
	})

	require.FileExists(t, "testdata/calculation.csv")
}
