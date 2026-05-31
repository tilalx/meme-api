package redis

import (
	"context"
	"encoding/json"
	"math/rand"

	"Meme_Api/data"
	"Meme_Api/models"
)

// WritePostsToCache shuffles memes randomly before storing so the round-robin
// serving order is different on every cache refresh.
func WritePostsToCache(sub string, memes []models.Meme) error {
	if Client == nil {
		return nil
	}

	// Shuffle a copy so callers aren't surprised by mutation
	shuffled := make([]models.Meme, len(memes))
	copy(shuffled, memes)
	rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

	memesBinary, err := json.Marshal(shuffled)
	if err != nil {
		return err
	}

	ctx := context.Background()
	// Reset the index counter so round-robin starts from 0 on the new shuffled list
	Client.Del(ctx, sub+indexKeySuffix)

	return Client.Set(ctx, sub, memesBinary, data.CacheExpirationTime).Err()
}
