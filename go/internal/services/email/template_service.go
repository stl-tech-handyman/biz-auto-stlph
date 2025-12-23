package email

import "fmt"

// TemplateService handles email template generation
type TemplateService struct{}

// NewTemplateService creates a new template service
func NewTemplateService() *TemplateService {
	return &TemplateService{}
}

// GenerateFinalInvoiceEmail generates HTML for final invoice email
func (s *TemplateService) GenerateFinalInvoiceEmail(name string, totalAmount, depositPaid, remainingBalance float64, invoiceURL string) string {
	// Only show deposit info if deposit was actually paid (> 0)
	depositSection := ""
	if depositPaid > 0 {
		depositSection = fmt.Sprintf(`<p><strong>Deposit Paid:</strong> $%.2f</p>
            `, depositPaid)
	}
	
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Final Invoice - STL Party Helpers</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #0047ab;">Hello %s!</h1>
        
        <p>Thank you for your business with STL Party Helpers!</p>
        
        <p>Your event has been completed. Please find your final invoice below for the remaining balance.</p>
        
        <div style="background-color: #f9f9f9; padding: 15px; border-radius: 5px; margin: 20px 0;">
            <h2 style="margin-top: 0;">Invoice Details</h2>
            <p><strong>Total Event Cost:</strong> $%.2f</p>
            %s
            <p><strong>Remaining Balance:</strong> <strong style="color: #0047ab; font-size: 1.2em;">$%.2f</strong></p>
        </div>
        
        <p style="text-align: center; margin: 30px 0;">
            <a href="%s" 
               style="display: inline-block; background-color: #0047ab; color: #fff; padding: 12px 24px; text-decoration: none; border-radius: 5px; font-weight: bold;">
                Pay Final Invoice
            </a>
        </p>
        
        <p style="font-size: 0.9em; color: #666;">
            If you have any questions about this invoice, please don't hesitate to contact us.
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
</html>`, name, totalAmount, depositSection, remainingBalance, invoiceURL)
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

