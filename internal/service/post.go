package service

import (
	"encoding/json"
	"fmt"
	"golang-final-test/internal/database"
	"golang-final-test/internal/models"
	"golang-final-test/internal/repository"
	"golang-final-test/pkg/cache"
	"time"

	"gorm.io/gorm"
)

type PostService interface {
	CreatePost(post *models.Post) error
	SearchPostsByTag(tag string) ([]*models.Post, error)
	GetPostByID(id string) (*models.Post, error)
}

type postService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) PostService {
	return &postService{repo: repo}
}

func (s *postService) CreatePost(post *models.Post) error {
	return database.WithTransaction(func(tx *gorm.DB) error {
		repo := repository.NewPostRepository(tx)
		return repo.CreatePost(post)
	})
}

func (s *postService) SearchPostsByTag(tag string) ([]*models.Post, error) {
	return s.repo.SearchPostsByTag(tag)
}

func (s *postService) GetPostByID(id string) (*models.Post, error) {
	cacheKey := "post:" + id

	val, err := cache.Get(cacheKey)
	if err == nil && val != "" {
		var post models.Post
		if jsonErr := json.Unmarshal([]byte(val), &post); jsonErr == nil {
			return &post, nil
		}
	}

	post, err := s.repo.GetPostByID(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching post: %w", err)
	}

	if jsonData, jsonErr := json.Marshal(post); jsonErr == nil {
		if setErr := cache.Set(cacheKey, string(jsonData), 5*time.Minute); setErr != nil {
			fmt.Printf("warning: failed to set cache for key %s: %v\n", cacheKey, setErr)
		} else {
			fmt.Println("Cache saved to Redis:", cacheKey)
		}
	}

	return post, nil
}

func (s *postService) UpdatePost(post *models.Post) error {
	return database.WithTransaction(func(tx *gorm.DB) error {
		repo := repository.NewPostRepository(tx)
		if err := repo.UpdatePost(post); err != nil {
			return err
		}
		cacheKey := "post:" + post.ID.String()
		if err := cache.Delete(cacheKey); err != nil {
			fmt.Printf("warning: failed to delete cache for key %s: %v\n", cacheKey, err)
		} else {
			fmt.Println("Cache deleted from Redis:", cacheKey)
		}
		return nil
	})
}
