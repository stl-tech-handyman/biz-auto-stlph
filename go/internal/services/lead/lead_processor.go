package lead

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bizops360/go-api/internal/infra/calendar"
	"github.com/bizops360/go-api/internal/infra/email"
	"github.com/bizops360/go-api/internal/infra/geo"
	"github.com/bizops360/go-api/internal/infra/stripe"
	"github.com/bizops360/go-api/internal/ports"
	"github.com/bizops360/go-api/internal/services/pricing"
	"github.com/bizops360/go-api/internal/util"
)

// Processor orchestrates the lead processing workflow
// This is the main entry point that coordinates all steps
type Processor struct {
	calendarService  *calendar.CalendarService
	emailClient      *email.EmailServiceClient
	gmailSender      *email.GmailSender
	geocodingService *geo.GeocodingService
	logger           *slog.Logger
	calendarID       string
}

// NewProcessor creates a new lead processor
func NewProcessor(
	calendarService *calendar.CalendarService,
	emailClient *email.EmailServiceClient,
	gmailSender *email.GmailSender,
	geocodingService *geo.GeocodingService,
	logger *slog.Logger,
	calendarID string,
) *Processor {
	return &Processor{
		calendarService:  calendarService,
		emailClient:      emailClient,
		gmailSender:      gmailSender,
		geocodingService: geocodingService,
		logger:           logger,
		calendarID:       calendarID,
	}
}

// ProcessResult contains the result of processing a lead
type ProcessResult struct {
	ReferenceNumber string
	Success         bool
	EmailSent       bool
	EmailError      *string
	Estimate        float64
	CalendarCreated bool
	CalendarError   *string
	EventID         *string
	Lat             *float64
	Long            *float64
	FullAddress     *string
	GeoError        *string
}

// ProcessLead processes a transformed lead through the complete workflow:
// 1. Calculate estimate
// 2. Generate quote ID (commented out for now)
// 3. Create calendar event
// 4. Send quote email
// 5. Geocode address
func (p *Processor) ProcessLead(ctx context.Context, data *util.TransformedLeadData) (*ProcessResult, error) {
	result := &ProcessResult{
		Success: true,
	}

	// Step 1: Calculate estimate
	p.logger.Debug("calculating estimate",
		"eventDate", data.EventDate,
		"duration", data.Duration,
		"numHelpers", data.NumHelpers,
	)
	estimate, err := pricing.CalculateEstimate(data.EventDate, data.Duration, data.NumHelpers)
	if err != nil {
		p.logger.Error("failed to calculate estimate", "error", err)
		return nil, fmt.Errorf("failed to calculate estimate: %w", err)
	}
	result.Estimate = estimate.TotalCost
	p.logger.Info("estimate calculated", "totalCost", estimate.TotalCost)

	// Step 2: Generate quote ID (commented out for now as requested)
	// dateKey := data.EventDate.Format("2006-01-02")
	// referenceNumber := util.GenerateShortQuoteID(data.Email, dateKey)
	// result.ReferenceNumber = referenceNumber
	result.ReferenceNumber = "TBD" // Placeholder

	// Step 3: Create calendar event
	if p.calendarService != nil {
		p.logger.Debug("creating calendar event",
			"clientName", data.ClientName,
			"eventDate", data.EventDateStr,
			"eventTime", data.EventTime,
		)
		calendarReq := &calendar.CreateEventRequest{
			CalendarID: p.calendarID,
			ClientName: data.ClientName,
			Occasion:   data.Occasion,
			GuestCount: data.GuestCount,
			EventDate:  data.EventDateStr,
			EventTime:  data.EventTime,
			Phone:      data.Phone,
			Location:   data.EventLocation,
			NumHelpers: data.NumHelpers,
			Duration:   data.Duration,
			TotalCost:  estimate.TotalCost,
			EmailID:    data.Email,
			ThreadID:   "",
			DataSource: "zapier",
			Status:     "Pending",
		}

		calendarResult, err := p.calendarService.CreateEvent(ctx, calendarReq)
		if err != nil {
			errMsg := err.Error()
			result.CalendarError = &errMsg
			p.logger.Warn("failed to create calendar event", "error", err)
		} else if calendarResult.Error != "" {
			result.CalendarError = &calendarResult.Error
			p.logger.Warn("calendar event creation failed", "error", calendarResult.Error)
		} else {
			result.CalendarCreated = true
			result.EventID = &calendarResult.EventID
			p.logger.Info("calendar event created", "eventId", calendarResult.EventID)
		}
	} else {
		p.logger.Warn("calendar service not available, skipping calendar event creation")
	}

	// Step 4: Send quote email
	if p.emailClient != nil || p.gmailSender != nil {
		p.logger.Debug("sending quote email", "to", data.Email)
		emailSent, emailErr := p.sendQuoteEmail(ctx, data, estimate)
		result.EmailSent = emailSent
		if emailErr != "" {
			result.EmailError = &emailErr
			p.logger.Warn("failed to send quote email", "error", emailErr)
		} else {
			p.logger.Info("quote email sent successfully", "to", data.Email)
		}
	} else {
		errMsg := "email service not configured"
		result.EmailError = &errMsg
		p.logger.Warn("email service not available, skipping quote email")
	}

	// Step 5: Geocode address
	if p.geocodingService != nil && data.EventLocation != "" {
		p.logger.Debug("geocoding address", "location", data.EventLocation)
		geoResult, err := p.geocodingService.GetLatLng(ctx, data.EventLocation)
		if err != nil {
			errMsg := err.Error()
			result.GeoError = &errMsg
			p.logger.Warn("geocoding failed", "error", err, "address", data.EventLocation)
		} else {
			result.Lat = &geoResult.Lat
			result.Long = &geoResult.Lng
			fullAddr := geoResult.FullAddress
			result.FullAddress = &fullAddr
			p.logger.Info("address geocoded successfully",
				"lat", geoResult.Lat,
				"lng", geoResult.Lng,
				"fullAddress", geoResult.FullAddress,
			)
		}
	} else {
		p.logger.Debug("geocoding service not available or location empty, skipping geocoding")
	}

	return result, nil
}

// sendQuoteEmail sends the quote email
func (p *Processor) sendQuoteEmail(ctx context.Context, data *util.TransformedLeadData, estimate *pricing.EstimateResult) (bool, string) {
	// Determine rate label
	rateLabel := "Base Rate"
	if estimate.SpecialLabel != nil {
		rateLabel = *estimate.SpecialLabel
	}

	// Format date for email
	dateForEmail := formatDateForEmail(data.EventDate)

	// Calculate deposit from total cost
	estimateCents := util.DollarsToCents(estimate.TotalCost)
	depositCalc := stripe.CalculateDepositFromEstimate(estimateCents)
	depositAmount := util.CentsToDollars(depositCalc.Value)

	// Generate email HTML
	emailData := util.QuoteEmailData{
		ClientName:    data.ClientName,
		EventDate:     dateForEmail,
		EventTime:     data.EventTime,
		EventLocation: data.EventLocation,
		Occasion:      data.Occasion,
		GuestCount:    data.GuestCount,
		Helpers:       data.NumHelpers,
		Hours:         data.Duration,
		BaseRate:      estimate.BasePerHelper,
		HourlyRate:    estimate.ExtraPerHourPerHelper,
		TotalCost:     estimate.TotalCost,
		DepositAmount: depositAmount,
		RateLabel:     rateLabel,
	}

	htmlBody := util.GenerateQuoteEmailHTML(emailData)
	subject := fmt.Sprintf("Party Helpers for %s - %s - Estimate & Details for %s",
		data.Occasion, dateForEmail, data.ClientName)

	if data.DryRun {
		subject = "Dry Run - " + subject
	}

	emailReq := &ports.SendEmailRequest{
		To:       data.Email,
		Subject:  subject,
		HTMLBody: htmlBody,
		FromName: "STL Party Helpers Team",
	}

	var emailResult *ports.SendEmailResult
	var err error

	if p.gmailSender != nil {
		emailResult, err = p.gmailSender.SendEmail(ctx, emailReq)
	} else if p.emailClient != nil {
		emailResult, err = p.emailClient.SendEmail(ctx, emailReq)
	} else {
		return false, "no email service available"
	}

	if err != nil {
		return false, err.Error()
	}

	if !emailResult.Success {
		if emailResult.Error != nil {
			return false, *emailResult.Error
		}
		return false, "unknown error"
	}

	return true, ""
}

// formatDateForEmail formats a date for email display
func formatDateForEmail(date time.Time) string {
	return date.Format("January 2, 2006")
}

