package stripe

import (
	"context"
	"fmt"

	"github.com/bizops360/go-api/internal/domain"
	"github.com/bizops360/go-api/internal/ports"
	"github.com/bizops360/go-api/internal/services/pricing"
	"github.com/bizops360/go-api/internal/util"
	"time"
)

// InvoiceService handles invoice-related business logic
type InvoiceService struct {
	paymentsProvider ports.PaymentsProvider
}

// NewInvoiceService creates a new invoice service
func NewInvoiceService(paymentsProvider ports.PaymentsProvider) *InvoiceService {
	return &InvoiceService{
		paymentsProvider: paymentsProvider,
	}
}

// CalculateDepositFromEstimate calculates deposit from an estimate amount
func (s *InvoiceService) CalculateDepositFromEstimate(ctx context.Context, estimateCents int64) (*domain.Deposit, error) {
	return s.paymentsProvider.CalculateDeposit(ctx, estimateCents)
}

// CalculateDepositFromEventDetails calculates deposit from event details
func (s *InvoiceService) CalculateDepositFromEventDetails(ctx context.Context, eventDateStr string, durationHours float64, numHelpers int) (*domain.Deposit, *pricing.EstimateResult, error) {
	eventDate, err := time.Parse("2006-01-02", eventDateStr[:10])
	if err != nil {
		return nil, nil, fmt.Errorf("invalid event date format: %w", err)
	}

	estimateResult, err := pricing.CalculateEstimate(eventDate, durationHours, numHelpers)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to calculate estimate: %w", err)
	}

	estimateCents := util.DollarsToCents(estimateResult.TotalCost)
	deposit, err := s.paymentsProvider.CalculateDeposit(ctx, estimateCents)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to calculate deposit: %w", err)
	}

	return deposit, estimateResult, nil
}

// CreateDepositInvoice creates a deposit invoice
func (s *InvoiceService) CreateDepositInvoice(ctx context.Context, req *CreateDepositInvoiceRequest) (*ports.InvoiceResult, error) {
	// Determine amount
	var amountCents int64
	if req.DepositValueCents != nil {
		amountCents = *req.DepositValueCents
	} else if req.EstimateCents != nil {
		deposit, err := s.paymentsProvider.CalculateDeposit(ctx, *req.EstimateCents)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate deposit: %w", err)
		}
		amountCents = deposit.AmountCents
	} else {
		return nil, fmt.Errorf("either depositValue or estimate is required")
	}

	// Create invoice request
	invoiceReq := &ports.CreateInvoiceRequest{
		CustomerEmail: req.CustomerEmail,
		CustomerName:  req.CustomerName,
		AmountCents:   amountCents,
		Currency:      "usd",
		Description:   req.Description,
		Metadata:      req.Metadata,
	}

	if invoiceReq.Description == "" {
		invoiceReq.Description = "Booking Deposit Invoice"
	}

	return s.paymentsProvider.CreateInvoice(ctx, invoiceReq)
}

// CreateFinalInvoice creates a final invoice for remaining balance
func (s *InvoiceService) CreateFinalInvoice(ctx context.Context, req *CreateFinalInvoiceRequest) (*ports.InvoiceResult, error) {
	// Convert dollars to cents if needed
	var totalCents int64
	if req.TotalAmountCents != nil {
		totalCents = *req.TotalAmountCents
	} else if req.TotalAmount != nil {
		totalCents = util.DollarsToCents(*req.TotalAmount)
	} else {
		return nil, fmt.Errorf("totalAmount or totalAmountCents is required")
	}

	var depositPaidCents int64
	if req.DepositPaidCents != nil {
		depositPaidCents = *req.DepositPaidCents
	} else if req.DepositPaid != nil {
		depositPaidCents = util.DollarsToCents(*req.DepositPaid)
	}

	// Initialize metadata
	metadata := req.Metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}

	// Add deposit info to metadata if provided
	if depositPaidCents > 0 {
		metadata["deposit_paid_cents"] = fmt.Sprintf("%d", depositPaidCents)
		metadata["deposit_paid_dollars"] = fmt.Sprintf("%.2f", util.CentsToDollars(depositPaidCents))
	}

	// Create final invoice request
	finalInvoiceReq := &ports.CreateFinalInvoiceRequest{
		CustomerEmail:    req.CustomerEmail,
		CustomerName:     req.CustomerName,
		TotalAmountCents: totalCents,
		DepositPaidCents: depositPaidCents,
		Currency:         req.Currency,
		Description:      req.Description,
		Metadata:         metadata,
	}

	if finalInvoiceReq.Currency == "" {
		finalInvoiceReq.Currency = "usd"
	}

	return s.paymentsProvider.CreateFinalInvoice(ctx, finalInvoiceReq)
}

// SendInvoice sends an invoice via Stripe
func (s *InvoiceService) SendInvoice(ctx context.Context, invoiceID string, useTest bool) error {
	return s.paymentsProvider.SendInvoice(ctx, invoiceID, useTest)
}

// CreateDepositInvoiceRequest contains data needed to create a deposit invoice
type CreateDepositInvoiceRequest struct {
	CustomerEmail    string
	CustomerName     string
	DepositValueCents *int64
	EstimateCents    *int64
	Description      string
	Metadata         map[string]string
}

// CreateFinalInvoiceRequest contains data needed to create a final invoice
type CreateFinalInvoiceRequest struct {
	CustomerEmail    string
	CustomerName     string
	TotalAmountCents *int64
	TotalAmount      *float64
	DepositPaidCents *int64
	DepositPaid      *float64
	Currency         string
	Description      string
	Metadata         map[string]string
}

