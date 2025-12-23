package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/bizops360/go-api/internal/domain"
	"gopkg.in/yaml.v3"
)

// BusinessLoader loads and caches business configurations from YAML files
type BusinessLoader struct {
	config     *Config
	cache     map[string]*domain.BusinessConfig
	cacheMu   sync.RWMutex
	pipelines map[string]*domain.PipelineDefinition
	pipelinesMu sync.RWMutex
}

// NewBusinessLoader creates a new business loader
func NewBusinessLoader(cfg *Config) *BusinessLoader {
	return &BusinessLoader{
		config:    cfg,
		cache:     make(map[string]*domain.BusinessConfig),
		pipelines: make(map[string]*domain.PipelineDefinition),
	}
}

// LoadBusiness loads a business configuration by ID
func (bl *BusinessLoader) LoadBusiness(ctx context.Context, businessID string) (*domain.BusinessConfig, error) {
	// Check cache first
	bl.cacheMu.RLock()
	if cached, exists := bl.cache[businessID]; exists {
		bl.cacheMu.RUnlock()
		return cached, nil
	}
	bl.cacheMu.RUnlock()

	// Load from file
	path := bl.config.GetBusinessConfigPath(businessID)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read business config %s: %w", businessID, err)
	}

	var business domain.BusinessConfig
	if err := yaml.Unmarshal(data, &business); err != nil {
		return nil, fmt.Errorf("failed to parse business config %s: %w", businessID, err)
	}

	// Validate required fields
	if business.ID == "" {
		business.ID = businessID
	}

	// Cache it
	bl.cacheMu.Lock()
	bl.cache[businessID] = &business
	bl.cacheMu.Unlock()

	return &business, nil
}

// LoadPipeline loads a pipeline definition by key
func (bl *BusinessLoader) LoadPipeline(ctx context.Context, pipelineKey string) (*domain.PipelineDefinition, error) {
	// Check cache first
	bl.pipelinesMu.RLock()
	if cached, exists := bl.pipelines[pipelineKey]; exists {
		bl.pipelinesMu.RUnlock()
		return cached, nil
	}
	bl.pipelinesMu.RUnlock()

	// Load from file
	path := bl.config.GetPipelineConfigPath(pipelineKey)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read pipeline config %s: %w", pipelineKey, err)
	}

	var pipeline domain.PipelineDefinition
	if err := yaml.Unmarshal(data, &pipeline); err != nil {
		return nil, fmt.Errorf("failed to parse pipeline config %s: %w", pipelineKey, err)
	}

	// Validate required fields
	if pipeline.Key == "" {
		pipeline.Key = pipelineKey
	}

	// Cache it
	bl.pipelinesMu.Lock()
	bl.pipelines[pipelineKey] = &pipeline
	bl.pipelinesMu.Unlock()

	return &pipeline, nil
}

// LoadAllBusinesses loads all business configurations from the config directory
func (bl *BusinessLoader) LoadAllBusinesses(ctx context.Context) ([]*domain.BusinessConfig, error) {
	businessesDir := filepath.Join(bl.config.ConfigDir, "businesses")
	entries, err := os.ReadDir(businessesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read businesses directory: %w", err)
	}

	var businesses []*domain.BusinessConfig
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}

		businessID := entry.Name()[:len(entry.Name())-5] // Remove .yaml extension
		business, err := bl.LoadBusiness(ctx, businessID)
		if err != nil {
			// Log error but continue loading others
			continue
		}
		businesses = append(businesses, business)
	}

	return businesses, nil
}

