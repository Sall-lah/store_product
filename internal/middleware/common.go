package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/cors"
)

// responseWriterWrapper intercepts the HTTP status code written by downstream handlers
// so it can be recorded accurately in structured access logs.
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger returns middleware that emits structured HTTP access logs.
// Tracking latency and status codes is essential for microservice observability.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapper := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapper, r)

		duration := time.Since(start)
		log.Printf("[HTTP] %s %s | Status: %d | Duration: %v | Client: %s",
			r.Method,
			r.URL.Path,
			wrapper.statusCode,
			duration,
			r.RemoteAddr,
		)
	})
}

// Recovery catches unhandled panics within request lifecycles, logging the stack trace
// and returning an RFC-compliant 500 JSON response to avoid process termination.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("[PANIC] Recovered from panic in %s %s: %v\nStack: %s",
					r.Method, r.URL.Path, rec, string(debug.Stack()))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error":   "internal_server_error",
					"message": "An unexpected server error occurred.",
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// SetupCORS returns a standard CORS middleware allowing communication across microservices and frontend clients.
func SetupCORS() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-User-Id", "X-User-Role", "X-Forwarded-For"},
		ExposedHeaders:   []string{"Link", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset", "Retry-After"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}
