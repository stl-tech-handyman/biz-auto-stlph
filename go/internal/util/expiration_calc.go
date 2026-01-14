package util

import (
	"time"
)

// CalculateUrgencyLevel determines the urgency based on days until the event
func CalculateUrgencyLevel(daysUntilEvent int) string {
	if daysUntilEvent <= 3 {
		return "critical"
	} else if daysUntilEvent <= 7 {
		return "urgent"
	} else if daysUntilEvent <= 14 {
		return "high"
	} else if daysUntilEvent <= 30 {
		return "moderate"
	}
	return "normal"
}

// CalculateExpirationDate calculates the expiration date for a quote based on days until the event
// Rules:
// - Tomorrow (1 day): till midnight of the day when form is filled out
// - 2-3 days: till midnight of the day when form is filled out
// - 4-7 days: 48 hours
// - 8-14 days: 3 days
// - 15-180 days (6 months): 2 weeks
// - >180 days (6+ months): 2 weeks
func CalculateExpirationDate(daysUntilEvent int) (time.Time, string) {
	now := time.Now()
	location, _ := time.LoadLocation("America/Chicago")
	nowInLocation := now.In(location)
	today := time.Date(nowInLocation.Year(), nowInLocation.Month(), nowInLocation.Day(), 0, 0, 0, 0, location)

	var expirationDate time.Time
	var expirationFormatted string

	if daysUntilEvent <= 3 {
		// Tomorrow (1 day) or 2-3 days: expire at midnight of today
		expirationDate = time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 0, location)
		expirationFormatted = expirationDate.Format("Mon, Jan 2, 2006 3:04 PM CST")
	} else if daysUntilEvent <= 7 {
		// 4-7 days: expire in 48 hours
		expirationDate = now.Add(48 * time.Hour)
		expirationFormatted = expirationDate.In(location).Format("Mon, Jan 2, 2006 3:04 PM CST")
	} else if daysUntilEvent <= 14 {
		// 8-14 days: expire in 3 days
		expirationDate = now.Add(72 * time.Hour)
		expirationFormatted = expirationDate.In(location).Format("Mon, Jan 2, 2006 3:04 PM CST")
	} else {
		// 15+ days (including 6 months+): expire in 2 weeks
		expirationDate = now.Add(14 * 24 * time.Hour)
		expirationFormatted = expirationDate.In(location).Format("Mon, Jan 2, 2006 3:04 PM CST")
	}

	return expirationDate, expirationFormatted
}

// GetExpirationMessage returns a human-readable message about when the quote expires
func GetExpirationMessage(daysUntilEvent int) string {
	if daysUntilEvent <= 3 {
		return "This quote expires today at midnight"
	} else if daysUntilEvent <= 7 {
		return "This quote expires in 48 hours"
	} else if daysUntilEvent <= 14 {
		return "This quote expires in 3 days"
	} else {
		return "This quote expires in 2 weeks"
	}
}

// IsDepositNonRefundable returns true if the deposit is non-refundable (< 3 days until event)
func IsDepositNonRefundable(daysUntilEvent int) bool {
	return daysUntilEvent < 3
}

// GetNonRefundableDepositMessage returns the message explaining why deposit is non-refundable
func GetNonRefundableDepositMessage() string {
	return "Deposit is non-refundable. To fairly pay our helpers and maintain high-quality service, we reserve them when you book — which means they lose other opportunities."
}
