package client

import (
	"github.com/ukama/ukama/systems/api/api-gateway/pkg"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/db"
)

type Client interface {
	GetNetwork(string) (*NetworkInfo, *db.ResourceStatus, error)
	CreateNetwork(string, string) (*NetworkInfo, *db.ResourceStatus, error)
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

func (c *clients) GetNetwork(id string) (*NetworkInfo, *db.ResourceStatus, error) {
	net, err := c.network.Get(id)
	if err != nil {
		return nil, nil, err
	}

	res, err := c.resRepo.Get(net.Id)
	if err != nil {
		return nil, nil, err
	}

	return net, &res.Status, nil
}

func (c *clients) CreateNetwork(orgName, NetworkName string) (*NetworkInfo, *db.ResourceStatus, error) {
	net, err := c.network.Add(AddNetworkRequest{OrgName: orgName, NetName: NetworkName})
	if err != nil {
		return nil, nil, err
	}

	res := &db.Resource{
		Id:     net.Id,
		Status: db.ParseStatus("pending"),
	}

	err = c.resRepo.Add(res)
	if err != nil {
		return nil, nil, err
	}

	return net, &res.Status, nil
}
