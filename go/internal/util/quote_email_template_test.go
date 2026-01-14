package util

import (
	"strings"
	"testing"
)

// TestQuoteEmailTemplateCompatibility tests email template for cross-client compatibility
func TestQuoteEmailTemplateCompatibility(t *testing.T) {
	data := QuoteEmailData{
		ClientName:     "John Doe",
		EventDate:      "December 25, 2025",
		EventTime:      "6:00 PM",
		EventLocation:  "123 Main St, St. Louis, MO 63110",
		Occasion:       "Birthday Party",
		GuestCount:     50,
		Helpers:        2,
		Hours:          4.0,
		BaseRate:       275.0,
		HourlyRate:     50.0,
		TotalCost:      1100.0,
		DepositAmount:  400.0,
		RateLabel:      "Base Rate",
		ExpirationDate: "December 28, 2025 at 6:00 PM",
		DepositLink:    "https://invoice.stripe.com/i/test",
	}

	html := GenerateQuoteEmailHTML(data, nil)

	// Test 1: All required content is present
	requiredContent := []string{
		"Hi John Doe!",
		"Event Details",
		"EVENT QUOTE",
		"This quote expires in 72 hours",
		"Our Rates & Pricing",
		"Secure Your Date",
		"What Happens Next",
		"Services Included",
		"December 25, 2025",
		"6:00 PM",
		"123 Main St, St. Louis, MO 63110",
		"Birthday Party",
		"$275",
		"$400",
		"$50",
		"$1100",
	}

	for _, content := range requiredContent {
		if !strings.Contains(html, content) {
			t.Errorf("Required content missing: %s", content)
		}
	}

	// Test 2: Email-safe HTML structure (tables, not divs)
	if strings.Contains(html, "<div") {
		t.Error("Template should not use <div> tags - use tables for email compatibility")
	}

	// Test 3: Inline styles only (no external stylesheets or <style> blocks in body)
	if strings.Contains(html, "<style>") {
		t.Error("Template should use inline styles only for email compatibility")
	}

	// Test 4: No CSS classes that might not be supported
	if strings.Contains(html, "class=") {
		t.Error("Template should avoid CSS classes - use inline styles for email compatibility")
	}

	// Test 5: Proper table structure for email clients
	if !strings.Contains(html, "<table") {
		t.Error("Template must use table-based layout for email compatibility")
	}

	// Test 6: All styles are inline
	// Check that style attributes are present on key elements
	styleChecks := []string{
		`style="font-size:`,
		`style="padding:`,
		`style="background-color:`,
		`style="color:`,
	}

	for _, check := range styleChecks {
		if !strings.Contains(html, check) {
			t.Errorf("Missing inline styles - email should have inline styles: %s", check)
		}
	}

	// Test 7: No emoji in critical data fields (they can break in some clients)
	// Emojis in labels are OK, but not in data
	if strings.Contains(html, "📅") || strings.Contains(html, "📍") || strings.Contains(html, "🌐") {
		// Emojis in labels are acceptable, but we should verify they don't break layout
		// This is a warning, not an error
		t.Log("Warning: Template contains emojis which may not render consistently across all email clients")
	}

	// Test 8: Proper encoding
	if !strings.Contains(html, `charset="UTF-8"`) && !strings.Contains(html, `charset=UTF-8`) {
		t.Error("Template must specify UTF-8 encoding")
	}

	// Test 9: Viewport meta tag for mobile
	if !strings.Contains(html, "viewport") {
		t.Error("Template should include viewport meta tag for mobile compatibility")
	}

	// Test 10: No JavaScript
	if strings.Contains(html, "<script") || strings.Contains(html, "javascript:") {
		t.Error("Template must not contain JavaScript - email clients block it")
	}

	// Test 11: All links use http/https (no relative URLs that might break)
	// This is a basic check - in production you'd want more comprehensive link validation
	if strings.Contains(html, `href="//`) {
		t.Error("Template should use absolute URLs (http:// or https://) not protocol-relative URLs")
	}

	// Test 12: Proper parameter substitution (no format errors)
	errorPatterns := []string{
		"%!d",
		"%!s",
		"%!f",
		"%!v",
	}

	for _, pattern := range errorPatterns {
		if strings.Contains(html, pattern) {
			t.Errorf("Template contains format error: %s - check parameter order", pattern)
		}
	}

	// Test 13: Font fallbacks for better compatibility
	if !strings.Contains(html, "Arial") && !strings.Contains(html, "Helvetica") {
		t.Error("Template should specify font fallbacks (Arial, Helvetica, sans-serif)")
	}

	// Test 14: Maximum width constraint for desktop
	if !strings.Contains(html, "max-width:") && !strings.Contains(html, "max-width") {
		t.Error("Template should specify max-width for desktop email clients")
	}
}

// TestQuoteEmailTemplateMobileCompatibility tests mobile-specific compatibility
func TestQuoteEmailTemplateMobileCompatibility(t *testing.T) {
	data := QuoteEmailData{
		ClientName:    "Test User",
		EventDate:     "January 1, 2026",
		EventTime:     "12:00 PM",
		EventLocation: "Test Location",
		Occasion:      "Test Event",
		GuestCount:    25,
		Helpers:       1,
		Hours:         3.0,
		BaseRate:      200.0,
		HourlyRate:    40.0,
		TotalCost:     600.0,
		DepositAmount: 200.0,
		RateLabel:     "Base Rate",
	}

	html := GenerateQuoteEmailHTML(data, nil)

	// Test: Viewport meta tag
	if !strings.Contains(html, "viewport") {
		t.Error("Template must include viewport meta tag for mobile rendering")
	}

	// Test: Responsive table structure
	if !strings.Contains(html, `width="100%"`) {
		t.Error("Template should use width=\"100%\" for responsive tables")
	}

	// Test: Font sizes are readable on mobile (not too small)
	// Check that main text is at least 12px
	if strings.Contains(html, `font-size: 10px`) && !strings.Contains(html, `font-size: 12px`) {
		t.Error("Template should use readable font sizes (minimum 12px) for mobile")
	}
}

// TestQuoteEmailTemplateDataIntegrity tests that all data is correctly inserted
func TestQuoteEmailTemplateDataIntegrity(t *testing.T) {
	data := QuoteEmailData{
		ClientName:    "Jane Smith",
		EventDate:     "March 15, 2026",
		EventTime:     "5:00 PM",
		EventLocation: "456 Oak Avenue, Chicago, IL 60601",
		Occasion:      "Wedding Reception",
		GuestCount:    100,
		Helpers:       4,
		Hours:         6.5,
		BaseRate:      300.0,
		HourlyRate:    60.0,
		TotalCost:     2400.0,
		DepositAmount: 800.0,
		RateLabel:     "Premium Rate",
	}

	html := GenerateQuoteEmailHTML(data, nil)

	// Verify all data fields are correctly inserted
	tests := []struct {
		name    string
		content string
	}{
		{"Client Name", "Jane Smith"},
		{"Event Date", "March 15, 2026"},
		{"Event Time", "5:00 PM"},
		{"Event Location", "456 Oak Avenue, Chicago, IL 60601"},
		{"Occasion", "Wedding Reception"},
		{"Guest Count", "100"},
		{"Helpers", "4"},
		{"Hours", "6.5"},
		{"Base Rate", "$300"},
		{"Deposit Amount", "$800"},
		{"Hourly Rate", "$60"},
		{"Total Cost", "$2400"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(html, tt.content) {
				t.Errorf("Expected content '%s' not found in template", tt.content)
			}
		})
	}

	// Verify no format errors
	if strings.Contains(html, "%!") {
		t.Error("Template contains format errors - check parameter order and types")
	}
}
