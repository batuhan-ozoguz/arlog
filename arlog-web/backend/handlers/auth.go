package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"arlog/backend/services"
)

var authService *services.AuthService

func init() {
	authService = services.NewAuthService()
}

// OktaLogin initiates the Okta authentication flow
func OktaLogin(w http.ResponseWriter, r *http.Request) {
	// Check if authentication is disabled (dev mode)
	authMode := os.Getenv("AUTH_MODE")
	if authMode == "dev" {
		// In dev mode, create a dummy token and redirect
		dummyToken, err := createDevToken()
		if err != nil {
			log.Printf("Error creating dev token: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:5173"
		}

		redirectURL := fmt.Sprintf("%s/auth/callback?token=%s", frontendURL, dummyToken)
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
		return
	}

	// Generate state token for CSRF protection
	state, err := services.GenerateStateToken()
	if err != nil {
		log.Printf("Error generating state token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Store state in session/cookie (for MVP, we'll pass it through the flow)
	// In production, you should store this in a secure session
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		MaxAge:   600, // 10 minutes
	})

	// Get authorization URL and redirect
	authURL := authService.GetAuthorizationURL(state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// OktaCallback handles the Okta OAuth2 callback
func OktaCallback(w http.ResponseWriter, r *http.Request) {
	// Get state from query parameters
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	errorParam := r.URL.Query().Get("error")

	// Check for errors from Okta
	if errorParam != "" {
		errorDescription := r.URL.Query().Get("error_description")
		log.Printf("Okta authentication error: %s - %s", errorParam, errorDescription)
		http.Error(w, fmt.Sprintf("Authentication failed: %s", errorDescription), http.StatusUnauthorized)
		return
	}

	// Verify state token
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != state {
		log.Printf("Invalid state token")
		http.Error(w, "Invalid state token", http.StatusBadRequest)
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "oauth_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// Exchange code for token
	tokenResponse, err := authService.ExchangeCodeForToken(code)
	if err != nil {
		log.Printf("Error exchanging code for token: %v", err)
		http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
		return
	}

	// Get user information
	userInfo, err := authService.GetUserInfo(tokenResponse.AccessToken)
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		http.Error(w, "Failed to get user information", http.StatusInternalServerError)
		return
	}

	// Create session token
	sessionToken, err := authService.CreateSessionToken(userInfo)
	if err != nil {
		log.Printf("Error creating session token: %v", err)
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// For MVP, we'll return the token in JSON
	// In production, you might want to set it as an HTTP-only cookie and redirect to the frontend
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	// Redirect to frontend with token
	redirectURL := fmt.Sprintf("%s/auth/callback?token=%s", frontendURL, sessionToken)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// Logout handles user logout
func Logout(w http.ResponseWriter, r *http.Request) {
	// Clear any cookies
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}

// createDevToken creates a dummy JWT token for development mode
func createDevToken() (string, error) {
	// Create a dummy user for development
	dummyUserInfo := &services.UserInfo{
		Sub:        "dev-user-123",
		Email:      "dev@example.com",
		Name:       "Development User",
		Groups:     []string{"cosmos-team-okta-group"}, // Match seeded test data
		OktaUserID: "dev-okta-id",
	}

	return authService.CreateSessionToken(dummyUserInfo)
}

