package reddit

import (
	"encoding/base64"
	"io"
	"log"
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

	req.Header.Add("Authorization", "Bearer "+AccessToken)
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "oauth.reddit.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, err := httpClient.Do(req)

	if err != nil {
		sentry.CaptureException(err)
		log.Println("Error while making request", err)
		return nil, http.StatusInternalServerError
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		sentry.CaptureException(err)
		log.Println("Error while Parsing Response Body")
		return nil, res.StatusCode
	}

	return body, res.StatusCode
}

// GetSubredditAPIURL : Returns API Reddit URL with Limit
func GetSubredditAPIURL(subreddit string, limit int) (url string) {
	url = "https://oauth.reddit.com/r/" + subreddit + "/hot?limit=" + strconv.Itoa(limit)
	return
}
