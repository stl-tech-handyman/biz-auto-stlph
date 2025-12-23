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
	RateLabel     string
}

// GenerateQuoteEmailHTML generates the HTML email template for quotes
// Matches the Apps Script generateQuoteEmail function
func GenerateQuoteEmailHTML(data QuoteEmailData) string {
	// Format currency values
	totalFormatted := fmt.Sprintf("$%.0f", data.TotalCost)
	baseRateFormatted := fmt.Sprintf("$%.0f", data.BaseRate)
	hourlyRateFormatted := fmt.Sprintf("$%.0f", data.HourlyRate)

	// Format hours
	hoursFormatted := fmt.Sprintf("%.0f", data.Hours)
	if data.Hours != float64(int(data.Hours)) {
		hoursFormatted = fmt.Sprintf("%.1f", data.Hours)
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <title>STL Party Helpers - Quote</title>
  </head>
  <body style="margin:0; padding:0; font-family: Arial, sans-serif; background-color: #ffffff; color: #333;">
    <table width="100%%" cellpadding="0" cellspacing="2" border="0" style="background-color: #ffffff;">
      <tr>
        <td align="center" style="padding: 0 16px;">
          <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="max-width: 650px; border: 1px solid #ccc; padding: 20px;">
            <!-- Header -->
            <tr>
              <td align="center" style="padding: 5px;">
                <img src="https://stlpartyhelpers.com/wp-content/uploads/2025/08/stlph-logo-2.jpg" alt="STL Party Helpers Logo" />
              </td>
            </tr>
            <tr>
              <td align="center" style="font-size: 22px; font-weight: bold;">Hi %s!</td>
            </tr>
            <tr>
              <td align="center" style="padding-bottom: 10px;">
                Thank you for reaching out!<br />
                Below is your event quote, along with important details and next steps.
              </td>
            </tr>

            <!-- Pricing -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 15px;">Our Rates & Pricing</td>
            </tr>
            <tr>
              <td>
                <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color: #f9f9f9; margin-top: 10px;">
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">Base Rate:</td>
                    <td style="padding: 8px 10px;">%s / helper (first 4 hours)</td>
                  </tr>
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">Additional Hours:</td>
                    <td style="padding: 8px 10px;">%s per additional hour per helper</td>
                  </tr>
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">Estimated Total:</td>
                    <td style="padding: 8px 10px;">%s</td>
                  </tr>
                </table>
                <p style="font-size: 12px; color: #666; padding-top: 5px;">
                  Final total may adjust based on our call. Gratuity is not included but always appreciated!
                </p>
              </td>
            </tr>

            <!-- Event Details -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 20px;">Event Details</td>
            </tr>
            <tr>
              <td>
                <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color: #f9f9f9; margin-top: 10px;">
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">📅 When:</td>
                    <td style="padding: 8px 10px;">%s %s</td>
                  </tr>
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">📍 Where:</td>
                    <td style="padding: 8px 10px;">%s</td>
                  </tr>
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">🌐 Occasion:</td>
                    <td style="padding: 8px 10px;">%s</td>
                  </tr>
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">👥 Guest Count:</td>
                    <td style="padding: 8px 10px;">%d</td>
                  </tr>
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">🧑 Helpers Needed:</td>
                    <td style="padding: 8px 10px;">%d</td>
                  </tr>
                  <tr>
                    <td style="padding: 8px 10px; font-weight: bold;">⏰ For How Long:</td>
                    <td style="padding: 8px 10px;">%s Hours</td>
                  </tr>
                </table>
              </td>
            </tr>

            <!-- Services -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 15px;">Services Included</td>
            </tr>
            <tr>
              <td style="padding: 10px 0;">
                <ul style="padding-left: 20px; margin: 0;">
                  <li><strong>Setup & Presentation</strong>
                    <ul>
                      <li>Arranging tables, chairs, and decorations</li>
                      <li>Buffet setup & live buffet service</li>
                      <li>Butler-passed appetizers & cocktails</li>
                    </ul>
                  </li>
                  <li><strong>Dining & Guest Assistance</strong>
                    <ul>
                      <li>Multi-course plated dinners</li>
                      <li>General bussing (plates, silverware, glassware)</li>
                      <li>Beverage service (water, wine, champagne, coffee, etc.)</li>
                      <li>Special services (cake cutting, dessert plating, etc.)</li>
                    </ul>
                  </li>
                  <li><strong>Cleanup & End-of-Event Support</strong>
                    <ul>
                      <li>Washing dishes, managing trash, and keeping the event space tidy</li>
                      <li>Kitchen cleanup & end-of-event breakdown</li>
                      <li>Assisting with food storage & leftovers</li>
                    </ul>
                  </li>
                </ul>
                <p>Need something specific? Let us know! We'll do our best to accommodate your request.</p>
              </td>
            </tr>

            <!-- Payment Options -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 15px;">Payment Options</td>
            </tr>
            <tr>
              <td style="background-color: #f9f9f9; padding: 10px;">
                Check, Debit / Credit Cards (via Stripe), Venmo, Zelle, Apple Pay
              </td>
            </tr>

            <!-- Next Steps -->
            <tr>
              <td style="font-size: 14px; font-weight: bold; padding-top: 15px;">What Happens Next</td>
            </tr>
            <tr>
              <td style="background-color: #f9f9f9; padding: 10px;">
                <span style="text-align:center; font-size: 14px; font-weight: bold;">Booked already?</span><br />
                <table cellpadding="0" cellspacing="0" border="0" style="font-size: 14px;">
                  <tr>
                    <td valign="top" style="padding-right: 8px;">1.</td>
                    <td>We'll call you at your scheduled time to go over details.</td>
                  </tr>
                  <tr>
                    <td valign="top" style="padding-right: 8px;">2.</td>
                    <td>If all looks good after our call, we'll send a Stripe deposit link to proceed.</td>
                  </tr>
                  <tr>
                    <td valign="top" style="padding-right: 8px;">3.</td>
                    <td>Once the deposit is in, your reservation is locked in.</td>
                  </tr>
                </table>
                <p style="font-size: 13px; text-align: center; color: #666; margin-top: 8px;">
                  Deposit is 40–50%% of the estimate rounded for simplicity.
                </p>
                <p style="font-size: 13px; text-align: center; color: #666; margin-top: 5px;">
                  ❌ Required to confirm your reservation.
                </p>
              </td>
            </tr>
            <tr>
              <td style="background-color: #fff4e5; text-align: center; padding: 10px; margin-top: 5px; border: 1px solid #fddfb4;">
                <strong>Haven't scheduled a call yet?</strong><br />
                <strong>Book now to get started</strong><br />
                <span style="font-size: 0.9em;">(to confirm helpers, tasks, and setup)</span><br />
                <a href="https://calendly.com/stlpartyhelpers/quote-intake" style="display:inline-block; background-color:#0047ab; color:#fff; padding:8px 14px; margin-top: 12px; text-decoration:none; font-weight:bold; border-radius:4px;">Click Here to Schedule Appointment</a>
              </td>
            </tr>

            <!-- Footer -->
            <tr>
              <td align="center" style="font-size: 12px; padding-top: 20px; color: #666;">
                4220 Duncan Ave., Ste. 201, St. Louis, MO 63110<br />
                <a href="tel:+13147145514" style="display:inline-block;background-color:#ffffff;color:#000000;padding: 9px 10px;text-decoration:none;border-radius:4px;margin-top:12px;margin-bottom: 12px;border: 1px solid gray;" target="_blank">Tap to Call Us: (314) 714-5514</a><br />
                <a href="https://stlpartyhelpers.com" style="color:#0047ab; display: inline-block; margin-bottom: 8px;">stlpartyhelpers.com</a><br />
                &copy; 2025 STL Party Helpers<br />
                <span style="font-size: 0.55em;">v1.1</span>
              </td>
            </tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>`,
		data.ClientName,
		baseRateFormatted,
		hourlyRateFormatted,
		totalFormatted,
		data.EventDate,
		data.EventTime,
		data.EventLocation,
		data.Occasion,
		data.GuestCount,
		data.Helpers,
		hoursFormatted,
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

