package storage

import (
	"context"
	"fmt"
	"io"

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

// Put stores or replaces an object in the private bucket.
func (c *Client) Put(ctx context.Context, key, contentType string, body io.Reader, size int64) error {
	_, err := c.client.PutObject(ctx, c.bucket, key, body, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("minio: %w", err)
	}
	return nil
}

// Get opens an object from the private bucket.
func (c *Client) Get(ctx context.Context, key string) (io.ReadCloser, string, error) {
	obj, err := c.client.GetObject(ctx, c.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", fmt.Errorf("minio: %w", err)
	}
	info, err := obj.Stat()
	if err != nil {
		_ = obj.Close()
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return nil, "", ErrNotFound
		}
		return nil, "", fmt.Errorf("minio: %w", err)
	}
	return obj, info.ContentType, nil
}

// Delete removes an object. Missing keys succeed.
func (c *Client) Delete(ctx context.Context, key string) error {
	err := c.client.RemoveObject(ctx, c.bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("minio: %w", err)
	}
	return nil
}
