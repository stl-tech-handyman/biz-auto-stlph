package ports

import "context"

// Notifier defines the interface for notifications (Slack)
type Notifier interface {
	SendNotification(ctx context.Context, req *NotificationRequest) (*NotificationResult, error)
}

// NotificationRequest contains notification data
type NotificationRequest struct {
	Channel   string
	Message   string
	Title     string
	Color     string // for Slack attachments
	Fields    []NotificationField
}

// NotificationField represents a field in a notification
type NotificationField struct {
	Title string
	Value string
	Short bool
}

// NotificationResult contains the result of sending a notification
type NotificationResult struct {
	Success bool
	Error   *string
}

