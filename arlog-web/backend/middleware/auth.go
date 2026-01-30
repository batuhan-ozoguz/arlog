package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// UserContext key type for storing user information in request context
type contextKey string

const UserContextKey contextKey = "user"

// UserInfo represents authenticated user information
type UserInfo struct {
	Sub        string   `json:"sub"`
	Email      string   `json:"email"`
	Name       string   `json:"name"`
	Groups     []string `json:"groups"`
	OktaUserID string   `json:"okta_user_id"`
}

// AuthMiddleware validates JWT tokens and extracts user information
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if authentication is disabled (dev mode)
		authMode := os.Getenv("AUTH_MODE")
		if authMode == "dev" {
			// In dev mode, create a dummy user and skip authentication
			userInfo := &UserInfo{
				Sub:        "dev-user-123",
				Email:      "dev@example.com",
				Name:       "Development User",
				Groups:     []string{"cosmos-team-okta-group"},
				OktaUserID: "dev-okta-id",
			}
			ctx := context.WithValue(r.Context(), UserContextKey, userInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		// Extract Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondWithError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		tokenString := parts[1]

		// Parse and validate JWT token
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			respondWithError(w, http.StatusInternalServerError, "JWT secret not configured")
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Extract user information from claims
		userInfo := &UserInfo{
			Sub:        getStringClaim(claims, "sub"),
			Email:      getStringClaim(claims, "email"),
			Name:       getStringClaim(claims, "name"),
			Groups:     getStringArrayClaim(claims, "groups"),
			OktaUserID: getStringClaim(claims, "okta_token"),
		}

		// Add user info to request context
		ctx := context.WithValue(r.Context(), UserContextKey, userInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuthMiddleware is similar to AuthMiddleware but doesn't require authentication
// It adds user info to context if a valid token is provided
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString := parts[1]
				jwtSecret := os.Getenv("JWT_SECRET")
				
				claims := jwt.MapClaims{}
				token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(jwtSecret), nil
				})

				if err == nil && token.Valid {
					userInfo := &UserInfo{
						Sub:        getStringClaim(claims, "sub"),
						Email:      getStringClaim(claims, "email"),
						Name:       getStringClaim(claims, "name"),
						Groups:     getStringArrayClaim(claims, "groups"),
						OktaUserID: getStringClaim(claims, "okta_token"),
					}
					ctx := context.WithValue(r.Context(), UserContextKey, userInfo)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
		}

		// No valid token, continue without user context
		next.ServeHTTP(w, r)
	})
}

// GetUserFromContext retrieves user information from the request context
func GetUserFromContext(ctx context.Context) (*UserInfo, bool) {
	user, ok := ctx.Value(UserContextKey).(*UserInfo)
	return user, ok
}

// Helper functions to extract claims
func getStringClaim(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key].(string); ok {
		return val
	}
	return ""
}

func getStringArrayClaim(claims jwt.MapClaims, key string) []string {
	if val, ok := claims[key].([]interface{}); ok {
		result := make([]string, 0, len(val))
		for _, v := range val {
			if str, ok := v.(string); ok {
				result = append(result, str)
			}
		}
		return result
	}
	return []string{}
}

// respondWithError sends a JSON error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

