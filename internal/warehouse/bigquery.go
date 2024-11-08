package warehouse

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/bigquery"

	"github.com/dohernandez/horizon-blockchain-games/internal/entities"
)

// BigQueryType is the type of the BigQuery warehouse.
const BigQueryType = "bigquery"

// BigQueryConfig is the configuration for BigQuery.
type BigQueryConfig struct {
	ProjectID string
	Dataset   string
	Table     string
}

// BigQuery is a target for BigQuery.
type BigQuery struct {
	cfg BigQueryConfig

	client *bigquery.Client

	once sync.Once
}

// NewBigQuery creates a new BigQuery target.
func NewBigQuery(cfg BigQueryConfig) *BigQuery {
	return &BigQuery{
		cfg: cfg,
	}
}

// Save saves the flatten entity into BigQuery.
func (b *BigQuery) Save(ctx context.Context, f entities.Flatten) error {
	var err error

	b.once.Do(func() {
		b.client, err = bigquery.NewClient(ctx, b.cfg.ProjectID)
		if err != nil {
			err = fmt.Errorf("creating BigQuery client: %w", err)

			return
		}
	})

	if err != nil {
		return err
	}

	inserter := b.client.Dataset(b.cfg.Dataset).Table(b.cfg.Table).Inserter()

	if err = inserter.Put(ctx, f); err != nil {
		return fmt.Errorf("inserting data: %w", err)
	}

	return nil
}
