package handlers

import (
	"golang-final-test/internal/models"
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
	var post struct {
		Title   string   `json:"title" binding:"required"`
		Content string   `json:"content" binding:"required"`
		Tags    []string `json:"tags" binding:"required"`
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.CreatePost(&models.Post{
		Title:   post.Title,
		Content: post.Content,
		Tags:    post.Tags,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Post created successfully"})
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "id parameter is required"})
		return
	}

	post, err := h.svc.GetPostByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if post == nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(200, post)
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
