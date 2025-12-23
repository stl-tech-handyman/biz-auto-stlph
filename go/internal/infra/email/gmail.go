package email

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/bizops360/go-api/internal/ports"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// GmailSender implements Mailer interface using Gmail API
type GmailSender struct {
	service *gmail.Service
	from    string
}

// NewGmailSender creates a new Gmail sender
func NewGmailSender() (*GmailSender, error) {
	// Get Gmail credentials from environment (can be file path or JSON string)
	credentialsJSON := os.Getenv("GMAIL_CREDENTIALS_JSON")
	if credentialsJSON == "" {
		return nil, fmt.Errorf("GMAIL_CREDENTIALS_JSON environment variable is not set")
	}

	// Try to read from file if it's a path, otherwise use as JSON string
	var credsData []byte
	if _, err := os.Stat(credentialsJSON); err == nil {
		// It's a file path
		credsData, err = os.ReadFile(credentialsJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to read credentials file: %w", err)
		}
	} else {
		// It's JSON string
		credsData = []byte(credentialsJSON)
	}

	// Try JWT config first (for service accounts)
	config, err := google.JWTConfigFromJSON(credsData, gmail.GmailSendScope)
	if err != nil {
		// If JWT fails, try OAuth2 config
		_, oauthErr := google.ConfigFromJSON(credsData, gmail.GmailSendScope)
		if oauthErr != nil {
			return nil, fmt.Errorf("failed to parse Gmail credentials (tried JWT and OAuth2): %w", err)
		}
		// For OAuth2, we need a token - this is a simplified version
		// In production, you'd need to handle OAuth flow and store refresh token
		return nil, fmt.Errorf("OAuth2 credentials require token refresh flow - use service account with domain-wide delegation instead")
	}

	// Get the email address to impersonate (required for domain-wide delegation)
	impersonateEmail := os.Getenv("GMAIL_FROM")
	if impersonateEmail == "" {
		// Try to get from service account email (if it's a user email)
		impersonateEmail = config.Email
		if impersonateEmail == "" {
			return nil, fmt.Errorf("GMAIL_FROM environment variable must be set for service account with domain-wide delegation")
		}
	}

	// Set the subject (user to impersonate) for domain-wide delegation
	// This is required for service accounts to impersonate users
	config.Subject = impersonateEmail

	// Log configuration for debugging (without sensitive data)
	fmt.Printf("[Gmail] Using service account: %s\n", config.Email)
	fmt.Printf("[Gmail] Impersonating user: %s\n", impersonateEmail)
	fmt.Printf("[Gmail] Scope: %s\n", gmail.GmailSendScope)

	// Create Gmail service with JWT config
	ctx := context.Background()
	client := config.Client(ctx)

	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gmail service: %w", err)
	}

	from := impersonateEmail

	return &GmailSender{
		service: service,
		from:    from,
	}, nil
}

// SendEmail sends an email via Gmail API
func (g *GmailSender) SendEmail(ctx context.Context, req *ports.SendEmailRequest) (*ports.SendEmailResult, error) {
	// Build email message
	message := g.buildMessage(req)

	// Use the from email address as the user (for domain-wide delegation)
	userEmail := g.from

	// Send email
	sentMsg, err := g.service.Users.Messages.Send(userEmail, message).Context(ctx).Do()
	if err != nil {
		errorMsg := err.Error()
		return &ports.SendEmailResult{
			Success: false,
			Error:   &errorMsg,
		}, fmt.Errorf("failed to send email: %w", err)
	}

	return &ports.SendEmailResult{
		MessageID: sentMsg.Id,
		Success:   true,
	}, nil
}

// SendEmailDraft creates an email draft via Gmail API
func (g *GmailSender) SendEmailDraft(ctx context.Context, req *ports.SendEmailRequest) (*ports.SendEmailResult, error) {
	// Build email message
	message := g.buildMessage(req)

	// Use the from email address as the user (for domain-wide delegation)
	userEmail := g.from

	// Create draft
	draft := &gmail.Draft{
		Message: message,
	}

	createdDraft, err := g.service.Users.Drafts.Create(userEmail, draft).Context(ctx).Do()
	if err != nil {
		errorMsg := err.Error()
		return &ports.SendEmailResult{
			Success: false,
			Error:   &errorMsg,
		}, fmt.Errorf("failed to create draft: %w", err)
	}

	return &ports.SendEmailResult{
		MessageID: createdDraft.Message.Id,
		Success:   true,
	}, nil
}

// buildMessage builds a Gmail API message from SendEmailRequest
func (g *GmailSender) buildMessage(req *ports.SendEmailRequest) *gmail.Message {
	from := req.From
	if from == "" {
		from = g.from
	}

	fromName := req.FromName
	if fromName == "" {
		fromName = "BizOps360"
	}

	// Build email headers
	var headers []string
	headers = append(headers, fmt.Sprintf("From: %s <%s>", fromName, from))
	headers = append(headers, fmt.Sprintf("To: %s", req.To))
	headers = append(headers, fmt.Sprintf("Subject: %s", req.Subject))
	headers = append(headers, "MIME-Version: 1.0")

	// Build email body
	var body strings.Builder

	if req.HTMLBody != "" && req.TextBody != "" {
		// Multipart message with both HTML and text
		boundary := "boundary123456789"
		headers = append(headers, fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s", boundary))

		body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		body.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
		body.WriteString(req.TextBody)
		body.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
		body.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
		body.WriteString(req.HTMLBody)
		body.WriteString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	} else if req.HTMLBody != "" {
		headers = append(headers, "Content-Type: text/html; charset=UTF-8")
		body.WriteString(req.HTMLBody)
	} else {
		headers = append(headers, "Content-Type: text/plain; charset=UTF-8")
		body.WriteString(req.TextBody)
	}

	// Combine headers and body
	emailStr := strings.Join(headers, "\r\n") + "\r\n\r\n" + body.String()

	// Encode to base64url
	encoded := base64.URLEncoding.EncodeToString([]byte(emailStr))

	return &gmail.Message{
		Raw: encoded,
	}
}
