package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"golang-final-test/internal/config"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
)

var (
	ES   *elasticsearch.Client
	once sync.Once
)

func InitElasticsearch() *elasticsearch.Client {
	once.Do(func() {
		cfg := config.LoadElasticsearchConfig()
		address := fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port)

		var err error
		ES, err = elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{address},
		})
		if err != nil {
			log.Printf("Warning: failed to init elasticsearch client: %v", err)
			ES = nil
			return
		}

		res, err := ES.Info()
		if err != nil {
			log.Printf("Warning: elasticsearch not reachable at %s: %v", address, err)
			ES = nil
			return
		}
		defer res.Body.Close()

		log.Printf("Elasticsearch connected successfully at %s", address)

		if err := createPostsIndex(); err != nil {
			log.Printf("Warning: failed to create posts index: %v", err)
		}
	})
	return ES
}

func createPostsIndex() error {
	req := esapi.IndicesExistsRequest{
		Index: []string{"posts"},
	}

	res, err := req.Do(context.Background(), ES)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		log.Printf("Posts index already exists")
		return nil
	}

	mapping := `{
		"mappings": {
			"properties": {
				"id": {"type": "keyword"},
				"title": {"type": "text", "analyzer": "standard"},
				"content": {"type": "text", "analyzer": "standard"},
				"created_at": {"type": "date"},
				"updated_at": {"type": "date"}
			}
		},
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		}
	}`

	createReq := esapi.IndicesCreateRequest{
		Index: "posts",
		Body:  bytes.NewReader([]byte(mapping)),
	}

	createRes, err := createReq.Do(context.Background(), ES)
	if err != nil {
		return err
	}
	defer createRes.Body.Close()

	if createRes.IsError() {
		var response map[string]interface{}
		_ = json.NewDecoder(createRes.Body).Decode(&response)
		log.Printf("Failed to create posts index: %v", response)
		return nil
	}

	log.Printf("Posts index created successfully")
	return nil
}
