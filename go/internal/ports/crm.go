package ports

import "context"

// CRM defines the interface for CRM operations (Monday.com)
type CRM interface {
	CreateDeal(ctx context.Context, req *CreateDealRequest) (*CreateDealResult, error)
	UpdateDeal(ctx context.Context, req *UpdateDealRequest) (*UpdateDealResult, error)
	GetDeal(ctx context.Context, boardID int64, itemID int64) (*Deal, error)
}

// CreateDealRequest contains data to create a deal
type CreateDealRequest struct {
	BoardID int64
	Name    string
	Fields  map[string]any
}

// CreateDealResult contains the result of deal creation
type CreateDealResult struct {
	ItemID  int64
	Success bool
	Error   *string
}

// UpdateDealRequest contains data to update a deal
type UpdateDealRequest struct {
	BoardID int64
	ItemID  int64
	Fields  map[string]any
}

// UpdateDealResult contains the result of deal update
type UpdateDealResult struct {
	Success bool
	Error   *string
}

// Deal represents a CRM deal/item
type Deal struct {
	ItemID int64
	BoardID int64
	Name    string
	Fields  map[string]any
}

