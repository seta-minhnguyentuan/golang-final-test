package repository

import (
	"golang-final-test/internal/models"

	"gorm.io/gorm"
)

type PostRepository interface {
	CreatePost(post *models.Post) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) CreatePost(post *models.Post) error {
	return r.db.Create(post).Error
}
