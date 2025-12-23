package ports

import (
	"context"
	"github.com/bizops360/go-api/internal/domain"
)

// PipelinesRepo defines the interface for pipeline definition storage
type PipelinesRepo interface {
	GetByKey(ctx context.Context, key string) (*domain.PipelineDefinition, error)
	GetAll(ctx context.Context) ([]*domain.PipelineDefinition, error)
}

