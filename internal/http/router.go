package httpserver

import (
	"golang-final-test/internal/handlers"
	"golang-final-test/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouterDependencies struct {
	PostService service.PostService
}

func NewRouter(deps RouterDependencies) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

	posts := r.Group("/posts")
	{
		h := handlers.NewPostHandler(deps.PostService)
		posts.POST("/", h.CreatePost)
		posts.GET("/:id", h.GetPostByID)
		posts.GET("/search-by-tag", h.SearchPostsByTag)
		posts.PUT("/", h.UpdatePost)
		posts.GET("/search", h.SearchPosts)
	}

	return r
}
