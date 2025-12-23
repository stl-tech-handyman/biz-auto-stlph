package app

import (
	"context"
	"fmt"

	"github.com/bizops360/go-api/internal/config"
	"github.com/bizops360/go-api/internal/domain"
	"github.com/bizops360/go-api/internal/ports"
)

// FormEventsService handles form event processing
type FormEventsService struct {
	businessLoader *config.BusinessLoader
	pipelineRunner *domain.PipelineRunner
	jobsRepo       ports.JobsRepo
}

// NewFormEventsService creates a new form events service
func NewFormEventsService(
	businessLoader *config.BusinessLoader,
	pipelineRunner *domain.PipelineRunner,
	jobsRepo ports.JobsRepo,
) *FormEventsService {
	return &FormEventsService{
		businessLoader: businessLoader,
		pipelineRunner: pipelineRunner,
		jobsRepo:       jobsRepo,
	}
}

// Run executes a form event pipeline
func (s *FormEventsService) Run(ctx context.Context, req *FormEventsRequest) (*domain.PipelineResult, error) {
	// Load business config
	business, err := s.businessLoader.LoadBusiness(ctx, req.BusinessID)
	if err != nil {
		return nil, fmt.Errorf("failed to load business: %w", err)
	}

	// Resolve pipeline key
	pipelineKey := req.PipelineKey
	if pipelineKey == "" {
		pipelineKey = business.Pipelines.DefaultForm
	}
	if pipelineKey == "" {
		return nil, fmt.Errorf("pipeline key not specified and no default form configured")
	}

	// Load pipeline
	pipeline, err := s.businessLoader.LoadPipeline(ctx, pipelineKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load pipeline: %w", err)
	}

	// Build pipeline context
	pctx := &domain.PipelineContext{
		BusinessID:  req.BusinessID,
		PipelineKey: pipelineKey,
		Source:      req.Source,
		DryRun:      req.DryRun,
		Fields:      req.Fields,
		Options:     req.Options,
		RequestID:   req.RequestID,
	}

	// Run pipeline
	result, job := s.pipelineRunner.Run(ctx, pipeline, pctx)

	// Save job
	if job != nil {
		if err := s.jobsRepo.Save(ctx, job); err != nil {
			// Log error but don't fail the request
			_ = err
		}
	}

	return result, nil
}

// FormEventsRequest represents a form event request
type FormEventsRequest struct {
	BusinessID  string
	PipelineKey string
	Source      string
	DryRun      bool
	Fields      map[string]any
	Options     map[string]any
	RequestID   string
}

