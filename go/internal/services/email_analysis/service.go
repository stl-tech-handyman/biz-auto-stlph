package email_analysis

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bizops360/go-api/internal/infra/email"
	"github.com/bizops360/go-api/internal/infra/sheets"
	"google.golang.org/api/gmail/v1"
)

// Service handles email analysis
type Service struct {
	gmailClient   *email.GmailSender
	sheetsClient  *sheets.SheetsClient
	logger        *slog.Logger
	spreadsheetID string
}

// NewService creates a new email analysis service
func NewService(logger *slog.Logger) (*Service, error) {
	gmailClient, err := email.NewGmailSender()
	if err != nil {
		return nil, fmt.Errorf("failed to create gmail client: %w", err)
	}

	sheetsClient, err := sheets.NewSheetsClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create sheets client: %w", err)
	}

	return &Service{
		gmailClient:  gmailClient,
		sheetsClient: sheetsClient,
		logger:       logger,
	}, nil
}

// AnalysisState tracks processing state
type AnalysisState struct {
	LastProcessedIndex int       `json:"last_processed_index"`
	TotalProcessed     int       `json:"total_processed"`
	LastRun            time.Time `json:"last_run"`
	SpreadsheetID      string    `json:"spreadsheet_id"`
}

// AnalyzeEmailsRequest is the request to analyze emails
type AnalyzeEmailsRequest struct {
	MaxEmails     int    `json:"max_emails"`
	Query         string `json:"query"`
	Resume        bool   `json:"resume"`
	SpreadsheetID string `json:"spreadsheet_id"`
}

// AnalyzeEmailsResponse is the response
type AnalyzeEmailsResponse struct {
	Processed      int    `json:"processed"`
	Skipped        int    `json:"skipped"`
	Total          int    `json:"total"`
	SpreadsheetID  string `json:"spreadsheet_id"`
	SpreadsheetURL string `json:"spreadsheet_url"`
	NextIndex      int    `json:"next_index"`
	HasMore        bool   `json:"has_more"`
}

// AnalyzeEmails processes emails and writes to sheets
func (s *Service) AnalyzeEmails(ctx context.Context, req AnalyzeEmailsRequest) (*AnalyzeEmailsResponse, error) {
	s.logger.Info("starting email analysis",
		"max_emails", req.MaxEmails,
		"query", req.Query,
		"resume", req.Resume,
	)

	logToFile("Starting email analysis", map[string]interface{}{
		"max_emails": req.MaxEmails,
		"query":      req.Query,
		"resume":     req.Resume,
	})

	// Get or create spreadsheet
	spreadsheetID := req.SpreadsheetID
	if spreadsheetID == "" {
		var err error
		spreadsheetID, err = s.sheetsClient.GetOrCreateSpreadsheet(ctx, "Email Revenue Analytics")
		if err != nil {
			return nil, fmt.Errorf("failed to get/create spreadsheet: %w", err)
		}
		s.spreadsheetID = spreadsheetID
		logToFile("Created spreadsheet", map[string]interface{}{
			"spreadsheet_id": spreadsheetID,
		})
	}

	// Initialize sheets
	if err := s.initializeSheets(ctx, spreadsheetID); err != nil {
		return nil, fmt.Errorf("failed to initialize sheets: %w", err)
	}

	// Get state from Sheets
	state, err := s.getState(ctx, spreadsheetID)
	if err != nil {
		s.logger.Warn("failed to get state, starting fresh", "error", err)
		state = &AnalysisState{SpreadsheetID: spreadsheetID}
	}

	if !req.Resume {
		state.LastProcessedIndex = 0
		state.TotalProcessed = 0
	}

	// Search emails
	query := req.Query
	if query == "" {
		query = `from:zapier.com OR from:forms OR subject:"New Lead" OR subject:"Form Submission" OR subject:"Quote"`
	}

	userEmail := s.gmailClient.GetFromEmail()
	processed := 0
	skipped := 0
	startIndex := state.LastProcessedIndex

	// Process in batches
	batchSize := 50 // Smaller batches for testing
	maxBatches := 20
	if req.MaxEmails > 0 {
		maxBatches = (req.MaxEmails + batchSize - 1) / batchSize
		if maxBatches > 20 {
			maxBatches = 20 // Limit to 1000 emails per run for safety
		}
	}

	var processedEmails []*EmailData
	var hasMore bool

	for batchNum := 0; batchNum < maxBatches; batchNum++ {
		s.logger.Info("processing batch",
			"batch", batchNum+1,
			"max_batches", maxBatches,
		)

		logToFile("Processing batch", map[string]interface{}{
			"batch":      batchNum + 1,
			"max_batches": maxBatches,
		})

		// Search Gmail
		call := s.gmailClient.GetService().Users.Messages.List(userEmail).
			Q(query).
			MaxResults(int64(batchSize))

		resp, err := call.Context(ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("failed to search emails: %w", err)
		}

		if len(resp.Messages) == 0 {
			s.logger.Info("no more messages found")
			hasMore = false
			break
		}

		hasMore = resp.NextPageToken != ""

		// Process each message
		for _, msg := range resp.Messages {
			emailData, err := s.processMessage(ctx, msg.Id)
			if err != nil {
				s.logger.Warn("failed to process message", "id", msg.Id, "error", err)
				logToFile("Failed to process message", map[string]interface{}{
					"message_id": msg.Id,
					"error":      err.Error(),
				})
				skipped++
				continue
			}

			if emailData != nil && !emailData.IsTest {
				processedEmails = append(processedEmails, emailData)
				processed++
			} else {
				skipped++
			}
		}

		// Write batch to sheets when we have enough
		if len(processedEmails) >= 25 {
			if err := s.writeBatchToSheets(ctx, spreadsheetID, processedEmails); err != nil {
				s.logger.Error("failed to write batch", "error", err)
				logToFile("Failed to write batch", map[string]interface{}{
					"error": err.Error(),
				})
			} else {
				logToFile("Wrote batch to sheets", map[string]interface{}{
					"count": len(processedEmails),
				})
				processedEmails = []*EmailData{} // Clear batch
			}
		}

		// Update state
		state.LastProcessedIndex = startIndex + (batchNum+1)*batchSize
		state.TotalProcessed += processed
		state.LastRun = time.Now()
		state.SpreadsheetID = spreadsheetID

		if err := s.saveState(ctx, spreadsheetID, state); err != nil {
			s.logger.Warn("failed to save state", "error", err)
		}

		// Check if we've hit the limit
		if req.MaxEmails > 0 && processed >= req.MaxEmails {
			break
		}

		// Small delay to avoid rate limits
		time.Sleep(200 * time.Millisecond)
	}

	// Write remaining batch
	if len(processedEmails) > 0 {
		if err := s.writeBatchToSheets(ctx, spreadsheetID, processedEmails); err != nil {
			s.logger.Error("failed to write final batch", "error", err)
		} else {
			logToFile("Wrote final batch", map[string]interface{}{
				"count": len(processedEmails),
			})
		}
	}

	spreadsheetURL := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s", spreadsheetID)

	s.logger.Info("email analysis complete",
		"processed", processed,
		"skipped", skipped,
		"total", state.TotalProcessed,
	)

	logToFile("Analysis complete", map[string]interface{}{
		"processed": processed,
		"skipped":   skipped,
		"total":     state.TotalProcessed,
	})

	return &AnalyzeEmailsResponse{
		Processed:      processed,
		Skipped:        skipped,
		Total:          state.TotalProcessed,
		SpreadsheetID:  spreadsheetID,
		SpreadsheetURL: spreadsheetURL,
		NextIndex:      state.LastProcessedIndex,
		HasMore:        hasMore,
	}, nil
}

// EmailData represents extracted email data
type EmailData struct {
	EmailID        string
	ThreadID      string
	Date          time.Time
	FromEmail     string
	Subject       string
	BodyPreview   string
	IsTest        bool
	IsConfirmation bool
	ClientEmail   string
	EventDate     string
	TotalCost     float64
	Rate          float64
	Hours         string
	Helpers       string
	Occasion      string
	Status        string
	Guests        string
	Deposit       float64
	EmailType     string
	ConversationID string
	MessageNumber int
}

// processMessage extracts data from a Gmail message
func (s *Service) processMessage(ctx context.Context, messageID string) (*EmailData, error) {
	userEmail := s.gmailClient.GetFromEmail()
	msg, err := s.gmailClient.GetService().Users.Messages.Get(userEmail, messageID).Format("full").Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	// Extract headers
	var from, subject, dateStr string
	for _, header := range msg.Payload.Headers {
		switch header.Name {
		case "From":
			from = header.Value
		case "Subject":
			subject = header.Value
		case "Date":
			dateStr = header.Value
		}
	}

	// Parse date
	date, _ := time.Parse(time.RFC1123Z, dateStr)
	if date.IsZero() {
		date = time.Unix(msg.InternalDate/1000, 0)
	}

	// Extract body
	body := s.extractBody(msg.Payload)
	bodyPreview := body
	if len(bodyPreview) > 300 {
		bodyPreview = bodyPreview[:300]
	}

	// Extract data
	combined := strings.ToLower(subject + " " + body)
	isTest := s.detectTest(combined)
	isConfirmation := s.detectConfirmation(combined)
	clientEmail := s.extractClientEmail(from, body)
	eventData := s.extractEventData(body, subject)
	pricingData := s.extractPricingData(body, subject)
	emailType := s.classifyEmailType(subject, body, from)

	return &EmailData{
		EmailID:        messageID,
		ThreadID:       msg.ThreadId,
		Date:           date,
		FromEmail:      from,
		Subject:        subject,
		BodyPreview:    bodyPreview,
		IsTest:         isTest,
		IsConfirmation: isConfirmation,
		ClientEmail:    clientEmail,
		EventDate:      eventData.EventDate,
		TotalCost:      pricingData.TotalCost,
		Rate:           pricingData.Rate,
		Hours:          eventData.Hours,
		Helpers:        eventData.Helpers,
		Occasion:       eventData.Occasion,
		Status:         s.getStatus(isConfirmation, isTest),
		Guests:         eventData.Guests,
		Deposit:        pricingData.Deposit,
		EmailType:      emailType,
		ConversationID: msg.ThreadId + "_" + clientEmail,
		MessageNumber:  1,
	}, nil
}

// Helper methods for extraction
func (s *Service) extractBody(payload *gmail.MessagePart) string {
	if payload == nil {
		return ""
	}

	if payload.MimeType == "text/plain" && payload.Body != nil && payload.Body.Data != "" {
		data, _ := decodeBase64URL(payload.Body.Data)
		return string(data)
	}

	if payload.MimeType == "text/html" && payload.Body != nil && payload.Body.Data != "" {
		data, _ := decodeBase64URL(payload.Body.Data)
		return stripHTML(string(data))
	}

	for _, part := range payload.Parts {
		if body := s.extractBody(part); body != "" {
			return body
		}
	}

	return ""
}

func (s *Service) detectTest(text string) bool {
	testKeywords := []string{"test", "testing", "demo", "sample", "fake"}
	for _, keyword := range testKeywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

func (s *Service) detectConfirmation(text string) bool {
	confKeywords := []string{"confirm", "confirmed", "confirmation", "accepted", "approved"}
	for _, keyword := range confKeywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

func (s *Service) extractClientEmail(from, body string) string {
	emailRegex := regexp.MustCompile(`([a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\.[a-zA-Z0-9_-]+)`)
	matches := emailRegex.FindAllString(body, -1)
	for _, match := range matches {
		if !strings.Contains(match, "zapier") && !strings.Contains(match, "noreply") && !strings.Contains(match, "no-reply") {
			return match
		}
	}
	// Extract from "From" header if not found in body
	if from != "" {
		emailMatch := emailRegex.FindString(from)
		if emailMatch != "" {
			return emailMatch
		}
	}
	return ""
}

type EventData struct {
	EventDate string
	Hours     string
	Helpers   string
	Occasion  string
	Guests    string
}

func (s *Service) extractEventData(body, subject string) EventData {
	// Simplified - extract date patterns
	datePatterns := []string{
		`(\d{1,2}[/-]\d{1,2}[/-]\d{2,4})`,
		`(January|February|March|April|May|June|July|August|September|October|November|December)\s+\d{1,2}`,
	}
	
	var eventDate string
	for _, pattern := range datePatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(body + " " + subject)
		if len(matches) > 0 {
			eventDate = matches[0]
			break
		}
	}

	return EventData{
		EventDate: eventDate,
		Hours:     "",
		Helpers:   "",
		Occasion:  "",
		Guests:    "",
	}
}

type PricingData struct {
	TotalCost float64
	Rate      float64
	Deposit   float64
}

func (s *Service) extractPricingData(body, subject string) PricingData {
	// Extract dollar amounts
	dollarRegex := regexp.MustCompile(`\$?(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)`)
	matches := dollarRegex.FindAllString(body+" "+subject, -1)

	var totalCost, rate, deposit float64
	if len(matches) > 0 {
		totalCost, _ = parseDollarAmount(matches[0])
	}
	if len(matches) > 1 {
		rate, _ = parseDollarAmount(matches[1])
	}
	if len(matches) > 2 {
		deposit, _ = parseDollarAmount(matches[2])
	}

	return PricingData{
		TotalCost: totalCost,
		Rate:      rate,
		Deposit:   deposit,
	}
}

func (s *Service) classifyEmailType(subject, body, from string) string {
	stlphKeywords := []string{"party", "event", "helpers", "booking", "quote", "stl party"}
	combined := strings.ToLower(subject + " " + body)

	for _, keyword := range stlphKeywords {
		if strings.Contains(combined, keyword) {
			return "STLPH"
		}
	}
	return "Other"
}

func (s *Service) getStatus(isConfirmation, isTest bool) string {
	if isTest {
		return "Test"
	}
	if isConfirmation {
		return "Confirmed"
	}
	return "Pending"
}

// Utility functions
func decodeBase64URL(data string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(data)
}

func stripHTML(html string) string {
	re := regexp.MustCompile(`<[^>]+>`)
	return re.ReplaceAllString(html, " ")
}

func parseDollarAmount(str string) (float64, error) {
	cleaned := strings.ReplaceAll(strings.ReplaceAll(str, "$", ""), ",", "")
	return strconv.ParseFloat(cleaned, 64)
}

// State management using Sheets
func (s *Service) getState(ctx context.Context, spreadsheetID string) (*AnalysisState, error) {
	rows, err := s.sheetsClient.GetRows(ctx, spreadsheetID, "State")
	if err != nil || len(rows) < 2 {
		return &AnalysisState{SpreadsheetID: spreadsheetID}, nil
	}

	state := &AnalysisState{SpreadsheetID: spreadsheetID}
	if len(rows[1]) > 0 {
		if idx, err := strconv.Atoi(fmt.Sprintf("%v", rows[1][0])); err == nil {
			state.LastProcessedIndex = idx
		}
	}
	if len(rows[1]) > 1 {
		if total, err := strconv.Atoi(fmt.Sprintf("%v", rows[1][1])); err == nil {
			state.TotalProcessed = total
		}
	}
	if len(rows[1]) > 2 {
		if lastRun, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", rows[1][2])); err == nil {
			state.LastRun = lastRun
		}
	}
	if len(rows[1]) > 3 {
		state.SpreadsheetID = fmt.Sprintf("%v", rows[1][3])
	}

	return state, nil
}

func (s *Service) saveState(ctx context.Context, spreadsheetID string, state *AnalysisState) error {
	// Create or update State sheet
	_ = s.sheetsClient.CreateSheet(ctx, spreadsheetID, "State")
	
	// Set headers if needed
	headers := []string{"LastProcessedIndex", "TotalProcessed", "LastRun", "SpreadsheetID"}
	_ = s.sheetsClient.SetHeaders(ctx, spreadsheetID, "State", headers)

	// Update row 2 (index 1)
	values := []interface{}{
		state.LastProcessedIndex,
		state.TotalProcessed,
		state.LastRun.Format(time.RFC3339),
		state.SpreadsheetID,
	}

	return s.sheetsClient.UpdateRow(ctx, spreadsheetID, "State", 1, values)
}

// Initialize sheets
func (s *Service) initializeSheets(ctx context.Context, spreadsheetID string) error {
	sheets := []string{"Raw Data", "State", "Processing Log"}

	for _, sheetName := range sheets {
		_ = s.sheetsClient.CreateSheet(ctx, spreadsheetID, sheetName)

		// Set headers for Raw Data
		if sheetName == "Raw Data" {
			headers := []string{
				"Email ID", "Thread ID", "Date", "From Email", "Subject", "Body Preview",
				"Is Test", "Is Confirmation", "Client Email", "Event Date",
				"Total Cost", "Rate", "Hours", "Helpers", "Occasion", "Status",
				"Guests", "Deposit", "Email Type", "Conversation ID", "Message Number",
			}
			_ = s.sheetsClient.SetHeaders(ctx, spreadsheetID, sheetName, headers)
		}

		// Set headers for State
		if sheetName == "State" {
			headers := []string{"LastProcessedIndex", "TotalProcessed", "LastRun", "SpreadsheetID"}
			_ = s.sheetsClient.SetHeaders(ctx, spreadsheetID, sheetName, headers)
		}

		// Set headers for Processing Log
		if sheetName == "Processing Log" {
			headers := []string{"Timestamp", "Status", "Emails Processed", "Notes"}
			_ = s.sheetsClient.SetHeaders(ctx, spreadsheetID, sheetName, headers)
		}
	}

	return nil
}

// Write batch to sheets
func (s *Service) writeBatchToSheets(ctx context.Context, spreadsheetID string, emails []*EmailData) error {
	values := make([][]interface{}, len(emails))
	for i, email := range emails {
		values[i] = []interface{}{
			email.EmailID,
			email.ThreadID,
			email.Date.Format(time.RFC3339),
			email.FromEmail,
			email.Subject,
			email.BodyPreview,
			strconv.FormatBool(email.IsTest),
			strconv.FormatBool(email.IsConfirmation),
			email.ClientEmail,
			email.EventDate,
			email.TotalCost,
			email.Rate,
			email.Hours,
			email.Helpers,
			email.Occasion,
			email.Status,
			email.Guests,
			email.Deposit,
			email.EmailType,
			email.ConversationID,
			email.MessageNumber,
		}
	}

	return s.sheetsClient.AppendRows(ctx, spreadsheetID, "Raw Data", values)
}

// GetStatus returns current analysis status
func (s *Service) GetStatus(ctx context.Context, spreadsheetID string) (*AnalysisState, error) {
	if spreadsheetID == "" {
		return &AnalysisState{}, nil
	}
	return s.getState(ctx, spreadsheetID)
}

// Log helper
func logToFile(message string, data map[string]interface{}) {
	logPath := "c:\\Users\\Alexey\\Code\\biz-operating-system\\stlph\\.cursor\\debug.log"
	if logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		logEntry := map[string]interface{}{
			"sessionId":    "email-analysis",
			"runId":        "run1",
			"hypothesisId": "EMAIL_ANALYSIS",
			"location":     "service.go",
			"message":      message,
			"data":         data,
			"timestamp":    time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
}
