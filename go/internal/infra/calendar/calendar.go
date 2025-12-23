package calendar

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// CalendarService handles Google Calendar event creation
type CalendarService struct {
	service    *calendar.Service
	calendarID string
}

// NewCalendarService creates a new calendar service
func NewCalendarService(calendarID string) (*CalendarService, error) {
	// Get credentials from environment (same pattern as Gmail)
	credentialsJSON := os.Getenv("GMAIL_CREDENTIALS_JSON")
	if credentialsJSON == "" {
		// If no credentials, return a service that will fail gracefully
		return &CalendarService{
			calendarID: calendarID,
		}, nil
	}

	// Try to read from file if it's a path, otherwise use as JSON string
	var credsData []byte
	if _, err := os.Stat(credentialsJSON); err == nil {
		credsData, err = os.ReadFile(credentialsJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to read credentials file: %w", err)
		}
	} else {
		credsData = []byte(credentialsJSON)
	}

	// Try JWT config first (for service accounts)
	config, err := google.JWTConfigFromJSON(credsData, calendar.CalendarScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse calendar credentials: %w", err)
	}

	// Get the email address to impersonate (required for domain-wide delegation)
	impersonateEmail := os.Getenv("GMAIL_FROM")
	if impersonateEmail == "" {
		impersonateEmail = config.Email
		if impersonateEmail == "" {
			return nil, fmt.Errorf("GMAIL_FROM environment variable must be set for service account with domain-wide delegation")
		}
	}

	// Set the subject (user to impersonate) for domain-wide delegation
	config.Subject = impersonateEmail

	// Create Calendar service with JWT config
	ctx := context.Background()
	client := config.Client(ctx)

	service, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create Calendar service: %w", err)
	}

	return &CalendarService{
		service:    service,
		calendarID: calendarID,
	}, nil
}

// CreateEventRequest contains data for creating a calendar event
type CreateEventRequest struct {
	CalendarID string  // Optional, overrides default
	ClientName string
	Occasion   string
	GuestCount int
	EventDate  string
	EventTime  string
	Phone      string
	Location   string
	NumHelpers int
	Duration   float64
	TotalCost  float64
	EmailID    string
	ThreadID   string
	DataSource string
	Status     string // Optional, defaults to "Pending"
}

// CreateEventResult contains the result of creating a calendar event
type CreateEventResult struct {
	EventID string
	Error   string
	Success bool
}

// CreateEvent creates a calendar event
// Matches the Apps Script createEvent function behavior
func (c *CalendarService) CreateEvent(ctx context.Context, req *CreateEventRequest) (*CreateEventResult, error) {
	// If service is not initialized (no credentials), return error
	if c.service == nil {
		return &CreateEventResult{
			Error: "Calendar service not initialized - GMAIL_CREDENTIALS_JSON not configured",
		}, nil
	}

	// Validate required parameters
	if req.EventDate == "" || req.EventTime == "" || req.Duration < 1 {
		return &CreateEventResult{
			Error: "Missing required parameters: eventDate, eventTime, or duration",
		}, nil
	}

	// Use provided calendar ID or default
	calendarID := req.CalendarID
	if calendarID == "" {
		calendarID = c.calendarID
	}
	if calendarID == "" {
		return &CreateEventResult{
			Error: "Calendar ID is required",
		}, nil
	}

	// Parse date and time
	startTime, err := parseDateTime(req.EventDate, req.EventTime)
	if err != nil {
		return &CreateEventResult{
			Error: fmt.Sprintf("Invalid date/time format: %v", err),
		}, nil
	}

	// Calculate end time
	endTime := startTime.Add(time.Duration(req.Duration) * time.Hour)
	if startTime.After(endTime) || startTime.Equal(endTime) {
		return &CreateEventResult{
			Error: "Start time must be before end time",
		}, nil
	}

	// Set defaults
	clientName := req.ClientName
	if clientName == "" {
		clientName = "Unknown Client"
	}
	occasion := req.Occasion
	if occasion == "" {
		occasion = "No Occasion"
	}
	location := req.Location
	if location == "" {
		location = "TBD"
	}
	emailID := req.EmailID
	if emailID == "" {
		emailID = "Not Provided"
	}
	phone := req.Phone
	dataSource := req.DataSource
	if dataSource == "" {
		dataSource = "Not Provided"
	}
	status := req.Status
	if status == "" {
		status = "Pending"
	}

	// Build event title (matching Apps Script format)
	eventTitle := fmt.Sprintf("💲%.0f (%d-👧, %.0f-⏰) - %s - %s",
		req.TotalCost, req.NumHelpers, req.Duration, clientName, occasion)

	// Build event description (matching Apps Script format)
	var emailLink string
	if req.ThreadID != "" {
		emailLink = fmt.Sprintf("📧 <a href=\"https://mail.google.com/mail/u/0/#inbox/%s\" target=\"_blank\">Open Lead Email</a>", req.ThreadID)
	} else {
		emailLink = "⚠️ No Email Thread Available"
	}

	eventDescription := fmt.Sprintf(`🧑‍🤝‍🧑 Guest Count: %d
📌 Status: %s
📧 Email: %s
📧 Phone: %s
📧 Source: %s
🔗 %s`,
		req.GuestCount, status, emailID, phone, dataSource, emailLink)

	// Create calendar event
	event := &calendar.Event{
		Summary:     eventTitle,
		Description: eventDescription,
		Location:    location,
		Start: &calendar.EventDateTime{
			DateTime: startTime.Format(time.RFC3339),
			TimeZone: "America/Chicago", // STL timezone
		},
		End: &calendar.EventDateTime{
			DateTime: endTime.Format(time.RFC3339),
			TimeZone: "America/Chicago",
		},
	}

	// Add extended properties if threadId is provided
	if req.ThreadID != "" {
		event.ExtendedProperties = &calendar.EventExtendedProperties{
			Private: map[string]string{
				"emailThreadId": req.ThreadID,
			},
		}
	}

	// Insert event into calendar
	createdEvent, err := c.service.Events.Insert(calendarID, event).Context(ctx).Do()
	if err != nil {
		return &CreateEventResult{
			Error: fmt.Sprintf("Failed to create calendar event: %v", err),
		}, nil
	}

	return &CreateEventResult{
		EventID: createdEvent.Id,
		Success: true,
	}, nil
}

// parseDateTime parses eventDate and eventTime into a time.Time
// Handles various date/time formats like the Apps Script code
func parseDateTime(eventDate, eventTime string) (time.Time, error) {
	// Try combining date and time first
	combined := fmt.Sprintf("%s %s", strings.TrimSpace(eventDate), strings.TrimSpace(eventTime))
	
	// Try multiple formats
	formats := []string{
		"2006-01-02 3:04 PM",
		"2006-01-02 15:04",
		"2006-01-02 3:04PM",
		"2006-01-02 15:04:05",
		time.RFC3339,
		"January 2, 2006 3:04 PM",
		"Jan 2, 2006 3:04 PM",
		"2006/01/02 3:04 PM",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, combined); err == nil {
			// If no timezone, assume America/Chicago (STL)
			loc, _ := time.LoadLocation("America/Chicago")
			return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc), nil
		}
	}

	// Fallback: try parsing just the date
	dateFormats := []string{
		"2006-01-02",
		"January 2, 2006",
		"Jan 2, 2006",
		"2006/01/02",
	}

	for _, format := range dateFormats {
		if t, err := time.Parse(format, eventDate); err == nil {
			// Parse time separately
			timeFormats := []string{"3:04 PM", "15:04", "3:04PM", "15:04:05"}
			for _, tf := range timeFormats {
				if parsedTime, err := time.Parse(tf, eventTime); err == nil {
					loc, _ := time.LoadLocation("America/Chicago")
					return time.Date(t.Year(), t.Month(), t.Day(),
						parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, loc), nil
				}
			}
			// If time parsing fails, use noon as default
			loc, _ := time.LoadLocation("America/Chicago")
			return time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, loc), nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date/time: %s %s", eventDate, eventTime)
}

