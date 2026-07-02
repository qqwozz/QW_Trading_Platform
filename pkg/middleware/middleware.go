// Package middleware provides HTTP middleware functions for CORS, JWT authentication,
// and request logging.
package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/qw-trading/platform/pkg/logger"
)

// contextKey is a private type for context keys to avoid collisions.
type contextKey string

const (
	// UserIDKey is the context key used to store the authenticated user's ID.
	UserIDKey contextKey = "user_id"
	// RequestIDKey is the context key used to store the request ID.
	RequestIDKey contextKey = "request_id"
)

// CORS returns middleware that sets Cross-Origin Resource Sharing headers
// and handles preflight OPTIONS requests.
func CORS(allowedOrigins string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID")
			w.Header().Set("Access-Control-Max-Age", "86400")

			// Short-circuit preflight requests without calling downstream handlers.
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Auth returns middleware that validates JWT tokens from the Authorization header.
// Requests to public paths (register, login, health, docs) bypass authentication.
func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Allow unauthenticated access to public endpoints.
			publicPaths := []string{"/v1/auth/register", "/v1/auth/login", "/v1/auth/guest", "/health", "/docs"}
			for _, path := range publicPaths {
				if strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(w, r)
					return
				}
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"authorization header required"}`, http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			userID, err := validateToken(tokenStr, jwtSecret)
			if err != nil {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}

			// Inject the authenticated user ID into the request context.
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequestID returns middleware that generates a unique request ID and injects
// it into the request context and response header.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get("X-Request-ID")
		if rid == "" {
			rid = uuid.New().String()
		}
		w.Header().Set("X-Request-ID", rid)
		ctx := context.WithValue(r.Context(), RequestIDKey, rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logger returns middleware that logs each request using the structured logger.
func Logger(l *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			rid, _ := r.Context().Value(RequestIDKey).(string)
			l.Request(r.Method, r.URL.Path, wrapped.statusCode, time.Since(start), rid)
		})
	}
}

// GetUserID extracts the authenticated user ID from the request context.
// Returns the UUID and true if present and valid, or uuid.Nil and false otherwise.
func GetUserID(r *http.Request) (uuid.UUID, bool) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok || userID == "" {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

// validateToken parses and validates a JWT token, returning the user_id claim.
func validateToken(tokenStr, secret string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token was signed with HMAC (HS256/HS384/HS512).
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", jwt.ErrTokenInvalidClaims
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", jwt.ErrTokenInvalidClaims
	}

	return userID, nil
}

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before writing it to the underlying writer.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
