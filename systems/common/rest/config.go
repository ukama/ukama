package rest

import (
	cors "github.com/gin-contrib/cors"
)

// Use this if you don't use github.com/num30/config
func DefaultHTTPConfig() HttpConfig {
	return HttpConfig{
		Port: 8080,
		Cors: cors.Config{
			AllowOrigins: []string{"http://localhost", "https://localhost", "*"},
		},
	}
}

type HttpConfig struct {
	Port int         `default:"8080"`
	Cors cors.Config `default:"{\"AllowOrigins\": [\"http://localhost\", \"https://localhost\", \"*\"]}"`
}
