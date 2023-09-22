package client

import (
	"errors"
	"net/http"

	"github.com/ukama/ukama/systems/common/rest"

	log "github.com/sirupsen/logrus"
)

type Client interface {
	GetNetwork(string) (*NetworkInfo, error)
	CreateNetwork(string, string, []string, []string, bool) (*NetworkInfo, error)

	GetSim(string) (*SimInfo, error)
	ConfigureSim(string, string, string, string, string) (*SimInfo, error)
}

type clients struct {
	network    NetworkClient
	subscriber SubscriberClient
	sim        SimClient
}

func NewClientsSet(network NetworkClient, subscriber SubscriberClient, sim SimClient) Client {
	c := &clients{
		network:    network,
		subscriber: subscriber,
		sim:        sim,
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

func (c *clients) CreateNetwork(orgName, NetworkName string,
	allowedCountries, allowedNetworks []string, paymentLinks bool) (*NetworkInfo, error) {
	net, err := c.network.Add(AddNetworkRequest{
		OrgName:          orgName,
		NetName:          NetworkName,
		AllowedCountries: allowedCountries,
		AllowedNetworks:  allowedNetworks,
		PaymentLinks:     paymentLinks,
	})
	if err != nil {
		return nil, err
	}

	return net, nil
}

func (c *clients) GetSim(id string) (*SimInfo, error) {
	sim, err := c.sim.Get(id)
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

	if !sim.IsSynced {
		log.Warn("partial content. request is still ongoing")

		return sim, rest.HttpError{
			HttpCode: http.StatusPartialContent,
			Message:  "partial content. request is still ongoing",
		}
	}

	return sim, nil
}

func (c *clients) ConfigureSim(subscriberId, networkId, packageId,
	simType, simToken string) (*SimInfo, error) {
	sim, err := c.sim.Add(AddSimRequest{
		SubscriberId: subscriberId,
		NetworkId:    networkId,
		PackageId:    packageId,
		SimType:      simType,
		SimToken:     simToken,
	})
	if err != nil {
		return nil, err
	}

	return sim, nil
}
