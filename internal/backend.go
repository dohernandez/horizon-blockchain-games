package internal

import (
	"context"

	"github.com/bool64/ctxd"
	"github.com/bool64/zapctxd"
	"go.uber.org/zap/zapcore"

	"github.com/dohernandez/horizon-blockchain-games/internal/conversor"
	"github.com/dohernandez/horizon-blockchain-games/internal/storage"
	"github.com/dohernandez/horizon-blockchain-games/internal/warehouse"
)

// Config holds the configuration for the backend to create the providers dependencies.
type Config struct {
	// Environment is the environment to run the application.
	Environment string

	// Dir is the directory (filesystem) or bucket to store the data.
	Dir string
	// File is the file to load the data.
	File string
	// IsTest is to run the application in test mode.
	// When the application is in test mode, all dependencies are filesystem.
	IsTest bool

	// ConversorType is the type of conversor to use.
	ConversorType string

	// CoinGecko holds the configuration for the CoinGecko conversor.
	CoinGecko conversor.CoinGeckoConfig

	// StorageType is the type of storage to use.
	StorageType string

	// GCPBucketEndpoint is the endpoint of the Google Cloud Storage bucket.
	// It is mainly used for testing purposes.
	GCPBucketEndpoint string

	// WarehouseType is the type of warehouse to use.
	WarehouseType string

	BigQuery warehouse.BigQueryConfig

	// Logger is to enable logger.
	Logger bool
}

// Backend is the main struct that holds the providers dependencies.
type Backend struct {
	cfg Config

	extractProvider ExtractProvider
	conversor       Conversor
	loadProvider    WarehouseProvider
	stepProvider    StepProvider
}

// NewBackend creates a new backend with the given configuration.
func NewBackend(cfg Config) *Backend {
	ctx := context.Background()

	b := Backend{
		cfg: cfg,
	}

	var logger ctxd.Logger

	logger = ctxd.NoOpLogger{}

	if cfg.Logger {
		logger = zapctxd.New(zapctxd.Config{
			Level:   zapcore.DebugLevel,
			DevMode: false,
		})
	}

	st := storage.NewFileSystem(b.cfg.Dir, b.cfg.File)

	logger.Debug(ctx, "initializing extractProvider with filesystem storage")

	b.extractProvider = st

	logger.Debug(ctx, "initializing stepProvider with filesystem storage")

	b.stepProvider = st

	logger.Debug(ctx, "initializing conversor with hardcoded values")

	b.conversor = conversor.NewHardcoded()

	logger.Debug(ctx, "initializing loadProvider with print target")

	b.loadProvider = &warehouse.Print{}

	if cfg.IsTest {
		return &b
	}

	// Storage
	if cfg.StorageType == "bucket" {
		logger.Debug(ctx, "replacing extractProvider with google bucket storage")

		gbCfg := storage.GoogleBucketConfig{
			Bucket: cfg.Dir,
			File:   cfg.File,
		}

		opts := []storage.Option{storage.WithLogger(logger)}

		if cfg.GCPBucketEndpoint != "" {
			opts = append(opts, storage.WithEndpoint(cfg.GCPBucketEndpoint))
		}

		stgcp := storage.NewGoogleBucket(gbCfg, opts...)

		b.extractProvider = stgcp
		b.stepProvider = stgcp
	}

	// Conversor.
	if cfg.ConversorType == conversor.GoinGeckoType {
		logger.Debug(ctx, "replacing conversor with CoinGecko")

		cfg.CoinGecko.URL = conversor.DemoBaseURL

		if cfg.CoinGecko.KeyType == conversor.ProKeyType {
			cfg.CoinGecko.URL = conversor.ProBaseURL
		}

		logger.Debug(ctx, "CoinGecko conversor configuration", "config", cfg.CoinGecko)

		b.conversor = conversor.NewCoinGecko(cfg.CoinGecko, conversor.WithLogger(logger))
	}

	// Warehouse.
	if cfg.WarehouseType == warehouse.BigQueryType {
		logger.Debug(ctx, "replacing warehouse with BigQuery")

		b.loadProvider = warehouse.NewBigQuery(cfg.BigQuery)
	}

	return &b
}

// ExtractProvider returns the extract provider.
func (b *Backend) ExtractProvider() ExtractProvider {
	return b.extractProvider
}

// Conversor returns the conversor.
func (b *Backend) Conversor() Conversor {
	return b.conversor
}

// WarehouseProvider returns the load provider.
func (b *Backend) WarehouseProvider() WarehouseProvider {
	return b.loadProvider
}

// StepProvider returns the step provider.
func (b *Backend) StepProvider() StepProvider {
	return b.stepProvider
}
