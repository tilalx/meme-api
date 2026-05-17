package redis

import "Meme_Api/data"

const indexKeySuffix = ":idx"

// NextIndex atomically increments and returns the next serving index for the
// given subreddit. The index key shares the cache expiration so it resets
// naturally when the post list expires.
func NextIndex(sub string) int64 {
	if Client == nil {
		return 0
	}

	key := sub + indexKeySuffix
	idx, err := Client.Incr(key).Result()
	if err != nil {
		return 0
	}

	// Keep the index key alive for as long as the post cache.
	Client.Expire(key, data.CacheExpirationTime)

	return idx - 1 // return 0-based index before increment
}

// NextIndexBy atomically advances the counter by n and returns the 0-based
// starting position. Used for multi-meme requests.
func NextIndexBy(sub string, n int) int64 {
	if Client == nil {
		return 0
	}

	key := sub + indexKeySuffix
	result, err := Client.IncrBy(key, int64(n)).Result()
	if err != nil {
		return 0
	}

	Client.Expire(key, data.CacheExpirationTime)

	return result - int64(n) // 0-based start of this batch
}
