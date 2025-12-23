package email

import "fmt"

// TemplateService handles email template generation
type TemplateService struct{}

// NewTemplateService creates a new template service
func NewTemplateService() *TemplateService {
	return &TemplateService{}
}

// GenerateFinalInvoiceEmail generates HTML for final invoice email
func (s *TemplateService) GenerateFinalInvoiceEmail(name, eventType, eventDate string, helpersCount *int, originalQuote, depositPaid, remainingBalance float64, invoiceURL string) string {
	// Format helpers text
	helpersText := ""
	if helpersCount != nil {
		if *helpersCount == 1 {
			helpersText = "1 Helper"
		} else {
			helpersText = fmt.Sprintf("%d Helpers", *helpersCount)
		}
	}
	
	// Format deposit section (only show if > 0)
	depositSection := ""
	if depositPaid > 0 {
		depositSection = fmt.Sprintf(`<p><strong>Deposit Paid:</strong> $%.2f</p>`, depositPaid)
	}
	
	// Build gratuity URL (add query parameter to invoice URL)
	gratuityURL := invoiceURL
	if gratuityURL != "" {
		if len(gratuityURL) > 0 {
			separator := "?"
			if len(gratuityURL) > 0 {
				// Check if URL already has query parameters
				for _, char := range gratuityURL {
					if char == '?' {
						separator = "&"
						break
					}
				}
			}
			gratuityURL = gratuityURL + separator + "gratuity=true"
		}
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
            <p><strong>Staffing:</strong> %s</p>
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
        
        <div style="margin: 20px 0;">
            <p><strong>💛  Want to Include a Gratuity?</strong></p>
            <p>We're deeply grateful when clients choose to recognize our helpers' hard work (totally optional).<br>
            100%% of your gratuity goes directly to the event team.</p>
            <p style="text-align: center; margin: 15px 0;">
                <a href="%s" 
                   style="display: inline-block; background-color: #0047ab; color: #fff; padding: 10px 20px; text-decoration: none; border-radius: 5px;">
                    👉 Add a Tip for Our Team
                </a>
            </p>
        </div>
        
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
</html>`, name, eventType, eventDate, helpersText, originalQuote, depositSection, remainingBalance, invoiceURL, gratuityURL)
}

// GenerateDepositEmail generates HTML for deposit invoice email
func (s *TemplateService) GenerateDepositEmail(name string, depositAmount float64, invoiceURL string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Booking Deposit - STL Party Helpers</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #0047ab;">Hello %s!</h1>
        
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
</html>`, name, depositAmount, invoiceURL)
}

