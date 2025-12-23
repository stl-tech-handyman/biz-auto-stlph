package app

import (
	"context"

	"github.com/bizops360/go-api/internal/domain"
)

// NormalizeInputAction normalizes input fields
type NormalizeInputAction struct{}

func (a *NormalizeInputAction) Name() string {
	return "normalize_input"
}

func (a *NormalizeInputAction) Execute(ctx context.Context, pctx *domain.PipelineContext) domain.JobStep {
	// Stub implementation - just log and mark as ok
	return domain.JobStep{
		Name:     a.Name(),
		Status:   "ok",
		Critical: false,
		Details: map[string]any{
			"message": "input normalized (stub)",
			"fields":  pctx.Fields,
		},
	}
}

// SendSlackNotificationAction sends a Slack notification
type SendSlackNotificationAction struct{}

func (a *SendSlackNotificationAction) Name() string {
	return "send_slack_notification"
}

func (a *SendSlackNotificationAction) Execute(ctx context.Context, pctx *domain.PipelineContext) domain.JobStep {
	// Stub implementation - just log and mark as skipped for now
	return domain.JobStep{
		Name:     a.Name(),
		Status:   "skipped",
		Critical: false,
		Details: map[string]any{
			"message": "slack notification skipped (stub)",
		},
	}
}

