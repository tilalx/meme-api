package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Endpoint describes a single API route for documentation purposes.
type Endpoint struct {
	Path        string  `json:"path"`
	Method      string  `json:"method"`
	Description string  `json:"description"`
	Params      []Param `json:"params,omitempty"`
}

// Param describes a URL parameter.
type Param struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// DocsResponse is the JSON shape returned at GET /.
type DocsResponse struct {
	Title       string     `json:"title"`
	Version     string     `json:"version"`
	Description string     `json:"description"`
	BaseURL     string     `json:"baseUrl"`
	Endpoints   []Endpoint `json:"endpoints"`
}

var docs = DocsResponse{
	Title:       "Meme API",
	Version:     "1.0.0",
	Description: "A simple, free REST API that serves random memes scraped from Reddit.",
	BaseURL:     "https://meme-api.com",
	Endpoints: []Endpoint{
		{
			Path:        "/gimme",
			Method:      "GET",
			Description: "Returns a single random meme from a random subreddit.",
		},
		{
			Path:        "/gimme/:count",
			Method:      "GET",
			Description: "Returns N random memes from a random subreddit (max 50).",
			Params: []Param{
				{Name: "count", Type: "integer", Description: "Number of memes to return (1–50)."},
			},
		},
		{
			Path:        "/gimme/:subreddit",
			Method:      "GET",
			Description: "Returns a single random meme from the specified subreddit.",
			Params: []Param{
				{Name: "subreddit", Type: "string", Description: "Name of a Reddit subreddit (alphanumeric and underscores, max 50 chars)."},
			},
		},
		{
			Path:        "/gimme/:subreddit/:count",
			Method:      "GET",
			Description: "Returns N random memes from the specified subreddit (max 50).",
			Params: []Param{
				{Name: "subreddit", Type: "string", Description: "Name of a Reddit subreddit (alphanumeric and underscores, max 50 chars)."},
				{Name: "count", Type: "integer", Description: "Number of memes to return (1–50)."},
			},
		},
	},
}

// GetDocs returns the API documentation as JSON.
func GetDocs(c *gin.Context) {
	c.JSON(http.StatusOK, docs)
}
