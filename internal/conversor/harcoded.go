package conversor

import (
	"context"
	"fmt"
	"strings"
)

// Hardcoded is a hardcoded implementation of the conversor.
//
// It contains a map of exchange rates for some currencies.
type Hardcoded struct {
	exchangeRate map[string]float64
}

// NewHardcoded creates a new hardcoded conversor.
func NewHardcoded() *Hardcoded {
	return &Hardcoded{
		exchangeRate: map[string]float64{
			"SFL":    0.05649,
			"MATIC":  0.3264,
			"USDC":   1,
			"USDC.E": 1,
		},
	}
}

// ConvertUSD converts a value from a currency to USD based on the hardcoded exchange rates.
func (c *Hardcoded) ConvertUSD(_ context.Context, valueDecimal float64, symbol string) (float64, error) {
	upper := strings.ToUpper(symbol)

	rate, ok := c.exchangeRate[upper]
	if !ok {
		return 0, fmt.Errorf("unknown currency: %s", symbol)
	}

	return valueDecimal * rate, nil
}
