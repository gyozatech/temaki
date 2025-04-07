package middlewares

import (
	"net/http"
	"time"

	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/middleware/stdlib"
	"github.com/ulule/limiter/drivers/store/memory"
)

// SimpleRateLimitMiddleware is a middleware that limits the number of requests per IP address.
func SimpleRateLimitMiddleware(next http.Handler) http.Handler {
	// Create a new limiter store using in-memory storage.
	rate := limiter.Rate{
		Period: 1 * time.Second,
		Limit:  6, // 6 requests per second.
	}
	store := memory.NewStore()
	limitMiddleware := stdlib.NewMiddleware(limiter.New(store, rate))

	return limitMiddleware.Handler(next)
}
