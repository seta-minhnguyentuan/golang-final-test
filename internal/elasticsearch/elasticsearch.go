package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
)

var (
	ES   *elasticsearch.Client
	once sync.Once
)

func InitElasticsearch() *elasticsearch.Client {
	once.Do(func() {
		var err error
		ES, err = elasticsearch.NewDefaultClient()
		if err != nil {
			log.Printf("Warning: failed to init elasticsearch: %v", err)
			log.Printf("Search functionality will be disabled")
			ES = nil
		} else {
			_, err = ES.Info()
			if err != nil {
				log.Printf("Warning: elasticsearch not reachable: %v", err)
				log.Printf("Search functionality will be disabled")
				ES = nil
			} else {
				log.Printf("Elasticsearch connected successfully")
				if err := createPostsIndex(); err != nil {
					log.Printf("Warning: failed to create posts index: %v", err)
				}
			}
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
	res.Body.Close()

	if res.StatusCode == 200 {
		log.Printf("Posts index already exists")
		return nil
	}

	mapping := `{
		"mappings": {
			"properties": {
				"id": {"type": "keyword"},
				"title": {
					"type": "text",
					"analyzer": "standard"
				},
				"content": {
					"type": "text",
					"analyzer": "standard"
				},
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
		json.NewDecoder(createRes.Body).Decode(&response)
		log.Printf("Failed to create posts index: %v", response)
		return nil
	}

	log.Printf("Posts index created successfully")
	return nil
}
