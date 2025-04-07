package middlewares

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// IPRateLimiter stores rate limiters for individual IP addresses
type IPRateLimiter struct {
	ips    map[string]*rate.Limiter
	mu     sync.RWMutex
	rate   rate.Limit
	burst  int
	banned map[string]time.Time
	// Duration to ban IPs that exceed rate limits
	banDuration time.Duration
}

// NewIPRateLimiter creates a new rate limiter for IP addresses
// rate: requests per second
// burst: maximum burst size
// banDuration: how long to ban IPs that exceed the rate limit
func NewIPRateLimiter(r rate.Limit, b int, banDuration time.Duration) *IPRateLimiter {
	limiter := &IPRateLimiter{
		ips:         make(map[string]*rate.Limiter),
		rate:        r,
		burst:       b,
		banned:      make(map[string]time.Time),
		banDuration: banDuration,
	}
	// start the clanup job
	go limiter.CleanupJob(5 * time.Minute)

	return limiter
}

// AddIP creates a new rate limiter and adds it to the ips map,
// using the IP address as the key
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.rate, i.burst)
	i.ips[ip] = limiter

	return limiter
}

// GetLimiter returns the rate limiter for the provided IP address
// if it exists, otherwise calls AddIP to add a new one
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()

	// Check if IP is banned
	if bannedUntil, exists := i.banned[ip]; exists {
		if time.Now().Before(bannedUntil) {
			i.mu.RUnlock()
			return nil // IP is banned
		}
		// Ban has expired, remove from banned list
		i.mu.RUnlock()
		i.mu.Lock()
		delete(i.banned, ip)
		i.mu.Unlock()
		i.mu.RLock()
	}

	limiter, exists := i.ips[ip]
	i.mu.RUnlock()

	if !exists {
		return i.AddIP(ip)
	}

	return limiter
}

// BanIP bans an IP for the configured duration
func (i *IPRateLimiter) BanIP(ip string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.banned[ip] = time.Now().Add(i.banDuration)
}

// CleanupJob removes old IP entries to prevent memory leaks
// Should be run periodically, e.g., via a goroutine
func (i *IPRateLimiter) CleanupJob(cleanupInterval time.Duration) {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		i.mu.Lock()

		// Clean expired banned IPs
		for ip, bannedUntil := range i.banned {
			if time.Now().After(bannedUntil) {
				delete(i.banned, ip)
			}
		}

		// Could also implement cleanup of old rate limiters here
		// This would require tracking when they were last used

		i.mu.Unlock()
	}
}

// GetClientIP extracts the client IP address from the request
func GetClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header first (for proxies)
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0]); ip != "" {
		return ip
	}

	// Check for X-Real-IP header next
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // Return as is if there's an error
	}

	return ip
}

// MaliciousRateLimitMiddleware creates a middleware that limits requests based on IP
func MaliciousRateLimitMiddleware(limiter *IPRateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := GetClientIP(r)

			// Get the rate limiter for this IP
			ipLimiter := limiter.GetLimiter(ip)

			if ipLimiter == nil {
				// IP is banned
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			// Check if rate limit is exceeded
			if !ipLimiter.Allow() {
				// Ban IP for repeated violations
				consecutiveFailures := 0 // This could be stored and tracked per IP
				if consecutiveFailures > 5 {
					limiter.BanIP(ip)
				}

				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			// Check for suspicious headers or patterns that might indicate attacks
			if isSuspiciousRequest(r) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// Continue with the next handler if rate limit is not exceeded
			next.ServeHTTP(w, r)
		})
	}
}

// Additional security checks for common attacks
func isSuspiciousRequest(r *http.Request) bool {
	// Check for extremely large content length (potential DoS)
	if r.ContentLength > 10*1024*1024 { // 10MB limit
		return true
	}

	// Check for suspicious user agents
	userAgent := strings.ToLower(r.UserAgent())
	suspiciousAgents := []string{"nikto", "sqlmap", "nmap", "dirbuster", "nessus", "hydra"}
	for _, agent := range suspiciousAgents {
		if strings.Contains(userAgent, agent) {
			return true
		}
	}

	// Check for SQL injection attempts
	path := strings.ToLower(r.URL.Path)
	query := strings.ToLower(r.URL.RawQuery)
	sqlPatterns := []string{"union select", "order by", "group by", "1=1", "or 1=1", "--", ";--", "/*"}
	for _, pattern := range sqlPatterns {
		if strings.Contains(path, pattern) || strings.Contains(query, pattern) {
			return true
		}
	}

	// Check for XSS attempts
	xssPatterns := []string{"<script>", "javascript:", "onerror=", "onload=", "eval("}
	for _, pattern := range xssPatterns {
		if strings.Contains(path, pattern) || strings.Contains(query, pattern) {
			return true
		}
	}

	return false
}

// Usage example
func ExampleUsage() {
	// Create a new rate limiter allowing 6 requests per second with a burst of 10
	// and a ban duration of 1 hour for IPs that exceed the limit
	limiter := NewIPRateLimiter(6, 10, 1*time.Hour)

	// Create a simple handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Apply our middlewares
	secureHandler := SecurityHeaders(MaliciousRateLimitMiddleware(limiter)(handler))

	// Create and start the server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      secureHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	server.ListenAndServe()
}
