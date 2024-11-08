package storage

import (
	"context"
	"fmt"
	"io"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/bool64/ctxd"
	"google.golang.org/api/option"
)

// BucketType is the type of the GoogleBucket storage.
const BucketType = "bucket"

// GoogleBucketConfig holds the configuration for the Google Cloud Storage bucket.
type GoogleBucketConfig struct {
	Bucket string
	File   string
}

// GoogleBucket is a storage that save/loads data to/from a Google Cloud Storage bucket.
type GoogleBucket struct {
	cfg GoogleBucketConfig

	// client is the Google Cloud Storage client.
	client *storage.Client

	once sync.Once

	logger ctxd.Logger

	endpoint string
}

// Option is a convenience type which will be used to modify GoogleBucket private fields.
type Option func(b *GoogleBucket)

// WithLogger configures the logger of a GoogleBucket.
func WithLogger(logger ctxd.Logger) Option {
	return func(b *GoogleBucket) {
		if logger == nil {
			return
		}

		b.logger = logger
	}
}

// WithEndpoint configures the endpoint of a GoogleBucket.
func WithEndpoint(endpoint string) Option {
	return func(b *GoogleBucket) {
		if endpoint == "" {
			return
		}

		b.endpoint = endpoint
	}
}

// NewGoogleBucket creates a new GoogleBucket source.
func NewGoogleBucket(cfg GoogleBucketConfig, opts ...Option) *GoogleBucket {
	g := &GoogleBucket{
		cfg:    cfg,
		logger: ctxd.NoOpLogger{},
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

// Load loads the data from the Google bucket.
func (g *GoogleBucket) Load(ctx context.Context) ([]byte, error) {
	err := g.loadClient(ctx)
	if err != nil {
		return nil, err
	}

	g.logger.Debug(ctx, "creating reader", "bucket", g.cfg.Bucket, "file", g.cfg.File)

	reader, err := g.client.Bucket(g.cfg.Bucket).Object(g.cfg.File).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating reader: %w", err)
	}

	defer reader.Close() //nolint: errcheck

	return io.ReadAll(reader)
}

func (g *GoogleBucket) loadClient(ctx context.Context) error {
	var err error

	g.once.Do(func() {
		var opts []option.ClientOption

		g.logger.Debug(ctx, "creating Google Cloud Storage client")

		if g.endpoint != "" {
			g.logger.Debug(ctx, "using custom endpoint", "endpoint", g.endpoint)

			opts = append(opts, option.WithEndpoint(g.endpoint))
		}

		g.client, err = storage.NewClient(ctx, opts...)
		if err != nil {
			err = fmt.Errorf("creating Google Cloud Storage client: %w", err)

			return
		}

		g.logger.Debug(ctx, "Google Cloud Storage client created")
	})

	return err
}

// LoadStep loads the data from the Google bucket.
func (g *GoogleBucket) LoadStep(ctx context.Context, file string) ([]byte, error) {
	err := g.loadClient(ctx)
	if err != nil {
		return nil, err
	}

	g.logger.Debug(ctx, "creating reader", "bucket", g.cfg.Bucket, "file", file)

	reader, err := g.client.Bucket(g.cfg.Bucket).Object(file).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating reader: %w", err)
	}

	defer reader.Close() //nolint: errcheck

	g.logger.Debug(ctx, "reading file", "bucket", g.cfg.Bucket, "file", file)

	return io.ReadAll(reader)
}

// SaveStep saves the data to the Google bucket.
func (g *GoogleBucket) SaveStep(ctx context.Context, file string, data []byte) error {
	err := g.loadClient(ctx)
	if err != nil {
		return err
	}

	g.logger.Debug(ctx, "creating writer", "bucket", g.cfg.Bucket, "file", file)

	writer := g.client.Bucket(g.cfg.Bucket).Object(file).NewWriter(ctx)

	g.logger.Debug(ctx, "writing data", "bucket", g.cfg.Bucket, "file", file)

	if _, err := writer.Write(data); err != nil {
		return fmt.Errorf("writing data: %w", err)
	}

	return writer.Close()
}

// Close closes the Google Cloud Storage client.
func (g *GoogleBucket) Close() error {
	if g.client == nil {
		return nil
	}

	return g.client.Close()
}
