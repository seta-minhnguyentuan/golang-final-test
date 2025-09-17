package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang-final-test/internal/database"
	"golang-final-test/internal/elasticsearch"
	"golang-final-test/internal/models"
	"golang-final-test/internal/redis"
	"golang-final-test/internal/repository"
	"time"

	"github.com/elastic/go-elasticsearch/v9/esapi"
	"gorm.io/gorm"
)

type PostService interface {
	CreatePost(post *models.Post) error
	SearchPostsByTag(tag string) ([]*models.Post, error)
	GetPostByID(id string) (*models.Post, error)
	UpdatePost(post *models.Post) error
	SearchPosts(query string, offset, limit int) ([]*models.Post, int64, error)
}

type postService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) PostService {
	return &postService{repo: repo}
}

func (s *postService) indexPost(post *models.Post) error {
	data, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("failed to marshal post: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      "posts",
		DocumentID: post.ID.String(),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), elasticsearch.ES)
	if err != nil {
		return fmt.Errorf("failed to execute index request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch index error: %s", res.String())
	}

	return nil
}

func (s *postService) CreatePost(post *models.Post) error {
	if err := database.WithTransaction(func(tx *gorm.DB) error {
		repo := repository.NewPostRepository(tx)
		return repo.CreatePost(post)
	}); err != nil {
		return err
	}

	if err := s.indexPost(post); err != nil {
		fmt.Printf("warning: failed to index post to elasticsearch: %v\n", err)
	}

	return nil
}

func (s *postService) SearchPostsByTag(tag string) ([]*models.Post, error) {
	return s.repo.SearchPostsByTag(tag)
}

func (s *postService) GetPostByID(id string) (*models.Post, error) {
	cacheKey := "post:" + id

	val, err := redis.Get(cacheKey)
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
		if setErr := redis.Set(cacheKey, string(jsonData), 5*time.Minute); setErr != nil {
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

		if err := s.indexPost(post); err != nil {
			fmt.Printf("warning: failed to update post in elasticsearch: %v\n", err)
		}

		cacheKey := "post:" + post.ID.String()
		if err := redis.Delete(cacheKey); err != nil {
			fmt.Printf("warning: failed to delete cache for key %s: %v\n", cacheKey, err)
		} else {
			fmt.Println("Cache deleted from Redis:", cacheKey)
		}
		return nil
	})
}

type SearchResponse struct {
	Took int64 `json:"took"`
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			ID     string          `json:"_id"`
			Source json.RawMessage `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (s *postService) SearchPosts(query string, offset, limit int) ([]*models.Post, int64, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"title", "content"},
				"type":   "best_fields",
			},
		},
		"from": offset,
		"size": limit,
	}

	queryBody, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal search query: %w", err)
	}

	req := esapi.SearchRequest{
		Index: []string{"posts"},
		Body:  bytes.NewReader(queryBody),
	}

	res, err := req.Do(context.Background(), elasticsearch.ES)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute search request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("elasticsearch search error: %s", res.String())
	}

	var searchResponse SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
		return nil, 0, fmt.Errorf("failed to decode search response: %w", err)
	}

	posts := make([]*models.Post, len(searchResponse.Hits.Hits))
	for i, hit := range searchResponse.Hits.Hits {
		var post models.Post
		if err := json.Unmarshal(hit.Source, &post); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal post: %w", err)
		}
		posts[i] = &post
	}

	return posts, searchResponse.Hits.Total.Value, nil
}
