package client

import (
	"errors"
	"net/http"

	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client/rest"
	"github.com/ukama/ukama/systems/common/types"

	log "github.com/sirupsen/logrus"
	crest "github.com/ukama/ukama/systems/common/rest"
)

const (
	failedRequestMsg  = "invalid content. request has failed"
	pendingRequestMsg = "partial content. request is still ongoing"
)

func handleRestErrorStatus(err error) error {
	e := rest.ErrorStatus{}

	if errors.As(err, &e) {
		return crest.HttpError{
			HttpCode: e.StatusCode,
			Message:  err.Error(),
		}
	}

	return err
}

type Network interface {
	GetNetwork(string) (*rest.NetworkInfo, error)
	CreateNetwork(string, string, []string, []string, float64, float64, uint32, bool) (*rest.NetworkInfo, error)
}

type network struct {
	nc rest.NetworkClient
}

func NewNetworkClientSet(ntwk rest.NetworkClient) Network {
	n := &network{
		nc: ntwk,
	}

	return n
}

func (n *network) GetNetwork(id string) (*rest.NetworkInfo, error) {
	net, err := n.nc.Get(id)
	if err != nil {
		return nil, handleRestErrorStatus(err)
	}

	if net.SyncStatus == types.SyncStatusUnknown.String() || net.SyncStatus == types.SyncStatusFailed.String() {
		log.Error(failedRequestMsg)

		return nil, crest.HttpError{
			HttpCode: http.StatusUnprocessableEntity,
			Message:  failedRequestMsg,
		}
	}

	if net.SyncStatus == types.SyncStatusPending.String() {
		log.Warn(pendingRequestMsg)

		return net, crest.HttpError{
			HttpCode: http.StatusPartialContent,
			Message:  pendingRequestMsg,
		}
	}

	return net, nil
}

func (n *network) CreateNetwork(orgName, NetworkName string, allowedCountries,
	allowedNetworks []string, budget, overdraft float64, trafficPolicy uint32,
	paymentLinks bool) (*rest.NetworkInfo, error) {
	net, err := n.nc.Add(rest.AddNetworkRequest{
		OrgName:          orgName,
		NetName:          NetworkName,
		AllowedCountries: allowedCountries,
		AllowedNetworks:  allowedNetworks,
		Budget:           budget,
		Overdraft:        overdraft,
		TrafficPolicy:    trafficPolicy,
		PaymentLinks:     paymentLinks,
	})
	if err != nil {
		return nil, handleRestErrorStatus(err)
	}

	return net, nil
}
