package handlers

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bizops360/go-api/internal/infra/firestore"
	"github.com/bizops360/go-api/internal/infra/storage"
	"github.com/bizops360/go-api/internal/services/pdf"
	"github.com/bizops360/go-api/internal/util"
)

// PDFHandler handles PDF download requests
type PDFHandler struct {
	pdfService *pdf.Service
	logger     *slog.Logger
}

// NewPDFHandler creates a new PDF handler
func NewPDFHandler(logger *slog.Logger) (*PDFHandler, error) {
	// Initialize Firestore client (optional - will fail gracefully if not configured)
	var firestoreClient *firestore.Client
	ctx := context.Background()
	if projectID := os.Getenv("GCP_PROJECT_ID"); projectID != "" {
		if client, err := firestore.NewClient(ctx, projectID); err == nil {
			firestoreClient = client
			logger.Info("Firestore client initialized for PDF service")
		} else {
			logger.Warn("Firestore client not available", "error", err)
		}
	} else {
		logger.Warn("GCP_PROJECT_ID not set, PDF service will not be available")
	}

	// Initialize Storage client (optional - will fail gracefully if not configured)
	var storageClient *storage.Client
	if bucketName := os.Getenv("GCS_BUCKET_NAME"); bucketName != "" {
		if client, err := storage.NewClient(ctx, bucketName); err == nil {
			storageClient = client
			logger.Info("Cloud Storage client initialized for PDF service")
		} else {
			logger.Warn("Cloud Storage client not available", "error", err)
		}
	} else {
		logger.Warn("GCS_BUCKET_NAME not set, PDF service will not be available")
	}

	// Create PDF service if both clients are available
	var pdfService *pdf.Service
	if firestoreClient != nil && storageClient != nil {
		pdfService = pdf.NewService(firestoreClient, storageClient, logger)
	}

	return &PDFHandler{
		pdfService: pdfService,
		logger:     logger,
	}, nil
}

// HandlePDFDownload handles GET /api/quote/pdf?token=...
func (h *PDFHandler) HandlePDFDownload(w http.ResponseWriter, r *http.Request) {
	if h.pdfService == nil {
		http.Error(w, "PDF service is not configured", http.StatusServiceUnavailable)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		h.showExpiredPage(w, r, nil, "Token is required")
		return
	}

	ctx := r.Context()

	// First validate token signature
	confirmationNumber, email, expiresAt, err := util.ValidatePDFToken(token)
	if err != nil {
		// Token is invalid or expired - show expired page
		tokenPreview := token
	if len(token) > 10 {
		tokenPreview = token[:10] + "..."
	}
	h.logger.Warn("PDF token validation failed", "error", err, "token", tokenPreview)
		h.showExpiredPage(w, r, nil, "This quote link is invalid or has expired.")
		return
	}

	// Get token data from Firestore to check if quote itself expired
	tokenData, err := h.pdfService.GetTokenData(ctx, token)
	if err != nil {
		h.logger.Warn("Failed to get token data from Firestore", "error", err)
		// Still try to show expired page with what we have
		h.showExpiredPage(w, r, nil, "This quote link is invalid or has expired.")
		return
	}

	// Check if quote expired
	if time.Now().After(tokenData.ExpiresAt) {
		h.logger.Info("PDF token expired", "confirmationNumber", confirmationNumber, "expiresAt", expiresAt)
		h.showExpiredPage(w, r, tokenData, fmt.Sprintf("This quote expired on %s.", expiresAt.Format("January 2, 2006")))
		return
	}

	// Get signed URL for PDF download
	downloadURL, err := h.pdfService.GetPDFDownloadURL(ctx, token)
	if err != nil {
		tokenPreview := token
		if len(token) > 10 {
			tokenPreview = token[:10] + "..."
		}
		h.logger.Error("Failed to get PDF download URL", "error", err, "token", tokenPreview)
		h.showExpiredPage(w, r, tokenData, "Failed to generate download link. Please try again.")
		return
	}

	// Redirect to PDF
	h.logger.Info("PDF download requested", "confirmationNumber", confirmationNumber, "email", email)
	http.Redirect(w, r, downloadURL, http.StatusFound)
}

// showExpiredPage renders the expired quote page with regeneration form
func (h *PDFHandler) showExpiredPage(w http.ResponseWriter, r *http.Request, tokenData *util.PDFTokenData, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// If we have token data, show regeneration form with pre-filled data
	if tokenData != nil && tokenData.OriginalQuoteData != nil {
		html := h.renderExpiredQuotePageWithData(tokenData, message)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
		return
	}

	// Otherwise show simple expired message
	html := h.renderSimpleExpiredPage(message)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// renderExpiredQuotePageWithData renders the expired quote page with pre-filled form
func (h *PDFHandler) renderExpiredQuotePageWithData(tokenData *util.PDFTokenData, message string) string {
	originalData := tokenData.OriginalQuoteData

	// Helper function to safely get string from map
	getString := func(m map[string]interface{}, key string) string {
		if val, ok := m[key].(string); ok {
			return val
		}
		return ""
	}

	// Extract data with defaults
	occasion := getString(originalData, "occasion")
	eventDate := getString(originalData, "eventDate")
	eventTime := getString(originalData, "eventTime")
	eventLocation := getString(originalData, "eventLocation")
	clientName := getString(originalData, "clientName")
	clientEmail := tokenData.ClientEmail

	guestCount := 0
	if gc, ok := originalData["guestCount"].(int); ok {
		guestCount = gc
	} else if gc, ok := originalData["guestCount"].(float64); ok {
		guestCount = int(gc)
	}

	helpers := 0
	if h, ok := originalData["helpers"].(int); ok {
		helpers = h
	} else if h, ok := originalData["helpers"].(float64); ok {
		helpers = int(h)
	}

	hours := 4.0
	if h, ok := originalData["hours"].(float64); ok {
		hours = h
	}

	// Convert event date to HTML date input format (YYYY-MM-DD)
	eventDateHTML := eventDate
	if t, err := time.Parse("January 2, 2006", eventDate); err == nil {
		eventDateHTML = t.Format("2006-01-02")
	} else if t, err := time.Parse("Jan 2, 2006", eventDate); err == nil {
		eventDateHTML = t.Format("2006-01-02")
	}

	// Convert event time to HTML time input format (HH:MM)
	eventTimeHTML := eventTime
	if t, err := time.Parse("3:04 PM", eventTime); err == nil {
		eventTimeHTML = t.Format("15:04")
	} else if t, err := time.Parse("15:04", eventTime); err == nil {
		eventTimeHTML = t.Format("15:04")
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Quote Expired - STL Party Helpers</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            max-width: 600px;
            margin: 50px auto;
            padding: 20px;
            background-color: #f5f7fa;
            color: #333;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .expired-banner {
            background: #fff9c4;
            border-left: 4px solid #f59e0b;
            padding: 15px;
            margin: 20px 0;
            border-radius: 4px;
        }
        h1 {
            color: rgb(38, 37, 120);
            margin-top: 0;
        }
        .form-section {
            background: #fafafa;
            padding: 20px;
            border-radius: 8px;
            margin: 20px 0;
        }
        .field {
            margin: 15px 0;
        }
        label {
            display: block;
            font-weight: bold;
            margin-bottom: 5px;
            color: #333;
        }
        input, select {
            width: 100%%;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 14px;
            box-sizing: border-box;
        }
        button {
            background: rgb(38, 37, 120);
            color: white;
            padding: 12px 24px;
            border: none;
            border-radius: 4px;
            font-size: 16px;
            font-weight: bold;
            cursor: pointer;
            width: 100%%;
            margin-top: 10px;
        }
        button:hover {
            background: #1f1e5a;
        }
        .info {
            color: #666;
            font-size: 14px;
            margin-top: 10px;
            font-style: italic;
        }
        .message {
            color: #92400e;
            font-weight: 500;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Quote Expired</h1>
        
        <div class="expired-banner">
            <p class="message"><strong>%s</strong></p>
            <p>To get a new quote with current pricing, please review and update your details below.</p>
        </div>
        
        <p>Many of our clients need to submit quotes to accounting for approval. 
        Generate a fresh quote below to ensure you have the most current pricing.</p>
        
        <form action="/api/quote/regenerate" method="POST" id="regenerateForm">
            <input type="hidden" name="original_quote_id" value="%s">
            <input type="hidden" name="token" value="%s">
            
            <div class="form-section">
                <h3 style="margin-top: 0; color: rgb(38, 37, 120);">Event Details</h3>
                
                <div class="field">
                    <label for="occasion">Occasion *</label>
                    <input type="text" id="occasion" name="occasion" value="%s" required>
                </div>
                
                <div class="field">
                    <label for="event_date">Event Date *</label>
                    <input type="date" id="event_date" name="event_date" value="%s" required>
                </div>
                
                <div class="field">
                    <label for="event_time">Event Time *</label>
                    <input type="time" id="event_time" name="event_time" value="%s" required>
                </div>
                
                <div class="field">
                    <label for="event_location">Event Location *</label>
                    <input type="text" id="event_location" name="event_location" value="%s" required>
                </div>
                
                <div class="field">
                    <label for="guest_count">Guest Count *</label>
                    <input type="number" id="guest_count" name="guest_count" value="%d" min="1" required>
                </div>
                
                <div class="field">
                    <label for="helpers">Number of Helpers *</label>
                    <input type="number" id="helpers" name="helpers" value="%d" min="1" required>
                </div>
                
                <div class="field">
                    <label for="hours">Duration (Hours) *</label>
                    <input type="number" id="hours" name="hours" value="%.1f" min="1" step="0.5" required>
                </div>
                
                <div class="field">
                    <label for="client_name">Your Name *</label>
                    <input type="text" id="client_name" name="client_name" value="%s" required>
                </div>
                
                <div class="field">
                    <label for="email">Your Email *</label>
                    <input type="email" id="email" name="email" value="%s" required>
                </div>
            </div>
            
            <button type="submit">Generate New Quote</button>
        </form>
        
        <p class="info">All fields are pre-filled with your original quote details. 
        You can modify any information before generating a new quote.</p>
    </div>
</body>
</html>`,
		template.HTMLEscapeString(message),
		tokenData.ConfirmationNumber,
		tokenData.Token,
		template.HTMLEscapeString(occasion),
		template.HTMLEscapeString(eventDateHTML),
		template.HTMLEscapeString(eventTimeHTML),
		template.HTMLEscapeString(eventLocation),
		guestCount,
		helpers,
		hours,
		template.HTMLEscapeString(clientName),
		template.HTMLEscapeString(clientEmail),
	)

	return html
}

// renderSimpleExpiredPage renders a simple expired page without form data
func (h *PDFHandler) renderSimpleExpiredPage(message string) string {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Quote Expired - STL Party Helpers</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            max-width: 600px;
            margin: 50px auto;
            padding: 20px;
            background-color: #f5f7fa;
            color: #333;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            text-align: center;
        }
        .expired-banner {
            background: #fff9c4;
            border-left: 4px solid #f59e0b;
            padding: 15px;
            margin: 20px 0;
            border-radius: 4px;
        }
        h1 {
            color: rgb(38, 37, 120);
        }
        a {
            display: inline-block;
            background: rgb(38, 37, 120);
            color: white;
            padding: 12px 24px;
            text-decoration: none;
            border-radius: 4px;
            font-weight: bold;
            margin-top: 20px;
        }
        a:hover {
            background: #1f1e5a;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Quote Not Found</h1>
        <div class="expired-banner">
            <p>%s</p>
        </div>
        <p>This quote link is invalid or has expired.</p>
        <a href="/quote-request">Request a New Quote</a>
    </div>
</body>
</html>`, template.HTMLEscapeString(message))

	return html
}
