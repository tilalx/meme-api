package gimme

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"Meme_Api/data"
	"Meme_Api/libraries/reddit"
	"Meme_Api/libraries/redis"

	"Meme_Api/models/response"
	"Meme_Api/utils"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

// GetNPostsFromSub : Get N no. of posts from a specific subreddit
func GetNPostsFromSub(c *gin.Context) {

	sub := strings.ToLower(c.Param("interface"))
	count, err := strconv.Atoi(c.Param("count"))
	sort := c.DefaultQuery("sort", "hot")

	if err != nil || count <= 0 {
		c.JSON(http.StatusBadRequest, response.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid Count Value",
		})
		return
	}

	if !isValidSubreddit(sub) {
		c.JSON(http.StatusBadRequest, response.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid subreddit name",
		})
		return
	}

	if count > 50 {
		count = 50
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

	memesLen := len(memes)

	if memesLen == 0 {
		c.JSON(http.StatusNotFound, response.Error{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("r/%s has no posts matching the requested filters", sub),
		})
		return
	}

	if memesLen < count {
		count = memesLen
	}

	start := redis.NextIndexBy(sub, count)
	memes = utils.PickNMemesFrom(memes, start, count)

	var memesResponse []response.OneMeme

	for _, meme := range memes {
		memesResponse = append(memesResponse, response.OneMeme{
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

	c.JSON(http.StatusOK, response.MultipleMemes{
		Count: len(memesResponse),
		Memes: memesResponse,
	})
}
