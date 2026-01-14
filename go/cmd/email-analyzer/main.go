package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	oldEmail = "stlpartyhelpers@gmail.com"
	newEmail = "team@stlpartyhelpers.com"
	
	// Job configuration
	defaultJobID   = "JOB-1-INITIAL-INDEXING"
	defaultJobName = "Initial Indexing and Classification Systematization"
)

func main() {
	var (
		maxEmails     = flag.Int("max", 0, "Maximum emails to process (0 = process all available)")
		query         = flag.String("query", "", "Gmail search query (default: searches for form submissions)")
		spreadsheetID = flag.String("spreadsheet", "", "Existing spreadsheet ID (creates new if empty)")
		resume        = flag.Bool("resume", false, "Resume from last position")
		rebuild       = flag.Bool("rebuild", false, "Rebuild derived sheets (Email Mapping, State) from Raw Data without reading Gmail")
		verbose       = flag.Bool("v", false, "Verbose logging")
		batchSize     = flag.Int("batch", 50, "Batch size for processing")
		delay         = flag.Int("delay", 200, "Delay between batches in milliseconds")
		all           = flag.Bool("all", false, "Process all emails (equivalent to -max 0)")
		idempotent    = flag.Bool("idempotent", false, "Recreate sheets if they exist (idempotent mode)")
		agentID       = flag.String("agent", "", "Agent ID for concurrent processing (auto-generated if empty)")
		jobID         = flag.String("job", defaultJobID, "Job ID for this analysis run")
		jobName       = flag.String("job-name", defaultJobName, "Job name/description")
		workers       = flag.Int("workers", 5, "Number of concurrent workers (recommended: 3-5)")
		testMode      = flag.Bool("test", false, "Test mode: process only 10 emails and verify writes")
	)
	flag.Parse()

	if *all {
		*maxEmails = 0
	}
	
	// Test mode: process only 10 emails for verification
	if *testMode {
		*maxEmails = 10
		*verbose = true
		fmt.Printf("🧪 TEST MODE: Will process only 10 emails and verify writes\n")
	}
	
	// Test mode: process only 10 emails for verification
	if *testMode {
		*maxEmails = 10
		*verbose = true
		fmt.Printf("🧪 TEST MODE: Will process only 10 emails and verify writes\n")
	}
	
	// Generate agent ID if not provided
	if *agentID == "" {
		*agentID = fmt.Sprintf("agent-%d-%d", time.Now().Unix(), os.Getpid())
	}
	
	// Validate workers
	if *workers < 1 {
		*workers = 1
	}
	if *workers > 10 {
		fmt.Printf("⚠️  Warning: Workers > 10 may hit Gmail API rate limits. Recommended: 3-5\n")
	}

	ctx := context.Background()

	// Initialize clients
	gmailService, err := initGmail(ctx)
	if err != nil {
		log.Fatalf("Failed to init Gmail: %v", err)
	}

	sheetsService, err := initSheets(ctx)
	if err != nil {
		log.Fatalf("Failed to init Sheets: %v", err)
	}

	// Get or create spreadsheet
	if *spreadsheetID == "" {
		if *rebuild {
			log.Fatalf("rebuild mode requires -spreadsheet ID")
		}
		spreadsheet, err := sheetsService.Spreadsheets.Create(&sheets.Spreadsheet{
			Properties: &sheets.SpreadsheetProperties{
				Title: fmt.Sprintf("Email Analysis - %s", time.Now().Format("2006-01-02 15:04")),
			},
		}).Context(ctx).Do()
		if err != nil {
			log.Fatalf("Failed to create spreadsheet: %v", err)
		}
		*spreadsheetID = spreadsheet.SpreadsheetId
		fmt.Printf("\n✅ Created spreadsheet: %s\n", spreadsheet.SpreadsheetUrl)
		fmt.Printf("📊 Spreadsheet ID: %s\n", *spreadsheetID)
		fmt.Printf("🌐 Dashboard URL: http://localhost:8080/email-analysis-dashboard.html?spreadsheet_id=%s\n\n", *spreadsheetID)
	}

	// Acquire lock before processing
	lockAcquired, err := acquireLock(ctx, sheetsService, *spreadsheetID, *agentID, *idempotent)
	if err != nil {
		log.Fatalf("Failed to acquire lock: %v", err)
	}
	if !lockAcquired {
		log.Fatalf("❌ Lock already held by another agent. Wait for it to complete or cleanup expired locks.")
	}
	
	// Ensure lock is released on exit
	defer func() {
		if err := releaseLock(ctx, sheetsService, *spreadsheetID, *agentID); err != nil {
			fmt.Printf("⚠️  Warning: Failed to release lock: %v\n", err)
		} else {
			fmt.Printf("🔓 Lock released by agent: %s\n", *agentID)
		}
	}()

	// Initialize sheets
	fmt.Printf("📋 Initializing sheets...\n")
	if err := initializeSheets(ctx, sheetsService, *spreadsheetID, *idempotent); err != nil {
		log.Fatalf("Failed to initialize sheets: %v", err)
	}
	fmt.Printf("✅ Sheets initialized successfully\n")
	
	// Test write verification in test mode
	if *testMode {
		fmt.Printf("🧪 TEST: Verifying sheet write capability...\n")
		testRow := [][]interface{}{
			{"TEST_MESSAGE_ID", "TEST_THREAD_ID", time.Now().Format(time.RFC3339), "test@example.com", "test@example.com", "Test Email", "Test body", "false", "false", "false", "false", "STLPH", "1.0", "Website", "Test Page", "test@example.com", "test@example.com", "example.com", "Test User", "", "", "", "", "", "", "", "", "", "0.00", "0.00", "0.00", "0.00", "0.00", "Pending", "test_conversation", "1", "", "false", *jobID, *jobName},
		}
		if err := writeBatch(ctx, sheetsService, *spreadsheetID, testRow); err != nil {
			log.Fatalf("❌ TEST FAILED: Cannot write to spreadsheet: %v", err)
		}
		fmt.Printf("✅ TEST PASSED: Successfully wrote test row to spreadsheet\n")
		
		// Verify we can read it back
		readResp, readErr := sheetsService.Spreadsheets.Values.Get(*spreadsheetID, "Raw Data!A:A").Context(ctx).Do()
		if readErr != nil {
			log.Fatalf("❌ TEST FAILED: Cannot read from spreadsheet: %v", readErr)
		}
		found := false
		for _, row := range readResp.Values {
			if len(row) > 0 && fmt.Sprintf("%v", row[0]) == "TEST_MESSAGE_ID" {
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("❌ TEST FAILED: Test row not found in spreadsheet after write")
		}
		fmt.Printf("✅ TEST PASSED: Successfully verified read-back from spreadsheet\n")
		fmt.Printf("🧪 Test mode verification complete. Proceeding with email processing...\n\n")
	}

	// Rebuild mode: recompute derived sheets from Raw Data and exit
	if *rebuild {
		fmt.Printf("🔧 Rebuild mode: recomputing derived sheets from Raw Data...\n")
		if err := rebuildDerivedFromRawData(ctx, sheetsService, *spreadsheetID); err != nil {
			log.Fatalf("Failed to rebuild derived sheets: %v", err)
		}
		fmt.Printf("✅ Rebuild complete.\n")
		return
	}

	// Get state
	state := getState(ctx, sheetsService, *spreadsheetID)
	
	// Build set of already processed message IDs for fast O(1) lookup
	// Using Gmail message IDs as unique identifiers
	processedIDsSet := make(map[string]bool)
	if *resume && len(state.ProcessedIDs) > 0 {
		for _, id := range state.ProcessedIDs {
			processedIDsSet[id] = true
		}
		fmt.Printf("  🔄 Resuming: Found %d already processed message IDs\n", len(processedIDsSet))
		fmt.Printf("  📍 Will skip these and continue from last position\n")
	} else {
		state.LastIndex = 0
		state.TotalProcessed = 0
		state.ProcessedIDs = []string{}
		fmt.Printf("  🆕 Starting fresh - no previous state found\n")
	}

	// Default query - if -all flag is used, process ALL emails (empty query)
	// Test mode also uses empty query to get all emails
	// Otherwise default to form submissions for backward compatibility
	if *query == "" {
		if *all || *testMode {
			*query = "" // Empty query = all emails in Gmail API
			if *testMode {
				fmt.Printf("  ℹ️  Test mode: Using empty query to process ALL emails\n")
			} else {
				fmt.Printf("  ℹ️  Using empty query to process ALL emails\n")
			}
		} else {
			*query = `from:zapier.com OR from:forms OR subject:"New Lead" OR subject:"Form Submission" OR subject:"Quote"`
		}
	}

	fmt.Printf("Starting analysis...\n")
	fmt.Printf("  Job ID: %s\n", *jobID)
	fmt.Printf("  Job Name: %s\n", *jobName)
	fmt.Printf("  Agent ID: %s\n", *agentID)
	fmt.Printf("  Query: %s\n", *query)
	if *maxEmails == 0 {
		fmt.Printf("  Max emails: ALL (processing all available)\n")
	} else {
		fmt.Printf("  Max emails: %d\n", *maxEmails)
	}
	fmt.Printf("  Batch size: %d\n", *batchSize)
	fmt.Printf("  Resume from index: %d\n", state.LastIndex)
	fmt.Printf("  Total already processed: %d\n", state.TotalProcessed)
	fmt.Printf("  Idempotent mode: %v\n", *idempotent)
	fmt.Printf("  Workers: %d\n", *workers)
	fmt.Printf("  Spreadsheet: https://docs.google.com/spreadsheets/d/%s\n\n", *spreadsheetID)

	// Process emails
	startTime := time.Now()
	processed, skipped, err := processEmails(ctx, gmailService, sheetsService, *spreadsheetID, *query, *maxEmails, state.LastIndex, *batchSize, *delay, *verbose, processedIDsSet, &state.ProcessedIDs, *agentID, *jobID, *jobName, *workers)
	if err != nil {
		log.Fatalf("Failed to process emails: %v", err)
	}

	// Update state with final counts
	state.LastIndex += processed + skipped
	state.TotalProcessed += processed
	state.LastRun = time.Now()
	if err := saveState(ctx, sheetsService, *spreadsheetID, state); err != nil {
		fmt.Printf("⚠️  Warning: Failed to save final state: %v\n", err)
	} else {
		fmt.Printf("✅ Final state saved (%d processed IDs tracked)\n", len(state.ProcessedIDs))
	}

	elapsed := time.Since(startTime)
	rate := float64(0)
	if elapsed.Seconds() > 0 {
		rate = float64(processed) / elapsed.Seconds()
	}

	fmt.Printf("\n" + strings.Repeat("=", 60) + "\n")
	fmt.Printf("✅ PROCESSING COMPLETE!\n")
	fmt.Printf(strings.Repeat("=", 60) + "\n")
	fmt.Printf("  📧 Processed: %d emails\n", processed)
	fmt.Printf("  ⏭️  Skipped: %d emails\n", skipped)
	fmt.Printf("  📊 Total processed (all time): %d\n", state.TotalProcessed)
	fmt.Printf("  🆔 Unique message IDs tracked: %d\n", len(state.ProcessedIDs))
	fmt.Printf("  ⏱️  Time elapsed: %s\n", elapsed.Round(time.Second))
	if rate > 0 {
		fmt.Printf("  🚀 Processing rate: %.1f emails/second\n", rate)
		if processed > 0 {
			fmt.Printf("  ⏳ Average time per email: %s\n", (elapsed / time.Duration(processed)).Round(time.Millisecond))
		}
	}
	fmt.Printf("  📍 Next index: %d\n", state.LastIndex)
	fmt.Printf("  📄 Spreadsheet: https://docs.google.com/spreadsheets/d/%s\n", *spreadsheetID)
	fmt.Printf("\n  💡 Resume: Use -resume flag with this spreadsheet ID to continue safely\n")
	fmt.Printf(strings.Repeat("=", 60) + "\n")
}

type State struct {
	LastIndex      int
	TotalProcessed int
	LastRun        time.Time
	ProcessedIDs   []string `json:"processed_ids"` // Track processed email IDs
	CurrentJobID   string   `json:"current_job_id"`
	CurrentJobName string   `json:"current_job_name"`
}

func initGmail(ctx context.Context) (*gmail.Service, error) {
	credsJSON := os.Getenv("GMAIL_CREDENTIALS_JSON")
	if credsJSON == "" {
		return nil, fmt.Errorf("GMAIL_CREDENTIALS_JSON not set")
	}

	var credsData []byte
	if _, err := os.Stat(credsJSON); err == nil {
		credsData, _ = os.ReadFile(credsJSON)
	} else {
		credsData = []byte(credsJSON)
	}

	config, err := google.JWTConfigFromJSON(credsData, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, err
	}

	// Set subject for domain-wide delegation (same as Sheets)
	config.Subject = newEmail // team@stlpartyhelpers.com
	
	client := config.Client(ctx)
	return gmail.NewService(ctx, option.WithHTTPClient(client))
}

func initSheets(ctx context.Context) (*sheets.Service, error) {
	credsJSON := os.Getenv("GMAIL_CREDENTIALS_JSON")
	if credsJSON == "" {
		return nil, fmt.Errorf("GMAIL_CREDENTIALS_JSON not set")
	}

	var credsData []byte
	if _, err := os.Stat(credsJSON); err == nil {
		credsData, _ = os.ReadFile(credsJSON)
	} else {
		credsData = []byte(credsJSON)
	}

	config, err := google.JWTConfigFromJSON(credsData, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}

	// Use the same subject (email) as Gmail for consistency
	// This ensures the service account impersonates the same user
	config.Subject = newEmail // team@stlpartyhelpers.com
	
	client := config.Client(ctx)
	return sheets.NewService(ctx, option.WithHTTPClient(client))
}

// Lock represents a processing lock
type Lock struct {
	AgentID    string
	ExpiresAt  time.Time
	CreatedAt  time.Time
}

// acquireLock acquires a processing lock with expiration
func acquireLock(ctx context.Context, service *sheets.Service, spreadsheetID, agentID string, idempotent bool) (bool, error) {
	// Create Locks sheet if needed
	_ = createSheetIfNotExists(ctx, service, spreadsheetID, "Locks")
	
	// Set headers
	headers := []string{"Agent ID", "Created At", "Expires At", "Status"}
	vr := &sheets.ValueRange{
		Values: [][]interface{}{convertHeaders(headers)},
	}
	service.Spreadsheets.Values.Update(spreadsheetID, "Locks!A1", vr).ValueInputOption("RAW").Context(ctx).Do()
	
	// Clean up expired locks first
	if err := cleanupExpiredLocks(ctx, service, spreadsheetID); err != nil {
		fmt.Printf("⚠️  Warning: Failed to cleanup expired locks: %v\n", err)
	}
	
	// Check for active locks
	locks, err := getActiveLocks(ctx, service, spreadsheetID)
	if err != nil {
		return false, err
	}
	
	// If idempotent mode, clear all locks and recreate sheets
	if idempotent {
		fmt.Printf("🔄 Idempotent mode: Clearing existing locks and recreating sheets\n")
		clearAllLocks(ctx, service, spreadsheetID)
		// After clearing, treat as no active locks (idempotent override)
		locks = nil
		// Will recreate sheets in initializeSheets
	}
	
	// Check if any active locks exist
	if len(locks) > 0 {
		fmt.Printf("🔒 Active locks found:\n")
		for _, lock := range locks {
			fmt.Printf("  - Agent: %s | Expires: %s\n", lock.AgentID, lock.ExpiresAt.Format(time.RFC3339))
		}
		return false, nil
	}
	
	// Acquire lock (expires in 1 minute)
	expiresAt := time.Now().Add(1 * time.Minute)
	lockRow := []interface{}{
		agentID,
		time.Now().Format(time.RFC3339),
		expiresAt.Format(time.RFC3339),
		"ACTIVE",
	}
	
	vr = &sheets.ValueRange{
		Values: [][]interface{}{lockRow},
	}
	_, err = service.Spreadsheets.Values.Append(spreadsheetID, "Locks", vr).ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
	if err != nil {
		return false, err
	}
	
	fmt.Printf("🔒 Lock acquired by agent: %s (expires: %s)\n", agentID, expiresAt.Format(time.RFC3339))
	return true, nil
}

// releaseLock releases the lock for this agent
func releaseLock(ctx context.Context, service *sheets.Service, spreadsheetID, agentID string) error {
	locks, err := getActiveLocks(ctx, service, spreadsheetID)
	if err != nil {
		return err
	}
	
	// Find and remove this agent's lock
	for i, lock := range locks {
		if lock.AgentID == agentID {
			// Update lock status to RELEASED
			rangeStr := fmt.Sprintf("Locks!D%d", i+2) // +2 for header and 0-index
			vr := &sheets.ValueRange{
				Values: [][]interface{}{{"RELEASED"}},
			}
			_, err := service.Spreadsheets.Values.Update(spreadsheetID, rangeStr, vr).ValueInputOption("RAW").Context(ctx).Do()
			return err
		}
	}
	
	return nil
}

// getActiveLocks returns all active (non-expired) locks
func getActiveLocks(ctx context.Context, service *sheets.Service, spreadsheetID string) ([]Lock, error) {
	resp, err := service.Spreadsheets.Values.Get(spreadsheetID, "Locks!A2:D").Context(ctx).Do()
	if err != nil {
		return []Lock{}, nil // Sheet might not exist yet
	}
	
	var locks []Lock
	now := time.Now()
	
	for _, row := range resp.Values {
		if len(row) < 4 {
			continue
		}
		
		status := fmt.Sprintf("%v", row[3])
		if status != "ACTIVE" {
			continue
		}
		
		expiresStr := fmt.Sprintf("%v", row[2])
		expiresAt, err := time.Parse(time.RFC3339, expiresStr)
		if err != nil {
			continue
		}
		
		// Only return non-expired locks
		if expiresAt.After(now) {
			locks = append(locks, Lock{
				AgentID:   fmt.Sprintf("%v", row[0]),
				CreatedAt: time.Now(), // Could parse from row[1] if needed
				ExpiresAt: expiresAt,
			})
		}
	}
	
	return locks, nil
}

// cleanupExpiredLocks removes expired locks
func cleanupExpiredLocks(ctx context.Context, service *sheets.Service, spreadsheetID string) error {
	resp, err := service.Spreadsheets.Values.Get(spreadsheetID, "Locks!A2:D").Context(ctx).Do()
	if err != nil {
		return nil // Sheet might not exist
	}
	
	now := time.Now()
	rowsToUpdate := []int{}
	
	for i, row := range resp.Values {
		if len(row) < 4 {
			continue
		}
		
		status := fmt.Sprintf("%v", row[3])
		if status != "ACTIVE" {
			continue
		}
		
		expiresStr := fmt.Sprintf("%v", row[2])
		expiresAt, err := time.Parse(time.RFC3339, expiresStr)
		if err != nil {
			continue
		}
		
		// Mark expired locks as EXPIRED
		if expiresAt.Before(now) {
			rowsToUpdate = append(rowsToUpdate, i+2) // +2 for header and 0-index
		}
	}
	
	// Update expired locks
	for _, rowNum := range rowsToUpdate {
		rangeStr := fmt.Sprintf("Locks!D%d", rowNum)
		vr := &sheets.ValueRange{
			Values: [][]interface{}{{"EXPIRED"}},
		}
		service.Spreadsheets.Values.Update(spreadsheetID, rangeStr, vr).ValueInputOption("RAW").Context(ctx).Do()
	}
	
	if len(rowsToUpdate) > 0 {
		fmt.Printf("🧹 Cleaned up %d expired locks\n", len(rowsToUpdate))
	}
	
	return nil
}

// refreshLock extends the lock expiration time
func refreshLock(ctx context.Context, service *sheets.Service, spreadsheetID, agentID string) error {
	resp, err := service.Spreadsheets.Values.Get(spreadsheetID, "Locks!A2:D").Context(ctx).Do()
	if err != nil {
		return err
	}
	
	now := time.Now()
	newExpiresAt := now.Add(1 * time.Minute)
	
	for i, row := range resp.Values {
		if len(row) < 4 {
			continue
		}
		
		rowAgentID := fmt.Sprintf("%v", row[0])
		status := fmt.Sprintf("%v", row[3])
		if rowAgentID == agentID && status == "ACTIVE" {
			// Update expiration time
			rangeStr := fmt.Sprintf("Locks!C%d", i+2) // +2 for header and 0-index
			vr := &sheets.ValueRange{
				Values: [][]interface{}{{newExpiresAt.Format(time.RFC3339)}},
			}
			_, err := service.Spreadsheets.Values.Update(spreadsheetID, rangeStr, vr).ValueInputOption("RAW").Context(ctx).Do()
			return err
		}
	}
	
	return fmt.Errorf("lock not found for agent: %s", agentID)
}

// clearAllLocks clears all locks (for idempotent mode)
func clearAllLocks(ctx context.Context, service *sheets.Service, spreadsheetID string) {
	// Delete all rows except header
	resp, err := service.Spreadsheets.Values.Get(spreadsheetID, "Locks!A2:D").Context(ctx).Do()
	if err != nil || len(resp.Values) == 0 {
		return
	}
	
	// Clear all lock rows
	vr := &sheets.ValueRange{
		Values: [][]interface{}{},
	}
	rangeStr := fmt.Sprintf("Locks!A2:D%d", len(resp.Values)+1)
	service.Spreadsheets.Values.Update(spreadsheetID, rangeStr, vr).ValueInputOption("RAW").Context(ctx).Do()
}

// createSheetIfNotExists creates a sheet if it doesn't exist
func createSheetIfNotExists(ctx context.Context, service *sheets.Service, spreadsheetID, sheetName string) error {
	// Try to get the sheet to see if it exists
	_, err := service.Spreadsheets.Values.Get(spreadsheetID, fmt.Sprintf("%s!A1", sheetName)).Context(ctx).Do()
	if err == nil {
		return nil // Sheet exists
	}
	
	// Sheet doesn't exist, create it
	req := sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				AddSheet: &sheets.AddSheetRequest{
					Properties: &sheets.SheetProperties{Title: sheetName},
				},
			},
		},
	}
	_, err = service.Spreadsheets.BatchUpdate(spreadsheetID, &req).Context(ctx).Do()
	return err
}

func initializeSheets(ctx context.Context, service *sheets.Service, spreadsheetID string, idempotent bool) error {
	// Set headers - sheets will be created automatically when we write
	headers := map[string][]string{
		"Raw Data": {
			// Core email metadata
			"Email ID", "Thread ID", "Date", "From Email", "To Email", "Subject", "Body Preview", "Body Full",
			// Classification & Detection
			"Is Test", "Is Confirmation", "Is Form Submission", "Is Reply",
			"Email Type", "Email Type Confidence", "Source", "Source Page",
			// Client Information
			"Client Email", "Normalized Client Email", "Client Domain", "Client Name",
			// Event Details
			"Event Date", "Event Date Parsed", "Event Start Time", "Event End Time",
			"Event Type", "Hours", "Helpers", "Helper Amount", "Occasion", "Guests", "Event Notes",
			// Pricing Information
			"Total Cost", "Rate", "Deposit", "Payout Amount", "Self Payout Amount",
			// Status & Tracking
			"Status", "Conversation ID", "Message Number", "Forwarded From", "Migration Detected",
			// Job Tracking
			"Job ID", "Job Name",
		},
		"Email Mapping": {
			"Original Email", "Normalized Email", "First Seen", "Last Seen", "Count",
		},
		"State": {
			"LastIndex", "TotalProcessed", "LastRun", "CurrentJobID", "CurrentJobName",
		},
		"Job Stats": {
			"Job ID", "Job Name", "Started At", "Completed At", "Total Processed", "Total Skipped", "Status", "Agent IDs",
		},
	}

	for sheetName, headerRow := range headers {
		// Create sheet if it doesn't exist
		if err := createSheetIfNotExists(ctx, service, spreadsheetID, sheetName); err != nil {
			fmt.Printf("⚠️  Warning: Could not create sheet %s: %v\n", sheetName, err)
		}
		
		// If idempotent mode, clear existing data
		// NOTE: Use Values.Clear (Update with empty values does NOT reliably clear).
		if idempotent && sheetName != "Locks" {
			_, _ = service.Spreadsheets.Values.Clear(spreadsheetID, fmt.Sprintf("%s!A2:Z", sheetName), &sheets.ClearValuesRequest{}).Context(ctx).Do()
			fmt.Printf("🗑️  Cleared existing data in %s sheet (idempotent mode)\n", sheetName)
		}
		
		// Set headers
		vr := &sheets.ValueRange{
			Values: [][]interface{}{convertHeaders(headerRow)},
		}
		_, err := service.Spreadsheets.Values.Update(spreadsheetID, fmt.Sprintf("%s!A1", sheetName), vr).ValueInputOption("RAW").Context(ctx).Do()
		if err != nil {
			fmt.Printf("⚠️  Warning: Could not set headers for %s: %v\n", sheetName, err)
		}
	}

	return nil
}

func convertHeaders(headers []string) []interface{} {
	result := make([]interface{}, len(headers))
	for i, h := range headers {
		result[i] = h
	}
	return result
}

func processEmails(ctx context.Context, gmailService *gmail.Service, sheetsService *sheets.Service,
	spreadsheetID, query string, maxEmails, startIndex, batchSize, delayMs int, verbose bool,
	processedIDsSet map[string]bool, processedIDsList *[]string, agentID, jobID, jobName string, workers int) (processed, skipped int, err error) {

	userEmail := newEmail // Use new email as user
	startTime := time.Now()
	lastStateSave := time.Now()
	lastLockRefresh := time.Now()

	fmt.Printf("🚀 Starting email processing...\n")
	fmt.Printf("  Job: %s - %s\n", jobID, jobName)
	fmt.Printf("  Agent: %s\n", agentID)
	fmt.Printf("  Workers: %d\n\n", workers)
	
	// Initialize job stats
	updateJobStats(ctx, sheetsService, spreadsheetID, jobID, jobName, agentID, 0, 0, false)

	// Shared state protected by mutex
	var mu sync.Mutex
	var processedCount int
	var skippedCount int
	var batch [][]interface{}
	emailMapping := make(map[string]string) // original -> normalized
	mappingAgg := make(map[string]*mappingAggregate) // derived mapping rows (client->normalized)

	// Channel for distributing message IDs to workers
	messageChan := make(chan *gmail.Message, workers*2) // Buffered channel
	doneChan := make(chan bool, workers)
	
	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for msg := range messageChan {
				// Check if already processed (thread-safe)
				mu.Lock()
				alreadyProcessed := processedIDsSet[msg.Id]
				mu.Unlock()
				
				if alreadyProcessed {
					if verbose {
						fmt.Printf("  [Worker %d] ⏭️  Skipping already processed: %s\n", workerID, msg.Id)
					}
					mu.Lock()
					skippedCount++
					mu.Unlock()
					continue
				}

				// Process message
				emailData, err := processMessage(ctx, gmailService, userEmail, msg.Id, emailMapping, verbose && workerID == 0)
				if err != nil {
					if verbose {
						fmt.Printf("  [Worker %d] ⚠️  Error processing %s: %v\n", workerID, msg.Id, err)
					}
					mu.Lock()
					skippedCount++
					mu.Unlock()
					continue
				}

				if emailData == nil || emailData.IsTest {
					mu.Lock()
					skippedCount++
					mu.Unlock()
					continue
				}

				// Stamp with job ID
				emailData.JobID = jobID
				emailData.JobName = jobName

				// Thread-safe batch accumulation
				mu.Lock()
				// Track mapping (client -> normalized) when normalization changes the email
				if emailData.ClientEmail != "" && emailData.NormalizedClientEmail != "" && emailData.ClientEmail != emailData.NormalizedClientEmail {
					key := emailData.ClientEmail + "->" + emailData.NormalizedClientEmail
					agg := mappingAgg[key]
					if agg == nil {
						agg = &mappingAggregate{
							Original:   emailData.ClientEmail,
							Normalized: emailData.NormalizedClientEmail,
							FirstSeen:  emailData.Date,
							LastSeen:   emailData.Date,
							Count:      0,
						}
						mappingAgg[key] = agg
					}
					agg.Count++
					if emailData.Date.Before(agg.FirstSeen) {
						agg.FirstSeen = emailData.Date
					}
					if emailData.Date.After(agg.LastSeen) {
						agg.LastSeen = emailData.Date
					}
				}
				batch = append(batch, emailData.ToRow())
				*processedIDsList = append(*processedIDsList, msg.Id)
				processedIDsSet[msg.Id] = true
				processedCount++
				currentProcessed := processedCount
				currentSkipped := skippedCount
				shouldWriteBatch := len(batch) >= 25
				mu.Unlock()

				if verbose && workerID == 0 {
					fmt.Printf("  [Worker %d] ✅ Processed: %s | Thread: %s\n", workerID, msg.Id, msg.ThreadId)
				}

				// Write batch if needed (only one worker writes at a time)
				if shouldWriteBatch {
					mu.Lock()
					// Double-check batch size (another worker might have written)
						if len(batch) >= 25 {
						batchToWrite := make([][]interface{}, len(batch))
						copy(batchToWrite, batch)
						batch = [][]interface{}{} // Clear batch
						mu.Unlock()

						// Write batch with verification
						fmt.Printf("  💾 Writing batch of %d emails to spreadsheet...\n", len(batchToWrite))
						if err := writeBatch(ctx, sheetsService, spreadsheetID, batchToWrite); err != nil {
							fmt.Printf("  ❌ ERROR writing batch: %v\n", err)
							fmt.Printf("  🔍 Batch details: %d rows, first message ID: %v\n", len(batchToWrite), func() interface{} {
								if len(batchToWrite) > 0 && len(batchToWrite[0]) > 0 {
									return batchToWrite[0][0]
								}
								return "N/A"
							}())
							mu.Lock()
							batch = batchToWrite // Restore batch on error
							mu.Unlock()
							continue
						}

						fmt.Printf("  ✅ Successfully wrote batch of %d emails to spreadsheet\n", len(batchToWrite))
						// Lightweight progress reporting to Sheets (safe frequency).
						// We avoid writing State/Email Mapping here to stay under Sheets write quotas.
						updateJobStats(ctx, sheetsService, spreadsheetID, jobID, jobName, agentID, currentProcessed, currentSkipped, false)
					} else {
						mu.Unlock()
					}
				}

				// Progress update (every 100 processed)
				if currentProcessed%100 == 0 {
					elapsed := time.Since(startTime)
					rate := float64(currentProcessed) / elapsed.Seconds()
					remaining := ""
					if maxEmails > 0 {
						remaining = fmt.Sprintf(" | Remaining: %d", maxEmails-currentProcessed)
					}
					fmt.Printf("  📊 Progress: %d processed | %d skipped | %.1f emails/sec%s\n",
						currentProcessed, currentSkipped, rate, remaining)
				}
			}
			doneChan <- true
		}(i)
	}

	// Producer goroutine: fetch batches and send to workers
	go func() {
		defer close(messageChan)
		var pageToken string
		batchNum := 0

		for {
			// Check if we've hit the limit
			mu.Lock()
			currentProcessed := processedCount
			mu.Unlock()

			if maxEmails > 0 && currentProcessed >= maxEmails {
				fmt.Printf("\n✅ Reached max emails limit (%d)\n", maxEmails)
				break
			}

			if verbose {
				fmt.Printf("📦 Fetching batch %d...\n", batchNum+1)
			} else if batchNum%10 == 0 {
				mu.Lock()
				p := processedCount
				s := skippedCount
				mu.Unlock()
				fmt.Printf("📦 Batch %d | Processed: %d | Skipped: %d | Elapsed: %s\n",
					batchNum+1, p, s, time.Since(startTime).Round(time.Second))
			}

			// Search Gmail with pagination
			call := gmailService.Users.Messages.List(userEmail).
				Q(query).
				MaxResults(int64(batchSize))

			if pageToken != "" {
				call = call.PageToken(pageToken)
				if verbose {
					fmt.Printf("  📄 Using page token (batch %d, continuing pagination)\n", batchNum+1)
				}
			}

			resp, err := call.Context(ctx).Do()
			if err != nil {
				fmt.Printf("⚠️  Error fetching batch %d: %v\n", batchNum+1, err)
				// Log detailed error for debugging
				if verbose {
					fmt.Printf("  🔍 Query: %s\n", query)
					fmt.Printf("  🔍 Page token present: %v\n", pageToken != "")
				}
				break
			}

			if len(resp.Messages) == 0 {
				fmt.Printf("\n✅ No more emails found (batch %d returned 0 messages)\n", batchNum+1)
				if verbose {
					fmt.Printf("  🔍 Query: %s\n", query)
					fmt.Printf("  🔍 Total batches fetched: %d\n", batchNum+1)
				}
				break
			}

			// Log batch details
			if verbose || batchNum%50 == 0 {
				fmt.Printf("  📦 Batch %d: Fetched %d message IDs | Has next page: %v\n", 
					batchNum+1, len(resp.Messages), resp.NextPageToken != "")
			}

			// Send message IDs to workers
			for _, msg := range resp.Messages {
				messageChan <- msg
			}

			// Check if there are more pages
			pageToken = resp.NextPageToken
			if pageToken == "" {
				fmt.Printf("\n✅ Reached end of pagination (batch %d, no next page token)\n", batchNum+1)
				if verbose {
					fmt.Printf("  🔍 Total batches: %d\n", batchNum+1)
					fmt.Printf("  🔍 Total messages fetched: ~%d\n", (batchNum+1)*batchSize)
				}
				break
			}

			// Refresh lock every 30 seconds
			if time.Since(lastLockRefresh) > 30*time.Second {
				if err := refreshLock(ctx, sheetsService, spreadsheetID, agentID); err != nil {
					if verbose {
						fmt.Printf("  ⚠️  Warning: Failed to refresh lock: %v\n", err)
					}
				} else if verbose {
					fmt.Printf("  🔒 Lock refreshed (agent: %s)\n", agentID)
				}
				lastLockRefresh = time.Now()
			}

			// Auto-save state every 2 minutes
			if time.Since(lastStateSave) > 2*time.Minute {
				mu.Lock()
				state := &State{
					LastIndex:      startIndex + processedCount + skippedCount,
					TotalProcessed: processedCount,
					LastRun:        time.Now(),
					ProcessedIDs:   *processedIDsList,
					CurrentJobID:   jobID,
					CurrentJobName: jobName,
				}
				processedIDsCopy := make([]string, len(*processedIDsList))
				copy(processedIDsCopy, *processedIDsList)
				p := processedCount
				s := skippedCount
				mu.Unlock()

				if err := saveState(ctx, sheetsService, spreadsheetID, state); err == nil {
					lastStateSave = time.Now()
					if verbose {
						fmt.Printf("  💾 Auto-saved state (%d message IDs tracked, Job: %s)\n", len(processedIDsCopy), jobID)
					} else {
						fmt.Printf("  💾 State saved (%d message IDs, Job: %s)\n", len(processedIDsCopy), jobID)
					}
				}

				// Update job stats
				updateJobStats(ctx, sheetsService, spreadsheetID, jobID, jobName, agentID, p, s, false)
			}

			// Rate limiting delay
			if delayMs > 0 {
				time.Sleep(time.Duration(delayMs) * time.Millisecond)
			}

			batchNum++
		}
	}()

	// Wait for all workers to finish
	wg.Wait()

	// Write remaining batch
	mu.Lock()
	finalBatch := batch
	finalProcessed := processedCount
	finalSkipped := skippedCount
	mu.Unlock()

	if len(finalBatch) > 0 {
		fmt.Printf("  💾 Writing final batch of %d emails to spreadsheet...\n", len(finalBatch))
		if err := writeBatch(ctx, sheetsService, spreadsheetID, finalBatch); err != nil {
			return finalProcessed, finalSkipped, fmt.Errorf("failed to write final batch: %w", err)
		}
		fmt.Printf("  ✅ Successfully wrote final batch of %d emails to spreadsheet\n", len(finalBatch))

		// Save state after final batch
		mu.Lock()
		state := &State{
			LastIndex:      startIndex + finalProcessed + finalSkipped,
			TotalProcessed: finalProcessed,
			LastRun:        time.Now(),
			ProcessedIDs:   *processedIDsList,
			CurrentJobID:   jobID,
			CurrentJobName: jobName,
		}
		mu.Unlock()
		saveState(ctx, sheetsService, spreadsheetID, state)
	}

	// Final mapping snapshot (derived from processed rows)
	mu.Lock()
	finalMappingSnapshot := make([]*mappingAggregate, 0, len(mappingAgg))
	for _, v := range mappingAgg {
		c := *v
		finalMappingSnapshot = append(finalMappingSnapshot, &c)
	}
	mu.Unlock()
	writeEmailMappingSnapshot(ctx, sheetsService, spreadsheetID, finalMappingSnapshot)

	return finalProcessed, finalSkipped, nil
}

type mappingAggregate struct {
	Original   string
	Normalized string
	FirstSeen  time.Time
	LastSeen   time.Time
	Count      int
}

type EmailData struct {
	// Core email metadata
	EmailID            string
	ThreadID          string
	Date              time.Time
	FromEmail         string
	ToEmail           string
	Subject           string
	BodyPreview       string
	BodyFull          string // Full body for later analysis
	
	// Classification & Detection
	IsTest            bool
	IsConfirmation    bool
	IsFormSubmission  bool   // Detected as form submission (Zapier, Google Forms, etc.)
	IsReply           bool   // Is this a reply to a previous email
	EmailType         string // STLPH, Marketing, School, Other
	EmailTypeConfidence float64 // 0.0-1.0 confidence in classification
	Source            string // Form source: "Zapier", "Google Forms", "Website", "Direct", "Unknown"
	SourcePage        string // Which page/form: extracted from subject or body
	
	// Client Information
	ClientEmail       string
	NormalizedClientEmail string
	ClientDomain      string // Extracted domain for organization grouping
	ClientName        string // Extracted client name if available
	
	// Event Details
	EventDate         string // Raw extracted date string
	EventDateParsed   time.Time // Parsed event date (if successful)
	EventStartTime    string // Start time (e.g., "2:00 PM")
	EventEndTime      string // End time (e.g., "6:00 PM")
	EventType         string // Birthday, Corporate, Wedding, etc.
	Hours             string // Duration in hours
	Helpers           string // Number of helpers
	HelperAmount      float64 // Helper amount (if specified)
	Occasion          string // Event occasion/type
	Guests            string // Number of guests
	EventNotes        string // Specific notes about the event
	
	// Pricing Information
	TotalCost         float64
	Rate              float64 // Rate per hour/helper
	Deposit           float64
	PayoutAmount      float64 // Calculated: TotalCost * 0.45
	SelfPayoutAmount  float64 // Calculated: PayoutAmount * 0.10
	
	// Status & Tracking
	Status            string // Test, Confirmed, Pending, etc.
	ConversationID    string
	MessageNumber     int
	ForwardedFrom     string
	MigrationDetected bool
	
	// Job Tracking
	JobID             string // Job ID that processed this email
	JobName           string // Job name/description
}

func processMessage(ctx context.Context, service *gmail.Service, userEmail, messageID string,
	emailMapping map[string]string, verbose bool) (*EmailData, error) {

	msg, err := service.Users.Messages.Get(userEmail, messageID).Format("full").Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	// Extract headers
	var from, to, subject, dateStr string
	for _, header := range msg.Payload.Headers {
		switch header.Name {
		case "From":
			from = header.Value
		case "To":
			to = header.Value
		case "Subject":
			subject = header.Value
		case "Date":
			dateStr = header.Value
		}
	}
	
	// Use thread ID from Gmail API - this is the key for conversation grouping
	// Thread ID groups all emails in the same conversation (replies, forwards, etc.)
	threadID := msg.ThreadId

	// Normalize emails (migration handling)
	normalizedFrom := normalizeEmail(from)
	normalizedTo := normalizeEmail(to)
	
	// Detect forwarding/migration
	forwardedFrom := ""
	migrationDetected := false
	if strings.Contains(from, oldEmail) || strings.Contains(to, oldEmail) {
		migrationDetected = true
		if strings.Contains(from, oldEmail) {
			forwardedFrom = oldEmail
		}
	}

	// Track email mapping
	if normalizedFrom != from {
		emailMapping[from] = normalizedFrom
	}
	if normalizedTo != to && normalizedTo != "" {
		emailMapping[to] = normalizedTo
	}

	// Parse date
	date, _ := time.Parse(time.RFC1123Z, dateStr)
	if date.IsZero() {
		date = time.Unix(msg.InternalDate/1000, 0)
	}

	// Extract body (full and preview)
	body := extractBody(msg.Payload)
	bodyPreview := body
	if len(bodyPreview) > 500 {
		bodyPreview = bodyPreview[:500] + "..."
	}
	// Store full body (truncate if too long for Sheets - max ~50k chars per cell)
	bodyFull := body
	if len(bodyFull) > 10000 {
		bodyFull = bodyFull[:10000] + "... [truncated]"
	}

	// Extract data
	combined := strings.ToLower(subject + " " + body)
	isTest := detectTest(combined)
	isConfirmation := detectConfirmation(combined)
	isReply := detectReply(subject)
	
	// Form submission detection
	isFormSubmission, source, sourcePage := detectFormSubmission(from, subject, body)
	
	// Client information
	clientEmail := extractClientEmail(normalizedFrom, normalizedTo, body)
	normalizedClientEmail := normalizeEmail(clientEmail)
	clientDomain := extractDomain(normalizedClientEmail)
	clientName := extractClientName(from, body)
	
	// Event and pricing data
	eventData := extractEventData(body, subject)
	pricingData := extractPricingData(body, subject)
	
	// Email classification with confidence
	emailType, emailTypeConfidence := classifyEmailType(subject, body, normalizedFrom)

	// Use thread ID + normalized client email for conversation grouping
	// Gmail thread ID groups all replies/forwards in same conversation
	// This ensures we group conversations even when email addresses change
	conversationID := threadID + "_" + normalizedClientEmail
	
	// Detect conversation continuation (forwarding scenario)
	// If client wrote to old email, but this thread continues from new email
	if migrationDetected && normalizedClientEmail != "" {
		// This is a forwarded/migrated conversation
		// Thread ID will help us group it with related emails
		if verbose {
			fmt.Printf("  🔄 Migration detected: %s -> %s (Thread: %s, Conversation: %s)\n", 
				oldEmail, newEmail, threadID, conversationID)
		}
	}
	
	return &EmailData{
		// Core email metadata
		EmailID:              messageID,
		ThreadID:            msg.ThreadId,
		Date:                date,
		FromEmail:           from,
		ToEmail:             to,
		Subject:             subject,
		BodyPreview:         bodyPreview,
		BodyFull:            bodyFull,
		
		// Classification & Detection
		IsTest:              isTest,
		IsConfirmation:      isConfirmation,
		IsFormSubmission:    isFormSubmission,
		IsReply:             isReply,
		EmailType:           emailType,
		EmailTypeConfidence: emailTypeConfidence,
		Source:              source,
		SourcePage:          sourcePage,
		
		// Client Information
		ClientEmail:         clientEmail,
		NormalizedClientEmail: normalizedClientEmail,
		ClientDomain:        clientDomain,
		ClientName:           clientName,
		
		// Event Details
		EventDate:           eventData.EventDate,
		EventDateParsed:     eventData.EventDateParsed,
		EventStartTime:      eventData.EventStartTime,
		EventEndTime:        eventData.EventEndTime,
		EventType:           eventData.EventType,
		Hours:               eventData.Hours,
		Helpers:             eventData.Helpers,
		HelperAmount:        eventData.HelperAmount,
		Occasion:            eventData.Occasion,
		Guests:              eventData.Guests,
		EventNotes:          eventData.EventNotes,
		
		// Pricing Information
		TotalCost:           pricingData.TotalCost,
		Rate:                pricingData.Rate,
		Deposit:             pricingData.Deposit,
		PayoutAmount:        pricingData.PayoutAmount,
		SelfPayoutAmount:    pricingData.SelfPayoutAmount,
		
		// Status & Tracking
		Status:              getStatus(isConfirmation, isTest),
		ConversationID:      conversationID,
		MessageNumber:       1,
		ForwardedFrom:       forwardedFrom,
		MigrationDetected:   migrationDetected,
		
		// Job Tracking
		JobID:               "", // Will be set in processEmails
		JobName:             "", // Will be set in processEmails
	}, nil
}

func normalizeEmail(email string) string {
	if email == "" {
		return ""
	}
	
	// Extract email from "Name <email@domain.com>" format
	emailRegex := regexp.MustCompile(`([a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\.[a-zA-Z0-9_-]+)`)
	matches := emailRegex.FindStringSubmatch(email)
	if len(matches) > 0 {
		email = matches[1]
	}

	// Normalize: stlpartyhelpers@gmail.com -> team@stlpartyhelpers.com
	if strings.Contains(email, oldEmail) {
		email = strings.ReplaceAll(email, oldEmail, newEmail)
	}

	return strings.ToLower(strings.TrimSpace(email))
}

func extractBody(payload *gmail.MessagePart) string {
	if payload == nil {
		return ""
	}

	if payload.MimeType == "text/plain" && payload.Body != nil && payload.Body.Data != "" {
		data, _ := base64.URLEncoding.DecodeString(payload.Body.Data)
		return string(data)
	}

	if payload.MimeType == "text/html" && payload.Body != nil && payload.Body.Data != "" {
		data, _ := base64.URLEncoding.DecodeString(payload.Body.Data)
		return stripHTML(string(data))
	}

	for _, part := range payload.Parts {
		if body := extractBody(part); body != "" {
			return body
		}
	}

	return ""
}

func stripHTML(html string) string {
	re := regexp.MustCompile(`<[^>]+>`)
	return re.ReplaceAllString(html, " ")
}

func detectTest(text string) bool {
	keywords := []string{"test", "testing", "demo", "sample", "fake"}
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}

func detectConfirmation(text string) bool {
	keywords := []string{"confirm", "confirmed", "confirmation", "accepted", "approved"}
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}

func extractClientEmail(from, to, body string) string {
	emailRegex := regexp.MustCompile(`([a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\.[a-zA-Z0-9_-]+)`)
	
	// Try body first
	matches := emailRegex.FindAllString(body, -1)
	for _, match := range matches {
		match = strings.ToLower(match)
		if !strings.Contains(match, "zapier") &&
			!strings.Contains(match, "noreply") &&
			!strings.Contains(match, "no-reply") &&
			!strings.Contains(match, "stlpartyhelpers") &&
			!strings.Contains(match, "team@") {
			return match
		}
	}

	// Try from/to (but not our emails)
	fromEmail := normalizeEmail(from)
	toEmail := normalizeEmail(to)
	
	if fromEmail != "" && !strings.Contains(fromEmail, "stlpartyhelpers") && !strings.Contains(fromEmail, "team@") {
		return fromEmail
	}
	if toEmail != "" && !strings.Contains(toEmail, "stlpartyhelpers") && !strings.Contains(toEmail, "team@") {
		return toEmail
	}

	return ""
}

type EventData struct {
	EventDate      string
	EventDateParsed time.Time
	EventStartTime string
	EventEndTime   string
	EventType      string
	Hours          string
	Helpers        string
	HelperAmount   float64
	Occasion       string
	Guests         string
	EventNotes     string
}

func extractEventData(body, subject string) EventData {
	combined := body + " " + subject
	lower := strings.ToLower(combined)
	
	// Extract event date (multiple patterns)
	datePatterns := []string{
		`(\d{1,2}[/-]\d{1,2}[/-]\d{2,4})`,
		`(January|February|March|April|May|June|July|August|September|October|November|December)\s+\d{1,2}(?:st|nd|rd|th)?,?\s+\d{4}`,
		`(January|February|March|April|May|June|July|August|September|October|November|December)\s+\d{1,2}`,
		`\b(Mon|Tue|Wed|Thu|Fri|Sat|Sun)[a-z]*\s+\d{1,2}[/-]\d{1,2}[/-]\d{2,4}\b`,
	}
	
	var eventDate string
	var eventDateParsed time.Time
	for _, pattern := range datePatterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		matches := re.FindStringSubmatch(combined)
		if len(matches) > 0 {
			eventDate = matches[0]
			// Try to parse it
			if parsed, err := parseEventDate(eventDate); err == nil {
				eventDateParsed = parsed
			}
			break
		}
	}
	
	// Extract start/end time
	timePattern := regexp.MustCompile(`(?i)(\d{1,2}):?(\d{2})?\s*(am|pm|AM|PM)`)
	timeMatches := timePattern.FindAllString(combined, 2)
	var startTime, endTime string
	if len(timeMatches) > 0 {
		startTime = timeMatches[0]
	}
	if len(timeMatches) > 1 {
		endTime = timeMatches[1]
	}
	
	// Extract event type
	eventType := ""
	eventTypePatterns := map[string]string{
		"birthday":     "Birthday",
		"corporate":    "Corporate",
		"wedding":      "Wedding",
		"anniversary":  "Anniversary",
		"graduation":   "Graduation",
		"holiday":      "Holiday",
		"fundraiser":   "Fundraiser",
		"school":       "School",
		"party":        "Party",
	}
	for keyword, etype := range eventTypePatterns {
		if strings.Contains(lower, keyword) {
			eventType = etype
			break
		}
	}
	
	// Extract hours/duration
	var hours string
	hoursPattern := regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)\s*(?:hours?|hrs?|h)`)
	if matches := hoursPattern.FindStringSubmatch(combined); len(matches) > 0 {
		hours = matches[1]
	}
	
	// Extract helpers count
	helpersPattern := regexp.MustCompile(`(?i)(\d+)\s*(?:helpers?|staff|people|workers?)`)
	var helpers string
	var helperAmount float64
	if matches := helpersPattern.FindStringSubmatch(combined); len(matches) > 0 {
		helpers = matches[1]
		// Try to extract helper amount (e.g., "$50 per helper")
		helperAmountPattern := regexp.MustCompile(`(?i)(?:per helper|helper rate|each helper).*?\$?(\d+(?:\.\d+)?)`)
		if amtMatches := helperAmountPattern.FindStringSubmatch(combined); len(amtMatches) > 0 {
			if amt, err := parseDollarAmount(amtMatches[1]); err == nil {
				helperAmount = amt
			}
		}
	}
	
	// Extract guests
	guestsPattern := regexp.MustCompile(`(?i)(\d+)\s*(?:guests?|people|attendees?)`)
	var guests string
	if matches := guestsPattern.FindStringSubmatch(combined); len(matches) > 0 {
		guests = matches[1]
	}
	
	// Extract occasion/notes (look for common phrases)
	var eventNotes string
	notesPatterns := []string{
		`(?i)notes?[:\s]+([^\n]{10,200})`,
		`(?i)special\s+requests?[:\s]+([^\n]{10,200})`,
		`(?i)additional\s+info[:\s]+([^\n]{10,200})`,
	}
	for _, pattern := range notesPatterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(combined); len(matches) > 1 {
			eventNotes = strings.TrimSpace(matches[1])
			if len(eventNotes) > 500 {
				eventNotes = eventNotes[:500] + "..."
			}
			break
		}
	}
	
	return EventData{
		EventDate:      eventDate,
		EventDateParsed: eventDateParsed,
		EventStartTime: startTime,
		EventEndTime:   endTime,
		EventType:      eventType,
		Hours:          hours,
		Helpers:        helpers,
		HelperAmount:   helperAmount,
		Occasion:       eventType, // Use event type as occasion if no specific occasion found
		Guests:         guests,
		EventNotes:     eventNotes,
	}
}

func parseEventDate(dateStr string) (time.Time, error) {
	formats := []string{
		"1/2/2006",
		"01/02/2006",
		"1-2-2006",
		"01-02-2006",
		"January 2, 2006",
		"January 2nd, 2006",
		"Jan 2, 2006",
		"1/2/06",
		"01/02/06",
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	
	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

type PricingData struct {
	TotalCost        float64
	Rate             float64
	Deposit          float64
	PayoutAmount     float64 // Calculated: TotalCost * 0.45
	SelfPayoutAmount float64 // Calculated: PayoutAmount * 0.10
}

func extractPricingData(body, subject string) PricingData {
	combined := body + " " + subject
	
	// More intelligent pricing extraction with context
	var totalCost, rate, deposit float64
	
	// Extract total cost (look for "total", "amount", "cost", "$X")
	totalPatterns := []string{
		`(?i)total[:\s]+\$?(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)`,
		`(?i)amount[:\s]+\$?(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)`,
		`(?i)cost[:\s]+\$?(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)`,
		`(?i)\$(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)\s*(?:total|amount)`,
	}
	for _, pattern := range totalPatterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(combined); len(matches) > 1 {
			if val, err := parseDollarAmount(matches[1]); err == nil {
				totalCost = val
				break
			}
		}
	}
	
	// Extract rate (look for "rate", "per hour", "$X/hour")
	ratePatterns := []string{
		`(?i)rate[:\s]+\$?(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)`,
		`(?i)\$(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)\s*(?:per hour|/hour|per hr)`,
		`(?i)hourly[:\s]+\$?(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)`,
	}
	for _, pattern := range ratePatterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(combined); len(matches) > 1 {
			if val, err := parseDollarAmount(matches[1]); err == nil {
				rate = val
				break
			}
		}
	}
	
	// Extract deposit (look for "deposit", "down payment")
	depositPatterns := []string{
		`(?i)deposit[:\s]+\$?(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)`,
		`(?i)down payment[:\s]+\$?(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)`,
	}
	for _, pattern := range depositPatterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(combined); len(matches) > 1 {
			if val, err := parseDollarAmount(matches[1]); err == nil {
				deposit = val
				break
			}
		}
	}
	
	// Fallback: if no context found, try simple dollar amounts (largest is usually total)
	if totalCost == 0 {
		dollarRegex := regexp.MustCompile(`\$?(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)`)
		matches := dollarRegex.FindAllString(combined, -1)
		if len(matches) > 0 {
			// Find the largest amount (likely total)
			var amounts []float64
			for _, match := range matches {
				if val, err := parseDollarAmount(match); err == nil && val > 0 {
					amounts = append(amounts, val)
				}
			}
			if len(amounts) > 0 {
				sort.Float64s(amounts)
				totalCost = amounts[len(amounts)-1]
				if len(amounts) > 1 {
					deposit = amounts[0] // Smallest might be deposit
				}
			}
		}
	}
	
	// Calculate payouts (45% of total, 10% of payout to self)
	payoutAmount := totalCost * 0.45
	selfPayoutAmount := payoutAmount * 0.10
	
	return PricingData{
		TotalCost:        totalCost,
		Rate:             rate,
		Deposit:          deposit,
		PayoutAmount:     payoutAmount,
		SelfPayoutAmount: selfPayoutAmount,
	}
}

func parseDollarAmount(str string) (float64, error) {
	cleaned := strings.ReplaceAll(strings.ReplaceAll(str, "$", ""), ",", "")
	return strconv.ParseFloat(cleaned, 64)
}

// EmailClassificationResult contains classification and confidence
type EmailClassificationResult struct {
	Type       string
	Confidence float64
}

func classifyEmailType(subject, body, from string) (string, float64) {
	combined := strings.ToLower(subject + " " + body)
	fromLower := strings.ToLower(from)
	
	// STLPH keywords (business-related)
	stlphKeywords := []string{
		"party helpers", "stl party", "party helpers stl", "stlpartyhelpers",
		"booking", "quote", "estimate", "deposit", "event helpers",
		"party planning", "event staff", "helpers needed",
	}
	
	// Marketing keywords
	marketingKeywords := []string{
		"newsletter", "promotion", "sale", "discount", "offer",
		"unsubscribe", "marketing", "advertisement", "sponsor",
	}
	
	// School keywords
	schoolKeywords := []string{
		"school", "pta", "pta meeting", "parent teacher",
		"school event", "fundraiser", "school fundraiser",
		"grade", "teacher", "principal", "school board",
	}
	
	// Score each category
	stlphScore := 0.0
	marketingScore := 0.0
	schoolScore := 0.0
	
	// Check STLPH
	for _, kw := range stlphKeywords {
		if strings.Contains(combined, kw) {
			stlphScore += 1.0
		}
	}
	// Bonus for STLPH domain
	if strings.Contains(fromLower, "stlpartyhelpers") || strings.Contains(fromLower, "team@stlpartyhelpers") {
		stlphScore += 2.0
	}
	
	// Check Marketing
	for _, kw := range marketingKeywords {
		if strings.Contains(combined, kw) {
			marketingScore += 1.0
		}
	}
	
	// Check School
	for _, kw := range schoolKeywords {
		if strings.Contains(combined, kw) {
			schoolScore += 1.0
		}
	}
	
	// Determine winner
	maxScore := stlphScore
	resultType := "STLPH"
	if marketingScore > maxScore {
		maxScore = marketingScore
		resultType = "Marketing"
	}
	if schoolScore > maxScore {
		maxScore = schoolScore
		resultType = "School"
	}
	
	// Calculate confidence (normalize to 0.0-1.0)
	totalScore := stlphScore + marketingScore + schoolScore
	confidence := 0.5 // Default
	if totalScore > 0 {
		confidence = maxScore / totalScore
		if confidence > 1.0 {
			confidence = 1.0
		}
	}
	
	// If no strong match, default to "Other"
	if maxScore < 1.0 {
		resultType = "Other"
		confidence = 0.3
	}
	
	return resultType, confidence
}

func detectFormSubmission(from, subject, body string) (bool, string, string) {
	fromLower := strings.ToLower(from)
	subjectLower := strings.ToLower(subject)
	bodyLower := strings.ToLower(body)
	combined := subjectLower + " " + bodyLower
	
	// Zapier indicators
	if strings.Contains(fromLower, "zapier") || strings.Contains(combined, "zapier") {
		// Try to extract form/page name from subject
		page := extractPageFromSubject(subject)
		return true, "Zapier", page
	}
	
	// Google Forms indicators
	if strings.Contains(fromLower, "forms.gle") || strings.Contains(combined, "google forms") {
		return true, "Google Forms", extractPageFromSubject(subject)
	}
	
	// Form submission keywords
	formKeywords := []string{"form submission", "new lead", "new inquiry", "contact form"}
	for _, kw := range formKeywords {
		if strings.Contains(combined, kw) {
			return true, "Website", extractPageFromSubject(subject)
		}
	}
	
	return false, "Unknown", ""
}

func extractPageFromSubject(subject string) string {
	// Try to extract page/form name from subject
	// Common patterns: "New Lead from [Page]", "Form: [Page]", "[Page] - Contact Form"
	patterns := []string{
		`(?i)from\s+([^-:\n]+)`,
		`(?i)form[:\s]+([^-:\n]+)`,
		`(?i)^([^-:\n]+)\s*[-:]`,
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(subject); len(matches) > 1 {
			page := strings.TrimSpace(matches[1])
			if len(page) > 0 && len(page) < 100 {
				return page
			}
		}
	}
	
	return "Unknown"
}

func detectReply(subject string) bool {
	subjectLower := strings.ToLower(subject)
	return strings.HasPrefix(subjectLower, "re:") || strings.HasPrefix(subjectLower, "fwd:") || strings.HasPrefix(subjectLower, "fw:")
}

func extractClientName(from, body string) string {
	// Try to extract name from "Name <email>" format
	namePattern := regexp.MustCompile(`^([^<]+)\s*<`)
	if matches := namePattern.FindStringSubmatch(from); len(matches) > 1 {
		name := strings.TrimSpace(matches[1])
		// Remove quotes
		name = strings.Trim(name, `"'`)
		if len(name) > 0 && len(name) < 100 {
			return name
		}
	}
	
	// Try to extract from body (e.g., "Hi, [Name]")
	greetingPattern := regexp.MustCompile(`(?i)(?:hi|hello|dear)\s+([A-Z][a-z]+(?:\s+[A-Z][a-z]+)?)`)
	if matches := greetingPattern.FindStringSubmatch(body); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	
	return ""
}

func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return strings.ToLower(parts[1])
	}
	return ""
}

func getStatus(isConfirmation, isTest bool) string {
	if isTest {
		return "Test"
	}
	if isConfirmation {
		return "Confirmed"
	}
	return "Pending"
}

func (e *EmailData) ToRow() []interface{} {
	eventDateParsed := ""
	if !e.EventDateParsed.IsZero() {
		eventDateParsed = e.EventDateParsed.Format(time.RFC3339)
	}
	
	return []interface{}{
		// Core email metadata
		e.EmailID,
		e.ThreadID,
		e.Date.Format(time.RFC3339),
		e.FromEmail,
		e.ToEmail,
		e.Subject,
		e.BodyPreview,
		e.BodyFull, // Full body for later analysis
		
		// Classification & Detection
		strconv.FormatBool(e.IsTest),
		strconv.FormatBool(e.IsConfirmation),
		strconv.FormatBool(e.IsFormSubmission),
		strconv.FormatBool(e.IsReply),
		e.EmailType,
		fmt.Sprintf("%.2f", e.EmailTypeConfidence),
		e.Source,
		e.SourcePage,
		
		// Client Information
		e.ClientEmail,
		e.NormalizedClientEmail,
		e.ClientDomain,
		e.ClientName,
		
		// Event Details
		e.EventDate,
		eventDateParsed,
		e.EventStartTime,
		e.EventEndTime,
		e.EventType,
		e.Hours,
		e.Helpers,
		fmt.Sprintf("%.2f", e.HelperAmount),
		e.Occasion,
		e.Guests,
		e.EventNotes,
		
		// Pricing Information
		fmt.Sprintf("%.2f", e.TotalCost),
		fmt.Sprintf("%.2f", e.Rate),
		fmt.Sprintf("%.2f", e.Deposit),
		fmt.Sprintf("%.2f", e.PayoutAmount),
		fmt.Sprintf("%.2f", e.SelfPayoutAmount),
		
		// Status & Tracking
		e.Status,
		e.ConversationID,
		e.MessageNumber,
		e.ForwardedFrom,
		strconv.FormatBool(e.MigrationDetected),
		
		// Job Tracking
		e.JobID,
		e.JobName,
	}
}

func writeBatch(ctx context.Context, service *sheets.Service, spreadsheetID string, batch [][]interface{}) error {
	if len(batch) == 0 {
		return fmt.Errorf("cannot write empty batch")
	}
	
	vr := &sheets.ValueRange{
		Values: batch,
	}
	
	// Write batch
	resp, err := service.Spreadsheets.Values.Append(spreadsheetID, "Raw Data", vr).ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to append batch: %w", err)
	}
	
	// Verify write succeeded - check that we got an update response
	if resp.Updates == nil {
		return fmt.Errorf("write succeeded but no update response received")
	}
	
	// Log verification details
	rowsUpdated := resp.Updates.UpdatedRows
	if rowsUpdated != int64(len(batch)) {
		fmt.Printf("  ⚠️  Warning: Expected %d rows, but %d rows were updated\n", len(batch), rowsUpdated)
	}
	
	// Additional verification: read back the last row to confirm it was written
	// Get the first message ID from the batch (column A, first row)
	if len(batch) > 0 && len(batch[0]) > 0 {
		firstMessageID := fmt.Sprintf("%v", batch[0][0])
		
		// Read back from Raw Data to verify
		readResp, readErr := service.Spreadsheets.Values.Get(spreadsheetID, "Raw Data!A:A").Context(ctx).Do()
		if readErr == nil && len(readResp.Values) > 0 {
			// Check if our message ID is in the last few rows
			found := false
			checkRows := len(readResp.Values)
			if checkRows > 10 {
				checkRows = 10 // Check last 10 rows
			}
			startIdx := len(readResp.Values) - checkRows
			if startIdx < 0 {
				startIdx = 0
			}
			
			for i := startIdx; i < len(readResp.Values); i++ {
				if len(readResp.Values[i]) > 0 {
					if fmt.Sprintf("%v", readResp.Values[i][0]) == firstMessageID {
						found = true
						break
					}
				}
			}
			
			if !found {
				fmt.Printf("  ⚠️  Warning: Written message ID %s not found in sheet verification\n", firstMessageID)
			} else {
				fmt.Printf("  ✅ Verified: Message ID %s confirmed in sheet (total rows: %d)\n", firstMessageID, len(readResp.Values))
			}
		}
	}
	
	return nil
}

func writeEmailMappingSnapshot(ctx context.Context, service *sheets.Service, spreadsheetID string, rows []*mappingAggregate) {
	// Overwrite snapshot (avoids duplicate rows)
	_ = createSheetIfNotExists(ctx, service, spreadsheetID, "Email Mapping")
	// Ensure header row
	header := &sheets.ValueRange{Values: [][]interface{}{convertHeaders([]string{"Original Email", "Normalized Email", "First Seen", "Last Seen", "Count"})}}
	_, _ = service.Spreadsheets.Values.Update(spreadsheetID, "Email Mapping!A1", header).ValueInputOption("RAW").Context(ctx).Do()

	// Clear old values (keep header)
	_, _ = service.Spreadsheets.Values.Clear(spreadsheetID, "Email Mapping!A2:Z", &sheets.ClearValuesRequest{}).Context(ctx).Do()

	if len(rows) == 0 {
		return
	}

	values := make([][]interface{}, 0, len(rows))
	for _, r := range rows {
		values = append(values, []interface{}{
			r.Original,
			r.Normalized,
			r.FirstSeen.Format(time.RFC3339),
			r.LastSeen.Format(time.RFC3339),
			r.Count,
		})
	}
	vr := &sheets.ValueRange{Values: values}
	_, _ = service.Spreadsheets.Values.Update(spreadsheetID, "Email Mapping!A2", vr).ValueInputOption("RAW").Context(ctx).Do()
}

func getState(ctx context.Context, service *sheets.Service, spreadsheetID string) *State {
	// Try to get state from State sheet
	resp, err := service.Spreadsheets.Values.Get(spreadsheetID, "State!A2:E2").Context(ctx).Do()
	if err != nil || len(resp.Values) == 0 {
		return &State{ProcessedIDs: []string{}}
	}

	state := &State{ProcessedIDs: []string{}}
	if len(resp.Values[0]) > 0 {
		if idx, err := strconv.Atoi(fmt.Sprintf("%v", resp.Values[0][0])); err == nil {
			state.LastIndex = idx
		}
	}
	if len(resp.Values[0]) > 1 {
		if total, err := strconv.Atoi(fmt.Sprintf("%v", resp.Values[0][1])); err == nil {
			state.TotalProcessed = total
		}
	}
	
	// Load processed IDs from Raw Data sheet using message IDs (column A)
	// This is the source of truth - if message ID exists, it's been processed
	rawDataResp, err := service.Spreadsheets.Values.Get(spreadsheetID, "Raw Data!A2:A").Context(ctx).Do()
	if err == nil && len(rawDataResp.Values) > 0 {
		processedIDs := make([]string, 0, len(rawDataResp.Values))
		processedIDsSet := make(map[string]bool)
		
		for _, row := range rawDataResp.Values {
			if len(row) > 0 {
				if id, ok := row[0].(string); ok && id != "" {
					// Use message ID as unique identifier
					if !processedIDsSet[id] {
						processedIDs = append(processedIDs, id)
						processedIDsSet[id] = true
					}
				}
			}
		}
		state.ProcessedIDs = processedIDs
		if len(processedIDs) > 0 {
			fmt.Printf("  ✅ Loaded %d unique processed message IDs from Raw Data sheet\n", len(processedIDs))
		}
	}

	return state
}

func saveState(ctx context.Context, service *sheets.Service, spreadsheetID string, state *State) error {
	// Ensure State sheet has headers
	vr1 := &sheets.ValueRange{
		Values: [][]interface{}{{"LastIndex", "TotalProcessed", "LastRun", "ProcessedIDsCount", "CurrentJobID", "CurrentJobName"}},
	}
	_, err := service.Spreadsheets.Values.Update(spreadsheetID, "State!A1", vr1).ValueInputOption("RAW").Context(ctx).Do()
	if err != nil {
		return err
	}

	// Update state (processed IDs count for reference, actual IDs are in Raw Data)
	vr2 := &sheets.ValueRange{
		Values: [][]interface{}{{
			state.LastIndex,
			state.TotalProcessed,
			state.LastRun.Format(time.RFC3339),
			len(state.ProcessedIDs),
			state.CurrentJobID,
			state.CurrentJobName,
		}},
	}
	_, err = service.Spreadsheets.Values.Update(spreadsheetID, "State!A2", vr2).ValueInputOption("RAW").Context(ctx).Do()
	return err
}

// updateJobStats updates job statistics
func updateJobStats(ctx context.Context, service *sheets.Service, spreadsheetID, jobID, jobName, agentID string, processed, skipped int, completed bool) {
	_ = createSheetIfNotExists(ctx, service, spreadsheetID, "Job Stats")
	
	// Try to find existing job row
	resp, err := service.Spreadsheets.Values.Get(spreadsheetID, "Job Stats!A2:H").Context(ctx).Do()
	if err != nil {
		resp = &sheets.ValueRange{Values: [][]interface{}{}}
	}
	
	now := time.Now()
	status := "IN_PROGRESS"
	if completed {
		status = "COMPLETED"
	}
	
	// Find existing job or create new
	found := false
	rowIndex := -1
	for i, row := range resp.Values {
		if len(row) > 0 && fmt.Sprintf("%v", row[0]) == jobID {
			found = true
			rowIndex = i + 2 // +2 for header and 0-index
			break
		}
	}
	
	// Get existing stats or use new
	var totalProcessed, totalSkipped int
	var startedAt time.Time
	var agentIDs string
	
	if found && rowIndex > 0 {
		// Update existing
		existingRow := resp.Values[rowIndex-2]
		if len(existingRow) > 2 {
			if started, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", existingRow[2])); err == nil {
				startedAt = started
			}
		}
		if len(existingRow) > 4 {
			if tp, err := strconv.Atoi(fmt.Sprintf("%v", existingRow[4])); err == nil {
				totalProcessed = tp
			}
		}
		if len(existingRow) > 5 {
			if ts, err := strconv.Atoi(fmt.Sprintf("%v", existingRow[5])); err == nil {
				totalSkipped = ts
			}
		}
		if len(existingRow) > 7 {
			agentIDs = fmt.Sprintf("%v", existingRow[7])
		}
	} else {
		startedAt = now
		rowIndex = len(resp.Values) + 2
	}
	
	// IMPORTANT: `processed`/`skipped` passed in are ABSOLUTE counters from the current run,
	// so we should set totals, not add (adding would double-count on every update).
	totalProcessed = processed
	totalSkipped = skipped
	
	// Update agent IDs
	if agentIDs == "" {
		agentIDs = agentID
	} else if !strings.Contains(agentIDs, agentID) {
		agentIDs = agentIDs + ", " + agentID
	}
	
	completedAt := ""
	if completed {
		completedAt = now.Format(time.RFC3339)
	}
	
	// Update or append row
	rowData := []interface{}{
		jobID,
		jobName,
		startedAt.Format(time.RFC3339),
		completedAt,
		totalProcessed,
		totalSkipped,
		status,
		agentIDs,
	}
	
	vr := &sheets.ValueRange{
		Values: [][]interface{}{rowData},
	}
	
	if found {
		// Update existing row
		rangeStr := fmt.Sprintf("Job Stats!A%d", rowIndex)
		service.Spreadsheets.Values.Update(spreadsheetID, rangeStr, vr).ValueInputOption("RAW").Context(ctx).Do()
	} else {
		// Append new row
		service.Spreadsheets.Values.Append(spreadsheetID, "Job Stats", vr).ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
	}
}

// rebuildDerivedFromRawData recomputes derived sheets without reading Gmail.
// Useful when you want to "re-run checks" purely from what is already stored in the spreadsheet.
func rebuildDerivedFromRawData(ctx context.Context, service *sheets.Service, spreadsheetID string) error {
	// Read Raw Data columns:
	// A: Email ID, C: Date, J: Client Email, K: Normalized Client Email
	resp, err := service.Spreadsheets.Values.Get(spreadsheetID, "Raw Data!A2:K").Context(ctx).Do()
	if err != nil {
		return err
	}

	processedIDsSet := make(map[string]bool)
	mapping := make(map[string]*mappingAggregate)
	var rowCount int
	for _, row := range resp.Values {
		if len(row) == 0 {
			continue
		}
		rowCount++

		// Email ID
		if len(row) > 0 {
			id := fmt.Sprintf("%v", row[0])
			if id != "" {
				processedIDsSet[id] = true
			}
		}

		// Date
		var dt time.Time
		if len(row) > 2 {
			if t, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", row[2])); err == nil {
				dt = t
			}
		}

		// Client Email / Normalized Client Email
		if len(row) > 10 {
			client := fmt.Sprintf("%v", row[9])
			normalized := fmt.Sprintf("%v", row[10])
			if client != "" && normalized != "" && client != normalized {
				key := client + "->" + normalized
				agg := mapping[key]
				if agg == nil {
					agg = &mappingAggregate{Original: client, Normalized: normalized, FirstSeen: dt, LastSeen: dt, Count: 0}
					mapping[key] = agg
				}
				agg.Count++
				if !dt.IsZero() {
					if agg.FirstSeen.IsZero() || dt.Before(agg.FirstSeen) {
						agg.FirstSeen = dt
					}
					if agg.LastSeen.IsZero() || dt.After(agg.LastSeen) {
						agg.LastSeen = dt
					}
				}
			}
		}
	}

	// Write Email Mapping snapshot
	snap := make([]*mappingAggregate, 0, len(mapping))
	for _, v := range mapping {
		c := *v
		snap = append(snap, &c)
	}
	writeEmailMappingSnapshot(ctx, service, spreadsheetID, snap)

	// Update State (counts only)
	st := &State{
		LastIndex:      rowCount,
		TotalProcessed: rowCount,
		LastRun:        time.Now(),
		ProcessedIDs:   make([]string, 0, len(processedIDsSet)),
		CurrentJobID:   "",
		CurrentJobName: "rebuild",
	}
	for id := range processedIDsSet {
		st.ProcessedIDs = append(st.ProcessedIDs, id)
	}
	return saveState(ctx, service, spreadsheetID, st)
}
