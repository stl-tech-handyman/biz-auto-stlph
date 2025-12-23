package ports

import (
	"context"
	"github.com/bizops360/go-api/internal/domain"
)

// JobsRepo defines the interface for job storage
type JobsRepo interface {
	Save(ctx context.Context, job *domain.Job) error
	GetByID(ctx context.Context, id string) (*domain.Job, error)
	GetByBusinessID(ctx context.Context, businessID string, limit int) ([]*domain.Job, error)
}

