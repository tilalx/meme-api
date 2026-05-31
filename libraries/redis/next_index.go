package redis

import (
	"context"

	"Meme_Api/data"
)

const indexKeySuffix = ":idx"

// NextIndex atomically increments and returns the next serving index for the
// given subreddit. The index key shares the cache expiration so it resets
// naturally when the post list expires.
func NextIndex(sub string) int64 {
	if Client == nil {
		return 0
	}

	ctx := context.Background()
	key := sub + indexKeySuffix
	idx, err := Client.Incr(ctx, key).Result()
	if err != nil {
		return 0
	}

	// Keep the index key alive for as long as the post cache.
	Client.Expire(ctx, key, data.CacheExpirationTime)

	return idx - 1 // return 0-based index before increment
}

// NextIndexBy atomically advances the counter by n and returns the 0-based
// starting position. Used for multi-meme requests.
func NextIndexBy(sub string, n int) int64 {
	if Client == nil {
		return 0
	}

	ctx := context.Background()
	key := sub + indexKeySuffix
	result, err := Client.IncrBy(ctx, key, int64(n)).Result()
	if err != nil {
		return 0
	}

	Client.Expire(ctx, key, data.CacheExpirationTime)

	return result - int64(n) // 0-based start of this batch
}
