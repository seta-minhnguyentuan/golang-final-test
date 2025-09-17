package elasticsearch

import (
	"log"
	"sync"

	"github.com/elastic/go-elasticsearch/v9"
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
			log.Fatalf("failed to init elasticsearch: %v", err)
		}
	})
	return ES
}
