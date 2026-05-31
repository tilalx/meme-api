package gimme

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"

	"Meme_Api/data"
	"Meme_Api/libraries/reddit"
	"Meme_Api/libraries/redis"
	"Meme_Api/models/response"
	"Meme_Api/utils"
)

// GetOneRandomMeme returns a single meme from a random default subreddit.
//
// @Summary      Get a random meme
// @Description  Returns a single random meme from one of the default subreddits (memes, dankmemes, me_irl).
// @Tags         memes
// @Produce      json
// @Param        sort     query    string  false  "Reddit sort method"  Enums(hot, new, top)  default(hot)
// @Param        nsfw     query    boolean false  "Include NSFW posts"  default(true)
// @Param        spoiler  query    boolean false  "Include spoiler posts" default(true)
// @Success      200  {object}  response.OneMeme
// @Failure      404  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /gimme [get]
func GetOneRandomMeme(c *gin.Context) {

	sort := c.DefaultQuery("sort", "hot")

	// Choose Random Meme Subreddit
	sub := data.MemeSubreddits[utils.GetRandomN(len(data.MemeSubreddits))]

	// Check if the sub is present in the cache
	memes := redis.GetPostsFromCache(sub)

	// If it is not in Cache then get posts from Reddit
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
			Message: "No memes found matching the requested filters",
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
