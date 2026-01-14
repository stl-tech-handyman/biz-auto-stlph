package sheets

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// SheetsClient wraps Google Sheets API client
type SheetsClient struct {
	service *sheets.Service
}

// NewSheetsClient creates a new Sheets client
func NewSheetsClient() (*SheetsClient, error) {
	credentialsJSON := os.Getenv("GMAIL_CREDENTIALS_JSON")
	if credentialsJSON == "" {
		return nil, fmt.Errorf("GMAIL_CREDENTIALS_JSON environment variable is not set")
	}

	var credsData []byte
	if _, err := os.Stat(credentialsJSON); err == nil {
		credsData, err = os.ReadFile(credentialsJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to read credentials file: %w", err)
		}
	} else {
		credsData = []byte(credentialsJSON)
	}

	ctx := context.Background()
	config, err := google.JWTConfigFromJSON(credsData, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %w", err)
	}

	client := config.Client(ctx)
	service, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create sheets service: %w", err)
	}

	return &SheetsClient{
		service: service,
	}, nil
}

// GetOrCreateSpreadsheet gets or creates a spreadsheet
func (s *SheetsClient) GetOrCreateSpreadsheet(ctx context.Context, name string) (string, error) {
	// For now, always create a new one - you can enhance this to search Drive
	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: name,
		},
	}

	created, err := s.service.Spreadsheets.Create(spreadsheet).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create spreadsheet: %w", err)
	}

	return created.SpreadsheetId, nil
}

// AppendRows appends rows to a sheet
func (s *SheetsClient) AppendRows(ctx context.Context, spreadsheetID, sheetName string, values [][]interface{}) error {
	valueRange := &sheets.ValueRange{
		Values: convertToSheetsValues(values),
	}

	_, err := s.service.Spreadsheets.Values.Append(
		spreadsheetID,
		sheetName,
		valueRange,
	).ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Context(ctx).Do()

	if err != nil {
		return fmt.Errorf("failed to append rows: %w", err)
	}

	return nil
}

// GetRows gets all rows from a sheet
func (s *SheetsClient) GetRows(ctx context.Context, spreadsheetID, sheetName string) ([][]interface{}, error) {
	resp, err := s.service.Spreadsheets.Values.Get(spreadsheetID, sheetName).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}

	if len(resp.Values) == 0 {
		return [][]interface{}{}, nil
	}

	result := make([][]interface{}, len(resp.Values))
	for i, row := range resp.Values {
		result[i] = make([]interface{}, len(row))
		for j, val := range row {
			result[i][j] = val
		}
	}

	return result, nil
}

// UpdateRow updates a specific row in a sheet
func (s *SheetsClient) UpdateRow(ctx context.Context, spreadsheetID, sheetName string, rowIndex int, values []interface{}) error {
	rangeStr := fmt.Sprintf("%s!A%d", sheetName, rowIndex+1)
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{values},
	}

	_, err := s.service.Spreadsheets.Values.Update(
		spreadsheetID,
		rangeStr,
		valueRange,
	).ValueInputOption("RAW").Context(ctx).Do()

	return err
}

// CreateSheet creates a new sheet tab
func (s *SheetsClient) CreateSheet(ctx context.Context, spreadsheetID, sheetName string) error {
	request := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				AddSheet: &sheets.AddSheetRequest{
					Properties: &sheets.SheetProperties{
						Title: sheetName,
					},
				},
			},
		},
	}

	_, err := s.service.Spreadsheets.BatchUpdate(spreadsheetID, request).Context(ctx).Do()
	if err != nil {
		// Sheet might already exist, that's ok
		return nil
	}

	return nil
}

// SetHeaders sets header row in a sheet
func (s *SheetsClient) SetHeaders(ctx context.Context, spreadsheetID, sheetName string, headers []string) error {
	values := make([][]interface{}, 1)
	values[0] = make([]interface{}, len(headers))
	for i, h := range headers {
		values[0][i] = h
	}

	valueRange := &sheets.ValueRange{
		Values: convertToSheetsValues(values),
	}

	_, err := s.service.Spreadsheets.Values.Update(
		spreadsheetID,
		fmt.Sprintf("%s!A1", sheetName),
		valueRange,
	).ValueInputOption("RAW").Context(ctx).Do()

	return err
}

func convertToSheetsValues(values [][]interface{}) [][]interface{} {
	result := make([][]interface{}, len(values))
	for i, row := range values {
		result[i] = make([]interface{}, len(row))
		for j, val := range row {
			result[i][j] = val
		}
	}
	return result
}

// Log helper
func logToFile(message string, data map[string]interface{}) {
	logPath := "c:\\Users\\Alexey\\Code\\biz-operating-system\\stlph\\.cursor\\debug.log"
	if logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		logEntry := map[string]interface{}{
			"sessionId":    "email-analysis",
			"runId":        "run1",
			"hypothesisId":  "SHEETS",
			"location":     "sheets.go",
			"message":      message,
			"data":         data,
			"timestamp":    time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
}
