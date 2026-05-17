package gimme

import (
	"net/http"
	"strconv"

	"Meme_Api/data"
	"Meme_Api/libraries/reddit"
	"Meme_Api/libraries/redis"
	"Meme_Api/models/response"
	"Meme_Api/utils"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

// GetNRandomMemes returns N memes from a random default subreddit.
//
// @Summary      Get N random memes
// @Description  Returns up to 50 random memes from one of the default subreddits.
// @Tags         memes
// @Produce      json
// @Param        count    path     int     true   "Number of memes (1–50)"
// @Param        sort     query    string  false  "Reddit sort method"  Enums(hot, new, top)  default(hot)
// @Param        nsfw     query    boolean false  "Include NSFW posts"  default(true)
// @Param        spoiler  query    boolean false  "Include spoiler posts" default(true)
// @Success      200  {object}  response.MultipleMemes
// @Failure      400  {object}  response.Error
// @Failure      404  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /gimme/{count} [get]
func GetNRandomMemes(c *gin.Context) {

	count, _ := strconv.Atoi(c.Param("interface"))
	sort := c.DefaultQuery("sort", "hot")

    if count <= 0 {
        c.JSON(http.StatusBadRequest, response.Error{
            Code: http.StatusBadRequest,
            Message: "Invalid Count Value",
        })
        return
    }

	if count > 50 {
		count = 50
	}

	sub := data.MemeSubreddits[utils.GetRandomN(len(data.MemeSubreddits))]

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
			Message: "No memes found matching the requested filters",
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
