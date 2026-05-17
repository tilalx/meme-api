package reddit

import (
	"log/slog"
	"os"
	"sync"
)

// Reddit : Reddit structure with reddit credentials
var (
	AccessToken  string
	ClientID     string
	ClientSecret string
	UserAgent    string

	tokenMu sync.RWMutex
)

// Init : Initialize the Reddit Structure with App Credentials
func Init() {
	// Get Reddit Client Credentials from the environment variables
	clientID := os.Getenv("REDDIT_CLIENT_ID")
	clientSecret := os.Getenv("REDDIT_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		slog.Error("REDDIT_CLIENT_ID and REDDIT_CLIENT_SECRET must be set")
		os.Exit(1)
	}

	ClientID = clientID
	ClientSecret = clientSecret

	UserAgent = "MEME_API"

	accessToken := GetAccessToken()

	if accessToken == "" {
		slog.Error("failed to obtain Reddit access token at startup")
		os.Exit(1)
	}

	tokenMu.Lock()
	AccessToken = accessToken
	tokenMu.Unlock()
}

// GetNewAccessToken : Function to Generate New Access Token once the old one expires
func GetNewAccessToken() (ok bool) {
	newAccessToken := GetAccessToken()

	if newAccessToken == "" {
		slog.Warn("unable to refresh Reddit access token")
		return false
	}

	tokenMu.Lock()
	AccessToken = newAccessToken
	tokenMu.Unlock()
	return true
}
