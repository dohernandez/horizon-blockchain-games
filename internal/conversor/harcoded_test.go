package conversor_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dohernandez/horizon-blockchain-games/internal/conversor"
)

func TestHardcoded_ConvertUSD(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	valueDecimal := 1.0
	symbol := "SFL"

	c := conversor.NewHardcoded()

	got, err := c.ConvertUSD(ctx, valueDecimal, symbol)
	require.NoError(t, err)
	require.InEpsilon(t, 0.05649, got, 0)
}
