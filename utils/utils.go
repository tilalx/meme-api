package utils

import (
	"math/rand"
	"strings"

	"Meme_Api/models"

	"github.com/gin-gonic/gin"
)

// GetRandomN : Get a Random Number in the range [0, n)
func GetRandomN(n int) int {
	return rand.Intn(n)
}

// PickMemeAt returns the meme at position (index mod len), safe for any index value.
func PickMemeAt(memes []models.Meme, index int64) models.Meme {
	return memes[int(index)%len(memes)]
}

// PickNMemesFrom returns n memes starting at index (round-robin, with wraparound).
func PickNMemesFrom(memes []models.Meme, start int64, n int) []models.Meme {
	total := len(memes)
	result := make([]models.Meme, n)
	for i := 0; i < n; i++ {
		result[i] = memes[(int(start)+i)%total]
	}
	return result
}

// isImageURL returns true if the URL points to a supported image format
// or is hosted on Reddit's image CDN (i.redd.it).
func isImageURL(rawURL string) bool {
	lower := strings.ToLower(rawURL)

	// Reddit image CDN always serves images
	if strings.HasPrefix(lower, "https://i.redd.it/") {
		return true
	}

	// Strip query string before checking extension
	path := lower
	if idx := strings.Index(lower, "?"); idx != -1 {
		path = lower[:idx]
	}

	return strings.HasSuffix(path, ".jpg") ||
		strings.HasSuffix(path, ".jpeg") ||
		strings.HasSuffix(path, ".png") ||
		strings.HasSuffix(path, ".gif") ||
		strings.HasSuffix(path, ".webp")
}

// RemoveNonImagePosts removes all posts that don't have an image URL.
func RemoveNonImagePosts(memes []models.Meme) []models.Meme {
	var onlyImagePosts []models.Meme

	for _, meme := range memes {
		if isImageURL(meme.URL) {
			onlyImagePosts = append(onlyImagePosts, meme)
		}
	}

	return onlyImagePosts
}

// ApplyFilters applies optional ?nsfw and ?spoiler query param filters.
// By default NSFW and spoiler posts are included; pass nsfw=false or spoiler=false to exclude them.
func ApplyFilters(memes []models.Meme, c *gin.Context) []models.Meme {
	allowNSFW := c.DefaultQuery("nsfw", "true") != "false"
	allowSpoiler := c.DefaultQuery("spoiler", "true") != "false"

	if allowNSFW && allowSpoiler {
		return memes
	}

	var filtered []models.Meme
	for _, m := range memes {
		if !allowNSFW && m.NSFW {
			continue
		}
		if !allowSpoiler && m.Spoiler {
			continue
		}
		filtered = append(filtered, m)
	}
	return filtered
}