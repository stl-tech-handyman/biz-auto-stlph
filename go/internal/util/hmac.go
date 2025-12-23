package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// VerifyHMAC verifies an HMAC signature
func VerifyHMAC(secret, signature string, body []byte) bool {
	if signature == "" || secret == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// GenerateHMAC generates an HMAC signature for a body
func GenerateHMAC(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

// ExtractSignature extracts signature from header value (handles "sha256=..." format)
func ExtractSignature(headerValue string) (string, error) {
	if len(headerValue) > 7 && headerValue[:7] == "sha256=" {
		return headerValue[7:], nil
	}
	return headerValue, nil
}

