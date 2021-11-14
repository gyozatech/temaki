package middlewares

import (
	"net/http"
)

// CORSMiddleware function sets the headers for the response and returns directly if there is an OPTIONS request
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers for CORS Response
		w.Header().Set("access-control-allow-origin", "*")
		w.Header().Set("access-control-allow-credentials", "true")
		w.Header().Set("access-control-allow-methods", "GET, PUT, PATCH, POST, DELETE, HEAD, OPTIONS")
		w.Header().Set("access-control-allow-headers", "Authorization, Content-Type, session")
		w.Header().Set("access-control-expose-headers", "session")

		if r.Method == http.MethodOptions {
			return
		}
		// To next middleware
		next.ServeHTTP(w, r)
	})
}
