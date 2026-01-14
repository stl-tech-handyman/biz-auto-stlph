package stripe

import (
	"fmt"
	"math"
)

// Professional deposit amounts in cents (multiples of $50)
var professionalDepositAmounts = generateProfessionalAmounts()

func generateProfessionalAmounts() []int64 {
	var amounts []int64
	maxDollars := 5000
	increment := 50
	for dollars := increment; dollars <= maxDollars; dollars += increment {
		amounts = append(amounts, int64(dollars*100))
	}
	return amounts
}

// DepositCalculation represents deposit calculation result
type DepositCalculation struct {
	Value        int64
	Percentage   float64
	MinAmount    int64
	MaxAmount    int64
	TargetAmount int64
	FlooredAmount int64
	PickedBy     string
}

// CalculateDepositFromEstimate calculates deposit from estimate
// Rule: Try to stay between 15-30% and closest to 22.5%
func CalculateDepositFromEstimate(estimateCents int64) DepositCalculation {
	minPercent := 0.15 // 15% (reduced from 25%)
	maxPercent := 0.30 // 30% (reduced from 40%)
	targetPercent := 0.225 // 22.5% (reduced from 32.5%)

	minRange := float64(estimateCents) * minPercent
	maxRange := float64(estimateCents) * maxPercent
	target := float64(estimateCents) * targetPercent

	// Floor to nearest $0.50
	flooredAmount := int64(math.Floor(target/50) * 50)

	// Filter available amounts to those within 15-30% range
	var inRange []int64
	for _, amount := range professionalDepositAmounts {
		if float64(amount) >= minRange && float64(amount) <= maxRange {
			inRange = append(inRange, amount)
		}
	}

	// Use in-range amounts if available, otherwise use all available amounts
	candidates := inRange
	if len(candidates) == 0 {
		candidates = professionalDepositAmounts
	}

	finalAmount := RoundUpToProfessionalAmount(int64(target), candidates)
	percentage := (float64(finalAmount) / float64(estimateCents)) * 100

	return DepositCalculation{
		Value:        finalAmount,
		Percentage:   math.Round(percentage*10) / 10, // Round to 1 decimal
		MinAmount:    int64(minRange),
		MaxAmount:    int64(maxRange),
		TargetAmount: int64(target),
		FlooredAmount: flooredAmount,
		PickedBy:     fmt.Sprintf("calculated_%.1f%%_of_estimate", percentage),
	}
}

// RoundUpToProfessionalAmount rounds up to the next professional deposit amount
func RoundUpToProfessionalAmount(amount int64, candidates []int64) int64 {
	if len(candidates) == 0 {
		return amount
	}

	if amount <= candidates[0] {
		return candidates[0]
	}

	for _, professionalAmount := range candidates {
		if professionalAmount >= amount {
			return professionalAmount
		}
	}

	// Amount exceeds highest professional value
	return candidates[len(candidates)-1]
}

