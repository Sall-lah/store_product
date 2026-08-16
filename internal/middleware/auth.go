package middleware

import (
	"context"
	"encoding/json"
	"net/http"
)

type contextKey string

const (
	// UserCtxKey is the key used to store authenticated gateway user metadata in context.
	UserCtxKey contextKey = "gateway_user"
)

// GatewayUser holds identity attributes forwarded by the API Gateway.
type GatewayUser struct {
	ID   string
	Role string
}

// RequireAdmin verifies that the incoming request was authorized by the API Gateway
// with an admin role. In a microservices architecture, identity verification is offloaded
// to the gateway to maintain loose coupling across backend services.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("X-User-Role")
		if role != "admin" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error":   "forbidden",
				"message": "Admin privileges are required to perform this action.",
			})
			return
		}

		user := GatewayUser{
			ID:   r.Header.Get("X-User-Id"),
			Role: role,
		}

		ctx := context.WithValue(r.Context(), UserCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext retrieves the forwarded gateway user from the request context.
func GetUserFromContext(ctx context.Context) (GatewayUser, bool) {
	user, ok := ctx.Value(UserCtxKey).(GatewayUser)
	return user, ok
}
