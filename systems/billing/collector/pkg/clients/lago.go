package client

import (
	"fmt"

	"github.com/getlago/lago-go-client"
)

type LagoClient struct {
	L *lago.Client
}

func NewLagoClient(APIKey, Host string, Port uint) *LagoClient {
	lagoBaseURL := fmt.Sprintf("http://%s:%d", Host, Port)

	return &LagoClient{
		L: lago.New().SetBaseURL(lagoBaseURL).SetApiKey(APIKey).SetDebug(true),
	}
}
