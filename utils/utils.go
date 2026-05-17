package utils

import (
	"math/rand"
	"strings"

	"Meme_Api/models"
)

// GetRandomN : Get a Random Number in the range [0, n)
func GetRandomN(n int) int {
	return rand.Intn(n)
}

// RemoveNonImagePosts : Remove all posts from Memes List that don't have an image URL
func RemoveNonImagePosts(memes []models.Meme) []models.Meme {
	var onlyImagePosts []models.Meme

	for _, meme := range memes {
		url := strings.ToLower(meme.URL)
		if strings.HasSuffix(url, ".jpg") || strings.HasSuffix(url, ".png") || strings.HasSuffix(url, ".gif") {
			onlyImagePosts = append(onlyImagePosts, meme)
		}
	}

	return onlyImagePosts
}

// GetNRandomMemes : Get N no. of random memes from a list of Memes
func GetNRandomMemes(memes []models.Meme, n int) []models.Meme {
	rand.Shuffle(len(memes), func(i, j int) { memes[i], memes[j] = memes[j], memes[i] })
	return memes[:n]
}
