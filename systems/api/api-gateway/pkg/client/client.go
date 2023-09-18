package client

import (
	"errors"
	"net/http"

	"github.com/ukama/ukama/systems/common/rest"

	log "github.com/sirupsen/logrus"
)

type Client interface {
	GetNetwork(string) (*NetworkInfo, error)
	CreateNetwork(string, string) (*NetworkInfo, error)
}

type clients struct {
	network NetworkClient
}

func NewClientsSet(network NetworkClient) Client {
	c := &clients{
		network: network,
	}

	return c
}

func (c *clients) GetNetwork(id string) (*NetworkInfo, error) {
	net, err := c.network.Get(id)
	if err != nil {
		e := ErrorStatus{}

		if errors.As(err, &e) {
			return nil, rest.HttpError{
				HttpCode: e.StatusCode,
				Message:  err.Error(),
			}
		}

		return nil, err
	}

	if !net.IsSynced {
		log.Warn("partial content. request is still ongoing")

		return net, rest.HttpError{
			HttpCode: http.StatusPartialContent,
			Message:  "partial content. request is still ongoing",
		}
	}

	return net, nil
}

func (c *clients) CreateNetwork(orgName, NetworkName string) (*NetworkInfo, error) {
	net, err := c.network.Add(AddNetworkRequest{
		OrgName: orgName,
		NetName: NetworkName})
	if err != nil {
		return nil, err
	}

	return net, nil
}
