package rest

import (
	cors "github.com/gin-contrib/cors"
)

func DefaultHTTPConfig() HttpConfig {
	return HttpConfig{
		Port: 8080,
		Cors: cors.Config{
			AllowOrigins: []string{"http://localhost", "https://localhost", "*"},
		},
	}
}

type HttpConfig struct {
	Port int `default:"8080"`
	Cors cors.Config
}
