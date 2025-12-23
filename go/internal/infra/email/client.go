package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bizops360/go-api/internal/ports"
)

// EmailServiceClient implements Mailer interface via HTTP calls to email service
type EmailServiceClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewEmailServiceClient creates a new email service client
// If EMAIL_SERVICE_URL is not set, returns nil to use direct Gmail implementation
func NewEmailServiceClient() *EmailServiceClient {
	baseURL := os.Getenv("EMAIL_SERVICE_URL")
	
	// If EMAIL_SERVICE_URL is not set, use direct Gmail implementation (return nil)
	// The handler will use GmailSender instead
	if baseURL == "" {
		return nil
	}
	
	// If EMAIL_SERVICE_URL points to self (same service), return nil to use direct implementation
	// This allows the same API to handle email requests directly
	if strings.Contains(baseURL, "localhost") || strings.Contains(baseURL, "127.0.0.1") {
		return nil
	}

	apiKey := os.Getenv("EMAIL_SERVICE_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("SERVICE_API_KEY")
	}

	return &EmailServiceClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendEmail sends an email via the email service
func (c *EmailServiceClient) SendEmail(ctx context.Context, req *ports.SendEmailRequest) (*ports.SendEmailResult, error) {
	return c.callEmailService(ctx, "/api/email/send", req)
}

// SendEmailDraft creates an email draft via the email service
func (c *EmailServiceClient) SendEmailDraft(ctx context.Context, req *ports.SendEmailRequest) (*ports.SendEmailResult, error) {
	return c.callEmailService(ctx, "/api/email/draft", req)
}

func (c *EmailServiceClient) callEmailService(ctx context.Context, path string, req *ports.SendEmailRequest) (*ports.SendEmailResult, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("EMAIL_SERVICE_API_KEY is not configured")
	}

	payload := map[string]interface{}{
		"to":       req.To,
		"subject":  req.Subject,
		"html":     req.HTMLBody,
		"text":     req.TextBody,
		"from":     req.From,
		"fromName": req.FromName,
	}

	if len(req.Attachments) > 0 {
		attachments := make([]map[string]interface{}, len(req.Attachments))
		for i, att := range req.Attachments {
			attachments[i] = map[string]interface{}{
				"filename": att.Filename,
				"content":  att.Content,
				"mimeType": att.MimeType,
			}
		}
		payload["attachments"] = attachments
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Api-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call email service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		_ = json.Unmarshal(body, &errorResp)
		errorMsg := "unknown error"
		if errMsg, ok := errorResp["error"].(string); ok {
			errorMsg = errMsg
		}
		return &ports.SendEmailResult{
			Success: false,
			Error:   &errorMsg,
		}, fmt.Errorf("email service error (%d): %s", resp.StatusCode, errorMsg)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	messageID := ""
	if id, ok := result["messageId"].(string); ok {
		messageID = id
	}

	return &ports.SendEmailResult{
		MessageID: messageID,
		Success:   true,
	}, nil
}

// SendBookingDepositEmail sends a booking deposit email via the email service
func (c *EmailServiceClient) SendBookingDepositEmail(ctx context.Context, payload map[string]interface{}) (map[string]interface{}, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("EMAIL_SERVICE_API_KEY is not configured")
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/email/booking-deposit", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Api-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call email service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		_ = json.Unmarshal(body, &errorResp)
		return nil, fmt.Errorf("email service error (%d): %v", resp.StatusCode, errorResp)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}


