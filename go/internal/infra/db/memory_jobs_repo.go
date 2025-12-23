package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bizops360/go-api/internal/domain"
	"github.com/bizops360/go-api/internal/ports"
)

// MemoryJobsRepo is an in-memory implementation of JobsRepo
type MemoryJobsRepo struct {
	jobs  map[string]*domain.Job
	mu    sync.RWMutex
}

// NewMemoryJobsRepo creates a new in-memory jobs repository
func NewMemoryJobsRepo() ports.JobsRepo {
	return &MemoryJobsRepo{
		jobs: make(map[string]*domain.Job),
	}
}

// Save saves a job
func (r *MemoryJobsRepo) Save(ctx context.Context, job *domain.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if job.CreatedAt.IsZero() {
		job.CreatedAt = time.Now()
	}
	job.UpdatedAt = time.Now()

	r.jobs[job.ID] = job
	return nil
}

// GetByID retrieves a job by ID
func (r *MemoryJobsRepo) GetByID(ctx context.Context, id string) (*domain.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	job, exists := r.jobs[id]
	if !exists {
		return nil, fmt.Errorf("job not found: %s", id)
	}

	return job, nil
}

// GetByBusinessID retrieves jobs for a business
func (r *MemoryJobsRepo) GetByBusinessID(ctx context.Context, businessID string, limit int) ([]*domain.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var jobs []*domain.Job
	for _, job := range r.jobs {
		if job.BusinessID == businessID {
			jobs = append(jobs, job)
		}
	}

	// Simple limit (in production, would sort by CreatedAt desc)
	if limit > 0 && len(jobs) > limit {
		jobs = jobs[:limit]
	}

	return jobs, nil
}

