package service

import (
	"golang-final-test/internal/models"
	"golang-final-test/internal/repository"
)

type PostService interface {
	CreatePost(post *models.Post) error
	SearchPostsByTag(tag string) ([]*models.Post, error)
}

type postService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) PostService {
	return &postService{repo: repo}
}

func (s *postService) CreatePost(post *models.Post) error {
	return s.repo.CreatePost(post)
}

func (s *postService) SearchPostsByTag(tag string) ([]*models.Post, error) {
	return s.repo.SearchPostsByTag(tag)
}
