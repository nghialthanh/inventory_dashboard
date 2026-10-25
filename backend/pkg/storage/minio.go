package storage

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Settings is the private MinIO connection. Endpoint is host:port with no scheme.
type Settings struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// Client is the private MinIO bucket. The browser never receives this client.
type Client struct {
	client *minio.Client
	bucket string
}

// Open builds a MinIO client.
func Open(settings Settings) (*Client, error) {
	client, err := minio.New(settings.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(settings.AccessKey, settings.SecretKey, ""),
		Secure: settings.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio: %w", err)
	}
	return &Client{client: client, bucket: settings.Bucket}, nil
}

// Bucket is the private bucket name.
func (c *Client) Bucket() string {
	return c.bucket
}

// EnsureBucket creates the private bucket when it is missing.
func (c *Client) EnsureBucket(ctx context.Context) error {
	exists, err := c.client.BucketExists(ctx, c.bucket)
	if err != nil {
		return fmt.Errorf("minio: %w", err)
	}
	if exists {
		return nil
	}
	if err := c.client.MakeBucket(ctx, c.bucket, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("minio: %w", err)
	}
	return nil
}
