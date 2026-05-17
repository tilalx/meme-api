package config

import (
	"fmt"
	"os"
)

// required lists env vars that must be set for the app to start.
var required = []string{
	"REDDIT_CLIENT_ID",
	"REDDIT_CLIENT_SECRET",
}

// Validate checks that all required environment variables are set.
// Returns an error listing every missing variable so the operator sees
// all problems at once instead of one at a time.
func Validate() error {
	var missing []string
	for _, key := range required {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missing)
	}
	return nil
}
