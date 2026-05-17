package reddit

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	rm "Meme_Api/libraries/reddit/models"

	"github.com/getsentry/sentry-go"
)

// GetAccessToken : Get temporary Access Token based on App client ID and Secret
func GetAccessToken() (accessToken string) {

	encodedCredentials := EncodeCredentials()

	// Reddit URL to get access token
	url := "https://www.reddit.com/api/v1/access_token"

	// Set grant type to client_credentials as POST body
	payload := strings.NewReader("grant_type=client_credentials")

	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		sentry.CaptureException(err)
		slog.Error("error creating Reddit auth request", "error", err)
		return ""
	}

	// Set Headers including the User Agent and the Authorization with the encoded credentials
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("Authorization", "Basic "+encodedCredentials)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "www.reddit.com")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Connection", "keep-alive")

	// Make Request
	res, err := httpClient.Do(req)

	if err != nil {
		sentry.CaptureException(err)
		slog.Error("error connecting to Reddit auth endpoint", "url", url, "error", err)
		return ""
	}

	// Close the response body
	defer res.Body.Close()

	// Read the response
	body, err := io.ReadAll(res.Body)

	if err != nil {
		sentry.CaptureException(err)
		slog.Error("error reading Reddit auth response body", "error", err)
		return ""
	}

	var accessTokenBody rm.AccessTokenBody

	if err := json.Unmarshal(body, &accessTokenBody); err != nil {
		slog.Error("error parsing Reddit access token response", "error", err)
		sentry.CaptureException(err)
		return ""
	}

	return accessTokenBody.AccessToken
}
