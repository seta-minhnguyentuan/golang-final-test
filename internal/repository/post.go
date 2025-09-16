package repository

import (
	"golang-final-test/internal/models"

	"gorm.io/gorm"
)

type PostRepository interface {
	CreatePost(post *models.Post) error
	SearchPostsByTag(tag string) ([]*models.Post, error)
	GetPostByID(id string) (*models.Post, error)
	UpdatePost(post *models.Post) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) CreatePost(post *models.Post) error {
	if err := r.db.Create(post).Error; err != nil {
		return err
	}
	logEntry := models.ActivityLog{Action: "post_created", PostID: post.ID}
	if err := r.db.Create(&logEntry).Error; err != nil {
		return err
	}
	return nil
}

func (r *postRepository) SearchPostsByTag(tag string) ([]*models.Post, error) {
	var posts []*models.Post
	if err := r.db.Where("? = ANY (tags)", tag).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) GetPostByID(id string) (*models.Post, error) {
	var post models.Post
	if err := r.db.First(&post, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) UpdatePost(post *models.Post) error {
	_, err := r.GetPostByID(post.ID.String())
	if err != nil {
		return err
	}

	if err := r.db.Save(post).Error; err != nil {
		return err
	}
	logEntry := models.ActivityLog{Action: "post_updated", PostID: post.ID}
	if err := r.db.Create(&logEntry).Error; err != nil {
		return err
	}
	return nil
}
