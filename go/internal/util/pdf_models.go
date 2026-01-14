package util

import (
	"time"
)

// PDFTokenData represents a PDF token stored in Firestore
type PDFTokenData struct {
	Token              string                 `firestore:"token"`
	ConfirmationNumber string                 `firestore:"confirmationNumber"`
	ClientEmail        string                 `firestore:"clientEmail"`
	StoragePath        string                 `firestore:"storagePath"`
	StorageURL         string                 `firestore:"storageURL"`
	CreatedAt          time.Time              `firestore:"createdAt"`
	ExpiresAt          time.Time              `firestore:"expiresAt"`
	QuoteExpiresAt     time.Time              `firestore:"quoteExpiresAt"`
	AccessCount        int                    `firestore:"accessCount"`
	AccessedAt         *time.Time             `firestore:"accessedAt,omitempty"`
	OriginalQuoteData  map[string]interface{} `firestore:"originalQuoteData"`
}

// QuoteRecord represents a quote record in Firestore
type QuoteRecord struct {
	ConfirmationNumber string     `firestore:"confirmationNumber"`
	ClientEmail        string     `firestore:"clientEmail"`
	ClientName         string     `firestore:"clientName"`
	PDFToken           string     `firestore:"pdfToken"`
	PDFGeneratedAt     *time.Time `firestore:"pdfGeneratedAt,omitempty"`
	Status             string     `firestore:"status"` // "sent", "pdf_downloaded", "expired"
	CreatedAt          time.Time  `firestore:"createdAt"`
	LastActivityAt     *time.Time `firestore:"lastActivityAt,omitempty"`
}
