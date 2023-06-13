package client

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
	nodepb "github.com/ukama/ukama/systems/registry/node/pb/gen"
)

const NodeRegistry = "/v1/nodes/"

type NodeRegistryClient interface {
	GetNode(id string) (*NodeInfo, error)
}

type nodeRegistryClient struct {
	R *rest.RestClient
}

type NodeInfo struct {
	Id      string
	Network string
	Org     string
}

func NewRegistryClient(url string, debug bool) (*nodeRegistryClient, error) {
	f, err := rest.NewRestClient(url, debug)
	if err != nil {
		log.Errorf("Can't conncet to %s url. Error %s", url, err.Error())

		return nil, err
	}

	N := &nodeRegistryClient{
		R: f,
	}

	return N, nil
}

func (c *nodeRegistryClient) GetNode(id string) (*NodeInfo, error) {
	errStatus := &rest.ErrorMessage{}

	nodeInfo := &NodeInfo{}

	rnode := &nodepb.GetNodeResponse{}

	//TODO: Update request api
	resp, err := c.R.C.R().
		SetError(errStatus).
		Get(c.R.URL.String() + NodeRegistry + id)

	if err != nil {
		log.Errorf("Failed to send api request to registry/node. Error %s", err.Error())

		return nil, fmt.Errorf("api request to data plan system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to fetch node info. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)

		return nil, fmt.Errorf("node Info failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &rnode)
	if err != nil {
		log.Tracef("Failed to deserialize data package info. Error message is %s", err.Error())

		return nil, fmt.Errorf("data package info deserailization failure: %+v", err)
	} else {
		nodeInfo.Id = rnode.Node.Node
		nodeInfo.Network = rnode.Node.Network
		nodeInfo.Org = ""
	}

	log.Infof("Node Info: %+v", nodeInfo)

	return nodeInfo, nil
}
