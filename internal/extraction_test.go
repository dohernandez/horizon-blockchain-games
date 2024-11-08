package internal_test

import (
	"bytes"
	"context"
	"encoding/csv"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/dohernandez/horizon-blockchain-games/internal"
	"github.com/dohernandez/horizon-blockchain-games/internal/entities"
	"github.com/dohernandez/horizon-blockchain-games/internal/mocks"
)

func TestExtract(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Load sample data with limit 4 to load 3 data lines since offset is -1 which means load the header line.
	dataSample, err := internal.LoadSampleData(4, -1)
	require.NoError(t, err)

	// Create a bytes.Buffer to hold the CSV data from the sample data.
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	for _, record := range dataSample {
		err = writer.Write(record)
		require.NoError(t, err)
	}

	// Flush to ensure all data is written to the buffer.
	writer.Flush()
	require.NoError(t, writer.Error())

	provider := mocks.NewExtractProvider(t)
	provider.EXPECT().Load(mock.Anything).Return(buf.Bytes(), nil)

	output := make(chan entities.Transaction, 20)

	err = internal.Extract(ctx, provider, output)
	require.NoError(t, err)

	close(output)

	require.Len(t, output, 3)

	// Check if the records are the same as the sample data.
	// Skip the header line.
	i := 1

	for out := range output {
		require.Equal(t, dataSample[i][1], out.TS.Format("2006-01-02 15:04:05.000"))
		require.Equal(t, dataSample[i][2], out.Event)
		require.Equal(t, dataSample[i][3], out.ProjectID)

		i++
	}
}
