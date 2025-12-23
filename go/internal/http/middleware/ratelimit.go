package middleware

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/bizops360/go-api/internal/util"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	mu          sync.Mutex
	clients     map[string]*clientLimiter
	limit       int           // requests per window
	window      time.Duration // time window
	cleanupTick *time.Ticker
	logger      *slog.Logger
}

type clientLimiter struct {
	count     int
	windowEnd time.Time
}

// NewRateLimiter creates a new rate limiter
// limit: number of requests allowed per window
// window: time window duration
func NewRateLimiter(limit int, window time.Duration, logger *slog.Logger) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*clientLimiter),
		limit:   limit,
		window:  window,
		logger:  logger,
	}

	// Cleanup old entries every minute
	rl.cleanupTick = time.NewTicker(1 * time.Minute)
	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) cleanup() {
	for range rl.cleanupTick.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, cl := range rl.clients {
			if now.After(cl.windowEnd) {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cl, exists := rl.clients[ip]

	if !exists || now.After(cl.windowEnd) {
		// New client or window expired, reset
		rl.clients[ip] = &clientLimiter{
			count:     1,
			windowEnd: now.Add(rl.window),
		}
		return true
	}

	if cl.count >= rl.limit {
		return false
	}

	cl.count++
	return true
}

// RateLimitMiddleware limits requests per IP address
func RateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)
			requestID := util.GetRequestID(r.Context())

			if !limiter.Allow(ip) {
				limiter.logger.Warn("rate_limit_exceeded",
					"requestId", requestID,
					"ip", ip,
					"path", r.URL.Path,
					"method", r.Method,
				)

				util.WriteError(w, http.StatusTooManyRequests, "Rate limit exceeded")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for proxies/load balancers)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fallback to RemoteAddr
	return r.RemoteAddr
}


