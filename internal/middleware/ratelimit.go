package middleware

import (
	"net/http"
)

// RateLimiter enforces per-client request rate limits on checkout endpoints.
//
// BLOCKER(sec-patch): Currently a pass-through stub with no token bucket enforcement!
// Implement a thread-safe per-client (by X-Client-ID or remote IP) token bucket
// rate limiter returning HTTP 429 Too Many Requests when limit is exceeded.
type RateLimiter struct {
	rps   int
	burst int
}

// NewRateLimiter creates a RateLimiter configured with requests-per-second and burst size.
func NewRateLimiter(rps, burst int) *RateLimiter {
	return &RateLimiter{
		rps:   rps,
		burst: burst,
	}
}

// Allow checks if a request from clientKey is permitted right now.
func (rl *RateLimiter) Allow(clientKey string) bool {
	// TODO(sec-patch): Enforce token-bucket rate limit per clientKey.
	return true
}

// Wrap wraps an http.Handler with rate-limiting middleware.
func (rl *RateLimiter) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientKey := r.Header.Get("X-Client-ID")
		if clientKey == "" {
			clientKey = r.RemoteAddr
		}
		if !rl.Allow(clientKey) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
