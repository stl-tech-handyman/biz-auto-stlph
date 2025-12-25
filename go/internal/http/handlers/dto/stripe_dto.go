package dto

// DepositRequest represents a request to create a deposit
type DepositRequest struct {
	Email        string   `json:"email"`
	Name         string   `json:"name"`
	Estimate     *float64 `json:"estimate"`     // Direct estimate value - if provided, uses this to calculate deposit
	DepositValue *float64 `json:"depositValue"`
	Deposit      *float64 `json:"deposit"`
	HelpersCount *int     `json:"helpersCount"`
	Hours        *float64 `json:"hours"`
	UseTest      bool     `json:"useTest"`
	DryRun       bool     `json:"dryRun"`
	MockStripe   bool     `json:"mockStripe"`
}

// DepositCalculateRequest represents a request to calculate deposit
type DepositCalculateRequest struct {
	Estimate *float64 `json:"estimate"`
	Deposit  *float64 `json:"deposit"`
	ShowTable bool    `json:"showTable"`
}

// DepositWithEmailRequest represents a request to create deposit with email
type DepositWithEmailRequest struct {
	Name               string   `json:"name"`
	Email              string   `json:"email"`
	EventType          string   `json:"eventType"`
	EventDateTimeLocal string   `json:"eventDateTimeLocal"`
	EventDate          string   `json:"eventDate"`
	HelpersCount       *int     `json:"helpersCount"`
	Hours              *float64 `json:"hours"`
	Duration           *float64 `json:"duration"`
	Estimate     *float64 `json:"estimate"`     // Direct estimate value - if provided, skips event details calculation and uses this to calculate deposit
	DepositValue *float64 `json:"depositValue"`
	// Memo and Footer with toggles
	Memo            string `json:"memo"`
	ShowMemo        *bool  `json:"showMemo"`        // Toggle to show/hide memo (default: true if memo provided)
	Footer          string `json:"footer"`
	ShowFooter      *bool  `json:"showFooter"`      // Toggle to show/hide footer (default: true if footer provided)
	UseTest           bool   `json:"useTest"`
	DryRun            bool   `json:"dryRun"`
	SaveEmailAsDraft  *bool  `json:"saveEmailAsDraft"`  // If true, email is saved as draft and not sent (default: false - email is sent)
}

// CustomField represents a custom field for Stripe invoices
type CustomField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// FinalInvoiceRequest represents a request to create a final invoice
type FinalInvoiceRequest struct {
	Email            string            `json:"email"`
	Name             string            `json:"name"`
	Estimate         *float64          `json:"estimate"`         // Original estimate/quote amount - if provided, used as original quote
	TotalAmountCents *int64            `json:"totalAmountCents"`
	TotalAmount      *float64          `json:"totalAmount"`      // Actual total event cost
	DepositPaidCents *int64            `json:"depositPaidCents"`
	DepositPaid      *float64          `json:"depositPaid"`      // Deposit already paid (shown in Stripe metadata)
	Currency         string            `json:"currency"`
	Description      string            `json:"description"`
	Metadata         map[string]string `json:"metadata"`
	CustomFields     []CustomField     `json:"customFields"`
	// Fields for extracting custom fields if not explicitly provided
	EventType          string   `json:"eventType"`
	EventDateTimeLocal string   `json:"eventDateTimeLocal"`
	HelpersCount       *int     `json:"helpersCount"`
	Hours              *float64 `json:"hours"`
	Duration           *float64 `json:"duration"`
	// Memo and Footer with toggles
	Memo            string `json:"memo"`
	ShowMemo        *bool  `json:"showMemo"`        // Toggle to show/hide memo (default: true if memo provided)
	Footer          string `json:"footer"`
	ShowFooter      *bool  `json:"showFooter"`      // Toggle to show/hide footer (default: true if footer provided)
	InvoiceType        string  `json:"invoiceType"`        // "final" or "deposit" - used for stamp prefix
	SaveEmailAsDraft   *bool   `json:"saveEmailAsDraft"`  // If true, email is saved as draft and not sent (default: false - email is sent)
	ShowGratuity       *bool   `json:"showGratuity"`      // If true, show gratuity section in final invoice email (default: true)
	UseTemplate        *string `json:"useTemplate"`       // Template name from emltmpl folder (e.g., "invoice", "receipt") - if not set, uses default template
	UseTest            bool    `json:"useTest"`
	SendEmail          bool    `json:"sendEmail"`
}

// TestInvoiceRequest represents a request to test invoice creation
type TestInvoiceRequest struct {
	Email        string   `json:"email"`
	Name         string   `json:"name"`
	Estimate     *float64 `json:"estimate"`     // Direct estimate value - if provided, uses this to calculate deposit
	DepositValue *float64 `json:"depositValue"`
	UseTest      bool     `json:"useTest"`
	SendEmail    bool     `json:"sendEmail"`
}

// DepositAmountRequest represents a request to get deposit amount
type DepositAmountRequest struct {
	Estimate     *float64 `json:"estimate"`     // Direct estimate value - if provided, uses this to calculate deposit
	DepositValue *float64 `json:"depositValue"`
}

