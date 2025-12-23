package stripe

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bizops360/go-api/internal/domain"
	"github.com/bizops360/go-api/internal/ports"
)

// StripePayments implements PaymentsProvider using Stripe API
type StripePayments struct {
	client *http.Client
}

// NewStripePayments creates a new Stripe payments provider
func NewStripePayments() ports.PaymentsProvider {
	return &StripePayments{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// getAPIKey gets Stripe API key from environment based on business config and test mode
func (s *StripePayments) getAPIKey(businessID string, useTest bool) (string, error) {
	// For now, use environment variables directly (can be enhanced to use business config)
	envVar := "STRIPE_SECRET_KEY_PROD"
	if useTest {
		envVar = "STRIPE_SECRET_KEY_TEST"
	}
	
	apiKey := os.Getenv(envVar)
	if apiKey == "" {
		return "", fmt.Errorf("Stripe API key not set. Environment variable '%s' must be configured", envVar)
	}
	return apiKey, nil
}

// CreateInvoice creates a Stripe invoice
func (s *StripePayments) CreateInvoice(ctx context.Context, req *ports.CreateInvoiceRequest) (*ports.InvoiceResult, error) {
	// Determine if test mode (for now, default to production)
	useTest := false // Can be enhanced to read from request or business config
	
	apiKey, err := s.getAPIKey("", useTest)
	if err != nil {
		return nil, err
	}

	// Get or create customer
	customerID, err := s.getOrCreateCustomer(ctx, apiKey, req.CustomerEmail, req.CustomerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create customer: %w", err)
	}

	// Clear pending invoice items
	if err := s.clearPendingInvoiceItems(ctx, apiKey, customerID); err != nil {
		// Log but don't fail
		_ = err
	}

	// Add invoice item
	if err := s.addInvoiceItem(ctx, apiKey, customerID, req.AmountCents, req.Currency, req.Description); err != nil {
		return nil, fmt.Errorf("failed to add invoice item: %w", err)
	}

	// Create invoice
	invoice, err := s.createInvoice(ctx, apiKey, customerID, req.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	// Finalize invoice
	finalized, err := s.finalizeInvoice(ctx, apiKey, invoice.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to finalize invoice: %w", err)
	}

	return &ports.InvoiceResult{
		InvoiceID:       finalized.ID,
		HostedInvoiceURL: finalized.HostedInvoiceURL,
		AmountDue:       finalized.AmountDue,
		Status:          finalized.Status,
		InvoicePDF:     finalized.InvoicePDF,
	}, nil
}

// CalculateDeposit calculates deposit from estimate
func (s *StripePayments) CalculateDeposit(ctx context.Context, estimateTotalCents int64) (*domain.Deposit, error) {
	calc := CalculateDepositFromEstimate(estimateTotalCents)
	return &domain.Deposit{
		AmountCents:        calc.Value,
		AmountDollars:      float64(calc.Value) / 100,
		Percentage:         calc.Percentage,
		EstimateTotalCents: estimateTotalCents,
	}, nil
}

// CreateFinalInvoice creates a final invoice for the remaining balance after deposit
func (s *StripePayments) CreateFinalInvoice(ctx context.Context, req *ports.CreateFinalInvoiceRequest) (*ports.InvoiceResult, error) {
	useTest := false // Can be enhanced to read from request or business config
	
	apiKey, err := s.getAPIKey("", useTest)
	if err != nil {
		return nil, err
	}

	// Calculate remaining balance
	remainingCents := req.TotalAmountCents - req.DepositPaidCents
	if remainingCents <= 0 {
		return nil, fmt.Errorf("no remaining balance: total %d, deposit paid %d", req.TotalAmountCents, req.DepositPaidCents)
	}

	// Get or create customer
	customerID, err := s.getOrCreateCustomer(ctx, apiKey, req.CustomerEmail, req.CustomerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create customer: %w", err)
	}

	// Clear pending invoice items
	if err := s.clearPendingInvoiceItems(ctx, apiKey, customerID); err != nil {
		// Log but don't fail
		_ = err
	}

	// Add invoice item for remaining balance
	description := req.Description
	if description == "" {
		description = fmt.Sprintf("Final Payment - Remaining Balance (Total: $%.2f, Deposit Paid: $%.2f)", 
			float64(req.TotalAmountCents)/100, float64(req.DepositPaidCents)/100)
	}
	
	if err := s.addInvoiceItem(ctx, apiKey, customerID, remainingCents, req.Currency, description); err != nil {
		return nil, fmt.Errorf("failed to add invoice item: %w", err)
	}

	// Add metadata
	metadata := req.Metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}
	metadata["invoice_type"] = "final"
	metadata["total_amount_cents"] = strconv.FormatInt(req.TotalAmountCents, 10)
	metadata["deposit_paid_cents"] = strconv.FormatInt(req.DepositPaidCents, 10)
	metadata["remaining_balance_cents"] = strconv.FormatInt(remainingCents, 10)

	// Create invoice
	invoice, err := s.createInvoice(ctx, apiKey, customerID, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	// Finalize invoice
	finalized, err := s.finalizeInvoice(ctx, apiKey, invoice.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to finalize invoice: %w", err)
	}

	return &ports.InvoiceResult{
		InvoiceID:       finalized.ID,
		HostedInvoiceURL: finalized.HostedInvoiceURL,
		AmountDue:       finalized.AmountDue,
		Status:          finalized.Status,
		InvoicePDF:     finalized.InvoicePDF,
	}, nil
}

// GetInvoice retrieves an invoice by ID
func (s *StripePayments) GetInvoice(ctx context.Context, invoiceID string, useTest bool) (*ports.InvoiceResult, error) {
	apiKey, err := s.getAPIKey("", useTest)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.stripe.com/v1/invoices/"+invoiceID, nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("stripe API error: %s", string(body))
	}

	var invoice Invoice
	if err := json.NewDecoder(resp.Body).Decode(&invoice); err != nil {
		return nil, err
	}

	return &ports.InvoiceResult{
		InvoiceID:       invoice.ID,
		HostedInvoiceURL: invoice.HostedInvoiceURL,
		AmountDue:       invoice.AmountDue,
		Status:          invoice.Status,
		InvoicePDF:     invoice.InvoicePDF,
	}, nil
}

// SendInvoice sends an invoice to the customer via email
func (s *StripePayments) SendInvoice(ctx context.Context, invoiceID string, useTest bool) error {
	apiKey, err := s.getAPIKey("", useTest)
	if err != nil {
		return err
	}

	form := url.Values{}
	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.stripe.com/v1/invoices/"+invoiceID+"/send", strings.NewReader(form.Encode()))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("stripe API error: %s", string(body))
	}

	return nil
}

// getOrCreateCustomer gets or creates a Stripe customer
func (s *StripePayments) getOrCreateCustomer(ctx context.Context, apiKey, email, name string) (string, error) {
	// Try to find existing customer
	if email != "" {
		customers, err := s.listCustomers(ctx, apiKey, email)
		if err == nil && len(customers) > 0 {
			return customers[0].ID, nil
		}
	}

	// Create new customer
	customer, err := s.createCustomer(ctx, apiKey, email, name)
	if err != nil {
		return "", err
	}
	return customer.ID, nil
}

// listCustomers lists customers by email
func (s *StripePayments) listCustomers(ctx context.Context, apiKey, email string) ([]Customer, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.stripe.com/v1/customers", nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	
	q := req.URL.Query()
	q.Set("email", email)
	q.Set("limit", "1")
	req.URL.RawQuery = q.Encode()

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("stripe API error: %s", string(body))
	}

	var result struct {
		Data []Customer `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// createCustomer creates a new customer
func (s *StripePayments) createCustomer(ctx context.Context, apiKey, email, name string) (*Customer, error) {
	form := url.Values{}
	if email != "" {
		form.Set("email", email)
	}
	if name != "" {
		form.Set("name", name)
	}

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.stripe.com/v1/customers", strings.NewReader(form.Encode()))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("stripe API error: %s", string(body))
	}

	var customer Customer
	if err := json.NewDecoder(resp.Body).Decode(&customer); err != nil {
		return nil, err
	}
	return &customer, nil
}

// clearPendingInvoiceItems clears pending invoice items
func (s *StripePayments) clearPendingInvoiceItems(ctx context.Context, apiKey, customerID string) error {
	items, err := s.listInvoiceItems(ctx, apiKey, customerID)
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.Invoice == "" { // Pending item
			_ = s.deleteInvoiceItem(ctx, apiKey, item.ID) // Best effort
		}
	}
	return nil
}

// listInvoiceItems lists invoice items for a customer
func (s *StripePayments) listInvoiceItems(ctx context.Context, apiKey, customerID string) ([]InvoiceItem, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.stripe.com/v1/invoiceitems", nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	
	q := req.URL.Query()
	q.Set("customer", customerID)
	q.Set("limit", "100")
	req.URL.RawQuery = q.Encode()

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("stripe API error: %d", resp.StatusCode)
	}

	var result struct {
		Data []InvoiceItem `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// deleteInvoiceItem deletes an invoice item
func (s *StripePayments) deleteInvoiceItem(ctx context.Context, apiKey, itemID string) error {
	req, _ := http.NewRequestWithContext(ctx, "DELETE", "https://api.stripe.com/v1/invoiceitems/"+itemID, nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("stripe API error: %d", resp.StatusCode)
	}
	return nil
}

// addInvoiceItem adds an invoice item
func (s *StripePayments) addInvoiceItem(ctx context.Context, apiKey, customerID string, amountCents int64, currency, description string) error {
	form := url.Values{}
	form.Set("customer", customerID)
	form.Set("amount", strconv.FormatInt(amountCents, 10))
	form.Set("currency", currency)
	if description != "" {
		form.Set("description", description)
	}

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.stripe.com/v1/invoiceitems", strings.NewReader(form.Encode()))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("stripe API error: %s", string(body))
	}
	return nil
}

// createInvoice creates a draft invoice
func (s *StripePayments) createInvoice(ctx context.Context, apiKey, customerID string, metadata map[string]string) (*Invoice, error) {
	form := url.Values{}
	form.Set("customer", customerID)
	form.Set("collection_method", "send_invoice")
	form.Set("auto_advance", "false")
	form.Set("days_until_due", "7")
	
	for k, v := range metadata {
		form.Set("metadata["+k+"]", v)
	}

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.stripe.com/v1/invoices", strings.NewReader(form.Encode()))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("stripe API error: %s", string(body))
	}

	var invoice Invoice
	if err := json.NewDecoder(resp.Body).Decode(&invoice); err != nil {
		return nil, err
	}
	return &invoice, nil
}

// finalizeInvoice finalizes an invoice
func (s *StripePayments) finalizeInvoice(ctx context.Context, apiKey, invoiceID string) (*Invoice, error) {
	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.stripe.com/v1/invoices/"+invoiceID+"/finalize", nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("stripe API error: %s", string(body))
	}

	var invoice Invoice
	if err := json.NewDecoder(resp.Body).Decode(&invoice); err != nil {
		return nil, err
	}
	return &invoice, nil
}

// Customer represents a Stripe customer
type Customer struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// InvoiceItem represents a Stripe invoice item
type InvoiceItem struct {
	ID      string `json:"id"`
	Invoice string `json:"invoice"`
}

// Invoice represents a Stripe invoice
type Invoice struct {
	ID              string `json:"id"`
	Status          string `json:"status"`
	HostedInvoiceURL string `json:"hosted_invoice_url"`
	InvoicePDF      string `json:"invoice_pdf"`
	AmountDue       int64  `json:"amount_due"`
	Customer        string `json:"customer"`
	CustomerEmail   string `json:"customer_email"`
}

