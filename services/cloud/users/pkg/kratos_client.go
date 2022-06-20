package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type KratosClient interface {
	GetAccountName(networkOwnerId string) (string, error)
}
type kratosClient struct {
	apiUrl string
}

type Response struct {
	Traits struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"traits"`
}

func NewKratosClient(apiUrl string) *kratosClient {
	return &kratosClient{apiUrl: apiUrl}
}

func (i *kratosClient) GetAccountName(networkOwnerId string) (string, error) {
	if len(networkOwnerId) <= 0 {
		fmt.Println("Missing userId in the request")
		
	}
	resp, err := http.Get(i.apiUrl + networkOwnerId)
	if err != nil {
		return "", errors.Wrap(err,"failed to get a response")

	}

	defer resp.Body.Close()

	dataByte, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", errors.Wrap(err, "failed to decode json response")
	}

	var result Response

	if err := json.Unmarshal(dataByte, &result); err != nil {
		return "", errors.Wrap(err, "failed to marshal userName")

	}

	usr := result.Traits.Name

	return usr, nil
}
