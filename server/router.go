package server

import (
	"net/http"

	"Meme_Api/api"
	"Meme_Api/api/gimme"
	"Meme_Api/server/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewRouter : Function with routes
func NewRouter() *gin.Engine {
	router := gin.New()

	// Gin and CORS Middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	router.Use(middleware.RequestID())
	router.Use(middleware.RateLimit())

	// Health check
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API docs
	router.GET("/", api.GetDocs)

	// /gimme routes
	gimmeRouter := router.Group("gimme")
	gimme.Routes(gimmeRouter)

	return router
}
