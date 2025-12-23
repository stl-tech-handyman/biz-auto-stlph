package dto

// DepositRequest represents a request to create a deposit
type DepositRequest struct {
	Email          string   `json:"email"`
	Name           string   `json:"name"`
	EstimatedTotal *float64 `json:"estimatedTotal"`
	DepositValue   *float64 `json:"depositValue"`
	Deposit        *float64 `json:"deposit"`
	HelpersCount   *int     `json:"helpersCount"`
	Hours          *float64 `json:"hours"`
	UseTest        bool     `json:"useTest"`
	DryRun         bool     `json:"dryRun"`
	MockStripe     bool     `json:"mockStripe"`
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
	Estimate           *float64 `json:"estimate"`
	EstimatedTotal     *float64 `json:"estimatedTotal"`
	DepositValue       *float64 `json:"depositValue"`
	UseTest            bool     `json:"useTest"`
	DryRun             bool     `json:"dryRun"`
	SaveAsDraft        bool     `json:"saveAsDraft"`
}

// FinalInvoiceRequest represents a request to create a final invoice
type FinalInvoiceRequest struct {
	Email            string            `json:"email"`
	Name             string            `json:"name"`
	TotalAmountCents *int64            `json:"totalAmountCents"`
	TotalAmount      *float64          `json:"totalAmount"`
	DepositPaidCents *int64            `json:"depositPaidCents"`
	DepositPaid      *float64          `json:"depositPaid"`
	Currency         string            `json:"currency"`
	Description      string            `json:"description"`
	Metadata         map[string]string `json:"metadata"`
	UseTest          bool              `json:"useTest"`
	SendEmail        bool              `json:"sendEmail"`
}

// TestInvoiceRequest represents a request to test invoice creation
type TestInvoiceRequest struct {
	Email          string   `json:"email"`
	Name           string   `json:"name"`
	EstimatedTotal *float64 `json:"estimatedTotal"`
	DepositValue   *float64 `json:"depositValue"`
	UseTest        bool     `json:"useTest"`
	SendEmail      bool     `json:"sendEmail"`
}

// DepositAmountRequest represents a request to get deposit amount
type DepositAmountRequest struct {
	EstimatedTotal *float64 `json:"estimatedTotal"`
	Estimate       *float64 `json:"estimate"`
	DepositValue   *float64 `json:"depositValue"`
}

