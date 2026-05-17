package gimme

import (
	"fmt"
	"net/http"
	"strings"

	"Meme_Api/data"
	"Meme_Api/libraries/reddit"
	"Meme_Api/libraries/redis"

	"Meme_Api/models/response"
	"Meme_Api/utils"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

// GetOnePostFromSub : Get one post from a specific subreddit
func GetOnePostFromSub(c *gin.Context) {

	sub := strings.ToLower(c.Param("interface"))
	sort := c.DefaultQuery("sort", "hot")

	if !isValidSubreddit(sub) {
		c.JSON(http.StatusBadRequest, response.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid subreddit name",
		})
		return
	}

	memes := redis.GetPostsFromCache(sub)

	if memes == nil {
		freshMemes, res := reddit.GetNPosts(sub, data.RedditPostsLimit, sort)

		if freshMemes == nil {
			c.JSON(res.Code, res)
			return
		}

		freshMemes = utils.RemoveNonImagePosts(freshMemes)

		if err := redis.WritePostsToCache(sub, freshMemes); err != nil {
			sentry.CaptureException(err)
		}

		memes = freshMemes
	}

	memes = utils.ApplyFilters(memes, c)

	if len(memes) == 0 {
		c.JSON(http.StatusNotFound, response.Error{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("r/%s has no posts matching the requested filters", sub),
		})
		return
	}

	idx := redis.NextIndex(sub)
	meme := utils.PickMemeAt(memes, idx)

	c.JSON(http.StatusOK, response.OneMeme{
		PostLink:  meme.PostLink,
		Subreddit: meme.SubReddit,
		Title:     meme.Title,
		URL:       meme.URL,
		NSFW:      meme.NSFW,
		Spoiler:   meme.Spoiler,
		Author:    meme.Author,
		Ups:       meme.Ups,
		Preview:   meme.Preview,
	})
}
