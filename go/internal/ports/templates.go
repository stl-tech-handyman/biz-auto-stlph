package ports

import "context"

// TemplateProvider defines the interface for template rendering
type TemplateProvider interface {
	Render(ctx context.Context, templatePath string, data map[string]any) (string, error)
	RenderHTML(ctx context.Context, templatePath string, data map[string]any) (string, error)
}

