package ports

import "context"

// Mailer defines the interface for email sending (Gmail/SMTP)
type Mailer interface {
	SendEmail(ctx context.Context, req *SendEmailRequest) (*SendEmailResult, error)
	SendEmailDraft(ctx context.Context, req *SendEmailRequest) (*SendEmailResult, error)
}

// SendEmailRequest contains email data
type SendEmailRequest struct {
	To          string
	From        string
	FromName    string
	Subject     string
	HTMLBody    string
	TextBody    string
	Attachments []Attachment
}

// Attachment represents an email attachment
type Attachment struct {
	Filename string
	Content  []byte
	MimeType string
}

// SendEmailResult contains the result of sending an email
type SendEmailResult struct {
	MessageID string
	Success   bool
	Error     *string
}

