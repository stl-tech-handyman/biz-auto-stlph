package email

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

// TemplateService handles email template generation
type TemplateService struct {
	templates map[string]*template.Template
}

// NewTemplateService creates a new template service and loads templates
func NewTemplateService() *TemplateService {
	ts := &TemplateService{
		templates: make(map[string]*template.Template),
	}
	ts.loadTemplates()
	return ts
}

// loadTemplates loads all email templates from the templates/email directory
func (s *TemplateService) loadTemplates() {
	templateFiles := []string{
		"final_invoice.html",
		"deposit.html",
		"review_request.html",
	}

	// Try multiple possible paths for templates
	possiblePaths := []string{
		filepath.Join("templates", "email"),
		filepath.Join("go", "templates", "email"),
		filepath.Join(".", "templates", "email"),
	}

	for _, filename := range templateFiles {
		var tmpl *template.Template
		var err error
		
		// Try each possible path
		for _, basePath := range possiblePaths {
			templatePath := filepath.Join(basePath, filename)
			tmpl, err = template.ParseFiles(templatePath)
			if err == nil {
				s.templates[filename] = tmpl
				break
			}
		}
		
		// If all paths failed, template will use inline fallback
	}
}

// renderTemplate renders a template with the given data
func (s *TemplateService) renderTemplate(templateName string, data interface{}) (string, error) {
	if tmpl, ok := s.templates[templateName]; ok {
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
		}
		return buf.String(), nil
	}
	return "", fmt.Errorf("template %s not found", templateName)
}

// FinalInvoiceData holds data for final invoice email template
type FinalInvoiceData struct {
	Name             string
	EventType        string
	EventDate        string
	HelpersText      string
	OriginalQuote    float64
	DepositPaid      float64
	RemainingBalance float64
	InvoiceURL       string
	GratuityURL      string
	ShowGratuity     bool
}

// GenerateFinalInvoiceEmail generates HTML for final invoice email using template
func (s *TemplateService) GenerateFinalInvoiceEmail(name, eventType, eventDate string, helpersCount *int, originalQuote, depositPaid, remainingBalance float64, invoiceURL string, showGratuity bool) (string, error) {
	// Format helpers text
	helpersText := ""
	if helpersCount != nil {
		if *helpersCount == 1 {
			helpersText = "1 Helper"
		} else {
			helpersText = fmt.Sprintf("%d Helpers", *helpersCount)
		}
	}
	
	// Build gratuity URL (add query parameter to invoice URL)
	gratuityURL := invoiceURL
	if gratuityURL != "" {
		separator := "?"
		if strings.Contains(gratuityURL, "?") {
			separator = "&"
		}
		gratuityURL = gratuityURL + separator + "gratuity=true"
	}
	
	data := FinalInvoiceData{
		Name:             name,
		EventType:        eventType,
		EventDate:        eventDate,
		HelpersText:      helpersText,
		OriginalQuote:    originalQuote,
		DepositPaid:      depositPaid,
		RemainingBalance: remainingBalance,
		InvoiceURL:       invoiceURL,
		GratuityURL:      gratuityURL,
		ShowGratuity:     showGratuity,
	}
	
	// Try to use template file first
	if html, err := s.renderTemplate("final_invoice.html", data); err == nil {
		return html, nil
	}
	
	// Fallback to inline template if file not found
	return s.generateFinalInvoiceEmailInline(data), nil
}

// generateFinalInvoiceEmailInline generates final invoice email inline (fallback)
func (s *TemplateService) generateFinalInvoiceEmailInline(data FinalInvoiceData) string {
	depositSection := ""
	if data.DepositPaid > 0 {
		depositSection = fmt.Sprintf(`<p><strong>Deposit Paid:</strong> $%.2f</p>`, data.DepositPaid)
	}
	
	gratuitySection := ""
	if data.ShowGratuity {
		gratuitySection = fmt.Sprintf(`<div style="margin: 20px 0;">
            <p><strong>💛  Want to Include a Gratuity?</strong></p>
            <p>We're deeply grateful when clients choose to recognize our helpers' hard work (totally optional).<br>
            100%% of your gratuity goes directly to the event team.</p>
            <p style="text-align: center; margin: 15px 0;">
                <a href="%s" 
                   style="display: inline-block; background-color: #0047ab; color: #fff; padding: 10px 20px; text-decoration: none; border-radius: 5px;">
                    👉 Add a Tip for Our Team
                </a>
            </p>
        </div>`, data.GratuityURL)
	}
	
	helpersSection := ""
	if data.HelpersText != "" {
		helpersSection = fmt.Sprintf(`<p><strong>Staffing:</strong> %s</p>`, data.HelpersText)
	}
	
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Final Invoice - STL Party Helpers</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <p>Hi %s,</p>
        
        <p>We hope the celebration was everything you imagined — thank you for letting us be part of it.</p>
        
        <p>As agreed, here is your final invoice for staffing services.</p>
        
        <div style="margin: 20px 0;">
            <p><strong>Event:</strong> %s</p>
            <p><strong>Date:</strong> %s</p>
        </div>
        
        <div style="margin: 20px 0;">
            %s
            <p><strong>Original Quote (USD):</strong> $%.2f</p>
            %s
            <p><strong>Balance Due:</strong> $%.2f</p>
        </div>
        
        <p style="text-align: center; margin: 30px 0;">
            <a href="%s" 
               style="display: inline-block; background-color: #0047ab; color: #fff; padding: 12px 24px; text-decoration: none; border-radius: 5px; font-weight: bold;">
                👉 Click, to Pay Balance Now
            </a>
        </p>
        
        <p>We truly appreciate your trust in STL Party Helpers to support your event.</p>
        
        %s
        
        <p>Thank you again for choosing STL Party Helpers — your support means the world to us!<br>
        We hope to work with you again soon.</p>
        
        <p>Anna</p>
        
        <p style="font-size: 0.9em; color: #666;">
            <strong>Administrative Assistant</strong><br>
            STL Party Helpers Team<br><br>
            Phone: 314.714.5514<br>
            Email: team@stlpartyhelpers.com<br>
            Website: stlpartyhelpers.com
        </p>
    </div>
</body>
</html>`, data.Name, data.EventType, data.EventDate, helpersSection, data.OriginalQuote, depositSection, data.RemainingBalance, data.InvoiceURL, gratuitySection)
}

// DepositData holds data for deposit email template
type DepositData struct {
	Name          string
	DepositAmount float64
	InvoiceURL    string
}

// GenerateDepositEmail generates HTML for deposit invoice email using template
func (s *TemplateService) GenerateDepositEmail(name string, depositAmount float64, invoiceURL string) (string, error) {
	data := DepositData{
		Name:          name,
		DepositAmount: depositAmount,
		InvoiceURL:    invoiceURL,
	}
	
	// Try to use template file first
	if html, err := s.renderTemplate("deposit.html", data); err == nil {
		return html, nil
	}
	
	// Fallback to inline template if file not found
	return s.generateDepositEmailInline(data), nil
}

// generateDepositEmailInline generates deposit email inline (fallback)
func (s *TemplateService) generateDepositEmailInline(data DepositData) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Booking Deposit - STL Party Helpers</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <p>Hi %s,</p>
        
        <p>Thank you for choosing STL Party Helpers!</p>
        
        <p>Your booking deposit invoice has been created. Please find the details below.</p>
        
        <div style="background-color: #f9f9f9; padding: 15px; border-radius: 5px; margin: 20px 0;">
            <h2 style="margin-top: 0;">Deposit Details</h2>
            <p><strong>Deposit Amount:</strong> <strong style="color: #0047ab; font-size: 1.2em;">$%.2f</strong></p>
        </div>
        
        <p style="text-align: center; margin: 30px 0;">
            <a href="%s" 
               style="display: inline-block; background-color: #0047ab; color: #fff; padding: 12px 24px; text-decoration: none; border-radius: 5px; font-weight: bold;">
                Pay Deposit
            </a>
        </p>
        
        <p style="font-size: 0.9em; color: #666;">
            If you have any questions about this deposit, please don't hesitate to contact us.
        </p>
        
        <hr style="border: none; border-top: 1px solid #ddd; margin: 30px 0;">
        
        <p style="font-size: 0.85em; color: #666; text-align: center;">
            STL Party Helpers<br>
            4220 Duncan Ave., Ste. 201, St. Louis, MO 63110<br>
            <a href="tel:+13147145514" style="color: #0047ab;">(314) 714-5514</a><br>
            <a href="https://stlpartyhelpers.com" style="color: #0047ab;">stlpartyhelpers.com</a>
        </p>
    </div>
</body>
</html>`, data.Name, data.DepositAmount, data.InvoiceURL)
}

// ReviewRequestData holds data for review request email template
type ReviewRequestData struct {
	Name      string
	ReviewURL string
}

// GenerateReviewRequestEmail generates HTML for review request email using template
func (s *TemplateService) GenerateReviewRequestEmail(name, reviewURL string) (string, error) {
	data := ReviewRequestData{
		Name:      name,
		ReviewURL: reviewURL,
	}
	
	// Try to use template file first
	if html, err := s.renderTemplate("review_request.html", data); err == nil {
		return html, nil
	}
	
	// Fallback to inline template if file not found
	return s.generateReviewRequestEmailInline(data), nil
}

// generateReviewRequestEmailInline generates review request email inline (fallback)
func (s *TemplateService) generateReviewRequestEmailInline(data ReviewRequestData) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Review Request - STL Party Helpers</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <p>Hi %s,</p>
        
        <p>We hope you had a wonderful experience with STL Party Helpers!</p>
        
        <p>Your feedback means the world to us and helps us continue to provide exceptional service. We would be incredibly grateful if you could take a moment to share your experience.</p>
        
        <p style="text-align: center; margin: 30px 0;">
            <a href="%s" 
               style="display: inline-block; background-color: #0047ab; color: #fff; padding: 12px 24px; text-decoration: none; border-radius: 5px; font-weight: bold;">
                👉 Leave a Review
            </a>
        </p>
        
        <p>Thank you again for choosing STL Party Helpers — we truly appreciate your support!</p>
        
        <p>Anna</p>
        
        <p style="font-size: 0.9em; color: #666;">
            <strong>Administrative Assistant</strong><br>
            STL Party Helpers Team<br><br>
            Phone: 314.714.5514<br>
            Email: team@stlpartyhelpers.com<br>
            Website: stlpartyhelpers.com
        </p>
    </div>
</body>
</html>`, data.Name, data.ReviewURL)
}

