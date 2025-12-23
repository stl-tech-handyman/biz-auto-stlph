package ports

import (
	"context"
	"github.com/bizops360/go-api/internal/domain"
)

// BusinessesRepo defines the interface for business configuration storage
type BusinessesRepo interface {
	GetByID(ctx context.Context, id string) (*domain.BusinessConfig, error)
	GetAll(ctx context.Context) ([]*domain.BusinessConfig, error)
}

