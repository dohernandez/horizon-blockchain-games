package storage

import (
	"context"
	"fmt"
	"io"
	"sync"

	"cloud.google.com/go/storage"
)

// GoogleBucket is a storage that save/loads data to/from a Google Cloud Storage bucket.
type GoogleBucket struct {
	// bucket is the name of the bucket.
	bucket string

	file string

	// client is the Google Cloud Storage client.
	client *storage.Client

	once sync.Once
}

// NewGoogleBucket creates a new GoogleBucket source.
func NewGoogleBucket(bucket, file string) *GoogleBucket {
	return &GoogleBucket{
		bucket: bucket,
		file:   file,
	}
}

// Load loads the data from the Google bucket.
func (g *GoogleBucket) Load(ctx context.Context) ([]byte, error) {
	var err error

	g.once.Do(func() {
		var storageClient *storage.Client

		storageClient, err = storage.NewClient(ctx)
		if err != nil {
			err = fmt.Errorf("creating Google Cloud Storage client: %w", err)

			return
		}

		g.client = storageClient
	})

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating Google Cloud Storage client: %w", err)
	}

	reader, err := client.Bucket(g.bucket).Object(g.file).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating reader: %w", err)
	}

	defer reader.Close() //nolint: errcheck

	return io.ReadAll(reader)
}

// LoadStep loads the data from the Google bucket.
func (g *GoogleBucket) LoadStep(ctx context.Context, file string) ([]byte, error) {
	reader, err := g.client.Bucket(g.bucket).Object(file).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating reader: %w", err)
	}

	defer reader.Close() //nolint: errcheck

	return io.ReadAll(reader)
}

// SaveStep saves the data to the Google bucket.
func (g *GoogleBucket) SaveStep(ctx context.Context, file string, data []byte) error {
	writer := g.client.Bucket(g.bucket).Object(file).NewWriter(ctx)

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
