package source_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dohernandez/horizon-blockchain-games/internal/source"
)

func TestFile_Load(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	data, err := source.NewFile("../../resources", "sample_data.csv").Load(ctx)
	require.NoError(t, err)

	require.NotNil(t, data)
}
