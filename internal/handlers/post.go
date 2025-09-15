package handlers

import (
	"golang-final-test/internal/service"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	svc service.PostService
}

func NewPostHandler(svc service.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
}

func (h *PostHandler) SearchPostsByTag(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		c.JSON(400, gin.H{"error": "tag query parameter is required"})
		return
	}

	posts, err := h.svc.SearchPostsByTag(tag)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, posts)
}
