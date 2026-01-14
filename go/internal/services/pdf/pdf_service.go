package pdf

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bizops360/go-api/internal/infra/firestore"
	"github.com/bizops360/go-api/internal/infra/storage"
	"github.com/bizops360/go-api/internal/util"
	firestoreLib "cloud.google.com/go/firestore"
)

// Service handles PDF generation, storage, and token management
type Service struct {
	firestoreClient *firestore.Client
	storageClient   *storage.Client
	logger          *slog.Logger
}

// NewService creates a new PDF service
func NewService(firestoreClient *firestore.Client, storageClient *storage.Client, logger *slog.Logger) *Service {
	return &Service{
		firestoreClient: firestoreClient,
		storageClient:   storageClient,
		logger:          logger,
	}
}


// GenerateAndStorePDFAsync generates a PDF, uploads it to storage, and creates a token
// This should be called asynchronously after email is sent
func (s *Service) GenerateAndStorePDFAsync(ctx context.Context, quoteData util.QuoteEmailData, pdfData util.QuotePDFData, expirationDate time.Time) {
	go func() {
		// Use background context for async operation
		bgCtx := context.Background()
		
		if err := s.generateAndStorePDF(bgCtx, quoteData, pdfData, expirationDate); err != nil {
			s.logger.Error("async PDF generation failed", "error", err, "confirmationNumber", pdfData.ConfirmationNumber)
		} else {
			s.logger.Info("PDF generated and stored successfully", "confirmationNumber", pdfData.ConfirmationNumber)
		}
	}()
}

// generateAndStorePDF does the actual work of generating PDF, uploading, and storing token
func (s *Service) generateAndStorePDF(ctx context.Context, quoteData util.QuoteEmailData, pdfData util.QuotePDFData, expirationDate time.Time) error {
	// Generate PDF
	pdfBytes, err := util.GenerateQuotePDF(pdfData)
	if err != nil {
		return fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Upload to Cloud Storage
	storagePath := fmt.Sprintf("quotes/%s/quote-%s.pdf", pdfData.ConfirmationNumber, pdfData.ConfirmationNumber)
	storageURL, err := s.storageClient.UploadFile(ctx, storagePath, pdfBytes, "application/pdf")
	if err != nil {
		return fmt.Errorf("failed to upload PDF to storage: %w", err)
	}

	// Generate token
	token, err := util.GeneratePDFToken(pdfData.ConfirmationNumber, pdfData.ClientEmail, expirationDate)
	if err != nil {
		return fmt.Errorf("failed to generate PDF token: %w", err)
	}

	// Prepare original quote data for regeneration
	originalQuoteData := map[string]interface{}{
		"occasion":      quoteData.Occasion,
		"eventDate":     quoteData.EventDate,
		"eventTime":     quoteData.EventTime,
		"eventLocation": quoteData.EventLocation,
		"guestCount":    quoteData.GuestCount,
		"helpers":       quoteData.Helpers,
		"hours":         quoteData.Hours,
		"clientName":    quoteData.ClientName,
		"clientEmail":   pdfData.ClientEmail,
	}

	// Store token metadata in Firestore
	tokenData := util.PDFTokenData{
		Token:              token,
		ConfirmationNumber: pdfData.ConfirmationNumber,
		ClientEmail:        pdfData.ClientEmail,
		StoragePath:        storagePath,
		StorageURL:         storageURL,
		CreatedAt:          time.Now(),
		ExpiresAt:          expirationDate,
		QuoteExpiresAt:     expirationDate,
		AccessCount:        0,
		OriginalQuoteData:  originalQuoteData,
	}

	_, err = s.firestoreClient.GetClient().Collection("pdf_tokens").Doc(token).Set(ctx, tokenData)
	if err != nil {
		return fmt.Errorf("failed to store token in Firestore: %w", err)
	}

	// Store or update quote record
	quoteRecord := util.QuoteRecord{
		ConfirmationNumber: pdfData.ConfirmationNumber,
		ClientEmail:        pdfData.ClientEmail,
		ClientName:         pdfData.ClientName,
		PDFToken:           token,
		Status:             "sent",
		CreatedAt:          time.Now(),
	}
	
	now := time.Now()
	quoteRecord.PDFGeneratedAt = &now

	_, err = s.firestoreClient.GetClient().Collection("quotes").Doc(pdfData.ConfirmationNumber).Set(ctx, quoteRecord)
	if err != nil {
		s.logger.Warn("failed to store quote record", "error", err, "confirmationNumber", pdfData.ConfirmationNumber)
		// Don't fail the whole operation if quote record fails
	}

	return nil
}

// GetPDFDownloadURL returns the download URL for a PDF token
func (s *Service) GetPDFDownloadURL(ctx context.Context, token string) (string, error) {
	// Get token data from Firestore
	doc, err := s.firestoreClient.GetClient().Collection("pdf_tokens").Doc(token).Get(ctx)
	if err != nil {
		return "", fmt.Errorf("token not found: %w", err)
	}

	var tokenData util.PDFTokenData
	if err := doc.DataTo(&tokenData); err != nil {
		return "", fmt.Errorf("failed to parse token data: %w", err)
	}

	// Check expiration
	if time.Now().After(tokenData.ExpiresAt) {
		return "", fmt.Errorf("token expired")
	}

	// Generate signed URL (valid for 1 hour)
	signedURL, err := s.storageClient.GetSignedURL(ctx, tokenData.StoragePath, 1*time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	// Update access tracking
	now := time.Now()
	updates := []firestoreLib.Update{
		{Path: "accessCount", Value: firestoreLib.Increment(1)},
		{Path: "accessedAt", Value: now},
	}
	
	_, err = s.firestoreClient.GetClient().Collection("pdf_tokens").Doc(token).Update(ctx, updates)
	if err != nil {
		s.logger.Warn("failed to update token access tracking", "error", err)
	}

	// Update quote status
	_, err = s.firestoreClient.GetClient().Collection("quotes").Doc(tokenData.ConfirmationNumber).Update(ctx, []firestoreLib.Update{
		{Path: "status", Value: "pdf_downloaded"},
		{Path: "lastActivityAt", Value: now},
	})
	if err != nil {
		s.logger.Warn("failed to update quote status", "error", err)
	}

	return signedURL, nil
}

// GetTokenData retrieves token data from Firestore
func (s *Service) GetTokenData(ctx context.Context, token string) (*util.PDFTokenData, error) {
	doc, err := s.firestoreClient.GetClient().Collection("pdf_tokens").Doc(token).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("token not found: %w", err)
	}

	var tokenData util.PDFTokenData
	if err := doc.DataTo(&tokenData); err != nil {
		return nil, fmt.Errorf("failed to parse token data: %w", err)
	}

	return &tokenData, nil
}
