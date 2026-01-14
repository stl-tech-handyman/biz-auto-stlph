package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
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
		verbose       = flag.Bool("v", false, "Verbose logging")
		batchSize     = flag.Int("batch", 50, "Batch size for processing")
		delay         = flag.Int("delay", 200, "Delay between batches in milliseconds")
		all           = flag.Bool("all", false, "Process all emails (equivalent to -max 0)")
		idempotent    = flag.Bool("idempotent", false, "Recreate sheets if they exist (idempotent mode)")
		agentID       = flag.String("agent", "", "Agent ID for concurrent processing (auto-generated if empty)")
		jobID         = flag.String("job", defaultJobID, "Job ID for this analysis run")
		jobName       = flag.String("job-name", defaultJobName, "Job name/description")
		maxAgents     = flag.Int("max-agents", 5, "Maximum concurrent agents (recommended: 3-5)")
	)
	flag.Parse()

	if *all {
		*maxEmails = 0
	}
	
	// Generate agent ID if not provided
	if *agentID == "" {
		*agentID = fmt.Sprintf("agent-%d-%d", time.Now().Unix(), os.Getpid())
	}
	
	// Validate max agents
	if *maxAgents < 1 {
		*maxAgents = 1
	}
	if *maxAgents > 10 {
		fmt.Printf("⚠️  Warning: Max agents > 10 may hit Gmail API rate limits. Recommended: 3-5\n")
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
		spreadsheet, err := sheetsService.Spreadsheets.Create(&sheets.Spreadsheet{
			Properties: &sheets.SpreadsheetProperties{
				Title: fmt.Sprintf("Email Analysis - %s", time.Now().Format("2006-01-02 15:04")),
			},
		}).Context(ctx).Do()
		if err != nil {
			log.Fatalf("Failed to create spreadsheet: %v", err)
		}
		*spreadsheetID = spreadsheet.SpreadsheetId
		fmt.Printf("Created spreadsheet: %s\n", spreadsheet.SpreadsheetUrl)
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
	if err := initializeSheets(ctx, sheetsService, *spreadsheetID, *idempotent); err != nil {
		log.Fatalf("Failed to initialize sheets: %v", err)
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

	// Default query
	if *query == "" {
		*query = `from:zapier.com OR from:forms OR subject:"New Lead" OR subject:"Form Submission" OR subject:"Quote"`
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
	fmt.Printf("  Max agents: %d\n", *maxAgents)
	fmt.Printf("  Spreadsheet: https://docs.google.com/spreadsheets/d/%s\n\n", *spreadsheetID)

	// Process emails
	startTime := time.Now()
	processed, skipped, err := processEmails(ctx, gmailService, sheetsService, *spreadsheetID, *query, *maxEmails, state.LastIndex, *batchSize, *delay, *verbose, processedIDsSet, &state.ProcessedIDs, *agentID, *jobID, *jobName, *maxAgents)
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
			"Email ID", "Thread ID", "Date", "From Email", "To Email", "Subject", "Body Preview",
			"Is Test", "Is Confirmation", "Client Email", "Normalized Client Email",
			"Event Date", "Total Cost", "Rate", "Hours", "Helpers", "Occasion", "Status",
			"Guests", "Deposit", "Email Type", "Conversation ID", "Message Number",
			"Forwarded From", "Migration Detected", "Job ID", "Job Name",
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
		if idempotent && sheetName != "State" && sheetName != "Locks" {
			// Clear all data except header
			resp, err := service.Spreadsheets.Values.Get(spreadsheetID, fmt.Sprintf("%s!A2:Z", sheetName)).Context(ctx).Do()
			if err == nil && len(resp.Values) > 0 {
				// Clear all rows
				vr := &sheets.ValueRange{Values: [][]interface{}{}}
				rangeStr := fmt.Sprintf("%s!A2:Z%d", sheetName, len(resp.Values)+1)
				service.Spreadsheets.Values.Update(spreadsheetID, rangeStr, vr).ValueInputOption("RAW").Context(ctx).Do()
				fmt.Printf("🗑️  Cleared existing data in %s sheet (idempotent mode)\n", sheetName)
			}
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
	processedIDsSet map[string]bool, processedIDsList *[]string, agentID, jobID, jobName string, maxAgents int) (processed, skipped int, err error) {

	userEmail := newEmail // Use new email as user
	var batch [][]interface{}
	var pageToken string
	emailMapping := make(map[string]string) // original -> normalized
	startTime := time.Now()
	lastStateSave := time.Now()
	lastLockRefresh := time.Now()

	fmt.Printf("🚀 Starting email processing...\n")
	fmt.Printf("  Job: %s - %s\n", jobID, jobName)
	fmt.Printf("  Agent: %s\n", agentID)
	fmt.Printf("  Max concurrent agents: %d\n\n", maxAgents)
	
	// Initialize job stats
	updateJobStats(ctx, sheetsService, spreadsheetID, jobID, jobName, agentID, 0, 0, false)

	for batchNum := 0; ; batchNum++ {
		// Check if we've hit the limit
		if maxEmails > 0 && processed >= maxEmails {
			fmt.Printf("\n✅ Reached max emails limit (%d)\n", maxEmails)
			break
		}

		if verbose {
			fmt.Printf("📦 Processing batch %d...\n", batchNum+1)
		} else if batchNum%10 == 0 {
			fmt.Printf("📦 Batch %d | Processed: %d | Skipped: %d | Elapsed: %s\n",
				batchNum+1, processed, skipped, time.Since(startTime).Round(time.Second))
		}

		// Search Gmail with pagination
		call := gmailService.Users.Messages.List(userEmail).
			Q(query).
			MaxResults(int64(batchSize))

		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		resp, err := call.Context(ctx).Do()
		if err != nil {
			return processed, skipped, fmt.Errorf("failed to search emails: %w", err)
		}

		if len(resp.Messages) == 0 {
			fmt.Printf("\n✅ No more emails found\n")
			break
		}

		// Process each message in this batch
		for i, msg := range resp.Messages {
			// Skip if already processed (using message ID - Gmail's unique identifier)
			if processedIDsSet[msg.Id] {
				if verbose {
					fmt.Printf("  ⏭️  Skipping already processed message ID: %s (thread: %s)\n", msg.Id, msg.ThreadId)
				}
				skipped++
				continue
			}

			emailData, err := processMessage(ctx, gmailService, userEmail, msg.Id, emailMapping, verbose && i < 3)
			if err != nil {
				if verbose {
					fmt.Printf("  ⚠️  Error processing %s: %v\n", msg.Id, err)
				}
				skipped++
				continue
			}

			if emailData == nil || emailData.IsTest {
				skipped++
				continue
			}
			
			// Stamp with job ID
			emailData.JobID = jobID
			emailData.JobName = jobName

			batch = append(batch, emailData.ToRow())
			
			// Track by message ID (Gmail's unique identifier) - this is the key for resume
			*processedIDsList = append(*processedIDsList, msg.Id)
			processedIDsSet[msg.Id] = true
			processed++
			
			if verbose && i < 3 {
				fmt.Printf("  ✅ Processing message ID: %s | Thread: %s\n", msg.Id, msg.ThreadId)
			}

			// Write batch when it reaches size
			if len(batch) >= 25 {
				if err := writeBatch(ctx, sheetsService, spreadsheetID, batch); err != nil {
					return processed, skipped, fmt.Errorf("failed to write batch: %w", err)
				}
				if verbose {
					fmt.Printf("  💾 Wrote batch of %d emails\n", len(batch))
				}
				batch = [][]interface{}{}
				
				// Save state after each batch write (critical for resume)
				state := &State{
					LastIndex:      startIndex + processed + skipped,
					TotalProcessed: processed,
					LastRun:        time.Now(),
					ProcessedIDs:   *processedIDsList,
					CurrentJobID:   jobID,
					CurrentJobName: jobName,
				}
				if err := saveState(ctx, sheetsService, spreadsheetID, state); err == nil {
					lastStateSave = time.Now()
					if verbose {
						fmt.Printf("  💾 State saved (%d message IDs tracked, Job: %s)\n", len(*processedIDsList), jobID)
					}
				} else if verbose {
					fmt.Printf("  ⚠️  Warning: Failed to save state: %v\n", err)
				}
				
				// Update job stats incrementally
				updateJobStats(ctx, sheetsService, spreadsheetID, jobID, jobName, agentID, processed, skipped, false)
			}

			// Progress update
			if processed%100 == 0 {
				elapsed := time.Since(startTime)
				rate := float64(processed) / elapsed.Seconds()
				remaining := ""
				if maxEmails > 0 {
					remaining = fmt.Sprintf(" | Remaining: %d", maxEmails-processed)
				}
				fmt.Printf("  📊 Progress: %d processed | %d skipped | %.1f emails/sec%s\n",
					processed, skipped, rate, remaining)
			}
		}

		// Check if there are more pages
		pageToken = resp.NextPageToken
		if pageToken == "" {
			fmt.Printf("\n✅ Reached end of results\n")
			break
		}

		// Refresh lock every 30 seconds (extend expiration)
		if time.Since(lastLockRefresh) > 30*time.Second {
			// Extend lock expiration
			if err := refreshLock(ctx, sheetsService, spreadsheetID, agentID); err != nil {
				if verbose {
					fmt.Printf("  ⚠️  Warning: Failed to refresh lock: %v\n", err)
				}
			} else if verbose {
				fmt.Printf("  🔒 Lock refreshed (agent: %s)\n", agentID)
			}
			lastLockRefresh = time.Now()
		}
		
		// Auto-save state every 2 minutes (more frequent for safety)
		// Uses message IDs to track what's been processed
		if time.Since(lastStateSave) > 2*time.Minute {
			state := &State{
				LastIndex:      startIndex + processed + skipped,
				TotalProcessed: processed,
				LastRun:        time.Now(),
				ProcessedIDs:   *processedIDsList, // All message IDs processed so far
				CurrentJobID:   jobID,
				CurrentJobName: jobName,
			}
			if err := saveState(ctx, sheetsService, spreadsheetID, state); err == nil {
				lastStateSave = time.Now()
				if verbose {
					fmt.Printf("  💾 Auto-saved state (%d message IDs tracked, Job: %s)\n", len(*processedIDsList), jobID)
				} else {
					fmt.Printf("  💾 State saved (%d message IDs, Job: %s)\n", len(*processedIDsList), jobID)
				}
			}
			
			// Update job stats
			updateJobStats(ctx, sheetsService, spreadsheetID, jobID, jobName, agentID, processed, skipped, false)
		}

		// Rate limiting delay
		if delayMs > 0 {
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		}
	}

	// Write remaining batch
	if len(batch) > 0 {
		if err := writeBatch(ctx, sheetsService, spreadsheetID, batch); err != nil {
			return processed, skipped, fmt.Errorf("failed to write final batch: %w", err)
		}
		if verbose {
			fmt.Printf("  💾 Wrote final batch of %d emails\n", len(batch))
		}
		
		// Save state after final batch
		state := &State{
			LastIndex:      startIndex + processed + skipped,
			TotalProcessed: processed,
			LastRun:        time.Now(),
			ProcessedIDs:   *processedIDsList,
		}
		saveState(ctx, sheetsService, spreadsheetID, state)
	}

	// Write email mapping
	writeEmailMapping(ctx, sheetsService, spreadsheetID, emailMapping)

	return processed, skipped, nil
}

type EmailData struct {
	EmailID            string
	ThreadID          string
	Date              time.Time
	FromEmail         string
	ToEmail           string
	Subject           string
	BodyPreview       string
	IsTest            bool
	IsConfirmation    bool
	ClientEmail       string
	NormalizedClientEmail string
	EventDate         string
	TotalCost         float64
	Rate              float64
	Hours             string
	Helpers           string
	Occasion          string
	Status            string
	Guests            string
	Deposit           float64
	EmailType         string
	ConversationID    string
	MessageNumber     int
	ForwardedFrom     string
	MigrationDetected bool
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

	// Extract body
	body := extractBody(msg.Payload)
	bodyPreview := body
	if len(bodyPreview) > 300 {
		bodyPreview = bodyPreview[:300]
	}

	// Extract data
	combined := strings.ToLower(subject + " " + body)
	isTest := detectTest(combined)
	isConfirmation := detectConfirmation(combined)
	clientEmail := extractClientEmail(normalizedFrom, normalizedTo, body)
	normalizedClientEmail := normalizeEmail(clientEmail)
	
	eventData := extractEventData(body, subject)
	pricingData := extractPricingData(body, subject)
	emailType := classifyEmailType(subject, body, normalizedFrom)

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
		EmailID:              messageID,
		ThreadID:            msg.ThreadId,
		Date:                date,
		FromEmail:           from,
		ToEmail:             to,
		Subject:             subject,
		BodyPreview:         bodyPreview,
		IsTest:              isTest,
		IsConfirmation:      isConfirmation,
		ClientEmail:         clientEmail,
		NormalizedClientEmail: normalizedClientEmail,
		EventDate:           eventData.EventDate,
		TotalCost:           pricingData.TotalCost,
		Rate:                pricingData.Rate,
		Hours:               eventData.Hours,
		Helpers:             eventData.Helpers,
		Occasion:            eventData.Occasion,
		Status:              getStatus(isConfirmation, isTest),
		Guests:              eventData.Guests,
		Deposit:             pricingData.Deposit,
		EmailType:           emailType,
		ConversationID:      conversationID,
		MessageNumber:       1,
		ForwardedFrom:       forwardedFrom,
		MigrationDetected:   migrationDetected,
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
	EventDate string
	Hours     string
	Helpers   string
	Occasion  string
	Guests    string
}

func extractEventData(body, subject string) EventData {
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

	return EventData{EventDate: eventDate}
}

type PricingData struct {
	TotalCost float64
	Rate      float64
	Deposit   float64
}

func extractPricingData(body, subject string) PricingData {
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

	return PricingData{TotalCost: totalCost, Rate: rate, Deposit: deposit}
}

func parseDollarAmount(str string) (float64, error) {
	cleaned := strings.ReplaceAll(strings.ReplaceAll(str, "$", ""), ",", "")
	return strconv.ParseFloat(cleaned, 64)
}

func classifyEmailType(subject, body, from string) string {
	stlphKeywords := []string{"party", "event", "helpers", "booking", "quote", "stl party"}
	combined := strings.ToLower(subject + " " + body)

	for _, kw := range stlphKeywords {
		if strings.Contains(combined, kw) {
			return "STLPH"
		}
	}
	return "Other"
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
	return []interface{}{
		e.EmailID,
		e.ThreadID,
		e.Date.Format(time.RFC3339),
		e.FromEmail,
		e.ToEmail,
		e.Subject,
		e.BodyPreview,
		strconv.FormatBool(e.IsTest),
		strconv.FormatBool(e.IsConfirmation),
		e.ClientEmail,
		e.NormalizedClientEmail,
		e.EventDate,
		e.TotalCost,
		e.Rate,
		e.Hours,
		e.Helpers,
		e.Occasion,
		e.Status,
		e.Guests,
		e.Deposit,
		e.EmailType,
		e.ConversationID,
		e.MessageNumber,
		e.ForwardedFrom,
		strconv.FormatBool(e.MigrationDetected),
		e.JobID,   // Job ID that processed this
		e.JobName, // Job name/description
	}
}

func writeBatch(ctx context.Context, service *sheets.Service, spreadsheetID string, batch [][]interface{}) error {
	vr := &sheets.ValueRange{
		Values: batch,
	}
	_, err := service.Spreadsheets.Values.Append(spreadsheetID, "Raw Data", vr).ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
	return err
}

func writeEmailMapping(ctx context.Context, service *sheets.Service, spreadsheetID string, mapping map[string]string) {
	if len(mapping) == 0 {
		return
	}

	values := [][]interface{}{}
	for original, normalized := range mapping {
		values = append(values, []interface{}{original, normalized, time.Now().Format(time.RFC3339), time.Now().Format(time.RFC3339), 1})
	}

	vr := &sheets.ValueRange{
		Values: values,
	}
	service.Spreadsheets.Values.Append(spreadsheetID, "Email Mapping", vr).ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
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
	
	// Update totals
	totalProcessed += processed
	totalSkipped += skipped
	
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
