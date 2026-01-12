package email

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

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
	// #region agent log
	logPath := "c:\\Users\\Alexey\\Code\\biz-operating-system\\stlph\\.cursor\\debug.log"
	if logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		credentialsJSON := os.Getenv("GMAIL_CREDENTIALS_JSON")
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "H1,H2,H3",
			"location":     "gmail.go:NewGmailSender",
			"message":      "Checking GMAIL_CREDENTIALS_JSON env var",
			"data": map[string]interface{}{
				"gmailCredsSet":    credentialsJSON != "",
				"gmailCredsLength": len(credentialsJSON),
				"gmailCredsValue":  credentialsJSON,
			},
			"timestamp": time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion
	// Get Gmail credentials from environment (can be file path or JSON string)
	credentialsJSON := os.Getenv("GMAIL_CREDENTIALS_JSON")
	if credentialsJSON == "" {
		// #region agent log
		if logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			logEntry := map[string]interface{}{
				"sessionId":    "debug-session",
				"runId":        "run1",
				"hypothesisId": "H1,H5",
				"location":     "gmail.go:NewGmailSender",
				"message":      "GMAIL_CREDENTIALS_JSON not set - returning error",
				"data":          map[string]interface{}{},
				"timestamp":     time.Now().UnixMilli(),
			}
			json.NewEncoder(logFile).Encode(logEntry)
			logFile.Close()
		}
		// #endregion
		return nil, fmt.Errorf("GMAIL_CREDENTIALS_JSON environment variable is not set")
	}

	// #region agent log
	if logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		fileInfo, statErr := os.Stat(credentialsJSON)
		fileExists := statErr == nil
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "H2,H3,H4",
			"location":     "gmail.go:NewGmailSender",
			"message":      "Checking if credentials path is a file",
			"data": map[string]interface{}{
				"credentialsPath": credentialsJSON,
				"fileExists":      fileExists,
				"statError":       func() string { if statErr != nil { return statErr.Error() } else { return "" } }(),
				"isFile":          fileExists && !fileInfo.IsDir(),
			},
			"timestamp": time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion
	// Try to read from file if it's a path, otherwise use as JSON string
	var credsData []byte
	if _, err := os.Stat(credentialsJSON); err == nil {
		// It's a file path
		// #region agent log
		if logFile, logErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
			logEntry := map[string]interface{}{
				"sessionId":    "debug-session",
				"runId":        "run1",
				"hypothesisId": "H4",
				"location":     "gmail.go:NewGmailSender",
				"message":      "Attempting to read credentials file",
				"data": map[string]interface{}{
					"credentialsPath": credentialsJSON,
				},
				"timestamp": time.Now().UnixMilli(),
			}
			json.NewEncoder(logFile).Encode(logEntry)
			logFile.Close()
		}
		// #endregion
		credsData, err = os.ReadFile(credentialsJSON)
		if err != nil {
			// #region agent log
			if logFile, logErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
				logEntry := map[string]interface{}{
					"sessionId":    "debug-session",
					"runId":        "run1",
					"hypothesisId": "H4",
					"location":     "gmail.go:NewGmailSender",
					"message":      "Failed to read credentials file",
					"data": map[string]interface{}{
						"credentialsPath": credentialsJSON,
						"readError":       err.Error(),
					},
					"timestamp": time.Now().UnixMilli(),
				}
				json.NewEncoder(logFile).Encode(logEntry)
				logFile.Close()
			}
			// #endregion
			return nil, fmt.Errorf("failed to read credentials file: %w", err)
		}
		// #region agent log
		if logFile, logErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
			logEntry := map[string]interface{}{
				"sessionId":    "debug-session",
				"runId":        "run1",
				"hypothesisId": "H4",
				"location":     "gmail.go:NewGmailSender",
				"message":      "Successfully read credentials file",
				"data": map[string]interface{}{
					"credentialsPath": credentialsJSON,
					"fileSize":        len(credsData),
				},
				"timestamp": time.Now().UnixMilli(),
			}
			json.NewEncoder(logFile).Encode(logEntry)
			logFile.Close()
		}
		// #endregion
	} else {
		// It's JSON string
		// #region agent log
		if logFile, logErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
			logEntry := map[string]interface{}{
				"sessionId":    "debug-session",
				"runId":        "run1",
				"hypothesisId": "H2",
				"location":     "gmail.go:NewGmailSender",
				"message":      "Treating credentials as JSON string (file not found)",
				"data": map[string]interface{}{
					"credentialsPath": credentialsJSON,
					"statError":       err.Error(),
					"jsonLength":      len(credentialsJSON),
				},
				"timestamp": time.Now().UnixMilli(),
			}
			json.NewEncoder(logFile).Encode(logEntry)
			logFile.Close()
		}
		// #endregion
		credsData = []byte(credentialsJSON)
	}

	// #region agent log
	if logFile, logErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "A,B",
			"location":     "gmail.go:NewGmailSender",
			"message":      "Before JWT config creation - checking scope",
			"data": map[string]interface{}{
				"requestedScope": gmail.GmailModifyScope,
				"credsDataLength": len(credsData),
			},
			"timestamp": time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion
	// Try JWT config first (for service accounts)
	// Use GmailModifyScope to allow both sending and creating drafts
	config, err := google.JWTConfigFromJSON(credsData, gmail.GmailModifyScope)
	// #region agent log
	if logFile, logErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "A,B",
			"location":     "gmail.go:NewGmailSender",
			"message":      "After JWT config creation",
			"data": map[string]interface{}{
				"jwtConfigError": func() string { if err != nil { return err.Error() } else { return "" } }(),
				"jwtConfigSuccess": err == nil,
				"serviceAccountEmail": func() string { if config != nil { return config.Email } else { return "" } }(),
			},
			"timestamp": time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion
	if err != nil {
		// If JWT fails, try OAuth2 config
		_, oauthErr := google.ConfigFromJSON(credsData, gmail.GmailModifyScope)
		if oauthErr != nil {
			return nil, fmt.Errorf("failed to parse Gmail credentials (tried JWT and OAuth2): %w", err)
		}
		// For OAuth2, we need a token - this is a simplified version
		// In production, you'd need to handle OAuth flow and store refresh token
		return nil, fmt.Errorf("OAuth2 credentials require token refresh flow - use service account with domain-wide delegation instead")
	}

	// Get the email address to impersonate (required for domain-wide delegation)
	impersonateEmail := os.Getenv("GMAIL_FROM")
	// #region agent log
	if logFile, logErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "C",
			"location":     "gmail.go:NewGmailSender",
			"message":      "Checking impersonation email",
			"data": map[string]interface{}{
				"gmailFromEnv": impersonateEmail,
				"serviceAccountEmail": config.Email,
			},
			"timestamp": time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion
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
	fmt.Printf("[Gmail] Scope: %s\n", gmail.GmailModifyScope)

	// Create Gmail service with JWT config
	ctx := context.Background()
	// #region agent log
	if logFile, logErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "A,B",
			"location":     "gmail.go:NewGmailSender",
			"message":      "Before creating OAuth2 client - will trigger token fetch",
			"data": map[string]interface{}{
				"serviceAccountEmail": config.Email,
				"impersonateEmail":    impersonateEmail,
				"scope":               gmail.GmailModifyScope,
			},
			"timestamp": time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion
	// config.Client(ctx) triggers token fetch - this is where 401 error occurs if scope is not authorized
	client := config.Client(ctx)
	// #region agent log
	if logFile, logErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "A,B",
			"location":     "gmail.go:NewGmailSender",
			"message":      "After creating OAuth2 client - token fetch completed",
			"data": map[string]interface{}{
				"clientCreated": client != nil,
			},
			"timestamp": time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion
	// Note: If token fetch fails here with 401, it means the service account doesn't have
	// GmailModifyScope authorized in Google Workspace Admin Console domain-wide delegation

	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	// #region agent log
	if logFile, logErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "A,B",
			"location":     "gmail.go:NewGmailSender",
			"message":      "After creating Gmail service",
			"data": map[string]interface{}{
				"serviceCreated": service != nil,
				"error":          func() string { if err != nil { return err.Error() } else { return "" } }(),
			},
			"timestamp": time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion
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
	// #region agent log
	logPath := "c:\\Users\\Alexey\\Code\\biz-operating-system\\stlph\\.cursor\\debug.log"
	if logFile, logErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "C",
			"location":     "gmail.go:SendEmail",
			"message":      "Before sending email",
			"data": map[string]interface{}{
				"userEmail": g.from,
				"to":        req.To,
				"subject":   req.Subject,
			},
			"timestamp": time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion
	// Build email message
	message := g.buildMessage(req)

	// Use the from email address as the user (for domain-wide delegation)
	userEmail := g.from

	// #region agent log
	if logFile, logErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "C",
			"location":     "gmail.go:SendEmail",
			"message":      "Before calling Messages.Send API",
			"data": map[string]interface{}{
				"userEmail": userEmail,
			},
			"timestamp": time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion
	// Send email
	sentMsg, err := g.service.Users.Messages.Send(userEmail, message).Context(ctx).Do()
	// #region agent log
	if logFile, logErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
		messageID := ""
		if sentMsg != nil {
			messageID = sentMsg.Id
		}
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "C",
			"location":     "gmail.go:SendEmail",
			"message":      "After calling Messages.Send API",
			"data": map[string]interface{}{
				"error":     func() string { if err != nil { return err.Error() } else { return "" } }(),
				"success":   err == nil,
				"messageID": messageID,
				"to":        req.To,
				"subject":   req.Subject,
			},
			"timestamp": time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion
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
	// #region agent log
	logPath := "c:\\Users\\Alexey\\Code\\biz-operating-system\\stlph\\.cursor\\debug.log"
	if logFile, logErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "A,B",
			"location":     "gmail.go:SendEmailDraft",
			"message":      "Before creating draft",
			"data": map[string]interface{}{
				"userEmail": g.from,
				"to":        req.To,
				"subject":   req.Subject,
			},
			"timestamp": time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion
	// Build email message
	message := g.buildMessage(req)

	// Use the from email address as the user (for domain-wide delegation)
	userEmail := g.from

	// Create draft
	draft := &gmail.Draft{
		Message: message,
	}

	// #region agent log
	if logFile, logErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "A,B",
			"location":     "gmail.go:SendEmailDraft",
			"message":      "Before calling Drafts.Create API",
			"data": map[string]interface{}{
				"userEmail": userEmail,
			},
			"timestamp": time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion
	createdDraft, err := g.service.Users.Drafts.Create(userEmail, draft).Context(ctx).Do()
	// #region agent log
	if logFile, logErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
		draftID := ""
		draftMessageID := ""
		if createdDraft != nil {
			draftID = createdDraft.Id
			if createdDraft.Message != nil {
				draftMessageID = createdDraft.Message.Id
			}
		}
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "A,B",
			"location":     "gmail.go:SendEmailDraft",
			"message":      "After calling Drafts.Create API",
			"data": map[string]interface{}{
				"error":          func() string { if err != nil { return err.Error() } else { return "" } }(),
				"success":        err == nil,
				"draftCreated":   createdDraft != nil,
				"draftID":        draftID,
				"draftMessageID": draftMessageID,
				"to":             req.To,
				"subject":         req.Subject,
			},
			"timestamp": time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion
	if err != nil {
		errorMsg := err.Error()
		// Check if error is due to unauthorized scope
		if strings.Contains(errorMsg, "unauthorized_client") || strings.Contains(errorMsg, "401") {
			enhancedError := fmt.Sprintf("Service account is not authorized for GmailModifyScope. "+
				"Please add 'https://www.googleapis.com/auth/gmail.modify' to domain-wide delegation "+
				"in Google Workspace Admin Console (Security → API Controls → Domain-wide Delegation). "+
				"Find your service account's Client ID and add the scope. Original error: %s", errorMsg)
			return &ports.SendEmailResult{
				Success: false,
				Error:   &enhancedError,
			}, fmt.Errorf("failed to create draft: %w", err)
		}
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

// GetMessage retrieves an email message from Gmail by message ID
func (g *GmailSender) GetMessage(ctx context.Context, messageID string) (string, error) {
	userEmail := g.from
	
	// Get the message with full format to get HTML body
	msg, err := g.service.Users.Messages.Get(userEmail, messageID).Format("full").Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get message: %w", err)
	}
	
	// Extract HTML body from message parts
	var htmlBody string
	if msg.Payload != nil {
		htmlBody = g.extractHTMLFromPayload(msg.Payload)
	}
	
	if htmlBody == "" {
		return "", fmt.Errorf("no HTML body found in message")
	}
	
	return htmlBody, nil
}

// extractHTMLFromPayload recursively extracts HTML content from message payload
func (g *GmailSender) extractHTMLFromPayload(payload *gmail.MessagePart) string {
	if payload == nil {
		return ""
	}
	
	// Check if this part has HTML content
	if payload.MimeType == "text/html" && payload.Body != nil && payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err == nil {
			return string(data)
		}
	}
	
	// Recursively check parts
	if payload.Parts != nil {
		for _, part := range payload.Parts {
			if html := g.extractHTMLFromPayload(part); html != "" {
				return html
			}
		}
	}
	
	return ""
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

// GetMessage retrieves an email message from Gmail by message ID
func (g *GmailSender) GetMessage(ctx context.Context, messageID string) (string, error) {
	userEmail := g.from
	
	// Get the message with full format to get HTML body
	msg, err := g.service.Users.Messages.Get(userEmail, messageID).Format("full").Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get message: %w", err)
	}
	
	// Extract HTML body from message parts
	var htmlBody string
	if msg.Payload != nil {
		htmlBody = g.extractHTMLFromPayload(msg.Payload)
	}
	
	if htmlBody == "" {
		return "", fmt.Errorf("no HTML body found in message")
	}
	
	return htmlBody, nil
}

// extractHTMLFromPayload recursively extracts HTML content from message payload
func (g *GmailSender) extractHTMLFromPayload(payload *gmail.MessagePart) string {
	if payload == nil {
		return ""
	}
	
	// Check if this part has HTML content
	if payload.MimeType == "text/html" && payload.Body != nil && payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err == nil {
			return string(data)
		}
	}
	
	// Recursively check parts
	if payload.Parts != nil {
		for _, part := range payload.Parts {
			if html := g.extractHTMLFromPayload(part); html != "" {
				return html
			}
		}
	}
	
	return ""
}
}

// GetMessage retrieves an email message from Gmail by message ID


// extractHTMLFromPayload recursively extracts HTML content from message payload
func (g *GmailSender) extractHTMLFromPayload(payload *gmail.MessagePart) string {
	if payload == nil {
		return ""
	}
	
	// Check if this part has HTML content
	if payload.MimeType == "text/html" && payload.Body != nil && payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err == nil {
			return string(data)
		}
	}
	
	// Recursively check parts
	if payload.Parts != nil {
		for _, part := range payload.Parts {
			if html := g.extractHTMLFromPayload(part); html != "" {
				return html
			}
		}
	}
	
	return ""
}

// GetMessage retrieves an email message from Gmail by message ID


// extractHTMLFromPayload recursively extracts HTML content from message payload
func (g *GmailSender) extractHTMLFromPayload(payload *gmail.MessagePart) string {
	if payload == nil {
		return ""
	}
	
	// Check if this part has HTML content
	if payload.MimeType == "text/html" && payload.Body != nil && payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err == nil {
			return string(data)
		}
	}
	
	// Recursively check parts
	if payload.Parts != nil {
		for _, part := range payload.Parts {
			if html := g.extractHTMLFromPayload(part); html != "" {
				return html
			}
		}
	}
	
	return ""
}

// GetMessage retrieves an email message from Gmail by message ID


// extractHTMLFromPayload recursively extracts HTML content from message payload
func (g *GmailSender) extractHTMLFromPayload(payload *gmail.MessagePart) string {
	if payload == nil {
		return ""
	}
	
	// Check if this part has HTML content
	if payload.MimeType == "text/html" && payload.Body != nil && payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err == nil {
			return string(data)
		}
	}
	
	// Recursively check parts
	if payload.Parts != nil {
		for _, part := range payload.Parts {
			if html := g.extractHTMLFromPayload(part); html != "" {
				return html
			}
		}
	}
	
	return ""
}

// GetMessage retrieves an email message from Gmail by message ID


// extractHTMLFromPayload recursively extracts HTML content from message payload
func (g *GmailSender) extractHTMLFromPayload(payload *gmail.MessagePart) string {
	if payload == nil {
		return ""
	}
	
	// Check if this part has HTML content
	if payload.MimeType == "text/html" && payload.Body != nil && payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err == nil {
			return string(data)
		}
	}
	
	// Recursively check parts
	if payload.Parts != nil {
		for _, part := range payload.Parts {
			if html := g.extractHTMLFromPayload(part); html != "" {
				return html
			}
		}
	}
	
	return ""
}
// GetMessage retrieves an email message from Gmail by message ID
func (g *GmailSender) GetMessage(ctx context.Context, messageID string) (string, error) {
	userEmail := g.from
	
	// Get the message with full format to get HTML body
	msg, err := g.service.Users.Messages.Get(userEmail, messageID).Format("full").Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get message: %w", err)
	}
	
	// Extract HTML body from message parts
	var htmlBody string
	if msg.Payload != nil {
		htmlBody = g.extractHTMLFromPayload(msg.Payload)
	}
	
	if htmlBody == "" {
		return "", fmt.Errorf("no HTML body found in message")
	}
	
	return htmlBody, nil
}

// extractHTMLFromPayload recursively extracts HTML content from message payload
func (g *GmailSender) extractHTMLFromPayload(payload *gmail.MessagePart) string {
	if payload == nil {
		return ""
	}
	
	// Check if this part has HTML content
	if payload.MimeType == "text/html" && payload.Body != nil && payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err == nil {
			return string(data)
		}
	}
	
	// Recursively check parts
	if payload.Parts != nil {
		for _, part := range payload.Parts {
			if html := g.extractHTMLFromPayload(part); html != "" {
				return html
			}
		}
	}
	
	return ""
}

// GetMessage retrieves an email message from Gmail by message ID
func (g *GmailSender) GetMessage(ctx context.Context, messageID string) (string, error) {
	userEmail := g.from
	
	// Get the message with full format to get HTML body
	msg, err := g.service.Users.Messages.Get(userEmail, messageID).Format("full").Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get message: %w", err)
	}
	
	// Extract HTML body from message parts
	var htmlBody string
	if msg.Payload != nil {
		htmlBody = g.extractHTMLFromPayload(msg.Payload)
	}
	
	if htmlBody == "" {
		return "", fmt.Errorf("no HTML body found in message")
	}
	
	return htmlBody, nil
}

// extractHTMLFromPayload recursively extracts HTML content from message payload
func (g *GmailSender) extractHTMLFromPayload(payload *gmail.MessagePart) string {
	if payload == nil {
		return ""
	}
	
	// Check if this part has HTML content
	if payload.MimeType == "text/html" && payload.Body != nil && payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err == nil {
			return string(data)
		}
	}
	
	// Recursively check parts
	if payload.Parts != nil {
		for _, part := range payload.Parts {
			if html := g.extractHTMLFromPayload(part); html != "" {
				return html
			}
		}
	}
	
	return ""
}

// GetMessage retrieves an email message from Gmail by message ID
func (g *GmailSender) GetMessage(ctx context.Context, messageID string) (string, error) {
	userEmail := g.from
	
	// Get the message with full format to get HTML body
	msg, err := g.service.Users.Messages.Get(userEmail, messageID).Format("full").Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get message: %w", err)
	}
	
	// Extract HTML body from message parts
	var htmlBody string
	if msg.Payload != nil {
		htmlBody = g.extractHTMLFromPayload(msg.Payload)
	}
	
	if htmlBody == "" {
		return "", fmt.Errorf("no HTML body found in message")
	}
	
	return htmlBody, nil
}

// extractHTMLFromPayload recursively extracts HTML content from message payload
func (g *GmailSender) extractHTMLFromPayload(payload *gmail.MessagePart) string {
	if payload == nil {
		return ""
	}
	
	// Check if this part has HTML content
	if payload.MimeType == "text/html" && payload.Body != nil && payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err == nil {
			return string(data)
		}
	}
	
	// Recursively check parts
	if payload.Parts != nil {
		for _, part := range payload.Parts {
			if html := g.extractHTMLFromPayload(part); html != "" {
				return html
			}
		}
	}
	
	return ""
}

// GetMessage retrieves an email message from Gmail by message ID
func (g *GmailSender) GetMessage(ctx context.Context, messageID string) (string, error) {
	userEmail := g.from
	
	// Get the message with full format to get HTML body
	msg, err := g.service.Users.Messages.Get(userEmail, messageID).Format("full").Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get message: %w", err)
	}
	
	// Extract HTML body from message parts
	var htmlBody string
	if msg.Payload != nil {
		htmlBody = g.extractHTMLFromPayload(msg.Payload)
	}
	
	if htmlBody == "" {
		return "", fmt.Errorf("no HTML body found in message")
	}
	
	return htmlBody, nil
}

// extractHTMLFromPayload recursively extracts HTML content from message payload
func (g *GmailSender) extractHTMLFromPayload(payload *gmail.MessagePart) string {
	if payload == nil {
		return ""
	}
	
	// Check if this part has HTML content
	if payload.MimeType == "text/html" && payload.Body != nil && payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err == nil {
			return string(data)
		}
	}
	
	// Recursively check parts
	if payload.Parts != nil {
		for _, part := range payload.Parts {
			if html := g.extractHTMLFromPayload(part); html != "" {
				return html
			}
		}
	}
	
	return ""
}
