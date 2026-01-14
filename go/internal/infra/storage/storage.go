package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// Client wraps Cloud Storage client
type Client struct {
	client     *storage.Client
	bucketName string
}

// NewClient creates a new Cloud Storage client
func NewClient(ctx context.Context, bucketName string) (*Client, error) {
	// If bucket name is empty, try to get from environment
	if bucketName == "" {
		bucketName = os.Getenv("GCS_BUCKET_NAME")
	}

	if bucketName == "" {
		return nil, fmt.Errorf("bucket name must be provided or set via GCS_BUCKET_NAME environment variable")
	}

	// Create Storage client
	// It will use Application Default Credentials (ADC) automatically
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Storage client: %w", err)
	}

	return &Client{
		client:     client,
		bucketName: bucketName,
	}, nil
}

// NewClientWithCredentials creates a Storage client with explicit credentials
func NewClientWithCredentials(ctx context.Context, bucketName string, credentialsJSON string) (*Client, error) {
	if bucketName == "" {
		bucketName = os.Getenv("GCS_BUCKET_NAME")
	}

	if bucketName == "" {
		return nil, fmt.Errorf("bucket name must be provided or set via GCS_BUCKET_NAME environment variable")
	}

	// Try to read from file if it's a path, otherwise use as JSON string
	var credsData []byte
	if _, err := os.Stat(credentialsJSON); err == nil {
		var readErr error
		credsData, readErr = os.ReadFile(credentialsJSON)
		if readErr != nil {
			return nil, fmt.Errorf("failed to read credentials file: %w", readErr)
		}
	} else {
		credsData = []byte(credentialsJSON)
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(credsData))
	if err != nil {
		return nil, fmt.Errorf("failed to create Storage client: %w", err)
	}

	return &Client{
		client:     client,
		bucketName: bucketName,
	}, nil
}

// UploadFile uploads a file to Cloud Storage
func (c *Client) UploadFile(ctx context.Context, path string, data []byte, contentType string) (string, error) {
	bucket := c.client.Bucket(c.bucketName)
	obj := bucket.Object(path)

	writer := obj.NewWriter(ctx)
	writer.ContentType = contentType
	writer.CacheControl = "public, max-age=3600" // Cache for 1 hour

	if _, err := writer.Write(data); err != nil {
		writer.Close()
		return "", fmt.Errorf("failed to write to storage: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close storage writer: %w", err)
	}

	// Return public URL
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucketName, path), nil
}

// GetSignedURL generates a signed URL for temporary access to a file
func (c *Client) GetSignedURL(ctx context.Context, path string, duration time.Duration) (string, error) {
	bucket := c.client.Bucket(c.bucketName)
	obj := bucket.Object(path)

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(duration),
	}

	url, err := bucket.SignedURL(obj.ObjectName(), opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return url, nil
}

// DeleteFile deletes a file from Cloud Storage
func (c *Client) DeleteFile(ctx context.Context, path string) error {
	bucket := c.client.Bucket(c.bucketName)
	obj := bucket.Object(path)

	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// ReadFile reads a file from Cloud Storage
func (c *Client) ReadFile(ctx context.Context, path string) ([]byte, error) {
	bucket := c.client.Bucket(c.bucketName)
	obj := bucket.Object(path)

	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// Close closes the Storage client
func (c *Client) Close() error {
	return c.client.Close()
}
