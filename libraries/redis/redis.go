package redis

import (
	"context"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
)

// Redis : Redis structure with client
var (
	RedisURL string
	Client   *redis.Client
)

// Init : Initialize New Redis Connection (optional — app runs without caching if URL is unset)
func Init() {
	redisURL := os.Getenv("REDISCLOUD_URL")

	if redisURL == "" {
		slog.Warn("REDISCLOUD_URL not set — running without Redis cache")
		return
	}

	options, err := redis.ParseURL(redisURL)

	if err != nil {
		slog.Warn("invalid Redis URL, skipping Redis", "error", err)
		return
	}

	redisClient := redis.NewClient(options)

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		slog.Warn("unable to reach Redis — running without cache", "error", err)
		return
	}

	Client = redisClient
	slog.Info("Redis connected successfully")
}
