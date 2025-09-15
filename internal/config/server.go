package config

import "golang-final-test/pkg/utils"

type ServerConfig struct {
	Port string
}

func LoadServerConfig() ServerConfig {
	return ServerConfig{
		Port: utils.GetEnv("SERVER_PORT", "8080"),
	}
}
