package server

import (
	"net/http"

	"Meme_Api/api/gimme"
	"Meme_Api/server/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter : Function with routes
func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
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

	// Swagger UI (assets embedded, no CDN)
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// /gimme routes
	gimmeRouter := router.Group("gimme")
	gimme.Routes(gimmeRouter)

	return router
}
