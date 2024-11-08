package conversor_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bool64/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/dohernandez/horizon-blockchain-games/internal/conversor"
)

func TestCoinGecko_ConvertUSD(t *testing.T) {
	ctx := context.Background()

	// Prepare server mock.
	sm, url := httpmock.NewServer()
	defer sm.Close()

	cfg := conversor.CoinGeckoConfig{
		URL:     url,
		KeyType: conversor.DemoKeyType,
		Key:     "CG-UJ2zviozYVh558KpFDL7vR2m",
		TTL:     0,
	}

	// Set successful expectation.
	exp := httpmock.Expectation{
		Method:     http.MethodGet,
		RequestURI: "/simple/price?ids=sunflower-land&vs_currencies=usd",
		RequestHeader: map[string]string{
			"accept":    "application/json",
			cfg.KeyType: cfg.Key,
		},
		Status:       http.StatusOK,
		ResponseBody: []byte(`{"sunflower-land":{"usd":0.059499}}`),
	}
	sm.Expect(exp)

	c := conversor.NewCoinGecko(cfg)

	got, err := c.ConvertUSD(ctx, 1, "sfl")
	require.NoError(t, err)

	require.InEpsilon(t, 0.059499, got, 0)
}
