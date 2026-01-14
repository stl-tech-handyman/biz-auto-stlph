package firestore

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

// Client wraps Firestore client
type Client struct {
	client *firestore.Client
}

// NewClient creates a new Firestore client
func NewClient(ctx context.Context, projectID string) (*Client, error) {
	// If project ID is empty, try to get from environment
	if projectID == "" {
		projectID = os.Getenv("GCP_PROJECT_ID")
		if projectID == "" {
			projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
		}
	}

	if projectID == "" {
		return nil, fmt.Errorf("project ID must be provided or set via GCP_PROJECT_ID or GOOGLE_CLOUD_PROJECT environment variable")
	}

	// Create Firestore client
	// It will use Application Default Credentials (ADC) automatically
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create Firestore client: %w", err)
	}

	return &Client{
		client: client,
	}, nil
}

// NewClientWithCredentials creates a Firestore client with explicit credentials
func NewClientWithCredentials(ctx context.Context, projectID string, credentialsJSON string) (*Client, error) {
	if projectID == "" {
		projectID = os.Getenv("GCP_PROJECT_ID")
		if projectID == "" {
			projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
		}
	}

	if projectID == "" {
		return nil, fmt.Errorf("project ID must be provided or set via GCP_PROJECT_ID or GOOGLE_CLOUD_PROJECT environment variable")
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

	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsJSON(credsData))
	if err != nil {
		return nil, fmt.Errorf("failed to create Firestore client: %w", err)
	}

	return &Client{
		client: client,
	}, nil
}

// GetClient returns the underlying Firestore client
func (c *Client) GetClient() *firestore.Client {
	return c.client
}

// Close closes the Firestore client
func (c *Client) Close() error {
	return c.client.Close()
}
