package middleware

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/bizops360/go-api/internal/util"
)

// HMACAuthMiddleware verifies HMAC signature if X-Signature header is present
func HMACAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signature := r.Header.Get("X-Signature")
		if signature == "" {
			// No signature provided, continue (optional auth)
			next.ServeHTTP(w, r)
			return
		}

		// Read body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}
		r.Body.Close()

		// Create new reader for downstream handlers
		r.Body = io.NopCloser(bytes.NewReader(body))

		// Get secret from environment (can be per-business later)
		secret := os.Getenv("HMAC_SECRET")
		if secret == "" {
			// No secret configured, skip verification
			next.ServeHTTP(w, r)
			return
		}

		// Extract signature (handle "sha256=..." format)
		sig, err := util.ExtractSignature(signature)
		if err != nil {
			http.Error(w, "invalid signature format", http.StatusBadRequest)
			return
		}

		// Verify signature
		if !util.VerifyHMAC(secret, sig, body) {
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
