package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ukama/ukama/systems/api/api-gateway/pkg"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/db"
	"github.com/ukama/ukama/systems/common/rest"

	log "github.com/sirupsen/logrus"
)

type Client interface {
	GetNetwork(string) (*NetworkInfo, error)
	CreateNetwork(string, string) (*NetworkInfo, error)
}

type clients struct {
	resRepo db.ResourceRepo
	network *networkClient
}

func NewClientsSet(resRepo db.ResourceRepo, endpoints *pkg.HttpEndpoints) Client {
	c := &clients{resRepo: resRepo}

	c.network = NewNetworkClient(endpoints.Network)

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

	res, err := c.resRepo.Get(net.Id)
	if err != nil {
		log.Errorf("inconsistent state. failed to get network resource status: %s",
			err.Error())

		return nil, rest.HttpError{
			HttpCode: http.StatusUnprocessableEntity,
			Message:  "inconsistent state. failed to get network resource status: " + err.Error(),
		}
	}

	// TODO: IsComplete should be added to upstream resource in order to remove all this intermediary
	// resource state management.
	switch res.Status {
	case db.ResourceStatusPending:
		log.Warn("partial content. request is still ongoing")

		return net, rest.HttpError{
			HttpCode: http.StatusPartialContent,
			Message:  "partial content. request is still ongoing",
		}

	case db.ResourceStatusFailed:
		log.Warn("inconsistent state. request has failed")

		return nil, rest.HttpError{
			HttpCode: http.StatusUnprocessableEntity,
			Message:  "inconsistent state. request has failed",
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

	res := &db.Resource{
		Id:     net.Id,
		Status: db.ParseStatus("pending"),
	}

	err = c.resRepo.Add(res)
	if err != nil {
		return nil,
			fmt.Errorf("inconsistent state. failed to update network resource status: %w", err)

	}

	return net, nil
}
