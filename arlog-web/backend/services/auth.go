package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService handles OAuth2 authentication with Okta
type AuthService struct {
	clientID     string
	clientSecret string
	redirectURI  string
	oktaDomain   string
	issuer       string
}

// NewAuthService creates a new authentication service
func NewAuthService() *AuthService {
	oktaDomain := os.Getenv("OKTA_DOMAIN")
	return &AuthService{
		clientID:     os.Getenv("OKTA_CLIENT_ID"),
		clientSecret: os.Getenv("OKTA_CLIENT_SECRET"),
		redirectURI:  os.Getenv("OKTA_REDIRECT_URI"),
		oktaDomain:   oktaDomain,
		issuer:       fmt.Sprintf("https://%s/oauth2/default", oktaDomain),
	}
}

// GetAuthorizationURL generates the Okta authorization URL
func (a *AuthService) GetAuthorizationURL(state string) string {
	authURL := fmt.Sprintf("https://%s/oauth2/default/v1/authorize", a.oktaDomain)
	
	params := url.Values{}
	params.Add("client_id", a.clientID)
	params.Add("response_type", "code")
	params.Add("scope", "openid profile email groups")
	params.Add("redirect_uri", a.redirectURI)
	params.Add("state", state)
	
	return fmt.Sprintf("%s?%s", authURL, params.Encode())
}

// ExchangeCodeForToken exchanges an authorization code for an access token
func (a *AuthService) ExchangeCodeForToken(code string) (*TokenResponse, error) {
	tokenURL := fmt.Sprintf("https://%s/oauth2/default/v1/token", a.oktaDomain)
	
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", a.redirectURI)
	data.Set("client_id", a.clientID)
	data.Set("client_secret", a.clientSecret)
	
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}
	
	return &tokenResponse, nil
}

// GetUserInfo retrieves user information from Okta
func (a *AuthService) GetUserInfo(accessToken string) (*UserInfo, error) {
	userInfoURL := fmt.Sprintf("https://%s/oauth2/default/v1/userinfo", a.oktaDomain)
	
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+accessToken)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("userinfo request failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}
	
	return &userInfo, nil
}

// CreateSessionToken creates a JWT session token for the user
func (a *AuthService) CreateSessionToken(userInfo *UserInfo) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET not configured")
	}
	
	claims := jwt.MapClaims{
		"sub":        userInfo.Sub,
		"email":      userInfo.Email,
		"name":       userInfo.Name,
		"groups":     userInfo.Groups,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
		"iat":        time.Now().Unix(),
		"okta_token": userInfo.OktaUserID,
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	
	return signedToken, nil
}

// GenerateStateToken generates a random state token for OAuth2 flow
func GenerateStateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// TokenResponse represents the Okta token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token"`
	Scope        string `json:"scope"`
}

// UserInfo represents user information from Okta
type UserInfo struct {
	Sub        string   `json:"sub"`
	Email      string   `json:"email"`
	Name       string   `json:"name"`
	Groups     []string `json:"groups"`
	OktaUserID string   `json:"okta_user_id"`
}

