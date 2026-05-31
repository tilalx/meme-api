package redis

import (
	"context"
	"encoding/json"

	"Meme_Api/models"
)

// GetPostsFromCache : Get memes from Cache based on sub
func GetPostsFromCache(sub string) []models.Meme {

	if Client == nil {
		return nil
	}

	ctx := context.Background()
	memesBinary, err := Client.Get(ctx, sub).Bytes()

	if err != nil {
		return nil
	}

	var memes []models.Meme

	err = json.Unmarshal(memesBinary, &memes)

	if err != nil {
		return nil
	}

	return memes
}
