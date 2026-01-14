package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"
)

// GeneratePDFToken generates a signed token for PDF download
// Token format: base64(confirmationNumber|email|expiresAt|signature)
func GeneratePDFToken(confirmationNumber, email string, expiresAt time.Time) (string, error) {
	secretKey := getPDFTokenSecret()
	if secretKey == "" {
		return "", fmt.Errorf("PDF_TOKEN_SECRET environment variable not set")
	}

	// Create payload
	expiresAtUnix := expiresAt.Unix()
	payload := fmt.Sprintf("%s|%s|%d", confirmationNumber, email, expiresAtUnix)

	// Generate HMAC signature
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(payload))
	signature := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	// Combine payload and signature
	tokenData := fmt.Sprintf("%s|%s", payload, signature)

	// Encode to base64 URL-safe
	token := base64.URLEncoding.EncodeToString([]byte(tokenData))

	return token, nil
}

// ValidatePDFToken validates a PDF token and returns the confirmation number and email
func ValidatePDFToken(token string) (confirmationNumber, email string, expiresAt time.Time, err error) {
	secretKey := getPDFTokenSecret()
	if secretKey == "" {
		return "", "", time.Time{}, fmt.Errorf("PDF_TOKEN_SECRET environment variable not set")
	}

	// Decode base64
	tokenData, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("invalid token format: %w", err)
	}

	// Split payload and signature
	parts := strings.Split(string(tokenData), "|")
	if len(parts) != 4 {
		return "", "", time.Time{}, fmt.Errorf("invalid token structure")
	}

	confirmationNumber = parts[0]
	email = parts[1]
	expiresAtUnix := int64(0)
	if _, err := fmt.Sscanf(parts[2], "%d", &expiresAtUnix); err != nil {
		return "", "", time.Time{}, fmt.Errorf("invalid expiration timestamp: %w", err)
	}
	expiresAt = time.Unix(expiresAtUnix, 0)
	signature := parts[3]

	// Reconstruct payload
	payload := fmt.Sprintf("%s|%s|%d", confirmationNumber, email, expiresAtUnix)

	// Verify signature
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(payload))
	expectedSignature := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return "", "", time.Time{}, fmt.Errorf("invalid token signature")
	}

	// Check expiration
	if time.Now().After(expiresAt) {
		return confirmationNumber, email, expiresAt, fmt.Errorf("token expired")
	}

	return confirmationNumber, email, expiresAt, nil
}

// getPDFTokenSecret gets the secret key for PDF tokens from environment
func getPDFTokenSecret() string {
	secret := os.Getenv("PDF_TOKEN_SECRET")
	if secret == "" {
		// Fallback to service API key if PDF_TOKEN_SECRET not set
		secret = os.Getenv("SERVICE_API_KEY")
	}
	return secret
}
