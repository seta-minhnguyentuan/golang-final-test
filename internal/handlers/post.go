package handlers

import (
	"golang-final-test/internal/database"
	"golang-final-test/internal/models"
	"golang-final-test/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	err := database.WithTransaction(func(tx *gorm.DB) error {
		post := models.Post{Title: post.Title, Content: post.Content, Tags: post.Tags}
		if err := tx.Create(&post).Error; err != nil {
			return err
		}

		logEntry := models.ActivityLog{Action: "post_created", PostID: post.ID}
		if err := tx.Create(&logEntry).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Post created successfully"})
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
