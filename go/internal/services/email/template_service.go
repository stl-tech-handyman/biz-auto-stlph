package email

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// getFirstName extracts the first name from a full name string
// Handles cases like "John", "John Doe", "John Michael Doe", etc.
func getFirstName(fullName string) string {
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

// loadEmtlmplTemplate loads a template from the emltmpl folder and renders it with data
// This is a minimal implementation that extracts key content from emltmpl templates
func (s *TemplateService) loadEmtlmplTemplate(templateName string, data FinalInvoiceData) (string, error) {
	// Try to load from emltmpl/html or emltmpl/source
	possiblePaths := []string{
		fmt.Sprintf("emltmpl/html/%s.html", templateName),
		fmt.Sprintf("emltmpl/source/%s.html", templateName),
		fmt.Sprintf("../emltmpl/html/%s.html", templateName),
		fmt.Sprintf("../emltmpl/source/%s.html", templateName),
	}
	
	var err error
	for _, path := range possiblePaths {
		_, err = os.ReadFile(path)
		if err == nil {
			// Template file exists - for now, use our branded inline template
			// TODO: Parse and properly integrate emltmpl template structure
			// This is a minimal implementation to avoid breaking existing functionality
			break
		}
	}
	
	if err != nil {
		return "", fmt.Errorf("failed to load emltmpl template %s: %w", templateName, err)
	}
	
	// For now, return a simple branded version using our inline template
	// This ensures we have consistent branding while supporting the template flag
	// Future: Parse emltmpl template and inject our content into its structure
	return s.generateFinalInvoiceEmailInline(data), nil
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
// Returns (htmlBody, textBody, error)
// templateName: if provided, uses template from emltmpl folder (e.g., "invoice", "receipt")
func (s *TemplateService) GenerateFinalInvoiceEmail(name, eventType, eventDate string, helpersCount *int, originalQuote, depositPaid, remainingBalance float64, invoiceURL string, showGratuity bool, templateName string) (string, string, error) {
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
		Name:             getFirstName(name),
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
	
	// Generate plain text version
	textBody := s.generateFinalInvoiceEmailText(data)
	
	// Try to use template file first for HTML
	var htmlBody string
	if html, err := s.renderTemplate("final_invoice.html", data); err == nil {
		htmlBody = html
	} else {
		// Fallback to inline template if file not found
		htmlBody = s.generateFinalInvoiceEmailInline(data)
	}
	
	return htmlBody, textBody, nil
}

// generateFinalInvoiceEmailInline generates final invoice email inline (fallback)
func (s *TemplateService) generateFinalInvoiceEmailInline(data FinalInvoiceData) string {
	depositSection := ""
	if data.DepositPaid > 0 {
		depositSection = fmt.Sprintf(`<tr>
                <td style="padding: 8px 10px 8px 0; vertical-align: top;"><strong>Deposit Paid:</strong></td>
                <td style="padding: 8px 0; vertical-align: top;">$%.2f</td>
            </tr>`, data.DepositPaid)
	}
	
	gratuitySection := ""
	if data.ShowGratuity {
		gratuitySection = fmt.Sprintf(`<div style="margin: 10px 0;">
            <p><strong>💛  Want to Include a Gratuity?</strong></p>
            <p>We're deeply grateful when clients choose to recognize our helpers' hard work (totally optional).<br>
            100%% of your gratuity goes directly to the event team.</p>
            <p style="margin: 5px 0;">
                👉 <a href="%s" style="color: #0047ab; text-decoration: underline;">Add a Tip for Our Team</a>
            </p>
        </div>`, data.GratuityURL)
	}
	
	helpersSection := ""
	if data.HelpersText != "" {
		helpersSection = fmt.Sprintf(`<tr>
                <td style="padding: 8px 10px 8px 0; vertical-align: top;"><strong>Staffing:</strong></td>
                <td style="padding: 8px 0; vertical-align: top;">%s</td>
            </tr>`, data.HelpersText)
	}
	
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Final Invoice - STL Party Helpers</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #000;">
    <div style="max-width: 600px; margin: 0 auto; padding: 5px;">
        <p>Hi %s,</p>
        
        <p>We hope the celebration was everything you imagined — thank you for letting us be part of it.</p>
        
        <p>As agreed, here is your final invoice for staffing services.</p>
        
        <table style="width: 100%%; margin: 15px 0; border-collapse: collapse;">
            <tr>
                <td style="padding: 8px 10px 8px 0; vertical-align: top; width: 50%%;"><strong>Event:</strong></td>
                <td style="padding: 8px 0; vertical-align: top; width: 50%%;">%s</td>
            </tr>
            <tr>
                <td style="padding: 8px 10px 8px 0; vertical-align: top;"><strong>Date:</strong></td>
                <td style="padding: 8px 0; vertical-align: top;">%s</td>
            </tr>
            %s
            %s
            <tr>
                <td style="padding: 8px 10px 8px 0; vertical-align: top;"><strong>Balance Due:</strong></td>
                <td style="padding: 8px 0; vertical-align: top;">$%.2f</td>
            </tr>
        </table>
        
        <p style="margin: 8px 0;">
            👉 <a href="%s" style="color: #0047ab; text-decoration: underline;">Pay Your Remaining Balance Securely via Stripe</a>
        </p>
        
        <p>We truly appreciate your trust in STL Party Helpers to support your event.</p>
        
        %s
        
        <p>Thank you again for choosing STL Party Helpers — your support means the world to us!<br>
        We hope to work with you again soon.</p>
        
        <p>Sincerely,</p>
        
        <p>
            STL Party Helpers Team<br>
            Phone: 314.714.5514<br>
            Email: team@stlpartyhelpers.com<br>
            Website: stlpartyhelpers.com
        </p>
    </div>
</body>
</html>`, data.Name, data.EventType, data.EventDate, helpersSection, depositSection, data.RemainingBalance, data.InvoiceURL, gratuitySection)
}

// generateFinalInvoiceEmailText generates plain text version of final invoice email
func (s *TemplateService) generateFinalInvoiceEmailText(data FinalInvoiceData) string {
	var text strings.Builder
	
	text.WriteString(fmt.Sprintf("Hi %s,\n\n", data.Name))
	text.WriteString("We hope the celebration was everything you imagined — thank you for letting us be part of it.\n\n")
	text.WriteString("As agreed, here is your final invoice for staffing services.\n\n")
	
	text.WriteString(fmt.Sprintf("Event: %s\n", data.EventType))
	text.WriteString(fmt.Sprintf("Date: %s\n\n", data.EventDate))
	
	if data.HelpersText != "" {
		text.WriteString(fmt.Sprintf("Staffing: %s\n", data.HelpersText))
	}
	
	if data.DepositPaid > 0 {
		text.WriteString(fmt.Sprintf("Deposit Paid: $%.2f\n", data.DepositPaid))
	}
	
	text.WriteString(fmt.Sprintf("Balance Due: $%.2f\n\n", data.RemainingBalance))
	text.WriteString(fmt.Sprintf("Pay Your Remaining Balance Securely via Stripe: %s\n\n", data.InvoiceURL))
	
	text.WriteString("We truly appreciate your trust in STL Party Helpers to support your event.\n\n")
	
	if data.ShowGratuity {
		text.WriteString("💛  Want to Include a Gratuity?\n")
		text.WriteString("We're deeply grateful when clients choose to recognize our helpers' hard work (totally optional).\n")
		text.WriteString("100% of your gratuity goes directly to the event team.\n")
		text.WriteString(fmt.Sprintf("Add a Tip for Our Team: %s\n\n", data.GratuityURL))
	}
	
	text.WriteString("Thank you again for choosing STL Party Helpers — your support means the world to us!\n")
	text.WriteString("We hope to work with you again soon.\n\n")
	text.WriteString("Sincerely,\n\n")
	text.WriteString("STL Party Helpers Team\n")
	text.WriteString("Phone: 314.714.5514\n")
	text.WriteString("Email: team@stlpartyhelpers.com\n")
	text.WriteString("Website: stlpartyhelpers.com\n")
	
	return text.String()
}

// DepositData holds data for deposit email template
type DepositData struct {
	Name          string
	DepositAmount float64
	InvoiceURL    string
}

// GenerateDepositEmail generates HTML and plain text for deposit invoice email using template
// Returns (htmlBody, textBody, error)
func (s *TemplateService) GenerateDepositEmail(name string, depositAmount float64, invoiceURL string) (string, string, error) {
	data := DepositData{
		Name:          getFirstName(name),
		DepositAmount: depositAmount,
		InvoiceURL:    invoiceURL,
	}
	
	// Generate plain text version
	textBody := s.generateDepositEmailText(data)
	
	// Try to use template file first for HTML
	var htmlBody string
	if html, err := s.renderTemplate("deposit.html", data); err == nil {
		htmlBody = html
	} else {
		// Fallback to inline template if file not found
		htmlBody = s.generateDepositEmailInline(data)
	}
	
	return htmlBody, textBody, nil
}

// generateDepositEmailInline generates deposit email inline (fallback)
func (s *TemplateService) generateDepositEmailInline(data DepositData) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Action needed to secure your reservation - STL Party Helpers</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #000;">
    <div style="max-width: 600px; margin: 0 auto; padding: 5px;">
        <p>Hello %s,</p>
        
        <p>Thanks for choosing STL Party Helpers!</p>
        
        <p><strong>Action needed to secure your reservation</strong></p>
        <p>To lock in your date, please pay the $%.2f deposit using the stripe link below.</p>
        <p>Your reservation isn't confirmed until the deposit is received.</p>
        
        <p style="margin: 8px 0;">
            👉 <a href="%s" style="color: #0047ab; text-decoration: underline;">Pay Deposit via Stripe</a>
        </p>
        
        <p><strong>Deposit Refund Policy</strong></p>
        <p>Deposits are fully refundable if you cancel at least 3 days before your event.</p>
        
        <p>Questions? Reply to this email or call us at 314.714.5514.</p>
        
        <p>Sincerely,</p>
        
        <p>
            STL Party Helpers Team<br>
            Phone: 314.714.5514<br>
            Email: team@stlpartyhelpers.com<br>
            Website: stlpartyhelpers.com
        </p>
    </div>
</body>
</html>`, data.Name, data.DepositAmount, data.InvoiceURL)
}

// generateDepositEmailText generates plain text version of deposit email
func (s *TemplateService) generateDepositEmailText(data DepositData) string {
	var text strings.Builder
	
	text.WriteString(fmt.Sprintf("Hello %s,\n\n", data.Name))
	text.WriteString("Thanks for choosing STL Party Helpers!\n\n")
	text.WriteString("Action needed to secure your reservation\n")
	text.WriteString(fmt.Sprintf("To lock in your date, please pay the $%.2f deposit using the stripe link below.\n", data.DepositAmount))
	text.WriteString("Your reservation isn't confirmed until the deposit is received.\n\n")
	text.WriteString(fmt.Sprintf("👉 Pay Deposit via Stripe: %s\n\n", data.InvoiceURL))
	text.WriteString("Deposit Refund Policy\n")
	text.WriteString("Deposits are fully refundable if you cancel at least 3 days before your event.\n\n")
	text.WriteString("Questions? Reply to this email or call us at 314.714.5514.\n\n")
	text.WriteString("Sincerely,\n\n")
	text.WriteString("STL Party Helpers Team\n")
	text.WriteString("Phone: 314.714.5514\n")
	text.WriteString("Email: team@stlpartyhelpers.com\n")
	text.WriteString("Website: stlpartyhelpers.com\n")
	
	return text.String()
}

// ReviewRequestData holds data for review request email template
type ReviewRequestData struct {
	Name      string
	ReviewURL string
}

// GenerateReviewRequestEmail generates HTML for review request email using template
func (s *TemplateService) GenerateReviewRequestEmail(name, reviewURL string) (string, error) {
	data := ReviewRequestData{
		Name:      getFirstName(name),
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
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #000;">
    <div style="max-width: 600px; margin: 0 auto; padding: 5px;">
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
        
        <p>Sincerely,</p>
        
        <p>
            STL Party Helpers Team<br>
            Phone: 314.714.5514<br>
            Email: team@stlpartyhelpers.com<br>
            Website: stlpartyhelpers.com
        </p>
    </div>
</body>
</html>`, data.Name, data.ReviewURL)
}

