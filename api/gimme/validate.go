package gimme

import "regexp"

// validSubreddit matches Reddit subreddit names: 1–50 alphanumeric chars or underscores.
var validSubreddit = regexp.MustCompile(`^[a-zA-Z0-9_]{1,50}$`)

// isValidSubreddit returns true if the given name is a safe subreddit identifier.
func isValidSubreddit(name string) bool {
	return validSubreddit.MatchString(name)
}
