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

func TestCalculate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	dataSample, err := internal.LoadSampleData(3, 0)
	require.NoError(t, err)

	transactions := make([]entities.Transaction, 0, 20)

	conversor := mocks.NewConversor(t)

	for _, record := range dataSample {
		transaction, err := entities.TransactionNormalize(record)
		require.NoError(t, err)

		conversor.EXPECT().ConvertUSD(mock.Anything, transaction.CurrencyValueDecimal, transaction.CurrencySymbol).Return(1.0, nil)

		transactions = append(transactions, transaction)
	}

	input := make(chan entities.Transaction, 20)

	go func() {
		defer close(input)

		for i, transaction := range transactions {
			if i == 1 {
				transaction.Event = "SELL_ITEMS"
			}

			input <- transaction
		}
	}()

	output := make(chan entities.Flatten, 20)

	err = internal.Calculate(ctx, conversor, input, output)
	require.NoError(t, err)

	for out := range output {
		require.Equal(t, "2024-04-15", out.Date)
		require.Equal(t, "4974", out.ProjectID)

		if out.TotalVolume < 0 {
			require.InEpsilon(t, -1.0, out.TotalVolume, 0)

			return
		}

		require.InEpsilon(t, 1.0, out.TotalVolume, 0)
	}
}
