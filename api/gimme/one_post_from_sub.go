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

// GetOnePostFromSub returns a single meme from a specific subreddit.
//
// @Summary      Get a random meme from a subreddit
// @Description  Returns a single random meme from the specified subreddit.
// @Tags         memes
// @Produce      json
// @Param        subreddit  path     string  true   "Subreddit name (alphanumeric + underscores, max 50 chars)"
// @Param        sort       query    string  false  "Reddit sort method"  Enums(hot, new, top)  default(hot)
// @Param        nsfw       query    boolean false  "Include NSFW posts"  default(true)
// @Param        spoiler    query    boolean false  "Include spoiler posts" default(true)
// @Success      200  {object}  response.OneMeme
// @Failure      400  {object}  response.Error
// @Failure      404  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /gimme/{subreddit} [get]
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
		freshMemes, res := reddit.GetNPosts(c.Request.Context(), sub, data.RedditPostsLimit, sort)

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
