package internal

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/dohernandez/horizon-blockchain-games/internal/entities"
)

//go:generate mockery --name=ExtractProvider --outpkg=mocks --output=mocks --filename=extract_provider.go --with-expecter

// ExtractProvider is the interface that provides the ability to load the data from the provider.
type ExtractProvider interface {
	// Load loads the data from the provider.
	Load(ctx context.Context) ([]byte, error)
}

// Extract extracts the transactions from the provider and sends them to the output channel.
//
// It receives a provider which loads the data and an output channel to send the normalize transactions.
func Extract(ctx context.Context, provider ExtractProvider, output chan<- entities.Transaction) error {
	data, err := provider.Load(ctx)
	if err != nil {
		return err
	}

	reader := csv.NewReader(bytes.NewReader(data))

	// Read and discard the header line
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("reading header: %w", err)
	}

	if len(header) != entities.InputFieldNum {
		return fmt.Errorf("not enough fields in input: %d", len(header))
	}

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break // End.
		}

		if err != nil {
			return err
		}

		transaction, err := entities.TransactionNormalize(record)
		if err != nil {
			return err
		}

		output <- transaction
	}

	return nil
}
