package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/dohernandez/horizon-blockchain-games/internal"
	"github.com/dohernandez/horizon-blockchain-games/internal/conversor"
	"github.com/dohernandez/horizon-blockchain-games/internal/storage"
	"github.com/dohernandez/horizon-blockchain-games/internal/warehouse"
)

var conversorTypes = []string{conversor.GoinGeckoType, conversor.HardcodedType}

// isValidConversor checks if the input is a valid conversor.
func isValidConversor(conversorType string) bool {
	for _, v := range conversorTypes {
		if v == conversorType {
			return true
		}
	}

	return false
}

var conversorAPIKeyTypes = []string{conversor.DemoKeyType, conversor.ProKeyType}

// isValidConversorAPIKeyType checks if the input is a valid conversorAPIKeyTypes.
func isValidConversorAPIKeyType(convAPIKeyType string) bool {
	for _, v := range conversorAPIKeyTypes {
		if v == convAPIKeyType {
			return true
		}
	}

	return false
}

var storageTypes = []string{storage.FileSystemType, storage.BucketType}

// isValidStorage checks if the input is a valid storage.
func isValidStorage(storageType string) bool {
	for _, v := range storageTypes {
		if v == storageType {
			return true
		}
	}

	return false
}

var environments = []string{"dev", "prd"}

// isValidEnvironment checks if the input is a valid environment.
func isValidEnvironment(env string) bool {
	for _, v := range environments {
		if v == env {
			return true
		}
	}

	return false
}

var warehouseType = []string{warehouse.PrintType, warehouse.BigQueryType}

// isValidWarehouse checks if the input is a valid warehouse.
func isValidWarehouse(warehouse string) bool {
	for _, v := range warehouseType {
		if v == warehouse {
			return true
		}
	}

	return false
}

var sequenceFlags = []cli.Flag{
	&cli.StringFlag{
		Name:        "env",
		Required:    false,
		Usage:       "environment",
		DefaultText: "dev",
		Action: func(_ *cli.Context, s string) error {
			if !isValidEnvironment(s) {
				return fmt.Errorf("invalid environment %s", s)
			}

			return nil
		},
		EnvVars: []string{"ENVIRONMENT"},
	},
	&cli.BoolFlag{
		Name:        "extractor",
		Required:    false,
		Usage:       "run only pipeline step extractor",
		DefaultText: "false",
		Aliases:     []string{"e"},
		EnvVars:     []string{"EXTRACTOR_ENABLED"},
	},
	&cli.BoolFlag{
		Name:        "calculator",
		Required:    false,
		Usage:       "run only pipeline step calculator",
		DefaultText: "false",
		Aliases:     []string{"c"},
		EnvVars:     []string{"CALCULATOR_ENABLED"},
	},
	&cli.BoolFlag{
		Name:        "insertion",
		Required:    false,
		Usage:       "run only pipeline step insertion",
		DefaultText: "false",
		Aliases:     []string{"i"},
		EnvVars:     []string{"INSERTION_ENABLED"},
	},
	&cli.BoolFlag{
		Name:        "all",
		Required:    false,
		Usage:       "run all pipeline steps",
		DefaultText: "true",
		Value:       true,
		Aliases:     []string{"a"},
	},
	&cli.UintFlag{
		Name:        "workers",
		Required:    false,
		Usage:       "number of workers to run the pipeline",
		DefaultText: "1",
		Aliases:     []string{"w"},
		EnvVars:     []string{"CALCULATOR_WORKERS"},
	},
	&cli.StringFlag{
		Name:        "dir",
		Required:    false,
		Usage:       "folder or bucket to read/store the intermediate step data when required",
		DefaultText: time.Now().Format("2006-01-02"),
		Value:       time.Now().Format("2006-01-02"),
		EnvVars:     []string{"DIR", "BUCKET"},
	},
	&cli.StringFlag{
		Name:        "file",
		Required:    false,
		Usage:       "file to read the data from",
		DefaultText: "transactions.csv",
		Value:       "transactions.csv",
		EnvVars:     []string{"FILE", "DATA_FILE"},
	},
	&cli.BoolFlag{
		Name:        "test",
		Required:    false,
		Usage:       "run the pipeline in test mode using local file system as providers",
		DefaultText: "false",
	},
	&cli.StringFlag{
		Name:        "conversor",
		Required:    false,
		Usage:       fmt.Sprintf("conversor to use to convert the currency %s", conversorTypes),
		DefaultText: conversor.GoinGeckoType,
		Value:       conversor.GoinGeckoType,
		Action: func(_ *cli.Context, v string) error {
			if !isValidConversor(v) {
				return fmt.Errorf("invalid conversor %s", v)
			}

			return nil
		},
	},
	&cli.StringFlag{
		Name:        "coingecko-api-key-type",
		Required:    false,
		Usage:       "API key type to use with the coingecko conversor",
		DefaultText: conversor.DemoKeyType,
		Value:       conversor.DemoKeyType,
		Action: func(_ *cli.Context, s string) error {
			if !isValidConversorAPIKeyType(s) {
				return fmt.Errorf("invalid conversor API key type %s", s)
			}

			return nil
		},
		EnvVars: []string{"CG_API_KEY_TYPE"},
	},
	&cli.StringFlag{
		Name:     "coingecko-api-key",
		Required: false,
		Usage:    "API key to use with the coingecko conversor",
		EnvVars:  []string{"CG_API_KEY"},
	},
	&cli.StringFlag{
		Name:        "storage-type",
		Required:    false,
		Usage:       fmt.Sprintf("storage type to use to load/store %s", storageTypes),
		DefaultText: storage.BucketType,
		Value:       storage.BucketType,
		Action: func(_ *cli.Context, s string) error {
			if !isValidStorage(s) {
				return fmt.Errorf("invalid storage type %s", s)
			}

			return nil
		},
	},
	&cli.BoolFlag{
		Name:        "verbose",
		Required:    false,
		Usage:       "enable verbose output",
		DefaultText: "false",
		Aliases:     []string{"v"},
		EnvVars:     []string{"VERBOSE"},
	},
	&cli.StringFlag{
		Name:        "warehouse",
		Required:    false,
		Usage:       fmt.Sprintf("target type to use to load/store %s", storageTypes),
		DefaultText: warehouse.BigQueryType,
		Value:       warehouse.BigQueryType,
		Action: func(_ *cli.Context, s string) error {
			if !isValidWarehouse(s) {
				return fmt.Errorf("invalid storage type %s", s)
			}

			return nil
		},
	},
	&cli.StringFlag{
		Name:     "bigquery-dataset",
		Required: false,
		Usage:    "BigQuery dataset in the following format <project_id.dataset.table>",
		Action: func(_ *cli.Context, s string) error {
			parts := strings.Split(s, ".")
			if len(parts) != 3 {
				return fmt.Errorf("invalid BigQuery dataset %s", s)
			}

			return nil
		},
		EnvVars: []string{"BIGQUERY_DATASET"},
	},
}

func main() {
	app := &cli.App{
		Name:  "sequence",
		Usage: "Run a sequence of steps to process data",
		Commands: []*cli.Command{
			{
				Name:        "run",
				Description: "Run pipeline, or a specific step depending on options",
				Flags:       sequenceFlags,
				Action: func(c *cli.Context) error {
					// Backend
					// Configure backend
					cfg, err := loadConfig(c)
					if err != nil {
						return err
					}

					b := internal.NewBackend(cfg)

					// Pipeline
					// Configure pipeline
					cfgPipeline := loadPipelineConfig(c)

					// Run pipeline
					p := internal.NewPipeline(b, cfgPipeline)

					err = p.Run(c.Context)
					if err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func loadConfig(c *cli.Context) (internal.Config, error) {
	cfg := internal.Config{}

	cfg.Environment = c.String("env")
	cfg.Dir = c.String("dir")
	cfg.File = c.String("file")
	cfg.IsTest = c.Bool("test")
	cfg.Logger = c.Bool("verbose")

	if cfg.IsTest {
		return cfg, nil
	}

	cfg.ConversorType = c.String("conversor")

	if cfg.ConversorType == conversor.GoinGeckoType {
		if c.String("coingecko-api-key") == "" {
			return cfg, fmt.Errorf("coingecko api key is required")
		}

		geckoCfg := conversor.CoinGeckoConfig{
			KeyType: c.String("coingecko-api-key-type"),
			Key:     c.String("coingecko-api-key"),
		}

		cfg.CoinGecko = geckoCfg
	}

	cfg.StorageType = c.String("storage-type")

	cfg.WarehouseType = c.String("warehouse")

	if cfg.WarehouseType == warehouse.BigQueryType {
		parts := strings.Split(c.String("bigquery-dataset"), ".")

		cfg.BigQuery = warehouse.BigQueryConfig{
			ProjectID: parts[0],
			Dataset:   parts[1],
			Table:     parts[2],
		}
	}

	return cfg, nil
}

func loadPipelineConfig(c *cli.Context) internal.PipelineConfig {
	cfgPipeline := internal.PipelineConfig{}

	all := c.Bool("all")

	if c.Bool("extractor") {
		cfgPipeline.ExtractStepEnabled = true

		all = false
	}

	if c.Bool("calculator") {
		cfgPipeline.CalculateStepEnabled = true

		all = false
	}

	if c.Bool("insertion") {
		cfgPipeline.InsertStepEnabled = true

		all = false
	}

	if all {
		cfgPipeline.ExtractStepEnabled = true
		cfgPipeline.CalculateStepEnabled = true
		cfgPipeline.InsertStepEnabled = true
	}

	cfgPipeline.Workers = c.Int("workers")

	return cfgPipeline
}
