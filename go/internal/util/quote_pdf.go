package util

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf/v2"
)

// QuotePDFData contains all data needed for quote PDF generation
type QuotePDFData struct {
	ConfirmationNumber string
	Occasion           string
	ClientName         string
	ClientEmail        string
	EventDate          string
	EventTime          string
	HelpersCount       int
	Hours              float64
	TotalCost          float64
	DepositAmount      float64
	ExpirationDate     time.Time
	DepositLink        string
	IssueDate          time.Time
}

// GenerateQuotePDF generates a PDF quote document
func GenerateQuotePDF(data QuotePDFData) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 24)
	pdf.SetTextColor(38, 37, 120)
	pdf.Cell(0, 10, fmt.Sprintf("%s Quote", data.Occasion))
	pdf.Ln(15)

	// Company Info
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(0, 5, "STL Party Helpers LLC")
	pdf.Ln(5)
	pdf.Cell(0, 5, "4220 Duncan Ave., Ste. 201")
	pdf.Ln(5)
	pdf.Cell(0, 5, "St. Louis, Missouri 63110")
	pdf.Ln(5)
	pdf.Cell(0, 5, "United States")
	pdf.Ln(5)
	pdf.Cell(0, 5, "+1 314-714-5514")
	pdf.Ln(5)
	pdf.Cell(0, 5, "team@stlpartyhelpers.com")
	pdf.Ln(15)

	// Quote Number and Dates
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 5, fmt.Sprintf("Quote number: %s", data.ConfirmationNumber))
	pdf.Ln(7)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 5, fmt.Sprintf("Date of issue: %s", data.IssueDate.Format("January 2, 2006 at 3:04 PM")))
	pdf.Ln(5)
	pdf.Cell(0, 5, fmt.Sprintf("Expiration date: %s", data.ExpirationDate.Format("January 2, 2006 at 3:04 PM")))
	pdf.Ln(15)

	// Bill To
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 5, "Bill to")
	pdf.Ln(7)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 5, data.ClientName)
	pdf.Ln(5)
	pdf.Cell(0, 5, data.ClientEmail)
	pdf.Ln(15)

	// Event Details
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 5, "Event Details")
	pdf.Ln(7)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 5, fmt.Sprintf("Event Date & Time: %s %s", data.EventDate, data.EventTime))
	pdf.Ln(5)
	pdf.Cell(0, 5, fmt.Sprintf("Event Type: %s", data.Occasion))
	pdf.Ln(5)
	pdf.Cell(0, 5, fmt.Sprintf("Helpers Count: %d Helper%s", data.HelpersCount, pluralize(data.HelpersCount)))
	pdf.Ln(5)
	pdf.Cell(0, 5, fmt.Sprintf("Hours: %.0f Hour%s", data.Hours, pluralizeFloat(data.Hours)))
	pdf.Ln(15)

	// Pricing Table
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 5, "Pricing")
	pdf.Ln(7)

	// Table header
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(100, 7, "Description", "1", 0, "", false, 0, "")
	pdf.CellFormat(40, 7, "Qty", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 7, "Amount", "1", 0, "R", false, 0, "")
	pdf.Ln(7)

	// Table row
	pdf.SetFont("Arial", "", 10)
	description := fmt.Sprintf("Event Staffing Services - %s", data.Occasion)
	pdf.CellFormat(100, 7, description, "1", 0, "", false, 0, "")
	pdf.CellFormat(40, 7, "1", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 7, fmt.Sprintf("$%.2f", data.TotalCost), "1", 0, "R", false, 0, "")
	pdf.Ln(7)

	// Totals
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(100, 7, "Subtotal", "1", 0, "R", false, 0, "")
	pdf.CellFormat(80, 7, fmt.Sprintf("$%.2f", data.TotalCost), "1", 0, "R", false, 0, "")
	pdf.Ln(7)
	pdf.CellFormat(100, 7, "Total", "1", 0, "R", false, 0, "")
	pdf.CellFormat(80, 7, fmt.Sprintf("$%.2f", data.TotalCost), "1", 0, "R", false, 0, "")
	pdf.Ln(7)
	pdf.CellFormat(100, 7, "Deposit Amount", "1", 0, "R", false, 0, "")
	pdf.CellFormat(80, 7, fmt.Sprintf("$%.2f", data.DepositAmount), "1", 0, "R", false, 0, "")
	pdf.Ln(7)
	pdf.CellFormat(100, 7, "Amount due", "1", 0, "R", false, 0, "")
	pdf.CellFormat(80, 7, fmt.Sprintf("$%.2f USD", data.TotalCost), "1", 0, "R", false, 0, "")
	pdf.Ln(15)

	// Pay Online Link
	if data.DepositLink != "" && data.DepositLink != "#" {
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(0, 5, "Pay online")
		pdf.Ln(5)
		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(38, 37, 120)
		pdf.Cell(0, 5, data.DepositLink)
		pdf.Ln(10)
		pdf.SetTextColor(0, 0, 0)
	}

	// Footer
	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(128, 128, 128)
	pdf.Cell(0, 5, "This is a quote, not a confirmed reservation. Your reservation is confirmed only after deposit payment.")
	pdf.Ln(5)
	pdf.Cell(0, 5, fmt.Sprintf("This quote expires on %s", data.ExpirationDate.Format("January 2, 2006 at 3:04 PM")))

	// Generate PDF bytes
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}

func pluralizeFloat(count float64) string {
	if count == 1.0 {
		return ""
	}
	return "s"
}

