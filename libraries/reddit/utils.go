package reddit

import (
	"encoding/base64"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
)

// httpClient is a shared HTTP client with a timeout to prevent resource exhaustion.
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// EncodeCredentials : Return base64 Encoded client ID and Secret required for authentication
func EncodeCredentials() (encodedCredentials string) {
	data := ClientID + ":" + ClientSecret
	encodedCredentials = base64.StdEncoding.EncodeToString([]byte(data))
	return
}

// MakeGetRequest : Makes a GET Request to Reddit API with Access Token
func MakeGetRequest(url string) (responseBody []byte, errorCode int) {
	req, _ := http.NewRequest("GET", url, nil)

	tokenMu.RLock()
	token := AccessToken
	tokenMu.RUnlock()

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "oauth.reddit.com")
	req.Header.Add("Connection", "keep-alive")

	res, err := httpClient.Do(req)

	if err != nil {
		sentry.CaptureException(err)
		slog.Error("error making Reddit API request", "url", url, "error", err)
		return nil, http.StatusInternalServerError
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		sentry.CaptureException(err)
		slog.Error("error reading Reddit API response body", "error", err)
		return nil, res.StatusCode
	}

	return body, res.StatusCode
}

// validSorts is the set of Reddit sort modes accepted by the API.
var validSorts = map[string]bool{
	"hot":           true,
	"new":           true,
	"top":           true,
	"rising":        true,
	"controversial": true,
}

// GetSubredditAPIURL : Returns API Reddit URL for a subreddit with limit and sort order
func GetSubredditAPIURL(subreddit string, limit int, sort string) (url string) {
	if !validSorts[sort] {
		sort = "hot"
	}
	url = "https://oauth.reddit.com/r/" + subreddit + "/" + sort + "?limit=" + strconv.Itoa(limit)
	return
}
