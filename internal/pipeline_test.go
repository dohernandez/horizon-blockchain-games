package internal_test

import (
	"bytes"
	"context"
	"encoding/csv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/dohernandez/horizon-blockchain-games/internal"
	"github.com/dohernandez/horizon-blockchain-games/internal/entities"
	"github.com/dohernandez/horizon-blockchain-games/internal/mocks"
)

func TestPipeline_Run_all_in_one(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Load sample data with limit 4 to load 3 data lines since offset is -1 which means load the header line.
	dataSample, err := internal.LoadSampleData(4, -1)
	require.NoError(t, err)

	// Mock ExtractProvider.
	provider := mocks.NewExtractProvider(t)

	extBytes := encodeToBytes(t, dataSample, func(t *testing.T, record []string) []string {
		t.Helper()

		return record
	})

	provider.EXPECT().Load(mock.Anything).Return(extBytes, nil)

	// Mock Conversor.
	conversor := mocks.NewConversor(t)

	// Skip the header line.
	for _, record := range dataSample[1:] {
		transaction, err := entities.TransactionNormalize(record)
		require.NoError(t, err)

		conversor.EXPECT().ConvertUSD(mock.Anything, transaction.CurrencyValueDecimal, transaction.CurrencySymbol).Return(1.0, nil)
	}

	// Mock LoadProvider.
	storage := mocks.NewLoadProvider(t)
	storage.EXPECT().Save(mock.Anything, entities.Flatten{
		Date:        "2024-04-15",
		ProjectID:   "4974",
		NumTxs:      3,
		TotalVolume: 3.00,
	}).Return(nil)

	// Mock PipelineBackend.
	b := mocks.NewPipelineBackend(t)
	b.EXPECT().ExtractProvider().Return(provider)
	b.EXPECT().Conversor().Return(conversor)
	b.EXPECT().LoadProvider().Return(storage)

	// Run the pipeline.
	pipeline := internal.NewPipeline(b, internal.PipelineConfig{
		Workers:              1,
		ExtractStepEnabled:   true,
		CalculateStepEnabled: true,
		InsertStepEnabled:    true,
	})

	err = pipeline.Run(ctx)
	require.NoError(t, err)
}

func encodeToBytes(t *testing.T, dataSample [][]string, encoder func(*testing.T, []string) []string) []byte {
	t.Helper()

	// Create a bytes.Buffer to hold the CSV data flatten.
	buf := bytes.Buffer{}
	writer := csv.NewWriter(&buf)

	for _, record := range dataSample {
		err := writer.Write(encoder(t, record))
		require.NoError(t, err)
	}

	// Flush to ensure all data is written to the buffer.
	writer.Flush()
	require.NoError(t, writer.Error())

	return buf.Bytes()
}

func TestPipeline_Run_only_extract_step(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Load sample data with limit 4 to load 3 data lines since offset is -1 which means load the header line.
	dataSample, err := internal.LoadSampleData(4, -1)
	require.NoError(t, err)

	// Mock ExtractProvider.
	provider := mocks.NewExtractProvider(t)

	extBytes := encodeToBytes(t, dataSample, func(t *testing.T, record []string) []string {
		t.Helper()

		return record
	})

	provider.EXPECT().Load(mock.Anything).Return(extBytes, nil)

	// Mock StepProvider.
	stepProvider := mocks.NewStepProvider(t)

	saveBytes := encodeToBytes(t, dataSample[1:], func(t *testing.T, record []string) []string {
		t.Helper()

		tx, err := entities.TransactionNormalize(record)
		require.NoError(t, err)

		return tx.Encode()
	})

	stepProvider.EXPECT().SaveStep(mock.Anything, "extraction", saveBytes).Return(nil)

	// Mock PipelineBackend.
	b := mocks.NewPipelineBackend(t)
	b.EXPECT().ExtractProvider().Return(provider)
	b.EXPECT().StepProvider().Return(stepProvider)

	// Run the pipeline.
	pipeline := internal.NewPipeline(b, internal.PipelineConfig{
		Workers:            1,
		ExtractStepEnabled: true,
	})

	err = pipeline.Run(ctx)
	require.NoError(t, err)
}

func TestPipeline_Run_only_calculation_step(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	dataSample, err := internal.LoadSampleData(3, 0)
	require.NoError(t, err)

	// Mock Conversor.
	conversor := mocks.NewConversor(t)

	for _, record := range dataSample {
		tx, err := entities.TransactionNormalize(record)
		require.NoError(t, err)

		conversor.EXPECT().ConvertUSD(mock.Anything, tx.CurrencyValueDecimal, tx.CurrencySymbol).Return(1.0, nil)
	}

	// Mock StepProvider.
	stepProvider := mocks.NewStepProvider(t)

	loadBytes := encodeToBytes(t, dataSample, func(t *testing.T, record []string) []string {
		t.Helper()

		tx, err := entities.TransactionNormalize(record)
		require.NoError(t, err)

		return tx.Encode()
	})

	stepProvider.EXPECT().LoadStep(mock.Anything, "extraction").Return(loadBytes, nil)

	saveBytes := encodeToBytes(t, [][]string{{"2024-04-15", "4974", "3", "3"}}, func(t *testing.T, record []string) []string {
		t.Helper()

		return record
	})

	stepProvider.EXPECT().SaveStep(mock.Anything, "calculation", saveBytes).Return(nil)

	// Mock PipelineBackend.
	b := mocks.NewPipelineBackend(t)
	b.EXPECT().Conversor().Return(conversor)
	b.EXPECT().StepProvider().Return(stepProvider)

	// Run the pipeline.
	pipeline := internal.NewPipeline(b, internal.PipelineConfig{
		Workers:              1,
		CalculateStepEnabled: true,
	})

	err = pipeline.Run(ctx)
	require.NoError(t, err)
}

func TestPipeline_Run_only_insert_step(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Mock LoadProvider.
	storage := mocks.NewLoadProvider(t)
	storage.EXPECT().Save(mock.Anything, entities.Flatten{
		Date:        "2024-04-15",
		ProjectID:   "4974",
		NumTxs:      3,
		TotalVolume: 3.00,
	}).Return(nil)

	// Mock StepProvider.
	stepProvider := mocks.NewStepProvider(t)

	saveBytes := encodeToBytes(t, [][]string{{"2024-04-15", "4974", "3", "3"}}, func(t *testing.T, record []string) []string {
		t.Helper()

		return record
	})

	stepProvider.EXPECT().LoadStep(mock.Anything, "calculation").Return(saveBytes, nil)

	// Mock PipelineBackend.
	b := mocks.NewPipelineBackend(t)
	b.EXPECT().LoadProvider().Return(storage)
	b.EXPECT().StepProvider().Return(stepProvider)

	// Run the pipeline.
	pipeline := internal.NewPipeline(b, internal.PipelineConfig{
		Workers:           1,
		InsertStepEnabled: true,
	})

	err := pipeline.Run(ctx)
	require.NoError(t, err)
}

func TestPipeline_Run_all_split(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Insert sample data with limit 4 to load 3 data lines since offset is -1 which means load the header line.
	dataSample, err := internal.LoadSampleData(4, -1)
	require.NoError(t, err)

	// Mock ExtractProvider.
	provider := mocks.NewExtractProvider(t)

	extBytes := encodeToBytes(t, dataSample, func(t *testing.T, record []string) []string {
		t.Helper()

		return record
	})

	provider.EXPECT().Load(mock.Anything).Return(extBytes, nil)

	// Mock Conversor.
	conversor := mocks.NewConversor(t)

	// Skip the header line.
	for _, record := range dataSample[1:] {
		transaction, err := entities.TransactionNormalize(record)
		require.NoError(t, err)

		conversor.EXPECT().ConvertUSD(mock.Anything, transaction.CurrencyValueDecimal, transaction.CurrencySymbol).Return(1.0, nil)
	}

	// Mock LoadProvider.
	storage := mocks.NewLoadProvider(t)
	storage.EXPECT().Save(mock.Anything, entities.Flatten{
		Date:        "2024-04-15",
		ProjectID:   "4974",
		NumTxs:      3,
		TotalVolume: 3.00,
	}).Return(nil)

	// Mock StepProvider.
	stepProvider := mocks.NewStepProvider(t)

	// Extract step.
	txBytes := encodeToBytes(t, dataSample[1:], func(t *testing.T, record []string) []string {
		t.Helper()

		tx, err := entities.TransactionNormalize(record)
		require.NoError(t, err)

		return tx.Encode()
	})

	stepProvider.EXPECT().SaveStep(mock.Anything, "extraction", txBytes).Return(nil)
	stepProvider.EXPECT().LoadStep(mock.Anything, "extraction").Return(txBytes, nil)

	// Calculation step.
	conBytes := encodeToBytes(t, [][]string{{"2024-04-15", "4974", "3", "3"}}, func(t *testing.T, record []string) []string {
		t.Helper()

		return record
	})

	stepProvider.EXPECT().SaveStep(mock.Anything, "calculation", conBytes).Return(nil)
	stepProvider.EXPECT().LoadStep(mock.Anything, "calculation").Return(conBytes, nil)

	// Mock PipelineBackend.
	b := mocks.NewPipelineBackend(t)
	b.EXPECT().ExtractProvider().Return(provider)
	b.EXPECT().Conversor().Return(conversor)
	b.EXPECT().LoadProvider().Return(storage)
	b.EXPECT().StepProvider().Return(stepProvider)

	// Run the pipeline.
	wg := &sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer wg.Done()

		pipeline := internal.NewPipeline(b, internal.PipelineConfig{
			InsertStepEnabled: true,
		})

		err := pipeline.Run(ctx)
		assert.NoError(t, err)
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		pipeline := internal.NewPipeline(b, internal.PipelineConfig{
			CalculateStepEnabled: true,
		})

		err := pipeline.Run(ctx)
		assert.NoError(t, err)
	}()

	pipeline := internal.NewPipeline(b, internal.PipelineConfig{
		ExtractStepEnabled: true,
	})

	err = pipeline.Run(ctx)
	require.NoError(t, err)

	wg.Wait()
}
