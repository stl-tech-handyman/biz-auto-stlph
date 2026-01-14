package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
)

// RegenerateHandler handles quote regeneration requests
type RegenerateHandler struct {
	emailHandler *EmailHandler
	logger       *slog.Logger
}

// NewRegenerateHandler creates a new regenerate handler
func NewRegenerateHandler(emailHandler *EmailHandler, logger *slog.Logger) *RegenerateHandler {
	return &RegenerateHandler{
		emailHandler: emailHandler,
		logger:       logger,
	}
}

// HandleRegenerateQuote handles POST /api/quote/regenerate
// This creates a new quote based on the form data from the expired quote page
func (h *RegenerateHandler) HandleRegenerateQuote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Extract form values
	occasion := r.FormValue("occasion")
	eventDate := r.FormValue("event_date")
	eventTime := r.FormValue("event_time")
	eventLocation := r.FormValue("event_location")
	guestCountStr := r.FormValue("guest_count")
	helpersStr := r.FormValue("helpers")
	hoursStr := r.FormValue("hours")
	clientName := r.FormValue("client_name")
	email := r.FormValue("email")

	// Validate required fields
	if occasion == "" || eventDate == "" || eventTime == "" || eventLocation == "" ||
		guestCountStr == "" || helpersStr == "" || hoursStr == "" || clientName == "" || email == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Validate and parse numbers (for future use when actually calling email handler)
	_, err := strconv.Atoi(guestCountStr)
	if err != nil {
		http.Error(w, "Invalid guest count", http.StatusBadRequest)
		return
	}

	_, err = strconv.Atoi(helpersStr)
	if err != nil {
		http.Error(w, "Invalid helpers count", http.StatusBadRequest)
		return
	}

	_, err = strconv.ParseFloat(hoursStr, 64)
	if err != nil {
		http.Error(w, "Invalid hours", http.StatusBadRequest)
		return
	}

	// TODO: Actually call the email handler to generate and send the new quote
	// This would involve:
	// 1. Creating a proper HTTP request body with all the parsed values
	// 2. Calling h.emailHandler.HandleQuoteEmail(w, r) with the new data
	// 3. Or better yet, extracting the quote generation logic into a shared service

	h.logger.Info("Quote regeneration requested",
		"email", email,
		"occasion", occasion,
		"eventDate", eventDate,
		"eventTime", eventTime,
		"eventLocation", eventLocation,
	)

	// For now, show a success message
	// In production, you'd want to actually generate and send the new quote
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>New Quote Requested - STL Party Helpers</title>
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
        .success-banner {
            background: #d1fae5;
            border-left: 4px solid #10b981;
            padding: 15px;
            margin: 20px 0;
            border-radius: 4px;
        }
        h1 {
            color: rgb(38, 37, 120);
        }
        p {
            line-height: 1.6;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>New Quote Requested</h1>
        <div class="success-banner">
            <p><strong>Thank you!</strong></p>
            <p>We've received your request for a new quote. You'll receive an email with your updated quote shortly.</p>
        </div>
        <p>If you don't receive the email within a few minutes, please check your spam folder or contact us directly.</p>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
