package config

import (
	"golang-final-test/pkg/utils"
)

type ElasticsearchConfig struct {
	Host string
	Port string
}

func LoadElasticsearchConfig() *ElasticsearchConfig {
	return &ElasticsearchConfig{
		Host: utils.GetEnv("ELASTICSEARCH_HOST", "localhost"),
		Port: utils.GetEnv("ELASTICSEARCH_PORT", "9200"),
	}
}
