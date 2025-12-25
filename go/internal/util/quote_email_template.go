package util

import (
	"fmt"
	"strings"
)

// QuoteEmailData contains all data needed for quote email
type QuoteEmailData struct {
	ClientName    string
	EventDate     string
	EventTime     string
	EventLocation string
	Occasion      string
	GuestCount    int
	Helpers       int
	Hours         float64
	BaseRate      float64
	HourlyRate    float64
	TotalCost     float64
	DepositAmount float64 // Deposit amount in dollars
	RateLabel     string
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
            <!-- Header -->
            <tr>
              <td align="center" style="font-size: 16px; font-weight: bold; padding: 8px 0;">Hi %s!</td>
            </tr>
            <tr>
              <td align="center" style="padding-bottom: 8px; font-size: 12px; line-height: 1.5;">
                Thank you for reaching out!<br />
                Below is your event quote, along with important details and next steps.
              </td>
            </tr>

            <!-- Event Details - AIDA: Interest (build excitement about the event) -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 10px; padding-bottom: 6px; border-top: 1px solid #e0e0e0;">Event Details</td>
            </tr>
            <tr>
              <td>
                <table width="100%%" cellpadding="5" cellspacing="0" border="0" style="background-color: #f9f9f9; width: 100%%;">
                  <tr>
                    <td style="font-size: 12px; padding: 5px; width: 50%%;"><span style="font-weight: bold;">When:</span> %s %s</td>
                    <td style="font-size: 12px; padding: 5px; width: 50%%;"><span style="font-weight: bold;">Where:</span> %s</td>
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

            <!-- Services Included - AIDA: Desire (show value before price) -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 10px; padding-bottom: 6px; border-top: 1px solid #e0e0e0;">Services Included</td>
            </tr>
            <tr>
              <td style="padding: 4px 0 5px 0; font-size: 12px; line-height: 1.5;">
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
                <p style="margin-top: 6px; font-size: 12px;">Need something specific? Let us know! We'll do our best to accommodate your request.</p>
              </td>
            </tr>

            <!-- Rates & Pricing - AIDA: Action (price after value feels fair) -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 10px; padding-bottom: 6px; border-top: 1px solid #e0e0e0;">Our Rates & Pricing</td>
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
                <p style="font-size: 11px; color: #666666; padding-top: 8px; padding-bottom: 0; margin: 0;">
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
              <td style="background-color: #f0f7ff; padding: 10px; margin-top: 8px; border-left: 5px solid #0047ab; border-top: 1px solid #e0e0e0;">
                <p style="margin: 0 0 6px 0; font-size: 13px; font-weight: bold; color: #0047ab;">Secure Your Date</p>
                <p style="margin: 0; font-size: 12px; color: #333333; line-height: 1.6;">Paying the deposit (%s) locks in your event date and confirms your reservation. You can always cancel and receive a 100%% full refund as long as it's done at least 3 days prior to the start of your event.</p>
              </td>
            </tr>

            <!-- Payment Options - AIDA: Action (reduce friction, show it's easy) -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 10px; padding-bottom: 6px; border-top: 1px solid #e0e0e0;">Payment Options</td>
            </tr>
            <tr>
              <td style="background-color: #f9f9f9; padding: 8px; font-size: 12px; margin-top: 4px;">
                Check, Debit / Credit Cards (via Stripe), Venmo, Zelle
              </td>
            </tr>

            <!-- Next Steps - AIDA: Action (clear path forward) -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 10px; padding-bottom: 6px; border-top: 1px solid #e0e0e0;">What Happens Next</td>
            </tr>
            <tr>
              <td style="padding: 4px 0 5px 0; font-size: 12px; line-height: 1.5;">
                <p style="margin: 0 0 5px 0;">You can <a href="https://calendly.com/stlpartyhelpers/quote-intake" style="color: #0047ab; text-decoration: underline; font-weight: bold;">book an appointment with us</a> to learn more, or respond to this email if you are ready to proceed.</p>
                <p style="margin: 0 0 5px 0;">If you already booked an appointment, we will talk to you soon!</p>
                <p style="margin: 0;">If you are ready to proceed, we will send you a deposit link. Once the deposit is paid, your reservation is confirmed and your date is secured.</p>
              </td>
            </tr>

            <!-- Footer -->
            <tr>
              <td align="center" style="font-size: 11px; padding-top: 10px; color: #666666; line-height: 1.5;">
                <img src="https://stlpartyhelpers.com/wp-content/uploads/2025/08/stlph-logo-2.jpg" alt="STL Party Helpers - Professional Event Staffing Services in St. Louis" style="max-width: 150px; height: auto; margin-bottom: 8px;" /><br />
                4220 Duncan Ave., Ste. 201, St. Louis, MO 63110<br />
                <a href="tel:+13147145514" style="color: #000000; text-decoration: underline; display: inline-block; margin: 5px 0;">Tap to Call Us: (314) 714-5514</a><br />
                <a href="https://stlpartyhelpers.com" style="color: #0047ab; text-decoration: underline; display: inline-block; margin: 3px 0;">stlpartyhelpers.com</a><br />
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
		data.ClientName,
		data.EventDate,
		data.EventTime,
		data.EventLocation,
		data.Occasion,
		data.GuestCount,
		data.Helpers,
		hoursFormatted,
		baseRateFormatted,
		hourlyRateFormatted,
		depositFormatted,
		totalFormatted,
		depositFormatted, // For "Secure Your Date" section (right after pricing)
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

