package util

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// QuoteEmailData contains all data needed for quote email
type QuoteEmailData struct {
	ClientName     string
	EventDate      string
	EventTime      string
	EventLocation  string
	Occasion       string
	GuestCount     int
	Helpers        int
	Hours          float64
	BaseRate       float64
	HourlyRate     float64
	TotalCost      float64
	DepositAmount      float64 // Deposit amount in dollars
	RateLabel          string
	ExpirationDate     string // Expiration date formatted (e.g., "June 18, 2026 at 6:00 PM")
	DepositLink        string // Stripe payment link for deposit
	ConfirmationNumber string // 4-character unique confirmation number
}

// GetBookAppointmentURL returns the book appointment URL from environment variable
// Defaults to https://stlpartyhelpers.com/book-appointment if not set
func GetBookAppointmentURL() string {
	url := os.Getenv("BOOK_APPOINTMENT_URL")
	if url == "" {
		return "https://stlpartyhelpers.com/book-appointment"
	}
	return url
}

// GetFirstName extracts the first name from a full name string
// Handles cases like "John", "John Doe", "John Michael Doe", etc.
func GetFirstName(fullName string) string {
	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		return ""
	}
	parts := strings.Fields(fullName)
	if len(parts) > 0 {
		return parts[0]
	}
	return fullName
}

// GenerateQuoteEmailHTML generates the HTML email template for quotes
// Matches the Apps Script generateQuoteEmail function
func GenerateQuoteEmailHTML(data QuoteEmailData) string {
	// Format currency values - consistent formatting
	formatCurrency := func(amount float64) string {
		if amount == float64(int(amount)) {
			return fmt.Sprintf("$%.0f", amount)
		}
		return fmt.Sprintf("$%.2f", amount)
	}
	
	totalFormatted := formatCurrency(data.TotalCost)
	baseRateFormatted := formatCurrency(data.BaseRate)
	hourlyRateFormatted := formatCurrency(data.HourlyRate)
	depositFormatted := formatCurrency(data.DepositAmount)

	// Format hours
	hoursFormatted := fmt.Sprintf("%.0f", data.Hours)
	if data.Hours != float64(int(data.Hours)) {
		hoursFormatted = fmt.Sprintf("%.1f", data.Hours)
	}

	// Format expiration date as "Dec 27" (short format)
	expirationDateShort := data.ExpirationDate
	if data.ExpirationDate != "" {
		// Try to parse and format as "Jan 2" or "Dec 27"
		formats := []string{
			"January 2, 2006 at 3:04 PM",
			"Jan 2, 2006 at 3:04 PM",
			"January 2, 2006",
			"Jan 2, 2006",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, data.ExpirationDate); err == nil {
				expirationDateShort = t.Format("Jan 2")
				break
			}
		}
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>STL Party Helpers - Quote</title>
  </head>
  <body style="margin:0; padding:0; font-family: Arial, Helvetica, sans-serif; background-color: #ffffff; color: #333333;">
    <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color: #ffffff; width: 100%%;">
      <tr>
        <td align="center" style="padding: 8px;">
          <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="max-width: 600px; border: 1px solid #cccccc; padding: 12px; background-color: #ffffff;">
            <!-- QUOTE Header - Make it clear this is a QUOTE, not a reservation -->
            <tr>
              <td align="center" style="padding: 6px; margin-bottom: 10px;">
                <p style="margin: 0; font-size: 16px; font-weight: bold; color: rgb(38, 37, 120);">%s Quote</p>
                <p style="margin: 3px 0 0 0; font-size: 13px; font-weight: bold; color: rgb(38, 37, 120);">Quote ID: %s</p>
                <p style="margin: 3px 0 0 0; font-size: 12px; color: rgb(38, 37, 120);">This is a quote, not a confirmed reservation.</p>
                <p style="margin: 3px 0 0 0; font-size: 12px; color: rgb(38, 37, 120);">Your reservation is confirmed only after deposit payment.</p>
              </td>
            </tr>
            
            <!-- Expiration Notice -->
            <tr>
              <td align="center" style="background-color: #fff5e6; padding: 10px 6px; border-left: 3px solid #d97706; margin-bottom: 6px;">
                <p style="margin: 0; font-size: 13px; font-weight: bold; color: #d97706;">This quote expires in 72 hours / 3 days</p>
                <p style="margin: 2px 0 0 0; font-size: 11px; color: #d97706;">Valid until: %s</p>
                <p style="margin: 3px 0 0 0; font-size: 11px; color: #d97706;">
                  <a href="%s" style="color: #d97706; text-decoration: underline;">You can pay deposit (%s) until %s to secure your reservation.</a>
                </p>
              </td>
            </tr>

            <!-- Header -->
            <tr>
              <td align="center" style="font-size: 16px; font-weight: bold; padding: 8px 0;">Hi %s!</td>
            </tr>
            <tr>
              <td align="center" style="padding-bottom: 8px; font-size: 12px; line-height: 1.5;">
                Thank you for reaching out!<br />
                Below is your event quote, along with important details and next steps.<br />
                <span style="font-size: 11px; color: #666666; font-style: italic;">A PDF copy of this quote is attached for your records.</span>
              </td>
            </tr>

            <!-- Event Details - AIDA: Interest (build excitement about the event) -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 10px; padding-bottom: 6px; border-top: 1px solid #e0e0e0; text-align: center;">Event Details</td>
            </tr>
            <tr>
              <td>
                <table width="100%%" cellpadding="5" cellspacing="0" border="0" style="background-color: #f9f9f9; width: 100%%;">
                  <tr>
                    <td style="font-size: 12px; padding: 5px; width: 50%%; vertical-align: top;">
                      <span style="font-weight: bold;">When:</span><br />
                      %s %s
                    </td>
                    <td style="font-size: 12px; padding: 5px; width: 50%%; vertical-align: top;">
                      <span style="font-weight: bold;">Where:</span><br />
                      %s
                    </td>
                  </tr>
                  <tr>
                    <td style="font-size: 12px; padding: 5px; width: 50%%;"><span style="font-weight: bold;">Occasion:</span> %s</td>
                    <td style="font-size: 12px; padding: 5px; width: 50%%;"><span style="font-weight: bold;">Guest Count:</span> %d</td>
                  </tr>
                  <tr>
                    <td style="font-size: 12px; padding: 5px; width: 50%%;"><span style="font-weight: bold;">Helpers:</span> %d</td>
                    <td style="font-size: 12px; padding: 5px; width: 50%%;"><span style="font-weight: bold;">For How Long:</span> %s Hours</td>
                  </tr>
                </table>
              </td>
            </tr>

            <!-- Rates & Pricing - AIDA: Action (price after value feels fair) -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 10px; padding-bottom: 6px; border-top: 1px solid #e0e0e0; text-align: center;">Our Rates & Pricing</td>
            </tr>
            <tr>
              <td style="padding-top: 4px;">
                <table width="100%%" cellpadding="5" cellspacing="0" border="0" style="background-color: #f9f9f9; width: 100%%;">
                  <tr>
                    <td style="font-weight: bold; font-size: 12px; padding: 5px; width: 50%%;">Base Rate:</td>
                    <td style="font-size: 12px; padding: 5px; width: 50%%;">%s / helper (first 4 hours)</td>
                  </tr>
                  <tr>
                    <td style="font-weight: bold; font-size: 12px; padding: 5px; width: 50%%;">Additional Hours:</td>
                    <td style="font-size: 12px; padding: 5px; width: 50%%;">%s per hour per helper</td>
                  </tr>
                  <tr>
                    <td style="font-weight: bold; font-size: 12px; padding: 5px; width: 50%%;">Deposit Amount:</td>
                    <td style="font-size: 12px; padding: 5px; width: 50%%;">%s</td>
                  </tr>
                  <tr>
                    <td style="font-weight: bold; font-size: 12px; padding: 5px; width: 50%%;">Estimated Total:</td>
                    <td style="font-size: 12px; padding: 5px; width: 50%%;">%s</td>
                  </tr>
                </table>
                <p style="font-size: 11px; color: #666666; padding-top: 8px; padding-bottom: 0; margin: 0; text-align: center;">
                  Final total may adjust based on our call. Gratuity is not included but always appreciated!
                </p>
              </td>
            </tr>

            <!-- Spacing between pricing and Secure Your Date -->
            <tr>
              <td style="padding-top: 8px;"></td>
            </tr>

            <!-- Secure Your Date - AIDA: Action (CTA right after pricing, capitalize on the moment) -->
            <tr>
              <td style="background-color: #f0f0f7; padding: 10px; margin-top: 8px; border-left: 5px solid rgb(38, 37, 120); border-top: 1px solid #e0e0e0; text-align: center;">
                <p style="margin: 0 0 8px 0; font-size: 13px; font-weight: bold; color: rgb(38, 37, 120);">Ready to Secure Your Event?</p>
                <p style="margin: 0 0 6px 0; font-size: 12px; color: #333333; line-height: 1.6;">Make a deposit to confirm your reservation and lock in your date.</p>
                <p style="margin: 0 0 6px 0; text-align: center;">
                  <a href="%s" style="display: inline-block; background-color: rgb(38, 37, 120); color: #ffffff; padding: 6px 12px; text-decoration: none; font-weight: bold; font-size: 13px; border-radius: 4px;">Pay Deposit (%s) via Stripe</a>
                </p>
                <p style="margin: 8px 0 0 0; font-size: 11px; color: #666666;">100%% refund if cancelled 3+ days before the event.</p>
              </td>
            </tr>

            <!-- Next Steps - AIDA: Action (clear path forward) -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 10px; padding-bottom: 6px; border-top: 1px solid #e0e0e0; text-align: center;">Schedule an Appointment</td>
            </tr>
            <tr>
              <td style="padding: 4px 0 5px 0; font-size: 12px; line-height: 1.5; text-align: center;">
                <p style="margin: 0;">If you want to speak with us, <a href="%s" style="color: rgb(38, 37, 120); text-decoration: underline; font-weight: bold;">book an appointment</a> <span style="font-size: 10px;">(unless you already booked one)</span>.</p>
              </td>
            </tr>
            <tr>
              <td style="padding: 8px 0 5px 0; font-size: 12px; line-height: 1.5; text-align: center;">
                <p style="margin: 0;">Need something specific?</p>
                <p style="margin: 0;">Let us know!</p>
                <p style="margin: 0;">We'll do our best to accommodate your request.</p>
              </td>
            </tr>

            <!-- Services Included -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 10px; padding-bottom: 6px; border-top: 1px solid #e0e0e0; text-align: center;">Services Included</td>
            </tr>
            <tr>
              <td align="center" style="padding: 4px 0 5px 0;">
                <table width="50%%" cellpadding="0" cellspacing="0" border="0" style="width: 50%%; margin: 0 auto;">
                  <tr>
                    <td style="font-size: 12px; line-height: 1.5; text-align: left;">
                      <p style="margin: 5px 0; font-weight: bold; font-size: 12px;">Setup & Presentation</p>
                      <ul style="margin: 4px 0; padding-left: 25px;">
                        <li style="margin: 3px 0;">Arranging tables, chairs, and decorations</li>
                        <li style="margin: 3px 0;">Buffet setup & live buffet service</li>
                        <li style="margin: 3px 0;">Butler-passed appetizers & cocktails</li>
                      </ul>
                      <p style="margin: 6px 0 4px 0; font-weight: bold; font-size: 12px;">Dining & Guest Assistance</p>
                      <ul style="margin: 4px 0; padding-left: 25px;">
                        <li style="margin: 3px 0;">Multi-course plated dinners</li>
                        <li style="margin: 3px 0;">General bussing (plates, silverware, glassware)</li>
                        <li style="margin: 3px 0;">Beverage service (water, wine, champagne, coffee, etc.)</li>
                        <li style="margin: 3px 0;">Special services (cake cutting, dessert plating, etc.)</li>
                      </ul>
                      <p style="margin: 6px 0 4px 0; font-weight: bold; font-size: 12px;">Cleanup & End-of-Event Support</p>
                      <ul style="margin: 4px 0; padding-left: 25px;">
                        <li style="margin: 3px 0;">Washing dishes, managing trash, and keeping the event space tidy</li>
                        <li style="margin: 3px 0;">Kitchen cleanup & end-of-event breakdown</li>
                        <li style="margin: 3px 0;">Assisting with food storage & leftovers</li>
                      </ul>
                    </td>
                  </tr>
                </table>
              </td>
            </tr>

            <!-- Footer -->
            <tr>
              <td align="center" style="font-size: 11px; padding-top: 10px; color: #666666; line-height: 1.5;">
                <img src="https://stlpartyhelpers.com/wp-content/uploads/2025/08/stlph-logo-2.jpg" alt="STL Party Helpers - Professional Event Staffing Services in St. Louis" style="max-width: 150px; height: auto; margin-bottom: 8px;" /><br />
                4220 Duncan Ave., Ste. 201, St. Louis, MO 63110<br />
                <a href="tel:+13147145514" style="color: #000000; text-decoration: underline; display: inline-block; margin: 5px 0;">Tap to Call Us: (314) 714-5514</a><br />
                <a href="https://stlpartyhelpers.com" style="color: rgb(38, 37, 120); text-decoration: underline; display: inline-block; margin: 3px 0;">stlpartyhelpers.com</a><br />
                &copy; 2025 STL Party Helpers<br />
                <span style="font-size: 10px;">v1.1</span>
              </td>
            </tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>`,
		data.Occasion,           // EVENT QUOTE - %s Quote
		data.ConfirmationNumber, // Quote ID: %s
		data.ExpirationDate,     // Expiration notice - Valid until
		data.DepositLink,    // Expiration notice - deposit link
		depositFormatted,    // Expiration notice - deposit amount
		expirationDateShort, // Expiration notice - short date format
		GetFirstName(data.ClientName),     // Hi %s!
		data.EventDate,      // When: %s
		data.EventTime,      // When: %s (second)
		data.EventLocation, // Where: %s
		data.Occasion,       // Occasion: %s
		data.GuestCount,     // Guest Count: %d
		data.Helpers,        // Helpers: %d
		hoursFormatted,      // For How Long: %s Hours
		baseRateFormatted,   // Base Rate: %s
		hourlyRateFormatted, // Additional Hours: %s
		depositFormatted,    // Deposit Amount: %s
		totalFormatted,      // Estimated Total: %s
		data.DepositLink,    // Pay Deposit button link
		depositFormatted,    // Pay Deposit button text
		GetBookAppointmentURL(), // Book appointment link
	)

	return html
}

// FormatEventDate formats event date for display
func FormatEventDate(dateStr string) string {
	// If already formatted, return as-is
	if strings.Contains(dateStr, ",") {
		return dateStr
	}
	// Otherwise, try to parse and format
	// For now, return as-is - caller should provide formatted date
	return dateStr
}

