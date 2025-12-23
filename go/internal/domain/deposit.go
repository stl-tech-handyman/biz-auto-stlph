package domain

// Deposit represents deposit calculation and payment information
type Deposit struct {
	AmountCents int64  `json:"amountCents"`
	AmountDollars float64 `json:"amountDollars"`
	Percentage  float64  `json:"percentage"`
	EstimateTotalCents int64 `json:"estimateTotalCents,omitempty"`
}

