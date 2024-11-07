package conversor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bool64/ctxd"
)

// GoinGeckoType is the type of the CoinGecko conversor.
const GoinGeckoType = "coingecko"

const (
	// DemoKeyType is the key type for the CoinGecko demo API key.
	DemoKeyType = "x_cg_demo_api_key"
	// DemoBaseURL is the base URL for the CoinGecko demo API.
	DemoBaseURL = "https://api.coingecko.com/api/v3/"

	// ProKeyType is the key type for the CoinGecko pro API key.
	ProKeyType = "x-cg-pro-api-key"
	// ProBaseURL is the base URL for the CoinGecko pro API.
	ProBaseURL = "https://pro-api.coingecko.com/api/v3/"
)

var symbolToID = map[string]string{
	"sfl":    "sunflower-land",
	"matic":  "matic-network",
	"usdc":   "usd-coin",
	"usdc.e": "bridged-usdc-polygon-pos-bridge",
}

// CoinGeckoConfig holds the configuration for the CoinGecko client.
type CoinGeckoConfig struct {
	URL     string
	KeyType string
	Key     string
	TTL     time.Duration
}

// Option is a convenience type which will be used to modify Client private fields.
type Option func(client *CoinGecko)

// WithTransport configures the transport of a Client.
func WithTransport(transport http.RoundTripper) Option {
	return func(c *CoinGecko) {
		if transport != nil {
			return
		}

		c.transport = transport
	}
}

// WithLogger configures the logger of a Client.
func WithLogger(logger ctxd.Logger) Option {
	return func(c *CoinGecko) {
		if logger == nil {
			return
		}

		c.logger = logger
	}
}

// CoinGecko is a client for the CoinGecko API.
type CoinGecko struct {
	cfg CoinGeckoConfig

	transport http.RoundTripper

	mapRates map[string]float64
	sm       sync.RWMutex

	logger ctxd.Logger
}

// NewCoinGecko creates a new CoinGecko client with the given configuration.
func NewCoinGecko(cfg CoinGeckoConfig, opts ...Option) *CoinGecko {
	c := &CoinGecko{
		cfg:       cfg,
		transport: http.DefaultTransport,
		mapRates:  make(map[string]float64),
		logger:    ctxd.NoOpLogger{},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// ConvertUSD converts the given value in USD to the given currency.
func (c *CoinGecko) ConvertUSD(ctx context.Context, valueDecimal float64, symbol string) (float64, error) {
	symbol = strings.ToLower(symbol)

	c.sm.RLock()
	if price, ok := c.mapRates[symbol]; ok {
		c.sm.RUnlock()

		return valueDecimal * price, nil
	}
	c.sm.RUnlock()

	ctx = ctxd.AddFields(ctx, "symbol", symbol)

	id, ok := symbolToID[symbol]
	if !ok {
		return 0, fmt.Errorf("unknown currency: %s", symbol)
	}

	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd", c.cfg.URL, id)

	var (
		ctxc   = ctx
		cancel = func() {}
	)

	if c.cfg.TTL > 0 {
		c.logger.Debug(ctx, "setting timeout to request", "timeout", c.cfg.TTL)

		ctxc, cancel = context.WithTimeout(ctx, c.cfg.TTL)
	}

	defer cancel()

	req, err := http.NewRequestWithContext(ctxc, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Accept", "application/json")

	c.logger.Debug(ctx, "setting api key header", "api_key_type", c.cfg.KeyType)

	req.Header.Add(c.cfg.KeyType, c.cfg.Key)

	c.logger.Debug(ctx, "requesting price", "url", url)

	res, err := c.transport.RoundTrip(req)
	if err != nil {
		return 0, fmt.Errorf("doing request: %w", err)
	}

	defer res.Body.Close() //nolint:errcheck

	c.logger.Debug(ctx, "reading response body")

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("reading body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d, body: %s", res.StatusCode, string(body))
	}

	c.logger.Debug(ctx, "unmarshaling response body")

	var priceJSON map[string]map[string]float64

	err = json.Unmarshal(body, &priceJSON)
	if err != nil {
		return 0, fmt.Errorf("unmarshaling body: %w", err)
	}

	priceObj, ok := priceJSON[id]
	if !ok {
		return 0, fmt.Errorf("currency not found: %s", symbol)
	}

	price := priceObj["usd"]

	c.sm.Lock()
	c.mapRates[symbol] = price
	c.sm.Unlock()

	c.logger.Debug(ctx, "got price", "price", price)

	return valueDecimal * price, nil
}
