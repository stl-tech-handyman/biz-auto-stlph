package util

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bizops360/go-api/internal/domain"
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
	ExpirationDate     string               // Expiration date formatted (e.g., "June 18, 2026 at 6:00 PM")
	DepositLink        string               // Stripe payment link for deposit
	ConfirmationNumber string               // 4-character unique confirmation number
	IsHighDemand       bool                 // Whether this is a high-demand date (special date/holiday)
	UrgencyLevel       string               // Urgency level: "critical" (≤3 days), "urgent" (4-7 days), "high" (8-14 days), "moderate" (15-30 days), "normal" (>30 days)
	DaysUntilEvent     int                  // Number of days until the event
	IsReturningClient  bool                 // Whether this client has booked with us before
	WeatherForecast    *WeatherForecastData // Weather forecast (only for events < 10 days)
	TravelFeeInfo      *TravelFeeData       // Travel fee information (distance, fee, message)
	PDFDownloadLink    string               // PDF download link (token-based URL)
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
	Temperature    float64 // in Fahrenheit
	Condition      string  // e.g., "Clear", "Clouds", "Rain"
	Description    string  // e.g., "clear sky", "light rain"
	Recommendation string  // Weather-based recommendations
}

// GetContactInfo returns contact information from business config with smart defaults
// Uses convention over configuration - derives values from business ID and Gmail config if not explicitly set
func GetContactInfo(businessID string, businessConfig *domain.BusinessConfig) domain.ContactConfig {
	contact := domain.ContactConfig{}

	if businessConfig != nil {
		contact = businessConfig.Contact
	}

	// Smart defaults using convention over configuration
	if contact.SupportEmail == "" {
		// Try to extract from Gmail sender, or use support@{businessID}.com
		if businessConfig != nil && businessConfig.Gmail.Sender != "" {
			contact.SupportEmail = businessConfig.Gmail.Sender
		} else if gmailSender := os.Getenv("GMAIL_FROM"); gmailSender != "" {
			contact.SupportEmail = gmailSender
		} else {
			contact.SupportEmail = fmt.Sprintf("support@%s.com", businessID)
		}
	}

	if contact.WebsiteURL == "" {
		// Derive from Gmail sender domain or use https://{businessID}.com
		if businessConfig != nil && businessConfig.Gmail.Sender != "" {
			if parts := strings.Split(businessConfig.Gmail.Sender, "@"); len(parts) == 2 {
				contact.WebsiteURL = fmt.Sprintf("https://%s", parts[1])
			} else {
				contact.WebsiteURL = fmt.Sprintf("https://%s.com", businessID)
			}
		} else if gmailSender := os.Getenv("GMAIL_FROM"); gmailSender != "" {
			if parts := strings.Split(gmailSender, "@"); len(parts) == 2 {
				contact.WebsiteURL = fmt.Sprintf("https://%s", parts[1])
			} else {
				contact.WebsiteURL = fmt.Sprintf("https://%s.com", businessID)
			}
		} else {
			contact.WebsiteURL = fmt.Sprintf("https://%s.com", businessID)
		}
	}

	if contact.LogoURL == "" {
		// Default logo path convention
		contact.LogoURL = fmt.Sprintf("%s/wp-content/uploads/logo.jpg", contact.WebsiteURL)
	}

	if contact.BookAppointmentURL == "" {
		contact.BookAppointmentURL = fmt.Sprintf("%s/book-appointment", contact.WebsiteURL)
	}

	return contact
}

// GetBusinessDisplayName returns display name from config or falls back to business ID
func GetBusinessDisplayName(businessID string, businessConfig *domain.BusinessConfig) string {
	if businessConfig != nil && businessConfig.DisplayName != "" {
		return businessConfig.DisplayName
	}
	return businessID
}

// GetBusinessOfficeAddress returns office address from config or empty string
func GetBusinessOfficeAddress(businessConfig *domain.BusinessConfig) string {
	if businessConfig != nil && businessConfig.Location.OfficeAddress != "" {
		return businessConfig.Location.OfficeAddress
	}
	return ""
}

// GetBookAppointmentURL returns the book appointment URL (deprecated - use GetContactInfo instead)
// Defaults to https://stlpartyhelpers.com/book-appointment if not set
func GetBookAppointmentURL() string {
	contact := GetContactInfo("stlpartyhelpers", nil)
	return contact.BookAppointmentURL
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
		expirationText = `<p style="margin: 6px 0 0 0; font-size: 10.5px; color: #991b1b; font-style: italic; font-weight: bold;">Deposit Should Be Paid Now — ASAP to Proceed</p>`
	} else {
		expirationText = fmt.Sprintf(`<p style="margin: 6px 0 0 0; font-size: 10.5px; color: #999999; font-style: italic;">Expires: %s</p>`, expirationDate)
	}

	return fmt.Sprintf(`            <tr>
              <td style="background-color: #fafafa; padding: 12px; margin-top: 8px; border-left: 3px solid rgb(38, 37, 120); text-align: center;">
                <p style="margin: 0 0 8px 0; font-size: 10.5px; font-weight: bold; color: rgb(38, 37, 120);">Download PDF Quote</p>
                <p style="margin: 0 0 8px 0; font-size: 10.5px; color: #666666; line-height: 1.5;">
                  Many clients need to submit quotes to accounting for approval. Download your PDF quote below.
                </p>
                <p style="margin: 8px 0; text-align: center;">
                  <a href="%s" style="display: inline-block; background-color: rgb(38, 37, 120); color: #ffffff; padding: 8px 16px; text-decoration: none; font-weight: bold; font-size: 10.5px; border-radius: 4px;">Download PDF Quote (for accounting/approval)</a>
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
                <p style="font-size: 10.5px; color: rgb(38, 37, 120); padding: 8px 30px 0 30px; margin: 0; text-align: center; font-weight: 500;">
                  🌤️ Weather Forecast: %s, %s (%s)
                </p>`, weather.Condition, weather.Description, tempStr)

	// Add recommendation if available
	if weather.Recommendation != "" {
		weatherHTML += fmt.Sprintf(`
                <p style="font-size: 10.5px; color: #666666; padding: 4px 30px 0 30px; margin: 0; text-align: center; font-style: italic;">
                  %s
                </p>`, weather.Recommendation)
	}

	return weatherHTML
}

// GenerateQuoteEmailHTML generates the HTML email template for quotes
// Matches the Apps Script generateQuoteEmail function
// businessConfig is optional - if nil, uses smart defaults based on business ID
func GenerateQuoteEmailHTML(data QuoteEmailData, businessConfig *domain.BusinessConfig) string {
	// Get contact information with smart defaults
	businessID := "stlpartyhelpers" // Default business ID
	if businessConfig != nil && businessConfig.ID != "" {
		businessID = businessConfig.ID
	}
	contact := GetContactInfo(businessID, businessConfig)
	displayName := GetBusinessDisplayName(businessID, businessConfig)
	officeAddress := GetBusinessOfficeAddress(businessConfig)

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
		"Any hour after the initial 4 hours costs %s per helper.<br />With %d %s, each additional hour is %s.",
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
                          <td colspan="2" style="font-size: 10.5px; padding: 8px 5px; text-align: center; font-style: italic; color: #666666;">%s</td>
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
			"Mon, January 2, 2006 3:04 PM", // "Mon, Jan 2, 2006 6:00 PM"
			"Mon, Jan 2, 2006 3:04 PM",     // "Mon, Jan 2, 2006 6:00 PM"
			"January 2, 2006 3:04 PM",      // "January 2, 2006 6:00 PM"
			"Jan 2, 2006 3:04 PM",          // "Jan 2, 2006 6:00 PM"
			"Mon, January 2, 2006 15:04",   // 24-hour format
			"Mon, Jan 2, 2006 15:04",
			"January 2, 2006 15:04",
			"Jan 2, 2006 15:04",
			"2006-01-02 3:04 PM",
			"2006-01-02 15:04",
			"Mon, January 2, 2006 3:04PM", // Without space
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
                <p style="margin: 0; font-size: 10.5px; font-weight: bold; color: #92400e;">%s</p>
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
                <p style="margin: 0; font-size: 10.5px; font-weight: bold; color: rgb(38, 37, 120);">%s</p>
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
                <p style="margin: 0 0 4px 0; font-size: 10.5px; font-weight: bold; color: rgb(38, 37, 120);">%s</p>
                <p style="margin: 0; font-size: 10.5px; color: rgb(38, 37, 120); line-height: 1.4; font-style: italic;">
                  97%% of our clients book us again. Secure your date now.
                </p>
              </td>
            </tr>
`, expirationMessage)
	}

	// Build personalized greeting message
	greetingMessage := ""
	if data.IsReturningClient {
		greetingMessage = `                <strong style="color: rgb(38, 37, 120);">Welcome back! We're thrilled to work with you again.</strong><br />
                My name is Anna, and I am with Customer Success Team.<br />
                We'll handle setup, service, and cleanup so you can enjoy your event.<br />
                Below is your quote with all the details.`
	} else {
		greetingMessage = `                Thank you for your interest in our services!<br />
                My name is Anna, and I am with Customer Success Team.<br />
                We'll handle setup, service, and cleanup so you can enjoy your event.<br />
                Below is your quote with all the details.`
	}

	// Build deposit options message (different for returning vs new clients)
	depositOptionsHTML := ""
	if data.IsReturningClient {
		// For returning clients: optional deposit with peace of mind messaging
		depositOptionsHTML = `<p style="margin: 0 0 8px 0; font-size: 10.5px; color: rgb(38, 37, 120); font-style: italic;">Your peace of mind is our priority. However you prefer — we'll provide you convenient options.</p>
                <p style="margin: 0 0 8px 0; font-size: 10.5px; color: #666666;">As a returning client, you can proceed without a deposit if you prefer, or secure your spot with a deposit for added peace of mind.</p>`
	} else {
		// For new clients: standard deposit requirement
		depositOptionsHTML = `<p style="margin: 0 0 8px 0; font-size: 10.5px; color: rgb(38, 37, 120); font-style: italic;">Your peace of mind is our priority. Secure your spot with a deposit to guarantee your date.</p>`
	}

	// Get color for "(starts in X days)" text based on urgency
	daysUntilEventColor := getDaysUntilEventColor(data.DaysUntilEvent)

	// Get travel fee message and color
	travelFeeMessage := getTravelFeeMessage(data.TravelFeeInfo)
	travelFeeMessageColor := "#666666" // Default gray
	if data.TravelFeeInfo != nil {
		travelFeeMessageColor = getTravelFeeMessageColor(data.TravelFeeInfo.IsWithinServiceArea)
	}

	// Get service radius for travel fee message
	serviceRadiusMiles := 15.0 // Default
	if businessConfig != nil && businessConfig.Location.ServiceRadiusMiles > 0 {
		serviceRadiusMiles = businessConfig.Location.ServiceRadiusMiles
	}

	// Build travel fee row for pricing table (always show, even if $0)
	travelFeeRowHTML := ""
	if data.TravelFeeInfo != nil {
		travelFeeFormatted := formatCurrency(data.TravelFeeInfo.TravelFee)
		travelFeeMessageText := ""
		if data.TravelFeeInfo.TravelFee == 0 {
			if data.TravelFeeInfo.IsWithinServiceArea {
				travelFeeMessageText = fmt.Sprintf("within %.0f mile radius", serviceRadiusMiles)
			} else {
				travelFeeMessageText = fmt.Sprintf("within %.0f mile radius", serviceRadiusMiles)
			}
		}
		travelFeeRowHTML = fmt.Sprintf(`                  <tr>
                    <td style="font-size: 10.5px; padding: 5px;">- Travel Fee:</td>
                    <td style="font-size: 10.5px; padding: 5px; width: 120px;">%s</td>
                    <td style="font-size: 10.5px; padding: 5px; text-align: left; font-style: italic; color: #666666;">%s</td>
                  </tr>
`, travelFeeFormatted, travelFeeMessageText)
	} else {
		// If travel fee info is not available, show $0 with default message
		travelFeeRowHTML = fmt.Sprintf(`                  <tr>
                    <td style="font-size: 10.5px; padding: 5px;">- Travel Fee:</td>
                    <td style="font-size: 10.5px; padding: 5px; width: 120px;">$0</td>
                    <td style="font-size: 10.5px; padding: 5px; text-align: left; font-style: italic; color: #666666;">within %.0f mile radius</td>
                  </tr>
`, serviceRadiusMiles)
	}

	// Build refund notice HTML (conditional based on days until event)
	refundNoticeHTML := ""
	if IsDepositNonRefundable(data.DaysUntilEvent) {
		// Non-refundable for < 3 days
		refundNoticeHTML = `                <p style="margin: 8px 0 0 0; font-size: 10.5px; color: #991b1b; font-weight: bold;">Deposit is non-refundable.</p>
                <p style="margin: 4px 0 0 0; font-size: 10.5px; color: #666666; font-style: italic;">To fairly pay our helpers and maintain high-quality service, we reserve them when you book — which means they lose other opportunities. Our goal is to have the best helpers for you — to retain them we respect their time and commitment. 98 out of 100 show rate.</p>
`
	} else {
		// Refundable for 3+ days
		refundNoticeHTML = `<p style="margin: 8px 0 0 0; font-size: 10.5px; color: #666666;">100%% refund if cancelled 3+ days before the event.</p>
`
	}

	// Build phone HTML if available
	phoneHTML := ""
	if contact.Phone != "" {
		phoneHTML = fmt.Sprintf(`<a href="tel:%s" style="color: #000000; text-decoration: underline; display: inline-block; margin: 5px 0;">Tap to Call Us: %s</a><br />`,
			strings.ReplaceAll(contact.Phone, "-", ""), contact.Phone)
	}

	// Extract domain from website URL for display
	websiteDomain := contact.WebsiteURL
	if strings.HasPrefix(websiteDomain, "https://") {
		websiteDomain = strings.TrimPrefix(websiteDomain, "https://")
	} else if strings.HasPrefix(websiteDomain, "http://") {
		websiteDomain = strings.TrimPrefix(websiteDomain, "http://")
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>%s - Quote</title>
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
                      <img src="%s" alt="%s - Professional Event Staffing Services" style="max-width: 120px; height: auto;" />
                    </td>
                    <td style="vertical-align: middle; font-size: 10.5px; color: #666666; line-height: 1.6;">
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
                <p style="margin: 0; font-size: 10.5px; font-weight: normal; color: rgb(38, 37, 120);">%s Quote (requested by %s)</p>
                <p style="margin: 3px 0 0 0; font-size: 10.5px; font-weight: bold; color: rgb(38, 37, 120); white-space: nowrap;">Quote ID: %s</p>
                <p style="margin: 3px 0 0 0; font-size: 10.5px; color: rgb(38, 37, 120);">This is a quote to hold your details.</p>
              </td>
            </tr>
            <!-- Urgency Banner (if applicable) - placed above yellow reservation notice -->
%s
            <!-- Yellow Reservation Notice -->
            <tr>
              <td align="center" style="padding: 0 6px 10px 6px;">
                <p style="margin: 0; font-size: 10.5px; font-weight: bold; color: #92400e; background-color: #fff9c4; padding: 10px 6px; text-align: center; border-left: 3px solid #f59e0b;">Your Reservation is NOT confirmed until deposit is received.</p>
              </td>
            </tr>
            
            <!-- Expiration Notice -->
%s

            <!-- Header -->
            <tr>
              <td align="center" style="font-size: 17.5px; font-weight: bold; padding: 8px 0;">Hi %s!</td>
            </tr>
            <tr>
              <td align="center" style="padding-bottom: 16px; font-size: 10.5px; line-height: 1.5;">
%s
              </td>
            </tr>

            <!-- Event Details - AIDA: Interest (build excitement about the event) -->
            <tr>
              <td style="padding-top: 16px;"></td>
            </tr>
            <tr>
              <td style="font-size: 15.5px; font-weight: bold; padding-top: 10px; padding-bottom: 6px; border-top: 2px solid #d0d0d0; text-align: center;">Event Details</td>
            </tr>
            <tr>
              <td align="center" style="background-color: #fafafa; padding: 12px 4px;">
                <table width="540" cellpadding="5" cellspacing="0" border="0" style="background-color: #f9f9f9; width: 540px;">
                  <tr>
                    <td style="font-size: 10.5px; padding: 5px; width: 50%%; vertical-align: top;">
                      <span style="font-weight: bold;">When:</span><br />
                      %s %s<br />
                      <span style="color: %s; font-size: 10.5px; font-weight: 500; font-style: italic;">(starts in %s)</span>
                    </td>
                    <td style="font-size: 10.5px; padding: 5px; width: 50%%; vertical-align: top;">
                      <span style="font-weight: bold;">Where:</span><br />
                      %s<br />
                      <span style="color: %s; font-size: 10.5px; font-weight: 500;">(%s)</span>
                    </td>
                  </tr>
                  <tr>
                    <td style="font-size: 10.5px; padding: 5px; width: 50%%;"><span style="font-weight: bold;">Occasion:</span> %s</td>
                    <td style="font-size: 10.5px; padding: 5px; width: 50%%;"><span style="font-weight: bold;">Guest Count:</span> %d</td>
                  </tr>
                  <tr>
                    <td style="font-size: 10.5px; padding: 5px; width: 50%%;"><span style="font-weight: bold;">%s:</span> %d</td>
                    <td style="font-size: 10.5px; padding: 5px; width: 50%%;"><span style="font-weight: bold;">For How Long:</span> %s Hours</td>
                  </tr>
                </table>
                <p style="font-size: 10.5px; color: #666666; padding: 8px 30px 0 30px; margin: 0; text-align: center; font-style: italic;">
                  We advise our staff start time to be between %s to allow for setup and walk-through.
                </p>
%s
              </td>
            </tr>

            <!-- Services Included -->
            <tr>
              <td style="padding-top: 16px;"></td>
            </tr>
            <tr>
              <td style="font-size: 15.5px; font-weight: bold; padding-top: 10px; padding-bottom: 6px; border-top: 2px solid #d0d0d0; text-align: center;">Services Included</td>
            </tr>
            <tr>
              <td align="center" style="background-color: #fafafa; padding: 12px 4px;">
                <table width="540" cellpadding="0" cellspacing="0" border="0" style="width: 540px; font-size: 10.5px; line-height: 1.5;">
                  <tr>
                    <td style="text-align: left;">
                      <p style="margin: 5px 0; font-weight: bold; font-size: 10.5px;">Setup & Presentation</p>
                      <ul style="margin: 4px 0; padding-left: 25px;">
                        <li style="margin: 3px 0;">Arranging tables, chairs, and decorations</li>
                        <li style="margin: 3px 0;">Buffet setup & live buffet service</li>
                        <li style="margin: 3px 0;">Butler-passed appetizers & cocktails</li>
                      </ul>
                      <p style="margin: 6px 0 4px 0; font-weight: bold; font-size: 10.5px;">Dining & Guest Assistance</p>
                      <ul style="margin: 4px 0; padding-left: 25px;">
                        <li style="margin: 3px 0;">Multi-course plated dinners</li>
                        <li style="margin: 3px 0;">General bussing (plates, silverware, glassware)</li>
                        <li style="margin: 3px 0;">Beverage service (water, wine, champagne, coffee, etc.)</li>
                        <li style="margin: 3px 0;">Special services (cake cutting, dessert plating, etc.)</li>
                      </ul>
                      <p style="margin: 6px 0 4px 0; font-weight: bold; font-size: 10.5px;">Cleanup & End-of-Event Support</p>
                      <ul style="margin: 4px 0; padding-left: 25px;">
                        <li style="margin: 3px 0;">Washing dishes, managing trash, and keeping the event space tidy</li>
                        <li style="margin: 3px 0;">Kitchen cleanup & end-of-event breakdown</li>
                        <li style="margin: 3px 0;">Assisting with food storage & leftovers</li>
                      </ul>
                      <p style="margin: 12px 0 0 0; font-size: 10.5px; color: #666666; line-height: 1.5; text-align: center; font-style: italic;">
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
              <td style="padding-top: 16px;"></td>
            </tr>
            <tr>
              <td style="font-size: 15.5px; font-weight: bold; padding-top: 10px; padding-bottom: 6px; border-top: 2px solid #d0d0d0; text-align: center;">Our Rates & Pricing</td>
            </tr>
            <tr>
              <td align="center" style="background-color: #fafafa; padding: 12px 4px;">
                <table width="540" cellpadding="5" cellspacing="0" border="0" style="background-color: #f9f9f9; width: 540px;">
                  <!-- Breakdown Section -->
                  <tr>
                    <td colspan="3" style="padding: 0;">
                      <table width="100%%" cellpadding="5" cellspacing="0" border="0" style="margin: 5px 0; background-color: #ffffff;">
                        <!-- Estimated Total - First Row -->
                        <tr>
                          <td style="font-weight: bold; font-size: 10.5px; padding: 5px;">Estimated Total:</td>
                          <td style="font-size: 10.5px; padding: 5px; width: 120px;">%s</td>
                          <td style="font-size: 10.5px; padding: 5px; text-align: left;"></td>
                        </tr>
                        <tr>
                          <td style="font-size: 10.5px; padding: 5px;">- Service Rate:</td>
                          <td style="font-size: 10.5px; padding: 5px; width: 120px;">%s</td>
                          <td style="font-size: 10.5px; padding: 5px; text-align: left; font-style: italic; color: #666666;">per helper (first 4 hours)</td>
                        </tr>
%s
                        <tr>
                          <td colspan="3" style="font-size: 10.5px; padding: 8px 5px; text-align: center; font-style: italic; color: #666666;">%s</td>
                        </tr>
%s
                      </table>
                    </td>
                  </tr>
                  <!-- Deposit Amount - Separate from breakdown -->
                  <tr>
                    <td style="font-weight: bold; font-size: 13px; padding: 5px;">Deposit Amount:</td>
                    <td style="font-size: 13px; padding: 5px; width: 120px;">%s</td>
                    <td style="font-size: 13px; padding: 5px; text-align: left;"></td>
                  </tr>
                  <tr>
                    <td colspan="3" style="font-size: 11.5px; padding: 4px 5px; text-align: center; font-style: italic; color: #666666;">%s</td>
                  </tr>
                </table>
                <p style="font-size: 10.5px; color: #666666; padding-top: 8px; padding-bottom: 0; margin: 0; text-align: center; font-style: italic;">
                  Final total may adjust only if event details change. Gratuity is not included but always appreciated!
                </p>
              </td>
            </tr>

            <!-- How to Get Started -->
            <tr>
              <td style="padding-top: 16px;"></td>
            </tr>
            <tr>
              <td style="font-size: 15.5px; font-weight: bold; padding-top: 10px; padding-bottom: 6px; border-top: 2px solid #d0d0d0; text-align: center;">How to Get Started</td>
            </tr>
            <tr>
              <td align="center" style="background-color: #fafafa; padding: 12px 4px;">
                <table width="540" cellpadding="0" cellspacing="0" border="0" style="width: 540px; font-size: 10.5px; line-height: 1.6;">
                  <tr>
                    <td style="text-align: left; padding: 10px;">
                      <ol style="margin: 0; padding-left: 25px; color: #333333;">
                        <li style="margin: 8px 0;"><strong>Pay Your Deposit:</strong> Click the "Pay Deposit" button below to secure your staffing reservation. Your deposit reserves your date and confirms your event.</li>
                        <li style="margin: 8px 0;"><strong>We'll Confirm:</strong> Once your deposit is received, we'll send you a confirmation email with all the details and next steps.</li>
                        <li style="margin: 8px 0;"><strong>Final Payment:</strong> The remaining balance will be due in 7 days (Net 7). You'll receive your final invoice on the next business day after your event (%s) via email with a secure Stripe payment link. You can pay with a check, cash, or securely online via Stripe. This email will also include a separate link to provide gratuity for our helpers if you choose — many customers prefer this option as it keeps the paid reservation record clean and provides flexibility if you don't have cash on hand.</li>
                        <li style="margin: 8px 0;"><strong>Event Day:</strong> Our professional staff will arrive at %s to set up and ensure everything runs smoothly.</li>
                      </ol>
                      <p style="margin: 12px 0 0 0; font-size: 10.5px; color: #666666; line-height: 1.5; text-align: center; font-style: italic;">
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
                  <a href="%s" style="display: inline-block; background-color: rgb(38, 37, 120); color: #ffffff; padding: 6px 12px; text-decoration: none; font-weight: bold; font-size: 10.5px; border-radius: 4px;">Pay Deposit Securely (%s) via Stripe</a>
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
                      <p style="margin: 0 0 8px 0; font-size: 10.5px; line-height: 1.5; text-align: center;">
                        Need a call? Book one if:
                      </p>
                      <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin: 0 auto; border: none;">
                        <tr>
                          <td style="text-align: center; padding: 2px 0; border: none; font-size: 10.5px; line-height: 1.5;">• This is your first time</td>
                        </tr>
                        <tr>
                          <td style="text-align: center; padding: 2px 0; border: none; font-size: 10.5px; line-height: 1.5;">• You have complex food prep</td>
                        </tr>
                        <tr>
                          <td style="text-align: center; padding: 2px 0; border: none; font-size: 10.5px; line-height: 1.5;">• You have china plates</td>
                        </tr>
                        <tr>
                          <td style="text-align: center; padding: 2px 0; border: none; font-size: 10.5px; line-height: 1.5;">• Anything that might affect helper count or hours</td>
                        </tr>
                      </table>
                      <p style="margin: 8px 0; text-align: center;">
                        <a href="%s" style="display: inline-block; background-color: rgb(38, 37, 120); color: #ffffff; padding: 8px 16px; text-decoration: none; font-weight: bold; font-size: 10.5px; border-radius: 4px;">Schedule a Call</a>
                      </p>
                      <p style="margin: 6px 0 0 0; font-size: 10.5px; color: #666666; font-style: italic; text-align: left;">We never oversell — we recommend exactly what you need so your event runs smoothly and you shine. That's why 1,000+ clients trust us and call us back.</p>
                    </td>
                  </tr>
                </table>
              </td>
            </tr>

            <!-- Footer -->
            <tr>
              <td align="center" style="font-size: 10.5px; padding-top: 20px; color: #666666; line-height: 1.5;">
                %s<br />
                %s
                <a href="%s" style="color: rgb(38, 37, 120); text-decoration: underline; display: inline-block; margin: 3px 0;">%s</a><br />
                &copy; %d %s<br />
                <span style="font-size: 10px;">v1.1</span>
              </td>
            </tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>`,
		displayName,           // Title: %s - Quote (line 661)
		contact.LogoURL,       // Logo image src (line 674)
		displayName,           // Logo alt text (line 674)
		availabilityMeterHTML, // Availability meter (always shown, right below logo) (line 685)
		data.Occasion,         // EVENT QUOTE - %s Quote (line 689)
		GetFirstNameWithLastInitial(data.ClientName),   // (requested by First Name L.) (line 689)
		data.ConfirmationNumber,                        // Quote ID: %s (line 690)
		urgencyBannerHTML,                              // Urgency banner (conditional HTML based on time until event) (line 695)
		expirationNoticeHTML,                           // Expiration notice (conditional HTML - only for urgent cases) (line 704)
		GetFirstName(data.ClientName),                  // Hi %s! (line 708)
		greetingMessage,                                // Personalized greeting (returning vs new client) (line 712)
		data.EventDate,                                 // When: %s (Event Date) (line 729)
		data.EventTime,                                 // When: %s (Event Time) (line 729)
		daysUntilEventColor,                            // Color for "(starts in X days)" text (line 730)
		FormatDaysUntilEvent(data.DaysUntilEvent),      // Your event starts %s (after When) (line 730)
		data.EventLocation,                             // Where: %s (line 734)
		travelFeeMessageColor,                          // Color for travel fee message (line 735)
		travelFeeMessage,                               // Travel fee message (e.g., "within our service area - no travel fee") (line 735)
		data.Occasion,                                  // Occasion: %s (line 739)
		data.GuestCount,                                // Guest Count: %d (line 740)
		helpersText,                                    // Helpers/Helper label: %s (line 743)
		data.Helpers,                                   // Helpers count: %d (line 743)
		hoursFormatted,                                 // For How Long: %s Hours (line 744)
		recommendedArrivalTimeRange,                    // Recommended arrival time range: %s (line 748)
		buildWeatherHTML(data.WeatherForecast),         // Weather forecast HTML (only for < 10 days) (line 750)
		totalFormatted,                                 // Estimated Total: %s (first row) (line 812)
		baseRateFormatted,                              // Service Rate: %s (in breakdown) (line 817)
		travelFeeRowHTML,                               // Travel fee row (in breakdown) (line 820)
		additionalHoursText,                            // Additional Hours description: %s (in breakdown) (line 822)
		additionalHoursBeyondEndTimeHTML,               // Additional Hours beyond end time: %s (in breakdown) (line 824)
		depositFormatted,                               // Deposit Amount: %s (separate) (line 831)
		getDepositDeadlineMessage(data.DaysUntilEvent), // (due in X days to secure your staffing reservation) (line 835)
		data.EventDate,                                 // Final payment due date: %s (the day of your event) (line 859)
		recommendedArrivalTimeRange,                    // Staff arrival time: %s (line 860)
		buildPDFDownloadHTML(data.PDFDownloadLink, data.ExpirationDate, data.DaysUntilEvent), // PDF download link section (line 877)
		data.Occasion,              // Ready to Secure Your %s? (occasion) (line 882)
		depositOptionsHTML,         // Deposit options message (returning vs new client) (line 883)
		data.DepositLink,           // Pay Deposit button link (line 885)
		depositFormatted,           // Pay Deposit button text (line 885)
		refundNoticeHTML,           // Refund notice (conditional - non-refundable for < 3 days) (line 887)
		contact.BookAppointmentURL, // Book appointment link (line 916)
		officeAddress,              // Office address (line 928)
		phoneHTML,                  // Phone link HTML (if available) (line 929)
		contact.WebsiteURL,         // Website URL for link (line 930)
		websiteDomain,              // Website domain for display (line 930)
		time.Now().Year(),          // Copyright year (line 931)
		displayName,                // Copyright business name (line 931)
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

// GenerateQuoteEmailHTMLAppleStyle generates HTML using Apple-style design
// Uses the same structure and variables as GenerateQuoteEmailHTML but with Apple styling
func GenerateQuoteEmailHTMLAppleStyle(data QuoteEmailData, businessConfig *domain.BusinessConfig) string {
	// Get contact information with smart defaults
	businessID := "stlpartyhelpers" // Default business ID
	if businessConfig != nil && businessConfig.ID != "" {
		businessID = businessConfig.ID
	}
	contact := GetContactInfo(businessID, businessConfig)
	displayName := GetBusinessDisplayName(businessID, businessConfig)
	officeAddress := GetBusinessOfficeAddress(businessConfig)

	// Reuse all the data preparation logic from GenerateQuoteEmailHTML
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

	hoursFormatted := fmt.Sprintf("%.0f", data.Hours)
	if data.Hours != float64(int(data.Hours)) {
		hoursFormatted = fmt.Sprintf("%.1f", data.Hours)
	}

	helpersText := "Helper"
	if data.Helpers != 1 {
		helpersText = "Helpers"
	}

	costPerAdditionalHour := float64(data.Helpers) * data.HourlyRate
	helperWord := "helper"
	if data.Helpers != 1 {
		helperWord = "helpers"
	}

	additionalHoursText := fmt.Sprintf(
		"Any hour after the initial 4 hours costs %s per helper.<br />With %d %s, each additional hour is %s.",
		hourlyRateFormatted,
		data.Helpers,
		helperWord,
		formatCurrency(costPerAdditionalHour),
	)

	additionalHoursBeyondEndTimeText := ""
	if data.EventTime != "" && data.Hours > 0 {
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
				endTime = eventTime.Add(time.Duration(data.Hours) * time.Hour)
				break
			}
		}

		if !endTime.IsZero() {
			endTimeFormatted := endTime.Format("3:04 PM")
			additionalHoursBeyondEndTimeText = fmt.Sprintf(
				"If our helpers stay longer than anticipated (after %s), we'll add %s for each additional hour. The latest we can extend is 1:00 AM.",
				endTimeFormatted,
				formatCurrency(costPerAdditionalHour),
			)
		}
	}

	// additionalHoursBeyondEndTimeText is used directly in the template, no HTML wrapper needed

	recommendedArrivalTimeRange := ""
	if data.EventTime != "" {
		eventDateTimeStr := fmt.Sprintf("%s %s", data.EventDate, data.EventTime)
		formats := []string{
			"Mon, January 2, 2006 3:04 PM",
			"Mon, Jan 2, 2006 3:04 PM",
			"January 2, 2006 3:04 PM",
			"Jan 2, 2006 3:04 PM",
			"Mon, January 2, 2006 15:04",
			"Mon, Jan 2, 2006 15:04",
			"January 2, 2006 15:04",
			"Jan 2, 2006 15:04",
			"2006-01-02 3:04 PM",
			"2006-01-02 15:04",
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

		if parsed {
			earliestArrival := eventTime.Add(-1 * time.Hour)
			latestArrival := eventTime.Add(-30 * time.Minute)
			recommendedArrivalTimeRange = fmt.Sprintf("%s - %s",
				earliestArrival.Format("3:04 PM"),
				latestArrival.Format("3:04 PM"))
		}
	}

	availabilityMeterHTML := ""
	var availabilityMessage string
	if data.DaysUntilEvent == 0 {
		availabilityMessage = "Deposit Should Be Paid Now — ASAP to Proceed with The Staffing Reservation"
	} else {
		switch data.UrgencyLevel {
		case "critical":
			availabilityMessage = "Limited availability — secure with deposit today"
		case "urgent":
			availabilityMessage = "Filling fast — secure with deposit to guarantee your spot"
		case "high":
			availabilityMessage = "Popular date — secure with deposit soon"
		case "moderate":
			availabilityMessage = "Spots available — secure with deposit to reserve"
		default:
			if data.IsHighDemand {
				availabilityMessage = "Popular date — secure with deposit to guarantee availability"
			} else {
				availabilityMessage = "Secure your date with deposit"
			}
		}
	}

	availabilityMeterHTML = fmt.Sprintf(`            <tr>
              <td align="center" style="background-color: #fff9c4; padding: 8px 6px; margin-bottom: 6px; border-left: 3px solid #f59e0b;">
                <p style="margin: 0; font-size: 10.5px; font-weight: bold; color: #92400e;">%s</p>
              </td>
            </tr>
`, availabilityMessage)

	urgencyBannerHTML := ""
	var bannerMessage string
	if data.DaysUntilEvent == 0 {
		bannerMessage = "Deposit Should Be Paid Now — ASAP to Proceed with The Staffing Reservation"
	} else {
		switch data.UrgencyLevel {
		case "critical":
			bannerMessage = "Only a few spots left — secure your date today to avoid being left out"
		case "urgent":
			bannerMessage = "Dates fill up fast — secure your spot now to guarantee availability"
		case "high":
			bannerMessage = "Popular time period — secure your date soon to guarantee your spot"
		case "moderate":
			bannerMessage = "Spots are filling up — secure your date to guarantee availability"
		default:
			if data.IsHighDemand {
				bannerMessage = "Popular Date — Dates Fill Up Fast. Book Sooner to Secure Your Spot"
			}
		}
	}

	if bannerMessage != "" {
		urgencyBannerHTML = fmt.Sprintf(`            <tr>
              <td align="center" style="background-color: #f0f0f7; padding: 10px 6px; border-left: 3px solid rgb(38, 37, 120); margin-bottom: 6px;">
                <p style="margin: 0; font-size: 10.5px; font-weight: bold; color: rgb(38, 37, 120);">%s</p>
              </td>
            </tr>
`, bannerMessage)
	}

	expirationNoticeHTML := ""
	expirationMessage := getExpirationMessage(data.UrgencyLevel, data.DaysUntilEvent, data.ExpirationDate)
	if expirationMessage != "" {
		expirationNoticeHTML = fmt.Sprintf(`            <tr>
              <td align="center" style="background-color: #f0f0f7; padding: 10px 6px; border-left: 3px solid rgb(38, 37, 120); margin-bottom: 6px;">
                <p style="margin: 0 0 4px 0; font-size: 10.5px; font-weight: bold; color: rgb(38, 37, 120);">%s</p>
                <p style="margin: 0; font-size: 10.5px; color: rgb(38, 37, 120); line-height: 1.4; font-style: italic;">
                  97%% of our clients book us again. Secure your date now.
                </p>
              </td>
            </tr>
`, expirationMessage)
	}

	greetingMessage := ""
	if data.IsReturningClient {
		greetingMessage = `                <strong style="color: rgb(38, 37, 120);">Welcome back! We're thrilled to work with you again.</strong><br />
                My name is Anna, and I am with Customer Success Team.<br />
                We'll handle setup, service, and cleanup so you can enjoy your event.<br />
                Below is your quote with all the details.`
	} else {
		greetingMessage = `                Thank you for your interest in our services!<br />
                My name is Anna, and I am with Customer Success Team.<br />
                We'll handle setup, service, and cleanup so you can enjoy your event.<br />
                Below is your quote with all the details.`
	}

	// daysUntilEventColor and travelFeeMessageColor not used in Apple template (uses simpler styling)
	travelFeeMessage := getTravelFeeMessage(data.TravelFeeInfo)

	travelFeeRowText := ""
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
		travelFeeRowText = fmt.Sprintf("<strong>Travel Fee:</strong> %s%s<br />", travelFeeFormatted, travelFeeMessageText)
	} else {
		travelFeeRowText = "<strong>Travel Fee:</strong> $0 (Within Our Service Radius)<br />"
	}

	// Build PDF download HTML (used in template)
	pdfDownloadHTML := buildPDFDownloadHTML(data.PDFDownloadLink, data.ExpirationDate, data.DaysUntilEvent)
	// Weather HTML not used in Apple template (could be added later)
	_ = buildWeatherHTML(data.WeatherForecast)

	// Generate Apple-style HTML using the actual Apple email structure from the Apple Card email
	// This uses appl_topFeature, appl_body_copy, appl_cb_copy, appl_pill_button, etc. structure
	// Build phone text for contact message
	phoneText := ""
	if contact.Phone != "" {
		phoneText = " or call us directly"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html><head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="color-scheme" content="light">
    <meta name="format-detection" content="date=no">
    <meta name="format-detection" content="address=no">
    <meta name="format-detection" content="email=no">
    <title>Your Event Quote - %s</title>
<style>
:root{color-scheme:light;supported-color-schemes:light}html,body,div,span,applet,object,iframe,h1,h2,h3,h4,h5,h6,p,blockquote,pre,a,abbr,acronym,address,big,cite,code,del,dfn,em,img,ins,kbd,q,samp,small,strike,strong,sub,sup,tt,var,center,dl,dt,dd,ol,ul,li,fieldset,form,label,legend,table,caption,tbody,tfoot,thead,tr,th,td,article,aside,canvas,details,embed,figure,figcaption,footer,header,hgroup,menu,nav,output,ruby,section,summary,time,mark,audio,video{margin:0;padding:0;border:0;font-size:100%%;vertical-align:baseline}article,aside,details,figcaption,figure,footer,header,hgroup,menu,nav,section{display:block}
</style>
<style>body{line-height:1;margin:0;padding:0}ol, ul{list-style:none}blockquote, q{quotes:none}table{border-collapse:collapse;border-spacing:0}body, table, td, a{-webkit-text-size-adjust:100%%;-ms-text-size-adjust:100%%}body{font-family:system-ui, -apple-system, BlinkMacSystemFont, Segoe UI, Helvetica Neue, Helvetica, Arial, sans-serif;text-rendering:optimizeLegibility;background-color:#F5F5F7;color:#1D1D1F}img{-ms-interpolation-mode:bicubic;border:0;outline:none;text-decoration:none;font-size:10px}.appl_body_substitution .appl_main_table a, a, a:visited{color:#0066CC;text-decoration:none}.appl_body_copy p a{color:#0066CC!important;text-decoration:none}.appl_text_cta a{color:#0066CC!important}b, strong{font-weight:700}em, i{font-style:italic}sup{vertical-align:super;font-size:60%%;line-height:1;mso-text-raise:60%%}.appl_body_border{padding:24px}.appl_desktop{display:inline-block}.appl_desktop_block{display:block}.appl_desktop_nobr,.appl_mobile{display:none}.appl_mob_tv{width:200px;height:auto}.appl_legal_links{margin:0 auto;width:100%%}.appl_legal_copy,.appl_footer_links{font-size:10px;line-height:1.2}.appl-legal,.appl_legal_left_aligned{line-height:1.2;text-align:center;color:#6E6E73}.appl-legal a,.appl_legal_left_aligned a{color:#424245!important}.appl-legal p,.appl_legal_left_aligned p{margin:0 auto 20px}.appl_bento_blurb{display:table-cell}a[href^="x-apple-data-detectors:"]{color:inherit!important;text-decoration:none!important;font-size:inherit!important;font-family:inherit!important;font-weight:inherit!important;line-height:inherit!important}@media screen and (max-width:480px){.appl_100,.appl_mob_width{width:100%%!important;height:auto}.appl_body_border,.appl_paddingAdjust,.appl_footer_links_adjust{padding:0!important}.appl_end_wrapper,.appl_riverInner,.appl_transactional td,.appl_tf_gloss,.appl_copy_icon,.appl_topFeatureHeadline,.appl_cb_copy,.appl_cb_disclaimer,.appl_body_wrapper,.appl_disclamer,.appl_body_copy,.appl_1up_image,.appl_bento_1up_sides,.appl_bento_1up_mobile,.appl_bento_1up_sidePaddingAdjust,.appl_bento_inner{padding-left:27px!important;padding-right:27px!important}.appl_cso_none_p,.appl_legal_left_aligned,.appl-legal,.appl_3up_td,.appl_riverBlurb{padding:0 27px!important}.cbMobMargin,.appl_footer_legal_blurb,.appl_2up_bentoBlurb{margin:0 27px!important}.appl_mob_chiclet,.appl_mob_tv,.appl_mob_books,.appl_mob_default,.appl_mob_width,.appl_mob_books_irr{display:block;margin:0 auto}.appl_desktop,.appl_desktop_block,.appl_mobile_nobr,.appl_bento_blurb{display:none!important}.appl_mobile,.appl_pill_button_rtl{display:inline-block!important}.appl_mobile_block,.appl_desktop_nobr,.appl_bento_1up_img,.appl_twoInThree{display:block!important}.appl_topFeatureHeadline h1,.appl_topFeatureHeadline{font-size:32px!important;line-height:36px!important}h1.appl_oversize{font-size:36px!important;line-height:40px!important}.appl_topFeatureHeadline h2{font-size:28px!important;line-height:32px!important}.appl_topFeatureHeadline h3{font-size:22px!important;line-height:26px!important}.appl_matchTop, h1.appl_oversize.appl_matchTop{font-size:32px!important;line-height:36px!important}.appl_matchSecond, h1.appl_oversize.appl_matchSecond{font-size:28px!important;line-height:32px!important}.appl_matchThird, h1.appl_oversize.appl_matchThird{font-size:22px!important;line-height:26px!important}.appl_matchFourth, h1.appl_oversize.appl_matchFourth{font-size:18px!important;line-height:22px!important}.appl_m_body_copy{padding:36px 66px 0!important}.appl_2up_mo{float:none!important;margin:0 auto!important;width:100%%}.appl_2up_mo_inner,.appl_riverCTA_mobilePadding,.appl_blub td,.appl_cb_bodycopy{padding:0 27px 48px!important}.appl_2up_image{width:100%%!important}.apple_3up_locker{float:none!important;margin:0 auto!important;padding:0!important}.appl_3up_figure{padding-right:27px!important;height:auto!important}.appl_3up2up{padding-left:27px!important}.appl_3up_caption{padding-right:27px!important}.appl_notice{padding:24px!important}.appl_eyebrow{height:80px}.appl_eyebrow_copy,.appl_eyebrow_copy p{font-size:14px!important}.appl_m_right_pad{padding-right:24px!important}.appl_mod_chicklet_head{font-size:22px!important}.appl_eyebrow_chicklet{width:240px!important}.appl_eyebrow_chicklet_icn{width:32px!important;height:32px!important}.appl_eyebrow_chicklet_lob{padding:0 0 0 24px!important;vertical-align:middle!important}.appl_eyebrow_chicklet_lob img{width:auto!important;height:36px!important}.appl_eyebrow_chicklet_figure{padding:24px 24px 24px 6px!important;width:32px;line-height:0}.appl_eyebrow_chicklet_caption p{font-size:14px!important;line-height:1.1}.appl_eyebrow_chicklet_caption{padding:24px 24px 0 0}.appl_3up_figure img{width:100%%;height:auto}.appl_cb_eyebrow_inner{padding:22px 0 27px 24px!important}.appl_cb_eyebrow_copy{clear:both;text-align:left!important;margin:0!important;padding:0 6px!important;font-size:14px!important;line-height:1.2!important;float:none!important}.appl_cbs_bc,.appl_cbs_bottomPad,.appl_watchHeroMoCtaAdjust{padding-bottom:32px!important}.appl_cbs_topPad{padding-top:32px!important}.appl_feature_copy_cta p{font-size:18px;line-height:24px}.appl_end_head_below{padding:0 0 32px!important}.appl_hero_header{font-size:32px!important;line-height:36px!important}.appl_header1,.appl_header_above h1{font-size:28px!important;line-height:32px!important}.appl_legal_copy{font-size:10px!important;line-height:14px!important}.appl_ljust_block .appl_lr_padding{padding:0 28px 48px!important}.appl_footer_links{line-height:1.4}.appl_cso_none_m,.appl_list{margin:0!important}.appl_list ul,.appl_list ol{margin:0!important; padding:0 27px 0 54px!important;}.appl_bullet{ padding:0 10px 0 27px!important;}.appl_listItem{ padding-right:27px!important;}.appl-footer{line-height:1.7!important}.appl_mob_chiclet{width:96px!important;height:96px!important}.appl_mob_tv{width:267px!important;height:auto!important}.appl_mob_books{height:auto!important;width:160px!important}.appl_mob_books_irr{width:240px!important;height:auto!important}.appl_mob_default{width:160px!important;height:auto!important}.appl_cso_bottom_p,.appl_3up_header,.appl_bento_1up_headline_noTop{padding:0 27px 32px!important}.appl_cso_bottom_m{margin:0 0 32px!important}.appl_tv_wordmark{width:auto!important;height:36px!important}.appl_cta2_topPad{padding:32px 0 0!important}.appl_sxs_1{display:block!important;padding:0 0 32px!important}.appl_sxs_1 .appl_pill_button{margin:0 auto}.appl_sxs_1 .appl_pill_button_rtl{margin:0 auto;border-collapse:separate;display:table!important;display:revert!important}.appl_sxs_1 .appl_secondary{display:table!important;border-collapse:separate}.appl_sxs_2{padding-right:0!important;display:block!important}.appl_sxs_mobile{height:40px!important}.appl_sxs_mobile_auto,.appl_rmMobCapHeight{height:auto!important}.appl_footer a,.appl_f_router p,.appl_footer_pipe{font-size:10px!important}.appl_riverInner_fullBleed{padding-left:0!important;padding-right:0!important}.appl_riverImage_top10{border-top-left-radius:10px;border-top-right-radius:10px}.appl_riverImage_bottom10{border-bottom-left-radius:10px;border-bottom-right-radius:10px}.appl_riverImage_top30{border-top-left-radius:30px;border-top-right-radius:30px}.appl_riverImage_bottom30{border-bottom-left-radius:30px;border-bottom-right-radius:30px}.appl_river_blurbAlignment,.appl_riverBlurb_mo_center{text-align:center!important}.appl_riverBlurb_noTopPadding,.appl_cbs_rmMidPad,.appl_cta2_noTopPad{padding-top:0!important}.appl_riverBlurb_icon{margin:0 auto!important}.appl_riverBlurbContainer{padding:32px 0px!important}.appl_riverBlurb_bot24,.appl_padding_bot24,.appl_mobilePaddingAdjust{padding-bottom:24px!important}.appl_riverBlurb_bot0,.appl_cbs_rmBottomPad,.appl_moCtaPaddingAdjust,.appl_bento_padding{padding-bottom:0!important}.appl_bento_1up_headline{padding:27px 27px 32px!important}.appl_blurb_padding{padding-bottom:48px!important}.appl_bento_1up_topBotpaddingAdjust{padding-top:27px!important;padding-bottom:27px!important}.appl_bento_1up_top{padding-top:27px!important}.appl_bento_1up_bot{padding-bottom:27px!important}.appl_bento_1up_img_mobile{padding:0 20px 32px!important}.appl_2up_bentoBlurbCopy{margin:16px 27px 0!important}.appl_comp_grid{padding:0 27px!important;font-size:14px}.appl_comp_col1{width:60%%!important}.appl_comp_col2{width:20%%!important}}@media screen and (max-width:320px){.appl_eyebrow_chicklet_caption{display:none}.appl_eyebrow_chicklet{width:60px!important}}</style>
<!--[if gte mso 9]>
<xml> <o:OfficeDocumentSettings> <o:AllowPNG /> <o:PixelsPerInch>96</o:PixelsPerInch> </o:OfficeDocumentSettings> </xml>
<style type="text/css">
a {
    text-decoration: none!important;
}

h1, h2, h3, h4, h6, p {
    mso-line-height-rule: exactly!important;
}

.appl_pill_button_cell {
    padding: 0 20px!important;
}

.mso_block {display: block!important}
</style>
<![endif]-->

</head>
<body>

<table role="presentation" class="appl_body_substitution" cellspacing="0" cellpadding="0" border="0" bgcolor="#F5F5F7" style="background-color: #F5F5F7; color: #1D1D1F; font-family: system-ui, -apple-system, BlinkMacSystemFont, Segoe UI, Helvetica Neue, Helvetica, Arial, sans-serif; margin: 0 auto; text-rendering: optimizeLegibility; width: 100%%;">
    <tr>
        <td class="appl_body_border" style="padding: 24px;">
            <table role="presentation" dir="ltr" cellspacing="0" cellpadding="0" border="0" align="center" width="736px" class="appl_main_table appl_100" style="background-color: #FFFFFF; color: #1D1D1F; font-family: system-ui, -apple-system, BlinkMacSystemFont, Segoe UI, Helvetica Neue, Helvetica, Arial, sans-serif; margin: 0 auto; width: 736px;">
                <tr>
                    <td class="appl_modules" align="center">

                <table role="presentation" dir="ltr" cellspacing="0" cellpadding="0" border="0" width="736" height="88" style="background-color: #FFF" class="appl_eyebrow appl_100">
            <tr>
                <td style="vertical-align: middle; padding-left: 32px;" class="appl_eyebrow_chicklet_lob">
                    <img src="%s" alt="%s" style="max-width: 120px; height: auto;" />
                </td>
                <td style="vertical-align: middle; padding-left: 16px; font-size: 10.5px; color: #6E6E73; line-height: 1.6;">
                    <strong>Status:</strong> Quote<br />
                    <strong>Is My Reservation Confirmed?:</strong> Awaiting deposit
                </td>
            </tr>
        </table>

%s

               <table class="appl_topFeature" dir="ltr" border="0" cellspacing="0" cellpadding="0" role="presentation" width="100%%" style="text-align: center; background-color: transparent;">
        <tr>
        <td class="appl_cb_copy" style="padding: 0 32px 0 32px">
            <p style="display: block; margin: 0; padding: 0" align="center">
                <h1 style="font-size: 28px; line-height: 32px; color: #1D1D1F; font-weight: 600; margin: 0; padding: 0 0 8px;">%s Quote (requested by %s)</h1>
            </p>
            <p style="margin: 8px 0 0 0; font-size: 15.5px; font-weight: 600; color: #1D1D1F;">Quote ID: %s</p>
            <p style="margin: 4px 0 0 0; font-size: 10.5px; color: #6E6E73;">This is a quote to hold your details.</p>
        </td>
    </tr>
%s
        <tr>
            <td class="appl_body_copy appl_cbs_rmBottomPad" style="text-align: center; padding: 0 88px 16px 88px;">
                <p style="color: #92400e; font-size: 15.5px; font-weight: 600; line-height: 20px; margin: 0px auto; padding: 12px; background-color: #fff9c4; border-left: 3px solid #f59e0b;">
                    Your Reservation is NOT confirmed until deposit is received.
                </p>
            </td>
        </tr>
%s
        <tr>
            <td class="appl_body_copy appl_cbs_rmBottomPad" style="text-align: center; padding: 0 88px 0 88px;">
                <p style="color: #1D1D1F; font-size: 17.5px; font-weight: 400; line-height: 24px; margin: 0px auto; padding: 0">
                    Hi %s!<br /><br />
%s
                </p>
            </td>
        </tr>
    </table>

               <table class="appl_topFeature" dir="ltr" border="0" cellspacing="0" cellpadding="0" role="presentation" width="100%%" style="text-align: center; background-color: transparent;">
                        <tr>
                <td style="padding: 48px 0 0; font-size: 0">&nbsp;</td>
            </tr>
                <tr>
        <td class="appl_cb_copy" style="padding: 0 32px 0 32px">
        </td>
    </tr>
                        <tr>
            <td style="padding: 0;">
                                 <table role="presentation" cellspacing="0" cellpadding="0" border="0" class="appl_pill_button  appl_desktop" align="center" style="background: #0066CC; border-radius: 99px; color: #FFFFFF; margin: 0 auto">
                            <tr>
                                <td class="appl_pill_button_cell" height="40" style="color: #FFFFFF; font-size: 17.5px; font-weight: 700; height: 40px; vertical-align: middle; text-align: center; text-decoration: none;">
                                    <a style="color: #FFFFFF; display: inline-block; font-size: 17.5px; font-weight: 400; line-height: 40px; padding: 0 20px; text-align: center; text-decoration: none;" href="%s">View Your Quote</a>
                                </td>
                            </tr>
                        </table>
                        <table role="presentation" cellspacing="0" cellpadding="0" border="0" class="appl_mobile  appl_pill_button" align="center" style="display: none; background: #0066CC;border-radius: 99px; color: #FFFFFF; margin: 0 auto;">
                            <tr>
                                <td class="appl_pill_button_cell" style="color: #FFFFFF; font-size: 17.5px; font-weight: 700; height: 40px; vertical-align:middle;text-align: center; text-decoration: none;">
                                    <a style="color: #FFFFFF; display: inline-block; font-size: 17.5px; font-weight: 400; line-height: 40px; padding: 0 20px; text-align: center; text-decoration: none;" href="%s">View Your Quote</a>
                                </td>
                            </tr>
                        </table>
            </td>
        </tr>
                        <tr>
                <td style="padding: 0 0 48px; font-size: 0">&nbsp;</td>
            </tr>
            </table>

               <table class="appl_topFeature" dir="ltr" border="0" cellspacing="0" cellpadding="0" role="presentation" width="100%%" style="text-align: center; background-color: transparent;">
                    <tr>
        <td class="appl_cb_copy" style="padding: 0 32px 0 32px">
        </td>
    </tr>
                        <tr>
            <td class="appl_body_copy " style="text-align: center; padding: 0 88px 32px 88px;">
                <p style="color: #1D1D1F; font-size: 17.5px; font-weight: 400; line-height: 24px; margin: 0px auto; padding: 0">
                    %s isn't a typical event service. It's everything an event service should be and more. You get transparent pricing, professional setup, and reliable service. That's real value you get delivered directly to your event that you can trust every time.
<br><br>
View your quote in detail and see your event pricing and service details — all clearly presented.
                </p>
            </td>
        </tr>
    </table>

               <table class="appl_topFeature" dir="ltr" border="0" cellspacing="0" cellpadding="0" role="presentation" width="100%%" style="text-align: center; background-color: transparent;">
                    <tr>
        <td class="appl_cb_copy" style="padding: 0 32px 0 32px">
        </td>
    </tr>
            <tr>
            <td class="appl_topFeatureHeadline appl_cbs_bottomPad" style="padding: 0 32px 32px 32px; text-align: center; vertical-align: top; line-height: normal;">
                                    <h3 class="" style="font-size: 22px; letter-spacing: 0; line-height: 26px; color: #1D1D1F; font-weight: 600; margin: 0; padding: 0">No hidden fees.<br> Not even one.</h3>
                            </td>
        </tr>
                                                <tr>
            <td class="appl_body_copy appl_cbs_rmBottomPad" style="text-align: center; padding: 0 88px 0 88px;">
                <p style="color: #1D1D1F; font-size: 17.5px; font-weight: 400; line-height: 24px; margin: 0px auto; padding: 0">
                    Your event service should work for you, not against you. With %s, there are no hidden fees, surprise charges, or unexpected costs. You don't have to worry about fees, because everything is transparent.
                </p>
            </td>
        </tr>
                        <tr>
                <td style="padding: 0 0 48px; font-size: 0">&nbsp;</td>
            </tr>
            </table>

               <table class="appl_topFeature" dir="ltr" border="0" cellspacing="0" cellpadding="0" role="presentation" width="100%%" style="text-align: center; background-color: transparent;">
                    <tr>
        <td class="appl_cb_copy" style="padding: 0 32px 0 32px">
        </td>
    </tr>
            <tr>
            <td class="appl_topFeatureHeadline appl_cbs_bottomPad" style="padding: 0 32px 32px 32px; text-align: center; vertical-align: top; line-height: normal;">
                                    <h3 class="" style="font-size: 22px; letter-spacing: 0; line-height: 26px; color: #1D1D1F; font-weight: 600; margin: 0; padding: 0">Professional setup.<br>Every time.</h3>
                            </td>
        </tr>
                                                <tr>
            <td class="appl_body_copy appl_cbs_rmBottomPad" style="text-align: center; padding: 0 88px 0 88px;">
                <p style="color: #1D1D1F; font-size: 17.5px; font-weight: 400; line-height: 24px; margin: 0px auto; padding: 0">
                    You get professional setup and reliable service on all your event needs. It's real value that never expires or loses its quality. You can also count on our team to be there when you need us, whenever you need us.
                </p>
            </td>
        </tr>
                        <tr>
                <td style="padding: 0 0 48px; font-size: 0">&nbsp;</td>
            </tr>
            </table>

               <table class="appl_topFeature" dir="ltr" border="0" cellspacing="0" cellpadding="0" role="presentation" width="100%%" style="text-align: center; background-color: transparent;">
                    <tr>
        <td class="appl_cb_copy" style="padding: 0 32px 0 32px">
        </td>
    </tr>
            <tr>
            <td class="appl_topFeatureHeadline appl_cbs_bottomPad" style="padding: 0 32px 32px 32px; text-align: center; vertical-align: top; line-height: normal;">
                                    <h3 class="" style="font-size: 22px; letter-spacing: 0; line-height: 26px; color: #1D1D1F; font-weight: 600; margin: 0; padding: 0">Event Details</h3>
                            </td>
        </tr>
                                                <tr>
            <td class="appl_body_copy appl_cbs_rmBottomPad" style="text-align: center; padding: 0 88px 0 88px;">
                <p style="color: #1D1D1F; font-size: 17.5px; font-weight: 400; line-height: 24px; margin: 0px auto; padding: 0">
                    <strong>When:</strong> %s %s (starts in %s)<br />
                    <strong>Where:</strong> %s (%s)<br />
                    <strong>Occasion:</strong> %s<br />
                    <strong>Guest Count:</strong> %d<br />
                    <strong>%s:</strong> %d<br />
                    <strong>For How Long:</strong> %s Hours<br />
                    <br />
                    We advise our staff start time to be between %s to allow for setup and walk-through.
                </p>
            </td>
        </tr>
                        <tr>
                <td style="padding: 0 0 48px; font-size: 0">&nbsp;</td>
            </tr>
    </table>

               <table class="appl_topFeature" dir="ltr" border="0" cellspacing="0" cellpadding="0" role="presentation" width="100%%" style="text-align: center; background-color: transparent;">
                    <tr>
        <td class="appl_cb_copy" style="padding: 0 32px 0 32px">
        </td>
    </tr>
            <tr>
            <td class="appl_topFeatureHeadline appl_cbs_bottomPad" style="padding: 0 32px 32px 32px; text-align: center; vertical-align: top; line-height: normal;">
                                    <h3 class="" style="font-size: 22px; letter-spacing: 0; line-height: 26px; color: #1D1D1F; font-weight: 600; margin: 0; padding: 0">Services Included</h3>
                            </td>
        </tr>
                                                <tr>
            <td class="appl_body_copy appl_cbs_rmBottomPad" style="text-align: center; padding: 0 88px 0 88px;">
                <p style="color: #1D1D1F; font-size: 17.5px; font-weight: 400; line-height: 24px; margin: 0px auto; padding: 0">
                    <strong>Setup & Presentation:</strong> Arranging tables, chairs, and decorations. Buffet setup & live buffet service. Butler-passed appetizers & cocktails.<br /><br />
                    <strong>Dining & Guest Assistance:</strong> Multi-course plated dinners. General bussing (plates, silverware, glassware). Beverage service (water, wine, champagne, coffee, etc.). Special services (cake cutting, dessert plating, etc.).<br /><br />
                    <strong>Cleanup & End-of-Event Support:</strong> Washing dishes, managing trash, and keeping the event space tidy. Kitchen cleanup & end-of-event breakdown. Assisting with food storage & leftovers.<br /><br />
                    Need something specific or a small adjustment? Just let us know — we'll do our best to accommodate.
                </p>
            </td>
        </tr>
                        <tr>
                <td style="padding: 0 0 48px; font-size: 0">&nbsp;</td>
            </tr>
    </table>

               <table class="appl_topFeature" dir="ltr" border="0" cellspacing="0" cellpadding="0" role="presentation" width="100%%" style="text-align: center; background-color: transparent;">
                    <tr>
        <td class="appl_cb_copy" style="padding: 0 32px 0 32px">
        </td>
    </tr>
            <tr>
            <td class="appl_topFeatureHeadline appl_cbs_bottomPad" style="padding: 0 32px 32px 32px; text-align: center; vertical-align: top; line-height: normal;">
                                    <h3 class="" style="font-size: 22px; letter-spacing: 0; line-height: 26px; color: #1D1D1F; font-weight: 600; margin: 0; padding: 0">Our Rates & Pricing</h3>
                            </td>
        </tr>
                                                <tr>
            <td class="appl_body_copy appl_cbs_rmBottomPad" style="text-align: center; padding: 0 88px 0 88px;">
                <p style="color: #1D1D1F; font-size: 17.5px; font-weight: 400; line-height: 24px; margin: 0px auto; padding: 0">
                    <strong>Estimated Total:</strong> %s<br />
                    <strong>Service Rate:</strong> %s / helper (first 4 hours)<br />
%s
                    <br />
                    %s<br />
%s
                    <br />
                    <strong>Deposit Amount:</strong> %s<br />
                    %s<br />
                    <br />
                    Final total may adjust only if event details change. Gratuity is not included but always appreciated!
                </p>
            </td>
        </tr>
                        <tr>
                <td style="padding: 0 0 48px; font-size: 0">&nbsp;</td>
            </tr>
    </table>

               <table class="appl_topFeature" dir="ltr" border="0" cellspacing="0" cellpadding="0" role="presentation" width="100%%" style="text-align: center; background-color: transparent;">
                    <tr>
        <td class="appl_cb_copy" style="padding: 0 32px 0 32px">
        </td>
    </tr>
            <tr>
            <td class="appl_topFeatureHeadline appl_cbs_bottomPad" style="padding: 0 32px 32px 32px; text-align: center; vertical-align: top; line-height: normal;">
                                    <h3 class="" style="font-size: 22px; letter-spacing: 0; line-height: 26px; color: #1D1D1F; font-weight: 600; margin: 0; padding: 0">How to Get Started</h3>
                            </td>
        </tr>
                                                <tr>
            <td class="appl_body_copy appl_cbs_rmBottomPad" style="text-align: center; padding: 0 88px 0 88px;">
                <p style="color: #1D1D1F; font-size: 17.5px; font-weight: 400; line-height: 24px; margin: 0px auto; padding: 0">
                    <strong>1. Pay Your Deposit:</strong> Click the "Pay Deposit" button below to secure your staffing reservation. Your deposit reserves your date and confirms your event.<br /><br />
                    <strong>2. We'll Confirm:</strong> Once your deposit is received, we'll send you a confirmation email with all the details and next steps.<br /><br />
                    <strong>3. Final Payment:</strong> The remaining balance will be due in 7 days (Net 7). You'll receive your final invoice on the next business day after your event (%s) via email with a secure Stripe payment link. You can pay with a check, cash, or securely online via Stripe. This email will also include a separate link to provide gratuity for our helpers if you choose — many customers prefer this option as it keeps the paid reservation record clean and provides flexibility if you don't have cash on hand.<br /><br />
                    <strong>4. Event Day:</strong> Our professional staff will arrive at %s to set up and ensure everything runs smoothly.<br /><br />
                    Questions? We're here to help! Schedule a call or reach out anytime.
                </p>
            </td>
        </tr>
                        <tr>
                <td style="padding: 0 0 48px; font-size: 0">&nbsp;</td>
            </tr>
    </table>

%s

               <table cellspacing='0' cellpadding='0' border='0' align="center" valign="middle" class="appl_end appl_100" role="presentation" style="text-align: center;" width="736">
    <tr>
        <td class="appl_end_wrapper" style="padding: 48px 80px 0; text-align:center" align="center" valign="middle" >
                <h2 class="appl_cso_bottom_m" align="center" style="padding: 0; font-size: 17.5px; line-height: 1.2; font-weight: 400; text-align: center; margin: 0 50px 32px">Ready to Secure Your %s?<br class="appl_desktop">Pay your deposit to confirm your event.</h2>
                    </td>
    </tr>
        <tr>
            <td style="padding: 0" class="" align="center" valign="middle">
                                 <table role="presentation" cellspacing="0" cellpadding="0" border="0" class="appl_pill_button appl_desktop" align="center" style="background: #0066CC;border-radius: 99px; color: #FFFFFF; margin: 0 auto;">
                            <tr>
                                <td class="appl_pill_button_cell" height="40" style="color: #FFFFFF; font-size: 17.5px; font-weight: 700; height: 40px; vertical-align:middle;text-align: center; text-decoration: none;">
                                    <a style="color: #FFFFFF; display: block; font-size: 17.5px; font-weight: 400; line-height: 40px; padding: 0 20px; text-align: center; text-decoration: none;" href="%s">Pay Deposit Securely (%s) via Stripe</a>
                                </td>
                            </tr>
                        </table>
                            <table role="presentation" cellspacing="0" cellpadding="0" border="0" class="appl_pill_button appl_mobile" align="center" style="display:none; background: #0066CC; border-radius: 99px; color: #FFFFFF; margin: 0 auto;">
                                <tr>
                                    <td class="appl_pill_button_cell" height="40" style="color: #FFFFFF; font-size: 17.5px; font-weight: 700; height: 40px; vertical-align: middle; text-align: center; text-decoration: none;">
                                        <a style="color: #FFFFFF; display: block; font-size: 17.5px; font-weight: 400; line-height: 40px; padding: 0 20px; text-align: center; text-decoration: none;" href="%s">Pay Deposit Securely (%s) via Stripe</a>
                                    </td>
                                </tr>
                            </table>
            </td>
        </tr>
                                <tr>
            <td style="padding: 0 0 48px; font-size: 0">&nbsp;</td>
        </tr>
    </table>

               <table class="appl_topFeature" dir="ltr" border="0" cellspacing="0" cellpadding="0" role="presentation" width="100%%" style="text-align: center; background-color: transparent;">
                    <tr>
        <td class="appl_cb_copy" style="padding: 0 32px 0 32px">
        </td>
    </tr>
            <tr>
            <td class="appl_topFeatureHeadline appl_cbs_bottomPad" style="padding: 0 32px 32px 32px; text-align: center; vertical-align: top; line-height: normal;">
                                    <h3 class="" style="font-size: 22px; letter-spacing: 0; line-height: 26px; color: #1D1D1F; font-weight: 600; margin: 0; padding: 0">Questions?</h3>
                            </td>
        </tr>
                                                <tr>
            <td class="appl_body_copy appl_cbs_rmBottomPad" style="text-align: center; padding: 0 88px 0 88px;">
                <p style="color: #1D1D1F; font-size: 17.5px; font-weight: 400; line-height: 24px; margin: 0px auto; padding: 0">
                    Need a call? Book one if: This is your first time, you have complex food prep, you have china plates, or anything that might affect helper count or hours.<br /><br />
                    <a href="%s" style="color: #0066CC; text-decoration: none;">Schedule a Call</a><br /><br />
                    We never oversell — we recommend exactly what you need so your event runs smoothly and you shine. That's why 1,000+ clients trust us and call us back.
                </p>
            </td>
        </tr>
                        <tr>
                <td style="padding: 0 0 48px; font-size: 0">&nbsp;</td>
            </tr>
    </table>

<table cellspacing='0' cellpadding='0' border='0' width="736" class="appl_100" role="presentation" style="background-color: #F5F5F7">
    <tr>
        <td>
            <table dir="ltr" cellspacing='0' cellpadding='0' border='0' role="presentation" class="appl_footer appl_100" align="center" style="margin: 32px auto; width: 736;">
                <tr>
                    <td>
                        <table dir="ltr" class="appl_footer appl_legal_links" cellspacing='0' cellpadding='0' border='0' align="inherit" style="background-color: #F5F5F7;">
                            <tr>
                              <td class="appl_legal_left_aligned" style="padding: 0 88px;">
                               <p class="appl_legal_copy" style="color:#6E6E73; text-align: left; font-size: 10.5px; line-height: 1.3;">Hello %s,<br><br>
Your event quote is ready for %s. The quote includes all services and pricing details. Please review your quote and let us know if you have any questions.<br><br>
If you have any questions or need to make changes, please contact us at %s%s or call us directly.<br><br>
Thank you for choosing %s for your event needs.</p>
                                                                  <p class="appl_legal_copy" style="color:#6E6E73; text-align: left; font-size: 10.5px; line-height: 1.3;">All rights reserved. Copyright &copy; %d %s. %s.</p>
                              </td>
                            </tr>
                                                        <tr>
                                <td class="appl_footer" align="center" style="padding: 0 16px;">
                                    <table dir="ltr" cellspacing='0' cellpadding='0' border='0' style="background-color: #F5F5F7; color: #424245; text-align: center; margin: 0 auto;">
                                        <tr>
                                                                                 <td class="appl_footer_links" style="display: inline-block;">
                                                    <a href="%s" style="color: #424245; text-decoration: none; white-space: nowrap; font-size: 10.5px; line-height: 1.3">Privacy Policy</a>
                                                                                                 <span class="appl_footer_pipe" style="font-size: 10.5px">&nbsp;&nbsp;|&nbsp;&nbsp;</span>
                                                                                         </td>
                                                                                                                             <td class="appl_footer_links" style="display: inline-block;">
                                                    <a href="%s" style="color: #424245; text-decoration: none; white-space: nowrap; font-size: 10.5px; line-height: 1.3">Terms & Conditions</a>
                                                                                                 <span class="appl_footer_pipe" style="font-size: 10.5px">&nbsp;&nbsp;|&nbsp;&nbsp;</span>
                                                                                         </td>
                                                                                                                             <td class="appl_footer_links" style="display: inline-block;">
                                                    <a href="mailto:%s" style="color: #424245; text-decoration: none; white-space: nowrap; font-size: 10.5px; line-height: 1.3">Support</a>
                                                                                                 <span class="appl_footer_pipe" style="font-size: 10.5px">&nbsp;&nbsp;|&nbsp;&nbsp;</span>
                                                                                         </td>
                                                                                                                             <td class="appl_footer_links" style="display: inline-block;">
                                                    <a href="%s" style="color: #424245; text-decoration: none; white-space: nowrap; font-size: 10.5px; line-height: 1.3;">Contact</a>
                                                </td>
                                                                         </tr>
                                    </table>
                                </td>
                            </tr>
                        </table>
                    </td>
                </tr>
            </table>
        </td>
    </tr>
</table>

                    </td>
                </tr>
            </table>
        </td>
    </tr>
</table>
</body>
</html>`,
		availabilityMeterHTML,
		data.Occasion,
		GetFirstNameWithLastInitial(data.ClientName),
		data.ConfirmationNumber,
		urgencyBannerHTML,
		expirationNoticeHTML,
		GetFirstName(data.ClientName),
		greetingMessage,
		data.PDFDownloadLink,
		data.PDFDownloadLink,
		data.EventDate,
		data.EventTime,
		FormatDaysUntilEvent(data.DaysUntilEvent),
		data.EventLocation,
		travelFeeMessage,
		data.Occasion,
		data.GuestCount,
		helpersText,
		data.Helpers,
		hoursFormatted,
		recommendedArrivalTimeRange,
		totalFormatted,
		baseRateFormatted,
		travelFeeRowText,
		additionalHoursText,
		additionalHoursBeyondEndTimeText,
		depositFormatted,
		getDepositDeadlineMessage(data.DaysUntilEvent),
		data.EventDate,
		recommendedArrivalTimeRange,
		pdfDownloadHTML,
		data.Occasion,
		data.DepositLink,
		depositFormatted,
		data.DepositLink,
		depositFormatted,
		contact.BookAppointmentURL,
		displayName,          // Title: Your Event Quote - %s (line 1309)
		contact.LogoURL,      // Logo image src (line 1346)
		displayName,          // Logo alt text (line 1346)
		displayName,          // STL Party Helpers isn't... (line 1425)
		displayName,          // With STL Party Helpers... (line 1446)
		contact.SupportEmail, // Support email in contact message
		phoneText,            // Phone text (if available)
		displayName,          // Thank you for choosing %s
		time.Now().Year(),    // Copyright year
		displayName,          // Copyright business name
		officeAddress,        // Office address
		contact.WebsiteURL,   // Privacy Policy link
		contact.WebsiteURL,   // Terms & Conditions link
		contact.SupportEmail, // Support email (mailto)
		contact.WebsiteURL,   // Contact link
		GetFirstName(data.ClientName),
		data.EventDate,
		time.Now().Year(),
	)

	return html
}
