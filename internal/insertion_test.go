package internal_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/dohernandez/horizon-blockchain-games/internal"
	"github.com/dohernandez/horizon-blockchain-games/internal/entities"
	"github.com/dohernandez/horizon-blockchain-games/internal/mocks"
)

func TestInsert(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	flattens := map[string]entities.Flatten{
		"2024-04-15": {
			Date:        "2024-04-15",
			ProjectID:   "4974",
			NumTxs:      6,
			TotalVolume: 6.00,
		},
		"2024-04-01": {
			Date:        "2024-04-01",
			ProjectID:   "0",
			NumTxs:      10,
			TotalVolume: -10.00,
		},
	}

	storage := mocks.NewWarehouseProvider(t)

	for _, f := range flattens {
		storage.EXPECT().Save(mock.Anything, f).Return(nil)
	}

	for _, f := range flattens {
		err := internal.Insert(ctx, storage, f)
		require.NoError(t, err)
	}
}
