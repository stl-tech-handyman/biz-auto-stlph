package ports

import (
	"context"
	"github.com/bizops360/go-api/internal/domain"
)

// PaymentsProvider defines the interface for payment processing (Stripe)
type PaymentsProvider interface {
	CreateInvoice(ctx context.Context, req *CreateInvoiceRequest) (*InvoiceResult, error)
	CreateFinalInvoice(ctx context.Context, req *CreateFinalInvoiceRequest) (*InvoiceResult, error)
	CalculateDeposit(ctx context.Context, estimateTotalCents int64) (*domain.Deposit, error)
	GetInvoice(ctx context.Context, invoiceID string, useTest bool) (*InvoiceResult, error)
	SendInvoice(ctx context.Context, invoiceID string, useTest bool) error
}

// CreateInvoiceRequest contains data needed to create an invoice
type CreateInvoiceRequest struct {
	CustomerEmail string
	CustomerName  string
	AmountCents   int64
	Currency      string
	Description   string
	Metadata      map[string]string
}

// CustomField represents a custom field for Stripe invoices
type CustomField struct {
	Name  string
	Value string
}

// CreateFinalInvoiceRequest contains data needed to create a final invoice (remaining balance)
type CreateFinalInvoiceRequest struct {
	CustomerEmail   string
	CustomerName    string
	TotalAmountCents int64  // Total event cost
	DepositPaidCents int64  // Amount already paid as deposit
	Currency        string
	Description     string
	Metadata        map[string]string
	CustomFields    []CustomField
}

// InvoiceResult contains the result of invoice creation
type InvoiceResult struct {
	InvoiceID       string
	HostedInvoiceURL string
	AmountDue       int64
	Status          string
	InvoicePDF     string
}

