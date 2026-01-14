package util

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// QuoteEmailData contains all data needed for quote email
type QuoteEmailData struct {
	ClientName         string
	EventDate          string
	EventTime          string
	EventLocation      string
	Occasion           string
	GuestCount         int
	Helpers            int
	Hours              float64
	BaseRate           float64
	HourlyRate         float64
	TotalCost          float64
	DepositAmount      float64 // Deposit amount in dollars
	RateLabel          string
	ExpirationDate     string // Expiration date formatted (e.g., "June 18, 2026 at 6:00 PM")
	DepositLink        string // Stripe payment link for deposit
	ConfirmationNumber string // 4-character unique confirmation number
	IsHighDemand       bool   // Whether this is a high-demand date (special date/holiday)
	UrgencyLevel       string // Urgency level: "critical" (≤3 days), "urgent" (4-7 days), "high" (8-14 days), "moderate" (15-30 days), "normal" (>30 days)
	DaysUntilEvent     int    // Number of days until the event
	IsReturningClient  bool   // Whether this client has booked with us before
	WeatherForecast     *WeatherForecastData // Weather forecast (only for events < 10 days)
	TravelFeeInfo      *TravelFeeData // Travel fee information (distance, fee, message)
	PDFDownloadLink     string // PDF download link (token-based URL)
}

// TravelFeeData contains travel fee calculation details for email display
type TravelFeeData struct {
	IsWithinServiceArea bool    // True if within 15 miles
	DistanceMiles       float64 // Distance from office in miles
	TravelFee           float64 // Total travel fee (0 if within service area)
	Message             string  // Message to display (e.g., "within our service area - no travel fee")
}

// WeatherForecastData contains weather information for the event
type WeatherForecastData struct {
	Temperature   float64 // in Fahrenheit
	Condition     string  // e.g., "Clear", "Clouds", "Rain"
	Description   string  // e.g., "clear sky", "light rain"
	Recommendation string // Weather-based recommendations
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

// GetFirstNameWithLastInitial formats a name as "First Name L." (e.g., "John D.")
// Handles cases like "John", "John Doe", "John Michael Doe", etc.
func GetFirstNameWithLastInitial(fullName string) string {
	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		return ""
	}
	parts := strings.Fields(fullName)
	if len(parts) == 0 {
		return fullName
	}
	
	firstName := parts[0]
	if len(parts) > 1 {
		// Get last name (last part) and extract first character
		lastName := parts[len(parts)-1]
		if len(lastName) > 0 {
			lastInitial := strings.ToUpper(string(lastName[0]))
			return fmt.Sprintf("%s %s.", firstName, lastInitial)
		}
	}
	return firstName
}

// buildPDFDownloadHTML builds the PDF download link section for the email
func buildPDFDownloadHTML(pdfDownloadLink string, expirationDate string, daysUntilEvent int) string {
	if pdfDownloadLink == "" {
		return ""
	}

	// For same-day bookings, don't show expiration date
	expirationText := ""
	if daysUntilEvent == 0 {
		expirationText = `<p style="margin: 6px 0 0 0; font-size: 11px; color: #991b1b; font-style: italic; font-weight: bold;">Deposit Should Be Paid Now — ASAP to Proceed</p>`
	} else {
		expirationText = fmt.Sprintf(`<p style="margin: 6px 0 0 0; font-size: 11px; color: #999999; font-style: italic;">Expires: %s</p>`, expirationDate)
	}

	return fmt.Sprintf(`            <tr>
              <td style="background-color: #fafafa; padding: 12px; margin-top: 8px; border-left: 3px solid rgb(38, 37, 120); text-align: center;">
                <p style="margin: 0 0 8px 0; font-size: 14px; font-weight: bold; color: rgb(38, 37, 120);">Download PDF Quote</p>
                <p style="margin: 0 0 8px 0; font-size: 12.5px; color: #666666; line-height: 1.5;">
                  Many clients need to submit quotes to accounting for approval. Download your PDF quote below.
                </p>
                <p style="margin: 8px 0; text-align: center;">
                  <a href="%s" style="display: inline-block; background-color: rgb(38, 37, 120); color: #ffffff; padding: 8px 16px; text-decoration: none; font-weight: bold; font-size: 14px; border-radius: 4px;">Download PDF Quote (for accounting/approval)</a>
                </p>
                %s
              </td>
            </tr>`, pdfDownloadLink, expirationText)
}

// FormatDaysUntilEvent formats days until event in a clear, scannable way
// Returns formats like "3 days", "2 weeks", "1 month", "within 3 days", etc.
func FormatDaysUntilEvent(daysUntilEvent int) string {
	if daysUntilEvent < 0 {
		return "past due"
	}
	if daysUntilEvent == 0 {
		return "today"
	}
	if daysUntilEvent == 1 {
		return "tomorrow"
	}
	if daysUntilEvent <= 7 {
		return fmt.Sprintf("%d days", daysUntilEvent)
	}
	if daysUntilEvent <= 14 {
		weeks := daysUntilEvent / 7
		if weeks == 1 {
			return "1 week"
		}
		return fmt.Sprintf("%d weeks", weeks)
	}
	if daysUntilEvent <= 30 {
		weeks := daysUntilEvent / 7
		if weeks == 2 {
			return "2 weeks"
		}
		if weeks == 3 {
			return "3 weeks"
		}
		return fmt.Sprintf("%d weeks", weeks)
	}
	if daysUntilEvent <= 60 {
		months := daysUntilEvent / 30
		if months == 1 {
			return "1 month"
		}
		return fmt.Sprintf("%d months", months)
	}
	months := daysUntilEvent / 30
	if months == 1 {
		return "1 month"
	}
	return fmt.Sprintf("%d months", months)
}

// buildWeatherHTML builds HTML for weather forecast (only shown for events < 10 days)
func buildWeatherHTML(weather *WeatherForecastData) string {
	if weather == nil {
		return ""
	}

	// Format temperature
	tempStr := fmt.Sprintf("%.0f°F", weather.Temperature)

	// Build weather line
	weatherHTML := fmt.Sprintf(`
                <p style="font-size: 13px; color: rgb(38, 37, 120); padding: 8px 30px 0 30px; margin: 0; text-align: center; font-weight: 500;">
                  🌤️ Weather Forecast: %s, %s (%s)
                </p>`, weather.Condition, weather.Description, tempStr)

	// Add recommendation if available
	if weather.Recommendation != "" {
		weatherHTML += fmt.Sprintf(`
                <p style="font-size: 12px; color: #666666; padding: 4px 30px 0 30px; margin: 0; text-align: center; font-style: italic;">
                  %s
                </p>`, weather.Recommendation)
	}

	return weatherHTML
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

	// Format helpers text (singular/plural)
	helpersText := "Helper"
	if data.Helpers != 1 {
		helpersText = "Helpers"
	}

	// Build additional hours description with cost calculation
	// Make it clear: any hour after initial 4 hours = NumHelpers × HourlyRate per hour
	costPerAdditionalHour := float64(data.Helpers) * data.HourlyRate
	helperWord := "helper"
	if data.Helpers != 1 {
		helperWord = "helpers"
	}
	
	additionalHoursText := fmt.Sprintf(
		"Any hour after the initial 4 hours costs %s per helper. With %d %s, each additional hour is %s.",
		hourlyRateFormatted,
		data.Helpers,
		helperWord,
		formatCurrency(costPerAdditionalHour),
	)

	// Calculate end time and build additional hours beyond end time text
	// Additional hours beyond anticipated Staffing Reservation Time: from endTimeHour till 1 AM (max late)
	additionalHoursBeyondEndTimeText := ""
	if data.EventTime != "" && data.Hours > 0 {
		// Parse event time to calculate end time
		eventDateTimeStr := fmt.Sprintf("%s %s", data.EventDate, data.EventTime)
		formats := []string{
			"January 2, 2006 3:04 PM",
			"Jan 2, 2006 3:04 PM",
			"January 2, 2006 15:04",
			"2006-01-02 3:04 PM",
			"2006-01-02 15:04",
		}

		var endTime time.Time
		for _, format := range formats {
			if eventTime, err := time.Parse(format, eventDateTimeStr); err == nil {
				// Calculate end time: event time + duration
				endTime = eventTime.Add(time.Duration(data.Hours) * time.Hour)
				break
			}
		}

		// If we successfully calculated end time
		if !endTime.IsZero() {
			endTimeFormatted := endTime.Format("3:04 PM")
			additionalHoursBeyondEndTimeText = fmt.Sprintf(
				"If our helpers stay longer than anticipated (after %s), we'll add %s for each additional hour. The latest we can extend is 1:00 AM.",
				endTimeFormatted,
				formatCurrency(costPerAdditionalHour),
			)
		}
	}

	// Build HTML for additional hours beyond end time (italicized and centered like other section bottoms)
	additionalHoursBeyondEndTimeHTML := ""
	if additionalHoursBeyondEndTimeText != "" {
		additionalHoursBeyondEndTimeHTML = fmt.Sprintf(
			`                        <tr>
                          <td colspan="2" style="font-size: 13px; padding: 8px 5px; text-align: center; font-style: italic; color: #666666;">%s</td>
                        </tr>`,
			additionalHoursBeyondEndTimeText,
		)
	}

	// Calculate recommended arrival time range (1 hour to 30 minutes before event time)
	// Format as just time range (e.g., "7:00 PM - 7:30 PM") without date
	recommendedArrivalTimeRange := ""
	if data.EventTime != "" {
		// Try to parse event time with date
		eventDateTimeStr := fmt.Sprintf("%s %s", data.EventDate, data.EventTime)
		formats := []string{
			"Mon, January 2, 2006 3:04 PM",    // "Mon, Jan 2, 2006 6:00 PM"
			"Mon, Jan 2, 2006 3:04 PM",        // "Mon, Jan 2, 2006 6:00 PM"
			"January 2, 2006 3:04 PM",         // "January 2, 2006 6:00 PM"
			"Jan 2, 2006 3:04 PM",            // "Jan 2, 2006 6:00 PM"
			"Mon, January 2, 2006 15:04",      // 24-hour format
			"Mon, Jan 2, 2006 15:04",
			"January 2, 2006 15:04",
			"Jan 2, 2006 15:04",
			"2006-01-02 3:04 PM",
			"2006-01-02 15:04",
			"Mon, January 2, 2006 3:04PM",    // Without space
			"Mon, Jan 2, 2006 3:04PM",
			"January 2, 2006 3:04PM",
			"Jan 2, 2006 3:04PM",
		}

		var eventTime time.Time
		parsed := false
		for _, format := range formats {
			if t, err := time.Parse(format, eventDateTimeStr); err == nil {
				eventTime = t
				parsed = true
				break
			}
		}

		// If parsing failed, try parsing just the time part
		if !parsed && data.EventTime != "" {
			timeFormats := []string{
				"3:04 PM",
				"15:04",
				"3:04PM",
				"15:04:00",
			}
			for _, format := range timeFormats {
				if t, err := time.Parse(format, data.EventTime); err == nil {
					// Parse the date separately to get the actual event date
					dateFormats := []string{
						"Mon, January 2, 2006",
						"Mon, Jan 2, 2006",
						"January 2, 2006",
						"Jan 2, 2006",
						"2006-01-02",
					}
					var eventDate time.Time
					dateParsed := false
					for _, dateFormat := range dateFormats {
						if d, err := time.Parse(dateFormat, data.EventDate); err == nil {
							eventDate = d
							dateParsed = true
							break
						}
					}
					
					if dateParsed {
						// Combine parsed date with parsed time
						eventTime = time.Date(eventDate.Year(), eventDate.Month(), eventDate.Day(),
							t.Hour(), t.Minute(), t.Second(), 0, eventDate.Location())
						parsed = true
						break
					} else {
						// Use today's date as fallback
						now := time.Now()
						eventTime = time.Date(now.Year(), now.Month(), now.Day(),
							t.Hour(), t.Minute(), t.Second(), 0, now.Location())
						parsed = true
						break
					}
				}
			}
		}

		if parsed {
			// Calculate 1 hour before (earliest) and 30 minutes before (latest)
			earliestArrival := eventTime.Add(-1 * time.Hour)
			latestArrival := eventTime.Add(-30 * time.Minute)

			// Format as just time range "7:00 PM - 7:30 PM" (no date)
			recommendedArrivalTimeRange = fmt.Sprintf("%s - %s",
				earliestArrival.Format("3:04 PM"),
				latestArrival.Format("3:04 PM"))
		} else {
			// Last resort fallback: try to manually subtract from event time string
			// This should rarely happen, but provides a better fallback than using event time directly
			recommendedArrivalTimeRange = data.EventTime + " (check time calculation)"
		}
	}

	// Build availability meter HTML (shows right below logo)
	// This provides concise availability information to encourage deposit payment
	availabilityMeterHTML := ""
	var availabilityMessage string
	
	// Determine availability message based on urgency level and high demand
	// Special handling for same-day bookings
	if data.DaysUntilEvent == 0 {
		availabilityMessage = "Deposit Should Be Paid Now — ASAP to Proceed with The Staffing Reservation"
	} else {
		switch data.UrgencyLevel {
		case "critical": // ≤3 days
			availabilityMessage = "Limited availability — secure with deposit today"
		case "urgent": // 4-7 days
			availabilityMessage = "Filling fast — secure with deposit to guarantee your spot"
		case "high": // 8-14 days
			availabilityMessage = "Popular date — secure with deposit soon"
		case "moderate": // 15-30 days
			availabilityMessage = "Spots available — secure with deposit to reserve"
		default: // normal (>30 days)
			if data.IsHighDemand {
				availabilityMessage = "Popular date — secure with deposit to guarantee availability"
			} else {
				availabilityMessage = "Secure your date with deposit"
			}
		}
	}
	
	// Always show availability meter
	availabilityMeterHTML = fmt.Sprintf(`            <tr>
              <td align="center" style="background-color: #fff9c4; padding: 8px 6px; margin-bottom: 6px; border-left: 3px solid #f59e0b;">
                <p style="margin: 0; font-size: 14px; font-weight: bold; color: #92400e;">%s</p>
              </td>
            </tr>
`, availabilityMessage)

	// Calculate urgency banner HTML based on urgency level and high demand status
	// Priority: Urgency banner (time-based) takes precedence over high demand (date-based)
	// If both apply, show urgency banner only (it's more time-sensitive)
	urgencyBannerHTML := ""

	// Determine banner message based on urgency level
	// Special handling for same-day bookings
	var bannerMessage string
	if data.DaysUntilEvent == 0 {
		bannerMessage = "Deposit Should Be Paid Now — ASAP to Proceed with The Staffing Reservation"
	} else {
		switch data.UrgencyLevel {
		case "critical": // ≤3 days
			bannerMessage = "Only a few spots left — secure your date today to avoid being left out"
		case "urgent": // 4-7 days
			bannerMessage = "Dates fill up fast — secure your spot now to guarantee availability"
		case "high": // 8-14 days
			bannerMessage = "Popular time period — secure your date soon to guarantee your spot"
		case "moderate": // 15-30 days
			bannerMessage = "Spots are filling up — secure your date to guarantee availability"
		default: // normal (>30 days) or high demand date
			if data.IsHighDemand {
				// For high demand dates that are >30 days away
				bannerMessage = "Popular Date — Dates Fill Up Fast. Book Sooner to Secure Your Spot"
			}
		}
	}

	// Show urgency banner if we have a message
	if bannerMessage != "" {
		urgencyBannerHTML = fmt.Sprintf(`            <tr>
              <td align="center" style="background-color: #f0f0f7; padding: 10px 6px; border-left: 3px solid rgb(38, 37, 120); margin-bottom: 6px;">
                <p style="margin: 0; font-size: 15px; font-weight: bold; color: rgb(38, 37, 120);">%s</p>
              </td>
            </tr>
`, bannerMessage)
	}

	// Build expiration notice HTML (only show for urgent cases)
	expirationNoticeHTML := ""
	expirationMessage := getExpirationMessage(data.UrgencyLevel, data.DaysUntilEvent, data.ExpirationDate)
	if expirationMessage != "" {
		expirationNoticeHTML = fmt.Sprintf(`            <tr>
              <td align="center" style="background-color: #f0f0f7; padding: 10px 6px; border-left: 3px solid rgb(38, 37, 120); margin-bottom: 6px;">
                <p style="margin: 0 0 4px 0; font-size: 15px; font-weight: bold; color: rgb(38, 37, 120);">%s</p>
                <p style="margin: 0; font-size: 12px; color: rgb(38, 37, 120); line-height: 1.4; font-style: italic;">
                  97%% of our clients book us again. Secure your date now.
                </p>
              </td>
            </tr>
`, expirationMessage)
	}

	// Build personalized greeting message
	greetingMessage := ""
	if data.IsReturningClient {
		greetingMessage = `                <strong style="color: rgb(38, 37, 120);">Welcome back! We're thrilled to work with you again.</strong><br /><br />
                My name is Anna, and I am with Customer Success Team.<br /><br />
                We'll handle setup, service, and cleanup so you can enjoy your event.<br /><br />
                Below is your quote with all the details.`
	} else {
		greetingMessage = `                Thank you for your interest in our services!<br /><br />
                My name is Anna, and I am with Customer Success Team.<br /><br />
                We'll handle setup, service, and cleanup so you can enjoy your event.<br /><br />
                Below is your quote with all the details.`
	}

	// Build deposit options message (different for returning vs new clients)
	depositOptionsHTML := ""
	if data.IsReturningClient {
		// For returning clients: optional deposit with peace of mind messaging
		depositOptionsHTML = `<p style="margin: 0 0 8px 0; font-size: 13px; color: rgb(38, 37, 120); font-style: italic;">Your peace of mind is our priority. However you prefer — we'll provide you convenient options.</p>
                <p style="margin: 0 0 8px 0; font-size: 13px; color: #666666;">As a returning client, you can proceed without a deposit if you prefer, or secure your spot with a deposit for added peace of mind.</p>`
	} else {
		// For new clients: standard deposit requirement
		depositOptionsHTML = `<p style="margin: 0 0 8px 0; font-size: 13px; color: rgb(38, 37, 120); font-style: italic;">Your peace of mind is our priority. Secure your spot with a deposit to guarantee your date.</p>`
	}

	// Get color for "(starts in X days)" text based on urgency
	daysUntilEventColor := getDaysUntilEventColor(data.DaysUntilEvent)

	// Get travel fee message and color
	travelFeeMessage := getTravelFeeMessage(data.TravelFeeInfo)
	travelFeeMessageColor := "#666666" // Default gray
	if data.TravelFeeInfo != nil {
		travelFeeMessageColor = getTravelFeeMessageColor(data.TravelFeeInfo.IsWithinServiceArea)
	}

	// Build travel fee row for pricing table (always show, even if $0)
	travelFeeRowHTML := ""
	if data.TravelFeeInfo != nil {
		travelFeeFormatted := formatCurrency(data.TravelFeeInfo.TravelFee)
		travelFeeMessageText := ""
		if data.TravelFeeInfo.TravelFee == 0 {
			if data.TravelFeeInfo.IsWithinServiceArea {
				travelFeeMessageText = " (Within Our Service Radius)"
			} else {
				travelFeeMessageText = " (Within Our Service Area)"
			}
		}
		travelFeeRowHTML = fmt.Sprintf(`                  <tr>
                    <td style="font-size: 14px; padding: 5px; width: 50%%;">- Travel Fee:</td>
                    <td style="font-size: 14px; padding: 5px; width: 50%%;">%s%s</td>
                  </tr>
`, travelFeeFormatted, travelFeeMessageText)
	} else {
		// If travel fee info is not available, show $0 with default message
		travelFeeRowHTML = fmt.Sprintf(`                  <tr>
                    <td style="font-size: 14px; padding: 5px; width: 50%%;">- Travel Fee:</td>
                    <td style="font-size: 14px; padding: 5px; width: 50%%;">$0 (Within Our Service Radius)</td>
                  </tr>
`)
	}

	// Build refund notice HTML (conditional based on days until event)
	refundNoticeHTML := ""
	if IsDepositNonRefundable(data.DaysUntilEvent) {
		// Non-refundable for < 3 days
		refundNoticeHTML = `                <p style="margin: 8px 0 0 0; font-size: 12.5px; color: #991b1b; font-weight: bold;">Deposit is non-refundable.</p>
                <p style="margin: 4px 0 0 0; font-size: 11.5px; color: #666666; font-style: italic;">To fairly pay our helpers and maintain high-quality service, we reserve them when you book — which means they lose other opportunities. Our goal is to have the best helpers for you — to retain them we respect their time and commitment. 98 out of 100 show rate.</p>
`
	} else {
		// Refundable for 3+ days
		refundNoticeHTML = `<p style="margin: 8px 0 0 0; font-size: 12.5px; color: #666666;">100%% refund if cancelled 3+ days before the event.</p>
`
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
            <!-- Status Section with Logo - Office Only Fields -->
            <tr>
              <td style="padding-bottom: 8px; border-bottom: 1px solid #e0e0e0;">
                <table width="100%%" cellpadding="0" cellspacing="0" border="0">
                  <tr>
                    <td align="center" style="padding: 0 10px 0 0; vertical-align: middle;">
                      <img src="https://stlpartyhelpers.com/wp-content/uploads/2025/08/stlph-logo-2.jpg" alt="STL Party Helpers - Professional Event Staffing Services in St. Louis" style="max-width: 120px; height: auto;" />
                    </td>
                    <td style="vertical-align: middle; font-size: 11px; color: #666666; line-height: 1.6;">
                      <strong>Status:</strong> Quote<br />
                      <strong>Is My Reservation Confirmed?:</strong> Awaiting deposit
                    </td>
                  </tr>
                </table>
              </td>
            </tr>
            <!-- Availability Meter - shows availability status right below logo -->
%s
            <!-- QUOTE Header - Make it clear this is a QUOTE, not a reservation -->
            <tr>
              <td align="center" style="padding: 6px; margin-bottom: 10px;">
                <p style="margin: 0; font-size: 14px; font-weight: normal; color: rgb(38, 37, 120);">%s Quote (requested by %s)</p>
                <p style="margin: 3px 0 0 0; font-size: 14.5px; font-weight: bold; color: rgb(38, 37, 120); white-space: nowrap;">Quote ID: %s</p>
                <p style="margin: 3px 0 0 0; font-size: 13.5px; color: rgb(38, 37, 120);">This is a quote to hold your details.</p>
              </td>
            </tr>
            <!-- Urgency Banner (if applicable) - placed above yellow reservation notice -->
%s
            <!-- Yellow Reservation Notice -->
            <tr>
              <td align="center" style="padding: 0 6px 10px 6px;">
                <p style="margin: 0; font-size: 15px; font-weight: bold; color: #92400e; background-color: #fff9c4; padding: 10px 6px; text-align: center; border-left: 3px solid #f59e0b;">Your Reservation is NOT confirmed until deposit is received.</p>
              </td>
            </tr>
            
            <!-- Expiration Notice -->
%s

            <!-- Header -->
            <tr>
              <td align="center" style="font-size: 18px; font-weight: bold; padding: 8px 0;">Hi %s!</td>
            </tr>
            <tr>
              <td align="center" style="padding-bottom: 16px; font-size: 14px; line-height: 1.5;">
%s
              </td>
            </tr>

            <!-- Event Details - AIDA: Interest (build excitement about the event) -->
            <tr>
              <td style="font-size: 16px; font-weight: bold; padding-top: 10px; padding-bottom: 6px; border-top: 1px solid #e0e0e0; text-align: center;">Event Details</td>
            </tr>
            <tr>
              <td align="center" style="background-color: #fafafa; padding: 12px 4px;">
                <table width="540" cellpadding="5" cellspacing="0" border="0" style="background-color: #f9f9f9; width: 540px;">
                  <tr>
                    <td style="font-size: 14px; padding: 5px; width: 50%%; vertical-align: top;">
                      <span style="font-weight: bold;">When:</span><br />
                      %s %s<br />
                      <span style="color: %s; font-size: 13px; font-weight: 500; font-style: italic;">(starts in %s)</span>
                    </td>
                    <td style="font-size: 14px; padding: 5px; width: 50%%; vertical-align: top;">
                      <span style="font-weight: bold;">Where:</span><br />
                      %s<br />
                      <span style="color: %s; font-size: 13px; font-weight: 500;">(%s)</span>
                    </td>
                  </tr>
                  <tr>
                    <td style="font-size: 14px; padding: 5px; width: 50%%;"><span style="font-weight: bold;">Occasion:</span> %s</td>
                    <td style="font-size: 14px; padding: 5px; width: 50%%;"><span style="font-weight: bold;">Guest Count:</span> %d</td>
                  </tr>
                  <tr>
                    <td style="font-size: 14px; padding: 5px; width: 50%%;"><span style="font-weight: bold;">%s:</span> %d</td>
                    <td style="font-size: 14px; padding: 5px; width: 50%%;"><span style="font-weight: bold;">For How Long:</span> %s Hours</td>
                  </tr>
                </table>
                <p style="font-size: 13.5px; color: #666666; padding: 8px 30px 0 30px; margin: 0; text-align: center; font-style: italic;">
                  We advise our staff start time to be between %s to allow for setup and walk-through.
                </p>
%s
              </td>
            </tr>

            <!-- Services Included -->
            <tr>
              <td style="font-size: 16px; font-weight: bold; padding-top: 10px; padding-bottom: 6px; border-top: 1px solid #e0e0e0; text-align: center;">Services Included</td>
            </tr>
            <tr>
              <td align="center" style="background-color: #fafafa; padding: 12px 4px;">
                <table width="540" cellpadding="0" cellspacing="0" border="0" style="width: 540px; font-size: 14px; line-height: 1.5;">
                  <tr>
                    <td style="text-align: left;">
                      <p style="margin: 5px 0; font-weight: bold; font-size: 14px;">Setup & Presentation</p>
                      <ul style="margin: 4px 0; padding-left: 25px;">
                        <li style="margin: 3px 0;">Arranging tables, chairs, and decorations</li>
                        <li style="margin: 3px 0;">Buffet setup & live buffet service</li>
                        <li style="margin: 3px 0;">Butler-passed appetizers & cocktails</li>
                      </ul>
                      <p style="margin: 6px 0 4px 0; font-weight: bold; font-size: 14px;">Dining & Guest Assistance</p>
                      <ul style="margin: 4px 0; padding-left: 25px;">
                        <li style="margin: 3px 0;">Multi-course plated dinners</li>
                        <li style="margin: 3px 0;">General bussing (plates, silverware, glassware)</li>
                        <li style="margin: 3px 0;">Beverage service (water, wine, champagne, coffee, etc.)</li>
                        <li style="margin: 3px 0;">Special services (cake cutting, dessert plating, etc.)</li>
                      </ul>
                      <p style="margin: 6px 0 4px 0; font-weight: bold; font-size: 14px;">Cleanup & End-of-Event Support</p>
                      <ul style="margin: 4px 0; padding-left: 25px;">
                        <li style="margin: 3px 0;">Washing dishes, managing trash, and keeping the event space tidy</li>
                        <li style="margin: 3px 0;">Kitchen cleanup & end-of-event breakdown</li>
                        <li style="margin: 3px 0;">Assisting with food storage & leftovers</li>
                      </ul>
                      <p style="margin: 12px 0 0 0; font-size: 13.5px; color: #666666; line-height: 1.5; text-align: center; font-style: italic;">
                        Need something specific or a small adjustment?<br />
                        Just let us know — we'll do our best to accommodate.
                      </p>
                    </td>
                  </tr>
                </table>
              </td>
            </tr>

            <!-- Rates & Pricing - AIDA: Action (price after value feels fair) -->
            <tr>
              <td style="font-size: 16px; font-weight: bold; padding-top: 10px; padding-bottom: 6px; border-top: 1px solid #e0e0e0; text-align: center;">Our Rates & Pricing</td>
            </tr>
            <tr>
              <td align="center" style="background-color: #fafafa; padding: 12px 4px;">
                <table width="540" cellpadding="5" cellspacing="0" border="0" style="background-color: #f9f9f9; width: 540px;">
                  <!-- Breakdown Section with Border -->
                  <tr>
                    <td colspan="2" style="padding: 0;">
                      <table width="100%%" cellpadding="5" cellspacing="0" border="0" style="border: 1px solid #cccccc; margin: 5px 0; background-color: #ffffff;">
                        <!-- Estimated Total - First Row -->
                        <tr>
                          <td style="font-weight: bold; font-size: 14px; padding: 5px; width: 50%%;">Estimated Total:</td>
                          <td style="font-size: 14px; padding: 5px; width: 50%%;">%s</td>
                        </tr>
                        <tr>
                          <td style="font-size: 14px; padding: 5px; width: 50%%;">- Service Rate:</td>
                          <td style="font-size: 14px; padding: 5px; width: 50%%;">%s / helper (first 4 hours)</td>
                        </tr>
%s
                        <tr>
                          <td colspan="2" style="font-size: 13px; padding: 8px 5px; text-align: center; font-style: italic; color: #666666;">%s</td>
                        </tr>
%s
                      </table>
                    </td>
                  </tr>
                  <!-- Deposit Amount - Separate from breakdown -->
                  <tr>
                    <td style="font-weight: bold; font-size: 14px; padding: 5px; width: 50%%;">Deposit Amount:</td>
                    <td style="font-size: 14px; padding: 5px; width: 50%%;">%s<br /><span style="font-size: 12px; color: #666666; font-style: italic;">%s</span></td>
                  </tr>
                </table>
                <p style="font-size: 12.5px; color: #666666; padding-top: 8px; padding-bottom: 0; margin: 0; text-align: center; font-style: italic;">
                  Final total may adjust only if event details change. Gratuity is not included but always appreciated!
                </p>
              </td>
            </tr>

            <!-- How to Get Started -->
            <tr>
              <td style="font-size: 16px; font-weight: bold; padding-top: 10px; padding-bottom: 6px; border-top: 1px solid #e0e0e0; text-align: center;">How to Get Started</td>
            </tr>
            <tr>
              <td align="center" style="background-color: #fafafa; padding: 12px 4px;">
                <table width="540" cellpadding="0" cellspacing="0" border="0" style="width: 540px; font-size: 14px; line-height: 1.6;">
                  <tr>
                    <td style="text-align: left; padding: 10px;">
                      <ol style="margin: 0; padding-left: 25px; color: #333333;">
                        <li style="margin: 8px 0;"><strong>Pay Your Deposit:</strong> Click the "Pay Deposit" button below to secure your staffing reservation. Your deposit reserves your date and confirms your event.</li>
                        <li style="margin: 8px 0;"><strong>We'll Confirm:</strong> Once your deposit is received, we'll send you a confirmation email with all the details and next steps.</li>
                        <li style="margin: 8px 0;"><strong>Final Payment:</strong> The remaining balance will be due in 7 days (Net 7). You'll receive your final invoice on the next business day after your event (%s) via email with a secure Stripe payment link. You can pay with a check, cash, or securely online via Stripe. This email will also include a separate link to provide gratuity for our helpers if you choose — many customers prefer this option as it keeps the paid reservation record clean and provides flexibility if you don't have cash on hand.</li>
                        <li style="margin: 8px 0;"><strong>Event Day:</strong> Our professional staff will arrive at %s to set up and ensure everything runs smoothly.</li>
                      </ol>
                      <p style="margin: 12px 0 0 0; font-size: 13.5px; color: #666666; line-height: 1.5; text-align: center; font-style: italic;">
                        Questions? We're here to help! Schedule a call or reach out anytime.
                      </p>
                    </td>
                  </tr>
                </table>
              </td>
            </tr>

            <!-- Spacing between pricing and Secure Your Date -->
            <tr>
              <td style="padding-top: 8px;"></td>
            </tr>

            <!-- PDF Download Link -->
%s

            <!-- Secure Your Date - AIDA: Action (CTA right after pricing, capitalize on the moment) -->
            <tr>
              <td style="background-color: #f0f0f7; padding: 10px; margin-top: 8px; border-left: 5px solid rgb(38, 37, 120); border-top: 1px solid #e0e0e0; text-align: center;">
                <p style="margin: 0 0 8px 0; font-size: 17px; font-weight: bold; color: rgb(38, 37, 120);">Ready to Secure Your %s?</p>
%s
                <p style="margin: 0 0 6px 0; text-align: center;">
                  <a href="%s" style="display: inline-block; background-color: rgb(38, 37, 120); color: #ffffff; padding: 6px 12px; text-decoration: none; font-weight: bold; font-size: 14.5px; border-radius: 4px;">Pay Deposit Securely (%s) via Stripe</a>
                </p>
%s
              </td>
            </tr>

            <!-- Questions? -->
            <tr>
              <td style="background-color: #f0f0f7; padding: 10px; margin-top: 20px; border-left: 5px solid rgb(38, 37, 120); border-top: 1px solid #e0e0e0; text-align: center;">
                <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="max-width: 500px; margin: 0 auto;">
                  <tr>
                    <td style="text-align: left;">
                      <p style="margin: 0 0 8px 0; font-size: 17px; font-weight: bold; color: rgb(38, 37, 120); text-align: center;">Questions?</p>
                      <p style="margin: 0 0 8px 0; font-size: 12.5px; line-height: 1.5; text-align: center;">
                        Need a call? Book one if:
                      </p>
                      <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin: 0 auto; border: none;">
                        <tr>
                          <td style="text-align: center; padding: 2px 0; border: none; font-size: 12.5px; line-height: 1.5;">• This is your first time</td>
                        </tr>
                        <tr>
                          <td style="text-align: center; padding: 2px 0; border: none; font-size: 12.5px; line-height: 1.5;">• You have complex food prep</td>
                        </tr>
                        <tr>
                          <td style="text-align: center; padding: 2px 0; border: none; font-size: 12.5px; line-height: 1.5;">• You have china plates</td>
                        </tr>
                        <tr>
                          <td style="text-align: center; padding: 2px 0; border: none; font-size: 12.5px; line-height: 1.5;">• Anything that might affect helper count or hours</td>
                        </tr>
                      </table>
                      <p style="margin: 8px 0; text-align: center;">
                        <a href="%s" style="display: inline-block; background-color: rgb(38, 37, 120); color: #ffffff; padding: 8px 16px; text-decoration: none; font-weight: bold; font-size: 14.5px; border-radius: 4px;">Schedule a Call</a>
                      </p>
                      <p style="margin: 6px 0 0 0; font-size: 12.5px; color: #666666; font-style: italic; text-align: left;">We never oversell — we recommend exactly what you need so your event runs smoothly and you shine. That's why 1,000+ clients trust us and call us back.</p>
                    </td>
                  </tr>
                </table>
              </td>
            </tr>

            <!-- Footer -->
            <tr>
              <td align="center" style="font-size: 12.5px; padding-top: 20px; color: #666666; line-height: 1.5;">
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
		availabilityMeterHTML,        // Availability meter (always shown, right below logo)
		data.Occasion,                // EVENT QUOTE - %s Quote
		GetFirstNameWithLastInitial(data.ClientName), // (requested by First Name L.)
		urgencyBannerHTML,            // Urgency banner (conditional HTML based on time until event)
		data.ConfirmationNumber,       // Quote ID: %s
		expirationNoticeHTML,          // Expiration notice (conditional HTML - only for urgent cases)
		GetFirstName(data.ClientName), // Hi %s!
		greetingMessage,                // Personalized greeting (returning vs new client)
		data.EventDate,                // When: %s (Event Date)
		data.EventTime,                // When: %s (Event Time)
		daysUntilEventColor,           // Color for "(starts in X days)" text
		FormatDaysUntilEvent(data.DaysUntilEvent), // Your event starts %s (after When)
		data.EventLocation,            // Where: %s
		travelFeeMessageColor,         // Color for travel fee message
		travelFeeMessage,              // Travel fee message (e.g., "within our service area - no travel fee")
		data.Occasion,                 // Occasion: %s
		data.GuestCount,               // Guest Count: %d
		helpersText,                   // Helpers/Helper label: %s
		data.Helpers,                  // Helpers count: %d
		hoursFormatted,                // For How Long: %s Hours
		recommendedArrivalTimeRange,   // Recommended arrival time range: %s
		buildWeatherHTML(data.WeatherForecast), // Weather forecast HTML (only for < 10 days)
		totalFormatted,                // Estimated Total: %s (first row)
		baseRateFormatted,             // Service Rate: %s (in breakdown)
		travelFeeRowHTML,              // Travel fee row (in breakdown)
		additionalHoursText,           // Additional Hours description: %s (in breakdown)
		additionalHoursBeyondEndTimeHTML, // Additional Hours beyond end time: %s (in breakdown)
		depositFormatted,              // Deposit Amount: %s (separate)
		getDepositDeadlineMessage(data.DaysUntilEvent), // (due in X days to secure your staffing reservation)
		data.EventDate,                // Final payment due date: %s (the day of your event)
		recommendedArrivalTimeRange,   // Staff arrival time: %s
		buildPDFDownloadHTML(data.PDFDownloadLink, data.ExpirationDate, data.DaysUntilEvent), // PDF download link section
		data.Occasion,                 // Ready to Secure Your %s? (occasion)
		depositOptionsHTML,            // Deposit options message (returning vs new client)
		data.DepositLink,              // Pay Deposit button link
		depositFormatted,              // Pay Deposit button text
		refundNoticeHTML,              // Refund notice (conditional - non-refundable for < 3 days)
		GetBookAppointmentURL(),       // Book appointment link
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

// getExpirationMessage returns the expiration message based on urgency level
func getExpirationMessage(urgencyLevel string, daysUntilEvent int, expirationDate string) string {
	// Special handling for same-day bookings - no expiration date shown
	if daysUntilEvent == 0 {
		return "Deposit Should Be Paid Now — ASAP to Proceed with The Staffing Reservation"
	}
	
	switch urgencyLevel {
	case "critical": // ≤3 days
		// Break into two lines: deadline message and date
		return fmt.Sprintf("Pay deposit by deadline<br />deadline: %s", expirationDate)
	case "urgent": // 4-7 days
		return fmt.Sprintf("This quote expires in 48 hours — pay deposit by deadline<br />deadline: %s", expirationDate)
	case "high": // 8-14 days
		return fmt.Sprintf("This quote expires in 3 days — pay deposit by deadline<br />deadline: %s", expirationDate)
	case "moderate": // 15-30 days
		return fmt.Sprintf("This quote expires in 3 days — pay deposit by deadline<br />deadline: %s", expirationDate)
	default: // normal (>30 days)
		return "" // No message for normal urgency - deposit payment is in the "Ready to Secure" section
	}
}

// getDaysUntilEventColor returns the color code for the "(starts in X days)" text based on urgency
// Uses consistent color scheme across the email template:
// - ≤3 days (critical): #991b1b (dark red - matches non-refundable warning)
// - 4-7 days (urgent): #92400e (dark yellow/brown - matches availability meter)
// - 8-14 days (high): #f59e0b (orange - matches urgent border)
// - 15-30 days (moderate): #666666 (gray - neutral)
// - >30 days (normal): rgb(38, 37, 120) (brand color - default)
func getDaysUntilEventColor(daysUntilEvent int) string {
	if daysUntilEvent < 0 {
		return "#991b1b" // Past due - dark red
	}
	if daysUntilEvent <= 3 {
		return "#991b1b" // Critical (≤3 days) - dark red
	}
	if daysUntilEvent <= 7 {
		return "#92400e" // Urgent (4-7 days) - dark yellow/brown
	}
	if daysUntilEvent <= 14 {
		return "#f59e0b" // High (8-14 days) - orange
	}
	if daysUntilEvent <= 30 {
		return "#666666" // Moderate (15-30 days) - gray
	}
	return "rgb(38, 37, 120)" // Normal (>30 days) - brand color
}

// getTravelFeeMessageColor returns the color code for the travel fee message
// - Within service area: #059669 (green - positive message)
// - Outside service area: #f59e0b (orange - matches urgent border for attention)
func getTravelFeeMessageColor(isWithinServiceArea bool) string {
	if isWithinServiceArea {
		return "#059669" // Green for positive message
	}
	return "#f59e0b" // Orange for attention (outside service area)
}

// getTravelFeeMessage returns the travel fee message to display
// Returns empty string if travel fee info is not available
func getTravelFeeMessage(travelFeeInfo *TravelFeeData) string {
	if travelFeeInfo == nil {
		return "" // No travel fee info available
	}
	return travelFeeInfo.Message
}

// getDepositDeadlineMessage returns the message for deposit deadline based on days until event
// Format: "(due in X days to secure your staffing reservation)"
func getDepositDeadlineMessage(daysUntilEvent int) string {
	if daysUntilEvent <= 3 {
		return "(due today to secure your staffing reservation)"
	} else if daysUntilEvent <= 7 {
		return "(due in 2 days to secure your staffing reservation)"
	} else if daysUntilEvent <= 14 {
		return "(due in 3 days to secure your staffing reservation)"
	} else {
		return "(due in 14 days to secure your staffing reservation)"
	}
}
