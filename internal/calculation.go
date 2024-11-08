package internal

import (
	"context"
	"fmt"

	"github.com/dohernandez/horizon-blockchain-games/internal/entities"
)

//go:generate mockery --name=Conversor --outpkg=mocks --output=mocks --filename=conversor.go --with-expecter

// Conversor is the interface that provides the ability to convert the value in the given currency to USD.
type Conversor interface {
	// ConvertUSD converts the value in the given currency to USD.
	ConvertUSD(ctx context.Context, valueDecimal float64, symbol string) (float64, error)
}

// Calculate calculates the total volume of the transactions in USD.
//
// It receives a channel with the transactions and sends the flatten entities to the output channel.
func Calculate(ctx context.Context, conversor Conversor, input <-chan entities.Transaction, output chan<- entities.Flatten) error {
	for {
		var (
			transaction entities.Transaction
			ok          bool
		)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case transaction, ok = <-input:
			if !ok {
				return nil
			}
		}

		date := transaction.TS.Format("2006-01-02")

		valueUSD, err := conversor.ConvertUSD(ctx, transaction.CurrencyValueDecimal, transaction.CurrencySymbol)
		if err != nil {
			return err
		}

		if transaction.Event != "BUY_ITEMS" && transaction.Event != "SELL_ITEMS" {
			return fmt.Errorf("unknown event: %s", transaction.Event)
		}

		if transaction.Event == "SELL_ITEMS" {
			valueUSD *= -1
		}

		output <- entities.Flatten{
			ProjectID:   transaction.ProjectID,
			Date:        date,
			TotalVolume: valueUSD,
		}
	}
}
